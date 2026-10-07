package ui

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func buildSelect(c *ui.Context, selected *string, options []string) {
	t := c.Theme()
	sel := Select(c, selected)
	sel.Trigger.Label("size").Children(func() { ui.Text(c, *selected).SingleLine() })
	sel.Popup(func(panel *ui.Element) {
		panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		for _, opt := range options {
			item := sel.Item(opt).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Label("item " + opt)
			if item.Highlighted() {
				item.Background(t.Accent).TextColor(t.AccentText)
			}
			item.Children(func() { ui.Text(c, opt).SingleLine() })
		}
	})
}

func TestSelectChoose(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := theme.Light
		if dark {
			mode = theme.Dark
		}
		selected := "M"
		options := []string{"S", "M", "L"}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			buildSelect(c, &selected, options)
		}, 300, 300)
		if tt.HasText("item S") {
			t.Fatal("popup shows before the trigger is clicked")
		}
		if err := tt.Click("size"); err != nil {
			t.Fatal(err)
		}
		if !tt.HasText("item S") {
			t.Fatal("popup did not open")
		}
		if err := tt.Click("item L"); err != nil {
			t.Fatal(err)
		}
		if selected != "L" {
			t.Errorf("dark=%v: selected=%q, want L", dark, selected)
		}
		if tt.HasText("item S") {
			t.Errorf("dark=%v: popup stayed open after choosing", dark)
		}
	}
}

func TestSelectKeyboard(t *testing.T) {
	selected := "M"
	options := []string{"S", "M", "L"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		buildSelect(c, &selected, options)
	}, 300, 300)
	if err := tt.Click("size"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if selected != "L" {
		t.Errorf("Down+Enter chose %q, want L", selected)
	}
	// Escape closes an open popup.
	if err := tt.Click("size"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("item S") {
		t.Error("Escape did not close the popup")
	}
}
