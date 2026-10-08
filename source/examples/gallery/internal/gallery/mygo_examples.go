package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// demoState is owned by one page, survives navigation, and is never persisted.
type demoState struct {
	terminal        terminalDemo
	text            string
	buffer          *ui.TextBuffer
	readOnly        bool
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

func exampleTextInputBuffer(c *ui.Context, s *demoState) {
	ui.TextInputBuffer(c, s.buffer).Label("Name").Placeholder("Enter a name").ReadOnly(s.readOnly).FillWidth()
}

func exampleTextAreaBuffer(c *ui.Context, s *demoState) {
	ui.TextAreaBuffer(c, s.buffer).Label("Notes").Placeholder("Write some notes").ReadOnly(s.readOnly).Height(140).FillWidth()
}

func exampleContextMenu(c *ui.Context, s *demoState) {
	components.Button(c, "Right-click or Shift+F10", components.ButtonSecondary).ContextMenu(func(m *ui.Menu) {
		if m.Item("Increment").Chosen() {
			s.count++
		}
		m.Separator()
		m.Item("Unavailable").Disabled(true)
	})
	components.Textf(c, "Count: %d", s.count)
}

func exampleSegmented(c *ui.Context, s *demoState) {
	ui.Segmented(c, &s.index, "List", "Grid").Label("View mode")
	components.Text(c, []string{"List view selected", "Grid view selected"}[s.index])
}

func exampleToolbar(c *ui.Context, s *demoState) {
	ui.Toolbar(c, func() {
		components.ToggleButton(c, &s.on, "Bold")
		ui.Spacer(c)
		if components.Button(c, "Apply", components.ButtonSecondary).Clicked() {
			s.count++
		}
	}).Label("Formatting").FillWidth()
	components.Textf(c, "Applied %d times", s.count)
}

func exampleCheckboxGroup(c *ui.Context, s *demoState) {
	ui.CheckboxGroup(c, "Notifications", func() {
		components.Checkbox(c, &s.on, "Email")
		components.Checkbox(c, &s.other, "Sound")
	})
}

func exampleStepSlider(c *ui.Context, s *demoState) {
	ui.StepSlider(c, &s.value, 0, 100, 20).Label("Quality")
	components.Textf(c, "Quality: %.0f", s.value)
}

func exampleRangeSlider(c *ui.Context, s *demoState) {
	ui.RangeSlider(c, &s.value, &s.high, 0, 100, 5).Label("Range")
	components.Textf(c, "%.0f – %.0f", s.value, s.high)
}

func exampleStepper(c *ui.Context, s *demoState) {
	ui.Stepper(c, &s.value, 0, 100, 1).Label("Copies")
	components.Textf(c, "Copies: %.0f", s.value)
}

func exampleRating(c *ui.Context, s *demoState) {
	ui.Rating(c, &s.index, 5).Label("Rating")
	components.Textf(c, "%d of 5 stars", s.index)
}

func exampleNumberInput(c *ui.Context, s *demoState) {
	ui.NumberInput(c, &s.value, 0, 100, 1).Label("Quantity")
	components.Textf(c, "Quantity: %.0f", s.value)
}

func exampleSearchField(c *ui.Context, s *demoState) {
	ui.SearchField(c, &s.text).Label("Search samples").FillWidth()
	for _, name := range []string{"Hello MyGo", "Hello Go", "Native UI"} {
		if strings.Contains(strings.ToLower(name), strings.ToLower(s.text)) {
			components.Text(c, name)
		}
	}
}

func exampleCombobox(c *ui.Context, s *demoState) {
	ui.Combobox(c, &s.choice, []string{"Small", "Medium", "Large"}).Label("Size")
	components.Text(c, s.choice)
}

func exampleAutocomplete(c *ui.Context, s *demoState) {
	ui.Autocomplete(c, &s.text, []string{"Hello MyGo", "Hello Go", "Hello world"}).Label("Greeting")
	components.Text(c, s.text).Wrap()
}

func exampleTokenField(c *ui.Context, s *demoState) {
	ui.TokenField(c, &s.tags, []string{"Go", "UI", "Native", "Desktop"}).Label("Tags").FillWidth()
	components.Text(c, strings.Join(s.tags, ", ")).Wrap()
}

func exampleEditableText(c *ui.Context, s *demoState) {
	ui.EditableText(c, &s.text).Label("Editable name")
}

func exampleFindBar(c *ui.Context, s *demoState) {
	if components.Button(c, "Find", components.ButtonSecondary).Clicked() {
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
	components.Text(c, strings.Join(lines, "\n")).Wrap()
	if len(matches) > 0 {
		components.Text(c, "Match: "+matches[min(max(s.index, 0), len(matches)-1)]).Wrap()
	}
}

func exampleCalendar(c *ui.Context, s *demoState) {
	ui.Calendar(c, &s.day).Label("Day")
	components.Text(c, s.day.Format("2006-01-02"))
}

func exampleDateInput(c *ui.Context, s *demoState) {
	ui.DateInput(c, &s.day).Label("Date")
	components.Text(c, s.day.Format("2006-01-02"))
}

func exampleTimeInput(c *ui.Context, s *demoState) {
	ui.TimeInput(c, &s.day).Label("Time")
	components.Text(c, s.day.Format("15:04"))
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
	ui.List(c, &s.list, 1000, func(i int) { components.Textf(c, "Item %d", i+1).Padding(8) }).Height(200).FillWidth().Label("Items")
	components.Textf(c, "Selected item: %d", s.index+1)
}

func exampleTable(c *ui.Context, s *demoState) {
	s.list.Selected = &s.index
	cols := []ui.TableColumn{{Title: "Name"}, {Title: "Size", Width: 90}}
	ui.Table(c, &s.list, cols, 100, func(row, col int) {
		if col == 0 {
			components.Textf(c, "File %d", row+1)
		} else {
			components.Textf(c, "%d KB", row+1)
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
	}, func(name string) { components.Text(c, name) }).Height(180).FillWidth().Label("Project outline")
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
			components.Text(c, name)
		} else {
			components.Text(c, "File")
		}
	}).Height(180).FillWidth().Label("Project table")
}

func exampleGridView(c *ui.Context, s *demoState) {
	s.grid.Selected = &s.index
	s.grid.Label = func(i int) string { return fmt.Sprintf("Item %d", i+1) }
	ui.GridView(c, &s.grid, 1000, 90, 70, func(i int) {
		components.Textf(c, "Item %d", i+1).Padding(12)
	}).Height(220).FillWidth().Label("Items")
	components.Textf(c, "Selected item: %d", s.index+1)
}

func exampleScroll(c *ui.Context, s *demoState) {
	ui.Scroll(c).Height(180).FillWidth().Children(func() {
		for i := range 30 {
			components.Textf(c, "Line %d", i+1).Padding(4)
		}
	})
}

func exampleScrollHorizontal(c *ui.Context, s *demoState) {
	ui.ScrollHorizontal(c).Height(80).FillWidth().Children(func() {
		ui.Row(c).Gap(12).Children(func() {
			for i := range 20 {
				components.Textf(c, "Card %d", i+1).Width(120).Padding(16).Background(c.Theme().Surface)
			}
		})
	})
}

func exampleScrollBoth(c *ui.Context, s *demoState) {
	ui.ScrollBoth(c).Height(180).FillWidth().Children(func() {
		ui.Column(c).Width(800).Children(func() {
			for i := range 30 {
				components.Textf(c, "Row %d — Name                         Status                         Updated", i+1).Padding(8)
			}
		})
	})
}

func exampleGrid(c *ui.Context, s *demoState) {
	ui.Grid(c).Columns(2).Gap(8).FillWidth().Children(func() {
		for i := range 6 {
			components.Textf(c, "Cell %d", i+1).Padding(12).Background(c.Theme().Surface)
		}
	})
}

func exampleSplit(c *ui.Context, s *demoState) {
	ui.Split(c, &s.size, func() { components.Text(c, "Left pane").Wrap() }, func() { components.Text(c, "Right pane").Wrap() }).Height(180).FillWidth()
}

func exampleSplitVertical(c *ui.Context, s *demoState) {
	ui.SplitVertical(c, &s.size, func() { components.Text(c, "Top pane") }, func() { components.Text(c, "Bottom pane") }).Height(240).FillWidth()
}

func exampleCollapsible(c *ui.Context, s *demoState) {
	ui.Collapsible(c, "Details", &s.open, func() { components.Checkbox(c, &s.on, "Notifications") })
}

func exampleAccordion(c *ui.Context, s *demoState) {
	ui.Accordion(c, func() {
		ui.AccordionItem(c, "General", &s.open, func() { components.Checkbox(c, &s.on, "Notifications") })
		ui.AccordionItem(c, "About", &s.other, func() { components.Text(c, "Built with MyGo") })
	}).FillWidth()
}

func exampleForm(c *ui.Context, s *demoState) {
	ui.Form(c, func() {
		ui.Field(c, "Name", func() { components.TextField(c, &s.text) })
		ui.Field(c, "Sound", func() { components.Switch(c, &s.on, "Sound") })
	}).FillWidth()
}

func exampleField(c *ui.Context, s *demoState) {
	message := ""
	if s.text != "" && !strings.Contains(s.text, "@") {
		message = "Enter an email address."
	}
	ui.Field(c, "Email", func() { components.TextField(c, &s.text).FillWidth() }).Description("Email address").Error(message)
}

func exampleFieldset(c *ui.Context, s *demoState) {
	ui.Form(c, func() {
		ui.Fieldset(c, "Notifications", func() {
			ui.Field(c, "Email", func() { components.Checkbox(c, &s.on, "Weekly summary") })
			ui.Field(c, "Sound", func() { components.Switch(c, &s.other, "Sound") })
		})
	})
}

func exampleSidebar(c *ui.Context, s *demoState) {
	ui.Sidebar(c, &s.choice, func() {
		ui.SidebarSection(c, "Locations", nil, func() {
			ui.SidebarItem(c, "Home", lucide.House(), "Home")
			ui.SidebarItem(c, "Files", lucide.Folder(), "Files")
		})
	}).Height(180).FillWidth()
	components.Text(c, "Selected: "+s.choice)
}

func exampleBreadcrumbs(c *ui.Context, s *demoState) {
	path := []string{"Home", "Documents", "Project"}
	ui.Breadcrumbs(c, path, &s.index).Label("Path")
	components.Text(c, "Chosen: "+path[s.index])
}

func exampleBackButton(c *ui.Context, s *demoState) {
	ui.Toolbar(c, func() {
		ui.BackButton(c, s.router)
		ui.ForwardButton(c, s.router)
	}).Label("History")
	if components.Button(c, "Visit next page", components.ButtonSecondary).Clicked() {
		s.count++
		s.router.Push(fmt.Sprintf("/page/%d", s.count))
	}
	s.router.View(c, func(r *ui.Route) { r.Title(r.Path()); components.Text(c, r.Path()) }).Height(60).FillWidth()
}

func exampleForwardButton(c *ui.Context, s *demoState) {
	if components.Button(c, "Prepare forward history", components.ButtonSecondary).Clicked() {
		s.count++
		s.router.Push(fmt.Sprintf("/page/%d", s.count))
		s.router.Back()
	}
	ui.ForwardButton(c, s.router)
	s.router.View(c, func(r *ui.Route) {
		r.Title(r.Path())
		components.Text(c, r.Path())
	}).Height(60).FillWidth()
}

func exampleRadioGroup(c *ui.Context, s *demoState) {
	ui.RadioGroup(c, func() {
		components.Radio(c, &s.choice, "Small", "Small")
		components.Radio(c, &s.choice, "Large", "Large")
	}).Row().Gap(16).Label("Size options")
	components.Text(c, s.choice)
}

func exampleModal(c *ui.Context, s *demoState) {
	if components.Button(c, "Open dialog", components.ButtonSecondary).Clicked() {
		s.open = true
	}
	ui.Modal(c, &s.open, func() {
		components.Text(c, "Edit name").Bold()
		components.TextField(c, &s.text).Label("Name")
		if components.Button(c, "Done", components.ButtonSecondary).Clicked() {
			s.open = false
		}
	})
}

func exampleAlertDialog(c *ui.Context, s *demoState) {
	if components.Button(c, "Ask for confirmation", components.ButtonSecondary).Clicked() {
		s.open = true
	}
	if ui.AlertDialog(c, &s.open, "Apply change?", "Apply these settings?", "Cancel", "Apply") == 1 {
		s.count++
	}
	components.Textf(c, "Confirmed: %d", s.count)
}

func examplePopover(c *ui.Context, s *demoState) {
	anchor := components.Button(c, "Show popover", components.ButtonSecondary)
	if anchor.Clicked() {
		s.open = !s.open
	}
	ui.Popover(c, anchor, &s.open, func() {
		components.Text(c, "Popover content")
		if components.Button(c, "Close", components.ButtonSecondary).Clicked() {
			s.open = false
		}
	})
}

func exampleTooltip(c *ui.Context, s *demoState) {
	components.Button(c, "Hover or focus me", components.ButtonSecondary).Tooltip("A helpful description of this action")
}

func exampleToast(c *ui.Context, s *demoState) {
	if components.Button(c, "Show toast", components.ButtonSecondary).Clicked() {
		c.Toast("Saved")
	}
	if components.Button(c, "Show undo toast", components.ButtonSecondary).Clicked() {
		s.count++
		c.ToastAction("Count increased", "Undo", func() { s.count-- })
	}
	components.Textf(c, "Count: %d", s.count)
}

func exampleProgress(c *ui.Context, s *demoState) {
	components.Slider(c, &s.value, 0, 100).Label("Progress value")
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
	components.Slider(c, &s.value, 0, 100).Label("Usage value")
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

func exampleRichText(c *ui.Context, s *demoState) {
	ui.RichText(c).Children(func() {
		components.Text(c, "Combine ")
		components.Text(c, "bold text").Bold()
		components.Text(c, " with ")
		if components.Button(c, "an inline action", components.ButtonLink).Clicked() {
			s.count++
		}
	}).FillWidth()
	components.Textf(c, "Clicks: %d", s.count)
}

func exampleBox(c *ui.Context, s *demoState) {
	ui.Box(c).Size(120, 80).Padding(12).Radius(8).Background(c.Theme().Surface).Center().Children(func() { components.Text(c, "Box") })
}

func exampleRow(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(8).Children(func() { ui.Badge(c, "First"); ui.Badge(c, "Second") })
}

func exampleColumn(c *ui.Context, s *demoState) {
	ui.Column(c).Gap(8).Children(func() { components.Text(c, "First"); components.Text(c, "Second") })
}

func exampleSpacer(c *ui.Context, s *demoState) {
	ui.Row(c).FillWidth().Children(func() { components.Text(c, "Start"); ui.Spacer(c); components.Text(c, "End") })
}

func exampleDivider(c *ui.Context, s *demoState) {
	components.Text(c, "Above")
	ui.Divider(c)
	components.Text(c, "Below")
}
