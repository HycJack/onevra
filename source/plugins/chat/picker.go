package chat

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// Option is one choice in a Picker. ID is stable; Group, when set, gathers
// options under a heading in slice order; Detail is secondary text such as a
// context size; Description is a muted second line.
type Option struct {
	ID, Label, Detail, Description, Group string
}

// PickerState retains whether the popover is open and the search text.
type PickerState struct {
	Open  bool
	Query string
}

// PickerOptions describes a chip-triggered popover choice. Label is the
// trigger text and accessible name; SearchLabel enables filtering when set
// and there are more than Searchable options (default 8).
type PickerOptions struct {
	Label, SearchLabel, EmptyText string
	Icon                          func() *ui.SVG
	Options                       []Option
	Selected                      string
	Width                         float32
	Searchable                    int
}

// PickerResult reports a frame-local choice.
type PickerResult struct {
	Trigger ui.Element
	Chosen  string
}

// Picker shows a Chip that opens a popover of grouped, optionally searchable
// options with a mark on the selection. Choosing closes the popover; the
// application applies and persists the choice.
func Picker(c *ui.Context, state *PickerState, options PickerOptions) PickerResult {
	var result PickerResult
	t := c.Theme()
	result.Trigger = Chip(c, options.Icon, options.Label, true)
	if state.Open {
		result.Trigger.Background(overlayStrong(t))
	}
	if result.Trigger.Clicked() {
		state.Open, state.Query = !state.Open, ""
	}
	width := options.Width
	if width <= 0 {
		width = 260
	}
	threshold := options.Searchable
	if threshold <= 0 {
		threshold = 8
	}
	ui.PopoverBase(c, result.Trigger, &state.Open, func(panel ui.Element) {
		panel.Column().Width(width).Padding(4).Gap(2).Radius(10).Border(1, strongBorder(t)).Background(t.Surface).TextColor(t.Text).
			Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, .16)).Children(func() {
			if options.SearchLabel != "" && len(options.Options) > threshold {
				ui.Row(c).FillWidth().Height(30).Padding(0, 8).Gap(6).AlignItems(ui.Center).Margin(0, 0, 2, 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
					components.Icon(c, lucide.Search).Size(12, 12).Shrink(0).TextColor(tertiary(t))
					components.TextBox(c, &state.Query, components.TextBoxOptions{}).Grow(1).MinWidth(0).AutoFocus().
						Label(options.SearchLabel).Placeholder(options.SearchLabel).Border(0, ui.Color{}).Background(ui.Color{}).Padding(0).FontSize(12.5)
				})
			}
			query := strings.ToLower(strings.TrimSpace(state.Query))
			var visible []Option
			for _, o := range options.Options {
				if query == "" || strings.Contains(strings.ToLower(o.Label+" "+o.Group+" "+o.Detail+" "+o.ID), query) {
					visible = append(visible, o)
				}
			}
			if len(visible) == 0 {
				components.Text(c, options.EmptyText).FontSize(12.5).TextColor(tertiary(t)).Padding(8)
				return
			}
			components.ScrollView(c, components.ScrollViewOptions{}).FillWidth().MaxHeight(340).Gap(1).Children(func() {
				group := ""
				for i, o := range visible {
					if o.Group != "" && (i == 0 || o.Group != group) {
						components.Text(c, o.Group).FontSize(11.5).FontWeight(500).TextColor(tertiary(t)).Padding(8, 8, 3, 8)
					}
					group = o.Group
					if pickerRow(c, o, o.ID == options.Selected) {
						result.Chosen = o.ID
						state.Open = false
					}
				}
			})
		})
	})
	return result
}

func pickerRow(c *ui.Context, o Option, selected bool) bool {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(o.ID).Label(o.Label).Selected(selected).FillWidth().MinWidth(0).Padding(6, 8).Radius(6).
		Justify(ui.Start).AlignItems(ui.Center).Gap(8).Background(ui.Color{}).TextColor(t.Text)
	if b.Hovered() || b.Focused() {
		b.Background(overlay(t))
	}
	b.Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
			components.Text(c, o.Label).MinWidth(0).SingleLine().Ellipsis("…").FontSize(13)
			if o.Description != "" {
				components.Text(c, o.Description).MinWidth(0).Wrap().FontSize(11.5).TextColor(tertiary(t))
			}
		})
		if o.Detail != "" {
			components.Text(c, o.Detail).Shrink(0).FontSize(11.5).TextColor(tertiary(t))
		}
		icon := ui.Box(c).Size(12, 12).Shrink(0)
		if selected {
			icon.Children(func() { components.Icon(c, lucide.Check).Size(12, 12).TextColor(t.Accent) })
		}
	})
	return b.Clicked()
}
