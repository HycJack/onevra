package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// demoState is owned by one page, survives navigation, and is never persisted.
type demoState struct {
	variant         string
	flyout          flyoutDemo
	glass           glassDemo
	terminal        terminalDemo
	webView         webViewDemo
	audio           audioDemo
	canvas          canvasDemo
	game            gameDemo
	lawn            lawnDemo
	code            *editor.State
	diff            *editor.DiffState
	text            string
	buffer          *ui.TextBuffer
	readOnly        bool
	choice          string
	on, other, open bool
	index, count    int
	value, high     float64
	size            float32
	offset          ui.Point
	tags            []string
	day             time.Time
	color           ui.Color
	list            ui.ListState
	grid            ui.GridState
	outline         ui.OutlineState[string]
}

func newDemoState() *demoState {
	return &demoState{text: "Hello MyGo", choice: "Small", value: 40, high: 80,
		size: 140, day: time.Date(2026, 6, 15, 9, 30, 0, 0, time.Local),
		color: ui.RGB(80, 120, 220), tags: []string{"Go"}}
}

func exampleTextInputBuffer(c *ui.Context, s *demoState) {
	ui.TextInputBuffer(c, s.buffer).Label(tr(c, "gallery.name")).Placeholder(tr(c, "gallery.enter_a_name")).ReadOnly(s.readOnly).FillWidth()
}

func exampleTextAreaBuffer(c *ui.Context, s *demoState) {
	ui.TextAreaBuffer(c, s.buffer).Label(tr(c, "gallery.notes")).Placeholder(tr(c, "gallery.write_some_notes")).ReadOnly(s.readOnly).Height(140).FillWidth()
}

func exampleSegmented(c *ui.Context, s *demoState) {
	ui.Segmented(c, &s.index, tr(c, "gallery.list"), tr(c, "gallery.grid")).Label(tr(c, "gallery.view_mode"))
	components.Text(c, tr(c, []string{"List view selected", "Grid view selected"}[s.index]))
}

func exampleCheckboxGroup(c *ui.Context, s *demoState) {
	ui.CheckboxGroup(c, tr(c, "gallery.notifications"), func() {
		components.Checkbox(c, &s.on, tr(c, "gallery.email"))
		components.Checkbox(c, &s.other, tr(c, "gallery.sound"))
	})
}

func exampleStepSlider(c *ui.Context, s *demoState) {
	ui.StepSlider(c, &s.value, 0, 100, 20).Label(tr(c, "gallery.quality"))
	ui.Textf(c, tr(c, "gallery.quality_0f"), s.value)
}

func exampleRangeSlider(c *ui.Context, s *demoState) {
	ui.RangeSlider(c, &s.value, &s.high, 0, 100, 5).Label(tr(c, "gallery.range"))
	ui.Textf(c, tr(c, "gallery.0f_0f"), s.value, s.high)
}

func exampleStepper(c *ui.Context, s *demoState) {
	ui.Stepper(c, &s.value, 0, 100, 1).Label(tr(c, "gallery.copies"))
	ui.Textf(c, tr(c, "gallery.copies_0f"), s.value)
}

func exampleRating(c *ui.Context, s *demoState) {
	ui.Rating(c, &s.index, 5).Label(tr(c, "gallery.rating"))
	ui.Textf(c, tr(c, "gallery.d_of_5_stars"), s.index)
}

func exampleNumberInput(c *ui.Context, s *demoState) {
	ui.NumberInput(c, &s.value, 0, 100, 1).Label(tr(c, "gallery.quantity"))
	ui.Textf(c, tr(c, "gallery.quantity_0f"), s.value)
}

func exampleTokenField(c *ui.Context, s *demoState) {
	ui.TokenField(c, &s.tags, []string{"Go", "UI", "Native", "Desktop"}).Label(tr(c, "gallery.tags")).FillWidth()
	components.Text(c, strings.Join(s.tags, ", ")).Wrap()
}

func exampleEditableText(c *ui.Context, s *demoState) {
	ui.EditableText(c, &s.text).Label(tr(c, "gallery.editable_name"))
}

func exampleCalendar(c *ui.Context, s *demoState) {
	ui.Calendar(c, &s.day).Label(tr(c, "gallery.day"))
	components.Text(c, s.day.Format("2006-01-02"))
}

func exampleDateInput(c *ui.Context, s *demoState) {
	ui.DateInput(c, &s.day).Label(tr(c, "gallery.date"))
	components.Text(c, s.day.Format("2006-01-02"))
}

func exampleTimeInput(c *ui.Context, s *demoState) {
	ui.TimeInput(c, &s.day).Label(tr(c, "gallery.time"))
	components.Text(c, s.day.Format("15:04"))
}

func exampleColorWell(c *ui.Context, s *demoState) {
	ui.ColorWell(c, &s.color).Label(tr(c, "gallery.tint"))
	ui.Box(c).Size(80, 40).Background(s.color)
}

func exampleColorPicker(c *ui.Context, s *demoState) {
	ui.ColorPicker(c, &s.color).Label(tr(c, "gallery.color"))
	ui.Box(c).Size(80, 40).Background(s.color)
}

func exampleList(c *ui.Context, s *demoState) {
	s.list.Selected = &s.index
	s.list.Label = func(i int) string { return fmt.Sprintf(tr(c, "gallery.item_d"), i+1) }
	ui.List(c, &s.list, 1000, func(i int) { ui.Textf(c, tr(c, "gallery.item_d"), i+1).Padding(8) }).Height(200).FillWidth().Label(tr(c, "gallery.items"))
	ui.Textf(c, tr(c, "gallery.selected_item_d"), s.index+1)
}

func exampleTable(c *ui.Context, s *demoState) {
	s.list.Selected = &s.index
	cols := []ui.TableColumn{{Title: tr(c, "gallery.name")}, {Title: tr(c, "gallery.size"), Width: 90}}
	ui.Table(c, &s.list, cols, 100, func(row, col int) {
		if col == 0 {
			ui.Textf(c, tr(c, "gallery.file_d"), row+1)
		} else {
			ui.Textf(c, tr(c, "gallery.d_kb"), row+1)
		}
	}).Height(200).FillWidth().Label(tr(c, "gallery.files"))
}

func exampleTree(c *ui.Context, s *demoState) {
	ui.Tree(c, func() {
		ui.TreeItem(c, "source", &s.open, func() {
			if ui.TreeItem(c, "main.go", nil, nil).Selected(s.on).Clicked() {
				s.on = true
			}
		})
	}).Label(tr(c, "gallery.project"))
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
	}, func(name string) { components.Text(c, name) }).Height(180).FillWidth().Label(tr(c, "gallery.project_outline"))
}

func exampleOutlineTable(c *ui.Context, s *demoState) {
	s.outline.List.Selected = &s.index
	cols := []ui.TableColumn{{Title: tr(c, "gallery.name")}, {Title: tr(c, "gallery.kind"), Width: 90}}
	ui.OutlineTable(c, &s.outline, cols, []string{"source"}, func(name string) []string {
		if name == "source" {
			return []string{"main.go", "go.mod"}
		}
		return nil
	}, func(name string, col int) {
		if col == 0 {
			components.Text(c, name)
		} else {
			components.Text(c, tr(c, "gallery.file"))
		}
	}).Height(180).FillWidth().Label(tr(c, "gallery.project_table"))
}

func exampleGridView(c *ui.Context, s *demoState) {
	s.grid.Selected = &s.index
	s.grid.Label = func(i int) string { return fmt.Sprintf(tr(c, "gallery.item_d"), i+1) }
	ui.GridView(c, &s.grid, 1000, 90, 70, func(i int) {
		ui.Textf(c, tr(c, "gallery.item_d"), i+1).Padding(12)
	}).Height(220).FillWidth().Label(tr(c, "gallery.items"))
	ui.Textf(c, tr(c, "gallery.selected_item_d"), s.index+1)
}

func exampleScroll(c *ui.Context, s *demoState) {
	ui.Scroll(c).Height(180).FillWidth().Children(func() {
		for i := range 30 {
			ui.Textf(c, tr(c, "gallery.line_d"), i+1).Padding(4)
		}
	})
}

func exampleScrollHorizontal(c *ui.Context, s *demoState) {
	ui.ScrollHorizontal(c).Height(80).FillWidth().Children(func() {
		ui.Row(c).Gap(12).Children(func() {
			for i := range 20 {
				ui.Textf(c, tr(c, "gallery.card_d"), i+1).Width(120).Padding(16).Background(c.Theme().Surface)
			}
		})
	})
}

func exampleScrollBoth(c *ui.Context, s *demoState) {
	ui.ScrollBoth(c).Height(180).FillWidth().Children(func() {
		ui.Column(c).Width(800).Children(func() {
			for i := range 30 {
				ui.Textf(c, tr(c, "gallery.row_d_name_status_updated"), i+1).Padding(8)
			}
		})
	})
}

func exampleGrid(c *ui.Context, s *demoState) {
	ui.Grid(c).Columns(2).Gap(8).FillWidth().Children(func() {
		for i := range 6 {
			ui.Textf(c, tr(c, "gallery.cell_d"), i+1).Padding(12).Background(c.Theme().Surface)
		}
	})
}

func exampleSplit(c *ui.Context, s *demoState) {
	ui.Split(c, &s.size, func() { components.Text(c, tr(c, "gallery.left_pane")).Wrap() }, func() { components.Text(c, tr(c, "gallery.right_pane")).Wrap() }).Height(180).FillWidth()
}

func exampleSplitVertical(c *ui.Context, s *demoState) {
	ui.SplitVertical(c, &s.size, func() { components.Text(c, tr(c, "gallery.top_pane")) }, func() { components.Text(c, tr(c, "gallery.bottom_pane")) }).Height(240).FillWidth()
}

func exampleCollapsible(c *ui.Context, s *demoState) {
	ui.Collapsible(c, tr(c, "gallery.details"), &s.open, func() { components.Checkbox(c, &s.on, tr(c, "gallery.notifications")) })
}

func exampleAccordion(c *ui.Context, s *demoState) {
	ui.Accordion(c, func() {
		ui.AccordionItem(c, tr(c, "gallery.general"), &s.open, func() { components.Checkbox(c, &s.on, tr(c, "gallery.notifications")) })
		ui.AccordionItem(c, tr(c, "gallery.about"), &s.other, func() { components.Text(c, tr(c, "gallery.built_with_mygo")) })
	}).FillWidth()
}

func exampleSidebar(c *ui.Context, s *demoState) {
	ui.Sidebar(c, &s.choice, func() {
		ui.SidebarSection(c, tr(c, "gallery.locations"), nil, func() {
			ui.SidebarItem(c, tr(c, "gallery.home"), lucide.House(), "Home")
			ui.SidebarItem(c, tr(c, "gallery.files"), lucide.Folder(), "Files")
		})
	}).Height(180).FillWidth()
	components.Text(c, fmt.Sprintf(tr(c, "gallery.selected_s"), tr(c, s.choice)))
}

func exampleBreadcrumbs(c *ui.Context, s *demoState) {
	path := []string{tr(c, "gallery.home"), tr(c, "gallery.documents"), tr(c, "gallery.project")}
	ui.Breadcrumbs(c, path, &s.index).Label(tr(c, "gallery.path"))
	components.Text(c, fmt.Sprintf(tr(c, "gallery.chosen_s"), path[s.index]))
}

func exampleRadioGroup(c *ui.Context, s *demoState) {
	ui.RadioGroup(c, func() {
		components.Radio(c, &s.choice, "Small", tr(c, "gallery.small"))
		components.Radio(c, &s.choice, "Large", tr(c, "gallery.large"))
	}).Row().Gap(16).Label(tr(c, "gallery.size_options"))
	components.Text(c, tr(c, s.choice))
}

func exampleModal(c *ui.Context, s *demoState) {
	if components.Button(c, tr(c, "gallery.open_dialog"), components.ButtonSecondary).Clicked() {
		s.open = true
	}
	ui.Modal(c, &s.open, func() {
		components.Text(c, tr(c, "gallery.edit_name")).Bold()
		components.TextBox(c, &s.text, components.TextBoxOptions{}).Label(tr(c, "gallery.name"))
		if components.Button(c, tr(c, "gallery.done"), components.ButtonSecondary).Clicked() {
			s.open = false
		}
	})
}

func exampleAlertDialog(c *ui.Context, s *demoState) {
	if components.Button(c, tr(c, "gallery.ask_for_confirmation"), components.ButtonSecondary).Clicked() {
		s.open = true
	}
	if ui.AlertDialog(c, &s.open, tr(c, "gallery.apply_change"), tr(c, "gallery.apply_these_settings"), tr(c, "gallery.cancel"), tr(c, "gallery.apply")) == 1 {
		s.count++
	}
	ui.Textf(c, tr(c, "gallery.confirmed_d"), s.count)
}

func exampleTooltip(c *ui.Context, s *demoState) {
	components.Button(c, tr(c, "gallery.hover_or_focus_me"), components.ButtonSecondary).Tooltip(tr(c, "gallery.a_helpful_description_of_this_action"))
}

func exampleToast(c *ui.Context, s *demoState) {
	if components.Button(c, tr(c, "gallery.show_toast"), components.ButtonSecondary).Clicked() {
		c.Toast(tr(c, "gallery.saved"))
	}
	if components.Button(c, tr(c, "gallery.show_undo_toast"), components.ButtonSecondary).Clicked() {
		s.count++
		c.ToastAction(tr(c, "gallery.count_increased"), tr(c, "gallery.undo"), func() { s.count-- })
	}
	ui.Textf(c, tr(c, "gallery.count_d"), s.count)
}

func exampleProgress(c *ui.Context, s *demoState) {
	components.Slider(c, &s.value, 0, 100).Label(tr(c, "gallery.progress_value"))
	ui.Progress(c, s.value/100).Label(tr(c, "gallery.progress"))
	ui.Progress(c, -1).Label(tr(c, "gallery.indeterminate_progress"))
}

func exampleSpinner(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(20).Children(func() {
		for _, size := range []float32{16, 24, 32} {
			ui.Spinner(c).Label(tr(c, "gallery.loading")).Size(size, size)
		}
	})
}

func exampleBadge(c *ui.Context, s *demoState) {
	ui.Badge(c, tr(c, "gallery.new"))
	ui.Badge(c, tr(c, "gallery.12")).Background(c.Theme().Accent).TextColor(c.Theme().AccentText)
}

func exampleImage(c *ui.Context, s *demoState) {
	ui.Image(c, demoImage).Size(180, 100).Fit(ui.Contain).Label(tr(c, "gallery.contain"))
	ui.Image(c, demoImage).Size(120, 120).Fit(ui.Cover).Radius(12).Label(tr(c, "gallery.cover"))
	ui.Image(c, demoImage).Size(180, 100).Grayscale().Label(tr(c, "gallery.grayscale"))
}

func exampleRichText(c *ui.Context, s *demoState) {
	ui.RichText(c).Children(func() {
		components.Text(c, tr(c, "gallery.combine"))
		components.Text(c, tr(c, "gallery.bold_text")).Bold()
		components.Text(c, tr(c, "gallery.with"))
		if components.Button(c, tr(c, "gallery.an_inline_action"), components.ButtonLink).Clicked() {
			s.count++
		}
	}).FillWidth()
	ui.Textf(c, tr(c, "gallery.clicks_d"), s.count)
}

func exampleBox(c *ui.Context, s *demoState) {
	ui.Box(c).Size(120, 80).Padding(12).Radius(8).Background(c.Theme().Surface).Center().Children(func() { components.Text(c, tr(c, "gallery.box")) })
}

func exampleRow(c *ui.Context, s *demoState) {
	ui.Row(c).Gap(8).Children(func() { components.Badge(c, tr(c, "gallery.first")); components.Badge(c, tr(c, "gallery.second")) })
}

func exampleColumn(c *ui.Context, s *demoState) {
	ui.Column(c).Gap(8).Children(func() { components.Text(c, tr(c, "gallery.first")); components.Text(c, tr(c, "gallery.second")) })
}

func exampleSpacer(c *ui.Context, s *demoState) {
	ui.Row(c).FillWidth().Children(func() {
		components.Text(c, tr(c, "gallery.start"))
		ui.Spacer(c)
		components.Text(c, tr(c, "gallery.end"))
	})
}

func exampleDivider(c *ui.Context, s *demoState) {
	components.Text(c, tr(c, "gallery.above"))
	ui.Divider(c)
	components.Text(c, tr(c, "gallery.below"))
}
