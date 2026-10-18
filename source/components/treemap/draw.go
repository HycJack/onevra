package treemap

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// input records primary clicks for the next build. It takes nothing, so
// pressing still focuses the view.
func (v *view) input(ev ui.InputEvent) bool {
	switch ev.Kind {
	case ui.InputPointerDown:
		if ev.Button == 0 {
			v.pressed, v.pressX, v.pressY = true, ev.X, ev.Y
		}
	case ui.InputPointerUp:
		if v.pressed && ev.Button == 0 {
			v.pressed = false
			if math.Hypot(float64(ev.X-v.pressX), float64(ev.Y-v.pressY)) < tapSlop {
				v.clicks = append(v.clicks, click{v.pressX, v.pressY})
			}
		}
	}
	return false
}

// hoverFade is how long tile highlights take to cross-fade, at most the
// view's transition duration.
const hoverFade = 160 * time.Millisecond

// draw paints the tiles at the canvas r, through a drill or data transition
// while one runs.
func (v *view) draw(p *ui.Painter, r ui.Rect, t *ui.Theme, active int, focused bool) {
	full := ui.Rect{W: v.width, H: v.height}
	none := func(string) float32 { return 0 }
	if v.zooming {
		e := float32(p.Now().Sub(v.zoomStart)) / float32(v.duration)
		if e < 1 {
			e = ui.EaseInOut(max(e, 0))
			p.AnimationFrame()
			// One affine zoom carries both pictures: the old one maps the
			// drilled tile onto the growing box, the new one fills it.
			if v.zoomIn {
				box := lerp(v.zoomRect, full, e)
				v.paint(p, r, t, v.old, mapping(r, v.zoomRect, box), 1, 1, none)
				p.Clip(offset(box, r), 0, func() { v.paint(p, r, t, v.tiles, mapping(r, full, box), 1, 1, none) })
			} else {
				box := lerp(full, v.zoomRect, e)
				v.paint(p, r, t, v.tiles, mapping(r, v.zoomRect, box), 1, 1, none)
				v.paint(p, r, t, v.old, mapping(r, full, box), 1-e, 1, none)
			}
			return
		}
		v.zooming, v.old = false, nil
	}
	e := v.morphProgress(p.Now())
	if v.morphing {
		if e < 1 {
			p.AnimationFrame()
		} else {
			v.morphing = false
		}
	}
	// The highlight moves from the tile the pointer left to the one it is on.
	h := float32(1)
	if v.duration > 0 {
		if h = float32(p.Now().Sub(v.hoverStart)) / float32(min(hoverFade, v.duration)); h < 1 {
			h = ui.EaseInOut(max(h, 0))
			p.AnimationFrame()
		} else {
			h = 1
		}
	}
	lift := func(id string) float32 {
		switch id {
		case "":
			return 0
		case v.hoverID:
			return h
		case v.hoverOut:
			return 1 - h
		}
		return 0
	}
	v.paint(p, r, t, v.tiles, ui.Transform{}, 1, e, lift)
	outline := func(i int, width float32, c ui.Color) {
		tl := v.tiles[i]
		b := offset(inset(lerp(tl.from, tl.rect, e), -.5, -.5), r)
		radius := v.radius(b, tl.depth)
		if c == (ui.Color{}) {
			p.FocusRing(b, [4]float32{radius, radius, radius, radius})
			return
		}
		p.Stroke(b, c, radius, width)
	}
	// The shown level's tile under the pointer is what a click drills into.
	for _, id := range [...]string{v.hoverOut, v.hoverID} {
		if i := v.index(id); i >= 0 {
			if top := v.topOf(i); top != active {
				amount := lift(id)
				outline(top, 1.5, t.Background.Mix(v.tiles[top].color, .7).Alpha(amount))
			}
		}
	}
	if active >= 0 && active < len(v.tiles) {
		if focused {
			outline(active, 0, ui.Color{})
		} else {
			outline(active, 2, v.tiles[active].color.Mix(t.Text, .2))
		}
	}
}

// paint draws tiles through transform at opacity, e of the way through their
// data transition, deepening each by lift of the hover highlight.
func (v *view) paint(p *ui.Painter, r ui.Rect, t *ui.Theme, tiles []tile, transform ui.Transform, opacity, e float32, lift func(string) float32) {
	m := v.metrics
	text, muted := t.Text.Alpha(opacity), t.TextMuted.Alpha(opacity)
	p.Transform(transform, 0, 0, func() {
		for _, tl := range tiles {
			moving := e < 1 && tl.from != tl.rect
			b := offset(lerp(tl.from, tl.rect, e), r)
			if b.W < .5 || b.H < .5 {
				continue
			}
			amount := tl.fromAmount + (tl.amount-tl.fromAmount)*e + .08*lift(tl.id)
			radius := v.radius(b, tl.depth)
			// A soft diagonal gradient, lighter at the top left, with a faint
			// edge in the same hue.
			p.FillGradient(b, ui.LinearGradient{
				From:  t.Background.Mix(tl.color, amount*.6).Alpha(opacity),
				To:    t.Background.Mix(tl.color, min(amount*1.3, .7)).Alpha(opacity),
				Angle: 135,
			}, radius)
			p.Stroke(b, t.Background.Mix(tl.color, min(amount+.2, .8)).Alpha(.5*opacity), radius, 1)
			labels := func() {
				x, y := b.X+m.padX, b.Y+m.padY
				if tl.label != "" {
					p.RichText(x, y, 0, ui.Span{Text: tl.label, Size: m.small, Weight: labelWeight, Color: text})
				}
				if tl.value != "" {
					p.RichText(b.X+b.W-m.padX-tl.valueWidth, y, 0, ui.Span{Text: tl.value, Size: m.small, Weight: labelWeight, Color: muted})
				}
				if tl.big != "" {
					y += m.line + m.padY
					p.RichText(x, y, 0, ui.Span{Text: tl.big, Size: m.big, Weight: bigWeight, Color: text})
					if tl.caption != "" {
						p.RichText(x, y+m.bigLine, 0, ui.Span{Text: tl.caption, Size: m.small, Weight: labelWeight, Color: muted})
					}
				}
			}
			if tl.label == "" && tl.big == "" {
				continue
			}
			if moving {
				// Labels are fitted to where the tile settles; keep them inside
				// it on the way.
				p.Clip(b, radius, labels)
			} else {
				labels()
			}
		}
	})
}

// radius rounds a tile, less with depth so nested corners stay concentric
// with their parents', and never more than a small tile can take.
func (v *view) radius(b ui.Rect, depth int) float32 {
	m := v.metrics
	r := max(m.radius-m.radiusStep*float32(depth), m.radiusStep)
	return max(0, min(r, b.W/3, b.H/3))
}

// index returns the tile with id, or -1.
func (v *view) index(id string) int {
	if id == "" {
		return -1
	}
	for i := range v.tiles {
		if v.tiles[i].id == id {
			return i
		}
	}
	return -1
}

// mapping is the transform that draws the canvas rectangle from onto to, for
// a canvas at r in window coordinates.
func mapping(r, from, to ui.Rect) ui.Transform {
	if from.W <= 0 || from.H <= 0 {
		return ui.Transform{}
	}
	return ui.Translation(-(r.X + from.X), -(r.Y + from.Y)).
		Then(ui.Scaling(to.W/from.W, to.H/from.H)).
		Then(ui.Translation(r.X+to.X, r.Y+to.Y))
}

// offset moves a canvas rectangle into window coordinates.
func offset(b, r ui.Rect) ui.Rect { b.X, b.Y = b.X+r.X, b.Y+r.Y; return b }
