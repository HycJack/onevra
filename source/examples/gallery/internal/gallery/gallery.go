// Package gallery holds the standalone showcase, not public component API.
package gallery

import (
	"fmt"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
	"github.com/oneo-me/onevra/source/theme"
)

// Gallery owns demonstration state and application preferences.
type Gallery struct {
	commands       commandState
	settings       settingsState
	tray           trayDemo
	dock           dockDemo
	notification   notificationDemo
	language       *translator
	systemLanguage string
	config         Config
	configPath     string
	Clicks         int
	Shell          *sidebar.Shell
	page           galleryPage
	grid           components.GridState
	mygoStates     map[string]*demoState
	controlStates  map[string]*controlState
	toasts         components.ToastOptions
	// Native startup supplies the factory; headless views never load libraries.
	createTerminal func() (*terminal.Terminal, error)
	createWebView  func() (*mygo.WebView, error)
	audioOutput    bool
	gamepads       bool
	// stateKey asks MyGo to remember the window; empty keeps it disk-free.
	stateKey string
	openURL  func(string) error
}

// New creates the standalone shell demonstration.
func (g *Gallery) New() *sidebar.Shell {
	return g.newWindow(false)
}

func (g *Gallery) newWindow(hidden bool) *sidebar.Shell {
	g.createTerminal = newShellTerminal
	g.createWebView = g.newWebView
	g.audioOutput, g.gamepads = true, true
	g.openURL = mygo.Shell.OpenExternal
	g.Shell = sidebar.NewShell(sidebar.ShellOptions{
		Title:     "Onevra Gallery",
		Hidden:    hidden,
		StateKey:  g.stateKey,
		LeftWidth: 220, LeftMinWidth: 220, LeftMaxWidth: 520,
		Left: g.Left, Content: g.View, LeftToolbar: g.LeftToolbar,
		ContentTitle: g.Titlebar,
		Theme:        theme.Onevra,
	})
	g.updateLanguageWindows()
	g.Shell.Window().OnClose(func(*mygo.CloseEvent) { g.commands.palette.ReleaseInputSource() })
	g.Shell.Window().OnClosed(g.Close)
	return g.Shell
}

func (g *Gallery) Titlebar(c *ui.Context) {
	g.withLanguage(c, func() { g.titlebarLocalized(c) })
}

func (g *Gallery) titlebarLocalized(c *ui.Context) {
	title := "Shell"
	if g.page == welcomePage {
		title = "Welcome"
	}
	if g.page == lucidePage {
		title = "Icon library / Lucide"
	}
	if g.page == trayPage {
		title = "Desktop / Tray"
	}
	if g.page == notificationsPage {
		title = "Desktop / System notifications"
	}
	if g.page == dockPage {
		title = "Desktop / Dock menu"
	}
	if example, ok := mygoExampleFor(g.page); ok {
		title = tr(c, example.group) + " / " + tr(c, example.name)
	}
	if example, ok := controlExampleFor(g.page); ok {
		title = tr(c, example.group) + " / " + tr(c, example.name)
	}
	components.Text(c, tr(c, title)).Bold().SingleLine()
}

func (g *Gallery) Left(c *ui.Context) {
	g.withLanguage(c, func() { g.leftLocalized(c) })
}

func (g *Gallery) leftLocalized(c *ui.Context) {
	ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() { g.navigation(c) })
}

func (g *Gallery) navigation(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).FillWidth().Shrink(0).Padding(0, t.Space(3), t.Space(2), t.Space(3)).Children(func() {
		_, shortcut := commandShortcut()
		trigger := components.Button(c, "", components.ButtonSecondary).
			Label(tr(c, "gallery.search_components")).FillWidth().MinWidth(0).
			Padding(t.Space(1.5), t.Space(2.5)).Gap(t.Space(2)).TextColor(t.TextMuted)
		tint, borderAlpha := float32(.08), float32(.18)
		if trigger.Pressed() {
			tint, borderAlpha = .18, .4
		} else if trigger.Hovered() {
			tint, borderAlpha = .12, .3
		}
		trigger.Background(t.Background.Mix(t.Text, tint).Alpha(.3)).Border(1, t.Text.Alpha(borderAlpha))
		trigger.Children(func() {
			components.Icon(c, lucide.Search).Size(t.Space(4), t.Space(4)).Shrink(0)
			components.Text(c, tr(c, "gallery.search_2")).Grow(1).MinWidth(0).SingleLine()
			components.Text(c, shortcut).FontSize(t.Rem(.8)).Shrink(0).SingleLine()
		})
		if trigger.Clicked() {
			g.openCommands()
		}
	})
	groups := navigationCatalog()
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		components.Sidebar(c).Children(func() {
			for _, group := range groups {
				components.SidebarSection(c, tr(c, group.label), func() {
					for _, entry := range group.entries {
						g.navigationItem(c, entry.label, entry.icon(), entry.page)
					}
				})
			}
		})
	})
}

func (g *Gallery) LeftToolbar(c *ui.Context) {
	g.withLanguage(c, func() { g.leftToolbarLocalized(c) })
}

func (g *Gallery) leftToolbarLocalized(c *ui.Context) {
	if sidebar.IconButton(c, lucide.Settings(), tr(c, "gallery.settings")).Clicked() {
		g.settings.open = true
	}
	ui.Spacer(c)
	sidebar.IconButton(c, themeIcon(), tr(c, "gallery.theme")).Menu(g.themeMenu)
}

// iconSamples uses only named accessors, not a full-library runtime registry.
func iconSamples(c *ui.Context) {
	components.Text(c, tr(c, "gallery.lucide_icons")).Bold()
	components.Text(c, tr(c, "gallery.original_strokes_at_16_24_and_32_dip_colors_inherit_from_mygo")).Wrap()
	for _, sample := range []struct {
		name string
		get  func() *ui.SVG
	}{
		{"Search", lucide.Search},
		{"Settings", lucide.Settings},
		{"Cloud rain", lucide.CloudRain},
		{"Fingerprint", lucide.FingerprintPattern},
		{"Accessibility", lucide.Accessibility},
	} {
		ui.Column(c).Gap(4).Children(func() {
			components.Text(c, tr(c, sample.name))
			ui.Row(c).Gap(12).Children(func() {
				for _, size := range []float32{16, 24, 32} {
					components.Icon(c, sample.get).Size(size, size)
				}
			})
		})
	}
	ui.Row(c).Gap(12).TextColor(c.Theme().TextMuted).Children(func() {
		components.Icon(c, lucide.Sun)
		components.Icon(c, lucide.Moon)
		components.Text(c, tr(c, "gallery.inherited_muted_color"))
	})
}

// View builds the default page's main content.
func (g *Gallery) View(c *ui.Context) {
	g.withLanguage(c, func() { g.viewLocalized(c) })
}

func (g *Gallery) viewLocalized(c *ui.Context) {
	var content ui.Element
	_, control := controlExampleFor(g.page)
	example, upstream := mygoExampleFor(g.page)
	// The lawn game's background music belongs to its page alone.
	if example.name != "Lawn defense" {
		if s := g.mygoStates["Lawn defense"]; s != nil {
			s.lawn.stopMusic()
		}
	}
	if !g.config.DisablePixelGrid && (control || upstream && example.name != "Terminal") {
		content = components.PixelGrid(c, components.PixelGridOptions{Spacing: 22, MajorEvery: -1}).Column()
	} else {
		content = ui.Column(c)
	}
	content.Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0)
	if g.page == welcomePage {
		welcomeBackground(c, content)
	}
	content.Children(func() {
		if g.page == welcomePage {
			welcomeIconField(c)
		}
		g.body(c)
	})
	g.Settings(c)
	g.Commands(c)
}

func (g *Gallery) body(c *ui.Context) {
	if g.page == notificationsPage {
		g.Notifications(c)
		return
	}
	if g.page == dockPage {
		g.DockMenu(c)
		return
	}
	if g.page == trayPage {
		g.Tray(c)
		return
	}
	if g.page == welcomePage {
		g.Welcome(c)
		return
	}
	if example, ok := controlExampleFor(g.page); ok {
		g.Controls(c, example)
		return
	}
	if example, ok := mygoExampleFor(g.page); ok {
		g.MyGo(c, example)
		return
	}
	if g.page == lucidePage {
		g.Icons(c)
		return
	}
	t := c.Theme()
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(4)).Children(func() {
			components.Text(c, tr(c, "gallery.a_sidebar_and_content_one_native_window")).TextColor(t.TextMuted).Wrap()
			components.Text(c, tr(c, "gallery.toggle_or_resize_the_sidebar_inside_the_window_double_click_its_separa")).Wrap()
			if components.Button(c, tr(c, "gallery.count"), components.ButtonSecondary).Clicked() {
				g.Clicks++
			}
			components.Text(c, fmt.Sprintf(tr(c, "gallery.clicks_d"), g.Clicks))
			iconSamples(c)
		})
	})
}
