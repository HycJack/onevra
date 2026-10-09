package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/plugins/chat"
)

func chatTurnSummary(c *ui.Context, turn *chatDemoTurn) string {
	end := turn.ended
	if end.IsZero() {
		end = c.Now()
	}
	seconds := max(0, int(end.Sub(turn.started).Seconds()))
	if turn.status == chat.Completed {
		return fmt.Sprintf(tr(c, "gallery.chat_worked"), seconds/60, seconds%60)
	}
	return fmt.Sprintf(tr(c, "gallery.chat_working"), chatStatusLabel(c, turn.status), seconds/60, seconds%60)
}

var chatDemoArtifacts = []chat.SummaryItem{
	{ID: "component", Label: "component.go", Metadata: chatDemoMetrics(32, 8)},
	{ID: "tests", Label: "component_test.go", Metadata: chatDemoMetrics(24, 0)},
	{ID: "example", Label: "examples/demo.go", Metadata: chatDemoMetrics(12, 2)},
	{ID: "readme", Label: "README.md", Metadata: chatDemoMetrics(4, 1)},
	{ID: "theme", Label: "theme.go", Metadata: chatDemoMetrics(3, 3)},
	{ID: "styles", Label: "styles.go", Metadata: chatDemoMetrics(7, 2)},
}

func chatDemoResults(c *ui.Context, turn *chatDemoTurn) {
	if turn.status != chat.Completed {
		return
	}
	result := chat.ResultSummary(c, &turn.summary, chat.SummaryOptions{
		Title: tr(c, "gallery.chat_changed"), Metadata: chatDemoMetrics(82, 16), ActionLabel: tr(c, "gallery.chat_review"),
		MoreLabel: tr(c, "gallery.chat_more_files"), LessLabel: tr(c, "gallery.chat_less_files"), Items: chatDemoArtifacts,
	})
	if result.Action || result.SelectedID != "" {
		turn.showReview = true
	}
	if turn.showReview {
		ui.Column(c).FillWidth().MinWidth(0).Gap(c.Theme().Space(2)).Children(func() {
			components.Text(c, tr(c, "gallery.chat_sample_diff")).TextColor(c.Theme().TextMuted)
			turn.review.Render(c).FillWidth()
			if components.Button(c, tr(c, "gallery.chat_close_review"), components.ButtonGhost).Clicked() {
				turn.showReview = false
			}
		})
	}
}

func chatDemoToolbar(c *ui.Context, s *chatDemoState) {
	if chat.IconAction(c, tr(c, "gallery.chat_attach"), lucide.Paperclip).Clicked() {
		s.attachment = true
	}
	ui.Box(c).Key("model").Children(func() {
		selectModel := components.Select(c, &s.model)
		selectModel.Trigger.Label(tr(c, "gallery.chat_model")).Children(func() { components.Text(c, s.model).SingleLine() })
		selectModel.Popup(func(ui.Element) {
			for _, name := range []string{"Demo model", "Fast demo"} {
				selectModel.Item(name).Children(func() { components.Text(c, name) })
			}
		})
	})
	ui.Box(c).Key("effort").Children(func() {
		selectEffort := components.Select(c, &s.effort)
		selectEffort.Trigger.Label(tr(c, "gallery.chat_effort")).Children(func() { components.Text(c, s.effort).SingleLine() })
		selectEffort.Popup(func(ui.Element) {
			for _, name := range []string{"Low", "Medium", "High"} {
				selectEffort.Item(name).Children(func() { components.Text(c, name) })
			}
		})
	})
}

func chatDemoMetrics(added, removed int) chat.View {
	return func(c *ui.Context) {
		components.Text(c, fmt.Sprintf("+%d", added)).TextColor(c.Theme().Success)
		components.Text(c, fmt.Sprintf("−%d", removed)).TextColor(c.Theme().Danger)
	}
}
