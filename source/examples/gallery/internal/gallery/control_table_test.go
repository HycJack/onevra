package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTablePageVariants(t *testing.T) {
	g := &Gallery{}
	var table controlExample
	for _, example := range controlCatalog {
		if example.name == "TableVirtualized" {
			t.Fatal("virtualized table must not have a separate page")
		}
		if example.name == "Table" {
			table = example
		}
	}
	if len(table.variants) != 2 {
		t.Fatal("Table page must contain both table variants")
	}
	page := ui.NewTester(func(c *ui.Context) { g.Controls(c, table) }, 540, 900)
	if !page.HasText("Table") || !page.HasText("TableVirtualized") {
		t.Fatal("missing table section headings")
	}
	for _, variant := range table.variants {
		for _, disabled := range []bool{false, true} {
			s := (&Gallery{}).controlState("Table", variant)
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Disabled(disabled).Children(func() {
					table.build(c, s, variant)
				})
			}, 540, 360)
			// Size, then reversed: the smallest file comes last.
			for range 2 {
				r, ok := tt.Find("Size")
				if !ok {
					t.Fatal("missing Size header")
				}
				tt.ClickAt(r.X, r.Y+r.H/2)
			}
			sorted := s.tableSort == ui.SortOrder{Column: "size", Descending: true}
			if sorted == disabled {
				t.Fatalf("%s: sort=%+v, disabled=%v", variant, s.tableSort, disabled)
			}
			if variant == "TableVirtualized" && !disabled && s.tableRows[0] == 0 {
				t.Fatal("virtualized rows were not re-sorted")
			}
			if variant == "Table" {
				if err := tt.Click("go.mod"); err != nil {
					t.Fatal(err)
				}
				if (s.choice == "go.mod") == disabled {
					t.Fatalf("row selection = %q, disabled=%v", s.choice, disabled)
				}
			}
		}
	}
	normal, virtual := g.controlState("Table", "Table"), g.controlState("Table", "TableVirtualized")
	if normal == virtual {
		t.Fatal("table variants share state")
	}
	if mygoHandled["Table"].Replacement == "" {
		t.Fatal("missing replacement for Table")
	}
}
