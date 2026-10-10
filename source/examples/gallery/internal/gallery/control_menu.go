package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlMenu(c *ui.Context, s *controlState, variant string) {
	labels := map[string]string{"Open": tr(c, "gallery.open"), "Duplicate": tr(c, "gallery.duplicate"), "Delete": tr(c, "gallery.delete")}
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
		components.Button(c, tr(c, "gallery.right_click_here"), components.ButtonOutline).ContextMenu(build)
	} else {
		components.Button(c, tr(c, "gallery.actions"), components.ButtonSecondary).Menu(build)
	}
	if s.index > 0 {
		ui.Textf(c, tr(c, "gallery.s_d"), labels[s.text], s.index)
	} else {
		components.Text(c, tr(c, "gallery.no_action_yet"))
	}
}
