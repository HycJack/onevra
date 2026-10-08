package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSelect(c *ui.Context, s *controlState, _ string) {
	t := c.Theme()
	selectParts := components.Select(c, &s.choice)
	selectParts.Trigger.Label("Size").Children(func() { components.Text(c, s.choice).SingleLine() })
	selectParts.Popup(func(panel *ui.Element) {
		panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		for _, option := range []string{"Small", "Medium", "Large"} {
			item := selectParts.Item(option).Label(option).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
			if item.Highlighted() {
				item.Background(t.Accent).TextColor(t.AccentText)
			}
			item.Children(func() { components.Text(c, option).SingleLine() })
		}
	})
}
