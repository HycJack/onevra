package gallery

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGlassPreview(t *testing.T) {
	bitmap, err := loadGalleryImage("../../../../resources")
	if err != nil {
		t.Fatal(err)
	}
	for _, width := range []int{280, 540} {
		for _, dark := range []bool{false, true} {
			s := newDemoState()
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Gap(12).Children(func() { glassPreview(c, s, bitmap) })
			}, width, 720)
			tt.SetDark(dark)
			find := func(label string) ui.Rect {
				t.Helper()
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing %s", label)
				}
				return r
			}
			find("Author avatar")
			for i, label := range []string{"Regular", "Clear", "Tinted"} {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				if s.index != i {
					t.Fatalf("material selection: got %d, want %d", s.index, i)
				}
			}
			for i, label := range []string{"Glass bar", "Soft edge", "Hard edge", "Progressive blur"} {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				if s.glass.edge != i {
					t.Fatalf("edge selection: got %d, want %d", s.glass.edge, i)
				}
				find("Toolbar")
				find("Bottom bar")
			}
			r := find("Draggable glass")
			before := bytes.Clone(tt.Image().Pix)
			x, y := r.X+r.W/2, r.Y+r.H/2
			tt.Press(x, y)
			tt.Move(x+20, y+60)
			tt.Release(x+20, y+60)
			moved := find("Draggable glass")
			if moved.X != r.X+20 || moved.Y != r.Y+60 || bytes.Equal(before, tt.Image().Pix) {
				t.Fatalf("glass did not follow pointer: %v -> %v", r, moved)
			}
			tt.Press(moved.X+10, moved.Y+10)
			tt.Move(2000, 2000)
			tt.Release(2000, 2000)
			tt.SetSize(260, 720)
			pane, canvas := find("Draggable glass"), find("Glass canvas")
			if pane.X < canvas.X || pane.Y < canvas.Y || pane.X+pane.W > canvas.X+canvas.W || pane.Y+pane.H > canvas.Y+canvas.H {
				t.Fatalf("glass escaped canvas: %v in %v", pane, canvas)
			}
			if err := tt.Click("Reset position"); err != nil {
				t.Fatal(err)
			}
			if s.glass.x != 0 || s.glass.y != 0 || s.glass.edge != 3 {
				t.Fatalf("reset did not clear position: %+v", s.glass)
			}
		}
	}
}
