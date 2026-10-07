// Package gallery holds the standalone showcase, not public component API.
package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
	oneoshell "github.com/oneo-me/oneo-go/shell"
	oneoui "github.com/oneo-me/oneo-go/ui"
)

// Gallery owns demonstration state and application preferences.
type Gallery struct {
	config        Config
	configPath    string
	Clicks        int
	Shell         *oneoshell.Shell
	page          galleryPage
	grid          ui.GridState
	mygoStates    map[string]*demoState
	controlStates map[string]*controlState
}

// New creates the standalone shell demonstration.
func (g *Gallery) New() *oneoshell.Shell {
	g.Shell = oneoshell.NewShell(oneoshell.ShellOptions{
		Title:     "oneo-go Gallery",
		LeftWidth: 220, LeftMinWidth: 220, LeftMaxWidth: 520,
		Left: g.Left, Content: g.View, LeftToolbar: g.LeftToolbar,
		ContentTitle: g.Titlebar,
	})
	g.Shell.Window().OnClosed(g.Close)
	return g.Shell
}

func (g *Gallery) Titlebar(c *ui.Context) {
	title := "Shell"
	if g.page == lucidePage {
		title = "Icon library / Lucide"
	}
	if example, ok := mygoExampleFor(g.page); ok {
		title = example.group + " / " + example.name
	}
	if example, ok := controlExampleFor(g.page); ok {
		title = "Components / " + example.name
	}
	ui.Text(c, title).Bold().SingleLine()
}

func (g *Gallery) Left(c *ui.Context) {
	ui.Column(c).Fill().Padding(oneoshell.HeaderHeight(c), 0, 0, 0).Children(func() { g.navigation(c) })
}

func (g *Gallery) navigation(c *ui.Context) {
	ui.Scroll(c).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(0, 12, 12, 12).Gap(4).Children(func() {
			g.navigationGroup(c, "Icon library")
			g.navigationItem(c, "Lucide", lucide.Shapes(), lucidePage)
			ui.Box(c).Height(4)
			g.navigationGroup(c, "Components")
			for i, example := range controlCatalog {
				g.navigationItem(c, example.label, lucide.Component(), controlPageStart()+galleryPage(i))
			}
			ui.Box(c).Height(4)
			g.navigationGroup(c, "Plugins")
			for i, example := range mygoCatalog {
				if example.group == "Plugins" {
					g.navigationItem(c, example.name, lucide.Component(), mygoPageStart+galleryPage(i))
				}
			}
			ui.Box(c).Height(4)
			g.navigationGroup(c, "MyGo components")
			for _, group := range mygoGroups {
				ui.Box(c).Height(4)
				g.navigationGroup(c, group)
				for i, example := range mygoCatalog {
					if example.group == group {
						g.navigationItem(c, example.name, lucide.Component(), mygoPageStart+galleryPage(i))
					}
				}
			}
		})
	})
}

func (g *Gallery) LeftToolbar(c *ui.Context) {
	oneoshell.IconButton(c, themeIcon(), "Theme").Menu(g.themeMenu)
	ui.Spacer(c)
	if oneoshell.IconButton(c, lucide.Settings(), "Settings").Clicked() && g.Shell != nil {
		g.Shell.SetPage(g.Settings)
	}
}

func (g *Gallery) Settings(c *ui.Context) {
	ui.Column(c).Fill().Padding(oneoshell.HeaderHeight(c)+24, 24, 24, 24).Gap(16).Children(func() {
		ui.Text(c, "Settings").FontSize(c.Theme().Rem(2)).Bold()
		ui.Text(c, "A replacement page outside the sidebar-and-content layout.").Wrap()
		if oneoui.Button(c, "Back", oneoui.ButtonSecondary).Clicked() && g.Shell != nil {
			g.Shell.SetPage(nil)
		}
	})
}

// iconSamples uses only named accessors, not a full-library runtime registry.
func iconSamples(c *ui.Context) {
	ui.Text(c, "Lucide icons").Bold()
	ui.Text(c, "Original strokes at 16, 24 and 32 DIP; colors inherit from MyGo.").Wrap()
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
			ui.Text(c, sample.name)
			ui.Row(c).Gap(12).Children(func() {
				for _, size := range []float32{16, 24, 32} {
					ui.Icon(c, sample.get()).Size(size, size)
				}
			})
		})
	}
	ui.Row(c).Gap(12).TextColor(c.Theme().TextMuted).Children(func() {
		ui.Icon(c, lucide.Sun())
		ui.Icon(c, lucide.Moon())
		ui.Text(c, "Inherited muted color")
	})
}

// View builds the default page's main content.
func (g *Gallery) View(c *ui.Context) {
	ui.Column(c).Fill().Padding(oneoshell.HeaderHeight(c), 0, 0, 0).Children(func() { g.body(c) })
}

func (g *Gallery) body(c *ui.Context) {
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
	ui.Scroll(c).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(4)).Children(func() {
			ui.Text(c, "A sidebar and content, one native window.").TextColor(t.TextMuted).Wrap()
			ui.Text(c, "Toggle or resize the sidebar inside the window. Double-click its separator to restore the default width. Open Settings to replace the entire page.").Wrap()
			if oneoui.Button(c, "Count", oneoui.ButtonSecondary).Clicked() {
				g.Clicks++
			}
			ui.Text(c, fmt.Sprintf("Clicks: %d", g.Clicks))
			iconSamples(c)
		})
	})
}
