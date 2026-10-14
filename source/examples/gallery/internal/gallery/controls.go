package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

type controlExample struct {
	name, label, group string
	variants           []string
	build              func(*ui.Context, *controlState, string)
}

// Control examples keep their state separate from the upstream widget catalog.
type controlState struct {
	chat                   *chatDemoState
	chart                  *chartDemoState
	graph                  *graphDemoState
	dialogOpen, nestedOpen bool
	tree                   components.TreeVirtualizedState[string]
	table                  components.TableVirtualizedState
	tableSort              ui.SortOrder
	tableRows              []int
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
	{"Bar chart", "Bar chart", "Charts", []string{"Default"}, controlBarChart},
	{"Line chart", "Line chart", "Charts", []string{"Default"}, controlLineChart},
	{"Area chart", "Area chart", "Charts", []string{"Default"}, controlAreaChart},
	{"Stacked bar chart", "Stacked bar chart", "Charts", []string{"Default"}, controlStackedBarChart},
	{"Donut chart", "Donut chart", "Charts", []string{"Default"}, controlDonutChart},
	{"Relationship graph", "Relationship graph", "Charts", []string{"Default"}, controlGraph},
	{"TimePicker", "TimePicker", "Date & time", []string{"Default"}, controlTimePicker},
	{"Calendar", "Calendar", "Date & time", []string{"Default"}, controlCalendar},
	{"CalendarPicker", "CalendarPicker", "Date & time", []string{"Default"}, controlCalendarPicker},
	{"ColorView", "ColorView", "Color", []string{"Default"}, controlColorView},
	{"ColorPicker", "ColorPicker", "Color", []string{"Default", "Outlined", "Rounded"}, controlColorPicker},
	{"Tooltip", "Tooltip", "Feedback & overlays", []string{"Default"}, controlTooltip},
	{"Popover", "Popover", "Feedback & overlays", []string{"Default"}, controlPopover},
	{"Dialog", "Dialog", "Feedback & overlays", []string{"Default", "Protected"}, controlDialog},
	{"AlertDialog", "AlertDialog", "Feedback & overlays", []string{"Confirmation", "Notice"}, controlAlertDialog},
	{"Badge", "Badge", "Feedback & overlays", []string{"Default", "Accent"}, controlBadge},
	{"Toast", "Toast", "Feedback & overlays", []string{"Status", "Actions", "Progress"}, controlToast},
	{"Code", "Code", "Code & diffs", []string{"Default"}, controlCode},
	{"Markdown", "Markdown", "Content display", []string{"Default"}, controlMarkdown},
	{"Chat", "Chat", "Content display", []string{"Default"}, controlChat},
	{"RichText", "RichText", "Content display", []string{"Default"}, controlRichText},
	{"Breadcrumbs", "Breadcrumbs", "Actions & navigation", []string{"Default"}, controlBreadcrumbs},
	{"Rating", "Rating", "Selection & values", []string{"Default"}, controlRating},
	{"Text", "Text", "Content display", []string{"Default", "Editable"}, controlText},
	{"TextBox", "TextBox", "Text input", []string{"TextBox", "TextBox (multiline)"}, controlTextBox},
	{"TextBoxBuffer", "TextBoxBuffer", "Text input", []string{"TextBoxBuffer", "TextBoxBuffer (multiline)"}, controlTextBoxBuffer},
	{"Button", "Buttons", "Actions & navigation", []string{"Button", "ToggleButton"}, controlButton},
	{"Selection", "Selection", "Selection & values", []string{"Checkbox", "Switch", "Radio"}, controlSelection},
	{"Select", "Selects", "Selection & values", []string{"Default", "Selected", "Searchable", "NativeSelect"}, controlSelect},
	{"Menu", "Menus", "Feedback & overlays", []string{"Left-click", "Right-click"}, controlMenu},
	{"NumberInput", "NumberInput", "Selection & values", []string{"Integer", "Decimal", "Scrub", "Without stepper"}, controlNumberInput},
	{"Slider", "Sliders", "Selection & values", []string{"Default", "Step", "Range", "Stepped range"}, controlSlider},
	{"Progress", "Progress", "Feedback & overlays", []string{"Progress", "ProgressRing"}, controlProgress},
	{"Spinner", "Spinner", "Feedback & overlays", []string{"Ring", "Spokes", "Tinted"}, controlSpinner},
	{"Tabs", "Tabs", "Actions & navigation", []string{"Default"}, controlTabs},
	{"Accordion", "Accordion", "Layout & containers", []string{"Default"}, controlAccordion},
	{"Collapsible", "Collapsible", "Layout & containers", []string{"Right arrow", "Left arrow"}, controlCollapsible},
	{"Icon", "Icons", "Content display", []string{"Default"}, controlIcon},
	{"Image", "Images", "Content display", []string{"Contain", "Cover", "Stretch", "Rounded", "Grayscale"}, controlImage},
	{"List", "List", "Content display", []string{"Default"}, controlList},
	{"GridView", "GridView", "Content display", []string{"Default"}, controlGridView},
	{"Tree", "Tree", "Content display", []string{"Tree", "TreeVirtualized"}, controlTree},
	{"Table", "Table", "Content display", []string{"Table", "TableVirtualized"}, controlTable},
	{"Sidebar", "Sidebar", "Actions & navigation", []string{"Default"}, controlSidebar},
	{"PixelGrid", "PixelGrid", "Layout & containers", []string{"Uniform", "Major lines"}, controlPixelGrid},
	{"Split", "Split panes", "Layout & containers", []string{"Horizontal", "Vertical"}, controlSplit},
	{"ScrollView", "Scroll views", "Layout & containers", []string{"Vertical", "Horizontal", "Both", "Insets"}, controlScrollView},
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
	g.withLanguage(c, func() { g.controlsLocalized(c, example) })
}

func (g *Gallery) controlsLocalized(c *ui.Context, example controlExample) {
	t := c.Theme()
	if example.name == "Chat" {
		example.build(c, g.controlState(example.name, "Default"), "Default")
		return
	}
	components.ScrollView(c, components.ScrollViewOptions{}).Key("controls-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			if example.name == "Markdown" {
				example.build(c, g.controlState(example.name, "Default"), "Default")
				return
			}
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
			if example.name == "Toast" {
				g.controlSection(c, "Placement", func() { g.toastPlacement(c) })
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
			if example.name == "Button" || example.name == "Selection" || example.name == "Progress" || example.name == "Tree" || example.name == "Table" {
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
			components.Text(c, tr(c, example.name)).Bold().TextColor(t.Text)
			for _, variant := range example.variants {
				ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
					components.Text(c, tr(c, variant)).TextColor(t.TextMuted)
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
				case "Text", "GridView", "Markdown", "Code", "Bar chart", "Line chart", "Area chart", "Stacked bar chart", "Donut chart", "Relationship graph":
					cell.FillWidth()
				case "TextBox", "TextBox (multiline)", "TextBoxBuffer", "TextBoxBuffer (multiline)", "NumberInput", "Slider", "Sidebar", "List", "Split", "ScrollView", "Progress", "Accordion", "Collapsible", "PixelGrid":
					cell.Width(240)
				}
				cell.Children(func() {
					if len(example.variants) > 1 {
						components.Text(c, tr(c, variant)).TextColor(t.TextMuted)
					}
					example.build(c, g.controlState(example.name, variant), variant)
				})
			}
		})
	})
}

func (g *Gallery) controlSection(c *ui.Context, name string, build func()) {
	ui.Column(c).Key(name).FillWidth().MinWidth(0).Gap(12).Children(func() {
		components.Text(c, tr(c, name)).Bold().TextColor(c.Theme().Text)
		build()
	})
}
