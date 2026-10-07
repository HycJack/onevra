//go:build darwin && arm64

package components

import (
	"errors"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// This experimental bridge relies on MyGo v0.2.15's AppKit view hierarchy.
// It uses public NSWindow handles, not access to MyGo's internal Go structs.
type point struct{ X, Y float64 }
type size struct{ W, H float64 }
type rect struct {
	Origin point
	Size   size
}

func send(id objc.ID, name string, args ...any) objc.ID {
	return id.Send(objc.RegisterName(name), args...)
}
func bounds(id objc.ID) rect {
	return objc.Send[rect](id, objc.RegisterName("bounds"))
}

type macAttachment struct {
	host, donor, root, surface, parent objc.ID
	closed                             bool
}

func supported() error { return nil }

func attach(host, donor *mygo.Window) (attachment, error) {
	h, d := objc.ID(host.NativeHandle()), objc.ID(donor.NativeHandle())
	parent, root := send(h, "contentView"), send(d, "contentView")
	if parent == 0 || root == 0 {
		return nil, errors.New("missing AppKit content views")
	}
	children := send(parent, "subviews")
	// Shell adds a vibrancy view; other WebView instances add their own roots.
	// Locate the shared native Surface without assuming an only-child layout.
	var surface objc.ID
	for i := uintptr(0); i < uintptr(send(children, "count")); i++ {
		child := send(children, "objectAtIndex:", i)
		if send(child, "isKindOfClass:", objc.GetClass("MyGoSurfaceView")) != 0 {
			if surface != 0 {
				return nil, errors.New("ambiguous MyGo native surface")
			}
			surface = child
		}
	}
	if surface == 0 {
		return nil, errors.New("expected MyGoSurfaceView; review the pinned backend")
	}
	webChildren := send(root, "subviews")
	if send(webChildren, "count") != 1 || send(send(webChildren, "objectAtIndex:", uintptr(0)), "isKindOfClass:", objc.GetClass("WKWebView")) == 0 {
		return nil, errors.New("expected a single MyGo-managed WKWebView")
	}
	// Hold an explicit reference until restored, independent of NSWindow and
	// NSView's ownership. MyGo retains its own references and browser delegates.
	send(root, "retain")
	send(root, "setHidden:", true)
	placeholder := send(send(objc.ID(objc.GetClass("NSView")), "alloc"), "init")
	send(d, "setContentView:", placeholder)
	send(placeholder, "release")
	send(root, "setAutoresizingMask:", uintptr(0))
	send(parent, "addSubview:", root)
	return &macAttachment{host: h, donor: d, root: root, surface: surface, parent: parent}, nil
}

func (a *macAttachment) place(r ui.Rect, visible bool) {
	if a.closed {
		return
	}
	if !visible || r.W <= 0 || r.H <= 0 {
		// Do not leave keyboard focus in an invisible browser.
		first := send(a.host, "firstResponder")
		if first != 0 && send(first, "isKindOfClass:", objc.GetClass("NSView")) != 0 && send(first, "isDescendantOf:", a.root) != 0 {
			send(a.host, "makeFirstResponder:", a.surface)
		}
		send(a.root, "setHidden:", true)
		return
	}
	// MyGo uses top-left DIPs. Its containing NSView is normally unflipped.
	y := float64(r.Y)
	if send(a.parent, "isFlipped") == 0 {
		y = bounds(a.parent).Size.H - y - float64(r.H)
	}
	frame := rect{point{float64(r.X), y}, size{float64(r.W), float64(r.H)}}
	send(a.root, "setFrame:", frame)
	send(a.root, "setHidden:", false)
}

func (a *macAttachment) close() {
	if a.closed {
		return
	}
	a.closed = true
	// No access to host: OnClosed can run after its native window was freed.
	send(a.root, "removeFromSuperview")
	send(a.donor, "setContentView:", a.root)
	send(a.root, "release")
}
