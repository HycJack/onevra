package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

func controlTextBox(c *ui.Context, s *controlState, variant string) {
	options := components.TextBoxOptions{ReadOnly: variant == "Read only"}
	label, placeholder := "Name", "Your name"
	if variant == "Search" {
		options.LeadingIcon = lucide.Search()
		label, placeholder = "Search samples", "Search…"
	}
	components.TextBox(c, &s.text, options).Label(label).Placeholder(placeholder).FillWidth().MinWidth(0)
}
