package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTextArea(c *ui.Context, s *controlState, variant string) {
	components.TextArea(c, &s.text).Label("Notes").Placeholder("Write a note").ReadOnly(variant == "Read only").Height(100).FillWidth().MinWidth(0)
}
