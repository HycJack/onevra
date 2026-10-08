package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlNumberInput(c *ui.Context, s *controlState, variant string) {
	step := 1.0
	if variant == "Decimal" {
		step = 0.1
	}
	components.NumberInput(c, &s.value, -100, 100, step).Label("Quantity").FillWidth()
}
