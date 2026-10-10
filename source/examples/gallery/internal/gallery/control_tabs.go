package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlTabs(c *ui.Context, s *controlState, _ string) {
	components.Tabs(c, &s.index, tr(c, "gallery.general"), tr(c, "gallery.style"))
}
