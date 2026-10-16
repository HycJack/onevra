package main

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/audio"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/games/lawn/locales"
	"github.com/oneo-me/onevra/source/theme"
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

// lawnTester shows the game view and lets its loop run a few steps.
func lawnTester(t *testing.T, lang string) (*lawnDemo, *ui.Tester, func(x, y float32)) {
	t.Helper()
	d := &lawnDemo{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).FillWidth().Padding(16).Gap(12).Children(func() { d.view(c) })
	}, 400, 480)
	tt.SetLocale(lang)
	// No sun or zombie arrives by itself.
	d.nextSun, d.nextSpawn = 1<<30, 1<<30
	advance(tt)
	view, ok := tt.Find(locales.In(lang).LawnScene())
	if !ok {
		t.Fatal("missing lawn view")
	}
	// click presses at a pixel of the lawn's logical screen, which the
	// game's next steps take.
	click := func(x, y float32) {
		tt.ClickAt(view.X+d.fitX+x/d.fitK, view.Y+d.fitY+y/d.fitK)
		advance(tt)
	}
	return d, tt, click
}

// advance runs the game's frames for a few of its steps.
func advance(tt *ui.Tester) {
	tt.Advance(50 * time.Millisecond)
}

func TestLawnPlanting(t *testing.T) {
	d, tt, click := lawnTester(t, "en")
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
	d, _, _ := lawnTester(t, "en")
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
	d, tt, _ := lawnTester(t, "zh-Hans")
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
	if !tt.HasText(locales.In("zh-Hans").LawnDefeated(1, lawnZombies)) {
		t.Fatalf("missing defeated count: %q", tt.Texts())
	}
}

func TestLawnMowersAndDefeat(t *testing.T) {
	d, tt, _ := lawnTester(t, "en")
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
	if d.state != lawnLost || !tt.HasText(locales.In("en").LawnLost()) {
		t.Fatalf("state %v: %q", d.state, tt.Texts())
	}
	if err := tt.Click(locales.In("en").LawnRestart()); err != nil {
		t.Fatal(err)
	}
	advance(tt)
	// The button took the keyboard, so the new round waits for the lawn.
	if d.state != lawnPaused || d.sun != 150 || d.mowers[0].gone || len(d.zombies) != 0 {
		t.Fatalf("restart left state %v, sun %d, mower %+v", d.state, d.sun, d.mowers[0])
	}
}

func TestLawnLevel(t *testing.T) {
	d, _, _ := lawnTester(t, "en")
	d.nextSpawn = newLawnWorld(false).nextSpawn
	for range 12 * 60 * lawnSecond {
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

func TestLawnEndlessMode(t *testing.T) {
	d, tt, _ := lawnTester(t, "en")
	// Turning on endless mode after losing restarts the round in the new
	// mode, since the switch takes the keyboard.
	d.setState(lawnLost)
	advance(tt)
	if err := tt.Click(locales.In("en").LawnEndless()); err != nil {
		t.Fatal(err)
	}
	if !d.endless || d.state == lawnLost {
		t.Fatalf("endless=%v state=%v after turning on endless mode", d.endless, d.state)
	}
	// A level stops at its total; endless mode keeps sending waves.
	d.setState(lawnPlaying)
	d.spawned, d.nextSpawn = lawnZombies, 1
	for range 13 * lawnSecond {
		d.step()
	}
	if d.state != lawnPlaying || d.spawned <= lawnZombies || len(d.zombies) == 0 {
		t.Fatalf("endless: spawned %d, %d zombies, state %v", d.spawned, len(d.zombies), d.state)
	}
	advance(tt)
	if !tt.HasText(locales.In("en").LawnDefeatedEndless(d.defeated)) {
		t.Fatalf("missing endless defeated count: %q", tt.Texts())
	}
}

func TestLawnAudio(t *testing.T) {
	d, _, _ := lawnTester(t, "en")
	if d.output {
		t.Fatal("headless lawn opened audio output")
	}
	if d.sfx == nil {
		t.Fatal("lawn effects were not synthesized")
	}
	effects := map[string]*audio.Sound{
		"plant": d.sfx.plant, "shoot": d.sfx.shoot, "collect": d.sfx.collect,
		"defeat": d.sfx.defeat, "mower": d.sfx.mower, "win": d.sfx.win, "lose": d.sfx.lose,
	}
	for name, s := range effects {
		if s == nil || s.Duration() <= 0 {
			t.Fatalf("%s effect is empty", name)
		}
	}
}

func TestLawnMusic(t *testing.T) {
	d, _, _ := lawnTester(t, "en")
	if !d.sfx.musicOn {
		t.Fatal("background music did not default on")
	}
	if d.sfx.music != nil {
		t.Fatal("headless lawn opened a music player")
	}
	d.syncMusic()
	if d.sfx.music != nil {
		t.Fatal("syncMusic opened a music player without output")
	}
	tune := &lawnTune{}
	if tune.Channels() != 2 || tune.Len() != -1 || tune.SampleRate() != lawnTuneRate {
		t.Fatalf("tune format: %d channels, %d frames, %d Hz", tune.Channels(), tune.Len(), tune.SampleRate())
	}
	buf := make([]float32, tune.SampleRate()) // half a second, stereo
	n, err := tune.Read(buf)
	if err != nil || n != len(buf) {
		t.Fatalf("tune read: %d of %d, %v", n, len(buf), err)
	}
	var peak, sum float64
	stereo := false
	for i := 0; i+1 < len(buf); i += 2 {
		peak = max(peak, math.Abs(float64(buf[i])), math.Abs(float64(buf[i+1])))
		stereo = stereo || buf[i] != buf[i+1]
		sum += float64(buf[i])*float64(buf[i]) + float64(buf[i+1])*float64(buf[i+1])
	}
	if peak == 0 || peak > 1 {
		t.Fatalf("tune peak %v", peak)
	}
	if rms := math.Sqrt(sum / float64(len(buf))); rms < .05 || rms > .5 {
		t.Fatalf("tune rms %v", rms)
	}
	if !stereo {
		t.Fatal("tune is mono")
	}
}

func TestLawnMusicToggle(t *testing.T) {
	d, _, _ := lawnTester(t, "en")
	if !d.musicShouldPlay() {
		t.Fatal("music should play while a round is underway")
	}
	// Turning the switch off must rest the loop, not just stop new plays.
	d.sfx.musicOn = false
	if d.musicShouldPlay() {
		t.Fatal("music still plays after the switch went off")
	}
	d.sfx.musicOn = true
	d.setState(lawnWon)
	if d.musicShouldPlay() {
		t.Fatal("music still plays after the round ended")
	}
	d.setState(lawnPlaying)
	if !d.musicShouldPlay() {
		t.Fatal("music did not resume with the new round")
	}
}

func TestLawnPausesWithoutFocus(t *testing.T) {
	d, tt, click := lawnTester(t, "en")
	tt.SetFocused(false)
	advance(tt)
	steps := d.steps
	advance(tt)
	if d.state != lawnPaused || d.steps != steps || !tt.HasText(locales.In("en").LawnPaused()) {
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
	if err := tt.Click(locales.In("en").LawnRestart()); err != nil {
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

// TestLawnMessages checks that the example's catalog covers both languages.
func TestLawnMessages(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		m := locales.In(lang)
		if m.Title() == "" || m.LawnHint() == "" || m.LawnScene() == "" {
			t.Fatalf("%s: missing messages", lang)
		}
		if got := m.LawnDefeated(2, 100); got == "" {
			t.Fatalf("%s: empty defeated count", lang)
		}
	}
	if locales.In("zh-Hans").LawnRestart() == locales.In("en").LawnRestart() {
		t.Fatal("translations are not distinct")
	}
}
