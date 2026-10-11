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
				names = append(names, entry.label)
			}
		}
	}
	if len(names) != 5 {
		t.Fatalf("chart group = %v", names)
	}
	if pageFromKey("mygo/Charts") != pageFromKey("controls/Bar chart") {
		t.Fatal("legacy chart route was not restored")
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
				tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c, theme.Auto); g.View(c) }, 280, 900)
				for _, dark := range []bool{false, true} {
					tt.SetDark(dark)
					if !tt.HasText(g.text(name)) || tt.HasText("gallery.chart_intro") {
						t.Fatal("chart page is not localized")
					}
					for _, label := range []string{"gallery.chart_update", "gallery.chart_empty"} {
						if err := tt.Click(g.text(label)); err != nil {
							t.Fatal(err)
						}
					}
					if !tt.HasText(g.text("gallery.chart_no_data")) {
						t.Fatal("empty state is missing")
					}
					if err := tt.Click(g.text("gallery.chart_restore")); err != nil {
						t.Fatal(err)
					}
					r, ok := tt.Find(g.text("gallery.chart_selection"))
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
