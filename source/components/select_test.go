package components

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func buildSelect(c *ui.Context, selected *string, options []string) {
	sel := Select(c, selected)
	sel.Trigger.Label("size").Children(func() { ui.Text(c, *selected).SingleLine() })
	sel.Popup(func(panel ui.Element) {
		for _, opt := range options {
			item := sel.Item(opt).Label("item " + opt)
			item.Children(func() { ui.Text(c, opt).SingleLine() })
		}
	})
}

func TestSelectItemBackgrounds(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		selected := "M"
		var accent, hover ui.Color
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			accent, hover = c.Theme().Accent, c.Theme().SurfaceHover
			buildSelect(c, &selected, []string{"S", "M", "L"})
		}, 320, 320)
		tt.Click("size")
		assertFill := func(label string, want ui.Color) {
			t.Helper()
			r, ok := tt.Find(label)
			if !ok {
				t.Fatalf("missing %s", label)
			}
			pixel := tt.Image().RGBAAt(int(r.X+4), int(r.Y+r.H/2))
			if pixel.R != want.R || pixel.G != want.G || pixel.B != want.B {
				t.Fatalf("mode=%v %s fill=%v, want %v", mode, label, pixel, want)
			}
		}
		moveTo := func(label string) {
			t.Helper()
			r, _ := tt.Find(label)
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
		}
		assertFill("item M", accent)
		moveTo("item M")
		assertFill("item M", accent)
		moveTo("item S")
		assertFill("item S", hover)
		assertFill("item M", accent)
		tt.Move(310, 310)
		tt.Key(0, ui.KeyDown)
		tt.Key(0, ui.KeyDown)
		assertFill("item L", hover)
		assertFill("item M", accent)
		tt.Click("item L")
		tt.Click("size")
		assertFill("item L", accent)
		moveTo("item M")
		assertFill("item M", hover)
		assertFill("item L", accent)
	}
}

func TestSelectPopupShadow(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		selected := "M"
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			ui.Column(c).Fill().Padding(40).Background(c.Theme().Background).Children(func() {
				sel := Select(c, &selected)
				sel.Trigger.Label("size").Children(func() { Text(c, selected) })
				sel.Popup(func(panel ui.Element) {
					panel.Label("options").Width(180)
					sel.Item("L").Padding(10).Children(func() { Text(c, "Large") })
				})
			})
		}, 320, 320)
		tt.Click("size")
		panel, ok := tt.Find("options")
		if !ok {
			t.Fatal("popup missing")
		}
		img := tt.Image()
		x, y := int(panel.X+panel.W/2), int(panel.Y+panel.H+3)
		if img.RGBAAt(x, y) == img.RGBAAt(310, 310) {
			t.Fatalf("mode=%v: no shadow below popup", mode)
		}
	}
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

func TestSelectSearch(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		selected, query := "M", ""
		var changed bool
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			sel := Select(c, &selected)
			sel.Trigger.Label("size").Children(func() { Text(c, selected) })
			sel.Popup(func(panel ui.Element) {
				panel.MinWidth(200)
				sel.SearchField(&query).Label("Search sizes").Placeholder("Search…")
				matches := 0
				for _, option := range []string{"S", "M", "L"} {
					if !strings.Contains(strings.ToLower(option), strings.ToLower(query)) {
						continue
					}
					matches++
					sel.Item(option).Label("item " + option).Padding(8).Children(func() { Text(c, option) })
				}
				if matches == 0 {
					Text(c, "No matches")
				}
			})
			changed = changed || sel.Trigger.Changed()
		}, 320, 320)
		if err := tt.Click("size"); err != nil {
			t.Fatal(err)
		}
		search, ok := tt.Find("Search sizes")
		if !ok {
			t.Fatal("search input missing")
		}
		first, ok := tt.Find("item S")
		if !ok || first.Y <= search.Y+search.H {
			t.Fatal("missing spacing below search input")
		}
		tt.Type("l")
		if query != "l" || !tt.HasText("item L") || tt.HasText("item M") || selected != "M" || changed {
			t.Fatalf("search did not filter independently: query=%q selected=%q changed=%t", query, selected, changed)
		}
		if err := tt.Click("item L"); err != nil {
			t.Fatal(err)
		}
		if selected != "L" || !changed || tt.HasText("Search sizes") {
			t.Fatal("filtered selection did not choose and close")
		}
		query = "missing"
		tt.Click("size")
		if !tt.HasText("No matches") {
			t.Fatal("empty results message missing")
		}
		tt.Key(0, ui.KeyEnter)
		if selected != "L" {
			t.Fatal("empty results changed selection")
		}
		tt.Key(0, ui.KeyEscape)
		if tt.HasText("Search sizes") {
			t.Fatal("Escape from search did not close")
		}
	}
}

func TestSelectSearchKeyboard(t *testing.T) {
	selected, query := "M", ""
	tt := ui.NewTester(func(c *ui.Context) {
		sel := Select(c, &selected)
		sel.Trigger.Label("size").Children(func() { Text(c, selected) })
		sel.Popup(func(panel ui.Element) {
			panel.MinWidth(200)
			sel.SearchField(&query).Label("Search sizes")
			for _, option := range []string{"S", "M", "L"} {
				if strings.Contains(strings.ToLower(option), strings.ToLower(query)) {
					sel.Item(option).Children(func() { Text(c, option) })
				}
			}
		})
	}, 320, 320)
	tt.Click("size")
	tt.Type("l")
	tt.Key(0, ui.KeyDown) // Return focus to the trigger.
	tt.Key(0, ui.KeyDown) // Navigate the filtered options with MyGo.
	tt.Key(0, ui.KeyEnter)
	if selected != "L" || tt.HasText("Search sizes") {
		t.Fatal("keyboard did not choose the filtered result")
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
