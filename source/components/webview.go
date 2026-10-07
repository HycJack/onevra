package components

import (
	"errors"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// ErrWebViewUnsupported means this platform has no embedded WebView backend.
var ErrWebViewUnsupported = errors.New("embedded WebView currently requires macOS ARM64")

// attachment moves MyGo-owned native views. MyGo owns browser rendering and
// input; this bridge does not introduce another renderer or event loop.
type attachment interface {
	place(ui.Rect, bool)
	close()
}

// WebView is an experimental system WebView embedded in a native MyGo window.
// Create one per layout slot, outside a view builder, and reuse it across frames.
// All its methods except Page must run on the UI thread. Close is automatic when
// the host closes, or can be called earlier to release browser resources.
//
// Only macOS ARM64 is implemented. The system view sits above MyGo's Surface:
// ancestor clipping/scrolling, arbitrary transforms, native overlays over the
// page and unified Tab navigation are not supported. Use a non-scrolling,
// rectangular slot. Window-bound browser dialogs, menus, popups and fullscreen
// remain associated with the hidden owner and are outside this component's scope.
type WebView struct {
	host, donor                      *mygo.Window
	view                             attachment
	post                             func(func())
	offClosed                        func()
	closed, pending, visible, placed bool
	rect                             ui.Rect
}

// NewWebView attaches a hidden, MyGo-managed browser to host, which must be a
// live native-content window. Call after App.WhenReady, on the UI thread, never
// during build/layout/paint. options configures MyGo's original Page; no content
// is loaded until the caller uses Page().LoadHTML or Page().LoadURL.
//
// Multiple instances can occupy distinct slots in one window. Unsupported
// platforms return ErrWebViewUnsupported without creating another window.
func NewWebView(host *mygo.Window, options mygo.PageOptions) (*WebView, error) {
	if err := supported(); err != nil {
		return nil, err
	}
	if host == nil || host.IsDestroyed() || host.Page() != nil || host.NativeHandle() == 0 {
		return nil, errors.New("WebView requires a live native MyGo window")
	}
	w := &WebView{host: host, visible: true, post: host.Update}
	w.donor = mygo.NewWindow(mygo.WindowOptions{Title: "WebView owner", Width: 640, Height: 480, Hidden: true, Page: options})
	view, err := attach(host, w.donor)
	if err != nil {
		w.donor.Destroy()
		return nil, err
	}
	w.view = view
	// window.close() must not destroy the owner while its view is attached.
	w.donor.OnClose(func(e *mygo.CloseEvent) { e.PreventDefault() })
	w.offClosed = host.OnClosed(w.Close)
	return w, nil
}

// Page returns MyGo's original page for navigation, JavaScript and load events.
// It is safe to use its page operations from any goroutine as documented by
// MyGo. Do not show, destroy or otherwise manipulate Page().Window(): it is a
// hidden lifecycle owner, not the visible host. Use WebView.Close instead.
func (w *WebView) Page() *mygo.Page { return w.donor.Page() }

// View reserves a slot and synchronizes its painted bounds with the system
// view. Style the returned element with Fill, Size, etc. Do not replace its Draw
// callback or use it twice in one frame. Ancestors must not scroll or clip the
// slot itself; the browser does not inherit their clipping.
//
// SetVisible(false) before omitting View on navigation, and SetVisible(true)
// when returning. Omission alone cannot unmount a system view. A hidden view
// retains its document; Close releases it. After Close this is only an empty slot.
func (w *WebView) View(c *ui.Context) *ui.Element {
	return ui.Box(c).Draw(func(_ *ui.Painter, r ui.Rect) {
		if !w.placed || r != w.rect {
			w.rect, w.placed = r, true
			w.schedule()
		}
	})
}

// SetVisible hides or restores the browser, preserving its document. Hiding
// returns keyboard focus to the native Surface if the browser owned it. Showing
// waits for a new View paint so a returning page never uses stale slot geometry.
func (w *WebView) SetVisible(visible bool) {
	if w.closed || visible == w.visible {
		return
	}
	w.visible = visible
	if visible {
		w.placed = false
	}
	w.schedule()
}

func (w *WebView) schedule() {
	if w.closed || w.pending {
		return
	}
	w.pending = true
	// AppKit hierarchy changes can reenter layout. Coalesce them in MyGo's
	// deferred UI queue, never mutate native frames inside a paint transaction.
	w.post(func() {
		w.pending = false
		if !w.closed {
			w.view.place(w.rect, w.visible && w.placed)
		}
	})
}

// Close releases the browser without closing the host. It is idempotent and
// also runs when the host closes. Queued placements become harmless no-ops.
func (w *WebView) Close() {
	if w.closed {
		return
	}
	w.closed = true
	if w.offClosed != nil {
		w.offClosed()
		w.offClosed = nil
	}
	if w.host != nil && !w.host.IsDestroyed() {
		w.view.place(w.rect, false) // Return focus before detaching a live view.
	}
	w.view.close()
	if w.donor != nil {
		w.donor.Destroy()
	}
}
