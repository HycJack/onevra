package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlPopover(c *ui.Context, s *controlState, variant string) {
	anchor := components.Button(c, tr(c, "showPopover"), components.ButtonOutline)
	if anchor.Clicked() {
		s.dialogOpen = !s.dialogOpen
	}
	components.Popover(c, anchor, &s.dialogOpen, func(panel ui.Element) {
		panel.Label(tr(c, "popoverContent"))
		components.Text(c, tr(c, "popoverContent")).Bold()
		components.Text(c, tr(c, "chooseAValue")).TextColor(c.Theme().TextMuted)
		components.Checkbox(c, &s.on, tr(c, "enableNotifications"))
		if components.Button(c, tr(c, "close"), components.ButtonSecondary).Clicked() {
			s.dialogOpen = false
		}
	})
}
