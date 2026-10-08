package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

type controlExample struct {
	name, label, group string
	variants           []string
	build              func(*ui.Context, *controlState, string)
}

// Control examples keep their state separate from the upstream widget catalog.
type controlState struct {
	tree    components.TreeVirtualizedState[string]
	list    components.ListState
	on      bool
	choice  string
	text    string
	query   string
	index   int
	button  buttonExampleState
	toggles [2]bool
	value   float64
	high    float64
	split   float32
}

var textBoxControlCatalog = []controlExample{
	{"TextBox", "Text boxes", "Text input", []string{"Empty", "Filled", "Read only", "Search"}, controlTextBox},
	{"TextBox (multiline)", "Text areas", "Text input", []string{"Empty", "Filled", "Read only"}, controlTextBoxMultiline},
}

var controlCatalog = []controlExample{
	{"Badge", "Badge", "Status", []string{"Default", "Accent"}, controlBadge},
	{"RichText", "RichText", "Media & text", []string{"Default"}, controlRichText},
	{"Breadcrumbs", "Breadcrumbs", "Navigation", []string{"Default"}, controlBreadcrumbs},
	{"Rating", "Rating", "Selection & values", []string{"Default"}, controlRating},
	{"Text", "Text", "Media & text", []string{"Default", "Editable"}, controlText},
	{"TextBox", "TextBox", "Text input", []string{"TextBox", "TextBox (multiline)"}, controlTextBox},
	{"Button", "Buttons", "Actions", []string{"Button", "ToggleButton"}, controlButton},
	{"Selection", "Selection", "Selection & values", []string{"Checkbox", "Switch", "Radio"}, controlSelection},
	{"Select", "Selects", "Selection & values", []string{"Default", "Selected", "Searchable", "NativeSelect"}, controlSelect},
	{"Menu", "Menus", "Overlays", []string{"Left-click", "Right-click"}, controlMenu},
	{"NumberInput", "NumberInput", "Text input", []string{"Integer", "Decimal"}, controlNumberInput},
	{"Slider", "Sliders", "Selection & values", []string{"Default", "Step", "Range", "Stepped range"}, controlSlider},
	{"Progress", "Progress", "Status", []string{"Progress", "ProgressRing"}, controlProgress},
	{"Tabs", "Tabs", "Navigation", []string{"Default"}, controlTabs},
	{"Accordion", "Accordion", "Containers", []string{"Default"}, controlAccordion},
	{"Collapsible", "Collapsible", "Containers", []string{"Right arrow", "Left arrow"}, controlCollapsible},
	{"Icon", "Icons", "Media & text", []string{"Default"}, controlIcon},
	{"Image", "Images", "Media & text", []string{"Contain", "Cover", "Stretch", "Rounded", "Grayscale"}, controlImage},
	{"List", "List", "Collections", []string{"Default"}, controlList},
	{"Tree", "Tree", "Collections", []string{"Tree", "TreeVirtualized"}, controlTree},
	{"Sidebar", "Sidebar", "Navigation", []string{"Default"}, controlSidebar},
	{"Split", "Split panes", "Layout", []string{"Horizontal", "Vertical"}, controlSplit},
	{"ScrollView", "Scroll views", "Layout", []string{"Vertical", "Horizontal", "Both", "Insets"}, controlScrollView},
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
	s := &controlState{on: variant == "On", choice: "Small", value: 40, high: 80, split: 100, toggles: [2]bool{false, true}}
	if variant == "Selected" {
		s.choice, s.index = "Large", 1
	}
	if variant != "Empty" && variant != "Search" {
		s.text = "Hello Onevra"
	}
	if name == "TextBox (multiline)" && variant != "Empty" {
		s.text += "\nA multiline note."
	}
	if name == "Rating" && variant != "Empty" {
		s.index = 3
	}
	if name == "Breadcrumbs" {
		s.text = "Choose an ancestor"
	}
	g.controlStates[key] = s
	return s
}

func (g *Gallery) Controls(c *ui.Context, example controlExample) {
	t := c.Theme()
	components.ScrollView(c, components.ScrollViewOptions{}).Key("controls-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			if example.name == "TextBox" {
				for _, textExample := range textBoxControlCatalog {
					g.controlVariants(c, textExample)
				}
				return
			}
			if example.name == "Select" {
				selectExample := example
				selectExample.variants = example.variants[:len(example.variants)-1]
				g.controlVariants(c, selectExample)
				g.controlSection(c, "NativeSelect", func() {
					controlSelect(c, g.controlState("Select", "NativeSelect"), "NativeSelect")
				})
				return
			}
			if example.name == "Button" || example.name == "Selection" || example.name == "Progress" || example.name == "Tree" {
				for _, variant := range example.variants {
					g.controlSection(c, variant, func() {
						example.build(c, g.controlState(example.name, variant), variant)
					})
				}
				return
			}
			g.controlVariants(c, example)
		})
	})
}

func (g *Gallery) controlVariants(c *ui.Context, example controlExample) {
	t := c.Theme()
	if example.name == "Text" || example.name == "Collapsible" {
		ui.Column(c).Key(example.name).FillWidth().MinWidth(0).Gap(t.Space(4)).Children(func() {
			components.Text(c, example.name).Bold().TextColor(t.Text)
			for _, variant := range example.variants {
				ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
					components.Text(c, variant).TextColor(t.TextMuted)
					controlPreview(c).Children(func() {
						example.build(c, g.controlState(example.name, variant), variant)
					})
				})
			}
		})
		return
	}
	g.controlSection(c, example.name, func() {
		ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Start).Gap(t.Space(4)).Children(func() {
			for _, variant := range example.variants {
				cell := ui.Column(c).Key(variant).MinWidth(0).MaxWidthPercent(100).Gap(12)
				switch example.name {
				case "Text":
					cell.FillWidth()
				case "TextBox", "TextBox (multiline)", "NumberInput", "Slider", "Sidebar", "List", "Split", "ScrollView", "Progress", "Accordion", "Collapsible":
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
}

func (g *Gallery) controlSection(c *ui.Context, name string, build func()) {
	ui.Column(c).Key(name).FillWidth().MinWidth(0).Gap(12).Children(func() {
		components.Text(c, name).Bold().TextColor(c.Theme().Text)
		controlPreview(c).Children(build)
	})
}

// controlPreview provides a shared canvas for component demonstrations.
func controlPreview(c *ui.Context) ui.Element {
	t := c.Theme()
	return components.PixelGrid(c, components.PixelGridOptions{Spacing: 22}).
		Key("control-preview").Column().FillWidth().MinWidth(0).
		Padding(20).Gap(12).Radius(t.Radius).Border(1, t.Border)
}
