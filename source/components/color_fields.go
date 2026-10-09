package components

import (
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

type colorFieldState struct {
	start, end int
	held       bool
	scrub      float64
}

func colorField(c *ui.Context, st *colorViewState, value *ui.Color, i int, label string, set func(ui.Color, colorFormat)) {
	t := c.Theme()
	weight := float32(1)
	if i == 3 {
		weight = 2
	}
	box := ui.Box(c).MinWidth(0).Grow(weight).Width(0)
	edit := ui.Local(box, "color-field", func() colorFieldState { return colorFieldState{} })
	box.Children(func() {
		names := [5]string{"Red", "Green", "Blue", "Hex color", "Opacity percent"}
		if st.format == colorHSBA {
			names[0], names[1], names[2] = "Hue degrees", "Saturation percent", "Brightness percent"
		}
		if st.format == colorHSLA {
			names[0], names[1], names[2] = "Hue degrees", "Saturation percent", "Lightness percent"
		}
		previous := st.fields[i]
		in := TextBox(c, &st.fields[i], TextBoxOptions{LeadingText: label}).Label(names[i]).MinWidth(0).FillWidth()
		if in.Changed() && !in.IsDisabled() {
			draft := st.fields[i]
			if !st.acceptDraft(i, draft, *value) {
				st.fields[i] = previous
				in.SetTextSelection(edit.start, edit.end)
			} else {
				if i == 3 {
					if next, ok := parseColorHex(draft, value.A); ok {
						set(next, colorRGBA)
						st.fields[i] = draft
					}
				} else if n, err := strconv.ParseFloat(strings.TrimSuffix(draft, "%"), 64); err == nil {
					next, source := st.setChannel(i, n, *value)
					set(next, source)
					st.fields[i] = draft
				}
			}
		}
		if !in.Focused() || in.Submitted() {
			old := st.fields
			st.sync(*value)
			canonical := st.fields[i]
			st.fields = old
			st.fields[i] = canonical
		}
		edit.start, edit.end = in.TextSelection()
		// A sibling hit region covers only the prefix, not the text editor. MyGo
		// owns pointer capture; the editor retains ordinary selection and IME input.
		prefixWidth, _ := c.MeasureText(0, ui.Span{Text: label})
		grip := ui.Box(c).Absolute().Left(1).Top(1).Bottom(1).
			Width(t.Space(2.5) + prefixWidth + t.Space(2)).Cursor(ui.CursorResizeEW).
			Focusable().Role(ui.RoleSlider).Label("Adjust " + names[i])
		current, hi := st.channel(i, *value)
		grip.Range(0, hi, current).Step(1)
		dx, _, held := grip.Dragged()
		next := current
		if !grip.IsDisabled() {
			if held {
				grip.Focus()
				if !edit.held {
					edit.scrub = current
				}
				edit.scrub = max(0, min(hi, edit.scrub+float64(dx/t.Space(.25))))
				precision := 1.0
				if i == 4 || (i < 3 && st.format != colorRGBA) {
					precision = 10
				}
				if dx != 0 {
					next = math.Round(edit.scrub*precision) / precision
				}
			}
			switch {
			case grip.Shortcut(0, ui.KeyRight), grip.Shortcut(0, ui.KeyUp):
				next = current + 1
			case grip.Shortcut(0, ui.KeyLeft), grip.Shortcut(0, ui.KeyDown):
				next = current - 1
			case grip.Shortcut(0, ui.KeyHome):
				next = 0
			case grip.Shortcut(0, ui.KeyEnd):
				next = hi
			}
			next = max(0, min(hi, next))
			if next != current {
				color, source := st.setChannel(i, next, *value)
				set(color, source)
			}
		}
		edit.held = held && !grip.IsDisabled()
	})
}
