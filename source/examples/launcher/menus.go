package main

import (
	"runtime"

	"github.com/egoist/mygo"
)

func (a *application) action(id string, fn func()) *mygo.MenuItem {
	return &mygo.MenuItem{ID: id, Label: a.text(id), Click: func(*mygo.MenuItem, *mygo.Window) { fn() }}
}

func (a *application) installMenus() {
	appSettings := a.action("settings", a.openSettings)
	appSettings.Accelerator = "CmdOrCtrl+,"
	if runtime.GOOS == "darwin" {
		mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
			{Role: mygo.RoleAppMenu, Submenu: []*mygo.MenuItem{appSettings, mygo.Separator(), {Role: mygo.RoleQuit}}},
			{Role: mygo.RoleEditMenu},
		}))
	}
	if a.tray == nil {
		return
	}
	open := a.action("openLauncher", a.show)
	open.Accelerator = a.accelerator
	traySettings := a.action("settings", a.openSettings)
	login := &mygo.MenuItem{ID: "login", Label: a.text("openAtLogin"), Type: mygo.MenuItemCheckbox, Checked: a.login,
		Click: func(item *mygo.MenuItem, _ *mygo.Window) {
			a.login = item.Checked
			if err := mygo.App.SetOpenAtLogin(item.Checked); err != nil {
				a.report("loginError", err)
			}
		}}
	a.trayMenu = mygo.NewMenu([]*mygo.MenuItem{
		open,
		mygo.Separator(),
		login,
		mygo.Separator(),
		traySettings,
		mygo.Separator(),
		{ID: "quit", Role: mygo.RoleQuit},
	})
	// AppIndicator only reports menu actions, not left/right button events.
	if runtime.GOOS == "linux" {
		a.tray.SetMenu(a.trayMenu)
	}
}
