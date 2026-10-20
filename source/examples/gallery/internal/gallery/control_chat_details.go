package gallery

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func chatTurnSummary(c *ui.Context, turn *chatDemoTurn) string {
	end := turn.ended
	if end.IsZero() {
		end = c.Now()
	}
	seconds := max(0, int(end.Sub(turn.started).Seconds()))
	if turn.status == chat.Completed {
		return msg(c).ChatWorked(seconds/60, seconds%60)
	}
	return msg(c).ChatWorking(chatStatusLabel(c, turn.status), seconds/60, seconds%60)
}

type chatDemoModel struct {
	id, name, provider string
	window             int
}

var chatDemoModels = []chatDemoModel{
	{"demo", "Demo model", "chatProviderLocal", 200_000},
	{"fast", "Fast demo", "chatProviderLocal", 128_000},
	{"deep", "Deep demo", "chatProviderRemote", 1_000_000},
}

var chatDemoFiles = []chatDemoFile{{name: "context.md", sample: true}, {name: "screenshot.png", image: true, sample: true}, {name: "design.pdf", sample: true}}

func chatDemoModelFor(id string) chatDemoModel {
	for _, m := range chatDemoModels {
		if m.id == id {
			return m
		}
	}
	return chatDemoModels[0]
}

// chatDemoAccessory holds the model and reasoning effort pickers inside the card.
func chatDemoAccessory(c *ui.Context, s *chatDemoState) {
	model := chatDemoModelFor(s.model)
	models := make([]chat.Option, len(chatDemoModels))
	for i, m := range chatDemoModels {
		models[i] = chat.Option{ID: m.id, Label: m.name, Detail: chatDemoTokens(m.window), Group: tr(c, m.provider), Icon: lucide.Sparkles}
	}
	if r := chat.Picker(c, &s.modelPicker, chat.PickerOptions{Label: model.name, Icon: lucide.Sparkles, Options: models, Selected: s.model,
		SearchLabel: tr(c, "chatSearchModels"), EmptyText: tr(c, "chatNoResults"), Searchable: 2, Width: 280}); r.Chosen != "" {
		s.model = r.Chosen
	}
	efforts := []chat.Option{{ID: "off", Label: tr(c, "chatEffortOff")}, {ID: "low", Label: tr(c, "chatEffortLow")},
		{ID: "medium", Label: tr(c, "chatEffortMedium")}, {ID: "high", Label: tr(c, "chatEffortHigh")}, {ID: "max", Label: tr(c, "chatEffortMax")}}
	label := tr(c, "chatEffort")
	tooltip := label
	for _, o := range efforts {
		if o.ID == s.effort {
			tooltip = label + " · " + o.Label
		}
	}
	if r := chat.EffortPicker(c, &s.effortPicker, chat.EffortPickerOptions{Label: label, Tooltip: tooltip, Icon: lucide.Brain,
		Levels: efforts, Selected: s.effort, Off: true}); r.Chosen != "" {
		s.effort = r.Chosen
	}
}

// chatDemoTrailing attaches sample files and shows simulated context use.
func chatDemoTrailing(c *ui.Context, s *chatDemoState) {
	// Each sample attaches once; the action rests when all are attached.
	next := -1
	for i, f := range chatDemoFiles {
		if next < 0 && !slices.Contains(s.files, f) {
			next = i
		}
	}
	if attach := chat.IconAction(c, tr(c, "chatAttach"), lucide.Paperclip).Disabled(next < 0); attach.Clicked() && next >= 0 {
		s.files = append(s.files, chatDemoFiles[next])
	}
	window := chatDemoModelFor(s.model).window
	used := 18_400 + 6_200*len(s.turns)
	if turn := s.running(); turn != nil {
		used += turn.received / 4
	}
	fraction := math.Min(1, float64(used)/float64(window))
	value := fmt.Sprintf("%s / %s (%d%%)", chatDemoTokens(used), chatDemoTokens(window), int(math.Round(fraction*100)))
	chat.ContextMeter(c, &s.meter, chat.ContextMeterOptions{Label: tr(c, "chatContextWindow"), Tooltip: value,
		Title: tr(c, "chatContextWindow"), Value: value, Fraction: fraction})
}

// chatDemoFooter shows the project, branch and working mode under the card.
func chatDemoFooter(c *ui.Context, s *chatDemoState) {
	t := c.Theme()
	chat.Chip(c, "onevra", chat.ChipOptions{Icon: lucide.Folder, ReadOnly: true})
	var branches []chat.Option
	for _, b := range []string{"main", "feature/chat-ui", "release/0.2026.10"} {
		branches = append(branches, chat.Option{ID: b, Label: b})
	}
	if r := chat.Picker(c, &s.branchPicker, chat.PickerOptions{Label: s.branch, Icon: lucide.GitBranch, Options: branches, Selected: s.branch,
		SearchLabel: tr(c, "chatSearchBranches"), EmptyText: tr(c, "chatNoResults"), Searchable: 2}); r.Chosen != "" {
		s.branch = r.Chosen
	}
	// The label matches the chips beside it.
	components.Switch(c, &s.worktree, tr(c, "chatWorktree"), components.SwitchOptions{Compact: true}).
		AlignItems(ui.Center).Margin(0, 0, 0, 7).FontSize(12.5).TextColor(t.Text.Mix(t.TextMuted, .45))
	ui.Spacer(c)
	if s.running() != nil {
		chat.Shimmer(c, tr(c, "chatStreaming"), 12.5, t.TextMuted, t.Text)
	}
	quota := chatDemoQuota(c, s)
	row := quota.Rows[0]
	r := chat.QuotaSummary(c, &s.quota, chat.QuotaSummaryOptions{Label: tr(c, "chatUsageLimits"),
		Tooltip: quota.Title + "\n" + row.Label + " · " + row.Value, Row: row, Quotas: []chat.QuotaOptions{quota}})
	if r.Refresh != "" {
		// A host would start the read off the UI thread; the sample only waits.
		s.quotaRefreshed = c.Now().Add(chatQuotaRead)
	}
	if r.Open != "" {
		c.OpenURL(chatDemoUsagePage)
	}
}

const (
	chatQuotaEvery    = 5 * time.Minute
	chatQuotaRead     = 1200 * time.Millisecond
	chatDemoUsagePage = "https://github.com/oneo-me/onevra"
)

// chatDemoQuota is a sample account; the summary glances at its shortest
// limit. Its countdown and simulated refresh follow the frame clock.
func chatDemoQuota(c *ui.Context, s *chatDemoState) chat.QuotaOptions {
	now := c.Now()
	if !s.quotaRefreshed.IsZero() && !now.Before(s.quotaRefreshed) {
		s.quotaRefreshed, s.quotaNext = time.Time{}, time.Time{}
	}
	if s.quotaNext.IsZero() || now.After(s.quotaNext) {
		s.quotaNext = now.Add(chatQuotaEvery)
	}
	refreshing := !s.quotaRefreshed.IsZero()
	countdown := tr(c, "chatQuotaRefreshing")
	if !refreshing {
		// Round up so the countdown never shows zero before a read.
		seconds := int((s.quotaNext.Sub(now) + time.Second - 1) / time.Second)
		countdown = fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
	}
	if s.quota.Open || refreshing {
		c.After(time.Second / 4)
	}
	return chat.QuotaOptions{
		ID: "demo", Title: tr(c, "chatQuotaAccount"), Countdown: countdown, Refreshing: refreshing,
		RefreshLabel: tr(c, "chatQuotaRefresh"), OpenLabel: tr(c, "chatQuotaOpen"),
		Leading: func(c *ui.Context, size float32) {
			components.Icon(c, lucide.Gauge).Size(size, size).Shrink(0).TextColor(c.Theme().TextMuted)
		},
		Rows: []chat.QuotaRow{
			{Label: tr(c, "chatQuotaFiveHour"), Detail: tr(c, "chatQuotaFiveHourReset"), Value: msg(c).ChatQuotaLeft(63), Fraction: .63},
			{Label: tr(c, "chatQuotaWeekly"), Detail: tr(c, "chatQuotaWeeklyReset"), Value: msg(c).ChatQuotaLeft(18), Fraction: .18},
			{Label: tr(c, "chatQuotaCredits"), Value: "$12.40", Fraction: -1},
		},
	}
}

func chatDemoTokens(n int) string {
	switch {
	case n >= 999_500:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1000:
		return fmt.Sprintf("%.0fk", float64(n)/1e3)
	}
	return fmt.Sprint(n)
}
