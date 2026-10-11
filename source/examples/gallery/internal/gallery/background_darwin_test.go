package gallery

import (
	"bytes"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

func nativeBackground() bool { return os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" }

// AppKit suspends display-link ticks for hidden windows. Ask its view tree to
// display on the main thread so MyGo consumes input and builds fresh captures.
// This belongs to the native test harness, not the shell's runtime lifecycle.
func startNativeFrames() func() {
	if !nativeBackground() {
		return func() {}
	}
	stop, done := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second / 30)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				mygo.RunOnMain(func() {
					for _, w := range mygo.Windows() {
						if handle := w.NativeHandle(); handle != 0 && !w.IsVisible() {
							displayNativeContent(objc.ID(handle).Send(objc.RegisterName("contentView")))
						}
					}
				})
			}
		}
	}()
	return func() { close(stop); <-done }
}

func displayNativeContent(view objc.ID) {
	if view.Send(objc.RegisterName("isKindOfClass:"), objc.GetClass("MyGoSurfaceView")) != 0 {
		view.Send(objc.RegisterName("setNeedsDisplay:"), true)
		view.Send(objc.RegisterName("display"))
		return
	}
	children := view.Send(objc.RegisterName("subviews"))
	for i, n := uint(0), uint(children.Send(objc.RegisterName("count"))); i < n; i++ {
		displayNativeContent(children.Send(objc.RegisterName("objectAtIndex:"), i))
	}
}

// Check repeated captures, not just creation: CapturePage alone can return an
// old frame when the window is hidden and its display link has stopped.
func TestNativeHiddenShells(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || !nativeBackground() {
		t.Skip("requires background native smoke tests")
	}
	for _, kind := range []string{"sidebar", "plain"} {
		t.Run(kind, func(t *testing.T) {
			var w *mygo.Window
			color := ui.RGB(255, 0, 0)
			shown, focused := false, false
			mygo.RunOnMain(func() {
				content := func(c *ui.Context) { ui.Box(c).Fill().Background(color) }
				if kind == "sidebar" {
					w = sidebar.NewShell(sidebar.ShellOptions{Hidden: true, Content: content}).Window()
				} else {
					w = plain.NewShell(plain.ShellOptions{Hidden: true, Content: content}).Window()
				}
				shown, focused = w.IsVisible(), w.IsFocused()
				w.OnShow(func() { shown = true })
				w.OnFocus(func() { focused = true })
			})
			defer w.Destroy()
			for _, blue := range []bool{false, true} {
				mygo.RunOnMain(func() {
					if blue {
						color = ui.RGB(0, 0, 255)
					}
					w.Invalidate()
				})
				deadline := time.Now().Add(3 * time.Second)
				for {
					data, err := w.CapturePage()
					if err != nil {
						t.Fatal(err)
					}
					img, err := png.Decode(bytes.NewReader(data))
					if err != nil {
						t.Fatal(err)
					}
					b := img.Bounds()
					r, g, bval, a := img.At(b.Dx()*3/4, b.Dy()*2/3).RGBA()
					correct := r == 65535 && g == 0 && bval == 0 && a == 65535
					if blue {
						correct = r == 0 && g == 0 && bval == 65535 && a == 65535
					}
					if correct {
						break
					}
					if time.Now().After(deadline) {
						t.Fatalf("hidden capture did not update: blue=%v pixel=%d,%d,%d,%d bounds=%v", blue, r, g, bval, a, b)
					}
					time.Sleep(30 * time.Millisecond)
				}
				mygo.RunOnMain(func() {
					app := objc.ID(objc.GetClass("NSApplication")).Send(objc.RegisterName("sharedApplication"))
					if shown || focused || w.IsVisible() || w.IsFocused() || objc.Send[bool](app, objc.RegisterName("isActive")) {
						t.Error("background capture showed a window or activated the test app")
					}
				})
			}
		})
	}
}
