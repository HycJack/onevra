package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSlider(c *ui.Context, s *controlState, variant string) {
	opts := components.SliderOptions{}
	if variant == "Step" || variant == "Stepped range" {
		opts.Step = 20
	}
	if variant == "Range" || variant == "Stepped range" {
		opts.Range = &s.high
	}
	components.Slider(c, &s.value, 0, 100, opts).Label(tr(c, "gallery.volume")).FillWidth().MinWidth(0)
}
