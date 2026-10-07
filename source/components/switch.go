package components

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Switch creates a switch toggling *on, with a label. The label is also
// the accessible name.
func Switch(c *ui.Context, on *bool, label string) *ui.Element {
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(2))
	row.Children(func() {
		sw := ui.SwitchBase(c, on).Size(t.Space(9), t.Space(5)).Radius(t.Space(2.5))
		if label != "" {
			sw.Label(label)
		}
		pos := sw.Animate("knob", boolToFloat(*on), 140*time.Millisecond)
		off := t.Border.Mix(t.Text, 0.15)
		sw.Background(off.Mix(t.Accent, pos))
		sw.Draw(func(p *ui.Painter, r ui.Rect) {
			in := t.Space(0.5)
			d := r.H - 2*in
			knob := ui.Rect{X: r.X + in + pos*(r.W-r.H), Y: r.Y + in, W: d, H: d}
			p.Shadow(knob, d/2, 0, 1, 3, 0, ui.RGBA(0, 0, 0, 0.25))
			p.Fill(knob, ui.RGB(255, 255, 255), d/2)
		})
		if label != "" {
			Text(c, label)
		}
	})
	return dimmed(row)
}

func boolToFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}
