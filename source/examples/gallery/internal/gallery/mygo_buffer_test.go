package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMyGoBufferExamples(t *testing.T) {
	g := &Gallery{}
	for _, name := range []string{"TextInputBuffer", "TextAreaBuffer"} {
		t.Run(name, func(t *testing.T) {
			var example mygoExample
			for _, entry := range mygoCatalog {
				if entry.name == name {
					example = entry
					break
				}
			}
			if example.build == nil || example.group != "Text input" {
				t.Fatal("missing buffer example")
			}
			label := "Name"
			if name == "TextAreaBuffer" {
				label = "Notes"
			}
			for i, variant := range mygoVariants(name) {
				t.Run(variant, func(t *testing.T) {
					s := g.mygoState(name, variant, i == 0)
					before := s.buffer.String()
					if (before == "") != (variant == "Empty") {
						t.Fatal("incorrect initial buffer")
					}
					tt := ui.NewTester(func(c *ui.Context) {
						ui.Column(c).Fill().Disabled(variant == "Disabled").Children(func() { example.build(c, s) })
					}, 540, 220)
					if err := tt.Click(label); err != nil {
						t.Fatal(err)
					}
					tt.Key(ui.Cmd, ui.KeyA)
					tt.Type("Edited 世界 🌍")
					if variant == "Read only" || variant == "Disabled" {
						if s.buffer.String() != before {
							t.Fatal("non-editable buffer changed")
						}
					} else {
						if s.buffer.String() != "Edited 世界 🌍" {
							t.Fatalf("edit not published: %q", s.buffer.String())
						}
						tt.Key(ui.Cmd, ui.KeyZ)
						if s.buffer.String() != before {
							t.Fatal("undo did not restore buffer")
						}
						s.buffer.Set("External update")
						tt.Frame()
						tt.Key(ui.Cmd, ui.KeyA)
						tt.Type("Retained edit")
						if s.buffer.String() != "Retained edit" {
							t.Fatal("external update lost editability")
						}
					}
					if got := g.mygoState(name, variant, i == 0); got != s || got.buffer != s.buffer {
						t.Fatal("buffer lost on revisit")
					}
				})
			}
		})
	}
	seen := map[*ui.TextBuffer]bool{}
	for _, state := range g.mygoStates {
		if seen[state.buffer] {
			t.Fatal("buffer shared between pages or variants")
		}
		seen[state.buffer] = true
	}
}
