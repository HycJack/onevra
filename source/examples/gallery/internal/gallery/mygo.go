package gallery

import "github.com/egoist/mygo/ui"

const mygoPageStart galleryPage = 2

type mygoExample struct {
	name, group string
	build       func(*ui.Context, *demoState)
}

var mygoGroups = []string{
	"Actions", "Selection", "Text input", "Date & color", "Collections",
	"Layout", "Containers", "Navigation", "Overlays", "Status", "Media & text",
}

var demoImage = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 180 100"><rect width="180" height="100" rx="12" fill="#dbeafe"/><circle cx="140" cy="28" r="15" fill="#fbbf24"/><path d="M0 100L55 25L110 100Z" fill="#2563eb"/><path d="M65 100L120 45L180 100Z" fill="#059669"/></svg>`))

func mygoExampleFor(page galleryPage) (mygoExample, bool) {
	i := int(page - mygoPageStart)
	if i < 0 || i >= len(mygoCatalog) {
		return mygoExample{}, false
	}
	return mygoCatalog[i], true
}

// Variants are starting states, not frozen mockups: enabled controls remain interactive.
func mygoVariants(name string) []string {
	switch name {
	case "Toggle", "ToggleGroup", "Checkbox", "Switch":
		return []string{"Off", "On", "Disabled"}
	case "CheckboxGroup":
		return []string{"Unchecked", "Mixed", "Checked", "Disabled"}
	case "TextInput", "TextArea":
		return []string{"Empty", "Filled", "Read only", "Disabled"}
	case "SearchField", "Autocomplete", "TokenField":
		return []string{"Empty", "Filled", "Disabled"}
	case "Field":
		return []string{"Valid", "Invalid", "Disabled"}
	case "Tree", "Collapsible", "Accordion":
		return []string{"Collapsed", "Expanded"}
	case "Rating":
		return []string{"Unrated", "Rated", "Disabled"}
	case "Radio", "RadioGroup", "Select", "Combobox", "Segmented", "Tabs", "Sidebar":
		return []string{"Default", "Selected", "Disabled"}
	case "Button", "PrimaryButton", "Link", "MenuButton", "ContextMenu", "Toolbar",
		"Slider", "StepSlider", "RangeSlider", "Stepper", "NumberInput", "EditableText",
		"Calendar", "DateInput", "TimeInput", "ColorWell", "ColorPicker", "Form", "Fieldset":
		return []string{"Default", "Disabled"}
	default:
		return []string{"Default"}
	}
}

func (g *Gallery) mygoState(name, variant string, first bool) *demoState {
	if g.mygoStates == nil {
		g.mygoStates = make(map[string]*demoState)
	}
	key := name + "/" + variant
	if first {
		key = name
	}
	if s := g.mygoStates[key]; s != nil {
		return s
	}
	s := newDemoState()
	switch variant {
	case "Empty":
		s.text = ""
		s.tags = nil
	case "On", "Checked":
		s.on, s.other = true, true
	case "Mixed":
		s.on = true
	case "Selected":
		s.index = 1
		s.choice = "Large"
	case "Rated":
		s.index = 4
	case "Expanded":
		s.open, s.other = true, true
	case "Valid":
		s.text = "ada@example.com"
	case "Invalid":
		s.text = "ada"
	}
	if name == "Sidebar" {
		s.choice = "Home"
		if variant == "Selected" {
			s.choice = "Files"
		}
	}
	g.mygoStates[key] = s
	return s
}

func (g *Gallery) MyGo(c *ui.Context, example mygoExample) {
	t := c.Theme()
	variants := mygoVariants(example.name)
	ui.Scroll(c).Key("mygo-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			for i, variant := range variants {
				s := g.mygoState(example.name, variant, i == 0)
				ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
					if len(variants) > 1 {
						ui.Text(c, variant).Bold().TextColor(t.TextMuted)
					}
					ui.Column(c).FillWidth().MinWidth(0).Gap(12).Disabled(variant == "Disabled").Children(func() {
						if variant == "Read only" {
							if example.name == "TextInput" {
								ui.TextInput(c, &s.text).Label("Name").ReadOnly(true).FillWidth()
							} else {
								ui.TextArea(c, &s.text).Label("Notes").ReadOnly(true).Height(140).FillWidth()
							}
						} else {
							example.build(c, s)
						}
					})
				})
			}
		})
	})
}
