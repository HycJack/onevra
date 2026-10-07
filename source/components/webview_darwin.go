//go:build darwin

package components

import (
	"errors"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type point struct{ X, Y float64 }
type size struct{ W, H float64 }
type rect struct {
	Origin point
	Size   size
}

func send(id objc.ID, name string, args ...any) objc.ID {
	return id.Send(objc.RegisterName(name), args...)
}

var nativeRect func(objc.ID, objc.SEL) rect
var nativeConvert func(objc.ID, objc.SEL, point, objc.ID) point
var nativeHit func(objc.ID, objc.SEL, point) objc.ID

func bounds(id objc.ID) rect { return nativeRect(id, objc.RegisterName("bounds")) }

type macAttachment struct {
	owner                                           *WebView
	host, donor, root, surface, parent, web, router objc.ID
	closed, shown                                   bool
	lastDown                                        objc.ID
}

var macOnce sync.Once
var macInitErr error
var macRouters = map[objc.ID]*macAttachment{}
var macRouterClass objc.Class
var macMonitor objc.ID
var macMonitorBlock objc.Block
var macCursorHost objc.ID

func initMacBridge() {
	lib, err := purego.Dlopen("/usr/lib/libobjc.A.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		macInitErr = err
		return
	}
	addr, err := purego.Dlsym(lib, "objc_msgSend")
	if err != nil {
		macInitErr = err
		return
	}
	stret := addr
	if runtime.GOARCH == "amd64" {
		stret, err = purego.Dlsym(lib, "objc_msgSend_stret")
		if err != nil {
			macInitErr = err
			return
		}
	}
	purego.RegisterFunc(&nativeRect, stret)
	purego.RegisterFunc(&nativeConvert, addr)
	purego.RegisterFunc(&nativeHit, addr)
	macRouterClass, macInitErr = objc.RegisterClass("OneoHybridRouter", objc.GetClass("NSView"), nil, nil, []objc.MethodDef{{Cmd: objc.RegisterName("hitTest:"), Fn: func(self objc.ID, _ objc.SEL, p point) objc.ID {
		a := macRouters[self]
		if a == nil || a.closed || !a.shown {
			return 0
		}
		p = nativeConvert(a.surface, objc.RegisterName("convertPoint:fromView:"), p, a.parent)
		ev := send(send(objc.ID(objc.GetClass("NSApplication")), "sharedApplication"), "currentEvent")
		kind := send(ev, "type")
		down := kind == 1 || kind == 3 || kind == 25
		if !down {
			a.lastDown = 0
		}
		if !a.owner.route(p.X, p.Y, down && a.lastDown != ev) {
			return 0
		}
		if down {
			a.lastDown = ev
		}
		// WebKit docks its inspector beside the inspected WKWebView, not
		// inside it. Hit-test the whole browser subtree so inspector controls
		// and its split-view divider can receive the original native event.
		p = nativeConvert(a.parent, objc.RegisterName("convertPoint:fromView:"), p, a.surface)
		return nativeHit(a.root, objc.RegisterName("hitTest:"), p)
	}}})
	macMonitorBlock = objc.NewBlock(func(_ objc.Block, ev objc.ID) objc.ID {
		win := send(ev, "window")
		var over *macAttachment
		for _, a := range macRouters {
			if a.closed || !a.shown || a.host != win {
				continue
			}
			if send(ev, "type") == 5 && send(win, "firstResponder") != a.surface {
				send(a.surface, "mouseMoved:", ev)
			}
			p := objc.Send[point](ev, objc.RegisterName("locationInWindow"))
			p = nativeConvert(a.surface, objc.RegisterName("convertPoint:fromView:"), p, 0)
			if a.owner.route(p.X, p.Y, false) {
				over = a
				break
			}
		}
		if over == nil || macCursorHost != 0 && macCursorHost != win {
			restoreMacCursor()
		}
		if over != nil && macCursorHost == 0 {
			send(win, "disableCursorRects")
			send(send(objc.ID(objc.GetClass("NSCursor")), "arrowCursor"), "set")
			macCursorHost = send(win, "retain")
		}
		return ev
	})
}
func restoreMacCursor() {
	if macCursorHost == 0 {
		return
	}
	send(macCursorHost, "enableCursorRects")
	send(macCursorHost, "resetCursorRects")
	send(macCursorHost, "release")
	macCursorHost = 0
}

func attachWebView(w *WebView, n mygo.NativeComposition) (attachment, error) {
	macOnce.Do(initMacBridge)
	if macInitErr != nil {
		return nil, macInitErr
	}
	a := &macAttachment{owner: w, host: objc.ID(w.host.NativeHandle()), donor: objc.ID(w.donor.NativeHandle()), surface: objc.ID(n.View)}
	a.parent, a.root = send(a.host, "contentView"), send(a.donor, "contentView")
	children := send(a.root, "subviews")
	if a.surface == 0 || send(children, "count") != 1 {
		return nil, errors.New("unexpected MyGo native hierarchy")
	}
	a.web = send(children, "objectAtIndex:", uintptr(0))
	if send(a.web, "isKindOfClass:", objc.GetClass("WKWebView")) == 0 {
		return nil, errors.New("missing WKWebView")
	}
	send(a.root, "retain")
	send(a.root, "setHidden:", true)
	placeholder := send(send(objc.ID(objc.GetClass("NSView")), "alloc"), "init")
	send(a.donor, "setContentView:", placeholder)
	send(placeholder, "release")
	send(a.root, "setAutoresizingMask:", uintptr(0))
	send(a.web, "setAutoresizingMask:", uintptr(0))
	send(a.root, "setWantsLayer:", true)
	send(send(a.root, "layer"), "setMasksToBounds:", true)
	send(a.parent, "addSubview:positioned:relativeTo:", a.root, int64(-1), a.surface)
	a.router = send(send(objc.ID(macRouterClass), "alloc"), "init")
	send(a.router, "setHidden:", true)
	send(a.parent, "addSubview:positioned:relativeTo:", a.router, int64(1), a.surface)
	macRouters[a.router] = a
	if macMonitor == 0 {
		macMonitor = send(send(objc.ID(objc.GetClass("NSEvent")), "addLocalMonitorForEventsMatchingMask:handler:", uint64(0x3fffffff), macMonitorBlock), "retain")
	}
	return a, nil
}
func (a *macAttachment) ownsFocus() bool {
	first := send(a.host, "firstResponder")
	return first != 0 && send(first, "isKindOfClass:", objc.GetClass("NSView")) != 0 && send(first, "isDescendantOf:", a.root) != 0
}
func (a *macAttachment) place(r ui.Rect, visible bool) {
	if a.closed || a.owner.host.IsDestroyed() {
		return
	}
	c := a.owner.layout.Clip
	visible = visible && r.W > 0 && r.H > 0 && c.W > 0 && c.H > 0
	a.shown = visible
	if (!visible || !a.owner.layout.Focused) && a.ownsFocus() {
		send(a.host, "makeFirstResponder:", a.surface)
	}
	tx := objc.ID(objc.GetClass("CATransaction"))
	send(tx, "begin")
	send(tx, "setDisableActions:", true)
	if visible {
		h := bounds(a.parent).Size.H
		send(a.root, "setFrame:", rect{point{float64(c.X), h - float64(c.Y+c.H)}, size{float64(c.W), float64(c.H)}})
		send(a.web, "setFrame:", rect{point{float64(r.X - c.X), float64(c.Y + c.H - r.Y - r.H)}, size{float64(r.W), float64(r.H)}})
		send(a.router, "setFrame:", bounds(a.parent))
	}
	send(a.root, "setHidden:", !visible)
	send(a.router, "setHidden:", !visible)
	send(tx, "commit")
	if visible && a.owner.layout.Focused && !a.ownsFocus() {
		send(a.host, "makeFirstResponder:", a.web)
	}
	if !visible && macCursorHost == a.host {
		restoreMacCursor()
	}
}
func (a *macAttachment) close() {
	if a.closed {
		return
	}
	a.place(ui.Rect{}, false)
	a.closed = true
	delete(macRouters, a.router)
	if macCursorHost == a.host {
		restoreMacCursor()
	}
	send(a.router, "removeFromSuperview")
	send(a.router, "release")
	send(a.root, "removeFromSuperview")
	send(a.donor, "setContentView:", a.root)
	send(a.web, "setAutoresizingMask:", uintptr(2|16))
	send(a.root, "release")
	if len(macRouters) == 0 && macMonitor != 0 {
		send(objc.ID(objc.GetClass("NSEvent")), "removeMonitor:", macMonitor)
		send(macMonitor, "release")
		macMonitor = 0
	}
}
