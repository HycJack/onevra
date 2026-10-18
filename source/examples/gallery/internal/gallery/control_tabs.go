package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func controlTabs(c *ui.Context, s *controlState, variant string) {
	tabs := components.TabsUnderline
	switch variant {
	case "Segmented control":
		tabs = components.TabsSegmented
	case "Highlight":
		tabs = components.TabsHighlight
	case "Custom leading":
		// Each tab leads with an application-built icon.
		icons := []func() *ui.SVG{lucide.Settings, lucide.Palette}
		components.TabsWithLeading(c, &s.index, components.TabsSegmented, func(c *ui.Context, i int, size float32) {
			components.Icon(c, icons[i]).Size(size, size).Shrink(0)
		}, tr(c, "general"), tr(c, "style"))
		return
	}
	components.Tabs(c, &s.index, tabs, tr(c, "general"), tr(c, "style"))
}
