package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func controlSidebar(c *ui.Context, s *controlState, _ string) {
	components.Sidebar(c).Children(func() {
		components.SidebarSection(c, tr(c, "gallery.locations"), func() {
			if components.SidebarItem(c, tr(c, "gallery.home"), lucide.House(), s.index == 0).Key("home").Clicked() {
				s.index = 0
			}
			if components.SidebarItem(c, tr(c, "gallery.files"), lucide.Folder(), s.index == 1).Key("files").Clicked() {
				s.index = 1
			}
		})
		components.SidebarSection(c, tr(c, "gallery.account"), func() {
			if components.SidebarItem(c, tr(c, "gallery.preferences"), nil, s.index == 2).Key("preferences").Clicked() {
				s.index = 2
			}
		})
	})
}
