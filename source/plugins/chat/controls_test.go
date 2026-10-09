package chat

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestComposerSubmitOnEnterAndSteering(t *testing.T) {
	draft := ""
	sends, stops := 0, 0
	running := false
	tt := ui.NewTester(func(c *ui.Context) {
		p := Composer(c, &draft, ComposerOptions{Label: "Message", SendLabel: "Send", StopLabel: "Stop",
			Running: running})
		if p.Send {
			sends++
		}
		if p.Stop {
			stops++
		}
	}, 400, 240)
	tt.Click("Message")
	tt.Key(0, ui.KeyEnter)
	if sends != 0 {
		t.Fatal("empty draft sent")
	}
	tt.Compose("ni", 2)
	tt.Key(0, ui.KeyEnter)
	if sends != 0 {
		t.Fatal("composition sent")
	}
	tt.Type("你")
	lines := strings.Count(draft, "\n")
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Type("b")
	if sends != 0 || strings.Count(draft, "\n") != lines+1 {
		t.Fatalf("shift+enter: sends=%d draft=%q", sends, draft)
	}
	before := draft
	tt.Key(0, ui.KeyEnter)
	if sends != 1 || draft != before {
		t.Fatalf("enter: sends=%d draft=%q", sends, draft)
	}
	running = true
	tt.Frame()
	if !tt.HasText("Stop") && tt.Click("Stop") != nil {
		t.Fatal("stop missing")
	}
	if err := tt.Click("Send"); err != nil || sends != 2 {
		t.Fatalf("steer send unavailable: %v %d", err, sends)
	}
	draft = ""
	tt.Frame()
	if _, ok := tt.Find("Send"); ok {
		t.Fatal("send shown for an empty draft while running")
	}
	if err := tt.Click("Stop"); err != nil || stops != 1 {
		t.Fatal("stop failed")
	}
}

func TestContextMeterPanel(t *testing.T) {
	var state MeterState
	tt := ui.NewTester(func(c *ui.Context) {
		ContextMeter(c, &state, ContextMeterOptions{Label: "Usage", Title: "Context window", Value: "42k / 100k (42%)", Fraction: .42})
	}, 500, 400)
	if tt.HasText("Context window") {
		t.Fatal("panel built while closed")
	}
	if err := tt.Click("Usage"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !state.Open || !tt.HasText("42k / 100k (42%)") {
		t.Fatalf("panel not shown: %v", tt.Texts())
	}
}

func TestAttachmentAndThreadItem(t *testing.T) {
	removed, clicked := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Attachment(c, AttachmentOptions{Label: "notes.md", Detail: "2 KB", RemoveLabel: "Remove notes.md"}).Remove {
			removed++
		}
		Attachment(c, AttachmentOptions{Label: "sent.md"})
		if ThreadItem(c, ThreadItemOptions{ID: "a", Title: "Fix login", Detail: "neorio", Time: "2m", Status: Failed}).Clicked() {
			clicked++
		}
	}, 500, 400)
	if err := tt.Click("Remove notes.md"); err != nil || removed != 1 {
		t.Fatal("remove failed")
	}
	if _, ok := tt.Find("Remove sent.md"); ok {
		t.Fatal("read-only attachment offered removal")
	}
	if err := tt.Click("Fix login"); err != nil || clicked != 1 || !tt.HasText("neorio") || !tt.HasText("2m") {
		t.Fatal("thread item activation failed")
	}
}

func TestPickerGroupsFiltersAndChooses(t *testing.T) {
	var state PickerState
	chosen := ""
	var options []Option
	for i := range 10 {
		options = append(options, Option{ID: fmt.Sprint("a/", i), Label: fmt.Sprint("Alpha ", i), Group: "alpha", Detail: "128k"})
	}
	options = append(options, Option{ID: "b/1", Label: "Beta One", Group: "beta"})
	tt := ui.NewTester(func(c *ui.Context) {
		r := Picker(c, &state, PickerOptions{Label: "Model", SearchLabel: "Search models", EmptyText: "None", Options: options, Selected: "a/1"})
		if r.Chosen != "" {
			chosen = r.Chosen
		}
	}, 600, 800)
	if tt.HasText("Beta One") {
		t.Fatal("closed picker built options")
	}
	if err := tt.Click("Model"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("alpha") || !tt.HasText("Alpha 1") {
		t.Fatalf("groups missing: %v", tt.Texts())
	}
	tt.Type("beta")
	tt.Frame()
	if tt.HasText("Alpha 1") || !tt.HasText("Beta One") {
		t.Fatalf("filter failed: %v", tt.Texts())
	}
	if err := tt.Click("Beta One"); err != nil {
		t.Fatal(err)
	}
	if chosen != "b/1" || state.Open {
		t.Fatalf("choice: %q open=%v", chosen, state.Open)
	}
}
