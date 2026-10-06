package gallery

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestExampleActions(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.View, 820, 1100)
	if err := tt.Click("Create item"); err != nil {
		t.Fatal(err)
	}
	if !g.Created || !tt.HasText("Item created") {
		t.Fatal("missing created state")
	}
	if err := tt.Click("Reset example"); err != nil {
		t.Fatal(err)
	}
	if g.Created || tt.HasText("Item created") {
		t.Fatal("reset retained success state")
	}
	if err := tt.Click("Disable example action"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Create item"); err != nil {
		t.Fatal(err)
	}
	if g.Created {
		t.Fatal("disabled action activated")
	}
}
