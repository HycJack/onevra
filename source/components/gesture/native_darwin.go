package gesture

import "github.com/ebitengine/purego/objc"

const (
	eventTypeScrollWheel = 22
	eventTypeMagnify     = 30
	eventMaskMagnify     = 1 << eventTypeMagnify
)

type nsPoint struct{ X, Y float64 }

// currentEvent returns the AppKit event being dispatched, or 0 off the main
// thread, without AppKit or outside dispatch.
func currentEvent() objc.ID {
	app := objc.ID(objc.GetClass("NSApplication"))
	if app == 0 || !mainThread() {
		return 0
	}
	shared := app.Send(objc.RegisterName("sharedApplication"))
	if shared == 0 {
		return 0
	}
	return shared.Send(objc.RegisterName("currentEvent"))
}

// nativeScrollTouch reports whether the scroll being handled came from a
// touch surface: those report gesture and momentum phases, wheels none.
func nativeScrollTouch() (touch, ok bool) {
	ev := currentEvent()
	if ev == 0 || objc.Send[uint](ev, objc.RegisterName("type")) != eventTypeScrollWheel {
		return false, false
	}
	phase := objc.Send[uint](ev, objc.RegisterName("phase"))
	momentum := objc.Send[uint](ev, objc.RegisterName("momentumPhase"))
	return phase != 0 || momentum != 0, true
}

func pointerOf(ev objc.ID) (nativePointer, bool) {
	if ev == 0 {
		return nativePointer{}, false
	}
	window := ev.Send(objc.RegisterName("window"))
	if window == 0 {
		return nativePointer{}, false
	}
	p := objc.Send[nsPoint](ev, objc.RegisterName("locationInWindow"))
	return nativePointer{window: uintptr(window), x: p.X, y: p.Y}, true
}

func currentPointer() (nativePointer, bool) { return pointerOf(currentEvent()) }

// pinchInstalled is only touched on the main thread.
var pinchInstalled bool

// installPinch observes magnify events for the application's lifetime, once
// called on the main thread with AppKit loaded. The monitor hands every
// event back to AppKit unchanged; MyGo's views ignore them.
func installPinch() {
	event := objc.ID(objc.GetClass("NSEvent"))
	if pinchInstalled || event == 0 || !mainThread() {
		return
	}
	pinchInstalled = true
	block := objc.NewBlock(func(_ objc.Block, ev objc.ID) objc.ID {
		if pointer, ok := pointerOf(ev); ok {
			pinch(pointer, objc.Send[float64](ev, objc.RegisterName("magnification")))
		}
		return ev
	})
	event.Send(objc.RegisterName("addLocalMonitorForEventsMatchingMask:handler:"), uint64(eventMaskMagnify), block)
}

func mainThread() bool {
	thread := objc.ID(objc.GetClass("NSThread"))
	return thread != 0 && objc.Send[bool](thread, objc.RegisterName("isMainThread"))
}
