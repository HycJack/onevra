package chat

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// shimmerPeriod is how long the highlight takes to cross the text and rest.
const shimmerPeriod = 1800 * time.Millisecond

// Shimmer draws a single line of text of size in color with a highlight
// sweeping across it, marking ongoing work such as a live activity title.
// Only drawing animates; the text is shaped once per build. Reduced Motion
// shows plain text.
func Shimmer(c *ui.Context, text string, size float32, color, highlight ui.Color) ui.Element {
	if c.Preferences().ReduceMotion {
		return components.Text(c, text).Shrink(0).SingleLine().FontSize(size).TextColor(color)
	}
	font := ui.Font{Size: size}
	glyphs := ui.Shape(text, font)
	var width float32
	if n := len(glyphs); n > 0 {
		width = glyphs[n-1].X + glyphs[n-1].Advance
	}
	m := font.Metrics()
	return ui.Box(c).Label(text).Size(width, m.Ascent+m.Descent).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
		p.After(33 * time.Millisecond)
		// The highlight's center travels from before the first glyph to past
		// the last, so each sweep enters and leaves softly.
		band := max(r.W*.25, size*1.25)
		phase := float32(p.Now().UnixMilli()%shimmerPeriod.Milliseconds()) / float32(shimmerPeriod.Milliseconds())
		center := -band + phase*(r.W+3*band)
		baseline := r.Y + m.Ascent
		for i := range glyphs {
			g := glyphs[i : i+1]
			d := float64((g[0].X + g[0].Advance/2 - center) / band)
			p.Glyphs(g, r.X, baseline, color.Mix(highlight, float32(math.Exp(-d*d*2))*.75))
		}
	})
}
