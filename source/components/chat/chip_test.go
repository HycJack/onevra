package chat

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func TestChipAndSwitchChip(t *testing.T) {
	clicks, on := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Children(func() {
			if Chip(c, "Branch", ChipOptions{Icon: lucide.GitBranch, Caret: true}).Clicked() {
				clicks++
			}
			Chip(c, "Project", ChipOptions{ReadOnly: true})
			SwitchChip(c, &on, "Worktree")
		})
	}, 500, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if tt.Focused("Project") {
		t.Fatal("read-only chip took focus")
	}
	for _, label := range []string{"Branch", "Worktree"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	if clicks != 1 || !on {
		t.Fatalf("clicks=%d on=%v", clicks, on)
	}
}

func TestPickerGroupsSearchAndChoice(t *testing.T) {
	var state PickerState
	selected, chosen := "b", ""
	options := []Option{{ID: "a", Label: "Alpha", Group: "One"}, {ID: "b", Label: "Beta", Group: "One", Detail: "200k"}}
	for i := range 8 {
		options = append(options, Option{ID: fmt.Sprint("x", i), Label: fmt.Sprint("Extra ", i), Group: "Two"})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		r := Picker(c, &state, PickerOptions{Label: selected, SearchLabel: "Search", EmptyText: "Nothing", Options: options, Selected: selected})
		if r.Chosen != "" {
			chosen, selected = r.Chosen, r.Chosen
		}
	}, 500, 600)
	if err := tt.Click("b"); err != nil || !state.Open {
		t.Fatalf("open: %v %v", err, state.Open)
	}
	if !tt.HasText("One") || !tt.HasText("Two") || !tt.HasText("200k") {
		t.Fatal("groups or detail missing")
	}
	if !tt.Focused("Search") {
		t.Fatal("search not focused")
	}
	tt.Type("zzz")
	if !tt.HasText("Nothing") || tt.HasText("Alpha") {
		t.Fatal("empty search result missing")
	}
	tt.Command("selectAll")
	tt.Type("alp")
	if tt.HasText("Beta") {
		t.Fatal("search did not filter")
	}
	if err := tt.Click("Alpha"); err != nil || chosen != "a" || state.Open {
		t.Fatalf("choice: %v chosen=%q open=%v", err, chosen, state.Open)
	}
}

func TestContextMeterPanel(t *testing.T) {
	var state MeterState
	tt := ui.NewTester(func(c *ui.Context) {
		ContextMeter(c, &state, ContextMeterOptions{Label: "Context", Title: "Context window", Value: "61k / 272k (22%)", Fraction: .22})
	}, 500, 300)
	if tt.HasText("61k") {
		t.Fatal("panel shown before opening")
	}
	if err := tt.Click("Context"); err != nil || !state.Open || !tt.HasText("61k / 272k (22%)") {
		t.Fatalf("panel: %v open=%v", err, state.Open)
	}
}

func TestQuotaSummaryPanel(t *testing.T) {
	var state MeterState
	var refresh, open string
	quota := QuotaOptions{ID: "acct", Title: "Account", Countdown: "04:59", RefreshLabel: "Refresh", OpenLabel: "Open page", Rows: []QuotaRow{
		{Label: "5-hour limit", Detail: "2h 14m", Value: "63% left", Fraction: .63},
		{Label: "Credits", Value: "$12.40", Fraction: -1},
	}}
	row := quota.Rows[0]
	tt := ui.NewTester(func(c *ui.Context) {
		r := QuotaSummary(c, &state, QuotaSummaryOptions{Label: "Usage limits", Row: row, Quotas: []QuotaOptions{quota}})
		if r.Refresh != "" {
			refresh = r.Refresh
		}
		if r.Open != "" {
			open = r.Open
		}
	}, 500, 400)
	if !tt.HasText("2h 14m") || tt.HasText("63% left") {
		t.Fatal("summary should show only the reset time and bar")
	}
	if err := tt.Click("Usage limits"); err != nil || !state.Open {
		t.Fatalf("open: %v", err)
	}
	for _, text := range []string{"Account", "04:59", "5-hour limit", "63% left", "$12.40"} {
		if !tt.HasText(text) {
			t.Fatalf("panel missing %q", text)
		}
	}
	if err := tt.Click("Refresh"); err != nil || refresh != "acct" {
		t.Fatalf("refresh: %v %q", err, refresh)
	}
	if err := tt.Click("Open page"); err != nil || open != "acct" {
		t.Fatalf("open: %v %q", err, open)
	}
	quota.Refreshing, refresh = true, ""
	tt.Frame()
	_ = tt.Click("Refresh")
	if refresh != "" {
		t.Fatal("refresh while refreshing")
	}
	state.Open, row = false, quota.Rows[1]
	tt.Frame()
	if !tt.HasText("$12.40") {
		t.Fatal("balance summary should show its value")
	}
}
