package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlScrollView(c *ui.Context, _ *controlState, variant string) {
	options := components.ScrollViewOptions{}
	switch variant {
	case "Horizontal":
		options.Axes = components.ScrollHorizontal
	case "Both":
		options.Axes = components.ScrollBoth
	}
	view := components.ScrollView(c, options).FillWidth().MinWidth(0).Height(160).
		Background(c.Theme().Background).Border(1, c.Theme().Border)
	if variant == "Insets" {
		view.ScrollbarInsets(32, 4, 24, 4)
	}
	view.Children(func() {
		if options.Axes == components.ScrollHorizontal {
			for i := 1; i <= 8; i++ {
				components.Textf(c, "Column %d", i).Width(100).Shrink(0).Padding(12)
			}
			return
		}
		content := ui.Column(c).Shrink(0).Gap(8).Padding(12)
		if options.Axes == components.ScrollBoth {
			content.Width(480)
		} else {
			content.FillWidth()
		}
		content.Children(func() {
			for i := 1; i <= 12; i++ {
				components.Textf(c, "Row %d", i).Height(28).Shrink(0)
			}
		})
	})
}
