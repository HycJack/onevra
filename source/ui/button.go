package ui

import "github.com/egoist/mygo/ui"

// ButtonVariant selects a button's visual weight.
type ButtonVariant int

const (
	// ButtonPrimary is the accent-filled main action.
	ButtonPrimary ButtonVariant = iota
	// ButtonSecondary is a bordered action with a light fill in light mode
	// and a high-contrast fill in dark mode.
	ButtonSecondary
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
// its children must be text. For stable identity, key a parent container rather
// than this element: interaction styling is resolved during construction.
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
	case ButtonSecondary:
		bg, fg, border = t.Text, t.Background, t.Border
		hover, pressed = bg.Mix(t.Background, .12), bg.Mix(t.Background, .25)
		if !t.Dark {
			bg, fg = t.Background.Mix(t.Text, .08), t.Text
			hover, pressed = t.Background.Mix(t.Text, .12), t.Background.Mix(t.Text, .18)
		}
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
	if variant == ButtonPrimary {
		border = bg
	}
	// A transparent Ghost border reserves the same space as the filled variants.
	b.Radius(t.Radius).Background(bg).TextColor(fg).Border(1, border)
	if label != "" {
		b.Label(label)
		b.Children(func() { ui.Text(c, label).SingleLine() })
	}
	return dimmed(b)
}
