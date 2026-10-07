package ui

import "github.com/egoist/mygo/ui"

// ButtonVariant selects a button's visual weight.
type ButtonVariant int

const (
	// ButtonPrimary is the accent-filled main action.
	ButtonPrimary ButtonVariant = iota
	// ButtonDefault is the neutral surface action.
	ButtonDefault
	// ButtonOutline is the bordered quiet action.
	ButtonOutline
	// ButtonGhost is the borderless quiet action.
	ButtonGhost
	// ButtonLink is inline text with an underline on hover.
	ButtonLink
)

// Button creates a labeled button in the oneo style with a single standard
// size. Ask Clicked whether it was clicked; give it other content with Children
// and an empty label. The label is also the accessible name. ButtonLink uses
// text height without padding or a background and supports RichText composition;
// its children must be text.
func Button(c *ui.Context, label string, variant ButtonVariant) *ui.Element {
	if variant == ButtonLink {
		return dimmed(ui.Link(c, label, ""))
	}
	t := c.Theme()
	b := ui.ButtonBase(c).Padding(t.Space(1.5), t.Space(3.5))
	bg, hover, pressed, fg, border := t.Surface, t.SurfaceHover, t.SurfacePressed, t.Text, t.Border
	switch variant {
	case ButtonPrimary:
		bg, hover, pressed, fg, border = t.Accent, t.AccentHover, t.AccentPressed, t.AccentText, ui.Color{}
	case ButtonOutline:
		bg, border = ui.Color{}, t.Border
	case ButtonGhost:
		bg, border = ui.Color{}, ui.Color{}
	}
	if !b.IsDisabled() {
		switch {
		case b.Pressed():
			bg = pressed
		case b.Hovered():
			bg = hover
		}
	}
	b.Radius(t.Radius).Background(bg).TextColor(fg)
	if border.A > 0 {
		b.Border(1, border)
	}
	if label != "" {
		b.Label(label)
		b.Children(func() { ui.Text(c, label).SingleLine() })
	}
	return dimmed(b)
}
