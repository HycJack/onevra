package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// controlScrollEdgeBar docks a toolbar and a bottom bar over a scrolling list
// in one style, uncolored on bare glass; the bottom bar picks the section the
// toolbar names.
func controlScrollEdgeBar(c *ui.Context, s *controlState, variant string) {
	t := c.Theme()
	style := components.ScrollEdgeGlass
	switch variant {
	case "Soft edge":
		style = components.ScrollEdgeSoft
	case "Hard edge":
		style = components.ScrollEdgeHard
	case "Progressive blur":
		style = components.ScrollEdgeBlur
	}
	// Over a glass bar the buttons stay plain; otherwise each carries its own face.
	face := components.IconButtonToolbar
	if style == components.ScrollEdgeGlass {
		face = components.IconButtonPlain
	}
	top := components.ScrollEdgeBarOptions{Style: style}
	top.Uncolored = variant == "Uncolored bar"
	bottom := top
	bottom.Bottom = true
	inset := components.ScrollEdgeBarExtent(c, top)
	sections := []glassAction{{lucide.House, tr(c, "home")}, {lucide.Search, tr(c, "search")}, {lucide.Settings, tr(c, "settings")}}
	s.index = min(max(0, s.index), len(sections)-1)
	ui.Box(c).FillWidth().MinWidth(0).Height(320).Radius(t.Radius*2).Clip().
		Background(t.Background).Border(1, t.Border).Children(func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Padding(inset+8, 16).
			ScrollbarInsets(inset, 0).Gap(10).Children(func() {
			for i := range 12 {
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					a, b := glassTiles[i%len(glassTiles)], glassTiles[(i+2)%len(glassTiles)]
					ui.Box(c).Size(56, 40).Radius(8).Shrink(0).Gradient(a, b, 135)
					components.Text(c, msg(c).RowD(i+1))
				})
			}
		})
		components.ScrollEdgeBar(c, top, func() {
			components.IconButton(c, lucide.ChevronLeft(), tr(c, "stepperBack"), face)
			components.Text(c, sections[s.index].label).Bold().Grow(1)
			components.IconButton(c, lucide.Share(), tr(c, "share"), face)
		}).Label(tr(c, "toolbar"))
		components.ScrollEdgeBar(c, bottom, func() {
			ui.Box(c).Grow(1)
			for i, section := range sections {
				if components.IconButton(c, section.icon(), section.label, face).Clicked() {
					s.index = i
				}
			}
			ui.Box(c).Grow(1)
		}).Label(tr(c, "glassBottomBar"))
	})
}
