package components

import "github.com/egoist/mygo/ui"

// Sidebar creates a transparent column of SidebarSections. Place it inside a
// ui.Scroll when scrolling is needed. Applications own its width, titlebar
// inset and selection; it can be used independently of Shell.
func Sidebar(c *ui.Context) *ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().MinWidth(0).
		Padding(0, t.Space(3), t.Space(3), t.Space(3)).Gap(t.Space(3))
}

// SidebarSection builds a muted heading followed by application-owned items.
// Sections are always expanded; fn may be nil for a heading alone.
func SidebarSection(c *ui.Context, label string, fn func()) *ui.Element {
	t := c.Theme()
	section := ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(1))
	section.Children(func() {
		Text(c, label).FontSize(t.Rem(.8)).Bold().TextColor(t.TextMuted).
			Padding(t.Space(1), t.Space(2), t.Space(1), t.Space(2))
		if fn != nil {
			fn()
		}
	})
	return section
}

// SidebarItem creates a navigation button with a neutral selected face and an
// optional accent-colored icon. Hover does not change its face. The caller
// supplies selection and handles Clicked; MyGo provides Tab focus and activation
// with Enter or Space, not tree-style arrow navigation. Key can distinguish
// items with identical labels. A nil icon leaves no icon space.
func SidebarItem(c *ui.Context, label string, icon *ui.SVG, selected bool) *ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(label).Label(label).FillWidth().Justify(ui.Start).
		Padding(t.Space(1.75), t.Space(2)).Gap(t.Space(2)).Radius(t.Radius).
		Selected(selected).Background(ui.Color{}).TextColor(t.Text)
	if selected {
		b.Background(t.Text.Alpha(.1))
	}
	b.Children(func() {
		if icon != nil {
			Icon(c, icon).Size(t.Space(4.5), t.Space(4.5)).TextColor(t.Accent)
		}
		Text(c, label).SingleLine()
	})
	return dimmed(b)
}
