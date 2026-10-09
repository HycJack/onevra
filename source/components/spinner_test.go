package components

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestSpinnerRendering(t *testing.T) {
	for _, variant := range []SpinnerVariant{SpinnerRing, SpinnerSpokes, SpinnerVariant(99)} {
		for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
			for _, scale := range []float32{1, 1.5} {
				t.Run(fmt.Sprint(variant, "/", mode, "/", scale), func(t *testing.T) {
					var diameter float32
					tt := ui.NewTester(func(c *ui.Context) {
						theme.Apply(c, mode)
						diameter = c.Theme().Space(8)
						ui.Column(c).Fill().Children(func() {
							Spinner(c, variant).Label("default")
							Spinner(c, variant).Size(80, 40).Label("wide")
						})
					}, 200, 150)
					tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
					if r, ok := tt.Find("default"); !ok || r.W != diameter || r.H != diameter {
						t.Fatalf("default bounds = %+v, want diameter %v", r, diameter)
					}
					if r, ok := tt.Find("wide"); !ok || r.W != 80 || r.H != 40 {
						t.Fatalf("custom bounds = %+v", r)
					}
					before := append([]byte(nil), tt.Image().Pix...)
					tt.Frame()
					if !bytes.Equal(before, tt.Image().Pix) {
						t.Fatal("reduced-motion rendering changed between frames")
					}
				})
			}
		}
	}
}

func TestSpinnerRingMatchesProgressRing(t *testing.T) {
	render := func(spinner bool) []byte {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, theme.Light)
			if spinner {
				Spinner(c, SpinnerRing).Size(80, 40)
			} else {
				ProgressRing(c, -1).Size(80, 40)
			}
		}, 100, 60)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		return append([]byte(nil), tt.Image().Pix...)
	}
	if !bytes.Equal(render(true), render(false)) {
		t.Fatal("ring spinner differs from indeterminate ProgressRing")
	}
}
