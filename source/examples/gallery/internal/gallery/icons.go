package gallery

import (
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
	ui.Text(c, label).FontSize(c.Theme().Rem(.8)).Bold().TextColor(c.Theme().TextMuted).Padding(4, 8, 4, 8)
}

func (g *Gallery) navigationItem(c *ui.Context, label string, icon *ui.SVG, page galleryPage) {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(label).Label(label).FillWidth().Justify(ui.Start).Padding(7, 8).Gap(8).Radius(t.Radius)
	if g.page == page {
		b.Background(t.Text.Alpha(.1))
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
	ui.Column(c).Key("lucide-browser").Fill().Children(func() {
		ui.Row(c).FillWidth().Shrink(0).Padding(t.Space(3), t.Space(4)).Gap(8).Children(func() {
			ui.Box(c).Grow(1).MinWidth(0).Children(func() {
				if ui.TextInput(c, &g.query).Label("Search icons").Placeholder("Search…").FillWidth().
					Padding(t.Space(1.5), t.Space(2.5), t.Space(1.5), t.Space(2.5)+18+8).Changed() {
					g.grid.ScrollTo(0, ui.Start)
					g.selectedIcon = -1
				}
				ui.Row(c).Absolute().Left(t.Space(2.5)).Top(0).Bottom(0).Width(18).Center().PassThrough().Children(func() {
					ui.Icon(c, lucide.Search()).Size(18, 18).TextColor(t.TextMuted).PassThrough()
				})
			})
			if g.query != "" && ui.Button(c, "Clear").Clicked() {
				g.query = ""
				g.grid.ScrollTo(0, ui.Start)
				g.selectedIcon = -1
			}
		})
		icons := filterIcons(g.query)
		if len(icons) == 0 {
			ui.Text(c, "No icons found. Try another name.").TextColor(t.TextMuted).Wrap().Padding(t.Space(4))
			return
		}
		g.grid.Key = func(i int) any { return icons[i].name }
		g.grid.Label = func(i int) string { return icons[i].name }
		g.grid.Selected = &g.selectedIcon
		previewSize := t.Rem(6)
		rowHeight := previewSize + t.Rem(1.5) + t.Space(7)
		ui.GridView(c, &g.grid, len(icons), t.Rem(8), rowHeight, func(i int) {
			icon := icons[i]
			ui.Column(c).Fill().Padding(t.Space(2)).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(previewSize, previewSize).Shrink(0).Center().Children(func() {
					ui.Icon(c, icon.get()).Size(28, 28)
				})
				ui.Text(c, icon.name).FillWidth().MinWidth(0).FontSize(t.Rem(.8)).TextAlign(ui.Center).SingleLine()
			})
		}).Label("Icon grid").FillWidth().Grow(1).MinHeight(0).Padding(t.Space(4))
	})
}
