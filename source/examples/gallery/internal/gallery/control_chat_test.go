package gallery

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/plugins/chat"
)

func TestChatStreamLifecycle(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 640, 800)
	if tt.HasText("read_source") {
		t.Fatal("completed demo not folded")
	}
	if err := tt.Click("Replay stream"); err != nil {
		t.Fatal(err)
	}
	turn := s.chat.turns[0]
	if !tt.HasText("read_source") || turn.status != chat.Running {
		t.Fatal("replay missing mixed blocks")
	}
	for range 4 {
		turn.advance(time.Now())
	}
	tt.Frame()
	if turn.answer == nil {
		t.Fatal("no partial answer")
	}
	partial := turn.answer
	tt.Click("Stop")
	turn.advance(time.Now())
	if turn.status != chat.Cancelled || turn.answer != partial {
		t.Fatal("stop changed partial response")
	}
	tt.Click("Retry")
	tt.Click("Simulate failure")
	if turn.status != chat.Failed || !tt.HasText("read_source") {
		t.Fatal("failed process hidden")
	}
	tt.Click("Retry")
	for turn.status == chat.Running {
		turn.advance(time.Now())
	}
	tt.Frame()
	if tt.HasText("read_source") || !tt.HasText("The review is complete.") {
		t.Fatal("final response did not replace execution")
	}
	tt.Scroll(120, 140, 0, -5000)
	seconds := max(0, int(turn.ended.Sub(turn.started).Seconds()))
	if err := tt.Click(fmt.Sprintf("Worked for %dm %ds", seconds/60, seconds%60)); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("read_source") {
		t.Fatal("completed details unavailable")
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
	click("Copy answer")
	if tt.Clipboard() != chatDemoAnswer {
		t.Fatal("copy included hidden execution or artifacts")
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
	if !s.chat.turns[len(s.chat.turns)-1].attached || s.chat.attachment || s.chat.draft != "" {
		t.Fatal("attachment was not transferred with draft")
	}
}
