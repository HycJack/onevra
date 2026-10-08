package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTabs(c *ui.Context, s *controlState, _ string) {
	components.Tabs(c, &s.index, "General", "Style")
}
