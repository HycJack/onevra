package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTimePicker(c *ui.Context, s *controlState, _ string) {
	components.TimePicker(c, &s.day).Label(tr(c, "gallery.time"))
}
