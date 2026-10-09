package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"testing"
)

func TestIndependentLazyToolDetails(t *testing.T) {
	var a, b BlockState
	builds := [2]int{}
	var turn TurnState
	opts := TurnOptions{ID: "turn", Summary: "Work", Blocks: []Block{
		{ID: "a", Title: "Read", State: &a, Content: func(c *ui.Context) { builds[0]++; components.Text(c, "Read output") }},
		{ID: "b", Title: "Run", State: &b, DetailHeight: 100, Content: func(c *ui.Context) { builds[1]++; ui.Box(c).Height(900).Shrink(0) }},
	}}
	tt := ui.NewTester(func(c *ui.Context) { Turn(c, &turn, opts) }, 400, 400)
	if builds != [2]int{} {
		t.Fatal("closed tool content built")
	}
	if err := tt.Click("Run"); err != nil {
		t.Fatal(err)
	}
	if a.Expanded || !b.Expanded || builds[0] != 0 || b.Scroll.MaxY <= 0 {
		t.Fatal("independent bounded tool view failed")
	}
	turn.SetExpanded(false)
	builds = [2]int{}
	tt.Frame()
	if builds != [2]int{} {
		t.Fatal("outer fold built tools")
	}
	turn.SetExpanded(true)
	tt.Frame()
	if !b.Expanded || builds[1] == 0 {
		t.Fatal("outer fold discarded tool state")
	}
}

func TestResultSummaryLimitsAndActions(t *testing.T) {
	var state SummaryState
	items := make([]SummaryItem, 1000)
	for i := range items {
		items[i] = SummaryItem{ID: fmt.Sprint(i), Label: fmt.Sprintf("Artifact %d", i)}
	}
	reviews := 0
	selected := ""
	tt := ui.NewTester(func(c *ui.Context) {
		p := ResultSummary(c, &state, SummaryOptions{Title: "Results", ActionLabel: "Review", MoreLabel: "Show first 100", LessLabel: "Less", Items: items})
		if p.Action {
			reviews++
		}
		if p.SelectedID != "" {
			selected = p.SelectedID
		}
	}, 400, 5000)
	if tt.HasText("Artifact 3") {
		t.Fatal("preview exceeds three items")
	}
	tt.Click("Artifact 1")
	tt.Click("Review")
	if selected != "1" || reviews != 1 {
		t.Fatal("result actions missing")
	}
	tt.Click("Show first 100")
	if !tt.HasText("Artifact 99") || tt.HasText("Artifact 100") {
		t.Fatal("expanded result cap failed")
	}
}

func TestTurnFooterSurvivesFoldWithoutAnswer(t *testing.T) {
	var state TurnState
	state.SetExpanded(false)
	tt := ui.NewTester(func(c *ui.Context) {
		Turn(c, &state, TurnOptions{ID: "turn", Summary: "Stopped", Status: Cancelled,
			Blocks:  []Block{{ID: "tool", Content: func(*ui.Context) { t.Fatal("hidden tool built") }}},
			Footer:  func(c *ui.Context) { components.Text(c, "Artifacts") },
			Actions: func(c *ui.Context) { components.Button(c, "Retry", components.ButtonGhost) },
		})
	}, 400, 300)
	if !tt.HasText("Artifacts") || !tt.HasText("Retry") {
		t.Fatal("fold hid footer")
	}
}
