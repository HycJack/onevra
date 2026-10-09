package components

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// SpinnerVariant selects the appearance of an indeterminate activity indicator.
type SpinnerVariant int

const (
	// SpinnerRing uses the same track and rotating arc as ProgressRing.
	SpinnerRing SpinnerVariant = iota
	// SpinnerSpokes uses twelve fading spokes, like MyGo's built-in spinner.
	SpinnerSpokes
)

// Spinner creates an indicator of work of unknown length. Name it with Label
// for accessibility. Both variants default to a theme-scaled diameter and stay
// circular and centered in non-square layouts. Reduced Motion freezes activity.
// Unknown variants use SpinnerRing.
func Spinner(c *ui.Context, variant SpinnerVariant) ui.Element {
	if variant != SpinnerSpokes {
		return ProgressRing(c, -1).Shrink(0)
	}
	t := c.Theme()
	reduced := c.Preferences().ReduceMotion
	e := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).
		Role(ui.RoleProgress).Range(0, 1, -1)
	// Spoke geometry and fading follow MyGo's indicator (MIT; see LICENSE.mygo),
	// composed on public painting primitives with reduced-motion support.
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		d := min(r.W, r.H)
		if d <= 0 {
			return
		}
		const spokes, turn = 12, 900
		lead := 0
		if !reduced {
			ms := int(p.Now().UnixMilli() % turn)
			lead = ms * spokes / turn
			p.After(time.Duration((lead+1)*turn/spokes-ms) * time.Millisecond)
		}
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		for i := range spokes {
			a := float64(i)*2*math.Pi/spokes - math.Pi/2
			x, y := float32(math.Cos(a)), float32(math.Sin(a))
			age := (lead - i + spokes) % spokes
			alpha := 1 - float32(age)/spokes*0.85
			var path ui.Path
			path.MoveTo(cx+x*d*0.22, cy+y*d*0.22).
				LineTo(cx+x*d*0.46, cy+y*d*0.46)
			p.StrokePath(&path, d*0.09, t.TextMuted.Alpha(alpha))
		}
	})
	return e
}
