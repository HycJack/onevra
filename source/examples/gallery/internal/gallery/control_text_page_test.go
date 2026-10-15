package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextBoxPageGrouping(t *testing.T) {
	page := pageFromKey("controls/TextBox")
	textPage := pageFromKey("controls/Text")
	if page == textPage || pageKey(page) != "controls/TextBox" {
		t.Fatal("TextBox must have its own persistent route")
	}
	display := ui.NewTester((&Gallery{page: textPage}).View, 900, 1800)
	if !display.HasText("Count: 42") || display.HasText("TextBox") {
		t.Fatal("Text page must contain only display text examples")
	}
	g := &Gallery{page: page}
	nav := ui.NewTester(g.Left, 220, 1800)
	if !nav.HasText("Text") || !nav.HasText("TextBox") || nav.HasText("Text fields") || nav.HasText("Text areas") || nav.HasText("Display text") {
		t.Fatal("Text and TextBox must have separate sidebar entries")
	}
	tt := ui.NewTester(g.View, 900, 1800)
	if tt.HasText("Count: 42") {
		t.Fatal("display text example remains on TextBox page")
	}
	for _, label := range []string{"TextBox", "TextBox (multiline)", "Name", "Notes", "Search"} {
		if !tt.HasText(label) {
			t.Fatalf("missing grouped text example %s", label)
		}
	}
	for _, example := range []struct{ name, label string }{{"TextBox", "Name"}, {"TextBox (multiline)", "Notes"}} {
		if err := tt.Click(example.label); err != nil {
			t.Fatal(err)
		}
		tt.Type(example.name)
		if g.controlState(example.name, "Empty").text != example.name {
			t.Fatalf("%s is not editable independently", example.name)
		}
		for _, variant := range []string{"Filled", "Read only"} {
			if g.controlState(example.name, variant).text == example.name {
				t.Fatal("text variants share state")
			}
		}
	}
	if g.controlState("TextBox", "Empty").text != "TextBox" {
		t.Fatal("editing multiline TextBox changed single-line TextBox")
	}
}
