package gallery

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/audio"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

const audioRate = 48000

// audioNotes are the pluck's playback speeds relative to its A4 recording.
var audioNotes = []struct {
	name  string
	speed float64
}{{"C4", 261.63 / 440}, {"E4", 329.63 / 440}, {"G4", 392.00 / 440}, {"C5", 523.25 / 440}}

// audioDemo synthesizes its sounds in memory on first use, so the page ships
// no resources. Output is false in headless views, which keep the controls'
// state without opening an audio device.
type audioDemo struct {
	output      bool
	pluck, pad  *audio.Sound
	loop        *audio.Player
	looping     bool
	played      int
	volume, pan float64
}

func (d *audioDemo) sounds() {
	if d.pluck == nil {
		d.pluck, d.pad = audio.NewSound(audioRate, 1, pluckSamples()), audio.NewSound(audioRate, 1, padSamples())
	}
}

func (d *audioDemo) play(speed float64) {
	d.played++
	if !d.output {
		return
	}
	d.sounds()
	d.pluck.NewPlayer().SetSpeed(speed).SetVolume(d.volume).SetPan(d.pan).Play()
}

func (d *audioDemo) toggleLoop() {
	d.looping = !d.looping
	if !d.output {
		return
	}
	if d.loop == nil {
		d.sounds()
		d.loop = d.pad.NewPlayer().SetLoop(true)
	}
	if d.looping {
		d.loop.SetVolume(d.volume).SetPan(d.pan).Play()
	} else {
		d.loop.Pause()
	}
}

// mix applies the shared volume and pan to the loop as the sliders move.
func (d *audioDemo) mix() {
	if d.loop != nil {
		d.loop.SetVolume(d.volume).SetPan(d.pan)
	}
}

func (d *audioDemo) close() {
	if d.loop != nil {
		d.loop.Stop()
	}
	d.loop, d.looping = nil, false
}

// pluckSamples is a decaying A4 with two harmonics, about 0.8 seconds long.
func pluckSamples() []float32 {
	samples := make([]float32, audioRate*8/10)
	for i := range samples {
		t := float64(i) / audioRate
		v := math.Sin(2*math.Pi*440*t) + .4*math.Sin(2*math.Pi*880*t) + .15*math.Sin(2*math.Pi*1320*t)
		attack := min(1, t/.005)
		samples[i] = float32(.35 * v * attack * math.Exp(-5*t))
	}
	return samples
}

// padSamples is two seconds of a soft A-major chord. Every partial and the
// swell complete whole cycles, so the loop has no seam.
func padSamples() []float32 {
	samples := make([]float32, audioRate*2)
	for i := range samples {
		t := float64(i) / audioRate
		v := math.Sin(2*math.Pi*220*t) + .8*math.Sin(2*math.Pi*277*t) + .7*math.Sin(2*math.Pi*330*t)
		swell := .75 + .25*math.Sin(2*math.Pi*.5*t)
		samples[i] = float32(.12 * v * swell)
	}
	return samples
}

func exampleAudio(c *ui.Context, s *demoState) {
	t := c.Theme()
	d := &s.audio
	components.Text(c, tr(c, "gallery.audio_notes")).Bold().TextColor(t.TextMuted)
	ui.Row(c).Gap(t.Space(2)).Wrap().Children(func() {
		for _, note := range audioNotes {
			if components.Button(c, note.name, components.ButtonSecondary).Key(note.name).Clicked() {
				d.play(note.speed)
			}
		}
	})
	label := tr(c, "gallery.audio_play_loop")
	if d.looping {
		label = tr(c, "gallery.audio_pause_loop")
	}
	ui.Row(c).Children(func() {
		if components.Button(c, label, components.ButtonPrimary).Key("loop").Clicked() {
			d.toggleLoop()
		}
	})
	changed := false
	for _, slider := range []struct {
		id     string
		value  *float64
		lo, hi float64
	}{{"gallery.volume", &d.volume, 0, 1}, {"gallery.audio_pan", &d.pan, -1, 1}} {
		ui.Row(c).Key(slider.id).FillWidth().MaxWidth(t.Space(100)).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			name := tr(c, slider.id)
			components.Text(c, name).Width(t.Space(14)).Shrink(0).SingleLine()
			control := components.Slider(c, slider.value, slider.lo, slider.hi).Label(name)
			control.Grow(1).MinWidth(0)
			changed = changed || control.Changed()
		})
	}
	if changed {
		d.mix()
	}
	if err := audio.Err(); d.output && err != nil {
		components.Text(c, fmt.Sprintf(tr(c, "gallery.audio_unavailable"), err)).TextColor(t.Danger).Wrap()
	}
}
