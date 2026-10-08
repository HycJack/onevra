package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSplit(c *ui.Context, s *controlState, variant string) {
	options := components.SplitOptions{}
	first, second := "Left pane", "Right pane"
	if variant == "Vertical" {
		options.Orientation = components.Vertical
		first, second = "Top pane", "Bottom pane"
	}
	pane := func(label string) func() {
		return func() {
			ui.Column(c).Fill().Padding(12).Children(func() {
				components.Text(c, label).Wrap()
			})
		}
	}
	components.Split(c, &s.split, pane(first), pane(second), options).
		Label(variant + " divider").FillWidth().Height(240)
}
