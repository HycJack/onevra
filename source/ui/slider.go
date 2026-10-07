package ui

import "github.com/egoist/mygo/ui"

// Slider creates a slider setting *value between lo and hi. Dragging across
// it, or the arrows, Home and End while it has the focus, change the value;
// Changed reports a new one.
func Slider(c *ui.Context, value *float64, lo, hi float64) *ui.Element {
	t := c.Theme()
	thumb := t.Space(4)
	s := ui.SliderBase(c, value, lo, hi).Height(t.Space(5)).MinWidth(t.Space(20)).PaddingX(thumb / 2).FocusRing(false)
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		frac := float32(0)
		if hi > lo {
			frac = float32((*value - lo) / (hi - lo))
		}
		frac = max(0, min(frac, 1))
		h := t.Space(1)
		cy := r.Y + r.H/2
		tx, tw := r.X+thumb/2, r.W-thumb
		p.Fill(ui.Rect{X: tx, Y: cy - h/2, W: tw, H: h}, t.Border, h/2)
		p.Fill(ui.Rect{X: tx, Y: cy - h/2, W: tw * frac, H: h}, t.Accent, h/2)
		cx := tx + tw*frac
		knob := ui.Rect{X: cx - thumb/2, Y: cy - thumb/2, W: thumb, H: thumb}
		p.Shadow(knob, thumb/2, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.2))
		p.Fill(knob, ui.RGB(255, 255, 255), thumb/2)
	})
	s.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if s.FocusVisible() {
			p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
		}
	})
	return dimmed(s)
}
