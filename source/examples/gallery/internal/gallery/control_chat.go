package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/plugins/chat"
	"github.com/oneo-me/onevra/plugins/markdown"
)

const chatDemoAnswer = "The review is complete.\n\n- **Two tools** ran concurrently and completed in a different order.\n- Execution details remain available above.\n- The final answer stays visible after the process folds.\n\n```go\n// Keep parsed documents outside the frame builder.\ndocument.Render(c)\n```"

type chatDemoTurn struct {
	id, prompt         string
	state              chat.TurnState
	status             chat.Status
	step               int
	answer             *markdown.MarkdownDocument
	next               time.Time
	started, ended     time.Time
	blocks             [3]chat.BlockState
	summary            chat.SummaryState
	review             *markdown.MarkdownDocument
	showReview, copied bool
	attached           bool
}

type chatDemoState struct {
	conversation  chat.ConversationState
	turns         []*chatDemoTurn
	draft         string
	serial        int
	model, effort string
	attachment    bool
}

func (s *chatDemoState) start(prompt string, now time.Time) {
	s.serial++
	s.turns = append(s.turns, &chatDemoTurn{id: fmt.Sprint(s.serial), prompt: prompt, next: now.Add(350 * time.Millisecond), started: now, review: markdown.ParseMarkdown("```diff\n- renderEveryBlock()\n+ renderVisibleBlocks()\n```"), blocks: [3]chat.BlockState{{Expanded: true}, {}, {}}})
	// The demonstration bounds retained history rather than claiming unlimited
	// transcript performance. A production host owns its retention policy.
	if len(s.turns) > 8 {
		s.turns = s.turns[len(s.turns)-8:]
	}
	s.conversation.JumpToLatest()
}

func (s *chatDemoTurn) advance(now time.Time) {
	if s.status != chat.Running {
		return
	}
	s.step++
	s.next = now.Add(350 * time.Millisecond)
	if s.step >= 4 {
		s.copied = false
		text := []rune(chatDemoAnswer)
		s.answer = markdown.ParseMarkdown(string(text[:min(len(text), (s.step-3)*48)]))
		if (s.step-3)*48 >= len(text) {
			s.status = chat.Completed
			s.ended = now
		}
	}
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
		state.chat = &chatDemoState{model: "Demo model", effort: "Medium"}
		state.chat.start("Review this component and explain the results.", c.Now())
		turn := state.chat.turns[0]
		turn.step, turn.status, turn.answer = 20, chat.Completed, markdown.ParseMarkdown(chatDemoAnswer)
		turn.started, turn.ended = c.Now().Add(-561*time.Second), c.Now()
	}
	s := state.chat
	last := s.turns[len(s.turns)-1]
	if last.status == chat.Running {
		if !c.Now().Before(last.next) {
			last.advance(c.Now())
		}
		if last.status == chat.Running {
			c.After(max(time.Millisecond, last.next.Sub(c.Now())))
		}
	}
	t := c.Theme()
	ui.Column(c).Fill().MinWidth(0).Gap(t.Space(3)).Children(func() {
		ui.Row(c).FillWidth().Shrink(0).Wrap().AlignItems(ui.Center).Gap(t.Space(2)).Padding(0, t.Space(4)).Children(func() {
			components.Text(c, tr(c, "gallery.chat")).Bold().Grow(1)
			if components.Button(c, tr(c, "gallery.chat_replay"), components.ButtonSecondary).Clicked() {
				s.turns = nil
				s.start("Review this component and explain the results.", c.Now())
				c.Invalidate()
			}
			if last.status == chat.Running {
				if components.Button(c, tr(c, "gallery.chat_fail"), components.ButtonOutline).Clicked() {
					last.status = chat.Failed
					last.ended = c.Now()
				}
			} else if last.status == chat.Failed || last.status == chat.Cancelled {
				if components.Button(c, tr(c, "gallery.chat_retry"), components.ButtonOutline).Clicked() {
					last.status, last.next = chat.Running, c.Now()
					last.state.ResetExpansion()
					last.ended = time.Time{}
					c.Invalidate()
				}
			}
		})
		chat.Conversation(c, &s.conversation, tr(c, "gallery.chat_latest"), func(c *ui.Context) {
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
				blocks := []chat.Block{
					{ID: "summary", State: &turn.blocks[0], Kind: "reasoning-summary", Title: tr(c, "gallery.chat_reasoning"), Status: chat.Completed,
						Content: func(c *ui.Context) {
							components.Text(c, "I will inspect the source and check the examples in parallel.").Selectable()
						}},
					{ID: "call-source", State: &turn.blocks[1], Preview: "component.go", DetailHeight: 120, Kind: "tool", Title: "read_source", Status: toolA, StatusLabel: chatStatusLabel(c, toolA),
						Content: func(c *ui.Context) {
							components.Text(c, fmt.Sprintf("component.go · %d / 3 chunks", min(turn.step, 3))).Selectable()
						}},
					{ID: "call-examples", State: &turn.blocks[2], Preview: "examples/", DetailHeight: 100, Kind: "tool", Title: "check_examples", Status: toolB, StatusLabel: chatStatusLabel(c, toolB),
						Content: func(c *ui.Context) {
							if turn.step >= 2 {
								components.Text(c, "3 examples checked. No errors.").Selectable()
							}
						}},
				}
				var final chat.View
				if turn.answer != nil {
					final = func(c *ui.Context) {
						turn.answer.Render(c, markdown.MarkdownOptions{CopyCodeLabel: tr(c, "gallery.copy_code"), CodeHeaders: true}).FillWidth()
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
					Final: final, Blocks: blocks,
					Footer: func(c *ui.Context) { chatDemoResults(c, turn) },
					Actions: func(c *ui.Context) {
						if turn.answer != nil {
							label := tr(c, "gallery.chat_copy_answer")
							if turn.copied {
								label = tr(c, "gallery.chat_copied")
							}
							if chat.IconAction(c, label, lucide.Copy).Clicked() {
								text := []rune(chatDemoAnswer)
								c.WriteClipboard(string(text[:min(len(text), max(0, turn.step-3)*48)]))
								turn.copied = true
							}
						}
					}, HoldOpen: !s.conversation.AtLatest(),
				})
			}
		})
		ui.Column(c).FillWidth().Shrink(0).MinWidth(0).Padding(0, t.Space(4), t.Space(4)).Children(func() {
			result := chat.Composer(c, &s.draft, chat.ComposerOptions{
				Label: tr(c, "gallery.chat_message"), Placeholder: tr(c, "gallery.chat_placeholder"),
				SendLabel: tr(c, "gallery.chat_send"), StopLabel: tr(c, "gallery.chat_stop"), Running: last.status == chat.Running,
				Accessory: func(c *ui.Context) { chatDemoToolbar(c, s) },
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
			if result.Stop {
				last.status = chat.Cancelled
				last.ended = c.Now()
			}
			if result.Send {
				s.start(s.draft, c.Now())
				s.turns[len(s.turns)-1].attached = s.attachment
				s.attachment = false
				s.draft = ""
				c.Invalidate()
			}
		})
	})
}
