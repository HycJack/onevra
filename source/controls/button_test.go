package controls

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestButtonVariantsClick(t *testing.T) {
	variants := []ButtonVariant{ButtonPrimary, ButtonDefault, ButtonOutline, ButtonGhost, ButtonLink}
	sizes := []ButtonSize{ButtonSmall, ButtonMedium, ButtonLarge}
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		for _, variant := range variants {
			for _, size := range sizes {
				clicked := false
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c, mode)
					if Button(c, "OK", variant, size).Clicked() {
						clicked = true
					}
				}, 300, 100)
				if err := tt.Click("OK"); err != nil {
					t.Fatal(err)
				}
				if !clicked {
					t.Errorf("dark=%v variant=%d size=%d: click not reported", dark, variant, size)
				}
			}
		}
	}
}

func TestButtonKeyboard(t *testing.T) {
	clicked := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		if Button(c, "Save", ButtonPrimary, ButtonMedium).Clicked() {
			clicked = true
		}
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !clicked {
		t.Error("Space did not click the focused button")
	}
}

func TestButtonDisabledDims(t *testing.T) {
	sample := func(disabled bool) (r, g, b uint8) {
		var rect ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, theme.Light)
			ui.Column(c).Disabled(disabled).Fill().Children(func() {
				Button(c, "OK", ButtonPrimary, ButtonMedium)
			})
		}, 300, 100)
		var ok bool
		rect, ok = tt.Find("OK")
		if !ok {
			t.Fatal("button not found")
		}
		// Inside the left padding, clear of the label glyphs.
		p := tt.Image().RGBAAt(int(rect.X+6), int(rect.Y+rect.H/2))
		return p.R, p.G, p.B
	}
	er, eg, eb := sample(false)
	dr, dg, db := sample(true)
	if er == dr && eg == dg && eb == db {
		t.Fatal("disabled button renders identically to enabled")
	}
	// Dimming blends the accent toward the light background.
	if dr <= er || dg <= eg || db <= eb {
		t.Errorf("disabled pixel (%d,%d,%d) not lighter than enabled (%d,%d,%d)", dr, dg, db, er, eg, eb)
	}
}

func TestButtonLabel(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Button(c, "Delete", ButtonDefault, ButtonMedium)
	}, 300, 100)
	if !tt.HasText("Delete") {
		t.Error("button label not shown")
	}
}
