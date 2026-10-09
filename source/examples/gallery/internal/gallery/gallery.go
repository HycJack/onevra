// Package gallery holds the standalone showcase, not public component API.
package gallery

import (
	"fmt"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/shell"
)

// Gallery owns demonstration state and application preferences.
type Gallery struct {
	settings        settingsState
	welcome         welcomeState
	tray            trayDemo
	dock            dockDemo
	config          Config
	configPath      string
	Clicks          int
	Shell           *shell.Shell
	page            galleryPage
	navigationQuery string
	grid            components.GridState
	mygoStates      map[string]*demoState
	controlStates   map[string]*controlState
	// Native startup supplies the factory; headless views never load libraries.
	createTerminal func() (*terminal.Terminal, error)
}

// New creates the standalone shell demonstration.
func (g *Gallery) New() *shell.Shell {
	g.createTerminal = newShellTerminal
	g.Shell = shell.NewShell(shell.ShellOptions{
		Title:     "Onevra Gallery",
		LeftWidth: 220, LeftMinWidth: 220, LeftMaxWidth: 520,
		Left: g.Left, Content: g.View, LeftToolbar: g.LeftToolbar,
		ContentTitle: g.Titlebar,
	})
	g.Shell.SetPage(g.WelcomePage)
	g.Shell.Window().OnClosed(g.Close)
	return g.Shell
}

func (g *Gallery) Titlebar(c *ui.Context) {
	title := "Shell"
	if g.page == welcomePage {
		title = "Welcome"
	}
	if g.page == lucidePage {
		title = "Icon library / Lucide"
	}
	if g.page == componentStatusPage {
		title = "Component status"
	}
	if g.page == trayPage {
		title = "Desktop / Tray"
	}
	if g.page == dockPage {
		title = "Desktop / Dock menu"
	}
	if example, ok := mygoExampleFor(g.page); ok {
		title = example.group + " / " + example.name
	}
	if example, ok := controlExampleFor(g.page); ok {
		title = example.group + " / " + example.name
	}
	components.Text(c, title).Bold().SingleLine()
}

func (g *Gallery) Left(c *ui.Context) {
	ui.Column(c).Fill().Padding(shell.HeaderHeight(c), 0, 0, 0).Children(func() { g.navigation(c) })
}

func (g *Gallery) navigation(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).FillWidth().Shrink(0).Padding(0, t.Space(3), t.Space(2), t.Space(3)).Children(func() {
		components.TextBox(c, &g.navigationQuery, components.TextBoxOptions{LeadingIcon: lucide.Search()}).
			Label("Search components").Placeholder("Search components…").FillWidth().MinWidth(0)
	})
	groups := filteredNavigation(g.navigationQuery)
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		components.Sidebar(c).Children(func() {
			for _, group := range groups {
				components.SidebarSection(c, group.label, func() {
					for _, entry := range group.entries {
						icon := lucide.Component()
						switch entry.page {
						case welcomePage:
							icon = lucide.House()
						case lucidePage:
							icon = lucide.Shapes()
						}
						g.navigationItem(c, entry.label, icon, entry.page)
					}
				})
			}
			if len(groups) == 0 {
				components.Text(c, "No matching components.").TextColor(t.TextMuted).Wrap()
			}
		})
	})
}

func (g *Gallery) LeftToolbar(c *ui.Context) {
	shell.IconButton(c, themeIcon(), "Theme").Menu(g.themeMenu)
	ui.Spacer(c)
	if shell.IconButton(c, lucide.Settings(), "Settings").Clicked() {
		g.settings.open = true
	}
}

// iconSamples uses only named accessors, not a full-library runtime registry.
func iconSamples(c *ui.Context) {
	components.Text(c, "Lucide icons").Bold()
	components.Text(c, "Original strokes at 16, 24 and 32 DIP; colors inherit from MyGo.").Wrap()
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
			components.Text(c, sample.name)
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
		components.Text(c, "Inherited muted color")
	})
}

// View builds the default page's main content.
func (g *Gallery) View(c *ui.Context) {
	var content ui.Element
	_, control := controlExampleFor(g.page)
	example, upstream := mygoExampleFor(g.page)
	if control || upstream && example.name != "Terminal" {
		content = components.PixelGrid(c, components.PixelGridOptions{Spacing: 22, MajorEvery: -1}).Column()
	} else {
		content = ui.Column(c)
	}
	content.Fill().Padding(shell.HeaderHeight(c), 0, 0, 0)
	if g.page == welcomePage {
		welcomeBackground(c, content)
	}
	content.Children(func() { g.body(c) })
	g.Settings(c)
}

func (g *Gallery) body(c *ui.Context) {
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
	if g.page == componentStatusPage {
		g.ComponentStatus(c)
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
			components.Text(c, "A sidebar and content, one native window.").TextColor(t.TextMuted).Wrap()
			components.Text(c, "Toggle or resize the sidebar inside the window. Double-click its separator to restore the default width. Open Settings in a dialog without leaving this page.").Wrap()
			if components.Button(c, "Count", components.ButtonSecondary).Clicked() {
				g.Clicks++
			}
			components.Text(c, fmt.Sprintf("Clicks: %d", g.Clicks))
			iconSamples(c)
		})
	})
}
