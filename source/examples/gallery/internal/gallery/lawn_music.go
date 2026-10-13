package gallery

import (
	"io"
	"math"

	"github.com/egoist/mygo/audio"
)

// The lawn tune is a short, brooding loop in A minor for the game: a
// pulsing bass, a soft chord pad, a music-box arpeggio, a whistled lead
// and a light drum kit. It is composed and played as it goes, so the page
// ships no audio resources, and its decoder is endless: the same eight
// bars repeat until the player stops.
const (
	lawnTuneRate = 32000
	lawnTuneBPM  = 96
)

var (
	lawnBeat = 60 / float64(lawnTuneBPM) // seconds
	lawnBar  = lawnBeat * 4              // 4/4, seconds
	lawnLoop = lawnBar * 8               // eight bars, seconds
)

// lawnHarmony is the loop's bars: the bass root and the triad over it, as
// MIDI notes, over a i–VI–VII–i progression.
var lawnHarmony = [8]struct {
	root  int
	triad [3]int
}{
	{45, [3]int{57, 60, 64}}, // Am
	{41, [3]int{53, 57, 60}}, // F
	{43, [3]int{55, 59, 62}}, // G
	{45, [3]int{57, 60, 64}}, // Am
	{45, [3]int{57, 60, 64}}, // Am
	{41, [3]int{53, 57, 60}}, // F
	{40, [3]int{52, 56, 59}}, // E
	{45, [3]int{57, 60, 64}}, // Am
}

// lawnMelody is the lead, a note or a rest (0) per beat, eight bars long.
// It leans on F, the flat sixth, and G♯, the raised seventh of the
// harmonic minor, for its unease, and comes to rest on the loop's A.
var lawnMelody = [32]int{
	69, 72, 76, 0,
	74, 72, 71, 0,
	69, 72, 77, 0,
	76, 74, 72, 0,
	76, 79, 76, 0,
	74, 76, 72, 0,
	71, 74, 79, 0,
	76, 80, 83, 0,
}

// lawnTune is the music as a Decoder, endless so it loops without a seam.
type lawnTune struct{ pos int64 }

func (t *lawnTune) SampleRate() int { return lawnTuneRate }
func (t *lawnTune) Channels() int   { return 2 }
func (t *lawnTune) Len() int64      { return -1 }

func (t *lawnTune) SetPosition(f int64) error {
	t.pos = f
	return nil
}

func (t *lawnTune) Read(buf []float32) (int, error) {
	n := 0
	for ; n+2 <= len(buf); n += 2 {
		at := math.Mod(float64(t.pos)/lawnTuneRate, lawnLoop)
		bar := int(at / lawnBar)
		inBar := at - float64(bar)*lawnBar
		beat := inBar / lawnBeat
		chord := lawnHarmony[bar]

		bass := lawnBass(midiFreq(chord.root), beat)
		var pad float64
		for _, m := range chord.triad {
			pad += lawnPad(midiFreq(m), inBar)
		}
		triad := [3]float64{midiFreq(chord.triad[0] + 12), midiFreq(chord.triad[1] + 12), midiFreq(chord.triad[2] + 12)}
		arp := lawnMusicArp(triad, inBar)
		var lead float64
		if note := lawnMelody[bar*4+int(beat)]; note > 0 {
			lead = lawnLead(midiFreq(note), (beat-math.Floor(beat))*lawnBeat)
		}
		drums := lawnDrums(beat)

		left := bass + pad + lead*.75 + arp*.35 + drums
		right := bass + pad + lead*.85 + arp*.55 + drums
		buf[n] = float32(math.Tanh(left * .8))
		buf[n+1] = float32(math.Tanh(right * .8))
		t.pos++
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

// midiFreq is a MIDI note's frequency in Hz, A4 (69) at 440.
func midiFreq(m int) float64 { return 440 * math.Pow(2, float64(m-69)/12) }

// lawnBass walks the root on the beats and its fifth between them, a
// rounded square that thumps.
func lawnBass(root, beat float64) float64 {
	f := root
	if int(math.Floor(beat))%2 == 1 {
		f *= math.Pow(2, 7.0/12)
	}
	dt := (beat - math.Floor(beat)) * lawnBeat
	wave := math.Tanh(3 * math.Sin(2*math.Pi*f*dt))
	return wave * math.Exp(-dt*3.5) * min(dt/.004, 1) * .42
}

// lawnPad swells a chord tone in under the melody.
func lawnPad(f, t float64) float64 {
	return math.Sin(2*math.Pi*f*t) * min(t/.08, 1) * math.Exp(-t*.7) * .085
}

// lawnMusicArp cycles the chord an octave up on sixteenths, as a music
// box, damping each note before the next.
func lawnMusicArp(triad [3]float64, t float64) float64 {
	six := t / (lawnBeat / 4)
	i := int(six) % 3
	dt := (six - math.Floor(six)) * (lawnBeat / 4)
	ph := math.Mod(triad[i]*dt, 1)
	tri := 4*math.Abs(ph-.5) - 1
	return tri * math.Exp(-dt*12) * min(dt/.003, 1)
}

// lawnLead is the tune's voice: a soft triangle with a slow vibrato.
func lawnLead(f, t float64) float64 {
	vib := 1 + .008*math.Sin(2*math.Pi*5.5*t)
	ph := math.Mod(f*vib*t, 1)
	tri := 4*math.Abs(ph-.5) - 1
	return tri * math.Exp(-t*2.2) * min(t/.01, 1) * .3
}

// lawnDrums puts a kick on beats 1 and 3, a snare on 2 and 4, and a hat on
// every eighth.
func lawnDrums(beat float64) float64 {
	b := int(beat)
	dt := (beat - float64(b)) * lawnBeat
	var s float64
	if b == 0 || b == 2 {
		s += math.Sin(2*math.Pi*(70+150*math.Exp(-dt*35))*dt) * math.Exp(-dt*12) * .6
	}
	if b == 1 || b == 3 {
		s += (lawnNoise(dt*977+float64(b)*31)*.5 + math.Sin(2*math.Pi*185*dt)*.3) * math.Exp(-dt*22) * .3
	}
	half := math.Mod(beat, .5)
	s += lawnNoise(float64(int(beat*2))*7.13+half*500) * math.Exp(-half*lawnBeat*55) * .06
	return s
}

// lawnNoise is noise that is the same at the same time.
func lawnNoise(x float64) float64 {
	v := math.Sin(x*12.9898) * 43758.5453
	return (v-math.Floor(v))*2 - 1
}

// playMusic starts or resumes the background loop, unless the view is
// headless or the player turned the music off.
func (d *lawnDemo) playMusic() {
	if d.sfx == nil || !d.output || !d.sfx.musicOn {
		return
	}
	if d.sfx.music == nil {
		d.sfx.music = audio.NewPlayer(&lawnTune{}).SetVolume(.3)
	}
	d.sfx.music.Play()
}

func (d *lawnDemo) stopMusic() {
	if d.sfx != nil && d.sfx.music != nil && d.sfx.music.Playing() {
		d.sfx.music.Stop()
	}
}

// syncMusic matches the loop to the round: it plays while a round is
// underway and the player keeps the switch on, and rests otherwise.
func (d *lawnDemo) syncMusic() {
	if d.musicShouldPlay() {
		d.playMusic()
		return
	}
	d.stopMusic()
}

// musicShouldPlay reports whether the loop belongs on now.
func (d *lawnDemo) musicShouldPlay() bool {
	return d.sfx != nil && d.sfx.musicOn && d.state != lawnWon && d.state != lawnLost
}

// closeMusic releases the loop, as when its page is no longer shown.
func (d *lawnDemo) closeMusic() {
	d.stopMusic()
	if d.sfx != nil {
		d.sfx.music = nil
	}
}
