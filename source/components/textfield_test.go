package components

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestTextFieldType(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		value := ""
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			TextField(c, &value).Label("Name").Placeholder("Ada").AutoFocus()
		}, 300, 100)
		tt.Type("Ada")
		if value != "Ada" {
			t.Errorf("dark=%v: typed %q", dark, value)
		}
	}
}

func TestTextAreaType(t *testing.T) {
	value := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Dark)
		TextArea(c, &value).Label("Notes").AutoFocus()
	}, 300, 120)
	tt.Type("line one")
	tt.Key(0, ui.KeyEnter)
	tt.Type("line two")
	if value != "line one\nline two" {
		t.Errorf("typed %q", value)
	}
}
