package oneogo

import "github.com/egoist/mygo/ui"

// EmptyStateOptions supplies presentation only. Actions builds native controls;
// the caller handles their events and decides whether to replace the empty state.
type EmptyStateOptions struct {
	Title       string
	Description string
	Actions     func()
}

// EmptyState renders a quiet, leading-aligned explanation and optional actions.
// It does not impose a page width or take ownership of the surrounding layout.
func EmptyState(c *ui.Context, options EmptyStateOptions) *ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(3)).Role(ui.RoleGroup).Label(options.Title).Children(func() {
		ui.Text(c, options.Title).FontSize(t.Rem(1.25)).Bold()
		if options.Description != "" {
			ui.Text(c, options.Description).TextColor(t.TextMuted)
		}
		if options.Actions != nil {
			options.Actions()
		}
	})
}
