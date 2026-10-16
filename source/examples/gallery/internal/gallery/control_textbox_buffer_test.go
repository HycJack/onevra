package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestTextBoxBufferExamples(t *testing.T) {
	g := &Gallery{}
	for _, name := range []string{"TextBoxBuffer", "TextBoxBuffer (multiline)"} {
		t.Run(name, func(t *testing.T) {
			example := sectionExample(name)
			if example.build == nil || example.group != "Text input" {
				t.Fatal("missing buffer example")
			}
			label := "Name"
			if name == "TextBoxBuffer (multiline)" {
				label = "Notes"
			}
			for _, variant := range append(append([]string(nil), example.variants...), "Disabled") {
				t.Run(variant, func(t *testing.T) {
					s := g.controlState(name, variant)
					before := s.buffer.String()
					if (before == "") != (variant == "Empty" || variant == "Search") {
						t.Fatal("incorrect initial buffer")
					}
					tt := ui.NewTester(func(c *ui.Context) {
						ui.Column(c).Fill().Disabled(variant == "Disabled").Children(func() { example.build(c, s, variant) })
					}, 540, 220)
					inputLabel := label
					if variant == "Search" {
						inputLabel = "Search samples"
					}
					if err := tt.Click(inputLabel); err != nil {
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
					if got := g.controlState(name, variant); got != s || got.buffer != s.buffer {
						t.Fatal("buffer lost on revisit")
					}
				})
			}
		})
	}
	seen := map[*ui.TextBuffer]bool{}
	for _, state := range g.controlStates {
		if seen[state.buffer] {
			t.Fatal("buffer shared between pages or variants")
		}
		seen[state.buffer] = true
	}
}
