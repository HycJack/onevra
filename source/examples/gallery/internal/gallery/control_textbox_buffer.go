package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

func controlTextBoxBuffer(c *ui.Context, s *controlState, variant string) {
	options := components.TextBoxOptions{ReadOnly: variant == "Read only"}
	label, placeholder := "Name", "Your name"
	if variant == "Search" {
		options.LeadingIcon = lucide.Search()
		label, placeholder = "Search samples", "Search…"
	}
	components.TextBoxBuffer(c, s.buffer, options).Label(label).Placeholder(placeholder).FillWidth().MinWidth(0)
}

func controlTextBoxBufferMultiline(c *ui.Context, s *controlState, variant string) {
	components.TextBoxBuffer(c, s.buffer, components.TextBoxOptions{Multiline: true, ReadOnly: variant == "Read only"}).
		Label("Notes").Placeholder("Write a note").Height(100).FillWidth().MinWidth(0)
}
