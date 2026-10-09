package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlColorView(c *ui.Context, s *controlState, _ string) {
	components.ColorView(c, &s.color).Label(tr(c, "gallery.color"))
}

func controlColorPicker(c *ui.Context, s *controlState, variant string) {
	p := components.ColorPicker(c, &s.color)
	p.Label(fmt.Sprintf(tr(c, "gallery.s_color"), tr(c, variant)))
	switch variant {
	case "Outlined":
		p.Size(72, 36).Border(2, c.Theme().Accent).Radius(2).Padding(4)
	case "Rounded":
		p.Size(64, 40).Border(0, ui.Color{}).Radius(16).Padding(0)
	}
}
