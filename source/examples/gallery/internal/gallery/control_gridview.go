package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlGridView(c *ui.Context, s *controlState, _ string) {
	s.grid.Selected = &s.index
	s.grid.Label = func(i int) string { return fmt.Sprintf("Item %d", i+1) }
	components.GridView(c, &s.grid, 1000, 90, 64, func(i int) {
		ui.Textf(c, "Item %d", i+1).Padding(c.Theme().Space(2))
	}).Height(240).FillWidth().Label("Grid items")
}
