package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chat"
)

func TestChatStaticPreview(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 640, 800)
	turn := s.chat.turns[0]
	if tt.HasText("read_source") || !tt.HasText("The review is complete.") {
		t.Fatal("completed preview missing or not folded")
	}
	for _, label := range []string{"Changed 6 files", "Sample change preview", "Chat"} {
		if tt.HasText(label) {
			t.Fatalf("unexpected page content %q", label)
		}
	}
	for range 10 {
		tt.Frame()
	}
	if len(s.chat.turns) != 1 || turn.status != chat.Completed || turn.answer == "" {
		t.Fatal("static preview changed across frames")
	}
	tt.Scroll(120, 140, 0, -5000)
	if err := tt.Click("Worked for 9m 21s"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("read_source") {
		t.Fatal("completed details unavailable")
	}
	if navigationIcons["Chat"] == nil || controlToolbars["Chat"] == nil {
		t.Fatal("chat navigation icon or titlebar action missing")
	}
}

func TestChatComposerControls(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 700, 900)
	click := func(label string) {
		t.Helper()
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	click("Demo model")
	if !tt.HasText("Remote") || !tt.HasText("1.0M") {
		t.Fatal("model groups or context sizes missing")
	}
	click("Deep demo")
	click("Reasoning effort")
	if !tt.HasText("Reasoning effort") {
		t.Fatal("effort panel missing")
	}
	for _, key := range []ui.Key{ui.KeyHome, ui.KeyRight, ui.KeyRight, ui.KeyRight, ui.KeyEscape} {
		tt.Key(0, key)
	}
	click("main")
	click("feature/chat-ui")
	click("Worktree")
	if s.chat.model != "deep" || s.chat.effort != "high" || s.chat.branch != "feature/chat-ui" || s.chat.worktree {
		t.Fatalf("composer choices: %s %s %s %v", s.chat.model, s.chat.effort, s.chat.branch, s.chat.worktree)
	}
	click("Attach a sample file")
	click("Attach a sample file")
	if len(s.chat.files) != 2 || !tt.HasText("context.md") || !tt.HasText("screenshot.png") {
		t.Fatal("attachments missing")
	}
	click("Remove attachment")
	if len(s.chat.files) != 1 || s.chat.files[0].name != "screenshot.png" {
		t.Fatalf("removal: %+v", s.chat.files)
	}
	click("Attach a sample file")
	click("Attach a sample file")
	click("Attach a sample file")
	if len(s.chat.files) != 3 {
		t.Fatalf("samples should attach once each: %+v", s.chat.files)
	}
	click("Usage limits")
	if !s.chat.quota.Open || !tt.HasText("Weekly limit") || !tt.HasText("$12.40") {
		t.Fatal("quota panel missing")
	}
	if !tt.HasText("05:00") {
		t.Fatal("refresh countdown missing")
	}
	click("Refresh limits")
	if !tt.HasText("Refreshing…") {
		t.Fatal("refresh state missing")
	}
	tt.Advance(1500 * time.Millisecond)
	if tt.HasText("Refreshing…") || !tt.HasText("05:00") {
		t.Fatal("refresh did not finish")
	}
	click("Open usage page")
	if urls := tt.OpenedURLs(); len(urls) != 1 || urls[0] != chatDemoUsagePage {
		t.Fatalf("opened: %v", urls)
	}
	tt.Key(0, ui.KeyEscape)
	click("Context window")
	if !s.chat.meter.Open || !tt.HasText("25k / 1.0M (2%)") {
		t.Fatal("context meter panel missing")
	}
}

func TestChatStreamsReply(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) { controlChat(c, s, "Default") }, 700, 900)
	if err := tt.Click("Message"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Continue")
	tt.Key(0, ui.KeyEnter)
	if len(s.chat.turns) != 2 || s.chat.draft != "" {
		t.Fatalf("send: turns=%d draft=%q", len(s.chat.turns), s.chat.draft)
	}
	turn := s.chat.turns[1]
	if turn.status != chat.Running || turn.prompt != "Continue" || !tt.HasText("Streaming a simulated reply") {
		t.Fatal("simulated run not started")
	}
	// Move the simulated start back instead of rendering every burst.
	at := func(elapsed time.Duration) {
		turn.started = tt.Now().Add(-elapsed)
		tt.Frame()
	}
	at(time.Second)
	if !tt.HasText("read_source component.go") {
		t.Fatal("live group title missing")
	}
	at(3 * time.Second)
	if turn.received == 0 || turn.received == len(turn.answer) || turn.status != chat.Running {
		t.Fatalf("reply not streaming: %d/%d", turn.received, len(turn.answer))
	}
	at(time.Minute)
	tt.Advance(time.Second) // let the reveal catch up
	if turn.status != chat.Completed || !tt.HasText("The application owns the request") {
		t.Fatalf("reply did not complete: status=%v %d/%d", turn.status, turn.received, len(turn.answer))
	}

	// The titlebar action starts another reply, which Stop cancels.
	toolbar := ui.NewTester(func(c *ui.Context) { chatToolbar(c, s) }, 100, 40)
	if err := toolbar.Click("Stream a reply"); err != nil {
		t.Fatal(err)
	}
	// Each tester keeps its own clock.
	s.chat.turns[len(s.chat.turns)-1].started = tt.Now()
	tt.Frame()
	if err := tt.Click("Stop"); err != nil {
		t.Fatal(err)
	}
	if last := s.chat.turns[len(s.chat.turns)-1]; len(s.chat.turns) != 3 || last.status != chat.Cancelled || last.prompt != chatStreamPrompt {
		t.Fatal("titlebar stream or stop missing")
	}
}
