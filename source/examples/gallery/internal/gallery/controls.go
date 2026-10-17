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
	shader                 *shaderDemoState
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
	tags                   []string
	index                  int
	button                 buttonExampleState
	toggles                [2]bool
	folders                [3]bool
	value                  float64
	high                   float64
	honeycomb              string
	honeycombItems         []honeycombDemoItem
	honeycombNext          int
	split, splitPercent    float32
	color                  ui.Color
	day                    time.Time
	badge                  struct{ status, unread int }
}

// controlSections lists the components a page gathers. Each section keeps its
// own variants and state under the component's name.
var controlSections = map[string][]controlExample{
	"Date & time pickers": {
		{"TimePicker", "TimePicker", "Date & time", []string{"Default"}, controlTimePicker},
		{"Calendar", "Calendar", "Date & time", []string{"Default"}, controlCalendar},
		{"CalendarPicker", "CalendarPicker", "Date & time", []string{"Default"}, controlCalendarPicker},
	},
	"Color pickers": {
		{"ColorView", "ColorView", "Color", []string{"Default"}, controlColorView},
		{"ColorPicker", "ColorPicker", "Color", []string{"Default", "Outlined", "Rounded"}, controlColorPicker},
	},
	"Dialogs": {
		{"Dialog", "Dialog", "Feedback & overlays", []string{"Default", "Protected"}, controlDialog},
		{"AlertDialog", "AlertDialog", "Feedback & overlays", []string{"Confirmation", "Notice"}, controlAlertDialog},
	},
	"TextBox": {
		{"TextBox", "TextBox", "Text input", []string{"Empty", "Filled", "Read only", "Search"}, controlTextBox},
		{"TextBox (multiline)", "TextBox (multiline)", "Text input", []string{"Empty", "Filled", "Read only"}, controlTextBoxMultiline},
		{"TextBoxBuffer", "TextBoxBuffer", "Text input", []string{"Empty", "Filled", "Read only", "Search"}, controlTextBoxBuffer},
		{"TextBoxBuffer (multiline)", "TextBoxBuffer (multiline)", "Text input", []string{"Empty", "Filled", "Read only"}, controlTextBoxBufferMultiline},
	},
	"Progress": {
		{"Progress", "Progress", "Feedback & overlays", []string{"Progress"}, controlProgress},
		{"ProgressRing", "ProgressRing", "Feedback & overlays", []string{"ProgressRing"}, controlProgress},
		{"Spinner", "Spinner", "Feedback & overlays", []string{"Ring", "Spokes", "Tinted"}, controlSpinner},
	},
	"Disclosure": {
		{"Accordion", "Accordion", "Layout & containers", []string{"Default"}, controlAccordion},
		{"Collapsible", "Collapsible", "Layout & containers", []string{"Right arrow", "Left arrow"}, controlCollapsible},
	},
	"Split panes": {
		{"Split", "Split", "Layout & containers", []string{"Horizontal", "Vertical"}, controlSplit},
		{"SplitPercent", "SplitPercent", "Layout & containers", []string{"Horizontal", "Vertical"}, controlSplitPercent},
	},
}

var controlCatalog = []controlExample{
	{"Shader backgrounds", "Shader backgrounds", "Shaders", []string{"Default"}, controlShaderBackgrounds},
	{"Data charts", "Data charts", "Charts", []string{"Bar chart", "Line chart", "Area chart", "Stacked bar chart", "Donut chart"}, controlChart},
	{"Relationship graph", "Relationship graph", "Charts", []string{"Default"}, controlGraph},
	{"Date & time pickers", "Date & time pickers", "Date & time", nil, nil},
	{"Color pickers", "Color pickers", "Color", nil, nil},
	{"Tooltip", "Tooltip", "Feedback & overlays", []string{"Default"}, controlTooltip},
	{"Popover", "Popover", "Feedback & overlays", []string{"Default"}, controlPopover},
	{"Dialogs", "Dialogs", "Feedback & overlays", nil, nil},
	{"Badge", "Badge", "Feedback & overlays", []string{"Tones", "Icons", "Live updates"}, controlBadge},
	{"Toast", "Toast", "Feedback & overlays", []string{"Status", "Actions", "Progress"}, controlToast},
	{"Chat", "Chat", "Conversation", []string{"Default"}, controlChat},
	{"RichText", "RichText", "Content display", []string{"Default"}, controlRichText},
	{"Breadcrumbs", "Breadcrumbs", "Actions & navigation", []string{"Default"}, controlBreadcrumbs},
	{"Rating", "Rating", "Selection & values", []string{"Default"}, controlRating},
	{"Text", "Text", "Content display", []string{"Default", "Editable"}, controlText},
	{"TextBox", "TextBox", "Text input", nil, nil},
	{"TagInput", "TagInput", "Text input", []string{"Default", "Empty"}, controlTagInput},
	{"Button", "Buttons", "Actions & navigation", []string{"Button", "Icon and text", "IconButton", "Toolbar", "ToggleButton"}, controlButton},
	{"HoldToConfirm", "HoldToConfirm", "Actions & navigation", []string{"Accent", "Danger", "Neutral"}, controlHoldToConfirm},
	{"Selection", "Selection", "Selection & values", []string{"Checkbox", "Switch", "Radio"}, controlSelection},
	{"Select", "Selects", "Selection & values", []string{"Default", "Selected", "Searchable", "NativeSelect"}, controlSelect},
	{"Menu", "Menus", "Feedback & overlays", []string{"Left-click", "Right-click"}, controlMenu},
	{"NumberInput", "NumberInput", "Selection & values", []string{"Integer", "Decimal", "Scrub", "Without stepper"}, controlNumberInput},
	{"Slider", "Sliders", "Selection & values", []string{"Default", "Step", "Range", "Stepped range"}, controlSlider},
	{"Progress", "Progress", "Feedback & overlays", nil, nil},
	{"Skeleton", "Skeleton", "Feedback & overlays", []string{"Default", "Avatar", "Content"}, controlSkeleton},
	{"Tabs", "Tabs", "Actions & navigation", []string{"Underline", "Segmented control", "Highlight"}, controlTabs},
	{"Disclosure", "Disclosure", "Layout & containers", nil, nil},
	{"Image", "Images", "Content display", []string{"Contain", "Cover", "Stretch", "Rounded", "Grayscale"}, controlImage},
	{"List", "List", "Content display", []string{"Default"}, controlList},
	{"GridView", "GridView", "Content display", []string{"Default"}, controlGridView},
	{"Tree", "Tree", "Content display", []string{"Tree", "TreeVirtualized"}, controlTree},
	{"Table", "Table", "Content display", []string{"Table", "TableVirtualized"}, controlTable},
	{"Sidebar", "Sidebar", "Actions & navigation", []string{"Default", "File tree", "Collapsible groups", "Custom leading"}, controlSidebar},
	{"Honeycomb", "Honeycomb", "Selection & values", []string{"Default"}, controlHoneycomb},
	{"Split panes", "Split panes", "Layout & containers", nil, nil},
	{"ScrollView", "Scroll views", "Layout & containers", []string{"Vertical", "Horizontal", "Both", "Insets"}, controlScrollView},
}

// controlToolbars hold page-level sample actions, which the titlebar shows
// after the page title so the content stays with the example itself.
var controlToolbars = map[string]func(*ui.Context, *controlState){
	"Data charts": chartToolbar, "Relationship graph": graphToolbar, "Honeycomb": honeycombToolbar, "Chat": chatToolbar,
}

func controlPageStart() galleryPage { return mygoPageStart + galleryPage(len(mygoCatalog)) }

func controlExampleFor(page galleryPage) (controlExample, bool) {
	i := int(page - controlPageStart())
	if i < 0 || i >= len(controlCatalog) {
		return controlExample{}, false
	}
	return controlCatalog[i], true
}

// pageComponents names the components a gathered page shows, so command
// search finds the page by any of them.
func pageComponents(example controlExample) []string {
	if example.name == "Data charts" {
		return example.variants
	}
	var names []string
	for _, section := range controlSections[example.name] {
		names = append(names, section.name)
	}
	return names
}

func (g *Gallery) controlState(name, variant string) *controlState {
	if g.controlStates == nil {
		g.controlStates = make(map[string]*controlState)
	}
	key := name + "/" + variant
	if s := g.controlStates[key]; s != nil {
		return s
	}
	s := &controlState{day: time.Now(), color: ui.RGB(110, 90, 230), on: variant == "On", choice: "Small", value: 40, high: 80, split: 100, splitPercent: 30, toggles: [2]bool{false, true}}
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
	if name == "TagInput" && variant != "Empty" {
		s.tags = []string{"design", "motion"}
	}
	if name == "Rating" && variant != "Empty" {
		s.index = 3
	}
	if name == "Breadcrumbs" {
		s.text = ""
	}
	if name == "Sidebar" {
		s.choice = map[string]string{"File tree": "main.go", "Collapsible groups": "home", "Custom leading": "personal2"}[variant]
		s.folders = [3]bool{true, variant == "Collapsible groups", false}
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
	if example.group == "Shaders" {
		// Backgrounds reach the page's edges, behind their sample content.
		ui.Box(c).Key("controls-" + example.name).Fill().MinWidth(0).MinHeight(0).Children(func() {
			example.build(c, g.controlState(example.name, "Default"), "Default")
		})
		return
	}
	if example.name == "Relationship graph" || example.name == "Honeycomb" {
		// These views fill the page instead of scrolling.
		ui.Column(c).Key("controls-" + example.name).Fill().MinWidth(0).MinHeight(0).Padding(t.Space(6)).Gap(12).Children(func() {
			example.build(c, g.controlState(example.name, "Default"), "Default")
		})
		return
	}
	components.ScrollView(c, components.ScrollViewOptions{}).Key("controls-" + example.name).Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			if sections := controlSections[example.name]; sections != nil {
				for _, section := range sections {
					g.controlVariants(c, section, true)
				}
				return
			}
			if example.name == "Data charts" {
				// Every chart follows the titlebar's data actions.
				s := g.controlState(example.name, "Default")
				for _, variant := range example.variants {
					g.controlSection(c, variant, func() { example.build(c, s, variant) })
				}
				return
			}
			if example.name == "Select" {
				selectExample := example
				selectExample.variants = example.variants[:len(example.variants)-1]
				g.controlVariants(c, selectExample, false)
				g.controlSection(c, "NativeSelect", func() {
					controlSelect(c, g.controlState("Select", "NativeSelect"), "NativeSelect")
				})
				return
			}
			if example.name == "Button" || example.name == "Selection" || example.name == "Tree" || example.name == "Table" {
				for _, variant := range example.variants {
					g.controlSection(c, variant, func() {
						example.build(c, g.controlState(example.name, variant), variant)
					})
				}
				return
			}
			g.controlVariants(c, example, false)
			if example.name == "Toast" {
				g.controlSection(c, "Placement", func() { g.toastPlacement(c) })
			}
		})
	})
}

// controlVariants lays out an example's variants. The titlebar already names
// the page, so only a section within a page gets its own heading.
func (g *Gallery) controlVariants(c *ui.Context, example controlExample, heading bool) {
	t := c.Theme()
	variants := func() {
		ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Start).Gap(t.Space(4)).Children(func() {
			for _, variant := range example.variants {
				cell := ui.Column(c).Key(variant).MinWidth(0).MaxWidthPercent(100).Gap(12)
				switch example.name {
				case "Text", "GridView":
					cell.FillWidth()
				case "TagInput":
					cell.Width(320)
				case "TextBox", "TextBox (multiline)", "TextBoxBuffer", "TextBoxBuffer (multiline)", "NumberInput", "Slider", "Sidebar", "List", "Split", "SplitPercent", "ScrollView", "Progress", "ProgressRing", "Skeleton", "Accordion", "Collapsible":
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
	}
	if example.name == "Text" || example.name == "Collapsible" || example.name == "Tabs" {
		variants = func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(4)).Children(func() {
				for _, variant := range example.variants {
					ui.Column(c).Key(variant).FillWidth().MinWidth(0).Gap(12).Children(func() {
						components.Text(c, tr(c, variant)).TextColor(t.TextMuted)
						example.build(c, g.controlState(example.name, variant), variant)
					})
				}
			})
		}
	}
	if !heading {
		ui.Column(c).Key(example.name).FillWidth().MinWidth(0).Children(variants)
		return
	}
	g.controlSection(c, example.name, variants)
}

func (g *Gallery) controlSection(c *ui.Context, name string, build func()) {
	ui.Column(c).Key(name).FillWidth().MinWidth(0).Gap(12).Children(func() {
		components.Text(c, tr(c, name)).Bold().TextColor(c.Theme().Text)
		build()
	})
}
