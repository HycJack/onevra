package chat

import (
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// Chip is a compact trigger for composer and workspace settings such as a
// model, reasoning level, project or branch. Attach a native menu with
// Element.Menu or read Clicked; disable it when the value is fixed. Icon may
// be nil; Caret adds a dropdown mark.
func Chip(c *ui.Context, icon func() *ui.SVG, label string, caret bool) ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Label(label).Height(26).MinWidth(0).Padding(0, 7).Radius(6).AlignItems(ui.Center).Gap(6).
		Background(ui.Color{}).TextColor(secondary(t))
	if b.Hovered() || b.Focused() {
		b.Background(overlay(t))
	}
	if b.IsDisabled() {
		b.Opacity(.7)
	}
	b.Children(func() {
		if icon != nil {
			components.Icon(c, icon).Size(12, 12).Shrink(0).TextColor(tertiary(t))
		}
		components.Text(c, label).MinWidth(0).SingleLine().Ellipsis("…").FontSize(12.5).TextColor(secondary(t))
		if caret {
			components.Icon(c, lucide.ChevronDown).Size(10.5, 10.5).Shrink(0).TextColor(tertiary(t).Alpha(.7))
		}
	})
	return b
}

// AttachmentOptions describes a file the application has already prepared.
// Image is an optional decoded thumbnail; Icon replaces the default file icon.
// An empty RemoveLabel makes the attachment read-only, as in sent messages.
type AttachmentOptions struct {
	Label, Detail, RemoveLabel string
	Image                      ui.ImageSource
	Icon                       func() *ui.SVG
}

// AttachmentResult reports a frame-local removal request.
type AttachmentResult struct {
	Root   ui.Element
	Remove bool
}

// Attachment shows a compact file or image chip. It performs no file I/O.
func Attachment(c *ui.Context, options AttachmentOptions) AttachmentResult {
	var result AttachmentResult
	t := c.Theme()
	result.Root = ui.Row(c).MinWidth(0).MaxWidth(220).Height(30).AlignItems(ui.Center).Gap(6).Padding(0, 4, 0, 4).
		Radius(7).Border(1, strongBorder(t)).Background(t.Background.Mix(t.Text, .035)).Children(func() {
		if options.Image != nil {
			ui.Image(c, options.Image).Size(22, 22).Fit(ui.Cover).Radius(4).Shrink(0)
		} else {
			icon := options.Icon
			if icon == nil {
				icon = lucide.FileText
			}
			ui.Row(c).Size(22, 22).Shrink(0).Center().Children(func() {
				components.Icon(c, icon).Size(13, 13).TextColor(tertiary(t))
			})
		}
		components.Text(c, options.Label).MinWidth(0).SingleLine().Ellipsis("…").FontSize(12.5).TextColor(secondary(t))
		if options.Detail != "" {
			components.Text(c, options.Detail).Shrink(0).FontSize(11.5).TextColor(tertiary(t))
		}
		if options.RemoveLabel != "" {
			b := ui.ButtonBase(c).Label(options.RemoveLabel).Size(18, 18).Shrink(0).Radius(4).Center().Background(ui.Color{}).TextColor(tertiary(t))
			if b.Hovered() || b.Focused() {
				b.Background(overlayStrong(t)).TextColor(t.Text)
			}
			b.Children(func() { components.Icon(c, lucide.X).Size(11, 11) })
			result.Remove = b.Clicked()
		} else {
			ui.Box(c).Width(2)
		}
	})
	return result
}

// ContextMeterOptions describes context-window occupancy. Fraction is 0..1,
// or negative when nothing is measured. Label is the trigger's accessible
// name and Tooltip its hover text; Title and Value head the panel, such as
// "Context window" and "61k / 272k (22%)".
type ContextMeterOptions struct {
	Label, Tooltip, Title, Value string
	Fraction                     float64
}

// MeterState retains whether the panel is open.
type MeterState struct{ Open bool }

// ContextMeter is a small ring gauge that opens a panel with the measured
// value and bar. It warns from 80% and alerts from 95%; it measures nothing.
func ContextMeter(c *ui.Context, state *MeterState, options ContextMeterOptions) ui.Element {
	t := c.Theme()
	fill := t.Accent
	switch {
	case options.Fraction >= .95:
		fill = t.Danger
	case options.Fraction >= .8:
		fill = t.Warning
	}
	b := ui.ButtonBase(c).Label(options.Label).Height(20).Padding(0, 5).Radius(5).Shrink(0).Center().Background(ui.Color{})
	if state.Open {
		b.Background(overlayStrong(t))
	} else if b.Hovered() || b.Focused() {
		b.Background(overlay(t))
	}
	b.Children(func() { gauge(c, options.Fraction, strongBorder(t).Mix(t.Text, .1), fill) })
	if options.Tooltip != "" {
		components.Tooltip(c, b, options.Tooltip)
	}
	if b.Clicked() {
		state.Open = !state.Open
	}
	ui.PopoverBase(c, b, &state.Open, func(panel ui.Element) {
		panel.Column().Width(320).Padding(14).Gap(7).Radius(10).Border(1, strongBorder(t)).Background(t.Surface).TextColor(t.Text).
			Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, .16)).Children(func() {
			ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(8).Children(func() {
				components.Text(c, options.Title).Grow(1).MinWidth(0).SingleLine().Ellipsis("…").FontSize(12.5)
				components.Text(c, options.Value).Shrink(0).FontSize(12.5).TextColor(tertiary(t))
			})
			value := options.Fraction
			if value < 0 {
				value = 0
			}
			ui.Box(c).FillWidth().Height(6).Radius(3).Background(overlayStrong(t)).Draw(func(p *ui.Painter, r ui.Rect) {
				if w := r.W * float32(min(value, 1)); w > 0 {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, fill, r.H/2)
				}
			})
		})
	})
	return b
}

func gauge(c *ui.Context, fraction float64, track, fill ui.Color) {
	ui.Box(c).Size(13, 13).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
		const stroke = 2.5
		radius := (r.W - stroke) / 2
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		var ring ui.Path
		ring.Circle(cx, cy, radius)
		p.StrokePath(&ring, stroke, track)
		if fraction <= 0 {
			return
		}
		var arc ui.Path
		if fraction >= 1 {
			arc.Circle(cx, cy, radius)
		} else {
			arcPath(&arc, cx, cy, radius, fraction)
		}
		p.StrokePath(&arc, stroke, fill)
	})
}

// ThreadItemOptions describes one conversation entry in a navigation list.
// ID is stable; Detail is typically the project and Time a relative label.
// Running shows activity, Waiting and Failed show attention marks; any other
// status is quiet.
type ThreadItemOptions struct {
	ID, Title, Detail, Time string
	DetailIcon              func() *ui.SVG
	Status                  Status
	Selected                bool
}

// ThreadItem builds a two-line, keyboard-operable conversation row. Read
// Clicked for activation; selection and ordering remain application-owned.
func ThreadItem(c *ui.Context, options ThreadItemOptions) ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Key(options.ID).Label(options.Title).FillWidth().MinWidth(0).Column().AlignItems(ui.Stretch).
		Padding(7, 8).Gap(4).Radius(7).Selected(options.Selected).Background(ui.Color{}).TextColor(t.Text)
	if options.Selected || b.Hovered() || b.Focused() {
		b.Background(t.Text.Alpha(.07))
	}
	b.Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(6).Children(func() {
			components.Text(c, options.Title).Grow(1).MinWidth(0).SingleLine().Ellipsis("...").FontSize(13.5)
			switch options.Status {
			case Running:
				components.ProgressRing(c, -1).Size(12, 12).Shrink(0)
			case Waiting:
				components.Icon(c, lucide.CircleAlert).Size(12, 12).Shrink(0).TextColor(t.Warning)
			case Failed:
				components.Icon(c, lucide.X).Size(12, 12).Shrink(0).TextColor(t.Danger)
			}
		})
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(5).Children(func() {
			if options.DetailIcon != nil {
				components.Icon(c, options.DetailIcon).Size(12.5, 12.5).Shrink(0).TextColor(tertiary(t))
			}
			components.Text(c, options.Detail).Grow(1).MinWidth(0).SingleLine().Ellipsis("…").FontSize(13).TextColor(tertiary(t))
			if options.Time != "" {
				color := tertiary(t).Alpha(.7)
				if options.Status == Running {
					color = tertiary(t)
				}
				components.Text(c, options.Time).Shrink(0).FontSize(12.5).TextColor(color)
			}
		})
	})
	return b
}

// arcPath draws a clockwise arc from 12 o'clock with cubic segments.
func arcPath(path *ui.Path, cx, cy, radius float32, fraction float64) {
	start, sweep := -math.Pi/2, fraction*2*math.Pi
	point := func(a float64) (float32, float32) {
		return cx + radius*float32(math.Cos(a)), cy + radius*float32(math.Sin(a))
	}
	x, y := point(start)
	path.MoveTo(x, y)
	steps := int(math.Ceil(sweep / (math.Pi / 2)))
	delta := sweep / float64(steps)
	k := float32(4.0/3*math.Tan(delta/4)) * radius
	for i := 0; i < steps; i++ {
		a, b := start+float64(i)*delta, start+float64(i+1)*delta
		x0, y0 := point(a)
		x1, y1 := point(b)
		path.CubeTo(x0-k*float32(math.Sin(a)), y0+k*float32(math.Cos(a)), x1+k*float32(math.Sin(b)), y1-k*float32(math.Cos(b)), x1, y1)
	}
}

// SidebarAction is a quiet navigation row such as "New task" or "Search",
// with an optional trailing hint like a keyboard shortcut.
func SidebarAction(c *ui.Context, icon func() *ui.SVG, label, hint string, selected bool) ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Label(label).FillWidth().MinWidth(0).Height(30).Padding(0, 4).Radius(7).Justify(ui.Start).
		AlignItems(ui.Center).Gap(10).Selected(selected).Background(ui.Color{}).TextColor(secondary(t))
	switch {
	case b.Pressed():
		b.Background(overlayStrong(t))
	case selected || b.Hovered() || b.Focused():
		b.Background(t.Text.Alpha(.07))
	}
	b.Children(func() {
		ui.Row(c).Size(20, 20).Shrink(0).Center().Children(func() {
			components.Icon(c, icon).Size(14, 14).TextColor(secondary(t))
		})
		components.Text(c, label).Grow(1).MinWidth(0).SingleLine().Ellipsis("…").FontSize(13).TextColor(secondary(t))
		if hint != "" {
			components.Text(c, hint).Shrink(0).FontSize(11.5).TextColor(tertiary(t))
		}
	})
	return b
}
