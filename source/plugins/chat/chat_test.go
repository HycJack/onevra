package chat

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

// group wraps activities in one step for tests.
func group(title string, activities ...Activity) []Step {
	return []Step{{ID: "group", Title: title, Activities: activities, Group: &GroupState{}}}
}

func TestTurnCompletionAndUserOverride(t *testing.T) {
	var state TurnState
	builds := 0
	opts := TurnOptions{ID: "turn", Summary: "Worked for 1m", Working: "Working", Status: Running,
		Steps: []Step{{ID: "m", Message: func(c *ui.Context) { builds++; components.Text(c, "Tool output").Selectable() }}}}
	tt := ui.NewTester(func(c *ui.Context) { Turn(c, &state, opts) }, 400, 400)
	if !tt.HasText("Tool output") || !tt.HasText("Working") {
		t.Fatal("running content missing")
	}
	opts.Status = Completed
	tt.Frame()
	if !tt.HasText("Tool output") || tt.HasText("Working") {
		t.Fatal("completion without final hid the only answer")
	}
	opts.Final = func(c *ui.Context) { components.Text(c, "Final response") }
	builds = 0
	tt.Frame()
	if builds != 0 || tt.HasText("Tool output") || !tt.HasText("Final response") {
		t.Fatal("completed steps were not lazy")
	}
	if err := tt.Click("Worked for 1m"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Tool output") {
		t.Fatal("cannot reopen completion")
	}
	tt.Frame()
	if !tt.HasText("Tool output") {
		t.Fatal("stream rebuild reset manual expansion")
	}
	tt.Key(0, ui.KeySpace)
	builds = 0
	tt.Frame()
	if builds != 0 {
		t.Fatal("closing still built hidden content")
	}
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Tool output") {
		t.Fatal("keyboard expansion failed")
	}
}

func TestTurnProtectsReadingAndExceptionalStates(t *testing.T) {
	var folded TurnState
	if folded.Expanded(Completed, true, false) || folded.Expanded(Completed, true, true) {
		t.Fatal("scrolling through history reopened a completed turn")
	}
	for _, status := range []Status{Running, Waiting, Failed, Cancelled} {
		var s TurnState
		if !s.Expanded(status, true, false) {
			t.Fatalf("status %v folded", status)
		}
	}
	var state TurnState
	opts := TurnOptions{ID: "turn", Summary: "Worked", Status: Running,
		Steps: []Step{{ID: "m", Message: func(c *ui.Context) { components.Button(c, "Review", components.ButtonSecondary) }}},
		Final: func(c *ui.Context) { components.Text(c, "Answer") }}
	tt := ui.NewTester(func(c *ui.Context) { Turn(c, &state, opts) }, 400, 400)
	if err := tt.Click("Review"); err != nil {
		t.Fatal(err)
	}
	opts.Status = Completed
	tt.Frame()
	if !tt.HasText("Review") || !tt.Focused("Review") {
		t.Fatal("completion removed focused content")
	}
	state.SetExpanded(false)
	opts.HoldOpen = true
	tt.Frame()
	if tt.HasText("Review") {
		t.Fatal("reading policy overrode explicit collapse")
	}
	state.ResetExpansion()
	tt.Frame()
	if !tt.HasText("Review") {
		t.Fatal("hold did not defer auto collapse")
	}
}

func TestActivityIdentityAcrossOutOfOrderResults(t *testing.T) {
	var state TurnState
	draftA, draftB := "A", "B"
	sa, sb := &ActivityState{Expanded: true}, &ActivityState{Expanded: true}
	a := Activity{ID: "call-a", Action: "Run", Detail: "a", Status: Running, State: sa, Content: func(c *ui.Context) { components.TextBox(c, &draftA, components.TextBoxOptions{}).Label("Input A") }}
	b := Activity{ID: "call-b", Action: "Run", Detail: "b", Status: Running, State: sb, Content: func(c *ui.Context) { components.TextBox(c, &draftB, components.TextBoxOptions{}).Label("Input B") }}
	opts := TurnOptions{ID: "turn", Status: Running, Steps: group("Running 2 commands", a, b)}
	tt := ui.NewTester(func(c *ui.Context) { Turn(c, &state, opts) }, 400, 500)
	if err := tt.Click("Input B"); err != nil {
		t.Fatal(err)
	}
	b.Status = Completed
	opts.Steps = group("Running 2 commands", b, a)
	tt.Frame()
	if !tt.Focused("Input B") {
		t.Fatal("reordering reassigned tool focus")
	}
	tt.Type("!")
	if draftA != "A" || draftB == "B" {
		t.Fatal("update targeted wrong tool")
	}
}

func TestConversationFollowAndJump(t *testing.T) {
	var state ConversationState
	height := float32(800)
	tt := ui.NewTester(func(c *ui.Context) {
		Conversation(c, &state, "Latest", func(c *ui.Context) { ui.Box(c).Height(height).Shrink(0) })
	}, 300, 300)
	if state.Scroll.MaxY <= 0 || !state.AtLatest() {
		t.Fatalf("initial position %+v", state.Scroll)
	}
	height += 100
	tt.Frame()
	if !state.AtLatest() {
		t.Fatal("growth lost bottom")
	}
	tt.Scroll(100, 100, 0, -180)
	before := state.Scroll.Y
	if state.AtLatest() {
		t.Fatal("wheel did not release follow")
	}
	height += 100
	tt.Frame()
	if state.Scroll.Y != before {
		t.Fatalf("reading moved from %v to %v", before, state.Scroll.Y)
	}
	if err := tt.Click("Latest"); err != nil {
		t.Fatal(err)
	}
	if !state.AtLatest() {
		t.Fatal("jump did not resume follow")
	}
	height = 400
	tt.Frame()
	if !state.AtLatest() {
		t.Fatal("fold lost bottom")
	}
}

func TestConversationExpansionKeepsHeaderAndComposer(t *testing.T) {
	var scroll ConversationState
	var state TurnState
	draft := ""
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			Conversation(c, &scroll, "Latest", func(c *ui.Context) {
				Turn(c, &state, TurnOptions{ID: "turn", Summary: "Details", Status: Completed,
					Steps: []Step{{ID: "large", Message: func(c *ui.Context) { ui.Box(c).Height(1500).Shrink(0) }}},
					Final: func(c *ui.Context) { components.Text(c, "Answer") }})
			})
			Composer(c, &draft, ComposerOptions{Label: "Draft", SendLabel: "Send"})
		})
	}, 360, 500)
	before, _ := tt.Find("Details")
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	after, _ := tt.Find("Details")
	if before.Y != after.Y {
		t.Fatalf("disclosure moved from %v to %v; scroll=%+v texts=%v", before.Y, after.Y, scroll.Scroll, tt.Texts())
	}
	send, ok := tt.Find("Send")
	if !ok || send.Y+send.H > 500 {
		t.Fatalf("composer pushed out of viewport: %+v", send)
	}
	if scroll.Scroll.MaxY <= 0 || scroll.AtLatest() {
		t.Fatal("expansion did not pause follow")
	}
}

func BenchmarkCompletedTurn(b *testing.B) {
	for _, count := range []int{10, 10000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			var state TurnState
			steps := make([]Step, count)
			for i := range steps {
				steps[i] = Step{ID: fmt.Sprint(i), Message: func(*ui.Context) { b.Fatal("hidden content built") }}
			}
			tt := ui.NewTester(func(c *ui.Context) {
				Turn(c, &state, TurnOptions{ID: "turn", Summary: "Completed", Status: Completed, Steps: steps,
					Final: func(c *ui.Context) { components.Text(c, "Answer") }})
			}, 400, 300)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				tt.Frame()
			}
		})
	}
}
