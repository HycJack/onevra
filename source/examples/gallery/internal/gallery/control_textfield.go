package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTextField(c *ui.Context, s *controlState, variant string) {
	components.TextField(c, &s.text).Label("Name").Placeholder("Your name").ReadOnly(variant == "Read only").FillWidth().MinWidth(0)
}
