package components

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// PixelGridOptions configures a canvas grid in device-independent pixels.
// Applications can implement zoom by changing Spacing and pan by changing
// OffsetX/OffsetY; the grid itself does not handle input or transform children.
type PixelGridOptions struct {
	// Spacing defaults to 16 for nonpositive or nonfinite values. Positive
	// values are clamped to at least 4 to keep dense grids legible.
	Spacing float32
	// MajorEvery emphasizes every nth line. Zero defaults to 8; negative
	// values disable emphasis. Main and minor lines share the same origin.
	MajorEvery int
	// Offsets move the grid origin relative to the element's top-left corner.
	// Nonfinite offsets are treated as zero.
	OffsetX, OffsetY float32
}

// PixelGrid creates a themed, noninteractive square-line background container.
// Set its size with the returned element's layout methods. Children render above
// the grid and are clipped to its bounds. Background can override the canvas
// color; grid lines use the frame theme and are one DIP wide. Lines coinciding
// with the top or left edge are omitted so the grid has no leading border.
func PixelGrid(c *ui.Context, options PixelGridOptions) *ui.Element {
	t := c.Theme()
	spacing := options.Spacing
	if !finiteGridValue(spacing) || spacing <= 0 {
		spacing = 16
	}
	spacing = max(4, spacing)
	major := options.MajorEvery
	if major == 0 {
		major = 8
	}
	return ui.Box(c).Background(t.Background).Clip().Draw(func(p *ui.Painter, r ui.Rect) {
		p.Clip(r, 0, func() {
			paint := func(step float64, color ui.Color) {
				gridPositions(r.W, step, options.OffsetX, func(x float32) {
					p.Fill(ui.Rect{X: r.X + x, Y: r.Y, W: min(1, r.W-x), H: r.H}, color, 0)
				})
				gridPositions(r.H, step, options.OffsetY, func(y float32) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y + y, W: r.W, H: min(1, r.H-y)}, color, 0)
				})
			}
			paint(float64(spacing), t.Text.Alpha(0.07))
			if major > 0 {
				paint(float64(spacing)*float64(major), t.Text.Alpha(0.14))
			}
		})
	})
}

func finiteGridValue(v float32) bool {
	return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0)
}

// Reduce the origin before iterating so work depends on visible size, not pan.
func gridPositions(length float32, step float64, offset float32, visit func(float32)) {
	if !finiteGridValue(length) || length <= 0 {
		return
	}
	if !finiteGridValue(offset) {
		offset = 0
	}
	start := math.Mod(float64(offset), step)
	if start < 0 {
		start += step
	}
	if start == 0 {
		start = step
	}
	for x := start; x < float64(length); x += step {
		visit(float32(x))
	}
}
