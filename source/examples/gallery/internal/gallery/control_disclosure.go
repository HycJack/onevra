package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlAccordion(c *ui.Context, s *controlState, _ string) {
	components.Accordion(c, func(group *components.AccordionGroup) {
		group.Item("General", &s.toggles[0], func() { components.Checkbox(c, &s.on, "Notifications") })
		group.Item("About", &s.toggles[1], func() { components.Text(c, "Built with Onevra") })
	})
}

func controlCollapsible(c *ui.Context, s *controlState, variant string) {
	position := components.ArrowRight
	if variant == "Left arrow" {
		position = components.ArrowLeft
	}
	components.Collapsible(c, "Details", &s.on, components.CollapsibleOptions{ArrowPosition: position}, func() {
		components.Checkbox(c, &s.toggles[0], "Notifications")
	})
}
