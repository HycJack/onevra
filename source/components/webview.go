package components

import (
	"errors"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// ErrWebViewUnsupported means this platform has no hybrid backend.
var ErrWebViewUnsupported = errors.New("hybrid WebView requires macOS or 64-bit Windows")

type attachment interface {
	place(ui.Rect, bool)
	close()
}

// WebView composites a live system browser below MyGo's GPU-capable Surface.
// Retain it outside builders. All methods except Page run on the UI thread.
// Rectangular scrolling/clipping and native overlays are supported. Rounded
// native clipping, transforms, opacity transitions and unified DOM/native Tab
// traversal are not. Browser dialogs and fullscreen still target a hidden owner.
// Requires the pinned MyGo integration patch prepared by the repository setup.
type WebView struct {
	host, donor                      *mygo.Window
	view                             attachment
	slot                             ui.NativeSlot
	layout                           ui.NativeLayout
	post                             func(func())
	offClosed                        func()
	closed, pending, visible, placed bool
	rect                             ui.Rect
	err                              error
}

func supported() error {
	if (runtime.GOOS == "darwin" || runtime.GOOS == "windows") && (runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64") {
		return nil
	}
	return ErrWebViewUnsupported
}

// NewWebView creates a retained, MyGo-managed browser. Call after App.WhenReady,
// on the UI thread, outside build/layout/paint. No document is loaded yet.
func NewWebView(host *mygo.Window, options mygo.PageOptions) (*WebView, error) {
	if err := supported(); err != nil {
		return nil, err
	}
	if host == nil || host.IsDestroyed() || host.Page() != nil || host.NativeHandle() == 0 {
		return nil, errors.New("WebView requires a live native MyGo window")
	}
	n, err := host.PrepareNativeComposition()
	if err != nil {
		return nil, err
	}
	w := &WebView{host: host, visible: true, post: host.Update}
	w.donor = mygo.NewWindow(mygo.WindowOptions{Title: "WebView owner", Width: 640, Height: 480, Hidden: true, Page: options})
	w.view, err = attachWebView(w, n)
	if err != nil {
		w.donor.Destroy()
		return nil, err
	}
	w.donor.OnClose(func(e *mygo.CloseEvent) { e.PreventDefault() })
	w.offClosed = host.OnClosed(w.Close)
	return w, nil
}

// Page is MyGo's original goroutine-safe Page API. Do not manipulate its hidden
// Window; call WebView.Close instead. Hiding or omitting View retains the page.
func (w *WebView) Page() *mygo.Page { return w.donor.Page() }

// Err reports an asynchronous native attachment failure (not navigation errors).
// Read on the UI thread. Windows creates its WebView2 controller asynchronously.
func (w *WebView) Err() error { return w.err }

// View reserves a GPU-composited slot. Its paint removes earlier UI backgrounds
// in that slot; later overlays blend above the live browser. Do not replace Draw
// or use the instance twice in one frame. Omission unmounts it on the next paint.
func (w *WebView) View(c *ui.Context) *ui.Element {
	w.slot.OnLayout = func(l ui.NativeLayout) {
		if w.closed {
			return
		}
		w.rect = l.Bounds
		if !w.placed || w.layout != l {
			w.layout = l
			w.placed = true
			w.schedule()
		}
	}
	return w.slot.View(c, w.visible && !w.closed && w.err == nil)
}

// SetVisible preserves the document. Showing waits for fresh slot geometry.
func (w *WebView) SetVisible(v bool) {
	if w.closed || v == w.visible {
		return
	}
	w.visible = v
	if v {
		w.placed = false
	}
	w.schedule()
}
func (w *WebView) route(x, y float64, down bool) bool {
	return !w.closed && w.visible && w.err == nil && w.slot.HitTest(x, y, down)
}
func (w *WebView) schedule() {
	if w.closed || w.pending {
		return
	}
	w.pending = true
	w.post(func() {
		w.pending = false
		if !w.closed {
			w.view.place(w.rect, w.visible && w.placed && w.layout.Visible && w.err == nil)
		}
	})
}

// Close releases the browser without closing its host. Host shutdown calls it
// automatically. Queued layouts and asynchronous controller callbacks are inert.
func (w *WebView) Close() {
	if w.closed {
		return
	}
	w.closed = true
	w.slot.Close()
	if w.offClosed != nil {
		w.offClosed()
		w.offClosed = nil
	}
	w.view.close()
	if w.donor != nil {
		w.donor.Destroy()
	}
	if w.host != nil && !w.host.IsDestroyed() {
		w.host.Invalidate()
	}
}
