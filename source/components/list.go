package components

import "github.com/egoist/mygo/ui"

// List creates a virtualized list with application-owned state and row content.
// Give it a height or Fill within a bounded parent. MyGo owns scrolling,
// selection, keyboard navigation and accessibility. Configure state.Selected
// for single selection or state.Selection and state.Key for keyed selection;
// nil state uses MyGo's local scroll state without selection.
//
// Rows use the Sidebar's neutral selection tint and have no hover face. Row
// builders inherit the original frame theme and supply their own content and
// spacing. They may run during layout: do not perform application work in them.
func List(c *ui.Context, state *ui.ListState, count int, row func(int)) ui.Element {
	original := c.Theme()
	// MyGo has no ListBase. Scope its row-face tokens to this list, restoring
	// the application theme for content, including rows built during layout.
	face := *original
	face.Accent = original.Text.Alpha(.1)
	face.AccentText = original.Text
	face.SurfaceHover = ui.Color{}
	c.SetTheme(&face)
	defer c.SetTheme(original)
	return ui.List(c, state, count, func(i int) {
		current := c.Theme()
		c.SetTheme(original)
		defer c.SetTheme(current)
		if row != nil {
			row(i)
		}
	}).TextColor(original.Text)
}
