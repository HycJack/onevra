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
	on      bool
	choice  string
	text    string
	index   int
	button  buttonExampleState
	toggles [2]bool
	value   float64
}

var controlCatalog = []controlExample{
	{"Text", "Display text", []string{"Default"}, controlText},
	{"Button", "Buttons", []string{"Button", "MenuButton"}, controlButton},
	{"Selection", "Selection", []string{"Checkbox", "Switch", "Radio"}, controlSelection},
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
	s := &controlState{on: variant == "On", choice: "Small", value: 40, toggles: [2]bool{false, true}}
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
			if example.name == "Button" || example.name == "Selection" {
				for _, variant := range example.variants {
					g.controlSection(c, variant, func() {
						example.build(c, g.controlState(example.name, variant), variant)
					})
				}
				return
			}
			g.controlSection(c, example.name, func() {
				ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(4)).Children(func() {
					for _, variant := range example.variants {
						cell := ui.Column(c).Key(variant).MinWidth(0).MaxWidthPercent(100).Gap(12)
						switch example.name {
						case "Text":
							cell.FillWidth()
						case "TextField", "TextArea", "Slider":
							cell.Width(240)
						}
						cell.Children(func() {
							if len(example.variants) > 1 {
								components.Text(c, variant).TextColor(t.TextMuted)
							}
							example.build(c, g.controlState(example.name, variant), variant)
						})
					}
				})
			})
		})
	})
}

func (g *Gallery) controlSection(c *ui.Context, name string, build func()) {
	ui.Column(c).Key(name).FillWidth().MinWidth(0).Gap(12).Children(func() {
		components.Text(c, name).Bold().TextColor(c.Theme().TextMuted)
		controlPreview(c).Children(build)
	})
}

// controlPreview provides a shared canvas for component demonstrations.
func controlPreview(c *ui.Context) *ui.Element {
	t := c.Theme()
	return components.PixelGrid(c, components.PixelGridOptions{Spacing: 22}).
		Key("control-preview").Column().FillWidth().MinWidth(0).
		Padding(20).Gap(12).Radius(t.Radius).Border(1, t.Border)
}
