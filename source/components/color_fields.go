package components

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

type colorFieldState struct {
	start, end int
}

func colorField(c *ui.Context, st *colorViewState, value *ui.Color, i int, label string, set func(ui.Color, colorFormat)) {
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
		editor := func() ui.Element {
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

			return in
		}
		// Hex is text notation, not an adjustable packed integer.
		if i == 3 {
			editor()
			return
		}
		current, hi := st.channel(i, *value)
		precision := 1.0
		if i == 4 || st.format != colorRGBA {
			precision = .1
		}
		n := numberInput(c, &current, 0, hi, 1, NumberInputOptions{LeadingText: label, HideStepper: true}, editor, 1, precision)
		n.Label(names[i]).FillWidth().MinWidth(0)
		if n.Changed() {
			next, source := st.setChannel(i, current, *value)
			set(next, source)
		}
	})
}
