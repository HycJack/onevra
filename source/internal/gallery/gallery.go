// Package gallery holds the standalone showcase, not public component API.
package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/components"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

// Gallery owns only demonstration state.
type Gallery struct {
	Clicks       int
	Shell        *components.Shell
	page         galleryPage
	query        string
	grid         ui.GridState
	selectedIcon int
	mygoStates   map[string]*demoState
}

// New creates the standalone shell demonstration.
func (g *Gallery) New() *components.Shell {
	g.Shell = components.NewShell(components.ShellOptions{
		Title: "oneo-go Gallery", Width: 820, Height: 900,
		MinWidth: 360, MinHeight: 480,
		Left: g.Left, Content: g.View, Right: g.Right, LeftToolbar: g.LeftToolbar,
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
	ui.Text(c, title).Bold().SingleLine()
}

func (g *Gallery) Left(c *ui.Context) {
	ui.Scroll(c).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(0, 12, 12, 12).Gap(4).Children(func() {
			g.navigationGroup(c, "Components")
			g.navigationItem(c, "Shell", lucide.PanelsTopLeft(), shellPage)
			ui.Box(c).Height(4)
			g.navigationGroup(c, "Icon library")
			g.navigationItem(c, "Lucide", lucide.Shapes(), lucidePage)
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
	ui.Spacer(c)
	if components.ToolbarIconButton(c, lucide.Settings(), "Settings").Clicked() && g.Shell != nil {
		g.Shell.SetPage(g.Settings)
	}
}

func (g *Gallery) Right(c *ui.Context) {
	ui.Column(c).Fill().Padding(16).Gap(12).Children(func() {
		ui.Text(c, "Inspector").Bold()
		ui.Text(c, "Opening this panel extends the window to the right.").Wrap()
	})
}

func (g *Gallery) Settings(c *ui.Context) {
	ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
		ui.Text(c, "Settings").FontSize(c.Theme().Rem(2)).Bold()
		ui.Text(c, "A replacement page outside the three-column layout.").Wrap()
		if ui.Button(c, "Back").Clicked() && g.Shell != nil {
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
			ui.Text(c, "Three columns, one native window.").TextColor(t.TextMuted).Wrap()
			ui.Text(c, "Toggle the sidebar inside the window, or extend it with the inspector. Open Settings to replace the entire page.").Wrap()
			if ui.Button(c, "Count").Clicked() {
				g.Clicks++
			}
			ui.Text(c, fmt.Sprintf("Clicks: %d", g.Clicks))
			iconSamples(c)
		})
	})
}
