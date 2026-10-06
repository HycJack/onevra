package gallery

import (
	"os"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/components"
)

// Sample every painted layout, not only the settled screenshots. This checks
// native viewport synchronization; it cannot observe OS compositor stretching.
func TestNativeInspectorFrames(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	type sample struct {
		left, main       ui.Rect
		viewport, native float32
	}
	var shell *components.Shell
	var left ui.Rect
	var latest sample
	var samples []sample
	recording := false
	mygo.RunOnMain(func() {
		shell = components.NewShell(components.ShellOptions{
			Title: "Inspector frame review", Width: 820, Height: 600,
			MinWidth: 360, MinHeight: 480,
			Left: func(c *ui.Context) {
				ui.Box(c).Fill().Draw(func(_ *ui.Painter, r ui.Rect) { left = r })
			},
			Content: func(c *ui.Context) {
				viewport, _ := c.Size()
				ui.Box(c).Fill().Draw(func(_ *ui.Painter, r ui.Rect) {
					if shell == nil {
						return
					}
					width, _ := shell.Window().ContentSize()
					latest = sample{left, r, viewport, float32(width)}
					if recording {
						samples = append(samples, latest)
					}
				})
			},
			Right: func(c *ui.Context) { ui.Text(c, "Inspector") },
		})
	})
	defer shell.Window().Destroy()
	for _, width := range []int{820, 360} {
		mygo.RunOnMain(func() { shell.Window().SetSize(width, 600) })
		time.Sleep(350 * time.Millisecond)
		var before sample
		mygo.RunOnMain(func() { before = latest })
		for _, open := range []bool{true, false} {
			mygo.RunOnMain(func() {
				samples = nil
				recording = true
				shell.SetRightExpanded(open)
			})
			time.Sleep(450 * time.Millisecond)
			mygo.RunOnMain(func() {
				recording = false
				if len(samples) < 2 {
					t.Errorf("width %d open=%v: too few transition frames (%d)", width, open, len(samples))
				}
				for i, frame := range samples {
					if frame.viewport != frame.native || frame.left != before.left || frame.main != before.main {
						t.Errorf("width %d open=%v frame %d: unstable layout %+v; baseline %+v", width, open, i, frame, before)
					}
				}
				want := width
				if open {
					want += 280
				}
				if got, _ := shell.Window().Size(); got != want {
					t.Errorf("width %d open=%v: final width %d, want %d", width, open, got, want)
				}
				t.Logf("width %d open=%v: sampled %d painted frames", width, open, len(samples))
			})
		}
	}
}
