package gallery

import (
	"image"
	"image/color"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

const spriteSize, spriteFrames = 8, 4

// spriteRows are four 8×8 frames of a bouncing slime: rest, squash,
// stretch and blink. Pixel art is fixture content, not themed.
var spriteRows = [spriteFrames][spriteSize]string{
	{"........", "..oooo..", ".oggggo.", "ogwggwgo", "ogkggkgo", "oggggggo", "oggggggo", ".oooooo."},
	{"........", "........", ".oooooo.", "ogwggwgo", "ogkggkgo", "oggggggo", "oggggggo", "oooooooo"},
	{"..oooo..", ".oggggo.", "ogwggwgo", "ogkggkgo", "oggggggo", ".oggggo.", ".oggggo.", "..oooo.."},
	{"........", "..oooo..", ".oggggo.", "oggggggo", "okkggkko", "oggggggo", "oggggggo", ".oooooo."},
}

var spritePalette = map[byte]color.NRGBA{
	'o': {20, 50, 40, 255}, 'g': {74, 200, 120, 255}, 'w': {255, 255, 255, 255}, 'k': {20, 30, 40, 255},
}

// spriteSheet lays the frames out in a row, as a game's sheet would be.
func spriteSheet() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, spriteSize*spriteFrames, spriteSize))
	for f, rows := range spriteRows {
		for y, row := range rows {
			for x := range spriteSize {
				img.SetNRGBA(f*spriteSize+x, y, spritePalette[row[x]])
			}
		}
	}
	return img
}

func spriteFrame(i int) ui.Rect {
	return ui.Rect{X: float32(i%spriteFrames) * spriteSize, W: spriteSize, H: spriteSize}
}

// canvasDemo keeps one smooth and one pixelated bitmap of the same sheet,
// and an offscreen canvas redrawn only when its frame or mirroring change.
type canvasDemo struct {
	smooth, pixelated *ui.Bitmap
	layer             *ui.Canvas
	drawn             int
}

func (d *canvasDemo) bitmaps() {
	if d.smooth == nil {
		sheet := spriteSheet()
		d.smooth, d.pixelated = ui.NewBitmap(sheet), ui.NewBitmap(sheet).Pixelated()
		d.layer, d.drawn = ui.NewCanvas(24, 24), -1
		d.layer.Bitmap().Pixelated()
	}
}

const canvasTile = 120

func inset(r ui.Rect, by float32) ui.Rect {
	return ui.Rect{X: r.X + by, Y: r.Y + by, W: r.W - 2*by, H: r.H - 2*by}
}

func exampleCanvas(c *ui.Context, s *demoState) {
	t := c.Theme()
	d := &s.canvas
	d.bitmaps()
	frame := s.index % spriteFrames
	options := ui.ImageOptions{FlipX: s.on}
	tile := func(key, label string, draw func(p *ui.Painter, r ui.Rect)) {
		ui.Column(c).Key(key).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(canvasTile, canvasTile).Radius(8).Clip().Background(t.Surface).Label(label).Draw(draw)
			components.Text(c, label).TextColor(t.TextMuted)
		})
	}
	ui.Row(c).Gap(t.Space(4)).Wrap().Children(func() {
		tile("smooth", tr(c, "gallery.canvas_smooth"), func(p *ui.Painter, r ui.Rect) {
			p.ImageRegion(d.smooth, spriteFrame(frame), inset(r, 12), options)
		})
		tile("pixelated", tr(c, "gallery.canvas_pixelated"), func(p *ui.Painter, r ui.Rect) {
			p.ImageRegion(d.pixelated, spriteFrame(frame), inset(r, 12), options)
		})
		tile("layer", tr(c, "gallery.canvas_layer"), func(p *ui.Painter, r ui.Rect) {
			// Redraw the offscreen layer only when its content changes.
			key := frame * 2
			if s.on {
				key++
			}
			if d.drawn != key {
				d.drawn = key
				d.layer.Draw(p, func(cp *ui.Painter) {
					for y := range 3 {
						for x := range 3 {
							shade := ui.RGB(232, 200, 140)
							if (x+y)%2 == 1 {
								shade = ui.RGB(210, 170, 110)
							}
							cp.Fill(ui.Rect{X: float32(x * 8), Y: float32(y * 8), W: 8, H: 8}, shade, 0)
						}
					}
					cp.ImageRegion(d.pixelated, spriteFrame(frame), ui.Rect{X: 8, Y: 8, W: 8, H: 8}, options)
				})
			}
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			p.Transform(ui.Rotation(float32(s.value)), cx, cy, func() {
				p.Image(d.layer.Bitmap(), ui.Rect{X: cx - 36, Y: cy - 36, W: 72, H: 72}, ui.FillBox)
			})
		})
	})
	ui.Row(c).Gap(t.Space(3)).Wrap().AlignItems(ui.Center).Children(func() {
		if components.Button(c, tr(c, "gallery.canvas_next_frame"), components.ButtonSecondary).Clicked() {
			s.index = (s.index + 1) % spriteFrames
		}
		components.Checkbox(c, &s.on, tr(c, "gallery.canvas_mirror"))
	})
	ui.Row(c).FillWidth().MaxWidth(t.Space(100)).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		name := tr(c, "gallery.canvas_rotation")
		components.Text(c, name).Shrink(0).SingleLine()
		components.Slider(c, &s.value, 0, 360).Label(name).Grow(1).MinWidth(0)
	})
}
