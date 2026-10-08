package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
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
	{"Text", "Display text", []string{"Default"}, controlText},
	{"Button", "Buttons", []string{"Default"}, controlButton},
	{"Checkbox", "Check boxes", []string{"Off", "On"}, controlCheckbox},
	{"Switch", "Switches", []string{"Off", "On"}, controlSwitch},
	{"Radio", "Radio buttons", []string{"Default", "Selected"}, controlRadio},
	{"TextField", "Text fields", []string{"Empty", "Filled", "Read only"}, controlTextField},
	{"TextArea", "Text areas", []string{"Empty", "Filled", "Read only"}, controlTextArea},
	{"Select", "Selects", []string{"Default", "Selected"}, controlSelect},
	{"Slider", "Sliders", []string{"Default"}, controlSlider},
	{"Tabs", "Tab rows", []string{"Default", "Selected"}, controlTabs},
	{"PixelGrid", "Pixel grid", []string{"Default", "Offset"}, controlPixelGrid},
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
		s.text = "Hello Onevra"
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
					components.Text(c, variant).Bold().TextColor(t.TextMuted)
					ui.Column(c).FillWidth().MinWidth(0).Gap(12).Children(func() {
						example.build(c, s, variant)
					})
				})
			}
		})
	})
}

func controlText(c *ui.Context, _ *controlState, _ string) {
	components.Text(c, "A heading").FontSize(c.Theme().Rem(1.5)).Bold()
	components.Text(c, "Hello, world. 你好，世界。こんにちは 🌍")
	components.Text(c, "Secondary text").TextColor(c.Theme().TextMuted)
	components.Text(c, "Long text wraps within the available width, including in narrow windows.").FillWidth()
	components.Text(c, "A single-line label").SingleLine()
	components.Textf(c, "%s: %d", "Count", 42)
	ui.RichText(c).Children(func() {
		components.Text(c, "Inline ")
		components.Text(c, "bold text").Bold()
	})
}

func controlButton(c *ui.Context, s *controlState, _ string) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		for _, variant := range []struct {
			name  string
			value components.ButtonVariant
		}{
			{"Primary", components.ButtonPrimary}, {"Secondary", components.ButtonSecondary},
			{"Outline", components.ButtonOutline}, {"Ghost", components.ButtonGhost}, {"Link", components.ButtonLink},
		} {
			if components.Button(c, variant.name, variant.value).Clicked() {
				s.clicks++
			}
		}
	})
	components.Textf(c, "Clicks: %d", s.clicks)
}

func controlCheckbox(c *ui.Context, s *controlState, _ string) {
	components.Checkbox(c, &s.on, "Notifications")
	components.Textf(c, "Enabled: %t", s.on)
}

func controlSwitch(c *ui.Context, s *controlState, _ string) {
	components.Switch(c, &s.on, "Sound")
	components.Textf(c, "Sound: %t", s.on)
}

func controlRadio(c *ui.Context, s *controlState, _ string) {
	ui.RadioGroup(c, func() {
		for _, option := range []string{"Small", "Medium", "Large"} {
			components.Radio(c, &s.choice, option, option)
		}
	}).Label("Size")
	components.Text(c, "Selected: "+s.choice)
}

func controlTextField(c *ui.Context, s *controlState, variant string) {
	components.TextField(c, &s.text).Label("Name").Placeholder("Your name").ReadOnly(variant == "Read only").FillWidth().MinWidth(0)
	components.Text(c, s.text).Wrap()
}

func controlTextArea(c *ui.Context, s *controlState, variant string) {
	components.TextArea(c, &s.text).Label("Notes").Placeholder("Write a note").ReadOnly(variant == "Read only").Height(100).FillWidth().MinWidth(0)
	components.Textf(c, "%d characters", len([]rune(s.text)))
}

func controlSelect(c *ui.Context, s *controlState, _ string) {
	t := c.Theme()
	selectParts := components.Select(c, &s.choice)
	selectParts.Trigger.Label("Size").Children(func() { components.Text(c, s.choice).SingleLine() })
	selectParts.Popup(func(panel *ui.Element) {
		panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		for _, option := range []string{"Small", "Medium", "Large"} {
			item := selectParts.Item(option).Label(option).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
			if item.Highlighted() {
				item.Background(t.Accent).TextColor(t.AccentText)
			}
			item.Children(func() { components.Text(c, option).SingleLine() })
		}
	})
	components.Text(c, "Selected: "+s.choice)
}

func controlSlider(c *ui.Context, s *controlState, _ string) {
	components.Slider(c, &s.value, 0, 100).Label("Volume").FillWidth().MinWidth(0)
	components.Textf(c, "Volume: %.0f", s.value)
}

func controlPixelGrid(c *ui.Context, s *controlState, variant string) {
	components.Textf(c, "Grid spacing: %.0f", s.value)
	components.Slider(c, &s.value, 4, 64).Label("Grid spacing").FillWidth().MinWidth(0)
	options := components.PixelGridOptions{Spacing: float32(s.value)}
	if variant == "Offset" {
		options.OffsetX, options.OffsetY = -13, 9
	}
	components.PixelGrid(c, options).FillWidth().MinWidth(0).Height(220).Padding(24).Children(func() {
		components.Button(c, "Canvas content", components.ButtonSecondary)
	})
}

func controlTabs(c *ui.Context, s *controlState, _ string) {
	components.Tabs(c, &s.index, "General", "Style")
	components.Text(c, []string{"General settings", "Style settings"}[s.index])
}
