package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

//go:generate go run ../cmd/galleryicons

type galleryPage int

const (
	shellPage galleryPage = iota
	lucidePage
)

type iconEntry struct {
	name string
	get  func() *ui.SVG
}

func (g *Gallery) navigationGroup(c *ui.Context, label string) {
	ui.Text(c, label).FontSize(c.Theme().Rem(.8)).Bold().TextColor(c.Theme().TextMuted).Padding(8, 8, 4, 8)
}

func (g *Gallery) navigationItem(c *ui.Context, label string, icon *ui.SVG, page galleryPage) {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(label).Label(label).FillWidth().Justify(ui.Start).Padding(7, 8).Gap(8).Radius(t.Radius)
	if g.page == page {
		b.Background(t.Background.Mix(t.Text, .1))
	} else if b.Hovered() || b.Pressed() {
		b.Background(t.Text.Alpha(.06))
	}
	b.Children(func() {
		ui.Icon(c, icon).Size(18, 18).TextColor(t.Accent)
		ui.Text(c, label).SingleLine()
	})
	if b.Clicked() {
		g.page = page
	}
}

func filterIcons(query string) []iconEntry {
	terms := strings.Fields(strings.ToLower(strings.ReplaceAll(query, "-", " ")))
	if len(terms) == 0 {
		return iconCatalog
	}
	var result []iconEntry
	for _, icon := range iconCatalog {
		name := strings.ReplaceAll(icon.name, "-", " ")
		match := true
		for _, term := range terms {
			if !strings.Contains(name, term) {
				match = false
				break
			}
		}
		if match {
			result = append(result, icon)
		}
	}
	return result
}

func (g *Gallery) Icons(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Key("lucide-browser").Fill().Padding(t.Space(6)).Gap(t.Space(3)).Children(func() {
		ui.Text(c, "Lucide").FontSize(t.Rem(2)).Bold()
		ui.Row(c).Gap(8).Children(func() {
			ui.Icon(c, lucide.Search()).Size(18, 18).TextColor(t.TextMuted)
			if ui.TextInput(c, &g.query).Label("Search icons").Placeholder("Search…").Grow(1).Changed() {
				g.grid.ScrollTo(0, ui.Start)
				g.selectedIcon = -1
			}
			if g.query != "" && ui.Button(c, "Clear").Clicked() {
				g.query = ""
				g.grid.ScrollTo(0, ui.Start)
				g.selectedIcon = -1
			}
		})
		icons := filterIcons(g.query)
		ui.Text(c, fmt.Sprintf("%d of %d icons", len(icons), len(iconCatalog))).TextColor(t.TextMuted)
		if len(icons) == 0 {
			ui.Text(c, "No icons found. Try another name.").TextColor(t.TextMuted).Wrap()
			return
		}
		g.grid.Key = func(i int) any { return icons[i].name }
		g.grid.Label = func(i int) string { return icons[i].name }
		g.grid.Selected = &g.selectedIcon
		ui.GridView(c, &g.grid, len(icons), t.Rem(8), t.Rem(7), func(i int) {
			icon := icons[i]
			ui.Column(c).Fill().Padding(8).Gap(8).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Icon(c, icon.get()).Size(28, 28)
				ui.Text(c, icon.name).FontSize(t.Rem(.8)).TextAlign(ui.Center).Wrap()
			})
		}).FillWidth().Grow(1).MinHeight(0)
	})
}
