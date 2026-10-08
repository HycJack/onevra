package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

type galleryPage int

const (
	lucidePage galleryPage = iota
	shellPage
)

type iconEntry struct {
	name string
	get  func() *ui.SVG
}

func (g *Gallery) navigationItem(c *ui.Context, label string, icon *ui.SVG, page galleryPage) {
	if components.SidebarItem(c, label, icon, g.page == page).Clicked() {
		g.page = page
	}
}

func (g *Gallery) Icons(c *ui.Context) {
	t := c.Theme()
	g.grid.Key = func(i int) any { return iconCatalog[i].name }
	buttonSize := t.Rem(4)
	// GridView gives each cell a stable key and a two-DIP border.
	ui.GridView(c, &g.grid, len(iconCatalog), buttonSize+4, buttonSize+4, func(i int) {
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
	}).Label("Icon grid").Fill().Padding(t.Space(4))
}
