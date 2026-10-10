package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlPopover(c *ui.Context, s *controlState, variant string) {
	anchor := components.Button(c, tr(c, "gallery.show_popover"), components.ButtonOutline)
	if anchor.Clicked() {
		s.dialogOpen = !s.dialogOpen
	}
	components.Popover(c, anchor, &s.dialogOpen, func(panel ui.Element) {
		panel.Label(tr(c, "gallery.popover_content"))
		components.Text(c, tr(c, "gallery.popover_content")).Bold()
		components.Text(c, tr(c, "gallery.choose_a_value")).TextColor(c.Theme().TextMuted)
		components.Checkbox(c, &s.on, tr(c, "gallery.enable_notifications"))
		if components.Button(c, tr(c, "gallery.close"), components.ButtonSecondary).Clicked() {
			s.dialogOpen = false
		}
	})
}
