package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

func controlSidebar(c *ui.Context, s *controlState, _ string) {
	components.Sidebar(c).Children(func() {
		components.SidebarSection(c, "Locations", func() {
			if components.SidebarItem(c, "Home", lucide.House(), s.index == 0).Clicked() {
				s.index = 0
			}
			if components.SidebarItem(c, "Files", lucide.Folder(), s.index == 1).Clicked() {
				s.index = 1
			}
		})
		components.SidebarSection(c, "Account", func() {
			if components.SidebarItem(c, "Preferences", nil, s.index == 2).Clicked() {
				s.index = 2
			}
		})
	})
}
