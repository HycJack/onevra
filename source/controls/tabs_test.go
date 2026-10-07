package controls

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestTabsClick(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		selected := 0
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			Tabs(c, &selected, "General", "Appearance", "Advanced")
		}, 400, 100)
		if err := tt.Click("Appearance"); err != nil {
			t.Fatal(err)
		}
		if selected != 1 {
			t.Errorf("dark=%v: selected=%d, want 1", dark, selected)
		}
		if err := tt.Click("Advanced"); err != nil {
			t.Fatal(err)
		}
		if selected != 2 {
			t.Errorf("dark=%v: selected=%d, want 2", dark, selected)
		}
	}
}

func TestTabsKeyboard(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Tabs(c, &selected, "General", "Appearance", "Advanced")
	}, 400, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if selected != 1 {
		t.Errorf("Right chose %d, want 1", selected)
	}
	tt.Key(0, ui.KeyLeft)
	if selected != 0 {
		t.Errorf("Left chose %d, want 0", selected)
	}
}
