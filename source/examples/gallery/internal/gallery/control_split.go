package gallery

import (
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSplit(c *ui.Context, s *controlState, variant string) {
	options, first, second := splitPanes(c, variant)
	// Double-clicking the divider restores the first pane's starting size.
	options.Default = 100
	components.Split(c, &s.split, first, second, options).
		Label(msg(c).DividerOf(tr(c, variant))).FillWidth().Height(240)
}

// controlSplitPercent keeps the panes' shares as the page resizes.
func controlSplitPercent(c *ui.Context, s *controlState, variant string) {
	options, first, second := splitPanes(c, variant)
	options.Default = 30
	components.SplitPercent(c, &s.splitPercent, first, second, options).
		Label(msg(c).DividerOf(tr(c, variant))).FillWidth().Height(240)
	components.Text(c, msg(c).SplitPercentShare(int(math.Round(float64(s.splitPercent))))).TextColor(c.Theme().TextMuted)
}

func splitPanes(c *ui.Context, variant string) (components.SplitOptions, func(), func()) {
	options := components.SplitOptions{}
	first, second := "Left pane", "Right pane"
	if variant == "Vertical" {
		options.Orientation = components.Vertical
		first, second = "Top pane", "Bottom pane"
	}
	pane := func(label string) func() {
		return func() {
			ui.Column(c).Fill().Padding(12).Children(func() {
				components.Text(c, tr(c, label)).Wrap()
			})
		}
	}
	return options, pane(first), pane(second)
}
