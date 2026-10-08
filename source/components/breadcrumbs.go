package components

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
)

// BreadcrumbsControl retains activation separately from Element styling.
type BreadcrumbsControl struct {
	ui.Element
	changed bool
}

// Changed reports activation of an ancestor, even if chosen already held its index.
func (b *BreadcrumbsControl) Changed() bool { return b.changed }

// Breadcrumbs displays a path with interactive ancestors and a non-interactive
// final item. Activation writes the ancestor index to chosen; applications own
// navigation and the path itself. A nil chosen binding is invalid.
// Keep the control to query Changed after applying Element layout methods.
func Breadcrumbs(c *ui.Context, items []string, chosen *int) *BreadcrumbsControl {
	if chosen == nil {
		panic("components.Breadcrumbs: nil chosen")
	}
	t := c.Theme()
	b := &BreadcrumbsControl{Element: ui.Row(c).MinWidth(0).AlignItems(ui.Center).Gap(t.Space(.5)).Role(ui.RoleGroup)}
	b.Element.Children(func() {
		for i, item := range items {
			if i > 0 {
				Icon(c, lucide.ChevronRight).Size(t.Space(3), t.Space(3)).Shrink(0).TextColor(t.TextMuted).Role(ui.RoleNone)
			}
			if i == len(items)-1 {
				Text(c, item).SingleLine().MinWidth(0).FontWeight(500)
				continue
			}
			ui.Box(c).Key(i).MinWidth(0).Children(func() {
				link := Button(c, item, ButtonGhost).Role(ui.RoleLink).MinWidth(0).
					Padding(t.Space(.5), t.Space(1.5)).TextColor(t.TextMuted)
				if link.Clicked() {
					*chosen, b.changed = i, true
				}
			})
		}
	})
	return b
}
