package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlCalendar(c *ui.Context, s *controlState, _ string) {
	components.Calendar(c, &s.day)
}

func controlCalendarPicker(c *ui.Context, s *controlState, _ string) {
	components.CalendarPicker(c, &s.day).Label(tr(c, "date"))
}
