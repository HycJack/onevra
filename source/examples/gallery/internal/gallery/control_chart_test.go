package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestChartPages(t *testing.T) {
	var names []string
	for _, group := range navigationCatalog() {
		if group.label == "Charts" {
			for _, entry := range group.entries {
				if entry.label != "Relationship graph" {
					names = append(names, entry.label)
				}
			}
		}
	}
	if len(names) != 5 {
		t.Fatalf("chart group = %v", names)
	}
	for _, language := range []string{"en", "zh-Hans"} {
		g := &Gallery{language: newTranslator(language, "en")}
		for _, name := range names {
			t.Run(language+"/"+name, func(t *testing.T) {
				g.page = pageFromKey("controls/" + name)
				example, ok := controlExampleFor(g.page)
				if !ok || example.group != "Charts" {
					t.Fatal("chart is not a public component page")
				}
				tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); g.View(c) }, 280, 900)
				for _, dark := range []bool{false, true} {
					tt.SetDark(dark)
					if !tt.HasText(g.text(name)) || tt.HasText("chartIntro") {
						t.Fatal("chart page is not localized")
					}
					for _, label := range []string{"chartUpdate", "chartEmpty"} {
						if err := tt.Click(g.text(label)); err != nil {
							t.Fatal(err)
						}
					}
					if !tt.HasText(g.text("chartNoData")) {
						t.Fatal("empty state is missing")
					}
					if err := tt.Click(g.text("chartRestore")); err != nil {
						t.Fatal(err)
					}
					r, ok := tt.Find(g.text("chartSelection"))
					if !ok || r.X < 0 || r.X+r.W > 280 || r.W < 80 {
						t.Fatalf("narrow plot bounds = %+v", r)
					}
					tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
					tt.Key(0, ui.KeyHome)
					tt.Key(0, ui.KeyRight)
					if g.controlState(name, "Default").chart.selection.Index != 1 {
						t.Fatal("chart is not keyboard selectable")
					}
				}
			})
		}
		for i, name := range names {
			state := g.controlState(name, "Default").chart
			for _, other := range names[i+1:] {
				if state == g.controlState(other, "Default").chart {
					t.Fatal("chart pages share state")
				}
			}
		}
	}
}
