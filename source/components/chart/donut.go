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
	indices, values, total := positive(data)
	if !slices.Contains(indices, state.Index) {
		state.Index = -1
	}
	e.Children(func() {
		arcs := animateSlices(c, e, data, indices, base.Pie(values, .025), opts)
		if len(indices) == 0 {
			empty(c, opts)
		} else {
			ring := ui.Box(c).FillWidth().MinWidth(0).Height(plotHeight(opts)).Label(opts.Label).Description(opts.Description)
			if !ring.IsDisabled() {
				shape := make([]float64, 0, 2*len(arcs))
				for _, arc := range arcs {
					shape = append(shape, arc.Start, arc.End)
				}
				shapeSelection(c, ring, state, indices, shape, func(x, y float32) int {
					r := ring.LayoutBounds()
					outer := max(0, min(r.W, r.H)/2-8)
					if j, ok := base.SliceAt(arcs, x-r.W/2, y-r.H/2, outer*.64, outer); ok {
						return j
					}
					return -1
				}, func(i int) string { return describe(opts, sliceSelection(c, data, i)) })
			}
			selected := state.Index
			ring.Draw(func(p *ui.Painter, r ui.Rect) {
				outer := max(0, min(r.W, r.H)/2-8)
				for j, slice := range arcs {
					i := indices[j]
					shade := color(c, i, data[i].Color)
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
		sliceLegend(c, state, data, indices, opts)
		detail(c, opts, sliceSelection(c, data, state.Index))
	})
	return e
}
