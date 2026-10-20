package gallery

import "github.com/egoist/mygo/ui"

// Selection groups related examples without moving their implementations or
// sharing their state. Controls supplies a separate state for each section.
func controlSelection(c *ui.Context, s *controlState, component string) {
	switch component {
	case "Checkbox":
		controlCheckbox(c, s, component)
	case "Switch", "Compact switch":
		controlSwitch(c, s, component)
	case "Radio":
		controlRadio(c, s, component)
	}
}
