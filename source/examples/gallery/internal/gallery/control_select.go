package gallery

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlSelect(c *ui.Context, s *controlState, variant string) {
	if variant == "NativeSelect" {
		components.NativeSelect(c, &s.choice, []components.NativeSelectOption[string]{
			{Value: "Small", Label: tr(c, "small")}, {Value: "Medium", Label: tr(c, "medium")}, {Value: "Large", Label: tr(c, "large")},
		}).Trigger.Label(tr(c, "chooseSize")).AlignSelf(ui.Start)
		return
	}
	t := c.Theme()
	selectParts := components.Select(c, &s.choice)
	selectParts.Trigger.Label(tr(c, "size")).Children(func() { components.Text(c, tr(c, s.choice)).SingleLine() })
	selectParts.Popup(func(panel ui.Element) {
		searchable := variant == "Searchable"
		if searchable {
			panel.MinWidth(t.Rem(14))
			selectParts.SearchField(&s.query).Label(tr(c, "searchSizes")).Placeholder(tr(c, "searchSizes2"))
		}
		matches := 0
		for _, option := range []string{"Small", "Medium", "Large"} {
			if searchable && !strings.Contains(strings.ToLower(option+" "+tr(c, option)), strings.ToLower(s.query)) {
				continue
			}
			matches++
			item := selectParts.Item(option).Label(tr(c, option))
			item.Children(func() { components.Text(c, tr(c, option)).SingleLine() })
		}
		if matches == 0 {
			components.Text(c, tr(c, "noMatchingSizes")).Padding(t.Space(2)).TextColor(t.TextMuted)
		}
	})
}
