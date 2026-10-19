package main

import (
	"log"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

// All application state and native callbacks are owned by the UI thread.
type application struct {
	config                                     config
	path, note                                 string
	preferred                                  []string // The user's languages, most preferred first.
	errorKey, errorDetail                      string
	ready, pendingOpen, quitting, settingsOpen bool
	main, panel, settings, welcome             *mygo.Window
	tray                                       *mygo.Tray
	trayMenu                                   *mygo.Menu
	settingsState                              components.SettingsState
	panelTab                                   int
}

func (a *application) report(key string, err error) {
	a.errorKey, a.errorDetail = key, err.Error()
	log.Printf("%s: %v", key, err)
}

func (a *application) start() {
	a.followSystemLanguages()
	if appKind != "desktop" {
		if err := a.createTray(); err != nil {
			a.report("trayError", err)
		}
	}
	a.ready = true
	a.refresh()
	switch {
	case appKind == "menubar":
		if !a.config.Welcomed {
			a.showWelcome()
		} else if a.tray == nil {
			a.openSettings()
		}
	default:
		a.showMain()
	}
	if a.pendingOpen {
		a.pendingOpen = false
		a.reopen()
	}
}

func show(win *mygo.Window) {
	if win.IsMinimized() {
		win.Restore()
	}
	win.Show()
	win.Focus()
}

func (a *application) showMain() {
	if appKind == "menubar" {
		a.openSettings()
		return
	}
	if a.main == nil {
		a.main = a.newWindow(plain.ShellOptions{
			Title: a.text("title"), Width: 820, Height: 740, MinWidth: 820, MinHeight: 740,
			Hidden: true, Content: a.mainView,
		})
		a.main.OnClose(a.closeMain)
	}
	if runtime.GOOS == "darwin" {
		mygo.App.Dock.Show()
	}
	show(a.main)
}

func (a *application) closeMain(e *mygo.CloseEvent) {
	if a.quitting {
		return
	}
	e.PreventDefault()
	if !keepRunning(a.config, runtime.GOOS, a.tray != nil) {
		mygo.App.Quit()
		return
	}
	a.settingsOpen = false
	a.main.Hide()
	a.hidePanel()
	a.updateDock()
}

func (a *application) updateDock() {
	if runtime.GOOS != "darwin" {
		return
	}
	if appKind == "desktop" || appKind == "menubar" && a.tray == nil {
		mygo.App.Dock.Show()
	} else {
		mygo.App.Dock.Hide()
	}
}

// Replace this action to open a document or another application-owned window.
func (a *application) reopen() {
	if a.quitting {
		return
	}
	if appKind != "menubar" {
		a.showMain()
		return
	}
	if !a.config.Welcomed {
		a.showWelcome()
		return
	}
	if a.settings != nil && a.settings.IsVisible() {
		show(a.settings)
		return
	}
	if a.config.Panel && a.tray != nil {
		a.showPanel()
	} else {
		a.openSettings()
	}
}

// Settings is an action, not a required window type. Desktop embeds the
// same settings surface in a Dialog; the windowless example gives it a window.
func (a *application) openSettings() {
	a.hidePanel()
	if appKind != "menubar" {
		a.settingsOpen = true
		a.showMain()
		a.main.Invalidate()
		return
	}
	if a.settings == nil {
		a.settings = a.newWindow(plain.ShellOptions{
			Title: a.text("settings"), Width: 680, Height: 540, MinWidth: 360, MinHeight: 400,
			Hidden: true, Content: func(c *ui.Context) {
				a.scope(c, func() { w, _ := c.Size(); a.settingsView(c, w) })
			},
		})
		a.settings.OnClose(func(e *mygo.CloseEvent) {
			if a.quitting {
				return
			}
			e.PreventDefault()
			if a.tray == nil {
				mygo.App.Quit()
			} else {
				a.settings.Hide()
			}
		})
	}
	a.updateDock()
	show(a.settings)
}

func (a *application) showWelcome() {
	if a.welcome == nil {
		a.welcome = a.newWindow(plain.ShellOptions{
			Title: a.text("welcome"), Width: 460, Height: 340, MinWidth: 360, MinHeight: 300,
			Hidden: true, Content: a.welcomeView, HeaderActions: a.headerActions,
		})
		a.welcome.OnClose(func(e *mygo.CloseEvent) {
			if a.quitting {
				return
			}
			e.PreventDefault()
			a.finishWelcome()
		})
	}
	a.updateDock()
	show(a.welcome)
}

func (a *application) finishWelcome() {
	a.config.Welcomed = true
	a.changed()
	a.welcome.Hide()
	if a.tray == nil {
		a.openSettings()
	}
}

func (a *application) changed() {
	if a.path != "" {
		if err := saveConfig(a.path, a.config); err != nil {
			a.report("saveError", err)
		} else if a.errorKey == "saveError" || a.errorKey == "loadError" {
			a.errorKey, a.errorDetail = "", ""
		}
	}
	if !a.config.Panel {
		a.hidePanel()
	}
	if a.ready {
		a.refresh()
		a.updateDock()
	}
}

func (a *application) refresh() {
	mygo.Theme.SetSource(a.themeSource())
	a.applyLanguage()
	a.installMenus()
	for _, item := range []struct {
		win   *mygo.Window
		title string
	}{
		{a.main, "title"}, {a.panel, "panelTitle"}, {a.settings, "settings"}, {a.welcome, "welcome"},
	} {
		if item.win != nil {
			item.win.SetTitle(a.text(item.title))
			item.win.Invalidate()
		}
	}
	if a.tray != nil {
		a.tray.SetToolTip(a.text("title"))
	}
}

func (a *application) newWindow(options plain.ShellOptions) *mygo.Window {
	options.Theme = theme.Onevra
	shell := plain.NewShell(options)
	return shell.Window()
}
