package components

import (
	"time"

	"github.com/egoist/mygo/ui"
)

const tabsIndicatorDuration = 180 * time.Millisecond

// Tabs creates a row of tabs showing labels, of which *selected is the
// index of the one chosen. A click chooses a tab, as do the arrows, Home
// and End while one has the keyboard focus, which follows the choice.
// A shared accent underline animates to the chosen tab without changing
// its font weight. Reduced Motion skips the transition. List's Changed
// reports a new choice.
func Tabs(c *ui.Context, selected *int, labels ...string) ui.TabsParts {
	t := c.Theme()
	tabs := ui.TabsBase(c, selected, len(labels))
	tabs.List.Gap(t.Space(1))
	var active ui.Element
	tabs.List.Children(func() {
		for i, label := range labels {
			tab := dimmed(tabs.Tab(i).Padding(t.Space(2), t.Space(3)).FocusRing(false))
			on := i == *selected
			fg := t.Text
			if !on {
				fg = fg.Alpha(0.6)
			}
			tab.TextColor(fg)
			if on {
				active = tab
			}
			tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if tab.FocusVisible() {
					p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
				}
			})
			tab.Children(func() { Text(c, label).SingleLine() })
		}
	})
	if len(labels) > 0 {
		tabs.List.DrawOver(func(p *ui.Painter, r ui.Rect) {
			// Read current layout during painting, and animate list-relative
			// geometry so moving the entire control does not move its underline.
			bounds := active.LayoutBounds()
			in, h := t.Space(1.5), t.Space(0.5)
			x := tabs.List.Animate("indicator-x", bounds.X-r.X+in, tabsIndicatorDuration)
			w := tabs.List.Animate("indicator-width", max(0, bounds.W-2*in), tabsIndicatorDuration)
			if c.Preferences().ReduceMotion {
				x, w = bounds.X-r.X+in, max(0, bounds.W-2*in)
			}
			p.Fill(ui.Rect{X: r.X + x, Y: bounds.Y + bounds.H - h, W: w, H: h}, t.Accent, h/2)
		})
	}
	return tabs
}
