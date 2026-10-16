package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/richtext"
)

func TestRichTextEditorExample(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		g, s, tt := featureTester("Rich text editor", "Default", lang, 540)
		if err := tt.Click(g.text("richTextLabel")); err != nil {
			t.Fatal(err)
		}
		mods, _ := commandShortcut()
		tt.Key(mods, ui.KeyA)
		tt.Frame()
		// The toolbar formats the editor's selection without taking focus.
		if err := tt.Click(g.text("bold")); err != nil {
			t.Fatal(err)
		}
		if f := s.richText.Formatting(); f.Bold != richtext.FormatOn {
			t.Fatalf("%s: bold = %v after formatting the selection", lang, f.Bold)
		}
		tt.Type("Shipped")
		tt.Frame()
		if got := s.richText.Text().PlainText(); got != "Shipped" {
			t.Fatalf("%s: text %q, want the selection replaced", lang, got)
		}
	}
}

func TestMarkdownStreamExample(t *testing.T) {
	g, s, tt := featureTester("Markdown stream", "Default", "zh-Hans", 540)
	// Reduced motion settles the reveal at once, showing the finished reply.
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click(g.text("markdownStreamReplay")); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if s.count != 1 || !strings.Contains(strings.Join(tt.Texts(), "\n"), "Finished blocks are never parsed again") {
		t.Fatalf("reply not shown: %q", tt.Texts())
	}
}

func TestTimelineScrubPauses(t *testing.T) {
	g, s, tt := featureTester("Motion", "Timeline", "en", 540)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	if s.on || s.value != 1 {
		t.Fatalf("reduced motion: paused %v, progress %v; want the timeline at rest", s.on, s.value)
	}
	if err := tt.Click(g.text("timelineScrub")); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	if !s.on || s.value != 0 || s.seek != 0 {
		t.Fatalf("scrubbing to the start: paused %v, progress %v, seek %v", s.on, s.value, s.seek)
	}
	if !tt.HasText(g.text("timelinePlay")) {
		t.Fatalf("paused timeline offers no Play: %q", tt.Texts())
	}
}
