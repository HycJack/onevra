//go:build darwin && arm64

package gallery

import (
	"os"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/components"
)

// Exercise the native pointer path as well as the headless resize tests.
func TestNativeShellResize(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a macOS desktop session")
	}
	var shell *components.Shell
	var main ui.Rect
	mygo.RunOnMain(func() {
		shell = components.NewShell(components.ShellOptions{
			Title: "Sidebar resize review", Width: 1000, Height: 600,
			LeftWidth: 240, LeftMinWidth: 180, LeftMaxWidth: 300,
			Content: func(c *ui.Context) {
				ui.Box(c).Fill().Draw(func(_ *ui.Painter, r ui.Rect) { main = r })
			},
		})
	})
	win := shell.Window()
	defer win.Destroy()
	// AppKit may still be presenting the window after creation returns.
	time.Sleep(400 * time.Millisecond)
	waitWidth := func(want float32) {
		t.Helper()
		deadline := time.Now().Add(3 * time.Second)
		var got float32
		for time.Now().Before(deadline) {
			if _, err := win.CapturePage(); err != nil {
				t.Fatal(err)
			}
			mygo.RunOnMain(func() { got = main.X })
			if got == want {
				return
			}
			time.Sleep(30 * time.Millisecond)
		}
		t.Fatalf("native sidebar width = %v, want %v", got, want)
	}
	send := func(id objc.ID, name string, args ...any) objc.ID {
		return id.Send(objc.RegisterName(name), args...)
	}
	mygo.RunOnMain(func() {
		// Keep physical hover movement from entering the synthetic drag.
		send(objc.ID(win.NativeHandle()), "setAcceptsMouseMovedEvents:", false)
	})
	mouse := func(kind uint64, x, y float64, clicks int64) {
		mygo.RunOnMain(func() {
			host := objc.ID(win.NativeHandle())
			_, height := win.ContentSize()
			event := send(objc.ID(objc.GetClass("NSEvent")), "mouseEventWithType:location:modifierFlags:timestamp:windowNumber:context:eventNumber:clickCount:pressure:",
				kind, struct{ X, Y float64 }{x, float64(height) - y}, uint64(0), float64(0), send(host, "windowNumber"), uintptr(0), int64(1), clicks, float32(1))
			send(host, "sendEvent:", event)
		})
		// Let MyGo process each native event before sending the next one;
		// otherwise queued press/move/release events can share one frame.
		time.Sleep(30 * time.Millisecond)
		if _, err := win.CapturePage(); err != nil {
			t.Fatal(err)
		}
	}
	drag := func(from, to float64) {
		mouse(1, from, 100, 1)
		mouse(6, to, 100, 1)
		mouse(2, to, 100, 1)
	}
	reset := func(x float64) {
		for _, count := range []int64{1, 2} {
			mouse(1, x, 200, count)
			mouse(2, x, 200, count)
		}
	}
	waitWidth(240)
	before := win.Bounds()
	drag(240, 360)
	waitWidth(300)
	reset(300)
	waitWidth(240)
	drag(240, 140)
	waitWidth(180)
	reset(180)
	waitWidth(240)
	if after := win.Bounds(); after != before {
		t.Fatalf("sidebar interaction resized window: %+v -> %+v", before, after)
	}
}
