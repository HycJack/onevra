package chart

import (
	"slices"

	"github.com/egoist/mygo/ui"
	base "github.com/egoist/mygo/ui/chart"
	"github.com/oneo-me/onevra/source/components"
)

// Donut draws proportional slices. Hover, arrow keys, Home/End and legend
// buttons select a slice; Escape clears it. State.Hidden is unused. Invalid
// and non-positive slices are omitted while State.Index refers to input order.
func Donut(c *ui.Context, state *State, data []Slice, opts Options) ui.Element {
	e, state := frame(c, state)
	t := c.Theme()
	colors := make([]ui.Color, len(data))
	indices, values := []int{}, []float64{}
	total, largest := 0.0, 0.0
	for i, s := range data {
		colors[i] = color(c, i, s.Color)
		if finite(s.Value) && s.Value > 0 {
			indices, values = append(indices, i), append(values, s.Value)
			total, largest = total+s.Value, max(largest, s.Value)
		}
	}
	for i := range values {
		values[i] /= largest // Keep the pie's sum finite even for very large values.
	}
	if !slices.Contains(indices, state.Index) {
		state.Index = -1
	}
	selection := func(i int) Selection {
		if i < 0 {
			return Selection{Index: -1}
		}
		s := data[i]
		return Selection{Index: i, Label: s.Label, Values: []Value{{s.ID, s.Label, s.Value, color(c, i, s.Color)}}}
	}
	e.Children(func() {
		arcs := animateSlices(c, e, data, indices, base.Pie(values, .025), opts)
		if len(indices) == 0 {
			empty(c, opts)
		} else {
			ring := ui.Box(c).FillWidth().MinWidth(0).Height(plotHeight(opts)).Label(opts.Label).Description(opts.Description)
			if !ring.IsDisabled() {
				radialSelection(c, ring, state, indices, arcs, func(i int) string { return describe(opts, selection(i)) })
			}
			selected := state.Index
			ring.Draw(func(p *ui.Painter, r ui.Rect) {
				outer := max(0, min(r.W, r.H)/2-8)
				for j, slice := range arcs {
					i := indices[j]
					shade := colors[i]
					edge := outer - 4
					if i == selected {
						edge = outer
					}
					p.FillPathGradient(base.Sector(r.X+r.W/2, r.Y+r.H/2, outer*.64, edge, slice), ui.LinearGradient{
						From: shade, To: shade.Mix(t.Background, .6), Angle: 180,
					})
				}
			})
			ring.Children(func() {
				ui.Column(c).Absolute().LeftPercent(25).TopPercent(35).WidthPercent(50).HeightPercent(30).
					AlignItems(ui.Center).Justify(ui.Center).Gap(4).Children(func() {
					v := total
					if selected >= 0 {
						v = data[selected].Value
					}
					text := "—"
					if finite(v) {
						text = format(opts, v)
					}
					components.Text(c, text).FontSize(t.FontSize * 1.7).Bold().NoWrap().MaxWidthPercent(100).Ellipsis("…")
				})
			})
		}
		if !opts.HideLegend {
			ui.Row(c).FillWidth().MinWidth(0).Wrap().Justify(ui.Center).Gap(c.Theme().Space(1)).Children(func() {
				for _, i := range indices {
					s := data[i]
					ui.Box(c).Key(s.ID).MaxWidthPercent(100).Children(func() {
						b := legendButton(c, s.Label, state.Index == i)
						b.Children(func() {
							swatch(c, color(c, i, s.Color))
							components.Text(c, s.Label).FontSize(c.Theme().FontSize - 1)
						})
						if b.Clicked() {
							state.Index = i
							c.Invalidate()
							c.Announce(describe(opts, selection(i)))
						}
					})
				}
			})
		}
		detail(c, opts, selection(state.Index))
	})
	return e
}

type radialHover struct {
	x, y    float32
	over    bool
	pointer bool
	arcs    []base.Slice
	indices []int
}

func radialSelection(c *ui.Context, ring ui.Element, state *State, indices []int, arcs []base.Slice, describe func(int) string) {
	ring.Focusable()
	position := slices.Index(indices, state.Index)
	keyboard := true
	switch {
	case ring.Shortcut(0, ui.KeyRight):
		state.Index = indices[min(position+1, len(indices)-1)]
	case ring.Shortcut(0, ui.KeyLeft):
		if position < 0 {
			position = len(indices)
		}
		state.Index = indices[max(0, position-1)]
	case ring.Shortcut(0, ui.KeyHome):
		state.Index = indices[0]
	case ring.Shortcut(0, ui.KeyEnd):
		state.Index = indices[len(indices)-1]
	case ring.Shortcut(0, ui.KeyEscape):
		state.Index = -1
	default:
		keyboard = false
	}
	hover := ui.Local(ring, "onevra.chart.radial", func() radialHover { return radialHover{} })
	changed := !slices.Equal(hover.arcs, arcs) || !slices.Equal(hover.indices, indices)
	hover.arcs, hover.indices = arcs, indices
	x, y, over := ring.PointerPosition()
	if keyboard || !over {
		hover.pointer = false
	} else if !hover.over || x != hover.x || y != hover.y {
		hover.pointer = true
	}
	if over && !keyboard && (!hover.over || x != hover.x || y != hover.y || (changed && hover.pointer)) {
		r := ring.LayoutBounds()
		outer := max(0, min(r.W, r.H)/2-8)
		state.Index = -1
		if j, ok := base.SliceAt(arcs, x-r.W/2, y-r.H/2, outer*.64, outer); ok {
			state.Index = indices[j]
		}
	} else if !over && hover.over && !keyboard {
		state.Index = -1
	}
	hover.x, hover.y, hover.over = x, y, over
	if state.Index >= 0 {
		value := describe(state.Index)
		ring.Value(value)
		if keyboard {
			c.Announce(value)
		}
	}
}
