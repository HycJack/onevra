package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlColorView(c *ui.Context, s *controlState, _ string) {
	components.ColorView(c, &s.color).Label("Color")
}

func controlColorPicker(c *ui.Context, s *controlState, variant string) {
	p := components.ColorPicker(c, &s.color)
	p.Label(variant + " color")
	switch variant {
	case "Outlined":
		p.Size(72, 36).Border(2, c.Theme().Accent).Radius(2).Padding(4)
	case "Rounded":
		p.Size(64, 40).Border(0, ui.Color{}).Radius(16).Padding(0)
	}
}
