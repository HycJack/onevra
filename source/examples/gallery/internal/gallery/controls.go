package gallery

import (
	"time"

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
	dialogOpen, nestedOpen bool
	tree                   components.TreeVirtualizedState[string]
	list                   components.ListState
	grid                   components.GridState
	on                     bool
	choice                 string
	text                   string
	buffer                 *ui.TextBuffer
	query                  string
	index                  int
	button                 buttonExampleState
	toggles                [2]bool
	value                  float64
	high                   float64
	split                  float32
	color                  ui.Color
	day                    time.Time
}

var textBoxControlCatalog = []controlExample{
	{"TextBox", "Text boxes", "Text input", []string{"Empty", "Filled", "Read only", "Search"}, controlTextBox},
	{"TextBox (multiline)", "Text areas", "Text input", []string{"Empty", "Filled", "Read only"}, controlTextBoxMultiline},
}

var textBoxBufferControlCatalog = []controlExample{
	{"TextBoxBuffer", "Buffer text boxes", "Text input", []string{"Empty", "Filled", "Read only", "Search"}, controlTextBoxBuffer},
	{"TextBoxBuffer (multiline)", "Buffer text areas", "Text input", []string{"Empty", "Filled", "Read only"}, controlTextBoxBufferMultiline},
}

var controlCatalog = []controlExample{
	{"TimePicker", "TimePicker", "Selection & values", []string{"Default"}, controlTimePicker},
	{"Calendar", "Calendar", "Selection & values", []string{"Default"}, controlCalendar},
	{"CalendarPicker", "CalendarPicker", "Selection & values", []string{"Default"}, controlCalendarPicker},
	{"ColorView", "ColorView", "Selection & values", []string{"Default"}, controlColorView},
	{"ColorPicker", "ColorPicker", "Selection & values", []string{"Default", "Outlined", "Rounded"}, controlColorPicker},
	{"Tooltip", "Tooltip", "Overlays", []string{"Default"}, controlTooltip},
	{"Popover", "Popover", "Overlays", []string{"Default"}, controlPopover},
	{"Dialog", "Dialog", "Overlays", []string{"Default", "Protected"}, controlDialog},
	{"AlertDialog", "AlertDialog", "Overlays", []string{"Confirmation", "Notice"}, controlAlertDialog},
	{"Badge", "Badge", "Status", []string{"Default", "Accent"}, controlBadge},
	{"RichText", "RichText", "Media & text", []string{"Default"}, controlRichText},
	{"Breadcrumbs", "Breadcrumbs", "Navigation", []string{"Default"}, controlBreadcrumbs},
	{"Rating", "Rating", "Selection & values", []string{"Default"}, controlRating},
	{"Text", "Text", "Media & text", []string{"Default", "Editable"}, controlText},
	{"TextBox", "TextBox", "Text input", []string{"TextBox", "TextBox (multiline)"}, controlTextBox},
	{"TextBoxBuffer", "TextBoxBuffer", "Text input", []string{"TextBoxBuffer", "TextBoxBuffer (multiline)"}, controlTextBoxBuffer},
	{"Button", "Buttons", "Actions", []string{"Button", "ToggleButton"}, controlButton},
	{"Selection", "Selection", "Selection & values", []string{"Checkbox", "Switch", "Radio"}, controlSelection},
	{"Select", "Selects", "Selection & values", []string{"Default", "Selected", "Searchable", "NativeSelect"}, controlSelect},
	{"Menu", "Menus", "Overlays", []string{"Left-click", "Right-click"}, controlMenu},
	{"NumberInput", "NumberInput", "Text input", []string{"Integer", "Decimal"}, controlNumberInput},
	{"Slider", "Sliders", "Selection & values", []string{"Default", "Step", "Range", "Stepped range"}, controlSlider},
	{"Progress", "Progress", "Status", []string{"Progress", "ProgressRing"}, controlProgress},
	{"Spinner", "Spinner", "Status", []string{"Ring", "Spokes"}, controlSpinner},
	{"Tabs", "Tabs", "Navigation", []string{"Default"}, controlTabs},
	{"Accordion", "Accordion", "Containers", []string{"Default"}, controlAccordion},
	{"Collapsible", "Collapsible", "Containers", []string{"Right arrow", "Left arrow"}, controlCollapsible},
	{"Icon", "Icons", "Media & text", []string{"Default"}, controlIcon},
	{"Image", "Images", "Media & text", []string{"Contain", "Cover", "Stretch", "Rounded", "Grayscale"}, controlImage},
	{"List", "List", "Collections", []string{"Default"}, controlList},
	{"GridView", "GridView", "Collections", []string{"Default"}, controlGridView},
	{"Tree", "Tree", "Collections", []string{"Tree", "TreeVirtualized"}, controlTree},
	{"Sidebar", "Sidebar", "Navigation", []string{"Default"}, controlSidebar},
	{"PixelGrid", "PixelGrid", "Layout", []string{"Uniform", "Major lines"}, controlPixelGrid},
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
	s := &controlState{day: time.Now(), color: ui.RGB(110, 90, 230), on: variant == "On", choice: "Small", value: 40, high: 80, split: 100, toggles: [2]bool{false, true}}
	if variant == "Selected" {
		s.choice, s.index = "Large", 1
	}
	if variant != "Empty" && variant != "Search" {
		s.text = "Hello Onevra"
	}
	if (name == "TextBox (multiline)" || name == "TextBoxBuffer (multiline)") && variant != "Empty" {
		s.text += "\nA multiline note."
	}
	if name == "TextBoxBuffer" || name == "TextBoxBuffer (multiline)" {
		s.buffer = ui.NewTextBuffer(s.text)
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
			if example.name == "TextBox" || example.name == "TextBoxBuffer" {
				catalog := textBoxControlCatalog
				if example.name == "TextBoxBuffer" {
					catalog = textBoxBufferControlCatalog
				}
				for _, textExample := range catalog {
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
					example.build(c, g.controlState(example.name, variant), variant)
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
				case "Text", "GridView":
					cell.FillWidth()
				case "TextBox", "TextBox (multiline)", "TextBoxBuffer", "TextBoxBuffer (multiline)", "NumberInput", "Slider", "Sidebar", "List", "Split", "ScrollView", "Progress", "Accordion", "Collapsible", "PixelGrid":
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
		build()
	})
}
