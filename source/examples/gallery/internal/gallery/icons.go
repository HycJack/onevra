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

func (g *Gallery) navigationGroup(c *ui.Context, label string) {
	components.Text(c, label).FontSize(c.Theme().Rem(.8)).Bold().TextColor(c.Theme().TextMuted).Padding(4, 8, 4, 8)
}

func (g *Gallery) navigationItem(c *ui.Context, label string, icon *ui.SVG, page galleryPage) {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(label).Label(label).FillWidth().Justify(ui.Start).Padding(7, 8).Gap(8).Radius(t.Radius)
	if g.page == page {
		b.Background(t.Text.Alpha(.1))
	}
	b.Children(func() {
		components.Icon(c, icon).Size(18, 18).TextColor(t.Accent)
		components.Text(c, label).SingleLine()
	})
	if b.Clicked() {
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
