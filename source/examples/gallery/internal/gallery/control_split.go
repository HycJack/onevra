package gallery

import (
	"fmt"
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
				components.Text(c, tr(c, label)).Wrap()
			})
		}
	}
	components.Split(c, &s.split, pane(first), pane(second), options).
		Label(fmt.Sprintf(tr(c, "gallery.s_divider"), tr(c, variant))).FillWidth().Height(240)
}
