package components

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
)

// CalendarControl exposes layout and frame-local date changes. Keep this handle
// when applying Element layout methods, which return ui.Element.
type CalendarControl struct {
	ui.Element
	changed bool
	chosen  bool
}

// Changed reports a selected date changed by input, not cursor navigation.
func (c *CalendarControl) Changed() bool { return c.changed }

func calendarDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func calendarMonth(d time.Time, delta int) time.Time {
	first := time.Date(d.Year(), d.Month()+time.Month(delta), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1).Day()
	return first.AddDate(0, 0, min(d.Day(), last)-1)
}

func calendarSetDay(value *time.Time, day time.Time) bool {
	next := time.Date(day.Year(), day.Month(), day.Day(), value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), value.Location())
	if next.Equal(*value) {
		return false
	}
	*value = next
	return true
}

type calendarState struct{ cursor, last time.Time }

// Calendar builds the fixed-size, Monday-first six-week view used inside
// CalendarPicker. Arrows move its cursor; Page Up/Down change months with day
// clamping; Home/End move to month boundaries. Enter/Space or a day click selects
// a date, preserving its time and location. Navigation alone does not select.
// The view keeps the same size inside CalendarPicker and when displayed directly.
// Key or disable a parent before construction.
func Calendar(c *ui.Context, value *time.Time) *CalendarControl {
	if value == nil {
		panic("components.Calendar: nil date")
	}
	t := c.Theme()
	out := &CalendarControl{Element: ui.Column(c).Width(t.Space(56)).Height(t.Space(62)).Shrink(0).Gap(t.Space(.5)).Focusable().Role(ui.RoleTable)}
	st := ui.Local(out.Element, "calendar", func() calendarState { return calendarState{calendarDay(*value), *value} })
	if st.last != *value {
		st.cursor, st.last = calendarDay(*value), *value
	}
	choose := func(d time.Time) {
		out.chosen = true
		out.changed = calendarSetDay(value, d)
		st.cursor, st.last = d, *value
	}
	if !out.IsDisabled() {
		switch {
		case out.Shortcut(0, ui.KeyLeft):
			st.cursor = st.cursor.AddDate(0, 0, -1)
		case out.Shortcut(0, ui.KeyRight):
			st.cursor = st.cursor.AddDate(0, 0, 1)
		case out.Shortcut(0, ui.KeyUp):
			st.cursor = st.cursor.AddDate(0, 0, -7)
		case out.Shortcut(0, ui.KeyDown):
			st.cursor = st.cursor.AddDate(0, 0, 7)
		case out.Shortcut(0, ui.KeyPageUp):
			st.cursor = calendarMonth(st.cursor, -1)
		case out.Shortcut(0, ui.KeyPageDown):
			st.cursor = calendarMonth(st.cursor, 1)
		case out.Shortcut(0, ui.KeyHome):
			st.cursor = st.cursor.AddDate(0, 0, 1-st.cursor.Day())
		case out.Shortcut(0, ui.KeyEnd):
			st.cursor = time.Date(st.cursor.Year(), st.cursor.Month()+1, 0, 0, 0, 0, 0, time.UTC)
		case out.Shortcut(0, ui.KeyEnter), out.Shortcut(0, ui.KeySpace):
			choose(st.cursor)
		}
	}
	out.Children(func() {
		ui.Row(c).FillWidth().Height(t.Space(8)).Margin(0, 0, t.Space(2), 0).AlignItems(ui.Center).Children(func() {
			nav := func(label string, delta int, icon *ui.SVG) {
				b := Button(c, "", ButtonGhost).Size(t.Space(8), t.Space(8)).Padding(t.Space(2)).Label(label)
				b.Children(func() { Icon(c, icon) })
				if !b.IsDisabled() && b.Clicked() {
					st.cursor = calendarMonth(st.cursor, delta)
					out.Focus()
				}
			}
			nav("Previous month", -1, lucide.ChevronLeft())
			Text(c, st.cursor.Format("January 2006")).SingleLine().Bold().Grow(1).TextAlign(ui.Center)
			nav("Next month", 1, lucide.ChevronRight())
		})
		month := time.Date(st.cursor.Year(), st.cursor.Month(), 1, 0, 0, 0, 0, time.UTC)
		out.Label(month.Format("January 2006")).Value(st.cursor.Format("January 2, 2006"))
		ui.Row(c).Height(t.Space(4)).Margin(0, 0, t.Space(2), 0).Children(func() {
			for _, label := range []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"} {
				Text(c, label).Width(t.Space(8)).TextAlign(ui.Center).TextColor(t.TextMuted).Role(ui.RoleColumnHeader)
			}
		})
		start := month.AddDate(0, 0, -((int(month.Weekday()) + 6) % 7))
		today := calendarDay(time.Now())
		for week := range 6 {
			ui.Row(c).Height(t.Space(7)).Role(ui.RoleRow).Children(func() {
				for col := range 7 {
					day := start.AddDate(0, 0, week*7+col)
					selected := day.Equal(calendarDay(*value))
					variant := ButtonGhost
					if selected {
						variant = ButtonPrimary
					}
					cell := Button(c, fmt.Sprint(day.Day()), variant).Size(t.Space(8), t.Space(7)).Padding(0).Label(day.Format("January 2, 2006")).Checked(selected)
					if !selected {
						if day.Month() != month.Month() {
							cell.TextColor(t.TextMuted)
						}
						if day.Equal(st.cursor) {
							cell.Background(t.SurfaceHover)
						}
						if day.Equal(today) {
							cell.Border(1, t.Accent)
						}
					}
					if day.Equal(st.cursor) {
						out.ActiveDescendant(cell)
					}
					if !cell.IsDisabled() && cell.Clicked() {
						choose(day)
					}
				}
			})
		}
	})
	dimmed(out.Element)
	return out
}

// CalendarPicker displays a leading calendar icon and an ISO date. Click,
// Enter or Space opens Calendar in a dismissible Popover; choosing a date closes
// it and restores trigger focus, including when choosing the existing date.
// Escape/outside press cancel navigation. The application owns value; selection
// preserves its time of day and location. Changed excludes opening/navigation.
// Key or disable a parent before construction; Label names the trigger.
func CalendarPicker(c *ui.Context, value *time.Time) *CalendarControl {
	if value == nil {
		panic("components.CalendarPicker: nil date")
	}
	t := c.Theme()
	out := &CalendarControl{Element: ui.ButtonBase(c).MinWidth(t.Space(36)).Justify(ui.Start).Role(ui.RolePopUpButton)}
	border := t.Border
	if !out.IsDisabled() {
		switch {
		case out.Focused():
			border = t.Accent
		case out.Hovered():
			border = t.Border.Mix(t.Text, .25)
		}
	}
	out.Padding(t.Space(1.5), t.Space(2.5), t.Space(1.5), t.Space(8)).
		Radius(t.Radius).Background(t.Surface).Border(1, border).TextColor(t.Text)
	out.DrawOver(func(p *ui.Painter, r ui.Rect) {
		size := t.Space(4)
		p.Icon(lucide.CalendarDays(), ui.Rect{
			X: r.X + t.Space(2.5), Y: r.Y + (r.H-size)/2, W: size, H: size,
		}, t.TextMuted)
	})
	dimmed(out.Element)
	open := ui.Local(out.Element, "calendar-open", func() bool { return false })
	if out.IsDisabled() {
		*open = false
	} else if out.Clicked() {
		*open = !*open
	}
	Popover(c, out.Element, open, func(panel ui.Element) {
		panel.Padding(t.Space(3)).Gap(0)
		view := Calendar(c, value)
		view.AutoFocus()
		out.changed = view.changed
		if view.chosen {
			*open = false
			out.Focus()
		}
	})
	out.Children(func() {
		Text(c, value.Format("2006-01-02")).SingleLine().FontFeatures("tnum")
	})
	out.Expanded(*open).Value(value.Format("2006-01-02"))
	return out
}
