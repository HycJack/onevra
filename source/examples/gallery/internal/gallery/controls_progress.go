package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSpinner(c *ui.Context, _ *controlState, variant string) {
	style := components.SpinnerRing
	if variant == "Spokes" {
		style = components.SpinnerSpokes
	}
	components.Spinner(c, style).Label(fmt.Sprintf(tr(c, "gallery.s_loading"), tr(c, variant)))
}

func controlProgress(c *ui.Context, s *controlState, variant string) {
	components.Slider(c, &s.value, 0, 100).Label(fmt.Sprintf(tr(c, "gallery.s_value"), tr(c, variant))).FillWidth()
	ui.Textf(c, tr(c, "gallery.0f"), s.value)
	build := components.Progress
	if variant == "ProgressRing" {
		build = components.ProgressRing
	}
	build(c, s.value/100).Label(fmt.Sprintf(tr(c, "gallery.s_completion"), tr(c, variant)))
	components.Text(c, tr(c, "gallery.indeterminate")).TextColor(c.Theme().TextMuted)
	build(c, -1).Label(fmt.Sprintf(tr(c, "gallery.s_activity"), tr(c, variant)))
}
