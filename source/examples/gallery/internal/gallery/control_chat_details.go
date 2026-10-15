package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
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

var chatDemoArtifacts = []chat.SummaryItem{
	{ID: "component", Label: "component.go", Metadata: chatDemoMetrics(32, 8)},
	{ID: "tests", Label: "component_test.go", Metadata: chatDemoMetrics(24, 0)},
	{ID: "example", Label: "examples/demo.go", Metadata: chatDemoMetrics(12, 2)},
	{ID: "readme", Label: "README.md", Metadata: chatDemoMetrics(4, 1)},
	{ID: "theme", Label: "theme.go", Metadata: chatDemoMetrics(3, 3)},
	{ID: "styles", Label: "styles.go", Metadata: chatDemoMetrics(7, 2)},
}

func chatDemoResults(c *ui.Context, s *chatDemoState, turn *chatDemoTurn) {
	if turn.status != chat.Completed {
		return
	}
	result := chat.ResultSummary(c, &turn.summary, chat.SummaryOptions{
		Title: tr(c, "chatChanged"), Metadata: chatDemoMetrics(82, 16), ActionLabel: tr(c, "chatReview"),
		MoreLabel: tr(c, "chatMoreFiles"), LessLabel: tr(c, "chatLessFiles"), Items: chatDemoArtifacts,
	})
	if result.Action || result.SelectedID != "" {
		turn.showReview = true
	}
	if turn.showReview {
		ui.Column(c).FillWidth().MinWidth(0).Gap(c.Theme().Space(2)).Children(func() {
			components.Text(c, tr(c, "chatSampleDiff")).TextColor(c.Theme().TextMuted)
			markdown.View(c, turn.review, &s.markdown).FillWidth()
			if components.Button(c, tr(c, "chatCloseReview"), components.ButtonGhost).Clicked() {
				turn.showReview = false
			}
		})
	}
}

func chatDemoToolbar(c *ui.Context, s *chatDemoState) {
	ui.Box(c).Key("model").Children(func() {
		selectModel := components.Select(c, &s.model)
		selectModel.Trigger.Label(tr(c, "chatModel")).Children(func() { components.Text(c, s.model).SingleLine() })
		selectModel.Popup(func(ui.Element) {
			for _, name := range []string{"Demo model", "Fast demo"} {
				selectModel.Item(name).Children(func() { components.Text(c, name) })
			}
		})
	})
	ui.Box(c).Key("effort").Children(func() {
		selectEffort := components.Select(c, &s.effort)
		selectEffort.Trigger.Label(tr(c, "chatEffort")).Children(func() { components.Text(c, s.effort).SingleLine() })
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
