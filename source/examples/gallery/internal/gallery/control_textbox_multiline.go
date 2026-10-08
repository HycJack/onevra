package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTextBoxMultiline(c *ui.Context, s *controlState, variant string) {
	components.TextBox(c, &s.text, components.TextBoxOptions{Multiline: true, ReadOnly: variant == "Read only"}).Label("Notes").Placeholder("Write a note").Height(100).FillWidth().MinWidth(0)
}
