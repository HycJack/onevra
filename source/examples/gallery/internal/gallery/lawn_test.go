package gallery

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestLawnSprites(t *testing.T) {
	img := pixelSheet(lawnTile, lawnPalette, lawnSprites...)
	if b := img.Bounds(); b.Dx() != lawnTile*len(lawnSprites) || b.Dy() != lawnTile {
		t.Fatalf("sheet bounds %v", b)
	}
	for f, rows := range lawnSprites {
		if len(rows) != lawnTile {
			t.Fatalf("frame %d has %d rows", f, len(rows))
		}
		for y, row := range rows {
			if len(row) != lawnTile {
				t.Fatalf("frame %d row %d has %d pixels", f, y, len(row))
			}
			for _, p := range []byte(row) {
				if _, ok := lawnPalette[p]; !ok && p != '.' {
					t.Fatalf("frame %d uses unknown pixel %q", f, p)
				}
			}
		}
	}
	for d, glyph := range lawnDigits {
		if len(glyph) != 15 {
			t.Fatalf("digit %d has %d pixels", d, len(glyph))
		}
	}
}

// lawnTester shows the game page and lets its loop run a few steps.
func lawnTester(t *testing.T, lang string) (*Gallery, *lawnDemo, *ui.Tester, func(x, y float32)) {
	g, s, tt := featureTester("Lawn defense", "Default", lang, 400)
	d := &s.lawn
	// Real frame times vary, so no sun or zombie arrives by itself.
	d.nextSun, d.nextSpawn = 1<<30, 1<<30
	advance(tt)
	view, ok := tt.Find(g.text("gallery.lawn_scene"))
	if !ok {
		t.Fatal("missing lawn view")
	}
	// click presses at a pixel of the lawn's logical screen.
	click := func(x, y float32) {
		tt.ClickAt(view.X+d.fitX+x/d.fitK, view.Y+d.fitY+y/d.fitK)
	}
	return g, d, tt, click
}

// advance runs a frame; the tester paints up to 20 for the running game,
// enough time for its steps.
func advance(tt *ui.Tester) {
	time.Sleep(20 * time.Millisecond)
	tt.Frame()
}

func TestLawnPlanting(t *testing.T) {
	_, d, tt, click := lawnTester(t, "en")
	if d.state != lawnPlaying || d.sun != 150 || d.fitK == 0 {
		t.Fatalf("start: state %v, sun %d, fit %v", d.state, d.sun, d.fitK)
	}
	tt.KeyDown(ui.Key2)
	advance(tt)
	tt.KeyUp(ui.Key2)
	if d.selected != lawnPeashooter {
		t.Fatalf("key 2 picked %d", d.selected)
	}
	click(lawnTileX(1)+8, lawnTileY(2)+8)
	if p := d.plants[2][1]; p == nil || p.seed != lawnPeashooter || d.sun != 50 || d.selected != -1 {
		t.Fatalf("planting: %+v, sun %d, selected %d", p, d.sun, d.selected)
	}
	// The peashooter is unaffordable; the sunflower is not.
	click(lawnCardX(lawnPeashooter)+10, 12)
	if d.selected != -1 {
		t.Fatal("picked an unaffordable seed")
	}
	click(lawnCardX(lawnSunflower)+10, 12)
	click(lawnTileX(1)+8, lawnTileY(2)+8)
	if d.selected != lawnSunflower || d.sun != 50 {
		t.Fatalf("planted on an occupied tile: selected %d, sun %d", d.selected, d.sun)
	}
	click(lawnTileX(0)+8, lawnTileY(2)+8)
	if p := d.plants[2][0]; p == nil || p.seed != lawnSunflower || d.sun != 0 {
		t.Fatalf("sunflower: %+v, sun %d", p, d.sun)
	}
	// The sunflower may have made a sun of its own meanwhile.
	sun := &lawnSunDrop{x: 100, y: 60, land: 60, life: lawnSunLife}
	d.suns = append(d.suns, sun)
	click(102, 58)
	if d.sun != lawnSunValue || slices.Contains(d.suns, sun) {
		t.Fatalf("collecting: sun %d, suns %v", d.sun, d.suns)
	}
}

func TestLawnRecharge(t *testing.T) {
	_, d, _, _ := lawnTester(t, "en")
	d.sun, d.selected = 500, lawnWallNut
	d.click(lawnTileX(2)+8, lawnTileY(1)+8)
	if d.recharge[lawnWallNut] != lawnSeeds[lawnWallNut].recharge || d.ready(lawnWallNut) {
		t.Fatalf("recharge %d", d.recharge[lawnWallNut])
	}
	d.click(lawnCardX(lawnWallNut)+10, 12)
	if d.selected != -1 {
		t.Fatal("picked a recharging seed")
	}
	for range lawnSeeds[lawnWallNut].recharge {
		d.step()
	}
	d.click(lawnCardX(lawnWallNut)+10, 12)
	if d.selected != lawnWallNut {
		t.Fatalf("recharged seed not picked: %d", d.recharge[lawnWallNut])
	}
}

func TestLawnDefense(t *testing.T) {
	g, d, tt, _ := lawnTester(t, "zh-Hans")
	d.plants[3][0] = &lawnPlant{lawnPeashooter, 300, 1}
	d.plants[3][5] = &lawnPlant{lawnWallNut, lawnSeeds[lawnWallNut].health, 0}
	d.zombies = []*lawnZombie{{row: 3, x: lawnRight, health: lawnZombieHealth}}
	d.spawned, d.nextSpawn = 1, 1<<30
	for range 60 * lawnSecond {
		d.step()
		if d.defeated == 1 {
			break
		}
	}
	if d.defeated != 1 || len(d.zombies) != 0 || d.plants[3][0] == nil {
		t.Fatalf("defense: defeated %d, %d zombies, shooter %v", d.defeated, len(d.zombies), d.plants[3][0])
	}
	if d.plants[3][5].health == lawnSeeds[lawnWallNut].health {
		t.Fatal("the zombie never reached the wall-nut")
	}
	advance(tt)
	if !tt.HasText(fmt.Sprintf(g.text("gallery.lawn_defeated"), 1, lawnZombies)) {
		t.Fatalf("missing defeated count: %q", tt.Texts())
	}
}

func TestLawnMowersAndDefeat(t *testing.T) {
	g, d, tt, _ := lawnTester(t, "en")
	d.nextSpawn = 1 << 30
	d.zombies = []*lawnZombie{{row: 0, x: lawnLeft, health: lawnZombieHealth}, {row: 0, x: lawnLeft + 60, health: lawnZombieHealth}}
	for range 2 * lawnSecond {
		d.step()
	}
	if len(d.zombies) != 0 || !d.mowers[0].gone || d.defeated != 2 || d.state != lawnPlaying {
		t.Fatalf("mower: %d zombies, %+v, defeated %d", len(d.zombies), d.mowers[0], d.defeated)
	}
	d.zombies = []*lawnZombie{{row: 0, x: 0, health: lawnZombieHealth}}
	for range 2 * lawnSecond {
		d.step()
	}
	advance(tt)
	if d.state != lawnLost || !tt.HasText(g.text("gallery.lawn_lost")) {
		t.Fatalf("state %v: %q", d.state, tt.Texts())
	}
	if err := tt.Click(g.text("gallery.lawn_restart")); err != nil {
		t.Fatal(err)
	}
	// The button took the keyboard, so the new round waits for the lawn.
	if d.state != lawnPaused || d.sun != 150 || d.mowers[0].gone || len(d.zombies) != 0 {
		t.Fatalf("restart left state %v, sun %d, mower %+v", d.state, d.sun, d.mowers[0])
	}
}

func TestLawnLevel(t *testing.T) {
	_, d, _, _ := lawnTester(t, "en")
	d.nextSpawn = newLawnWorld().nextSpawn
	for range 6 * 60 * lawnSecond {
		if d.spawned == lawnZombies {
			break
		}
		d.step()
		// Defeat arrivals at once so the level runs to its end.
		for _, z := range d.zombies {
			z.health = 0
		}
	}
	if d.spawned != lawnZombies || d.state != lawnWon {
		d.step()
	}
	if d.spawned != lawnZombies || d.state != lawnWon || d.defeated != lawnZombies {
		t.Fatalf("spawned %d, defeated %d, state %v", d.spawned, d.defeated, d.state)
	}
}

func TestLawnPausesWithoutFocus(t *testing.T) {
	g, d, tt, click := lawnTester(t, "en")
	tt.SetFocused(false)
	advance(tt)
	steps := d.steps
	advance(tt)
	if d.state != lawnPaused || d.steps != steps || !tt.HasText(g.text("gallery.lawn_paused")) {
		t.Fatalf("state %v, steps %d then %d", d.state, steps, d.steps)
	}
	// Returning to the window returns to the game.
	tt.SetFocused(true)
	advance(tt)
	if d.state != lawnPlaying || d.steps == steps {
		t.Fatalf("window focus: state %v", d.state)
	}
	// Another control takes the keyboard; the click back on the lawn
	// resumes without planting.
	if err := tt.Click(g.text("gallery.lawn_restart")); err != nil {
		t.Fatal(err)
	}
	advance(tt)
	if d.state != lawnPaused {
		t.Fatalf("state %v after focusing the button", d.state)
	}
	d.selected = lawnSunflower
	click(lawnTileX(4)+8, lawnTileY(0)+8)
	advance(tt)
	if d.state != lawnPlaying || d.plants[0][4] != nil {
		t.Fatalf("resuming: state %v, plant %+v", d.state, d.plants[0][4])
	}
}
