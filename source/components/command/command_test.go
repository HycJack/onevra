package command

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func TestFilter(t *testing.T) {
	items := []Item{
		{ID: "keyword", Title: "Preferences", Keywords: []string{"settings", "shezhi"}},
		{ID: "substring", Title: "Open Settings"},
		{ID: "prefix", Title: "Settings page"},
		{ID: "exact", Title: "Settings"},
		{ID: "chinese", Title: "设置"},
	}
	for _, tc := range []struct {
		query string
		ids   []string
	}{
		{" SETTINGS ", []string{"exact", "prefix", "substring", "keyword"}},
		{"settings page", []string{"prefix"}},
		{"设", []string{"chinese"}},
		{"shezhi", []string{"keyword"}},
		{"missing", nil},
		{"", []string{"keyword", "substring", "prefix", "exact", "chinese"}},
	} {
		var ids []string
		for _, item := range Filter(items, tc.query) {
			ids = append(ids, item.ID)
		}
		if !slices.Equal(ids, tc.ids) {
			t.Fatalf("%q: %v, want %v", tc.query, ids, tc.ids)
		}
	}
	if items[0].ID != "keyword" {
		t.Fatal("filter mutated input")
	}
}

func TestDialogInteraction(t *testing.T) {
	open := false
	state := State{}
	items := []Item{{ID: "one", Title: "First", Group: "Actions"}, {ID: "disabled", Title: "Unavailable", Disabled: true}, {ID: "two", Title: "Second", Group: "Pages"}}
	var activations []string
	tt := ui.NewTester(func(c *ui.Context) {
		if components.Button(c, "Open", components.ButtonPrimary).Clicked() {
			open = true
		}
		CommandDialog(c, &open, &state, items, Options{OnActivate: func(id string) { activations = append(activations, id) }}, components.DialogOptions{})
	}, 800, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Search commands") || state.ActiveID != "one" {
		t.Fatalf("opening focus/selection: %s", state.ActiveID)
	}
	tt.Key(0, ui.KeyDown)
	if state.ActiveID != "two" || !tt.Focused("Search commands") {
		t.Fatalf("navigation: %s", state.ActiveID)
	}
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(activations, []string{"two"}) {
		t.Fatalf("activation: %v", activations)
	}
	tt.Type("First")
	if state.Query != "First" || state.ActiveID != "one" || tt.HasText("Second") {
		t.Fatalf("filter: %+v", state)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || !tt.Focused("Open") {
		t.Fatalf("dismissal: open=%v opener=%v search=%v", open, tt.Focused("Open"), tt.Focused("Search commands"))
	}
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if state.Query != "First" || !tt.Focused("Search commands") {
		t.Fatal("reopening lost query or focus")
	}
	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !slices.Equal(activations, []string{"two", "one"}) || !tt.Focused("Search commands") {
		t.Fatalf("single click: %v", activations)
	}
}

func TestCompositionAndUnavailableResults(t *testing.T) {
	state := State{}
	open, loading, failure := true, false, ""
	activated := 0
	items := []Item{{ID: "settings", Title: "设置", Keywords: []string{"settings"}}}
	tt := ui.NewTester(func(c *ui.Context) {
		CommandDialog(c, &open, &state, items, Options{Loading: loading, Error: failure, OnActivate: func(string) { activated++ }}, components.DialogOptions{})
	}, 640, 520)
	tt.Compose("she", 3)
	tt.Key(0, ui.KeyDown)
	if activated != 0 || !open || state.Query != "" {
		t.Fatal("composition navigation escaped input")
	}
	tt.Key(0, ui.KeyEnter)
	if activated != 0 || !open {
		t.Fatal("composition confirmation activated command")
	}
	// The headless editor commits marked text on Enter, just as accepting an
	// IME candidate does. Start a new composition to test Escape independently.
	state.Query = ""
	tt.Frame()
	tt.Compose("she", 3)
	tt.Key(0, ui.KeyEscape)
	if activated != 0 || !open {
		t.Fatal("composition cancellation closed dialog")
	}
	tt.Compose("", 0)
	state.Query = ""
	tt.Frame()
	tt.Type("设置")
	tt.Key(0, ui.KeyEnter)
	if activated != 1 {
		t.Fatal("committed Chinese query did not activate")
	}
	for _, mode := range []string{"loading", "error", "empty", "disabled"} {
		loading, failure = mode == "loading", ""
		if mode == "error" {
			failure = "Search failed"
		}
		if mode == "empty" {
			items = nil
		}
		if mode == "disabled" {
			items = []Item{{ID: "settings", Title: "设置", Disabled: true}}
		}
		tt.Frame()
		tt.Key(0, ui.KeyEnter)
		if activated != 1 || state.ActiveID != "" {
			t.Fatalf("%s allowed activation", mode)
		}
	}
}

func TestExternalResultsAndScrolling(t *testing.T) {
	state := State{Query: "host query"}
	items := make([]Item, 200)
	for i := range items {
		items[i] = Item{ID: fmt.Sprint(i), Title: fmt.Sprintf("Command %03d", i)}
	}
	var activated string
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Size(400, 250).Children(func() {
			CommandPalette(c, &state, items, Options{PreFiltered: true, OnActivate: func(id string) { activated = id }})
		})
	}, 400, 250)
	if err := tt.Click("Search commands"); err != nil {
		t.Fatal(err)
	}
	for range 30 {
		tt.Key(0, ui.KeyDown)
	}
	if state.ActiveID != "30" || !tt.HasText("Command 030") {
		t.Fatal("navigation did not reveal active command")
	}
	items = []Item{items[30], items[1]}
	tt.Frame()
	if state.ActiveID != "30" {
		t.Fatal("reordering lost stable selection")
	}
	items = items[1:]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if activated != "1" {
		t.Fatalf("removed active command was not reconciled: %s", activated)
	}
}

func TestExpandOnSearch(t *testing.T) {
	for _, dialog := range []bool{false, true} {
		t.Run(fmt.Sprintf("dialog=%v", dialog), func(t *testing.T) {
			s := State{ActiveID: "settings"}
			open, activations := true, 0
			options := Options{ExpandOnSearch: true, OnActivate: func(string) { activations++ }}
			items := []Item{{ID: "settings", Title: "设置", Keywords: []string{"settings"}}}
			var bounds ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				var surface ui.Element
				if dialog {
					surface = CommandDialog(c, &open, &s, items, options, components.DialogOptions{})
				} else {
					ui.Column(c).Size(500, 440).Children(func() { surface = CommandPalette(c, &s, items, options) })
				}
				surface.DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
			}, 700, 600)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			if err := tt.Click("Search commands"); err != nil {
				t.Fatal(err)
			}
			compact := bounds.H
			if compact <= 0 || compact > 100 || tt.HasText("设置") || tt.HasText("No matching commands") || s.ActiveID != "" {
				t.Fatalf("not collapsed: bounds=%+v active=%q", bounds, s.ActiveID)
			}
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if activations != 0 {
				t.Fatal("hidden result activated")
			}
			tt.Type("  ")
			if bounds.H != compact || tt.HasText("设置") {
				t.Fatal("whitespace expanded results")
			}
			tt.Compose("she", 3)
			if bounds.H != compact {
				t.Fatal("uncommitted composition expanded results")
			}
			tt.Type("设置")
			if bounds.H <= compact || !tt.HasText("设置") || !tt.Focused("Search commands") {
				t.Fatal("committed query did not expand with retained focus")
			}
			tt.Key(0, ui.KeyEnter)
			if activations != 1 {
				t.Fatal("visible command did not activate")
			}
			tt.Command("selectAll")
			tt.Key(0, ui.KeyBackspace)
			if bounds.H != compact || s.ActiveID != "" || tt.HasText("设置") {
				t.Fatal("clearing query did not collapse")
			}
			tt.Key(0, ui.KeyEnter)
			if activations != 1 {
				t.Fatal("cleared query activated stale result")
			}
			options.PreFiltered = true
			for _, status := range []string{"loading", "error", "empty"} {
				options.Loading, options.Error = status == "loading", ""
				if status == "error" {
					options.Error = "Search failed"
				}
				if status == "empty" {
					items = nil
				}
				s.Query = ""
				tt.Frame()
				if bounds.H != compact {
					t.Fatalf("%s expanded without query", status)
				}
				s.Query = "query"
				tt.Frame()
				if bounds.H <= compact {
					t.Fatalf("%s remained collapsed after query", status)
				}
			}
			options = Options{}
			s.Query = ""
			tt.Frame()
			if bounds.H <= compact || !tt.HasText("No matching commands") {
				t.Fatal("default mode did not expand immediately")
			}
		})
	}
}
