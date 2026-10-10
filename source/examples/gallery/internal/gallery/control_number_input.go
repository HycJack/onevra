package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlNumberInput(c *ui.Context, s *controlState, variant string) {
	step := 1.0
	if variant == "Decimal" {
		step = 0.1
	}
	options := components.NumberInputOptions{HideStepper: variant == "Without stepper"}
	if variant == "Scrub" {
		options.LeadingText = "X"
	}
	components.NumberInput(c, &s.value, -100, 100, step, options).Label(tr(c, "gallery.quantity")).FillWidth()
}
