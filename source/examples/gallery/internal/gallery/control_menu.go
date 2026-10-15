package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlMenu(c *ui.Context, s *controlState, variant string) {
	labels := map[string]string{"Open": tr(c, "open"), "Duplicate": tr(c, "duplicate"), "Delete": tr(c, "delete")}
	build := func(m *ui.Menu) {
		for _, action := range []string{"Open", "Duplicate"} {
			if m.Item(labels[action]).Chosen() {
				s.text = action
				s.index++
			}
		}
		m.Separator()
		if m.Item(labels["Delete"]).Chosen() {
			s.text = "Delete"
			s.index++
		}
	}
	if variant == "Right-click" {
		components.Button(c, tr(c, "rightClickHere"), components.ButtonOutline).ContextMenu(build)
	} else {
		components.Button(c, tr(c, "actions"), components.ButtonSecondary).Menu(build)
	}
	if s.index > 0 {
		ui.Text(c, msg(c).LabelCount(labels[s.text], s.index))
	} else {
		components.Text(c, tr(c, "noActionYet"))
	}
}
