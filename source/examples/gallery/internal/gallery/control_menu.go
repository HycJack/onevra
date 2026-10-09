package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlMenu(c *ui.Context, s *controlState, variant string) {
	build := func(m *ui.Menu) {
		for _, action := range []string{"Open", "Duplicate"} {
			if m.Item(action).Chosen() {
				s.text = action
				s.index++
			}
		}
		m.Separator()
		if m.Item("Delete").Chosen() {
			s.text = "Delete"
			s.index++
		}
	}
	if variant == "Right-click" {
		components.Button(c, "Right-click here", components.ButtonOutline).ContextMenu(build)
	} else {
		components.Button(c, "Actions", components.ButtonSecondary).Menu(build)
	}
	if s.index > 0 {
		ui.Textf(c, "%s: %d", s.text, s.index)
	} else {
		components.Text(c, "No action yet")
	}
}
