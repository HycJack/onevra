package components

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/theme"
)

func TestIconButtonPlain(t *testing.T) {
	for _, dark := range []bool{false, true} {
		clicks := 0
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			c.Root().Background(ui.Transparent)
			ui.Row(c).Fill().Children(func() {
				if IconButton(c, lucide.Settings(), "Settings", IconButtonPlain).Clicked() {
					clicks++
				}
			})
		}, 200, 80)
		r, ok := tt.Find("Settings")
		if !ok || r.W != 32 || r.H != 32 {
			t.Fatalf("missing square icon button: %+v", r)
		}
		sample := func() uint8 { return tt.Image().RGBAAt(int(r.X+4), int(r.Y+r.H/2)).A }
		tt.Move(190, 70)
		if sample() != 0 {
			t.Fatal("resting button has a face")
		}
		tt.Move(r.X+r.W/2, r.Y+r.H/2)
		hover := sample()
		if hover == 0 || hover == 255 {
			t.Fatalf("hover must be a translucent tint: %d", hover)
		}
		tt.Press(r.X+r.W/2, r.Y+r.H/2)
		if sample() <= hover {
			t.Fatal("press must be a stronger translucent tint")
		}
		tt.Release(r.X+r.W/2, r.Y+r.H/2)
		if clicks != 1 {
			t.Fatalf("clicks=%d", clicks)
		}
	}
}

// TestIconButtonToolbarFace checks that the toolbar style keeps ButtonToolbar's
// raised, shadowed face instead of the plain tinted one.
func TestIconButtonToolbarFace(t *testing.T) {
	for _, dark := range []bool{false, true} {
		face, _, _, _ := toolbarFace(theme.Palette(dark))
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			IconButton(c, lucide.Settings(), "Settings", IconButtonToolbar)
		}, 200, 80)
		r, ok := tt.Find("Settings")
		if !ok || r.W != 32 || r.H != 32 {
			t.Fatalf("missing square icon button: %+v", r)
		}
		// Inside the horizontal padding, clear of the glyph and border.
		p := tt.Image().RGBAAt(int(r.X+4), int(r.Y+r.H/2))
		if p.R != face.R || p.G != face.G || p.B != face.B {
			t.Fatalf("dark=%v: face %v, want %v", dark, p, face)
		}
	}
}

// TestIconButtonToolbarExpanding checks that the expanding toolbar style rests
// as the toolbar square and reveals its label while hovered or keyboard-focused.
func TestIconButtonToolbarExpanding(t *testing.T) {
	for _, dark := range []bool{false, true} {
		clicks := 0
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			ui.Row(c).Fill().AlignItems(ui.Start).Children(func() {
				if IconButton(c, lucide.Upload(), "Upload", IconButtonToolbarExpanding).Clicked() {
					clicks++
				}
			})
		}, 300, 80)
		rest, ok := tt.Find("Upload")
		if !ok || rest.W != 32 || rest.H != 32 {
			t.Fatalf("dark=%v: resting button must be an icon-only square: %+v", dark, rest)
		}
		tt.Move(rest.X+rest.W/2, rest.Y+rest.H/2)
		tt.Advance(time.Second)
		open, _ := tt.Find("Upload")
		if open.W <= rest.W || open.H != rest.H {
			t.Fatalf("dark=%v: hover must reveal the label: %+v", dark, open)
		}
		if err := tt.Click("Upload"); err != nil || clicks != 1 {
			t.Fatalf("dark=%v: clicks=%d, err=%v", dark, clicks, err)
		}
		tt.Move(290, 70)
		tt.Advance(time.Second)
		if r, _ := tt.Find("Upload"); r.W-rest.W > .5 {
			t.Fatalf("dark=%v: leaving must collapse the button: %+v", dark, r)
		}
		tt.Key(0, ui.KeyTab)
		tt.Advance(time.Second)
		if r, _ := tt.Find("Upload"); r.W < open.W-.5 {
			t.Fatalf("dark=%v: keyboard focus must reveal the label: %+v", dark, r)
		}
	}
}
