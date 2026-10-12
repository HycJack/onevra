package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

const mygoPageStart galleryPage = notificationsPage + 1

type mygoExample struct {
	name, group string
	build       func(*ui.Context, *demoState)
}

var mygoGroups = []string{"Content display", "Feedback & overlays", "Layout & containers"}

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
	case "Motion":
		return []string{"Springs", "Keyframes", "Animated pose", "Draggable"}
	case "Transforms":
		return []string{"Editable controls", "Perspective"}
	case "Toggle", "ToggleGroup", "Checkbox", "Switch":
		return []string{"Off", "On", "Disabled"}
	case "CheckboxGroup":
		return []string{"Unchecked", "Mixed", "Checked", "Disabled"}
	case "TextInput", "TextArea", "TextInputBuffer", "TextAreaBuffer":
		return []string{"Empty", "Filled", "Read only", "Disabled"}
	case "SearchField", "TokenField":
		return []string{"Empty", "Filled", "Disabled"}
	case "Tree", "Collapsible", "Accordion":
		return []string{"Collapsed", "Expanded"}
	case "Rating":
		return []string{"Unrated", "Rated", "Disabled"}
	case "Radio", "RadioGroup", "Select", "Segmented", "Tabs", "Sidebar":
		return []string{"Default", "Selected", "Disabled"}
	case "Button", "PrimaryButton", "Link",
		"Slider", "StepSlider", "RangeSlider", "Stepper", "NumberInput", "EditableText",
		"Calendar", "DateInput", "TimeInput", "ColorWell", "ColorPicker":
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
	s.variant = variant
	if name == "Transforms" {
		s.value = 8
	}
	if name == "Audio" {
		s.audio.output, s.audio.volume = g.audioOutput, .8
	}
	s.flyout.gallery = g
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
	}
	if name == "Sidebar" {
		s.choice = "Home"
		if variant == "Selected" {
			s.choice = "Files"
		}
	}
	if name == "TextInputBuffer" || name == "TextAreaBuffer" {
		text := s.text
		if name == "TextAreaBuffer" && variant != "Empty" {
			text = "Hello MyGo\nEdit this document with indexed text storage.\n你好，世界 🌍"
		}
		s.buffer = ui.NewTextBuffer(text)
		s.readOnly = variant == "Read only"
	}
	g.mygoStates[key] = s
	return s
}

func (g *Gallery) MyGo(c *ui.Context, example mygoExample) {
	g.withLanguage(c, func() { g.myGoLocalized(c, example) })
}

func (g *Gallery) myGoLocalized(c *ui.Context, example mygoExample) {
	if example.name == "Terminal" {
		_, exists := g.mygoStates["Terminal"]
		s := g.mygoState("Terminal", "Default", true)
		if !exists && g.createTerminal != nil {
			s.terminal.start(g.createTerminal)
		}
		ui.Column(c).Key("mygo-Terminal").Fill().Children(func() {
			example.build(c, s)
		})
		return
	}
	if example.name == "WebView" {
		d := &g.mygoState("WebView", "Default", true).webView
		d.covered = g.settings.open || g.commands.open
		if g.createWebView != nil {
			d.start(g.createWebView)
		}
	}
	t := c.Theme()
	variants := mygoVariants(example.name)
	components.ScrollView(c, components.ScrollViewOptions{}).Key("mygo-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			for i, variant := range variants {
				s := g.mygoState(example.name, variant, i == 0)
				ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
					if len(variants) > 1 {
						components.Text(c, tr(c, variant)).Bold().TextColor(t.TextMuted)
					}
					ui.Column(c).FillWidth().MinWidth(0).Gap(12).Disabled(variant == "Disabled").Children(func() {
						if variant == "Read only" && s.buffer == nil {
							if example.name == "TextInput" {
								components.TextBox(c, &s.text, components.TextBoxOptions{ReadOnly: true}).Label(tr(c, "gallery.name")).FillWidth()
							} else {
								components.TextBox(c, &s.text, components.TextBoxOptions{Multiline: true, ReadOnly: true}).Label(tr(c, "gallery.notes")).Height(140).FillWidth()
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
