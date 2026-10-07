package gallery

import (
	"github.com/egoist/mygo/ui"
	oneoui "github.com/oneo-me/oneo-go/ui"
)

type controlExample struct {
	name, label string
	variants    []string
	build       func(*ui.Context, *controlState, string)
}

// Control examples keep their state separate from the upstream widget catalog.
type controlState struct {
	on            bool
	choice        string
	text          string
	index, clicks int
	value         float64
}

var controlCatalog = []controlExample{
	{"Button", "Buttons", []string{"Default"}, controlButton},
	{"Checkbox", "Check boxes", []string{"Off", "On"}, controlCheckbox},
	{"Switch", "Switches", []string{"Off", "On"}, controlSwitch},
	{"Radio", "Radio buttons", []string{"Default", "Selected"}, controlRadio},
	{"TextField", "Text fields", []string{"Empty", "Filled", "Read only"}, controlTextField},
	{"TextArea", "Text areas", []string{"Empty", "Filled", "Read only"}, controlTextArea},
	{"Select", "Selects", []string{"Default", "Selected"}, controlSelect},
	{"Slider", "Sliders", []string{"Default"}, controlSlider},
	{"Tabs", "Tab rows", []string{"Default", "Selected"}, controlTabs},
}

func controlPageStart() galleryPage { return mygoPageStart + galleryPage(len(mygoCatalog)) }

func controlExampleFor(page galleryPage) (controlExample, bool) {
	i := int(page - controlPageStart())
	if i < 0 || i >= len(controlCatalog) {
		return controlExample{}, false
	}
	return controlCatalog[i], true
}

func (g *Gallery) controlState(name, variant string) *controlState {
	if g.controlStates == nil {
		g.controlStates = make(map[string]*controlState)
	}
	key := name + "/" + variant
	if s := g.controlStates[key]; s != nil {
		return s
	}
	s := &controlState{on: variant == "On", choice: "Small", value: 40}
	if variant == "Selected" {
		s.choice, s.index = "Large", 1
	}
	if variant != "Empty" {
		s.text = "Hello oneo-go"
	}
	if name == "TextArea" && variant != "Empty" {
		s.text += "\nA multiline note."
	}
	g.controlStates[key] = s
	return s
}

func (g *Gallery) Controls(c *ui.Context, example controlExample) {
	t := c.Theme()
	ui.Scroll(c).Key("controls-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			for _, variant := range example.variants {
				s := g.controlState(example.name, variant)
				ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
					ui.Text(c, variant).Bold().TextColor(t.TextMuted)
					ui.Column(c).FillWidth().MinWidth(0).Gap(12).Children(func() {
						example.build(c, s, variant)
					})
				})
			}
		})
	})
}

func controlButton(c *ui.Context, s *controlState, _ string) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		for _, variant := range []struct {
			name  string
			value oneoui.ButtonVariant
		}{
			{"Primary", oneoui.ButtonPrimary}, {"Default", oneoui.ButtonDefault},
			{"Outline", oneoui.ButtonOutline}, {"Ghost", oneoui.ButtonGhost}, {"Link", oneoui.ButtonLink},
		} {
			if oneoui.Button(c, variant.name, variant.value).Clicked() {
				s.clicks++
			}
		}
	})
	ui.Textf(c, "Clicks: %d", s.clicks)
}

func controlCheckbox(c *ui.Context, s *controlState, _ string) {
	oneoui.Checkbox(c, &s.on, "Notifications")
	ui.Textf(c, "Enabled: %t", s.on)
}

func controlSwitch(c *ui.Context, s *controlState, _ string) {
	oneoui.Switch(c, &s.on, "Sound")
	ui.Textf(c, "Sound: %t", s.on)
}

func controlRadio(c *ui.Context, s *controlState, _ string) {
	ui.RadioGroup(c, func() {
		for _, option := range []string{"Small", "Medium", "Large"} {
			oneoui.Radio(c, &s.choice, option, option)
		}
	}).Label("Size")
	ui.Text(c, "Selected: "+s.choice)
}

func controlTextField(c *ui.Context, s *controlState, variant string) {
	oneoui.TextField(c, &s.text).Label("Name").Placeholder("Your name").ReadOnly(variant == "Read only").FillWidth().MinWidth(0)
	ui.Text(c, s.text).Wrap()
}

func controlTextArea(c *ui.Context, s *controlState, variant string) {
	oneoui.TextArea(c, &s.text).Label("Notes").Placeholder("Write a note").ReadOnly(variant == "Read only").Height(100).FillWidth().MinWidth(0)
	ui.Textf(c, "%d characters", len([]rune(s.text)))
}

func controlSelect(c *ui.Context, s *controlState, _ string) {
	t := c.Theme()
	selectParts := oneoui.Select(c, &s.choice)
	selectParts.Trigger.Label("Size").Children(func() { ui.Text(c, s.choice).SingleLine() })
	selectParts.Popup(func(panel *ui.Element) {
		panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		for _, option := range []string{"Small", "Medium", "Large"} {
			item := selectParts.Item(option).Label(option).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
			if item.Highlighted() {
				item.Background(t.Accent).TextColor(t.AccentText)
			}
			item.Children(func() { ui.Text(c, option).SingleLine() })
		}
	})
	ui.Text(c, "Selected: "+s.choice)
}

func controlSlider(c *ui.Context, s *controlState, _ string) {
	oneoui.Slider(c, &s.value, 0, 100).Label("Volume").FillWidth().MinWidth(0)
	ui.Textf(c, "Volume: %.0f", s.value)
}

func controlTabs(c *ui.Context, s *controlState, _ string) {
	oneoui.Tabs(c, &s.index, "General", "Style")
	ui.Text(c, []string{"General settings", "Style settings"}[s.index])
}
