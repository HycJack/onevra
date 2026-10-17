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

// controlSkeleton shows placeholders; the Content variant swaps one for a
// profile, with the loading state in the switch rather than background work.
func controlSkeleton(c *ui.Context, s *controlState, variant string) {
	t := c.Theme()
	switch variant {
	case "Avatar":
		components.Skeleton(c, components.SkeletonOptions{Lines: 2, Avatar: true})
		return
	case "Content":
	default:
		components.Skeleton(c, components.SkeletonOptions{})
		return
	}
	loading := !s.on
	components.Switch(c, &loading, tr(c, "loading"))
	s.on = !loading
	components.SkeletonContent(c, loading, components.SkeletonOptions{Lines: 2, Avatar: true}, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(t.Space(4)).Children(func() {
			ui.Box(c).Size(t.Space(9), t.Space(9)).Shrink(0).Radius(t.Space(4.5)).Background(t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				components.Text(c, tr(c, "skeletonSampleName")).Bold()
				components.Text(c, tr(c, "skeletonSampleBio")).TextColor(t.TextMuted)
			})
		})
	})
}
