package controls

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestSliderKeyboard(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		value := 50.0
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			Slider(c, &value, 0, 100).Label("volume").AutoFocus()
		}, 300, 100)
		tt.Key(0, ui.KeyRight)
		if value != 51 {
			t.Errorf("dark=%v: Right set %v, want 51", dark, value)
		}
		tt.Key(0, ui.KeyLeft)
		tt.Key(0, ui.KeyLeft)
		if value != 49 {
			t.Errorf("dark=%v: Left twice set %v, want 49", dark, value)
		}
		tt.Key(0, ui.KeyHome)
		if value != 0 {
			t.Errorf("dark=%v: Home set %v, want 0", dark, value)
		}
		tt.Key(0, ui.KeyEnd)
		if value != 100 {
			t.Errorf("dark=%v: End set %v, want 100", dark, value)
		}
	}
}

func TestSliderDrag(t *testing.T) {
	value := 0.0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Slider(c, &value, 0, 100).Label("volume")
	}, 300, 100)
	r, ok := tt.Find("volume")
	if !ok {
		t.Fatal("slider not found")
	}
	// Drag from the left end to the right end of the slider.
	tt.Press(r.X+4, r.Y+r.H/2)
	tt.Move(r.X+r.W-4, r.Y+r.H/2)
	tt.Release(r.X+r.W-4, r.Y+r.H/2)
	if value < 90 {
		t.Errorf("drag set %v, want near 100", value)
	}
}
