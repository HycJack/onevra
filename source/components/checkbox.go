package components

import "github.com/egoist/mygo/ui"

// Checkbox creates a check box toggling *checked, with a label. The label
// is also the accessible name.
func Checkbox(c *ui.Context, checked *bool, label string) *ui.Element {
	t := c.Theme()
	row := ui.CheckboxBase(c, checked).Gap(t.Space(2)).FocusRing(false)
	on := *checked
	row.Children(func() {
		box := ui.Box(c).Size(t.Space(4), t.Space(4)).Radius(t.Space(1)).Shrink(0)
		border := t.Border.Mix(t.Text, 0.25)
		if !on && row.Hovered() {
			border = t.Accent
		}
		if on {
			box.Background(t.Accent)
		} else {
			box.Background(t.Surface).Border(1, border)
		}
		box.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if on {
				var path ui.Path
				path.MoveTo(r.X+r.W*0.22, r.Y+r.H*0.52).
					LineTo(r.X+r.W*0.42, r.Y+r.H*0.71).
					LineTo(r.X+r.W*0.78, r.Y+r.H*0.31)
				p.StrokePath(&path, t.Space(0.5), t.AccentText)
			}
			if row.FocusVisible() {
				p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
			}
		})
		if label != "" {
			row.Label(label)
			ui.Text(c, label)
		}
	})
	return dimmed(row)
}
