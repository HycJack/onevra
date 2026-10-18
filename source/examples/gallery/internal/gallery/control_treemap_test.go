package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestTreemapPage(t *testing.T) {
	for _, language := range []string{"en", "zh-Hans"} {
		t.Run(language, func(t *testing.T) {
			g := &Gallery{language: newTranslator(language, "en")}
			g.page = pageFromKey("controls/Treemap")
			if example, ok := controlExampleFor(g.page); !ok || example.group != "Charts" {
				t.Fatal("treemap is not a public component page")
			}
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); g.shellFrame(c) }, 720, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			tt.Frame()
			s := g.controlState("Treemap", "Default").treemap
			if !tt.HasText(g.text("Treemap")) || !tt.HasText(g.text("treemapKeys")) {
				t.Fatal("treemap page is not localized")
			}
			if !tt.HasText(g.text("treemapHome")) || !tt.HasText("12.8 GB") {
				t.Fatal("breadcrumbs or total are missing")
			}

			// Drilling shows the folder's path, whose crumbs lead back.
			s.view.Root = "Pictures/2025"
			tt.Frame()
			if err := tt.Click("Pictures"); err != nil {
				t.Fatal(err)
			}
			if s.view.Root != "Pictures" || s.view.Active != "Pictures/2025" {
				t.Fatalf("crumb left root %q, active %q", s.view.Root, s.view.Active)
			}

			if err := tt.Click(g.text("chartUpdate")); err != nil {
				t.Fatal(err)
			}
			if !s.updated {
				t.Fatal("update did not change the sample")
			}
			if err := tt.Click(g.text("chartEmpty")); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText(g.text("treemapEmpty")) {
				t.Fatal("empty state is missing")
			}
		})
	}
}
