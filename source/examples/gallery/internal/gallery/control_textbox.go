package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func controlTextBox(c *ui.Context, s *controlState, variant string) {
	options := components.TextBoxOptions{ReadOnly: variant == "Read only"}
	label, placeholder := "Name", "Your name"
	if variant == "Search" {
		options.LeadingIcon = lucide.Search()
		label, placeholder = "Search samples", "Search…"
	}
	box := components.TextBox(c, &s.text, options).Label(tr(c, label)).Placeholder(tr(c, placeholder)).FillWidth().MinWidth(0)
	if variant == "Filled" {
		// The anchor stays where a selection started, after the caret when
		// it was made backwards.
		anchor, caret := box.TextAnchorCaret()
		components.Text(c, msg(c).TextBoxAnchorCaret(anchor, caret)).TextColor(c.Theme().TextMuted)
	}
}
