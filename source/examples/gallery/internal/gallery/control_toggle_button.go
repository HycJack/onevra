package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

func controlToggleButton(c *ui.Context, s *controlState) {
	ui.Row(c).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		components.ToggleButton(c, &s.on, tr(c, "gallery.bold"))
		components.ToggleButton(c, &s.toggles[1], "").Label(tr(c, "gallery.italic")).Children(func() {
			components.Icon(c, lucide.Italic).Size(c.Theme().Space(4), c.Theme().Space(4))
		})
	})
}
