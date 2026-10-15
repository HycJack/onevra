package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlAccordion(c *ui.Context, s *controlState, _ string) {
	components.Accordion(c, func(group *components.AccordionGroup) {
		group.Item(tr(c, "general"), &s.toggles[0], func() { components.Checkbox(c, &s.on, tr(c, "notifications")) })
		group.Item(tr(c, "about"), &s.toggles[1], func() { components.Text(c, tr(c, "builtWithOnevra")) })
	})
}

func controlCollapsible(c *ui.Context, s *controlState, variant string) {
	position := components.ArrowRight
	if variant == "Left arrow" {
		position = components.ArrowLeft
	}
	components.Collapsible(c, tr(c, "details"), &s.on, components.CollapsibleOptions{ArrowPosition: position}, func() {
		components.Checkbox(c, &s.toggles[0], tr(c, "notifications"))
	})
}
