package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/plugins/chat"
)

func TestChatStaticPreview(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 640, 800)
	turn := s.chat.turns[0]
	if tt.HasText("read_source") || !tt.HasText("The review is complete.") {
		t.Fatal("completed preview missing or not folded")
	}
	for _, label := range []string{"Replay stream", "Simulate failure", "Retry", "Stop"} {
		if tt.HasText(label) {
			t.Fatalf("unexpected simulation control %q", label)
		}
	}
	for range 10 {
		tt.Frame()
	}
	if len(s.chat.turns) != 1 || turn.status != chat.Completed || turn.answer == nil {
		t.Fatal("static preview changed across frames")
	}
	tt.Scroll(120, 140, 0, -5000)
	if err := tt.Click("Worked for 9m 21s"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("read_source") {
		t.Fatal("completed details unavailable")
	}
	if navigationIcons["Chat"] == nil {
		t.Fatal("chat navigation icon missing")
	}
}

func TestChatResultAndComposerActions(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 700, 1600)
	click := func(label string) {
		t.Helper()
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	click("Review")
	if !s.chat.turns[0].showReview || !tt.HasText("Sample change preview") {
		t.Fatal("review action missing")
	}
	click("Close preview")
	click("Show 3 more files")
	if !s.chat.turns[0].summary.Expanded || !tt.HasText("styles.go") {
		t.Fatal("result expansion missing")
	}
	click("Attach sample context")
	if !s.chat.attachment || !tt.HasText("context.md") {
		t.Fatal("attachment missing")
	}
	click("Model")
	click("Fast demo")
	if s.chat.model != "Fast demo" {
		t.Fatal("model selection missing")
	}
	click("Message")
	tt.Type("Continue")
	click("Send")
	if len(s.chat.turns) != 1 || s.chat.turns[0].status != chat.Completed || !s.chat.attachment || s.chat.draft != "Continue" {
		t.Fatal("composer started a simulated conversation or changed the draft")
	}
}
