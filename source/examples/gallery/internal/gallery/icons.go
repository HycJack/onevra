package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

type galleryPage int

const welcomePage galleryPage = -1

const (
	lucidePage galleryPage = iota
	shellPage
)

type iconEntry struct {
	name string
	get  func() *ui.SVG
}

func (g *Gallery) navigationItem(c *ui.Context, entry navigationEntry) {
	selected := g.page == entry.page || entry.mygo && isMyGoPage(g.page)
	if components.SidebarItem(c, tr(c, entry.label), entry.icon(), selected).Key(entry.page).Clicked() {
		g.page = entry.page
		if entry.mygo {
			g.page = g.lastMyGoPage()
		}
	}
}

func (g *Gallery) Icons(c *ui.Context) {
	g.withLanguage(c, func() { g.iconsLocalized(c) })
}

func (g *Gallery) iconsLocalized(c *ui.Context) {
	t := c.Theme()
	g.grid.Key = func(i int) any { return iconCatalog[i].name }
	buttonSize := t.Rem(4)
	components.GridView(c, &g.grid, len(iconCatalog), buttonSize, buttonSize, func(i int) {
		icon := iconCatalog[i]
		ui.Box(c).Fill().Center().Children(func() {
			b := components.Button(c, "", components.ButtonOutline).Label(icon.name).Size(buttonSize, buttonSize).Padding(t.Space(2))
			b.Children(func() {
				components.Icon(c, icon.get).Size(20, 20)
			})
			if b.Clicked() {
				c.WriteClipboard(icon.name)
			}
		})
	}).Label(tr(c, "iconGrid")).Fill().Padding(t.Space(4))
}
