package gallery

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
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
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); g.View(c) }, 520, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			tt.Frame()
			s := g.controlState("Relationship graph", "Default").graph
			if !tt.HasText(g.text("Relationship graph")) || !tt.HasText(g.text("gallery.graph_hint")) || tt.HasText("gallery.graph_intro") {
				t.Fatal("graph page is not localized")
			}
			if !s.view.Settled() || s.view.Hovered() != "" {
				t.Fatal("graph did not settle")
			}
			// The graph fills the page below its controls, leaving the hint
			// at the bottom edge.
			view, ok := tt.Find(fmt.Sprintf(g.text("gallery.graph_label"), len(s.cache.data.Nodes), len(s.cache.data.Edges)))
			hint, hintOK := tt.Find(g.text("gallery.graph_hint"))
			if !ok || !hintOK || view.H < 460 || hint.Y+hint.H < 860 {
				t.Fatalf("graph %v and hint %v do not fill the page", view, hint)
			}
			nodes := len(graphLibrary(s, false, "%d").Nodes)
			if nodes != len(graphNotes)+6 {
				t.Fatalf("sample has %d nodes", nodes)
			}

			if err := tt.Click(g.text("gallery.graph_settings")); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText(g.text("gallery.graph_attachments_hint")) {
				t.Fatal("settings popover did not open")
			}
			// The switch is the popover's only control; the keyboard reaches it.
			for range 2 {
				if tt.Focused(g.text("gallery.graph_attachments")) {
					break
				}
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, ui.KeySpace)
			if s.attachments || len(graphLibrary(s, false, "%d").Nodes) != len(graphNotes) {
				t.Fatal("switch did not hide attachments")
			}
			tt.Key(0, ui.KeyEscape)

			if err := tt.Click(g.text("gallery.graph_add")); err != nil {
				t.Fatal(err)
			}
			if s.added != 1 {
				t.Fatal("add note did not update the graph")
			}
			if err := tt.Click(fmt.Sprintf(g.text("gallery.graph_add_many"), 1000)); err != nil {
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
			if err := tt.Click(g.text("gallery.graph_reset")); err != nil {
				t.Fatal(err)
			}
			if s.added != 0 || tt.HasText(g.text("gallery.graph_reset")) {
				t.Fatal("reset did not remove the added notes")
			}
			if err := tt.Click(g.text("gallery.chart_empty")); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText(g.text("gallery.graph_empty")) {
				t.Fatal("empty state is missing")
			}
		})
	}
}
