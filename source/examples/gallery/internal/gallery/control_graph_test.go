package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery/locales"
	"github.com/oneo-me/onevra/source/theme"
)

func TestGraphPage(t *testing.T) {
	for _, language := range []string{"en", "zh-Hans"} {
		t.Run(language, func(t *testing.T) {
			g := &Gallery{language: newTranslator(language, "en")}
			g.page = pageFromKey("controls/Relationship graph")
			if example, ok := controlExampleFor(g.page); !ok || example.group != "Charts" {
				t.Fatal("graph is not a public component page")
			}
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); g.shellFrame(c) }, 520, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			// The layout settles on the next painted frame, which the
			// tester's clock must reach.
			tt.Advance(100 * time.Millisecond)
			s := g.controlState("Relationship graph", "Default").graph
			if !tt.HasText(g.text("Relationship graph")) {
				t.Fatal("graph page is not localized")
			}
			if !s.view.Settled() || s.view.Hovered() != "" {
				t.Fatal("graph did not settle")
			}
			// The graph fills the page below the titlebar, whose toolbar
			// holds its actions.
			view, ok := tt.Find(g.messages().GraphLabel(len(s.cache.data.Nodes), len(s.cache.data.Edges)))
			if !ok || view.Y > 100 || view.Y+view.H < 860 {
				t.Fatalf("graph %v does not fill the page", view)
			}
			nodes := len(graphLibrary(s, false, locales.In("")).Nodes)
			if nodes != len(graphNotes)+6 {
				t.Fatalf("sample has %d nodes", nodes)
			}

			if err := tt.Click(g.text("graphSettings")); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText(g.text("graphAttachments")) {
				t.Fatal("settings popover did not open")
			}
			// The switch is the popover's only control; the keyboard reaches it.
			for range 2 {
				if tt.Focused(g.text("graphAttachments")) {
					break
				}
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, ui.KeySpace)
			if s.attachments || len(graphLibrary(s, false, locales.In("")).Nodes) != len(graphNotes) {
				t.Fatal("switch did not hide attachments")
			}
			tt.Key(0, ui.KeyEscape)

			if err := tt.Click(g.text("graphAdd")); err != nil {
				t.Fatal(err)
			}
			if s.added != 1 {
				t.Fatal("add note did not update the graph")
			}
			if err := tt.Click(g.text("graphAddBulk")); err != nil {
				t.Fatal(err)
			}
			if err := tt.ChooseMenuItem(g.messages().GraphAddMany(1000)); err != nil {
				t.Fatal(err)
			}
			if s.added != 1001 || len(s.cache.data.Nodes) != len(graphNotes)+1001 {
				t.Fatalf("bulk add left %d added notes and %d nodes", s.added, len(s.cache.data.Nodes))
			}
			ids := map[string]bool{}
			for _, n := range s.cache.data.Nodes {
				if ids[n.ID] {
					t.Fatalf("duplicate node %q", n.ID)
				}
				ids[n.ID] = true
			}
			for _, e := range s.cache.data.Edges {
				if !ids[e.Source] || !ids[e.Target] {
					t.Fatalf("edge %s → %s leaves the sample", e.Source, e.Target)
				}
			}
			if err := tt.Click(g.text("graphReset")); err != nil {
				t.Fatal(err)
			}
			if s.added != 0 || tt.HasText(g.text("graphReset")) {
				t.Fatal("reset did not remove the added notes")
			}
			if err := tt.Click(g.text("chartEmpty")); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText(g.text("graphEmpty")) {
				t.Fatal("empty state is missing")
			}
		})
	}
}
