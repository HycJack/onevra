package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Brand colors follow the application icon rather than the interface theme.
var (
	brandDeep  = ui.RGB(0x3e, 0x6e, 0x7a)
	brandSky   = ui.RGB(0x51, 0x9d, 0xb3)
	brandMint  = ui.RGB(0x8b, 0xcb, 0xb0)
	brandPaper = ui.RGB(0xf4, 0xf3, 0xef)
)

// welcomeBrandMark redraws the application icon's panels so it stays crisp
// at any size without loading a bitmap.
func welcomeBrandMark(c *ui.Context, size float32) ui.Element {
	gap, inset := size*.07, size*.2
	panel := (size - 2*inset - gap) / 2
	return ui.Row(c).Size(size, size).Shrink(0).Padding(inset).Gap(gap).Radius(size*.24).
		Background(brandPaper).Border(1, ui.RGB(0, 0, 0).Alpha(.06)).
		Shadow(0, size*.06, size*.22, -size*.06, brandDeep.Alpha(.35)).Children(func() {
		ui.Box(c).Width(panel * .72).FillHeight().Radius(panel * .36).Background(brandDeep)
		ui.Column(c).Grow(1).FillHeight().Gap(gap).Children(func() {
			ui.Box(c).FillWidth().Grow(1).Radius(panel * .22).Background(brandSky)
			ui.Box(c).FillWidth().Grow(1).Radius(panel * .22).Background(brandMint)
		})
	})
}

type welcomeGlow struct {
	x, y, radius float32 // fractions of the painted width and height
	color        ui.Color
	alpha        float32
}

// welcomeBackdrop paints a theme-tinted gradient with soft light sources.
// Shadows of a point make the glows, so no image or radial gradient is needed.
func welcomeBackdrop(c *ui.Context, content ui.Element, glows ...welcomeGlow) ui.Element {
	t := c.Theme()
	strength := float32(1)
	if t.Dark {
		strength = .7
	}
	return content.Gradient(t.Background.Mix(t.Accent, .07), t.Background, 150).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.Clip(r, 0, func() {
				for _, g := range glows {
					size := g.radius * max(r.W, r.H)
					point := ui.Rect{X: r.X + g.x*r.W, Y: r.Y + g.y*r.H, W: 1, H: 1}
					p.Shadow(point, .5, 0, 0, size, size*.35, g.color.Alpha(g.alpha*strength))
				}
			})
		})
}

// welcomeEnter fades and lifts an element into place once, staggered by
// order. MyGo completes it at once while the desktop asks for less motion.
func welcomeEnter(c *ui.Context, e ui.Element, order int) ui.Element {
	v := e.Keyframes("enter", []ui.Keyframe{
		{Value: 1}, {At: 560 * time.Millisecond, Value: 0, Ease: ui.EaseOut},
	}, ui.KeyframeOptions{Delay: time.Duration(order) * 70 * time.Millisecond})
	return e.Opacity(1-v).Translate(0, v*c.Theme().Space(3))
}
