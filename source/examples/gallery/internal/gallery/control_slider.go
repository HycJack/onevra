package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSlider(c *ui.Context, s *controlState, _ string) {
	components.Slider(c, &s.value, 0, 100).Label("Volume").FillWidth().MinWidth(0)
}
