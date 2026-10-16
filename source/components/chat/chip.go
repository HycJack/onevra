package chat

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// ChipOptions describes a Chip's leading content and role. Leading replaces
// Image, which replaces Icon; each is optional.
type ChipOptions struct {
	Icon  func() *ui.SVG
	Image ui.ImageSource
	// Leading builds non-interactive content at the suggested size, such as
	// a product mark drawn in its own colors.
	Leading func(c *ui.Context, size float32)
	// Caret marks a chip that opens a choice.
	Caret bool
	// ReadOnly shows a fixed value, such as the current project, as an
	// undimmed, unfocusable label instead of a button.
	ReadOnly bool
}

// Chip is a compact composer or footer control for a setting such as a
// model, reasoning level, project or branch. Its label shortens with an
// ellipsis when the row is tight. Read Clicked or attach a native menu with
// Element.Menu on interactive chips; a ReadOnly chip is not a button.
func Chip(c *ui.Context, label string, options ChipOptions) ui.Element {
	t := c.Theme()
	var b ui.Element
	if options.ReadOnly {
		// Disabled dims the whole subtree in MyGo; a fixed value is text,
		// not an unavailable action.
		b = ui.Row(c).Justify(ui.Center).Shrink(0).Cursor(ui.CursorNotAllowed)
	} else {
		b = ui.ButtonBase(c)
	}
	b.Label(label).Height(26).MinWidth(0).Shrink(1).Padding(0, 7).Radius(6).AlignItems(ui.Center).Gap(6).
		Background(ui.Color{}).TextColor(secondary(t))
	if !options.ReadOnly && (b.Hovered() || b.Focused()) {
		b.Background(overlay(t))
	}
	if b.IsDisabled() {
		b.Opacity(.7)
	}
	b.Children(func() {
		switch {
		case options.Leading != nil:
			options.Leading(c, 14)
		case options.Image != nil:
			ui.Image(c, options.Image).Size(14, 14).Fit(ui.Contain).Radius(3).Shrink(0)
		case options.Icon != nil:
			components.Icon(c, options.Icon).Size(12, 12).Shrink(0).TextColor(tertiary(t))
		}
		components.Text(c, label).MinWidth(0).SingleLine().Ellipsis("…").FontSize(12.5).TextColor(secondary(t))
		if options.Caret && !options.ReadOnly {
			components.Icon(c, lucide.ChevronDown).Size(10.5, 10.5).Shrink(0).TextColor(tertiary(t).Alpha(.7))
		}
	})
	return b
}

// SwitchChip is a small switch with a Chip's proportions for composer
// options, such as working in a separate branch. Read Changed for the toggle.
func SwitchChip(c *ui.Context, on *bool, label string) ui.Element {
	t := c.Theme()
	b := ui.SwitchBase(c, on).Label(label).Height(26).MinWidth(0).Shrink(1).Padding(0, 7).Radius(6).AlignItems(ui.Center).Gap(6).
		Background(ui.Color{})
	if b.Hovered() || b.Focused() {
		b.Background(overlay(t))
	}
	if b.IsDisabled() {
		b.Opacity(.7)
	}
	var target float32
	if *on {
		target = 1
	}
	pos := b.Animate("knob", target, 140*time.Millisecond)
	b.Children(func() {
		off := t.Border.Mix(t.Text, .12)
		ui.Box(c).Size(22, 13).Shrink(0).Radius(6.5).Background(off.Mix(t.Accent, pos)).Draw(func(p *ui.Painter, r ui.Rect) {
			in := float32(1.5)
			d := r.H - 2*in
			knob := ui.Rect{X: r.X + in + pos*(r.W-r.H), Y: r.Y + in, W: d, H: d}
			p.Shadow(knob, d/2, 0, .5, 1.5, 0, ui.RGBA(0, 0, 0, .2))
			p.Fill(knob, ui.RGB(255, 255, 255), d/2)
		})
		components.Text(c, label).MinWidth(0).SingleLine().Ellipsis("…").FontSize(12.5).TextColor(secondary(t))
	})
	return b
}
