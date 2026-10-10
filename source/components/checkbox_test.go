package components

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestCheckboxToggle(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		checked := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			Checkbox(c, &checked, "Agree")
		}, 300, 100)
		if err := tt.Click("Agree"); err != nil {
			t.Fatal(err)
		}
		if !checked {
			t.Errorf("dark=%v: click did not check", dark)
		}
		if err := tt.Click("Agree"); err != nil {
			t.Fatal(err)
		}
		if checked {
			t.Errorf("dark=%v: click did not uncheck", dark)
		}
	}
}

func TestCheckboxKeyboard(t *testing.T) {
	checked := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Checkbox(c, &checked, "Agree")
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !checked {
		t.Error("Space did not toggle the focused check box")
	}
}
