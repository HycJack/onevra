package gallery

import (
	"math"
	"math/rand/v2"

	"github.com/egoist/mygo/audio"
)

// lawnAudio holds the game's effects, synthesized in memory so the page
// ships no audio resources, and the background music's player. It is
// silent in headless views, which never open a device. Effects that can
// sound many times a step, as a volley of peas or a mower's row, keep the
// step of the last one to throttle them.
type lawnAudio struct {
	plant, shoot, collect, defeat, mower, win, lose *audio.Sound
	music                                           *audio.Player
	musicOn                                         bool
	rng                                             *rand.Rand
	shootAt, defeatAt                               int
}

func newLawnAudio() *lawnAudio {
	return &lawnAudio{
		// A low thud of soil.
		plant: lawnSound(.14, func(t float64) float64 {
			return (lawnNote(150, t) + .5*lawnNote(100, t)) * math.Exp(-t*18) * .35
		}),
		// A short pop.
		shoot: lawnSound(.06, func(t float64) float64 {
			return lawnNote(900, t) * math.Exp(-t*50) * .2
		}),
		// A two-note chime.
		collect: lawnSound(.24, func(t float64) float64 {
			return (lawnNote(660, t) + .7*lawnNote(990, t-.07)) * math.Exp(-t*16) * .3
		}),
		// A zombie's crunch.
		defeat: lawnSound(.2, func(t float64) float64 {
			return (lawnNote(240, t) + .6*lawnNote(170, t-.05)) * math.Exp(-t*20) * .3
		}),
		// A mower's rumble.
		mower: lawnSound(.45, func(t float64) float64 {
			return (lawnNote(80, t) + .5*lawnNote(96, t) + .25*lawnNote(48, t)) * math.Exp(-t*8) * .35
		}),
		// A rising fanfare, and a falling one for defeat.
		win: lawnSound(.6, func(t float64) float64 {
			return lawnArp([]float64{523.25, 659.25, 783.99, 1046.5}, .13, t)
		}),
		lose: lawnSound(.85, func(t float64) float64 {
			return lawnArp([]float64{392, 311, 247, 196}, .18, t)
		}),
		rng:      rand.New(rand.NewPCG(3, 5)),
		musicOn:  true,
		shootAt:  -1 << 30,
		defeatAt: -1 << 30,
	}
}

// lawnSound renders d seconds of an effect from f, a function of time,
// fading in and out so it does not click.
func lawnSound(d float64, f func(t float64) float64) *audio.Sound {
	s := make([]float32, int(d*audioRate))
	for i := range s {
		t := float64(i) / audioRate
		edge := min(t/.002, (d-t)/.008, 1)
		s[i] = float32(f(t) * edge)
	}
	return audio.NewSound(audioRate, 1, s)
}

// lawnNote is a sine of frequency f at time t, silent before its start.
func lawnNote(f, t float64) float64 {
	if t < 0 {
		return 0
	}
	return math.Sin(2 * math.Pi * f * t)
}

// lawnArp plays freqs one after another, step apart.
func lawnArp(freqs []float64, step, t float64) float64 {
	var s float64
	for i, f := range freqs {
		at := t - float64(i)*step
		s += lawnNote(f, at) * math.Exp(-max(at, 0)*7)
	}
	return s * .28
}

// play plays one effect unless the view is headless.
func (d *lawnDemo) play(s *audio.Sound) {
	if s == nil || !d.output {
		return
	}
	s.Play()
}

// playVaried plays an effect a little higher or lower each time, as games
// vary the ones that sound often.
func (d *lawnDemo) playVaried(s *audio.Sound) {
	if s == nil || !d.output || d.sfx == nil {
		return
	}
	s.NewPlayer().SetSpeed(.94 + d.sfx.rng.Float64()*.12).Play()
}

// playShot throttles a volley of peas, so many peashooters do not stack
// into noise.
func (d *lawnDemo) playShot() {
	if d.sfx == nil || d.steps-d.sfx.shootAt < 5 {
		return
	}
	d.sfx.shootAt = d.steps
	d.playVaried(d.sfx.shoot)
}

// playDefeat throttles the crunch of a mower clearing a row.
func (d *lawnDemo) playDefeat() {
	if d.sfx == nil || d.steps-d.sfx.defeatAt < 6 {
		return
	}
	d.sfx.defeatAt = d.steps
	d.playVaried(d.sfx.defeat)
}
