package gallery

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlSelect(c *ui.Context, s *controlState, variant string) {
	if variant == "NativeSelect" {
		components.NativeSelect(c, &s.choice, []components.NativeSelectOption[string]{
			{Value: "Small", Label: tr(c, "gallery.small")}, {Value: "Medium", Label: tr(c, "gallery.medium")}, {Value: "Large", Label: tr(c, "gallery.large")},
		}).Trigger.Label(tr(c, "gallery.choose_size")).AlignSelf(ui.Start)
		return
	}
	t := c.Theme()
	selectParts := components.Select(c, &s.choice)
	selectParts.Trigger.Label(tr(c, "gallery.size")).Children(func() { components.Text(c, tr(c, s.choice)).SingleLine() })
	selectParts.Popup(func(panel ui.Element) {
		searchable := variant == "Searchable"
		if searchable {
			panel.MinWidth(t.Rem(14))
			selectParts.SearchField(&s.query).Label(tr(c, "gallery.search_sizes")).Placeholder(tr(c, "gallery.search_sizes_2"))
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
			components.Text(c, tr(c, "gallery.no_matching_sizes")).Padding(t.Space(2)).TextColor(t.TextMuted)
		}
	})
}
