package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSwitch(c *ui.Context, s *controlState, _ string) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4)).Children(func() {
		for i, label := range []string{"Off", "On"} {
			ui.Box(c).Key(label).Children(func() {
				components.Switch(c, &s.toggles[i], tr(c, label))
			})
		}
	})
}
