package controls

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
	// ButtonLink looks like a text link.
	ButtonLink
)

// ButtonSize selects a button's padding and text size.
type ButtonSize int

const (
	ButtonSmall ButtonSize = iota
	ButtonMedium
	ButtonLarge
)

// Button creates a labeled button in the oneo style. Ask Clicked whether it
// was clicked; give it other content with Children and an empty label. The
// label is also the accessible name.
func Button(c *ui.Context, label string, variant ButtonVariant, size ButtonSize) *ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c)
	switch size {
	case ButtonSmall:
		b.Padding(t.Space(1), t.Space(2.5)).FontSize(t.Rem(0.85))
	case ButtonLarge:
		b.Padding(t.Space(2), t.Space(4.5)).FontSize(t.Rem(1.125))
	default:
		b.Padding(t.Space(1.5), t.Space(3.5))
	}
	bg, hover, pressed, fg, border := t.Surface, t.SurfaceHover, t.SurfacePressed, t.Text, t.Border
	switch variant {
	case ButtonPrimary:
		bg, hover, pressed, fg, border = t.Accent, t.AccentHover, t.AccentPressed, t.AccentText, ui.Color{}
	case ButtonOutline:
		bg, border = ui.Color{}, t.Border
	case ButtonGhost:
		bg, border = ui.Color{}, ui.Color{}
	case ButtonLink:
		bg, border, fg = ui.Color{}, ui.Color{}, t.Accent
		b.Cursor(ui.CursorPointer)
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
