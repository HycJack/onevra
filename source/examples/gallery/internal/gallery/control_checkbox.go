package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlCheckbox(c *ui.Context, s *controlState, _ string) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4)).Children(func() {
		for i, label := range []string{"Unchecked", "Checked"} {
			ui.Box(c).Key(label).Children(func() {
				components.Checkbox(c, &s.toggles[i], label)
			})
		}
	})
}
