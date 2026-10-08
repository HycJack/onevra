package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlPopover(c *ui.Context, s *controlState, variant string) {
	anchor := components.Button(c, "Show popover", components.ButtonOutline)
	if anchor.Clicked() {
		s.dialogOpen = !s.dialogOpen
	}
	components.Popover(c, anchor, &s.dialogOpen, func(panel ui.Element) {
		panel.Label("Popover content")
		components.Text(c, "Popover content").Bold()
		components.Text(c, "Choose a value").TextColor(c.Theme().TextMuted)
		components.Checkbox(c, &s.on, "Enable notifications")
		if components.Button(c, "Close", components.ButtonSecondary).Clicked() {
			s.dialogOpen = false
		}
	})
}
