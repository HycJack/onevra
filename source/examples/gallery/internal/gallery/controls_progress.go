package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSpinner(c *ui.Context, _ *controlState, variant string) {
	style := components.SpinnerRing
	if variant == "Spokes" {
		style = components.SpinnerSpokes
	}
	components.Spinner(c, style).Label(variant + " loading")
}

func controlProgress(c *ui.Context, s *controlState, variant string) {
	components.Slider(c, &s.value, 0, 100).Label(variant + " value").FillWidth()
	components.Textf(c, "%.0f%%", s.value)
	build := components.Progress
	if variant == "ProgressRing" {
		build = components.ProgressRing
	}
	build(c, s.value/100).Label(variant + " completion")
	components.Text(c, "Indeterminate").TextColor(c.Theme().TextMuted)
	build(c, -1).Label(variant + " activity")
}
