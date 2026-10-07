package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

// demoState is owned by one page, survives navigation, and is never persisted.
type demoState struct {
	terminal        terminalDemo
	text            string
	choice          string
	on, other, open bool
	index, count    int
	value, high     float64
	size            float32
	tags            []string
	day             time.Time
	color           ui.Color
	list            ui.ListState
	grid            ui.GridState
	outline         ui.OutlineState[string]
	router          *ui.Router
}

func newDemoState() *demoState {
	return &demoState{text: "Hello MyGo", choice: "Small", value: 40, high: 80,
		size: 140, day: time.Date(2026, 6, 15, 9, 30, 0, 0, time.Local),
		color: ui.RGB(80, 120, 220), tags: []string{"Go"}, router: ui.NewRouter("/home")}
}

func exampleButton(c *ui.Context, s *demoState) {
	if ui.Button(c, "Count").Clicked() {
		s.count++
	}
	ui.Textf(c, "Clicks: %d", s.count)
	ui.Button(c, "").Label("Settings").Tooltip("Settings").Children(func() { ui.Icon(c, lucide.Settings()) })
}

func examplePrimaryButton(c *ui.Context, s *demoState) {
	if ui.PrimaryButton(c, "Apply").Clicked() {
		s.count++
	}
	ui.Textf(c, "Applied %d times", s.count)
}

func exampleLink(c *ui.Context, s *demoState) {
	if ui.Link(c, "Show details", "").Clicked() {
		s.on = !s.on
	}
	if s.on {
		ui.Text(c, "Details").Wrap()
	}
}

func exampleMenuButton(c *ui.Context, s *demoState) {
	ui.MenuButton(c, "Choose size", func(m *ui.Menu) {
		for _, name := range []string{"Small", "Medium", "Large"} {
			if m.Item(name).Checked(s.choice == name).Chosen() {
				s.choice = name
			}
		}
	})
	ui.Text(c, s.choice)
}

func exampleContextMenu(c *ui.Context, s *demoState) {
	ui.Button(c, "Right-click or Shift+F10").ContextMenu(func(m *ui.Menu) {
		if m.Item("Increment").Chosen() {
			s.count++
		}
		m.Separator()
		m.Item("Unavailable").Disabled(true)
	})
	ui.Textf(c, "Count: %d", s.count)
}

func exampleToggle(c *ui.Context, s *demoState) {
	ui.Toggle(c, &s.on, "Bold")
	preview := ui.Text(c, "Preview")
	if s.on {
		preview.Bold()
	}
}

func exampleToggleGroup(c *ui.Context, s *demoState) {
	ui.ToggleGroup(c, func() {
		ui.Toggle(c, &s.on, "Bold")
		ui.Toggle(c, &s.other, "Italic")
	}).Label("Formatting")
	preview := ui.Text(c, "Preview")
	if s.on {
		preview.Bold()
	}
	if s.other {
		preview.Italic()
	}
}

func exampleSegmented(c *ui.Context, s *demoState) {
	ui.Segmented(c, &s.index, "List", "Grid").Label("View mode")
	ui.Text(c, []string{"List view selected", "Grid view selected"}[s.index])
}

func exampleToolbar(c *ui.Context, s *demoState) {
	ui.Toolbar(c, func() {
		ui.Toggle(c, &s.on, "Bold")
		ui.Spacer(c)
		if ui.Button(c, "Apply").Clicked() {
			s.count++
		}
	}).Label("Formatting").FillWidth()
	ui.Textf(c, "Applied %d times", s.count)
}

func exampleCheckbox(c *ui.Context, s *demoState) {
	ui.Checkbox(c, &s.on, "Enable notifications")
	ui.Textf(c, "Enabled: %t", s.on)
}

func exampleCheckboxGroup(c *ui.Context, s *demoState) {
	ui.CheckboxGroup(c, "Notifications", func() {
		ui.Checkbox(c, &s.on, "Email")
		ui.Checkbox(c, &s.other, "Sound")
	})
}

func exampleSwitch(c *ui.Context, s *demoState) {
	ui.Switch(c, &s.on).Label("Sound")
	ui.Textf(c, "Sound: %t", s.on)
}

func exampleRadio(c *ui.Context, s *demoState) {
	ui.RadioGroup(c, func() {
		for _, name := range []string{"Small", "Medium", "Large"} {
			ui.Radio(c, &s.choice, name, name)
		}
	}).Label("Size")
	ui.Text(c, s.choice)
}

func exampleSelect(c *ui.Context, s *demoState) {
	ui.Select(c, &s.choice, []string{"Small", "Medium", "Large"}).Label("Size")
	ui.Text(c, s.choice)
}

func exampleSlider(c *ui.Context, s *demoState) {
	ui.Slider(c, &s.value, 0, 100).Label("Volume")
	ui.Textf(c, "Volume: %.0f", s.value)
}

func exampleStepSlider(c *ui.Context, s *demoState) {
	ui.StepSlider(c, &s.value, 0, 100, 20).Label("Quality")
	ui.Textf(c, "Quality: %.0f", s.value)
}

func exampleRangeSlider(c *ui.Context, s *demoState) {
	ui.RangeSlider(c, &s.value, &s.high, 0, 100, 5).Label("Range")
	ui.Textf(c, "%.0f – %.0f", s.value, s.high)
}

func exampleStepper(c *ui.Context, s *demoState) {
	ui.Stepper(c, &s.value, 0, 100, 1).Label("Copies")
	ui.Textf(c, "Copies: %.0f", s.value)
}

func exampleRating(c *ui.Context, s *demoState) {
	ui.Rating(c, &s.index, 5).Label("Rating")
	ui.Textf(c, "%d of 5 stars", s.index)
}

func exampleTextInput(c *ui.Context, s *demoState) {
	ui.TextInput(c, &s.text).Label("Name").Placeholder("Your name").FillWidth()
	ui.Text(c, s.text).Wrap()
}

func exampleTextArea(c *ui.Context, s *demoState) {
	ui.TextArea(c, &s.text).Label("Notes").Height(140).FillWidth()
	ui.Textf(c, "%d characters", len([]rune(s.text)))
}

func exampleNumberInput(c *ui.Context, s *demoState) {
	ui.NumberInput(c, &s.value, 0, 100, 1).Label("Quantity")
	ui.Textf(c, "Quantity: %.0f", s.value)
}

func exampleSearchField(c *ui.Context, s *demoState) {
	ui.SearchField(c, &s.text).Label("Search samples").FillWidth()
	for _, name := range []string{"Hello MyGo", "Hello Go", "Native UI"} {
		if strings.Contains(strings.ToLower(name), strings.ToLower(s.text)) {
			ui.Text(c, name)
		}
	}
}

func exampleCombobox(c *ui.Context, s *demoState) {
	ui.Combobox(c, &s.choice, []string{"Small", "Medium", "Large"}).Label("Size")
	ui.Text(c, s.choice)
}

func exampleAutocomplete(c *ui.Context, s *demoState) {
	ui.Autocomplete(c, &s.text, []string{"Hello MyGo", "Hello Go", "Hello world"}).Label("Greeting")
	ui.Text(c, s.text).Wrap()
}

func exampleTokenField(c *ui.Context, s *demoState) {
	ui.TokenField(c, &s.tags, []string{"Go", "UI", "Native", "Desktop"}).Label("Tags").FillWidth()
	ui.Text(c, strings.Join(s.tags, ", ")).Wrap()
}

func exampleEditableText(c *ui.Context, s *demoState) {
	ui.EditableText(c, &s.text).Label("Editable name")
}

func exampleFindBar(c *ui.Context, s *demoState) {
	if ui.Button(c, "Find").Clicked() {
		s.open = true
	}
	lines := []string{"Hello MyGo", "Hello Go", "Native UI"}
	var matches []string
	for _, line := range lines {
		if s.text != "" && strings.Contains(strings.ToLower(line), strings.ToLower(s.text)) {
			matches = append(matches, line)
		}
	}
	ui.FindBar(c, &s.open, &s.text, len(matches), &s.index)
	ui.Text(c, strings.Join(lines, "\n")).Wrap()
	if len(matches) > 0 {
		ui.Text(c, "Match: "+matches[min(max(s.index, 0), len(matches)-1)]).Wrap()
	}
}

func exampleCalendar(c *ui.Context, s *demoState) {
	ui.Calendar(c, &s.day).Label("Day")
	ui.Text(c, s.day.Format("2006-01-02"))
}

func exampleDateInput(c *ui.Context, s *demoState) {
	ui.DateInput(c, &s.day).Label("Date")
	ui.Text(c, s.day.Format("2006-01-02"))
}

func exampleTimeInput(c *ui.Context, s *demoState) {
	ui.TimeInput(c, &s.day).Label("Time")
	ui.Text(c, s.day.Format("15:04"))
}

func exampleColorWell(c *ui.Context, s *demoState) {
	ui.ColorWell(c, &s.color).Label("Tint")
	ui.Box(c).Size(80, 40).Background(s.color)
}

func exampleColorPicker(c *ui.Context, s *demoState) {
	ui.ColorPicker(c, &s.color).Label("Color")
	ui.Box(c).Size(80, 40).Background(s.color)
}

func exampleList(c *ui.Context, s *demoState) {
	s.list.Selected = &s.index
	s.list.Label = func(i int) string { return fmt.Sprintf("Item %d", i+1) }
	ui.List(c, &s.list, 1000, func(i int) { ui.Textf(c, "Item %d", i+1).Padding(8) }).Height(200).FillWidth().Label("Items")
	ui.Textf(c, "Selected item: %d", s.index+1)
}

func exampleTable(c *ui.Context, s *demoState) {
	s.list.Selected = &s.index
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 90}}
	ui.Table(c, &s.list, cols, 100, func(row, col int) {
		if col == 0 {
			ui.Textf(c, "File %d", row+1)
		} else {
			ui.Textf(c, "%d KB", row+1)
		}
	}).Height(200).FillWidth().Label("Files")
}

func exampleTree(c *ui.Context, s *demoState) {
	ui.Tree(c, func() {
		ui.TreeItem(c, "source", &s.open, func() {
			if ui.TreeItem(c, "main.go", nil, nil).Selected(s.on).Clicked() {
				s.on = true
			}
		})
	}).Label("Project")
}

func exampleOutline(c *ui.Context, s *demoState) {
	s.outline.List.Selected = &s.index
	ui.Outline(c, &s.outline, []string{"source", "docs"}, func(name string) []string {
		if name == "source" {
			return []string{"main.go", "go.mod"}
		}
		if name == "docs" {
			return []string{"README.md"}
		}
		return nil
	}, func(name string) { ui.Text(c, name) }).Height(180).FillWidth().Label("Project outline")
}

func exampleOutlineTable(c *ui.Context, s *demoState) {
	s.outline.List.Selected = &s.index
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Kind", Width: 90}}
	ui.OutlineTable(c, &s.outline, cols, []string{"source"}, func(name string) []string {
		if name == "source" {
			return []string{"main.go", "go.mod"}
		}
		return nil
	}, func(name string, col int) {
		if col == 0 {
			ui.Text(c, name)
		} else {
			ui.Text(c, "File")
		}
	}).Height(180).FillWidth().Label("Project table")
}

func exampleGridView(c *ui.Context, s *demoState) {
	s.grid.Selected = &s.index
	s.grid.Label = func(i int) string { return fmt.Sprintf("Item %d", i+1) }
	ui.GridView(c, &s.grid, 1000, 90, 70, func(i int) {
		ui.Textf(c, "Item %d", i+1).Padding(12)
	}).Height(220).FillWidth().Label("Items")
	ui.Textf(c, "Selected item: %d", s.index+1)
}

func exampleScroll(c *ui.Context, s *demoState) {
	ui.Scroll(c).Height(180).FillWidth().Children(func() {
		for i := range 30 {
			ui.Textf(c, "Line %d", i+1).Padding(4)
		}
	})
}

func exampleScrollHorizontal(c *ui.Context, s *demoState) {
	ui.ScrollHorizontal(c).Height(80).FillWidth().Children(func() {
		ui.Row(c).Gap(12).Children(func() {
			for i := range 20 {
				ui.Textf(c, "Card %d", i+1).Width(120).Padding(16).Background(c.Theme().Surface)
			}
		})
	})
}

func exampleScrollBoth(c *ui.Context, s *demoState) {
	ui.ScrollBoth(c).Height(180).FillWidth().Children(func() {
		ui.Column(c).Width(800).Children(func() {
			for i := range 30 {
				ui.Textf(c, "Row %d — Name                         Status                         Updated", i+1).Padding(8)
			}
		})
	})
}

func exampleGrid(c *ui.Context, s *demoState) {
	ui.Grid(c).Columns(2).Gap(8).FillWidth().Children(func() {
		for i := range 6 {
			ui.Textf(c, "Cell %d", i+1).Padding(12).Background(c.Theme().Surface)
		}
	})
}

func exampleSplit(c *ui.Context, s *demoState) {
	ui.Split(c, &s.size, func() { ui.Text(c, "Left pane").Wrap() }, func() { ui.Text(c, "Right pane").Wrap() }).Height(180).FillWidth()
}

func exampleSplitVertical(c *ui.Context, s *demoState) {
	ui.SplitVertical(c, &s.size, func() { ui.Text(c, "Top pane") }, func() { ui.Text(c, "Bottom pane") }).Height(240).FillWidth()
}

func exampleCollapsible(c *ui.Context, s *demoState) {
	ui.Collapsible(c, "Details", &s.open, func() { ui.Checkbox(c, &s.on, "Notifications") })
}

func exampleAccordion(c *ui.Context, s *demoState) {
	ui.Accordion(c, func() {
		ui.AccordionItem(c, "General", &s.open, func() { ui.Checkbox(c, &s.on, "Notifications") })
		ui.AccordionItem(c, "About", &s.other, func() { ui.Text(c, "Built with MyGo") })
	}).FillWidth()
}

func exampleForm(c *ui.Context, s *demoState) {
	ui.Form(c, func() {
		ui.Field(c, "Name", func() { ui.TextInput(c, &s.text) })
		ui.Field(c, "Sound", func() { ui.Switch(c, &s.on) })
	}).FillWidth()
}

func exampleField(c *ui.Context, s *demoState) {
	message := ""
	if s.text != "" && !strings.Contains(s.text, "@") {
		message = "Enter an email address."
	}
	ui.Field(c, "Email", func() { ui.TextInput(c, &s.text).FillWidth() }).Description("Email address").Error(message)
}

func exampleFieldset(c *ui.Context, s *demoState) {
	ui.Form(c, func() {
		ui.Fieldset(c, "Notifications", func() {
			ui.Field(c, "Email", func() { ui.Checkbox(c, &s.on, "Weekly summary") })
			ui.Field(c, "Sound", func() { ui.Switch(c, &s.other) })
		})
	})
}

func exampleTabs(c *ui.Context, s *demoState) {
	ui.Tabs(c, &s.index, "General", "Appearance")
	ui.Text(c, []string{"General settings", "Appearance settings"}[s.index])
}

func exampleSidebar(c *ui.Context, s *demoState) {
	ui.Sidebar(c, &s.choice, func() {
		ui.SidebarSection(c, "Locations", nil, func() {
			ui.SidebarItem(c, "Home", lucide.House(), "Home")
			ui.SidebarItem(c, "Files", lucide.Folder(), "Files")
		})
	}).Height(180).FillWidth()
	ui.Text(c, "Selected: "+s.choice)
}

func exampleBreadcrumbs(c *ui.Context, s *demoState) {
	path := []string{"Home", "Documents", "Project"}
	ui.Breadcrumbs(c, path, &s.index).Label("Path")
	ui.Text(c, "Chosen: "+path[s.index])
}

func exampleBackButton(c *ui.Context, s *demoState) {
	ui.Toolbar(c, func() {
		ui.BackButton(c, s.router)
		ui.ForwardButton(c, s.router)
	}).Label("History")
	if ui.Button(c, "Visit next page").Clicked() {
		s.count++
		s.router.Push(fmt.Sprintf("/page/%d", s.count))
	}
	s.router.View(c, func(r *ui.Route) { r.Title(r.Path()); ui.Text(c, r.Path()) }).Height(60).FillWidth()
}

func exampleForwardButton(c *ui.Context, s *demoState) {
	if ui.Button(c, "Prepare forward history").Clicked() {
		s.count++
		s.router.Push(fmt.Sprintf("/page/%d", s.count))
		s.router.Back()
	}
	ui.ForwardButton(c, s.router)
	s.router.View(c, func(r *ui.Route) {
		r.Title(r.Path())
		ui.Text(c, r.Path())
	}).Height(60).FillWidth()
}

func exampleRadioGroup(c *ui.Context, s *demoState) {
	ui.RadioGroup(c, func() {
		ui.Radio(c, &s.choice, "Small", "Small")
		ui.Radio(c, &s.choice, "Large", "Large")
	}).Row().Gap(16).Label("Size options")
	ui.Text(c, s.choice)
}

func exampleModal(c *ui.Context, s *demoState) {
	if ui.Button(c, "Open dialog").Clicked() {
		s.open = true
	}
	ui.Modal(c, &s.open, func() {
		ui.Text(c, "Edit name").Bold()
		ui.TextInput(c, &s.text).Label("Name")
		if ui.Button(c, "Done").Clicked() {
			s.open = false
		}
	})
}

func exampleAlertDialog(c *ui.Context, s *demoState) {
	if ui.Button(c, "Ask for confirmation").Clicked() {
		s.open = true
	}
	if ui.AlertDialog(c, &s.open, "Apply change?", "Apply these settings?", "Cancel", "Apply") == 1 {
		s.count++
	}
	ui.Textf(c, "Confirmed: %d", s.count)
}

func examplePopover(c *ui.Context, s *demoState) {
	anchor := ui.Button(c, "Show popover")
	if anchor.Clicked() {
		s.open = !s.open
	}
	ui.Popover(c, anchor, &s.open, func() {
		ui.Text(c, "Popover content")
		if ui.Button(c, "Close").Clicked() {
			s.open = false
		}
	})
}

func exampleTooltip(c *ui.Context, s *demoState) {
	ui.Button(c, "Hover or focus me").Tooltip("A helpful description of this action")
}

func exampleToast(c *ui.Context, s *demoState) {
	if ui.Button(c, "Show toast").Clicked() {
		c.Toast("Saved")
	}
	if ui.Button(c, "Show undo toast").Clicked() {
		s.count++
		c.ToastAction("Count increased", "Undo", func() { s.count-- })
	}
	ui.Textf(c, "Count: %d", s.count)
}

func exampleProgress(c *ui.Context, s *demoState) {
	ui.Slider(c, &s.value, 0, 100).Label("Progress value")
	ui.Progress(c, s.value/100).Label("Progress")
	ui.Progress(c, -1).Label("Indeterminate progress")
}

func exampleSpinner(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(20).Children(func() {
		for _, size := range []float32{16, 24, 32} {
			ui.Spinner(c).Label("Loading").Size(size, size)
		}
	})
}

func exampleMeter(c *ui.Context, s *demoState) {
	ui.Slider(c, &s.value, 0, 100).Label("Usage value")
	ui.Meter(c, s.value, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 90}).Label("Usage")
}

func exampleBadge(c *ui.Context, s *demoState) {
	ui.Badge(c, "New")
	ui.Badge(c, "12").Background(c.Theme().Accent).TextColor(c.Theme().AccentText)
}

func exampleImage(c *ui.Context, s *demoState) {
	ui.Image(c, demoImage).Size(180, 100).Fit(ui.Contain).Label("Contain")
	ui.Image(c, demoImage).Size(120, 120).Fit(ui.Cover).Radius(12).Label("Cover")
	ui.Image(c, demoImage).Size(180, 100).Grayscale().Label("Grayscale")
}

func exampleIcon(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(16).Children(func() {
		for _, size := range []float32{16, 24, 32} {
			ui.Icon(c, lucide.Sun()).Size(size, size).TextColor(c.Theme().Accent)
		}
	})
}

func exampleAvatar(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(12).Children(func() {
		ui.Avatar(c, "Ada Lovelace", nil)
		ui.Avatar(c, "Grace Hopper", nil)
	})
}

func exampleText(c *ui.Context, s *demoState) {
	ui.Text(c, "A heading").Bold().FontSize(c.Theme().Rem(1.5))
	ui.Text(c, "Hello, world. 你好，世界。こんにちは 🌍").Wrap()
	ui.Text(c, "Secondary text").TextColor(c.Theme().TextMuted)
}

func exampleRichText(c *ui.Context, s *demoState) {
	ui.RichText(c).Children(func() {
		ui.Text(c, "Combine ")
		ui.Text(c, "bold text").Bold()
		ui.Text(c, " with ")
		if ui.Link(c, "an inline action", "").Clicked() {
			s.count++
		}
	}).FillWidth()
	ui.Textf(c, "Clicks: %d", s.count)
}

func exampleBox(c *ui.Context, s *demoState) {
	ui.Box(c).Size(120, 80).Padding(12).Radius(8).Background(c.Theme().Surface).Center().Children(func() { ui.Text(c, "Box") })
}

func exampleRow(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(8).Children(func() { ui.Badge(c, "First"); ui.Badge(c, "Second") })
}

func exampleColumn(c *ui.Context, s *demoState) {
	ui.Column(c).Gap(8).Children(func() { ui.Text(c, "First"); ui.Text(c, "Second") })
}

func exampleSpacer(c *ui.Context, s *demoState) {
	ui.Row(c).FillWidth().Children(func() { ui.Text(c, "Start"); ui.Spacer(c); ui.Text(c, "End") })
}

func exampleDivider(c *ui.Context, s *demoState) {
	ui.Text(c, "Above")
	ui.Divider(c)
	ui.Text(c, "Below")
}
