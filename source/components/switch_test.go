package components

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestSwitchToggle(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := false
		if dark {
			mode = true
		}
		on := false
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(mode))
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
		c.SetTheme(theme.LightPalette())
		Switch(c, &on, "Wi-Fi")
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !on {
		t.Error("Space did not toggle the focused switch")
	}
}

func TestSwitchHiddenLabel(t *testing.T) {
	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		Switch(c, &on, "Smart notifications", SwitchOptions{HideLabel: true})
	}, 300, 100)
	if len(tt.Texts()) != 1 {
		t.Fatalf("expected only accessible label, got %v", tt.Texts())
	}
	if err := tt.Click("Smart notifications"); err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("named switch did not toggle")
	}
	tt.Key(0, ui.KeySpace)
	if on {
		t.Fatal("keyboard did not toggle hidden-label switch")
	}
}
