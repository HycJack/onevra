package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlTabs(c *ui.Context, s *controlState, variant string) {
	tabs := components.TabsUnderline
	switch variant {
	case "Segmented control":
		tabs = components.TabsSegmented
	case "Highlight":
		tabs = components.TabsHighlight
	}
	components.Tabs(c, &s.index, tabs, tr(c, "general"), tr(c, "style"))
}
