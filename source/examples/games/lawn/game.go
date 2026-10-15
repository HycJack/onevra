package main

import (
	"image/color"
	"math/rand/v2"
	"slices"
	"strconv"
	"time"

	"github.com/egoist/mygo/plugins/game"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/examples/games/lawn/locales"
)

// The lawn is a logical screen of whole pixels: a seed bar above five rows
// of nine tiles, with the patio to their left and the street to their
// right, shown pixel-perfect by game.IntegerFit.
const (
	lawnTile          = 16
	lawnRows          = 5
	lawnCols          = 9
	lawnLeft          = 24
	lawnTop           = 24
	lawnRight         = lawnLeft + lawnCols*lawnTile
	lawnW             = 192
	lawnH             = lawnTop + lawnRows*lawnTile
	lawnSecond        = game.Rate
	lawnZombies       = 100 // a level's total, its last lawnFinalWave together
	lawnFinalWave     = 5
	lawnZombieHealth  = 180
	lawnZombieSpeed   = 0.1
	lawnPeaDamage     = 20
	lawnPeaSpeed      = 2
	lawnShotInterval  = 84
	lawnBite          = 2
	lawnSunValue      = 25
	lawnSunLife       = 10 * lawnSecond
	lawnSkySunEvery   = 9 * lawnSecond
	lawnFlowerSunEach = 20 * lawnSecond
)

// Frames of the lawn sheet; the first three are also the seeds.
const (
	lawnSunflower = iota
	lawnPeashooter
	lawnWallNut
	lawnZombieStep
	lawnZombieStride
	lawnSun
	lawnMower
)

// lawnSprites are 16×16 pixel art; zombies face left, plants right.
var lawnSprites = [][]string{
	lawnSunflower: {
		"....y.yyyy.y....", "..yyyyyyyyyyyy..", ".yyyyooooooyyyy.", ".yyyobbbbbboyyy.",
		"yyyyobkbbkboyyyy", "yyyyobbbbbboyyyy", "yyyyobkbbkboyyyy", ".yyyobbkkbboyyy.",
		".yyyyooooooyyyy.", "..yyyyyyyyyyyy..", "....y.yGGy.y....", ".......GG.......",
		"..ggg..GG..ggg..", ".gggggGGGGggggg.", "..gggGGGGGGggg..", "....GGGGGGGG....",
	},
	lawnPeashooter: {
		"....oooooo......", "...oggggggo.....", "..oggggggggoooo.", "..oggggwkgggggGo",
		"..oggggggggggGko", "..oggggggggggGko", "..ogggggggggooo.", "...oggggggggo...",
		"....oooooooo....", ".......GG.......", "......GGG.......", ".......GG.......",
		"..gggg.GG.gggg..", ".ggggggGGgggggg.", "..ggggGGGGgggg..", "....GGGGGGGG....",
	},
	lawnWallNut: {
		".....oooooo.....", "...oonnnnnnoo...", "..onnnnnnnnnno..", ".onnnNnnnnnnnno.",
		".onnnnnnnnnnnno.", ".onnwwnnnnwwnno.", ".onnwknnnnwknno.", ".onnnnnnnnnnnno.",
		".onnnnnnnnnnnno.", ".onnnnNkkNnnnno.", ".onnnnnnnnnnnno.", ".onnNnnnnnnnnno.",
		".onnnnnnnnnnnNo.", "..onnnnnnnnnno..", "...oonnnnnnoo...", ".....oooooo.....",
	},
	lawnZombieStep: {
		".....ssssss.....", "....sssssssS....", "....swksswksS...", "....ssssssssS...",
		"....sskkkkssS...", ".....ssssssS....", "......ssss......", ".....cwwrwcc....",
		"ssscccccrccccc..", "....ccccrcccc...", "....ccccrcccC...", "....ppppppppp...",
		"....pppp.pppp...", "...ppp....ppp...", "...pp......pp...", "..kkk.....kkk...",
	},
	lawnZombieStride: {
		".....ssssss.....", "....sssssssS....", "....swksswksS...", "....ssssssssS...",
		"....sskkkkssS...", ".....ssssssS....", "......ssss......", ".....cwwrwcc....",
		"ssscccccrccccc..", "....ccccrcccc...", "....ccccrcccC...", "....ppppppppp...",
		".....pppppp.....", ".....pp..pp.....", ".....pp..pp.....", "....kkk.kkk.....",
	},
	lawnSun: {
		"................", "................", "......Y..Y......", "..Y...YyyY...Y..",
		"...YyyyyyyyyY...", "...yyyyyyyyyy...", ".YYyyyyyyyyyyYY.", "..yyyyyyyyyyyy..",
		"..yyyyyyyyyyyy..", ".YYyyyyyyyyyyYY.", "...yyyyyyyyyy...", "...YyyyyyyyyY...",
		"..Y...YyyY...Y..", "......Y..Y......", "................", "................",
	},
	lawnMower: {
		"................", "................", "................", "................",
		"................", "................", "..MM............", "...MM...........",
		"....MMmmmmmm....", "....mmmmmmmmmm..", "...mmmmmmmmmmmm.", "...mmmmmmmmmmmm.",
		"...MMMMMMMMMMMM.", "...kk......kk...", "..kkkk....kkkk..", "...kk......kk...",
	},
}

var lawnPalette = map[byte]color.NRGBA{
	'o': {30, 70, 30, 255}, 'g': {96, 196, 72, 255}, 'G': {52, 140, 52, 255},
	'y': {252, 214, 64, 255}, 'Y': {238, 152, 40, 255}, 'b': {124, 74, 32, 255},
	'n': {204, 144, 72, 255}, 'N': {150, 96, 46, 255}, 'w': {255, 255, 255, 255}, 'k': {24, 24, 30, 255},
	's': {152, 184, 142, 255}, 'S': {112, 142, 106, 255}, 'c': {112, 92, 72, 255}, 'C': {82, 66, 52, 255},
	'r': {190, 40, 40, 255}, 'p': {72, 82, 124, 255}, 'm': {204, 52, 52, 255}, 'M': {124, 124, 134, 255},
}

type lawnSeed struct{ cost, recharge, health int }

var lawnSeeds = []lawnSeed{
	lawnSunflower:  {50, 5 * lawnSecond, 300},
	lawnPeashooter: {100, 5 * lawnSecond, 300},
	lawnWallNut:    {50, 20 * lawnSecond, 2400},
}

// lawnDigits are a 3×5 pixel font, row by row.
var lawnDigits = [10]string{
	"111101101101111", "010110010010111", "111001111100111", "111001111001111", "101101111001001",
	"111100111001111", "111100111101111", "111001010010010", "111101111101111", "111101111001111",
}

type lawnState int

const (
	lawnPlaying lawnState = iota
	lawnPaused
	lawnLost
	lawnWon
)

type lawnPlant struct{ seed, health, timer int }

type lawnZombie struct {
	row           int
	x             float32
	health, flash int
	eating        bool
}

type lawnPea struct {
	row int
	x   float32
}

type lawnSunDrop struct {
	x, y, land float32
	life       int
}

type lawnMowerState struct {
	x             float32
	running, gone bool
}

// lawnWorld is one round, replaced whole on restart.
type lawnWorld struct {
	state              lawnState
	endless            bool
	steps, sun         int
	selected           int // a seed, or -1
	recharge           [3]int
	plants             [lawnRows][lawnCols]*lawnPlant
	zombies            []*lawnZombie
	peas               []lawnPea
	suns               []*lawnSunDrop
	mowers             [lawnRows]lawnMowerState
	spawned, defeated  int
	nextSpawn, nextSun int
}

func newLawnWorld(endless bool) lawnWorld {
	w := lawnWorld{endless: endless, sun: 150, selected: -1, nextSpawn: 20 * lawnSecond, nextSun: 4 * lawnSecond}
	for row := range w.mowers {
		w.mowers[row].x = 6
	}
	return w
}

// lawnDemo is the game.Game. Status shown by Onevra controls outside the
// view rebuilds it through invalidate, because the loop only repaints.
type lawnDemo struct {
	lawnWorld
	sheet      *ui.Bitmap
	screen     *ui.Canvas
	rng        *rand.Rand
	sfx        *lawnAudio
	output     bool
	invalidate func()
	// fitX, fitY and fitK map view points to screen pixels.
	fitX, fitY, fitK float32
	hoverX, hoverY   float32
	hover            bool
}

func (d *lawnDemo) init(c *ui.Context) {
	if d.sheet == nil {
		d.sheet = ui.NewBitmap(pixelSheet(lawnTile, lawnPalette, lawnSprites...)).Pixelated()
		d.screen = ui.NewCanvas(lawnW, lawnH)
		d.screen.Bitmap().Pixelated()
		d.rng = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
		d.sfx = newLawnAudio()
		d.lawnWorld = newLawnWorld(false)
	}
	d.invalidate = c.Invalidate
}

// restart starts a fresh round and keeps the chosen game mode.
func (d *lawnDemo) restart() {
	d.lawnWorld = newLawnWorld(d.endless)
	d.playMusic()
	d.changed()
}

func (d *lawnDemo) changed() {
	if d.invalidate != nil {
		d.invalidate()
	}
}

func (d *lawnDemo) setState(s lawnState) bool {
	if d.state == s {
		return false
	}
	d.state = s
	d.changed()
	return true
}

func (d *lawnDemo) ready(seed int) bool {
	return d.recharge[seed] == 0 && d.sun >= lawnSeeds[seed].cost
}

func (d *lawnDemo) Update(in *game.Input) {
	x, y, over := in.Pointer()
	d.hoverX, d.hoverY, d.hover = (x-d.fitX)*d.fitK, (y-d.fitY)*d.fitK, over
	if d.state == lawnLost || d.state == lawnWon {
		if in.Pressed(ui.KeyR) {
			d.restart()
		}
		return
	}
	if !in.Focused() {
		d.setState(lawnPaused)
		return
	}
	if d.state == lawnPaused {
		// The click that resumes plants nothing.
		d.setState(lawnPlaying)
		return
	}
	if in.Pressed(ui.KeyR) {
		d.restart()
		return
	}
	for seed, key := range []ui.Key{ui.Key1, ui.Key2, ui.Key3} {
		if in.Pressed(key) && d.ready(seed) {
			d.selected = seed
		}
	}
	if in.Pressed(ui.KeyEscape) || in.PointerPressed(1) {
		d.selected = -1
	}
	if in.PointerPressed(0) {
		d.click(d.hoverX, d.hoverY)
	}
	d.step()
}

// click collects a sun, picks a seed packet or plants the picked seed.
func (d *lawnDemo) click(x, y float32) {
	for i := len(d.suns) - 1; i >= 0; i-- {
		if s := d.suns[i]; x >= s.x-7 && x <= s.x+7 && y >= s.y-7 && y <= s.y+7 {
			d.sun += lawnSunValue
			d.suns = slices.Delete(d.suns, i, i+1)
			d.play(d.sfx.collect)
			return
		}
	}
	for seed := range lawnSeeds {
		if x >= lawnCardX(seed) && x < lawnCardX(seed)+20 && y >= 1 && y < 23 {
			if d.selected == seed {
				d.selected = -1
			} else if d.ready(seed) {
				d.selected = seed
			}
			return
		}
	}
	row, col, ok := lawnTileAt(x, y)
	if !ok || d.selected < 0 || d.plants[row][col] != nil {
		return
	}
	seed := lawnSeeds[d.selected]
	timer := 30
	if d.selected == lawnSunflower {
		timer = 6 * lawnSecond
	}
	d.plants[row][col] = &lawnPlant{d.selected, seed.health, timer}
	d.sun -= seed.cost
	d.recharge[d.selected] = seed.recharge
	d.selected = -1
	d.play(d.sfx.plant)
}

func (d *lawnDemo) step() {
	d.steps++
	for i := range d.recharge {
		d.recharge[i] = max(0, d.recharge[i]-1)
	}
	d.spawn()
	d.dropSun()
	d.updatePlants()
	d.updatePeas()
	d.updateZombies()
	d.updateMowers()
	d.updateSuns()
	for _, z := range d.zombies {
		if z.health > 0 && z.x+3 <= 0 {
			if d.setState(lawnLost) {
				d.stopMusic()
				d.play(d.sfx.lose)
			}
			return
		}
	}
	if !d.endless && d.spawned == lawnZombies && len(d.zombies) == 0 {
		if d.setState(lawnWon) {
			d.stopMusic()
			d.play(d.sfx.win)
		}
	}
}

// spawn sends zombies ever more often, then a final wave together. A
// level stops after its last wave; endless mode keeps sending waves.
func (d *lawnDemo) spawn() {
	if !d.endless && d.spawned == lawnZombies {
		return
	}
	if d.nextSpawn--; d.nextSpawn > 0 {
		return
	}
	if d.spawned < lawnZombies-lawnFinalWave {
		d.addZombie(0)
		d.nextSpawn = max(6, 14-d.spawned) * lawnSecond
		if d.spawned == lawnZombies-lawnFinalWave {
			d.nextSpawn = 12 * lawnSecond
		}
	} else {
		for i := range lawnFinalWave {
			d.addZombie(float32(i * 10))
		}
		d.nextSpawn = 12 * lawnSecond
	}
	d.changed()
}

func (d *lawnDemo) addZombie(behind float32) {
	d.zombies = append(d.zombies, &lawnZombie{row: d.rng.IntN(lawnRows), x: lawnW + behind, health: lawnZombieHealth})
	d.spawned++
}

func (d *lawnDemo) dropSun() {
	if d.nextSun--; d.nextSun > 0 {
		return
	}
	d.nextSun = lawnSkySunEvery
	x := lawnLeft + lawnTile/2 + float32(d.rng.IntN(lawnCols))*lawnTile
	land := lawnTop + lawnTile/2 + float32(d.rng.IntN(lawnRows))*lawnTile
	d.suns = append(d.suns, &lawnSunDrop{x: x, y: 8, land: land, life: lawnSunLife})
}

func (d *lawnDemo) updatePlants() {
	for row := range lawnRows {
		for col, p := range d.plants[row] {
			if p == nil {
				continue
			}
			x, y := lawnTileX(col), lawnTileY(row)
			switch p.seed {
			case lawnSunflower:
				if p.timer--; p.timer <= 0 {
					p.timer = lawnFlowerSunEach
					d.suns = append(d.suns, &lawnSunDrop{x: x + 11, y: y + 4, land: y + 10, life: lawnSunLife})
				}
			case lawnPeashooter:
				if d.zombieAhead(row, x) {
					if p.timer--; p.timer <= 0 {
						p.timer = lawnShotInterval
						d.peas = append(d.peas, lawnPea{row, x + 13})
						d.playShot()
					}
				}
			}
		}
	}
}

// zombieAhead reports a zombie on screen in row, right of x.
func (d *lawnDemo) zombieAhead(row int, x float32) bool {
	for _, z := range d.zombies {
		if z.row == row && z.health > 0 && z.x+3 > x && z.x < lawnW-4 {
			return true
		}
	}
	return false
}

func (d *lawnDemo) updatePeas() {
	peas := d.peas[:0]
	for _, p := range d.peas {
		p.x += lawnPeaSpeed
		if z := d.zombieAt(p.row, p.x+2); z != nil {
			z.health -= lawnPeaDamage
			z.flash = 6
			continue
		}
		if p.x < lawnW {
			peas = append(peas, p)
		}
	}
	d.peas = peas
}

func (d *lawnDemo) zombieAt(row int, x float32) *lawnZombie {
	for _, z := range d.zombies {
		if z.row == row && z.health > 0 && x >= z.x+3 && x <= z.x+13 {
			return z
		}
	}
	return nil
}

func (d *lawnDemo) updateZombies() {
	zombies := d.zombies[:0]
	for _, z := range d.zombies {
		if z.health <= 0 {
			d.defeated++
			d.playDefeat()
			d.changed()
			continue
		}
		z.flash = max(0, z.flash-1)
		z.eating = false
		// Zombies reach a plant with their arms, standing beside it.
		if col := lawnColAt(z.x + 6); col >= 0 {
			if p := d.plants[z.row][col]; p != nil {
				z.eating = true
				if p.health -= lawnBite; p.health <= 0 {
					d.plants[z.row][col] = nil
				}
			}
		}
		if !z.eating {
			z.x -= lawnZombieSpeed
		}
		zombies = append(zombies, z)
	}
	d.zombies = zombies
}

// updateMowers starts a row's mower when a zombie reaches the patio; it
// clears the row once.
func (d *lawnDemo) updateMowers() {
	for row := range d.mowers {
		m := &d.mowers[row]
		if m.gone {
			continue
		}
		for _, z := range d.zombies {
			if z.row == row && z.health > 0 && z.x+3 <= lawnLeft {
				if !m.running {
					d.play(d.sfx.mower)
				}
				m.running = true
			}
		}
		if !m.running {
			continue
		}
		m.x += 3
		for _, z := range d.zombies {
			if z.row == row && z.x+3 <= m.x+lawnTile && z.x+13 >= m.x {
				z.health = 0
			}
		}
		m.gone = m.x > lawnW
	}
}

func (d *lawnDemo) updateSuns() {
	suns := d.suns[:0]
	for _, s := range d.suns {
		if s.y < s.land {
			s.y = min(s.land, s.y+0.4)
		} else {
			s.life--
		}
		if s.life > 0 {
			suns = append(suns, s)
		}
	}
	d.suns = suns
}

func lawnTileX(col int) float32  { return float32(lawnLeft + col*lawnTile) }
func lawnTileY(row int) float32  { return float32(lawnTop + row*lawnTile) }
func lawnCardX(seed int) float32 { return float32(34 + seed*24) }

func lawnColAt(x float32) int {
	if x < lawnLeft || x >= lawnRight {
		return -1
	}
	return int(x-lawnLeft) / lawnTile
}

func lawnTileAt(x, y float32) (row, col int, ok bool) {
	col = lawnColAt(x)
	if col < 0 || y < lawnTop || y >= lawnH {
		return 0, 0, false
	}
	return int(y-lawnTop) / lawnTile, col, true
}

func lawnFrame(i int) ui.Rect {
	return ui.Rect{X: float32(i * lawnTile), W: lawnTile, H: lawnTile}
}

func (d *lawnDemo) Draw(p *ui.Painter, r ui.Rect, alpha float64) {
	d.screen.Draw(p, d.paint)
	p.Fill(r, ui.RGB(72, 52, 34), 0) // letterbox
	dst, _ := game.IntegerFit(p, lawnW, lawnH, r)
	d.fitX, d.fitY, d.fitK = dst.X-r.X, dst.Y-r.Y, lawnW/dst.W
	p.Image(d.screen.Bitmap(), dst, ui.FillBox)
	if d.state != lawnPlaying {
		p.Fill(dst, ui.RGBA(0, 0, 0, .35), 0)
	}
}

func (d *lawnDemo) paint(cp *ui.Painter) {
	tile := func(row, col int) ui.Rect {
		return ui.Rect{X: lawnTileX(col), Y: lawnTileY(row), W: lawnTile, H: lawnTile}
	}
	sprite := func(frame int, x, y float32, o ui.ImageOptions) {
		cp.ImageRegion(d.sheet, lawnFrame(frame), ui.Rect{X: x, Y: y, W: lawnTile, H: lawnTile}, o)
	}
	cp.Fill(ui.Rect{W: lawnW, H: lawnTop}, ui.RGB(96, 70, 44), 0)
	cp.Fill(ui.Rect{Y: lawnTop, W: lawnLeft, H: lawnH - lawnTop}, ui.RGB(198, 172, 130), 0)
	cp.Fill(ui.Rect{X: lawnRight, Y: lawnTop, W: lawnW - lawnRight, H: lawnH - lawnTop}, ui.RGB(118, 122, 128), 0)
	for row := range lawnRows {
		for col := range lawnCols {
			shade := ui.RGB(112, 180, 76)
			if (row+col)%2 == 1 {
				shade = ui.RGB(98, 164, 66)
			}
			cp.Fill(tile(row, col), shade, 0)
		}
	}
	if d.selected >= 0 && d.hover && d.state == lawnPlaying {
		if row, col, ok := lawnTileAt(d.hoverX, d.hoverY); ok && d.plants[row][col] == nil {
			r := tile(row, col)
			sprite(d.selected, r.X, r.Y, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, .5)})
		}
	}
	for row, m := range d.mowers {
		if !m.gone {
			sprite(lawnMower, m.x, lawnTileY(row), ui.ImageOptions{})
		}
	}
	for row := range lawnRows {
		for col, p := range d.plants[row] {
			if p == nil {
				continue
			}
			var o ui.ImageOptions
			if p.health*3 < lawnSeeds[p.seed].health {
				o.Tint = ui.RGB(200, 160, 140) // badly bitten
			}
			sprite(p.seed, lawnTileX(col), lawnTileY(row), o)
		}
	}
	for _, pea := range d.peas {
		y := lawnTileY(pea.row) + 3
		cp.Fill(ui.Rect{X: pea.x, Y: y, W: 4, H: 4}, ui.RGB(120, 220, 80), 0)
		cp.Fill(ui.Rect{X: pea.x + 2, Y: y + 2, W: 2, H: 2}, ui.RGB(60, 150, 50), 0)
	}
	// Nearer rows overlap farther ones.
	for row := range lawnRows {
		for _, z := range d.zombies {
			if z.row != row {
				continue
			}
			frame, pace := lawnZombieStep, 24
			if z.eating {
				pace = 10
			}
			if (d.steps/pace+int(z.x))%2 == 1 {
				frame = lawnZombieStride
			}
			var o ui.ImageOptions
			if z.flash > 0 {
				o.Tint = ui.RGB(255, 150, 150)
			}
			sprite(frame, float32(int(z.x)), lawnTileY(row)-1, o)
		}
	}
	d.paintSeedBar(cp, sprite)
	for _, s := range d.suns {
		// Fading suns blink before they go.
		if s.life > 2*lawnSecond || s.life/8%2 == 0 {
			sprite(lawnSun, float32(int(s.x))-8, float32(int(s.y))-8, ui.ImageOptions{})
		}
	}
}

func (d *lawnDemo) paintSeedBar(cp *ui.Painter, sprite func(int, float32, float32, ui.ImageOptions)) {
	cp.Fill(ui.Rect{X: 2, Y: 1, W: 28, H: 22}, ui.RGB(66, 48, 30), 2)
	sprite(lawnSun, 8, 0, ui.ImageOptions{})
	lawnNumber(cp, d.sun, 16, 16, ui.RGB(255, 236, 160))
	for seed, s := range lawnSeeds {
		x := lawnCardX(seed)
		card := ui.Rect{X: x, Y: 1, W: 20, H: 22}
		cp.Fill(card, ui.RGB(224, 208, 152), 2)
		sprite(seed, x+2, 1, ui.ImageOptions{})
		lawnNumber(cp, s.cost, x+10, 17, ui.RGB(60, 44, 28))
		if !d.ready(seed) {
			cp.Fill(card, ui.RGBA(0, 0, 0, .35), 2)
		}
		if left := d.recharge[seed]; left > 0 {
			h := float32(int(card.H * float32(left) / float32(s.recharge)))
			cp.Fill(ui.Rect{X: x, Y: 1, W: card.W, H: h}, ui.RGBA(0, 0, 0, .35), 0)
		}
		if d.selected == seed {
			cp.Stroke(card, ui.RGB(255, 232, 90), 2, 1)
		}
	}
	// Level progress, filling as zombies arrive.
	bar := ui.Rect{X: 120, Y: 9, W: 66, H: 6}
	cp.Fill(bar, ui.RGB(50, 36, 22), 0)
	if d.spawned > 0 {
		progress := min(float32(1), float32(d.spawned)/float32(lawnZombies))
		w := float32(int((bar.W - 2) * progress))
		cp.Fill(ui.Rect{X: bar.X + 1, Y: bar.Y + 1, W: w, H: bar.H - 2}, ui.RGB(150, 210, 80), 0)
	}
}

// lawnNumber draws n in the pixel font, centered on cx.
func lawnNumber(cp *ui.Painter, n int, cx, y float32, c ui.Color) {
	s := strconv.Itoa(n)
	x := float32(int(cx - float32(len(s)*4-1)/2))
	for i, digit := range s {
		for j, on := range lawnDigits[digit-'0'] {
			if on == '1' {
				cp.Fill(ui.Rect{X: x + float32(i*4+j%3), Y: y + float32(j/3), W: 1, H: 1}, c, 0)
			}
		}
	}
}

// view shows the round with the controls that drive it. The loop only
// repaints, so the status row is rebuilt through invalidate whenever the
// round changes.
func (d *lawnDemo) view(c *ui.Context) {
	t, m := c.Theme(), locales.In(c.Locale())
	d.init(c)
	d.syncMusic()
	components.Text(c, m.LawnHint()).TextColor(t.TextMuted).Wrap()
	game.View(c, d).Height(320).FillWidth().Radius(12).Clip().Label(m.LawnScene())
	status := map[lawnState]string{
		lawnPlaying: m.LawnPlaying(), lawnPaused: m.LawnPaused(),
		lawnLost: m.LawnLost(), lawnWon: m.LawnWon(),
	}[d.state]
	ui.Row(c).Gap(t.Space(3)).Wrap().AlignItems(ui.Center).Children(func() {
		components.Badge(c, status)
		if d.endless {
			components.Text(c, m.LawnDefeatedEndless(d.defeated)).TextColor(t.TextMuted)
		} else {
			components.Text(c, m.LawnDefeated(d.defeated, lawnZombies)).TextColor(t.TextMuted)
		}
		// Switching modes restarts, so the round matches the chosen mode.
		endless := d.endless
		components.Switch(c, &endless, m.LawnEndless()).Changed()
		if endless != d.endless {
			d.endless = endless
			d.restart()
		}
		music := d.sfx.musicOn
		components.Switch(c, &music, m.LawnMusic()).Changed()
		if music != d.sfx.musicOn {
			d.sfx.musicOn = music
			d.syncMusic()
		}
		if components.Button(c, m.LawnRestart(), components.ButtonSecondary).Clicked() {
			d.restart()
		}
	})
}
