package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
	"github.com/oneo-me/onevra/source/components/markdown"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

const chatDemoAnswer = "The review is complete.\n\n- **Two tools** ran concurrently and completed in a different order.\n- Execution details remain available above.\n- The final answer stays visible after the process folds.\n\n```go\n// Keep parsed documents outside the frame builder.\ndocument.Render(c)\n```"

type chatDemoTurn struct {
	id, prompt     string
	state          chat.TurnState
	status         chat.Status
	step           int
	answer         *markdown.MarkdownDocument
	started, ended time.Time
	activities     [3]chat.ActivityState
	group          chat.GroupState
	summary        chat.SummaryState
	review         *markdown.MarkdownDocument
	showReview     bool
	attached       bool
}

type chatDemoState struct {
	conversation  chat.ConversationState
	layout        chat.FloatingConversationState
	turns         []*chatDemoTurn
	draft         string
	model, effort string
	attachment    bool
}

func newChatDemoState(now time.Time) *chatDemoState {
	s := &chatDemoState{
		model: "Demo model", effort: "Medium",
		turns: []*chatDemoTurn{{
			id: "sample", prompt: "Review this component and explain the results.",
			step: 20, status: chat.Completed, answer: markdown.ParseMarkdown(chatDemoAnswer),
			started: now.Add(-561 * time.Second), ended: now,
			review:     markdown.ParseMarkdown("```diff\n- renderEveryBlock()\n+ renderVisibleBlocks()\n```"),
			activities: [3]chat.ActivityState{{Expanded: true}, {}, {}},
		}},
	}
	s.turns[0].group.SetExpanded(true)
	return s
}

func chatStatusLabel(c *ui.Context, status chat.Status) string {
	switch status {
	case chat.Completed:
		return tr(c, "gallery.chat_completed")
	case chat.Failed:
		return tr(c, "gallery.chat_failed")
	case chat.Cancelled:
		return tr(c, "gallery.chat_cancelled")
	case chat.Waiting:
		return tr(c, "gallery.chat_waiting")
	default:
		return tr(c, "gallery.chat_running")
	}
}

func controlChat(c *ui.Context, state *controlState, _ string) {
	if state.chat == nil {
		state.chat = newChatDemoState(c.Now())
	}
	s := state.chat
	t := c.Theme()
	ui.Column(c).Fill().MinWidth(0).Gap(t.Space(3)).Children(func() {
		ui.Row(c).FillWidth().Shrink(0).Wrap().AlignItems(ui.Center).Gap(t.Space(2)).Padding(0, t.Space(4)).Children(func() {
			components.Text(c, tr(c, "gallery.chat")).Bold().Grow(1)
		})
		chat.FloatingConversation(c, &s.layout, &s.conversation, tr(c, "gallery.chat_latest"), func(c *ui.Context) {
			for _, turn := range s.turns {
				toolA, toolB := chat.Running, chat.Running
				if turn.step >= 3 {
					toolA = chat.Completed
				}
				if turn.step >= 2 {
					toolB = chat.Completed
				}
				if turn.status == chat.Failed || turn.status == chat.Cancelled {
					if toolA != chat.Completed {
						toolA = turn.status
					}
					if toolB != chat.Completed {
						toolB = turn.status
					}
				}
				steps := []chat.Step{{ID: "work", Title: tr(c, "gallery.chat_reasoning"), Group: &turn.group, Activities: []chat.Activity{
					{ID: "summary", State: &turn.activities[0], Icon: lucide.Brain, Action: tr(c, "gallery.chat_think"), Status: chat.Completed,
						Content: func(c *ui.Context) {
							components.Text(c, "I will inspect the source and check the examples in parallel.").Selectable()
						}},
					{ID: "call-source", State: &turn.activities[1], Icon: lucide.FileText, Action: "read_source", Detail: "component.go", MaxHeight: 120, Mono: true, Status: toolA,
						Content: func(c *ui.Context) {
							components.Text(c, fmt.Sprintf("component.go · %d / 3 chunks", min(turn.step, 3))).Selectable()
						}},
					{ID: "call-examples", State: &turn.activities[2], Icon: lucide.Terminal, Action: "check_examples", Detail: "examples/", MaxHeight: 100, Mono: true, Status: toolB,
						Content: func(c *ui.Context) {
							if turn.step >= 2 {
								components.Text(c, "3 examples checked. No errors.").Selectable()
							}
						}},
				}}}
				var final chat.View
				if turn.answer != nil {
					final = func(c *ui.Context) {
						turn.answer.Render(c, markdown.MarkdownOptions{CopyCodeLabel: tr(c, "gallery.copy_code")}).FillWidth()
					}
				}
				chat.Turn(c, &turn.state, chat.TurnOptions{
					ID: turn.id, Summary: chatTurnSummary(c, turn), Status: turn.status,
					Prompt: func(c *ui.Context) {
						if turn.attached {
							components.Text(c, "context.md").TextColor(c.Theme().TextMuted)
						}
						components.Text(c, turn.prompt).Selectable()
					},
					Final: final, Steps: steps, Working: chatTurnSummary(c, turn),
					Footer:   func(c *ui.Context) { chatDemoResults(c, turn) },
					HoldOpen: !s.conversation.AtLatest(),
				})
			}
		}, func(c *ui.Context) {
			chat.GlassComposer(c, &s.layout, &s.draft, chat.ComposerOptions{
				Label: tr(c, "gallery.chat_message"), Placeholder: tr(c, "gallery.chat_placeholder"),
				SendLabel: tr(c, "gallery.chat_send"), StopLabel: tr(c, "gallery.chat_stop"),
				Accessory: func(c *ui.Context) { chatDemoToolbar(c, s) },
				Trailing: func(c *ui.Context) {
					if chat.IconAction(c, tr(c, "gallery.chat_attach"), lucide.Paperclip).Clicked() {
						s.attachment = true
					}
				},
				Attachments: func(c *ui.Context) {
					if s.attachment {
						components.Text(c, "context.md").TextColor(t.TextMuted)
						if chat.IconAction(c, tr(c, "gallery.chat_remove_attachment"), lucide.X).Clicked() {
							s.attachment = false
						}
					}
				},
				Footer: func(c *ui.Context) {
					components.Icon(c, lucide.Folder).Size(t.Space(4), t.Space(4)).TextColor(t.TextMuted)
					components.Text(c, "onevra").TextColor(t.TextMuted)
					components.Text(c, tr(c, "gallery.chat_demo_context")).TextColor(t.TextMuted)
				},
			})
		})
	})
}
