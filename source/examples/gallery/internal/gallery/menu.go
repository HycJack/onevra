package gallery

import (
	"runtime"

	"github.com/egoist/mygo"
)

// ApplicationMenu returns the gallery's macOS menu bar. The entry point owns
// installation and application lifetime. Other platforms retain their defaults.
func (g *Gallery) ApplicationMenu() *mygo.Menu {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu, Submenu: []*mygo.MenuItem{
			{Role: mygo.RoleAbout},
			mygo.Separator(),
			{ID: "settings", Label: "Settings…", Accelerator: "Cmd+,", Click: func(*mygo.MenuItem, *mygo.Window) {
				g.openSettingsFromMenu()
			}},
			mygo.Separator(),
			{Role: mygo.RoleServices},
			mygo.Separator(),
			{Role: mygo.RoleHide},
			{Role: mygo.RoleHideOthers},
			{Role: mygo.RoleUnhide},
			mygo.Separator(),
			{Role: mygo.RoleQuit},
		}},
		{Role: mygo.RoleFileMenu},
		{Role: mygo.RoleEditMenu},
		{Role: mygo.RoleViewMenu},
		{Role: mygo.RoleWindowMenu},
	})
}

// Application and Dock menus target their owner, not a possibly unrelated
// focused window. Repeated activation keeps the existing settings state.
func (g *Gallery) openSettingsFromMenu() {
	g.enterGallery()
	g.settings.open = true
	g.showFromDock()
}
