package gallery

import (
	"github.com/egoist/mygo/plugins/game"
	"github.com/egoist/mygo/plugins/game/gamepad"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// The world is tiles of 8 world units; the view shows a fixed logical
// screen, scaled by whole multiples so the pixel art stays crisp.
const (
	gameTiles           = 40
	gameWorld           = gameTiles * spriteSize
	gameScreenW         = 160
	gameScreenH         = 96
	gameSpeed   float32 = 1.25
)

// gameDemo is the page's game.Game. Its world survives navigation; the
// loop runs only while game.View is built. Gamepads are read only when
// Gallery opened a native window, so headless views load no frameworks.
type gameDemo struct {
	gamepads bool
	x, y     float32
	px, py   float32
	steps    int
	moving   bool
	left     bool
	sheet    *ui.Bitmap
	screen   *ui.Canvas
}

func (d *gameDemo) init() {
	if d.sheet == nil {
		d.sheet = ui.NewBitmap(spriteSheet()).Pixelated()
		d.screen = ui.NewCanvas(gameScreenW, gameScreenH)
		d.screen.Bitmap().Pixelated()
		d.x, d.y = gameWorld/2, gameWorld/2
		d.px, d.py = d.x, d.y
	}
}

func (d *gameDemo) Update(in *game.Input) {
	d.px, d.py = d.x, d.y
	var dx, dy float32
	held := func(keys ...ui.Key) bool {
		for _, k := range keys {
			if in.Down(k) {
				return true
			}
		}
		return false
	}
	if held(ui.KeyA, ui.KeyLeft) {
		dx--
	}
	if held(ui.KeyD, ui.KeyRight) {
		dx++
	}
	if held(ui.KeyW, ui.KeyUp) {
		dy--
	}
	if held(ui.KeyS, ui.KeyDown) {
		dy++
	}
	if d.gamepads {
		if pads := gamepad.Pads(); len(pads) > 0 {
			sx, sy := pads[0].State().LeftStick()
			dx, dy = dx+sx, dy+sy
		}
	}
	d.moving = dx != 0 || dy != 0
	if dx != 0 {
		d.left = dx < 0
	}
	d.x = min(gameWorld-spriteSize/2, max(spriteSize/2, d.x+dx*gameSpeed))
	d.y = min(gameWorld-spriteSize/2, max(spriteSize/2, d.y+dy*gameSpeed))
	d.steps++
}

func (d *gameDemo) Draw(p *ui.Painter, r ui.Rect, alpha float64) {
	a := float32(alpha)
	x, y := d.px+(d.x-d.px)*a, d.py+(d.y-d.py)*a
	// Keep the camera inside the world so its edges never show.
	cam := game.Camera{
		X: min(gameWorld-gameScreenW/2, max(gameScreenW/2, x)),
		Y: min(gameWorld-gameScreenH/2, max(gameScreenH/2, y)),
	}
	screen := ui.Rect{W: gameScreenW, H: gameScreenH}
	frame := 0
	if d.moving {
		frame = d.steps / 8 % 3
	} else if d.steps/30%8 == 0 {
		frame = 3
	}
	d.screen.Draw(p, func(cp *ui.Painter) {
		cam.Draw(cp, screen, func() {
			v := cam.Visible(screen)
			for ty := max(0, int(v.Y)/spriteSize); ty < min(gameTiles, int(v.Y+v.H)/spriteSize+1); ty++ {
				for tx := max(0, int(v.X)/spriteSize); tx < min(gameTiles, int(v.X+v.W)/spriteSize+1); tx++ {
					cp.Fill(ui.Rect{X: float32(tx * spriteSize), Y: float32(ty * spriteSize), W: spriteSize, H: spriteSize}, gameTile(tx, ty), 0)
				}
			}
			cp.ImageRegion(d.sheet, spriteFrame(frame), ui.Rect{X: x - spriteSize/2, Y: y - spriteSize/2, W: spriteSize, H: spriteSize},
				ui.ImageOptions{FlipX: d.left})
		})
	})
	p.Fill(r, ui.RGB(28, 40, 32), 0) // letterbox
	dst, _ := game.IntegerFit(p, gameScreenW, gameScreenH, r)
	p.Image(d.screen.Bitmap(), dst, ui.FillBox)
}

// gameTile is a checkered meadow with scattered flowers, fixed per tile.
func gameTile(x, y int) ui.Color {
	if (x*7+y*13)%23 == 0 {
		return ui.RGB(240, 200, 90)
	}
	if (x+y)%2 == 0 {
		return ui.RGB(118, 178, 92)
	}
	return ui.RGB(104, 164, 82)
}

func exampleGame(c *ui.Context, s *demoState) {
	t := c.Theme()
	d := &s.game
	d.init()
	components.Text(c, tr(c, "gallery.game_hint")).TextColor(t.TextMuted).Wrap()
	game.View(c, d).Height(300).FillWidth().Radius(12).Clip().Label(tr(c, "gallery.game_scene"))
}
