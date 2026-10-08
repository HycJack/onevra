package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTreePageVariants(t *testing.T) {
	g := &Gallery{}
	var tree controlExample
	for _, example := range controlCatalog {
		if example.name == "TreeVirtualized" {
			t.Fatal("virtualized tree must not have a separate page")
		}
		if example.name == "Tree" {
			tree = example
		}
	}
	if len(tree.variants) != 2 {
		t.Fatal("Tree page must contain both tree variants")
	}
	page := ui.NewTester(func(c *ui.Context) { g.Controls(c, tree) }, 540, 900)
	if !page.HasText("Tree") || !page.HasText("TreeVirtualized") {
		t.Fatal("missing tree section headings")
	}
	for _, variant := range tree.variants {
		for _, disabled := range []bool{false, true} {
			s := (&Gallery{}).controlState("Tree", variant)
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Disabled(disabled).Children(func() {
					tree.build(c, s, variant)
				})
			}, 540, 360)
			for range 2 {
				if err := tt.Click("source"); err != nil {
					t.Fatal(err)
				}
			}
			open := s.on
			if variant == "TreeVirtualized" {
				open = s.tree.Open.Has("source")
			}
			if open == disabled {
				t.Fatalf("%s: open=%v, disabled=%v", variant, open, disabled)
			}
		}
	}
	normal, virtual := g.controlState("Tree", "Tree"), g.controlState("Tree", "TreeVirtualized")
	if normal == virtual {
		t.Fatal("tree variants share state")
	}
	normal.on = true
	if virtual.on || virtual.tree.Open.Has("source") {
		t.Fatal("ordinary tree changed virtualized tree state")
	}
	for _, name := range []string{"Tree", "TreeItem", "Outline"} {
		if mygoHandled[name].Replacement == "" {
			t.Fatalf("missing replacement for %s", name)
		}
	}
}
