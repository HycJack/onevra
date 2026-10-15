package components

import "github.com/egoist/mygo/ui"

// IconButtonStyle selects an icon button's face.
type IconButtonStyle int

const (
	// IconButtonPlain tints a transparent face on hover and press.
	IconButtonPlain IconButtonStyle = iota
	// IconButtonToolbar raises the face, shadowed, as the chosen segment of a
	// segmented tab list is raised.
	IconButtonToolbar
)

// IconButton creates an icon-only button with an accessible label and a
// tooltip. It uses one theme-scaled square size. Applications own the label and
// handle Clicked on the result; MyGo owns its pointer, keyboard and focus
// behavior.
//
// IconButtonPlain keeps a transparent face and tints it on hover and press.
// IconButtonToolbar keeps the raised, shadowed face of ButtonToolbar, matching
// TabsSegmented's chosen segment.
func IconButton(c *ui.Context, icon *ui.SVG, label string, style IconButtonStyle) ui.Element {
	t := c.Theme()
	size := max(32, t.Rem(2))
	var b ui.Element
	switch style {
	case IconButtonToolbar:
		b = Button(c, "", ButtonToolbar).Label(label).Size(size, size).Padding(0)
		b.Children(func() { Icon(c, icon).Size(t.Rem(1.125), t.Rem(1.125)) })
	default:
		b = ui.ButtonBase(c).Size(size, size).Radius(t.Radius).TextColor(t.TextMuted).Label(label)
		// Paint after properties such as Disabled have been applied by the caller.
		b.Draw(func(p *ui.Painter, r ui.Rect) {
			if b.IsDisabled() {
				return
			}
			if b.Pressed() {
				p.Fill(r, t.Text.Alpha(0.14), t.Radius)
			} else if b.Hovered() {
				p.Fill(r, t.Text.Alpha(0.08), t.Radius)
			}
		}).Children(func() {
			Icon(c, icon).Size(t.Rem(1.25), t.Rem(1.25))
		})
	}
	Tooltip(c, b, label)
	return b
}
