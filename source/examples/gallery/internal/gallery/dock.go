package gallery

import (
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

const dockPage galleryPage = trayPage + 1

// Dock menus belong to the application, not to a window. Restore the previous
// menu on cleanup, but do not overwrite a menu another owner installed later.
type dockDemo struct {
	menu, previous *mygo.Menu
}

func (d *dockDemo) close() {
	if d.menu != nil && mygo.App.Dock.Menu() == d.menu {
		mygo.App.Dock.SetMenu(d.previous)
	}
	d.menu, d.previous = nil, nil
}

func (g *Gallery) showFromDock() {
	if g.Shell == nil {
		return
	}
	win := g.Shell.Window()
	if win.IsMinimized() {
		win.Restore()
	}
	win.Show()
	win.Focus()
}

func (g *Gallery) dockMenu() *mygo.Menu {
	return mygo.NewMenu([]*mygo.MenuItem{
		{ID: "show-gallery", Label: g.text("gallery.show_gallery"), Click: func(*mygo.MenuItem, *mygo.Window) {
			g.showFromDock()
		}},
		mygo.Separator(),
		{ID: "welcome", Label: g.text("gallery.open_welcome"), Click: func(*mygo.MenuItem, *mygo.Window) {
			g.page = welcomePage
			g.enterGallery()
			g.settings.open = false
			g.showFromDock()
		}},
		{ID: "settings", Label: g.text("gallery.open_settings"), Click: func(*mygo.MenuItem, *mygo.Window) {
			g.openSettingsFromMenu()
		}},
	})
}

func (g *Gallery) enableDockMenu() {
	if runtime.GOOS != "darwin" || g.Shell == nil || g.dock.menu != nil {
		return
	}
	g.dock.previous = mygo.App.Dock.Menu()
	g.dock.menu = g.dockMenu()
	mygo.App.Dock.SetMenu(g.dock.menu)
}

// DockMenu demonstrates MyGo's native Dock API without changing the normal
// application menu or replacing the operating system's standard Dock items.
func (g *Gallery) DockMenu(c *ui.Context) {
	g.withLanguage(c, func() { g.dockMenuLocalized(c) })
}

func (g *Gallery) dockMenuLocalized(c *ui.Context) {
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(c.Theme().Space(6)).Gap(16).Children(func() {
			components.Text(c, tr(c, "gallery.dock_menu")).FontSize(22).Bold()
			if runtime.GOOS != "darwin" {
				components.Text(c, tr(c, "gallery.dock_menus_are_available_only_on_macos_use_desktop_tray_for_a_tray_men")).Wrap()
				return
			}
			components.Text(c, tr(c, "gallery.enable_the_menu_then_right_click_gallery_s_dock_icon_show_gallery_open")).Wrap()
			components.Text(c, tr(c, "gallery.the_menu_stays_available_when_you_navigate_away_disable_it_here_or_clo")).Wrap().TextColor(c.Theme().TextMuted)
			if g.dock.menu == nil {
				components.Button(c, tr(c, "gallery.enable_dock_menu"), components.ButtonPrimary).
					Disabled(g.Shell == nil).OnClick(g.enableDockMenu)
			} else {
				components.Button(c, tr(c, "gallery.disable_dock_menu"), components.ButtonOutline).OnClick(g.dock.close)
			}
		})
	})
}
