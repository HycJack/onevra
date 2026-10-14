// Package gesture turns scroll and touchpad input into pan and zoom for
// canvases such as graphs, maps and image viewers: a mouse wheel and a
// touchpad pinch zoom, and two fingers dragged on a touchpad pan.
//
// MyGo reports scrolling, not the device behind it, and does not deliver
// touchpad magnification. On macOS a narrow AppKit adapter reads the scroll
// phase of the event being handled, which only touch surfaces report, and
// observes magnify events without consuming them. Elsewhere precise scrolling
// counts as a touchpad, and a pinch arrives as Control-scroll, as Windows
// precision touchpads send it.
package gesture

import (
	"math"
	"weak"

	"github.com/egoist/mygo/ui"
)

// Device is what produced a scroll.
type Device uint8

const (
	// Wheel is a mouse wheel, including smoothed ones that scroll by
	// pixels.
	Wheel Device = iota
	// Touchpad is a touch surface: a touchpad, or a Magic Mouse on macOS.
	Touchpad
)

// DeviceOf reports which device produced an InputScroll. Call it while
// handling the event, before returning from HandleInput.
func DeviceOf(ev ui.InputEvent) Device {
	if touch, ok := nativeScrollTouch(); ok {
		if touch {
			return Touchpad
		}
		return Wheel
	}
	if ev.Precise {
		return Touchpad
	}
	return Wheel
}

// Kind is the kind of a navigation Event.
type Kind uint8

const (
	// Pan moves the content by DX, DY DIPs: positive DX moves it right.
	Pan Kind = iota + 1
	// Zoom scales the content by Scale around X, Y.
	Zoom
)

// Event is a pan or zoom at X, Y, relative to the element's box.
type Event struct {
	Kind         Kind
	X, Y, DX, DY float32
	Scale        float64
	Device       Device
	// Pinch marks a zoom by a touchpad pinch, or Control-scroll on a
	// touchpad.
	Pinch bool
}

const (
	// wheelZoomDistance and pinchZoomDistance are the scroll distances in
	// DIPs that zoom by ×1.5.
	wheelZoomDistance = 120
	pinchZoomDistance = 48
)

// Classify maps a scroll to navigation. A touchpad pans unless Control or
// Command is held, which zooms as a pinch does; a wheel zooms, except a
// sideways wheel, or Shift with a vertical one on macOS, which pans across.
// Call it while handling the event; ok is false for other input.
func Classify(ev ui.InputEvent) (Event, bool) {
	if ev.Kind != ui.InputScroll {
		return Event{}, false
	}
	device := DeviceOf(ev)
	out := Event{X: ev.X, Y: ev.Y, Device: device}
	modified := ev.Mods&(ui.Ctrl|ui.Super) != 0
	switch {
	case device == Touchpad && !modified:
		out.Kind, out.DX, out.DY = Pan, -ev.DX, -ev.DY
	case device == Wheel && !modified && ev.DY == 0:
		out.Kind, out.DX = Pan, -ev.DX
	default:
		distance := float64(wheelZoomDistance)
		if device == Touchpad {
			distance, out.Pinch = pinchZoomDistance, true
		}
		out.Kind, out.Scale = Zoom, math.Exp(-float64(ev.DY)*math.Log(1.5)/distance)
	}
	return out, true
}

// Navigation delivers an element's pan and zoom gestures, including
// touchpad pinches that MyGo does not report. Keep one per element; its
// zero value is ready to use.
type Navigation struct {
	fn       func(Event)
	services ui.Services
	// Where the element last saw the pointer, in its own box and in the
	// native window, so a pinch there can be placed in the element.
	x, y    float32
	pointer nativePointer
}

// Handle has el take its input through fn, as HandleInput does, after
// sending scrolling and pinches to navigate. Navigation events reach
// navigate and are taken; other input goes to other, which may be nil.
// Call it on each build.
func (n *Navigation) Handle(c *ui.Context, el ui.Element, navigate func(Event), other func(ui.InputEvent) bool) ui.Element {
	n.fn, n.services = navigate, c.Services()
	installPinch()
	return el.HandleInput(func(ev ui.InputEvent) bool {
		switch ev.Kind {
		case ui.InputPointerMove, ui.InputPointerDown, ui.InputScroll:
			n.track(ev.X, ev.Y)
		}
		if g, ok := Classify(ev); ok && n.fn != nil {
			n.fn(g)
			return true
		}
		if other != nil {
			return other(ev)
		}
		return false
	})
}

// track remembers where the pointer is, making n the target of the next
// pinch made there.
func (n *Navigation) track(x, y float32) {
	pointer, ok := currentPointer()
	if !ok {
		return
	}
	n.x, n.y, n.pointer = x, y, pointer
	pinchTarget = weak.Make(n)
}

// pinchTarget is the navigation that last saw the pointer.
var pinchTarget weak.Pointer[Navigation]

// pinch delivers a magnification made at pointer to the navigation that
// last saw the pointer there; a pointer that has since moved elsewhere was
// seen by another element, or none.
func pinch(pointer nativePointer, magnification float64) bool {
	n := pinchTarget.Value()
	if n == nil || n.fn == nil || !n.pointer.near(pointer) || magnification <= -1 {
		return false
	}
	n.fn(Event{Kind: Zoom, X: n.x, Y: n.y, Scale: 1 + magnification, Device: Touchpad, Pinch: true})
	n.services.Invalidate()
	return true
}

// nativePointer is a pointer location in a native window's coordinates.
type nativePointer struct {
	window uintptr
	x, y   float64
}

func (p nativePointer) near(o nativePointer) bool {
	return p.window != 0 && p.window == o.window && math.Abs(p.x-o.x) <= 1 && math.Abs(p.y-o.y) <= 1
}
