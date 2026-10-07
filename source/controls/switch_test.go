package controls

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestSwitchToggle(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		on := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			Switch(c, &on, "Wi-Fi")
		}, 300, 100)
		// Find returns the switch track first: it is built and labeled
		// before the text. Click its center.
		r, ok := tt.Find("Wi-Fi")
		if !ok {
			t.Fatal("switch not found")
		}
		tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		if !on {
			t.Errorf("dark=%v: click did not turn on", dark)
		}
		tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		if on {
			t.Errorf("dark=%v: click did not turn off", dark)
		}
	}
}

func TestSwitchKeyboard(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Switch(c, &on, "Wi-Fi")
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !on {
		t.Error("Space did not toggle the focused switch")
	}
}
