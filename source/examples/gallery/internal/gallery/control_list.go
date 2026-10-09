package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlList(c *ui.Context, s *controlState, _ string) {
	s.list.Selected = &s.index
	components.List(c, &s.list, 1000, func(i int) {
		ui.Textf(c, "Item %d", i+1).Padding(c.Theme().Space(2))
	}).Height(240).FillWidth().Label("Items")
}
