package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlTabs(c *ui.Context, s *controlState, variant string) {
	tabs := components.TabsUnderline
	if variant == "Segmented control" {
		tabs = components.TabsSegmented
	}
	components.Tabs(c, &s.index, tabs, tr(c, "general"), tr(c, "style"))
}
