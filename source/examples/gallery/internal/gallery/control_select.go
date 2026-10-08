package gallery

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSelect(c *ui.Context, s *controlState, variant string) {
	if variant == "NativeSelect" {
		components.NativeSelect(c, &s.choice, []components.NativeSelectOption[string]{
			{Value: "Small", Label: "Small"}, {Value: "Medium", Label: "Medium"}, {Value: "Large", Label: "Large"},
		}).Trigger.Label("Choose size").AlignSelf(ui.Start)
		return
	}
	t := c.Theme()
	selectParts := components.Select(c, &s.choice)
	selectParts.Trigger.Label("Size").Children(func() { components.Text(c, s.choice).SingleLine() })
	selectParts.Popup(func(panel *ui.Element) {
		searchable := variant == "Searchable"
		if searchable {
			panel.MinWidth(t.Rem(14))
			selectParts.SearchField(&s.query).Label("Search sizes").Placeholder("Search sizes…")
		}
		matches := 0
		for _, option := range []string{"Small", "Medium", "Large"} {
			if searchable && !strings.Contains(strings.ToLower(option), strings.ToLower(s.query)) {
				continue
			}
			matches++
			item := selectParts.Item(option).Label(option)
			item.Children(func() { components.Text(c, option).SingleLine() })
		}
		if matches == 0 {
			components.Text(c, "No matching sizes").Padding(t.Space(2)).TextColor(t.TextMuted)
		}
	})
}
