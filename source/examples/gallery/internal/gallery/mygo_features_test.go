package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func featureTester(name, variant, language string, width int) (*Gallery, *demoState, *ui.Tester) {
	g := &Gallery{language: newTranslator(language, "en")}
	s := g.mygoState(name, variant, true)
	example, _ := mygoExampleFor(pageFromKey("mygo/" + name))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		g.withLanguage(c, func() {
			ui.Column(c).FillWidth().Padding(16).Gap(12).Children(func() { example.build(c, s) })
		})
	}, width, 480)
	return g, s, tt
}

func TestTransformedEditorAndPopover(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		for _, dark := range []bool{false, true} {
			g, s, tt := featureTester("Transforms", "Editable controls", lang, 280)
			tt.SetDark(dark)
			if err := tt.Click(g.text("gallery.transform_scale")); err != nil {
				t.Fatal(err)
			}
			if err := tt.Click(g.text("gallery.transform_note")); err != nil {
				t.Fatal(err)
			}
			mods, _ := commandShortcut()
			tt.Key(mods, ui.KeyA)
			tt.Type("Edited note")
			if s.text != "Edited note" {
				t.Fatalf("transformed input = %q", s.text)
			}
			if err := tt.Click(g.text("gallery.show_popover")); err != nil {
				t.Fatal(err)
			}
			if !s.open || !tt.HasText(g.text("gallery.popover_content")) {
				t.Fatal("transformed anchor did not open")
			}
			anchor, _ := tt.Find(g.text("gallery.show_popover"))
			panel, _ := tt.Find(g.text("gallery.popover_content"))
			if panel.Y < anchor.Y || panel.X < 0 || panel.X+panel.W > 280 {
				t.Fatalf("popover placement = %+v", panel)
			}
			tt.Key(0, ui.KeyEscape)
			if s.open {
				t.Fatal("Escape did not dismiss")
			}
		}
	}
}

func TestMotionReplayWithReducedMotion(t *testing.T) {
	g, s, tt := featureTester("Motion", "Keyframes", "en", 280)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	before, ok := tt.Find(g.text("gallery.motion_tile"))
	if !ok {
		t.Fatal("missing animated tile")
	}
	for range 2 {
		if err := tt.Click(g.text("gallery.motion_replay")); err != nil {
			t.Fatal(err)
		}
		after, _ := tt.Find(g.text("gallery.motion_tile"))
		if after != before {
			t.Fatalf("reduced motion replay moved tile: %+v -> %+v", before, after)
		}
	}
	if s.count != 2 {
		t.Fatal("replay was not retained")
	}
	g, s, tt = featureTester("Motion", "Springs", "zh-Hans", 280)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	for _, want := range []bool{true, false} {
		if err := tt.Click(g.text("gallery.motion_reverse")); err != nil {
			t.Fatal(err)
		}
		if s.on != want {
			t.Fatal("spring target did not reverse")
		}
	}
}

func TestCommandPhysicalShortcut(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.View, 820, 600)
	mods, _ := commandShortcut()
	tt.KeyAt(mods, ui.KeyE)
	if !g.commands.open || !tt.Focused("Command palette") {
		t.Fatal("non-Latin physical shortcut did not open command search")
	}
	tt.Type("Bar chart")
	tt.Key(0, ui.KeyEnter)
	if pageKey(g.page) != "controls/Bar chart" || g.commands.open {
		t.Fatal("command search did not navigate to the chart example")
	}
}
