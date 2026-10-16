package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestChartPages(t *testing.T) {
	page := pageFromKey("controls/Data charts")
	example, ok := controlExampleFor(page)
	if !ok || example.group != "Charts" || len(example.variants) != 5 {
		t.Fatal("charts are not one public component page")
	}
	for _, language := range []string{"en", "zh-Hans"} {
		t.Run(language, func(t *testing.T) {
			g := &Gallery{page: page, language: newTranslator(language, "en")}
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); g.shellFrame(c) }, 280, 2400)
			for _, dark := range []bool{false, true} {
				tt.SetDark(dark)
				for _, name := range example.variants {
					if !tt.HasText(g.text(name)) {
						t.Fatalf("chart section %s is not localized", name)
					}
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
			}
			if len(g.controlStates) != 1 {
				t.Fatal("charts do not share the page's data actions")
			}
		})
	}
}

func TestChartSelectionsIndependent(t *testing.T) {
	s := (&Gallery{}).controlState("Data charts", "Default")
	example, _ := controlExampleFor(pageFromKey("controls/Data charts"))
	names := example.variants
	for i, name := range names {
		tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); controlChart(c, s, name) }, 540, 400)
		r, ok := tt.Find("Chart selection")
		if !ok {
			t.Fatalf("%s: missing plot", name)
		}
		tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
		tt.Key(0, ui.KeyHome)
		tt.Key(0, ui.KeyRight)
		for j, other := range names {
			if want := map[bool]int{true: 1, false: -1}[j <= i]; chartDemo(s).selection(other).Index != want {
				t.Fatalf("after %s, %s selection = %d", name, other, chartDemo(s).selection(other).Index)
			}
		}
	}
}
