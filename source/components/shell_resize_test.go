package components

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestShellDeferredResizeKeepsSurfacesStable(t *testing.T) {
	for _, width := range []int{360, 820} {
		s, w := fixture()
		w.bounds.Width = width
		var pending []func()
		s.post = func(fn func()) { pending = append(pending, fn) }
		var left, main ui.Rect
		inFrame := false
		s.options.Left = func(c *ui.Context) {
			ui.Box(c).Fill().Draw(func(_ *ui.Painter, r ui.Rect) { left = r })
		}
		s.options.Content = func(c *ui.Context) {
			viewport, _ := c.Size()
			ui.Box(c).Fill().Draw(func(_ *ui.Painter, r ui.Rect) {
				main = r
				if viewport != float32(w.bounds.Width) {
					t.Errorf("painted stale viewport %v into native width %d", viewport, w.bounds.Width)
				}
			})
		}
		tt := ui.NewTester(func(c *ui.Context) {
			inFrame = true
			defer func() { inFrame = false }()
			s.view(c)
		}, width, 600)
		leftBefore, mainBefore := left, main
		w.onResize = func() {
			if inFrame {
				t.Fatal("native resize during view construction")
			}
			// AppKit may synchronously draw before SetBounds returns.
			tt.SetSize(w.bounds.Width, 600)
			if left != leftBefore || main != mainBefore {
				t.Fatalf("width %d: surfaces moved during native redraw: left %+v -> %+v, main %+v -> %+v", width, leftBefore, left, mainBefore, main)
			}
		}
		flush := func() {
			jobs := pending
			pending = nil
			if len(jobs) > 1 {
				t.Fatal("resize requests were not coalesced")
			}
			for _, fn := range jobs {
				fn()
			}
			if left != leftBefore || main != mainBefore {
				t.Fatal("surfaces moved between acknowledged resizes", left, main)
			}
		}
		for _, open := range []bool{true, false} {
			s.SetRightExpanded(open)
			tt.Frame()
			for range 12 {
				time.Sleep(25 * time.Millisecond)
				tt.Frame()
				flush()
			}
			want := width
			if open {
				want += 280
			}
			if w.bounds.Width != want {
				t.Fatalf("width %d: transition ended at %d, want %d", width, w.bounds.Width, want)
			}
		}
		if w.changes < 4 {
			t.Fatal("inspector no longer animates")
		}
	}
}

func TestShellQueuedResizeUsesLatestTarget(t *testing.T) {
	s, w := fixture()
	var pending []func()
	s.post = func(fn func()) { pending = append(pending, fn) }
	s.requestWidth(100)
	s.requestWidth(200)
	s.requestWidth(0)
	if len(pending) != 1 || w.changes != 0 {
		t.Fatal("requests were not deferred and coalesced")
	}
	pending[0]()
	if w.changes != 0 || s.resizeQueued {
		t.Fatal("cancelled request resized the window")
	}
	pending = nil
	s.requestWidth(280)
	w.maximized = true // State can change between scheduling and dispatch.
	pending[0]()
	if w.changes != 0 || s.resizeQueued {
		t.Fatal("queued request resized a maximized window")
	}
	pending = nil
	w.maximized = false
	s.requestWidth(280)
	pending[0]()
	if w.bounds.Width != 1100 || s.extension != 280 {
		t.Fatal("normal-window resize did not resume", w.bounds)
	}
}
