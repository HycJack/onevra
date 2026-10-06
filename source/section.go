package oneogo

import "github.com/egoist/mygo/ui"

// Section renders a heading, optional description, and caller-built content.
// The returned native element supports ordinary layout modifiers and Key.
// Content is built synchronously; nil leaves only the heading.
func Section(c *ui.Context, title, description string, content func()) *ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup).Label(title).Children(func() {
		ui.Text(c, title).Bold()
		if description != "" {
			ui.Text(c, description).TextColor(t.TextMuted)
		}
		if content != nil {
			content()
		}
	})
}
