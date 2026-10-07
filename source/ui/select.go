package ui

import "github.com/egoist/mygo/ui"

// Select creates a drop-down choosing *selected from options the caller
// builds in Popup with Item. The trigger has the oneo button look; build
// the popup and the items with the theme, highlighting the pointed option
// with the accent color:
//
//	sel := oneoui.Select(c, &choice)
//	sel.Trigger.Children(func() { ui.Text(c, string(choice)).SingleLine() })
//	sel.Popup(func(panel *ui.Element) {
//		t := c.Theme()
//		panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
//		for _, opt := range options {
//			item := sel.Item(opt).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
//			if item.Highlighted() {
//				item.Background(t.Accent).TextColor(t.AccentText)
//			}
//			item.Children(func() { ui.Text(c, string(opt)).SingleLine() })
//		}
//	})
//
// The trigger's Changed reports a new choice.
func Select[T comparable](c *ui.Context, selected *T) *ui.SelectParts[T] {
	t := c.Theme()
	sel := ui.SelectBase(c, selected)
	b := sel.Trigger
	b.Padding(t.Space(1.5), t.Space(6), t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
	b.Justify(ui.Start)
	bg := t.Surface
	if !b.IsDisabled() {
		switch {
		case b.Pressed():
			bg = t.SurfacePressed
		case b.Hovered():
			bg = t.SurfaceHover
		}
	}
	b.Background(bg).TextColor(t.Text).Border(1, t.Border)
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		s := t.Space(2)
		cx := r.X + r.W - t.Space(3)
		cy := r.Y + r.H/2
		var path ui.Path
		path.MoveTo(cx-s*0.35, cy-s*0.15).LineTo(cx, cy+s*0.2).LineTo(cx+s*0.35, cy-s*0.15)
		p.StrokePath(&path, 1.5, t.TextMuted)
	})
	dimmed(b)
	return sel
}
