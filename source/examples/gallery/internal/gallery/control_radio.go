package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlRadio(c *ui.Context, s *controlState, _ string) {
	ui.RadioGroup(c, func() {
		for _, option := range []string{"Small", "Medium", "Large"} {
			components.Radio(c, &s.choice, option, tr(c, option))
		}
	}).Label(tr(c, "gallery.size")).Row().FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4))
}
