package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSpinner(c *ui.Context, _ *controlState, variant string) {
	style := components.SpinnerRing
	if variant == "Spokes" || variant == "Tinted" {
		style = components.SpinnerSpokes
	}
	e := components.Spinner(c, style).Label(msg(c).LoadingOf(tr(c, variant)))
	if variant == "Tinted" {
		e.TextColor(c.Theme().Accent)
	}
}

func controlProgress(c *ui.Context, s *controlState, variant string) {
	components.Slider(c, &s.value, 0, 100).Label(msg(c).ValueOf(tr(c, variant))).FillWidth()
	ui.Text(c, msg(c).Percent(s.value))
	build := components.Progress
	if variant == "ProgressRing" {
		build = components.ProgressRing
	}
	build(c, s.value/100).Label(msg(c).CompletionOf(tr(c, variant)))
	components.Text(c, tr(c, "indeterminate")).TextColor(c.Theme().TextMuted)
	build(c, -1).Label(msg(c).ActivityOf(tr(c, variant)))
}
