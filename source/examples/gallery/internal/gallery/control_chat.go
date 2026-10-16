package gallery

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

const chatDemoAnswer = "The review is complete.\n\n- **Two tools** ran concurrently and completed in a different order.\n- Execution details remain available above.\n- The final answer stays visible after the process folds.\n\n```go\n// Share one Options value across a conversation.\nmarkdown.View(c, reply, &s.markdown)\n```"

const chatStreamPrompt = "Explain how the composer controls fit together."

const chatStreamAnswer = "The composer is a card with **slots** the application fills:\n\n1. *Attachments* sit above the editor as file chips or image tiles.\n2. *Accessory* holds pickers for the model and reasoning effort.\n3. *Trailing* keeps compact controls, such as the context meter, beside Send.\n4. *Footer* shows the project, branch and working mode below the card.\n\nEach reply streams through `markdown.Stream`, which reveals words at an even pace, and switches to `markdown.View` once complete:\n\n```go\nif turn.done {\n\tmarkdown.View(c, turn.text, &s.markdown)\n} else {\n\tmarkdown.Stream(c, turn.text, &s.markdown)\n}\n```\n\nThe application owns the request, cancellation and every label."

// The simulated timeline of a streamed turn: thinking, two tools that finish
// out of order, then the reply arriving in bursts as a network delivers it.
const (
	chatThinkFor    = 800 * time.Millisecond
	chatToolBDone   = 1600 * time.Millisecond
	chatToolADone   = 2200 * time.Millisecond
	chatReplyStart  = 2400 * time.Millisecond
	chatBurstEvery  = 80 * time.Millisecond
	chatBurstLength = 12
)

type chatDemoFile struct {
	name          string
	image, sample bool
}

type chatDemoTurn struct {
	id, prompt     string
	files          []chatDemoFile
	state          chat.TurnState
	status         chat.Status
	answer         string
	received       int
	started, ended time.Time
	activities     [3]chat.ActivityState
	group          chat.GroupState
}

type chatDemoState struct {
	conversation                            chat.ConversationState
	layout                                  chat.FloatingConversationState
	turns                                   []*chatDemoTurn
	next                                    int
	draft                                   string
	files                                   []chatDemoFile
	model, effort, branch                   string
	worktree                                bool
	modelPicker, effortPicker, branchPicker chat.PickerState
	meter, quota                            chat.MeterState
	quotaNext, quotaRefreshed               time.Time
	thumbnail                               *ui.Bitmap
	markdown                                markdown.Options
}

func newChatDemoState(now time.Time) *chatDemoState {
	s := &chatDemoState{
		model: "demo", effort: "medium", branch: "main", worktree: true, thumbnail: chatDemoThumbnail(),
		turns: []*chatDemoTurn{{
			id: "sample", prompt: "Review this component and explain the results.",
			status: chat.Completed, answer: chatDemoAnswer, received: len(chatDemoAnswer),
			started: now.Add(-561 * time.Second), ended: now,
			activities: [3]chat.ActivityState{{Expanded: true}, {}, {}},
		}},
	}
	s.turns[0].group.SetExpanded(true)
	return s
}

func (s *chatDemoState) running() *chatDemoTurn {
	if n := len(s.turns); n > 0 && s.turns[n-1].status == chat.Running {
		return s.turns[n-1]
	}
	return nil
}

// start begins a simulated reply, stopping one already running, as a new
// message interrupts an agent.
func (s *chatDemoState) start(prompt string, files []chatDemoFile, now time.Time) {
	s.stop(now)
	s.next++
	s.turns = append(s.turns, &chatDemoTurn{id: fmt.Sprint("turn-", s.next), prompt: prompt, files: files,
		status: chat.Running, answer: chatStreamAnswer, started: now})
	s.conversation.JumpToLatest()
}

func (s *chatDemoState) stop(now time.Time) {
	if turn := s.running(); turn != nil {
		turn.status, turn.ended = chat.Cancelled, now
	}
}

// advance delivers the reply received by now. It only computes from the
// clock, as a host would apply a network update on the UI thread.
func (turn *chatDemoTurn) advance(now time.Time) {
	if turn.status != chat.Running {
		return
	}
	elapsed := now.Sub(turn.started)
	if elapsed < chatReplyStart {
		return
	}
	n := min(len(turn.answer), int((elapsed-chatReplyStart)/chatBurstEvery)*chatBurstLength)
	for n < len(turn.answer) && !utf8.RuneStart(turn.answer[n]) {
		n++
	}
	turn.received = n
	if n == len(turn.answer) {
		turn.status, turn.ended = chat.Completed, now
	}
}

func chatStatusLabel(c *ui.Context, status chat.Status) string {
	switch status {
	case chat.Completed:
		return tr(c, "chatCompleted")
	case chat.Failed:
		return tr(c, "chatFailed")
	case chat.Cancelled:
		return tr(c, "chatCancelled")
	case chat.Waiting:
		return tr(c, "chatWaiting")
	default:
		return tr(c, "chatRunning")
	}
}

// chatToolbar starts a simulated streamed reply from the titlebar; the page
// itself only holds the conversation.
func chatToolbar(c *ui.Context, state *controlState) {
	if state.chat == nil {
		state.chat = newChatDemoState(c.Now())
	}
	if sidebar.IconButton(c, lucide.Play(), tr(c, "chatStreamReply")).Clicked() {
		state.chat.start(chatStreamPrompt, nil, c.Now())
	}
}

func controlChat(c *ui.Context, state *controlState, _ string) {
	if state.chat == nil {
		state.chat = newChatDemoState(c.Now())
	}
	s := state.chat
	if turn := s.running(); turn != nil {
		turn.advance(c.Now())
		c.After(chatBurstEvery)
	}
	ui.Column(c).Fill().MinWidth(0).Children(func() {
		chat.FloatingConversation(c, &s.layout, &s.conversation, tr(c, "chatLatest"), func(c *ui.Context) {
			for _, turn := range s.turns {
				chatDemoTurnView(c, s, turn)
			}
		}, func(c *ui.Context) { chatDemoComposer(c, s) })
	})
}

func chatDemoTurnView(c *ui.Context, s *chatDemoState, turn *chatDemoTurn) {
	end := turn.ended
	if turn.status == chat.Running {
		end = c.Now()
	}
	elapsed := end.Sub(turn.started)
	settled := func(at time.Duration) chat.Status {
		switch {
		case elapsed >= at:
			return chat.Completed
		case turn.status == chat.Running:
			return chat.Running
		}
		return turn.status
	}
	activities := []chat.Activity{{ID: "summary", State: &turn.activities[0], Icon: lucide.Brain, Action: tr(c, "chatThink"), Status: settled(chatThinkFor),
		Content: func(c *ui.Context) {
			components.Text(c, "I will inspect the source and check the examples in parallel.").Selectable()
		}}}
	if elapsed >= chatThinkFor {
		activities = append(activities,
			chat.Activity{ID: "call-source", State: &turn.activities[1], Icon: lucide.FileText, Action: "read_source", Detail: "component.go", MaxHeight: 120, Mono: true, Status: settled(chatToolADone),
				Content: func(c *ui.Context) {
					chunks := 3
					if elapsed < chatToolADone {
						chunks = int(3 * (elapsed - chatThinkFor) / (chatToolADone - chatThinkFor))
					}
					components.Text(c, fmt.Sprintf("component.go · %d / 3 chunks", chunks)).Selectable()
				}},
			chat.Activity{ID: "call-examples", State: &turn.activities[2], Icon: lucide.Terminal, Action: "check_examples", Detail: "examples/", MaxHeight: 100, Mono: true, Status: settled(chatToolBDone),
				Content: func(c *ui.Context) {
					if elapsed >= chatToolBDone {
						components.Text(c, "3 examples checked. No errors.").Selectable()
					}
				}})
	}
	title := tr(c, "chatReasoning")
	if live := activities[len(activities)-1]; turn.status == chat.Running && turn.received == 0 {
		// A live group names its latest activity.
		title = strings.TrimSpace(live.Action + " " + live.Detail)
	}
	steps := []chat.Step{{ID: "work", Title: title, Group: &turn.group, Activities: activities}}
	reply := turn.answer[:turn.received]
	var final chat.View
	switch {
	case turn.status == chat.Running && reply != "":
		// The reply streams as prose after the work and becomes the final
		// response once complete.
		steps = append(steps, chat.Step{ID: "reply", Message: func(c *ui.Context) {
			markdown.Stream(c, reply, &s.markdown).FillWidth()
		}})
	case reply != "":
		final = func(c *ui.Context) { markdown.View(c, reply, &s.markdown).FillWidth() }
	}
	chat.Turn(c, &turn.state, chat.TurnOptions{
		ID: turn.id, Summary: chatTurnSummary(c, turn), Status: turn.status,
		Prompt: func(c *ui.Context) {
			if len(turn.files) > 0 {
				ui.Row(c).Wrap().Gap(6).Children(func() {
					for i, f := range turn.files {
						ui.Box(c).Key(fmt.Sprint(i)).Children(func() { chatDemoAttachment(c, s, f, "") })
					}
				})
			}
			components.Text(c, turn.prompt).Selectable()
		},
		Final: final, Steps: steps, Working: chatTurnSummary(c, turn),
		HoldOpen: !s.conversation.AtLatest(),
	})
}

func chatDemoComposer(c *ui.Context, s *chatDemoState) {
	options := chat.ComposerOptions{
		Label: tr(c, "chatMessage"), Placeholder: tr(c, "chatPlaceholder"),
		SendLabel: tr(c, "chatSend"), StopLabel: tr(c, "chatStop"),
		Running: s.running() != nil, AcceptFiles: true,
		Accessory: func(c *ui.Context) { chatDemoAccessory(c, s) },
		Trailing:  func(c *ui.Context) { chatDemoTrailing(c, s) },
		Footer:    func(c *ui.Context) { chatDemoFooter(c, s) },
	}
	if len(s.files) > 0 {
		options.Attachments = func(c *ui.Context) {
			for i, f := range s.files {
				ui.Box(c).Key(fmt.Sprint(f.name, i)).Children(func() {
					if chatDemoAttachment(c, s, f, tr(c, "chatRemoveAttachment")) {
						s.files = append(s.files[:i:i], s.files[i+1:]...)
					}
				})
			}
		}
	}
	result := chat.GlassComposer(c, &s.layout, &s.draft, options)
	// The demonstration names dropped files without reading them.
	for _, path := range result.Dropped {
		s.files = append(s.files, chatDemoFile{name: filepath.Base(path), image: chatDemoImage(path)})
	}
	if result.Send {
		s.start(strings.TrimSpace(s.draft), s.files, c.Now())
		s.draft, s.files = "", nil
	}
	if result.Stop {
		s.stop(c.Now())
	}
}

func chatDemoImage(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".heic":
		return true
	}
	return false
}

// chatDemoAttachment shows every file as a chip of one height; an image
// shows its thumbnail, or an image icon until one is prepared. An empty
// removeLabel makes it read-only.
func chatDemoAttachment(c *ui.Context, s *chatDemoState, f chatDemoFile, removeLabel string) bool {
	options := chat.AttachmentOptions{Label: f.name, RemoveLabel: removeLabel}
	switch {
	case f.image && f.sample:
		options.Image = s.thumbnail
	case f.image:
		// The demonstration never reads dropped files to decode them.
		options.Icon = lucide.Image
	}
	return chat.Attachment(c, options).Remove
}

// chatDemoThumbnail paints a small sample image once, as a host would decode
// a thumbnail off the UI thread.
func chatDemoThumbnail() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 44, 44))
	for y := range 44 {
		for x := range 44 {
			img.Set(x, y, color.RGBA{R: uint8(90 + x*3), G: uint8(140 + y*2), B: 220, A: 255})
		}
	}
	return ui.NewBitmap(img)
}
