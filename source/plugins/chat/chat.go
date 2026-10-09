// Package chat composes native conversation views from Onevra controls and
// MyGo primitives. Applications own messages, stable IDs, localized labels,
// parsed documents and asynchronous work. Builders run on the UI thread and
// must not start work. Retain state per conversation and turn, not per frame.
package chat

import (
	"math"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// Status describes a turn or block's application-reported lifecycle. A tool
// finishing does not finish its turn. Only Completed permits automatic folding.
type Status uint8

const (
	Running Status = iota
	Waiting
	Completed
	Failed
	Cancelled
)

// View builds arbitrary content, including retained Markdown documents,
// attachments, confirmation controls or bounded tool-output previews.
type View func(*ui.Context)

// Block is one execution step. ID is nonempty and unique within the turn;
// update a tool's existing block by its call ID, even if results arrive out of
// order. Slice order is presentation order, not result-arrival order. Kind is
// application-defined metadata; it does not constrain the content builder.
type Block struct {
	ID, Kind, Title, StatusLabel string
	Status                       Status
	Content                      View
	// Preview is a short, single-line description shown beside the title.
	Preview string
	// State enables an independent lazy disclosure. Nil keeps content visible.
	State *BlockState
	// DetailHeight bounds an expanded output viewport in DIPs. Zero lets the
	// content determine height. The application still bounds retained output.
	DetailHeight float32
}

// BlockState belongs to one stable call ID and survives outer-turn folding.
type BlockState struct {
	Expanded bool
	Scroll   ui.ScrollState
}

// TurnState retains only presentation state. Do not share it between turns or
// copy it after first use. Its zero value follows the automatic fold policy.
type TurnState struct {
	choice uint8
	body   ui.Handle
	folded bool
}

// SetExpanded records an explicit preference that survives stream updates.
func (s *TurnState) SetExpanded(open bool) {
	s.choice = 2
	if open {
		s.choice = 1
	}
}

// ResetExpansion restores automatic folding, for example for a new attempt.
func (s *TurnState) ResetExpansion() { s.choice, s.folded = 0, false }

// Expanded resolves presentation without modifying application data. Completed
// turns fold only when a final response exists. HoldOpen temporarily protects
// reading/selection; explicit user choices always take precedence.
func (s *TurnState) Expanded(status Status, hasFinal, holdOpen bool) bool {
	if s.choice != 0 {
		return s.choice == 1
	}
	if status != Completed || !hasFinal {
		s.folded = false
		return true
	}
	if !holdOpen {
		s.folded = true
	}
	return !s.folded
}

// TurnOptions describes one request and its response. Supply a stable, nonempty
// ID and localized Summary (also the disclosure's accessible name). Final is
// explicit: the last tool or assistant block is never inferred to be final.
type TurnOptions struct {
	ID, Summary   string
	Status        Status
	Prompt, Final View
	Blocks        []Block
	// Footer holds artifacts or result summaries, even without a final answer.
	// Actions holds response actions such as copy. Neither folds with execution.
	Footer, Actions View
	// HoldOpen delays automatic folding while the user reads older content.
	// ConversationState.AtLatest is useful for this policy. Focus inside the
	// execution body also protects it automatically, including selectable text.
	HoldOpen bool
}

// Turn displays a user prompt, lazy execution disclosure and independent final
// response. A collapsed disclosure invokes no block builders, including during
// closing; folding is immediate to avoid relaying out large hidden documents.
func Turn(c *ui.Context, state *TurnState, options TurnOptions) ui.Element {
	return ui.Column(c).Key(options.ID).FillWidth().MinWidth(0).Shrink(0).Gap(c.Theme().Space(4)).Children(func() {
		if options.Prompt != nil {
			Message(c, true, func(c *ui.Context) { Bubble(c, options.Prompt) })
		}
		if len(options.Blocks) != 0 {
			ui.Column(c).Key("execution").FillWidth().MinWidth(0).Gap(c.Theme().Space(2)).Children(func() {
				scope := *conversationSlot(c)
				open := state.Expanded(options.Status, options.Final != nil, options.HoldOpen || scope.reading || state.body.FocusWithin(c))
				// The unstyled disclosure supplies native keyboard and expanded
				// semantics; content is built separately so hidden work is zero.
				var p ui.CollapsibleParts
				ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(c.Theme().Space(3)).Children(func() {
					divider(c)
					p = ui.CollapsibleBase(c, &open)
					p.Trigger.MaxWidthPercent(80).AlignItems(ui.Center).Gap(c.Theme().Space(1)).Padding(c.Theme().Space(1)).Radius(c.Theme().Radius).
						Label(options.Summary).Children(func() {
						components.Text(c, options.Summary).MinWidth(0).TextColor(c.Theme().TextMuted)
						arrow := lucide.ChevronRight
						if open {
							arrow = lucide.ChevronDown
						}
						components.Icon(c, arrow).Size(c.Theme().Space(4), c.Theme().Space(4)).Shrink(0).TextColor(c.Theme().TextMuted)
					})
					divider(c)
				})
				if p.Trigger.Changed() {
					state.SetExpanded(open)
					if scope.state != nil {
						scope.state.Scroll.Y = scope.previousY
						scope.state.paused, scope.state.pauseY = true, scope.previousY
					}
				}
				if open {
					ui.Column(c).Key("steps").Bind(&state.body).FillWidth().MinWidth(0).Gap(c.Theme().Space(3)).Children(func() {
						for _, block := range options.Blocks {
							ExecutionBlock(c, block)
						}
					})
				}
			})
		}
		if options.Final != nil {
			ui.Column(c).Key("final").FillWidth().MinWidth(0).Children(func() { Message(c, false, options.Final) })
		}
		if options.Footer != nil {
			ui.Column(c).Key("footer").FillWidth().MinWidth(0).Children(func() { options.Footer(c) })
		}
		if options.Actions != nil {
			ui.Row(c).Key("actions").FillWidth().Wrap().Gap(c.Theme().Space(2)).Children(func() { options.Actions(c) })
		}
	})
}

// ExecutionBlock renders an extensible execution step. Content costs and
// truncation are application-owned; do not build unbounded logs in a preview.
func ExecutionBlock(c *ui.Context, block Block) ui.Element {
	t := c.Theme()
	return ui.Column(c).Key(block.ID).FillWidth().MinWidth(0).Shrink(0).
		Radius(t.Radius).Border(1, t.Border).Background(t.Text.Alpha(.025)).Children(func() {
		header := func() {
			icon := lucide.LoaderCircle
			switch block.Status {
			case Completed:
				icon = lucide.Check
			case Failed:
				icon = lucide.CircleAlert
			case Waiting, Cancelled:
				icon = lucide.CirclePause
			}
			components.Icon(c, icon).Size(t.Space(4), t.Space(4)).Shrink(0).TextColor(t.TextMuted)
			components.Text(c, block.Title).Bold().Grow(2).MinWidth(0).SingleLine().Ellipsis("…")
			if block.Preview != "" {
				components.Text(c, block.Preview).Grow(1).MinWidth(0).SingleLine().Ellipsis("…").TextColor(t.TextMuted)
			}
			if block.StatusLabel != "" {
				components.Text(c, block.StatusLabel).TextColor(t.TextMuted)
			}
			if block.State != nil && block.Content != nil {
				icon := lucide.ChevronRight
				if block.State.Expanded {
					icon = lucide.ChevronDown
				}
				components.Icon(c, icon).Size(t.Space(3), t.Space(3)).Shrink(0)
			}
		}
		if block.State != nil && block.Content != nil {
			p := ui.CollapsibleBase(c, &block.State.Expanded)
			p.Trigger.FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(2)).Padding(t.Space(2)).Label(block.Title).Children(header)
			if p.Trigger.Changed() {
				pauseConversation(c)
			}
		} else {
			ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Children(header)
		}
		if block.Content != nil && (block.State == nil || block.State.Expanded) {
			body := func() { ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(3)).Children(func() { block.Content(c) }) }
			if block.DetailHeight > 0 {
				viewport := components.ScrollView(c, components.ScrollViewOptions{}).FillWidth().Height(block.DetailHeight)
				if block.State != nil {
					if block.Status == Running && block.State.Scroll.Y >= block.State.Scroll.MaxY-2 {
						block.State.Scroll.Y = math.MaxFloat32
					}
					viewport.TrackScroll(&block.State.Scroll)
				}
				viewport.Children(body)
			} else {
				body()
			}
		}
	})
}

func divider(c *ui.Context) { ui.Box(c).Grow(1).Height(1).Background(c.Theme().Border) }

func pauseConversation(c *ui.Context) {
	scope := *conversationSlot(c)
	if scope.state != nil {
		scope.state.Scroll.Y = scope.previousY
		scope.state.paused, scope.state.pauseY = true, scope.previousY
	}
}

// Message lays out a message independently of its surface. Outgoing messages
// align to the trailing edge; incoming content uses the available width.
func Message(c *ui.Context, outgoing bool, content View) ui.Element {
	e := ui.Column(c).FillWidth().MinWidth(0).Shrink(0)
	if outgoing {
		e.AlignItems(ui.End)
	}
	return e.Children(func() {
		if content != nil {
			content(c)
		}
	})
}

// Bubble supplies an optional muted message surface.
func Bubble(c *ui.Context, content View) ui.Element {
	t := c.Theme()
	return ui.Column(c).MinWidth(0).MaxWidthPercent(90).Padding(t.Space(3), t.Space(4)).Radius(t.Radius * 2).
		Background(t.Background.Mix(t.Text, .06)).Children(func() {
		if content != nil {
			content(c)
		}
	})
}

// ConversationState retains MyGo's scroll state. Use one per transcript;
// replacing a transcript should also replace this state. This first version
// uses ordinary scrolling to preserve selection; it does not virtualize history
// or anchor prepended rows. Applications must bound retained history themselves.
type ConversationState struct {
	Scroll ui.ScrollState
	paused bool
	pauseY float32
}

type conversationKey struct{}
type conversationScope struct {
	state     *ConversationState
	previousY float32
	reading   bool
}

func conversationSlot(c *ui.Context) *conversationScope {
	return ui.Local(c.Root(), conversationKey{}, func() conversationScope { return conversationScope{} })
}

// AtLatest reports whether the previously laid-out viewport is at the bottom.
func (s *ConversationState) AtLatest() bool { return s.Scroll.Y >= s.Scroll.MaxY-2 }

// JumpToLatest resumes bottom following on the next build.
func (s *ConversationState) JumpToLatest() {
	s.paused = false
	s.Scroll.Y = math.MaxFloat32
}

// Conversation owns a bounded scroll viewport and a jump button. Give its
// parent a bounded height. Supply a localized jumpLabel. Following uses MyGo's
// scroll state and yields when the user scrolls up. The button is overlaid so
// its appearance cannot change viewport height.
func Conversation(c *ui.Context, state *ConversationState, jumpLabel string, content View) ui.Element {
	if state.paused && state.Scroll.Y != state.pauseY && state.AtLatest() {
		state.paused = false
	}
	follow := state.AtLatest() && !state.paused
	previousY := state.Scroll.Y
	if follow {
		state.JumpToLatest()
	}
	return ui.Column(c).FillWidth().Grow(1).MinHeight(0).MinWidth(0).Children(func() {
		viewport := components.ScrollView(c, components.ScrollViewOptions{}).Key("viewport").FillWidth().Grow(1).MinHeight(0).TrackScroll(&state.Scroll)
		viewport.Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).Padding(c.Theme().Space(4), c.Theme().Space(4), c.Theme().Space(14)).Children(func() {
				ui.Column(c).FillWidth().MaxWidth(c.Theme().Rem(48)).MinWidth(0).Gap(c.Theme().Space(6)).Children(func() {
					slot := conversationSlot(c)
					saved := *slot
					*slot = conversationScope{state: state, previousY: previousY, reading: !follow}
					defer func() { *slot = saved }()
					if content != nil {
						content(c)
					}
				})
			})
		})
		// A disclosure can cancel this frame's follow without moving its header.
		viewport.TrackScroll(&state.Scroll)
		if !follow {
			if components.Button(c, jumpLabel, components.ButtonSecondary).Absolute().Attach(ui.AnchorBottom, ui.AnchorBottom).Margin(c.Theme().Space(2)).Clicked() {
				state.JumpToLatest()
				c.Invalidate()
			}
		}
	})
}

// ComposerOptions supplies localized labels and optional application controls.
// Enter inserts a newline using the native multiline editor. Sending is an
// explicit button action, so IME composition cannot accidentally submit a turn.
type ComposerOptions struct {
	Label, Placeholder, SendLabel, StopLabel string
	Running                                  bool
	Accessory                                View
	// Attachments appears above the editor; Footer appears below the surface.
	Attachments, Footer View
}

// ComposerResult contains frame-local actions. The application clears the
// draft after accepting Send and owns request cancellation for Stop.
type ComposerResult struct {
	Root, Input ui.Element
	Send, Stop  bool
}

// Composer builds a multiline editor with send/stop and an accessory slot.
// Keep draft outside the frame; a running request never discards the next draft.
func Composer(c *ui.Context, draft *string, options ComposerOptions) ComposerResult {
	var result ComposerResult
	t := c.Theme()
	result.Root = ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Gap(t.Space(2)).Children(func() {
		surface := ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(3)).Gap(t.Space(2)).
			Background(t.Surface).Border(1, t.Border).Radius(t.Radius * 2).Children(func() {
			if options.Attachments != nil {
				ui.Row(c).Key("attachments").FillWidth().Wrap().Gap(t.Space(2)).Children(func() { options.Attachments(c) })
			}
			result.Input = components.TextBox(c, draft, components.TextBoxOptions{Multiline: true}).FillWidth().Lines(1, 6).MinHeight(t.Space(8)).
				Label(options.Label).Placeholder(options.Placeholder).Border(0, ui.Color{}).Background(ui.Color{})

			ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
				ui.Row(c).Grow(1).MinWidth(0).Wrap().Gap(t.Space(1)).Children(func() {
					if options.Accessory != nil {
						options.Accessory(c)
					}
				})
				if options.Running {
					result.Stop = IconAction(c, options.StopLabel, lucide.Square).Clicked()
				} else {
					ui.Box(c).Disabled(strings.TrimSpace(*draft) == "" || result.Input.Composing()).Children(func() {
						result.Send = IconAction(c, options.SendLabel, lucide.ArrowUp).Background(t.Text).TextColor(t.Background).Clicked()
					})
				}
			})
		})
		if result.Input.Focused() {
			surface.BorderColor(t.Accent)
		}
		if options.Footer != nil {
			ui.Row(c).Key("composer-footer").FillWidth().MinWidth(0).Wrap().Gap(t.Space(2)).Children(func() { options.Footer(c) })
		}
	})
	return result
}

// IconAction is a keyboard-operable, named icon button with a focus tooltip.
func IconAction(c *ui.Context, label string, icon func() *ui.SVG) ui.Element {
	t := c.Theme()
	b := components.Button(c, "", components.ButtonGhost).Label(label).Padding(t.Space(2)).Radius(t.Space(5)).Shrink(0)
	b.Children(func() { components.Icon(c, icon).Size(t.Space(4), t.Space(4)) })
	components.Tooltip(c, b, label)
	return b
}
