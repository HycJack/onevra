package components

import "github.com/egoist/mygo/ui"

// Tabs creates a row of tabs showing labels, of which *selected is the
// index of the one chosen. A click chooses a tab, as do the arrows, Home
// and End while one has the keyboard focus, which follows the choice.
// The chosen tab is highlighted with the accent color; List's Changed
// reports a new choice.
func Tabs(c *ui.Context, selected *int, labels ...string) ui.TabsParts {
	t := c.Theme()
	tabs := ui.TabsBase(c, selected, len(labels))
	tabs.List.Gap(t.Space(1))
	tabs.List.Children(func() {
		for i, label := range labels {
			tab := dimmed(tabs.Tab(i).Padding(t.Space(2), t.Space(3)).FocusRing(false))
			on := i == *selected
			fg := t.TextMuted
			if on || tab.Hovered() {
				fg = t.Text
			}
			tab.TextColor(fg)
			if on {
				tab.FontWeight(600)
			}
			tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if on {
					in, h := t.Space(1.5), t.Space(0.5)
					p.Fill(ui.Rect{X: r.X + in, Y: r.Y + r.H - h, W: r.W - 2*in, H: h}, t.Accent, h/2)
				}
				if tab.FocusVisible() {
					p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
				}
			})
			tab.Children(func() { ui.Text(c, label).SingleLine() })
		}
	})
	return tabs
}
