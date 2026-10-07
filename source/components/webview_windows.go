//go:build windows && (amd64 || arm64)

package components

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// A DirectComposition-backed MyGo Surface sits above clipped WebView2 child
// controllers. This keeps WebView2's native input/IME rather than synthesizing
// DOM events or copying browser pixels into a MyGo bitmap.
var winUser = syscall.NewLazyDLL("user32.dll")
var winComctl = syscall.NewLazyDLL("comctl32.dll")
var winCreate = winUser.NewProc("CreateWindowExW")
var winDestroy = winUser.NewProc("DestroyWindow")
var winShow = winUser.NewProc("ShowWindow")
var winPos = winUser.NewProc("SetWindowPos")
var winDPI = winUser.NewProc("GetDpiForWindow")
var winFocus = winUser.NewProc("SetFocus")
var winGetFocus = winUser.NewProc("GetFocus")
var winIsChild = winUser.NewProc("IsChild")
var winIsWindow = winUser.NewProc("IsWindow")
var winScreenToClient = winUser.NewProc("ScreenToClient")
var winForeground = winUser.NewProc("GetForegroundWindow")
var winHook = winUser.NewProc("SetWindowsHookExW")
var winUnhook = winUser.NewProc("UnhookWindowsHookEx")
var winNextHook = winUser.NewProc("CallNextHookEx")
var winSend = winUser.NewProc("SendMessageW")
var winSubclass = winComctl.NewProc("SetWindowSubclass")
var winUnsubclass = winComctl.NewProc("RemoveWindowSubclass")
var winDefSubclass = winComctl.NewProc("DefSubclassProc")

type winPoint struct{ X, Y int32 }
type winRect struct{ Left, Top, Right, Bottom int32 }
type windowsAttachment struct {
	owner                           *WebView
	host, surface, clip, controller uintptr
	shown, closed                   bool
}

var winAttachments = map[uintptr][]*windowsAttachment{}
var winSubclassProc, winMouseProc, winMouseHook uintptr

func winPtr(p uintptr) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Pointer(&p)) }

//go:uintptrescapes
func webCOM(obj uintptr, index int, args ...uintptr) uintptr {
	vt := *(*uintptr)(winPtr(obj))
	fn := *(*uintptr)(winPtr(vt + uintptr(index)*unsafe.Sizeof(uintptr(0))))
	all := append([]uintptr{obj}, args...)
	r, _, _ := syscall.SyscallN(fn, all...)
	return r
}
func comError(hr uintptr) error {
	if int32(hr) < 0 {
		return fmt.Errorf("WebView2 HRESULT %#x", uint32(hr))
	}
	return nil
}
func (a *windowsAttachment) point(p winPoint) (float64, float64) {
	winScreenToClient.Call(a.surface, uintptr(unsafe.Pointer(&p)))
	dpi, _, _ := winDPI.Call(a.host)
	if dpi == 0 {
		dpi = 96
	}
	return float64(p.X) * 96 / float64(dpi), float64(p.Y) * 96 / float64(dpi)
}
func initWindowsBridge() {
	if winSubclassProc != 0 {
		return
	}
	winSubclassProc = syscall.NewCallback(func(hwnd uintptr, msg uint32, wp, lp, id, data uintptr) uintptr {
		if msg == 0x84 { // WM_NCHITTEST; leave actual press observation to the mouse hook.
			p := winPoint{int32(int16(lp & 0xffff)), int32(int16((lp >> 16) & 0xffff))}
			for _, a := range winAttachments[hwnd] {
				x, y := a.point(p)
				if a.shown && a.owner.route(x, y, false) {
					return ^uintptr(0)
				}
			}
		}
		r, _, _ := winDefSubclass.Call(hwnd, uintptr(msg), wp, lp)
		return r
	})
	// One process-wide observer, removed with the last attachment. It never
	// consumes input: OS hit testing still delivers original events to WebView2.
	winMouseProc = syscall.NewCallback(func(code int32, wp, lp uintptr) uintptr {
		if code >= 0 {
			host, _, _ := winForeground.Call()
			p := *(*winPoint)(winPtr(lp)) // MSLLHOOKSTRUCT begins with POINT.
			down := wp == 0x201 || wp == 0x204 || wp == 0x207
			for _, list := range winAttachments {
				for _, a := range list {
					if a.host != host || !a.shown || a.closed {
						continue
					}
					x, y := a.point(p)
					if a.owner.route(x, y, down) && wp == 0x200 {
						dpi, _, _ := winDPI.Call(host)
						packed := uintptr(uint16(int32(x*float64(dpi)/96))) | uintptr(uint16(int32(y*float64(dpi)/96)))<<16
						winSend.Call(a.surface, 0x200, 0, packed)
					}
				}
			}
		}
		r, _, _ := winNextHook.Call(0, uintptr(code), wp, lp)
		return r
	})
}
func attachWebView(w *WebView, n mygo.NativeComposition) (attachment, error) {
	if n.HWND == 0 {
		return nil, errors.New("missing MyGo composition HWND")
	}
	initWindowsBridge()
	a := &windowsAttachment{owner: w, host: w.host.NativeHandle(), surface: n.HWND}
	class, _ := syscall.UTF16PtrFromString("STATIC")
	a.clip, _, _ = winCreate.Call(0, uintptr(unsafe.Pointer(class)), 0, 0x40000000|0x04000000|0x02000000, 0, 0, 0, 0, a.host, 0, 0, 0)
	if a.clip == 0 {
		return nil, errors.New("cannot create WebView2 clipping HWND")
	}
	if len(winAttachments[a.surface]) == 0 {
		if ok, _, _ := winSubclass.Call(a.surface, winSubclassProc, 1, 0); ok == 0 {
			winDestroy.Call(a.clip)
			return nil, errors.New("cannot install native view routing")
		}
	}
	if winMouseHook == 0 {
		winMouseHook, _, _ = winHook.Call(14, winMouseProc, 0, 0) // WH_MOUSE_LL, callback runs on this UI thread.
		if winMouseHook == 0 {
			winUnsubclass.Call(a.surface, winSubclassProc, 1)
			winDestroy.Call(a.clip)
			return nil, errors.New("cannot observe browser pointer input")
		}
	}
	winAttachments[a.surface] = append(winAttachments[a.surface], a)
	w.donor.Page().WithNativeBrowser(func(controller uintptr, err error) {
		if a.closed {
			return
		}
		if err == nil {
			err = comError(webCOM(controller, 22, a.clip))
		} // put_ParentWindow
		if err != nil {
			w.err = err
			w.host.Invalidate()
			return
		}
		a.controller = controller
		w.schedule()
	})
	return a, nil
}
func (a *windowsAttachment) place(r ui.Rect, visible bool) {
	if a.closed || a.controller == 0 || a.owner.host.IsDestroyed() {
		return
	}
	c := a.owner.layout.Clip
	a.shown = visible && r.W > 0 && r.H > 0 && c.W > 0 && c.H > 0
	focus, _, _ := winGetFocus.Call()
	child, _, _ := winIsChild.Call(a.clip, focus)
	owns := focus == a.clip || child != 0
	if (!a.shown || !a.owner.layout.Focused) && owns {
		winFocus.Call(a.surface)
	}
	if !a.shown {
		webCOM(a.controller, 4, 0)
		winShow.Call(a.clip, 0)
		return
	}
	dpi, _, _ := winDPI.Call(a.host)
	if dpi == 0 {
		dpi = 96
	}
	px := func(x float32) int32 { return int32(math.Round(float64(x) * float64(dpi) / 96)) }
	x, y := px(c.X), px(c.Y)
	winPos.Call(a.clip, 0, uintptr(x), uintptr(y), uintptr(px(c.X+c.W)-x), uintptr(px(c.Y+c.H)-y), 0x54) // no activate, no z-order, show
	b := winRect{px(r.X) - x, px(r.Y) - y, px(r.X+r.W) - x, px(r.Y+r.H) - y}
	var hr uintptr
	if runtime.GOARCH == "arm64" {
		hr = webCOM(a.controller, 6, uintptr(uint64(uint32(b.Left))|uint64(uint32(b.Top))<<32), uintptr(uint64(uint32(b.Right))|uint64(uint32(b.Bottom))<<32))
	} else {
		hr = webCOM(a.controller, 6, uintptr(unsafe.Pointer(&b)))
	}
	if err := comError(hr); err != nil {
		a.owner.err = err
		a.shown = false
		webCOM(a.controller, 4, 0)
		a.owner.host.Invalidate()
		return
	}
	webCOM(a.controller, 23) // NotifyParentWindowPositionChanged
	webCOM(a.controller, 4, 1)
	// Insert the browser immediately behind the UI, preserving caption controls
	// already above the UI instead of raising the Surface over those controls.
	winPos.Call(a.clip, a.surface, 0, 0, 0, 0, 0x13)
	if a.owner.layout.Focused && !owns {
		webCOM(a.controller, 12, 0)
	}
}
func (a *windowsAttachment) close() {
	if a.closed {
		return
	}
	a.place(ui.Rect{}, false)
	a.closed = true
	if a.controller != 0 {
		webCOM(a.controller, 4, 0)
		webCOM(a.controller, 22, a.owner.donor.NativeHandle())
	}
	list := winAttachments[a.surface]
	for i, b := range list {
		if b == a {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	if len(list) == 0 {
		delete(winAttachments, a.surface)
		if ok, _, _ := winIsWindow.Call(a.surface); ok != 0 {
			winUnsubclass.Call(a.surface, winSubclassProc, 1)
		}
	} else {
		winAttachments[a.surface] = list
	}
	if ok, _, _ := winIsWindow.Call(a.clip); ok != 0 {
		winDestroy.Call(a.clip)
	}
	if len(winAttachments) == 0 && winMouseHook != 0 {
		winUnhook.Call(winMouseHook)
		winMouseHook = 0
	}
}
