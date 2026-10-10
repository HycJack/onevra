package chat

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// FloatingConversationState retains composer measurements. Keep one alongside
// each ConversationState and share it with the floating GlassComposer.
type FloatingConversationState struct {
	height       float32
	cardHeight   float32
	footerHeight float32
}

// FloatingConversation lets the transcript scroll behind the composer, but
// reserves enough trailing content space to leave the latest turn unobscured.
// Give the parent a bounded height and supply a localized jumpLabel. Build a
// GlassComposer with the same layout in composer to fade beneath its footer.
// Applications own both states and all content; builders run synchronously.
func FloatingConversation(c *ui.Context, layout *FloatingConversationState, state *ConversationState, jumpLabel string, content, composer View) ui.Element {
	return ui.Column(c).FillWidth().Grow(1).MinHeight(0).MinWidth(0).Clip().Children(func() {
		height := layout.height
		if height <= 0 {
			height = 140
		}
		if state.paused && state.Scroll.Y != state.pauseY && state.AtLatest() {
			state.paused = false
		}
		follow := state.AtLatest() && !state.paused
		previousY := state.Scroll.Y
		if follow {
			state.JumpToLatest()
		}
		viewport := components.ScrollView(c, components.ScrollViewOptions{}).Key("transcript").FillWidth().Grow(1).MinHeight(0).TrackScroll(&state.Scroll)
		viewport.Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).Padding(22, 20, 0, 20).Children(func() {
				ui.Column(c).FillWidth().MaxWidth(ContentWidth).MinWidth(0).Gap(48).Children(func() {
					slot := conversationSlot(c)
					saved := *slot
					*slot = conversationScope{state: state, previousY: previousY, reading: !follow}
					defer func() { *slot = saved }()
					if content != nil {
						content(c)
					}
				})
				ui.Box(c).Key("composer-clearance").Height(height + 22).Shrink(0)
			})
		})
		// Disclosures may cancel following while building the transcript.
		viewport.TrackScroll(&state.Scroll)
		if !follow {
			b := components.Button(c, "", components.ButtonSecondary).Label(jumpLabel).Absolute().Attach(ui.AnchorBottom, ui.AnchorBottom).
				Bottom(height+12).Size(28, 28).Padding(0).Radius(14)
			b.Children(func() { components.Icon(c, lucide.ArrowDown).Size(14, 14) })
			components.Tooltip(c, b, jumpLabel)
			if b.Clicked() {
				state.JumpToLatest()
				c.Invalidate()
			}
		}
		ui.Row(c).Key("floating-composer").Absolute().Attach(ui.AnchorBottom, ui.AnchorBottom).
			FillWidth().Justify(ui.Center).Padding(0, 20, 8, 20).Draw(func(p *ui.Painter, r ui.Rect) {
			// Fade behind the card's lower half, reaching an opaque surface
			// before the footer so scrolling messages cannot obscure its labels.
			footerHeight := layout.footerHeight + 8
			fadeHeight := layout.cardHeight / 2
			bottom := r.Y + r.H
			p.FillGradient(ui.Rect{X: r.X, Y: bottom - footerHeight - fadeHeight, W: r.W, H: fadeHeight},
				ui.LinearGradient{From: p.Theme().Background.Alpha(0), To: p.Theme().Background, Angle: 180}, 0)
			p.Fill(ui.Rect{X: r.X, Y: bottom - footerHeight, W: r.W, H: footerHeight}, p.Theme().Background, 0)
			if layout.height != r.H {
				layout.height = r.H
				c.Invalidate()
			}
		}).Children(func() {
			ui.Column(c).FillWidth().MaxWidth(ContentWidth).MinWidth(0).Children(func() {
				if composer != nil {
					composer(c)
				}
			})
		})
	})
}

// GlassComposer keeps Onevra's editor and actions, replacing only its surface.
// The footer stays outside the glass so the card remains visually lightweight.
// layout may be nil outside a floating conversation.
func GlassComposer(c *ui.Context, layout *FloatingConversationState, draft *string, options ComposerOptions) ComposerResult {
	theme := c.Theme()
	result := Composer(c, draft, options)
	tint := theme.Surface.Alpha(.6)
	if result.Dragging {
		tint = theme.Surface.Mix(theme.Text, .09).Alpha(.6)
	}
	result.Card.Material(glass.Glass{Style: glass.Regular, Tint: tint})
	if layout != nil {
		layout.footerHeight = 0
		if options.Footer != nil {
			layout.footerHeight = 28
		}
		result.Card.Draw(func(_ *ui.Painter, r ui.Rect) {
			if layout.cardHeight != r.H {
				layout.cardHeight = r.H
				c.Invalidate()
			}
		})
	}
	return result
}
