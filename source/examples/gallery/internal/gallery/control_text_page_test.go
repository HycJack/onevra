package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextPageGrouping(t *testing.T) {
	page := pageFromKey("controls/Text")
	for _, key := range []string{"controls/TextField", "controls/TextArea"} {
		if pageFromKey(key) != page {
			t.Fatalf("saved route %s did not migrate to Text", key)
		}
	}
	g := &Gallery{page: page}
	nav := ui.NewTester(g.Left, 220, 1800)
	if !nav.HasText("Text") || nav.HasText("Text fields") || nav.HasText("Text areas") || nav.HasText("Display text") {
		t.Fatal("text examples must share one sidebar entry")
	}
	tt := ui.NewTester(g.View, 900, 1800)
	for _, label := range []string{"Text", "TextField", "TextArea", "Count: 42", "Name", "Notes"} {
		if !tt.HasText(label) {
			t.Fatalf("missing grouped text example %s", label)
		}
	}
	for _, example := range []struct{ name, label string }{{"TextField", "Name"}, {"TextArea", "Notes"}} {
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
	if g.controlState("TextField", "Empty").text != "TextField" {
		t.Fatal("editing TextArea changed TextField")
	}
}
