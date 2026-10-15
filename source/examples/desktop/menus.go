package main

import (
	"github.com/egoist/mygo"
	"runtime"
)

func (a *application) action(id string, fn func()) *mygo.MenuItem {
	return &mygo.MenuItem{ID: id, Label: a.text(id), Click: func(*mygo.MenuItem, *mygo.Window) { fn() }}
}

func (a *application) installMenus() {
	settings := a.action("settings", a.openSettings)
	settings.Accelerator = "CmdOrCtrl+,"
	quit := &mygo.MenuItem{ID: "quit", Role: mygo.RoleQuit}
	if runtime.GOOS == "darwin" {
		mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
			{Role: mygo.RoleAppMenu, Submenu: []*mygo.MenuItem{settings, mygo.Separator(), quit}},
			{Role: mygo.RoleEditMenu},
		}))
		mygo.App.Dock.SetMenu(mygo.NewMenu([]*mygo.MenuItem{a.action("open", a.reopen), a.action("settings", a.openSettings)}))
	}
	if a.tray != nil {
		items := []*mygo.MenuItem{}
		if a.config.Panel {
			items = append(items, a.action("openPanel", a.showPanel))
		}
		if appKind == "hybrid" {
			items = append(items, a.action("open", a.showMain))
		}
		items = append(items, a.action("settings", a.openSettings), mygo.Separator(), a.action("quit", mygo.App.Quit))
		a.trayMenu = mygo.NewMenu(items)
		// AppIndicator only reports menu actions, not left/right button events.
		if runtime.GOOS == "linux" {
			a.tray.SetMenu(a.trayMenu)
		}
	}
}
