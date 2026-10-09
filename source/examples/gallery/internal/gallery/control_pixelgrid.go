package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlPixelGrid(c *ui.Context, _ *controlState, variant string) {
	majorEvery := -1
	if variant == "Major lines" {
		majorEvery = 4
	}
	components.PixelGrid(c, components.PixelGridOptions{Spacing: 16, MajorEvery: majorEvery}).
		FillWidth().MinWidth(0).Height(160).
		Radius(c.Theme().Radius).Border(1, c.Theme().Border)
}
