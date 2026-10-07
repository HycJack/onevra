package theme

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPalettes(t *testing.T) {
	light, dark := LightPalette(), DarkPalette()
	if light.Dark {
		t.Error("LightPalette must not set Dark")
	}
	if !dark.Dark {
		t.Error("DarkPalette must set Dark")
	}
	if light.Background == dark.Background {
		t.Error("light and dark backgrounds must differ")
	}
}

func TestScrollbarTextTint(t *testing.T) {
	for _, palette := range []*ui.Theme{LightPalette(), DarkPalette()} {
		bar := palette.Scrollbar
		if bar.R != palette.Text.R || bar.G != palette.Text.G || bar.B != palette.Text.B || bar.A == 0 || bar.A == 255 {
			t.Fatalf("dark=%v: scrollbar must be a translucent text tint: %v", palette.Dark, bar)
		}
		for _, bg := range []ui.Color{palette.Background, palette.Surface} {
			visible := bar.Over(bg)
			// Keep the thumb clearly separated from either scrolling surface.
			delta := func(a, b uint8) int { return max(int(a)-int(b), int(b)-int(a)) }
			if delta(visible.R, bg.R) < 70 || delta(visible.G, bg.G) < 70 || delta(visible.B, bg.B) < 70 {
				t.Fatalf("dark=%v: scrollbar %v too close to surface %v", palette.Dark, visible, bg)
			}
		}
	}
}

func TestApplyPinsMode(t *testing.T) {
	for _, tc := range []struct {
		mode Mode
		want ui.Color
	}{
		{Light, LightPalette().Background},
		{Dark, DarkPalette().Background},
	} {
		var got *ui.Theme
		ui.NewTester(func(c *ui.Context) {
			Apply(c, tc.mode)
			got = c.Theme()
		}, 200, 100)
		if got == nil || got.Background != tc.want {
			t.Errorf("Apply(%s): background = %v, want %v", tc.mode, got, tc.want)
		}
	}
}

func TestApplyAutoFollowsSystem(t *testing.T) {
	var got *ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		Apply(c, Auto)
		got = c.Theme()
	}, 200, 100)
	if got.Background != LightPalette().Background {
		t.Fatalf("Auto: initial background = %v, want light %v", got.Background, LightPalette().Background)
	}
	tt.SetDark(true)
	if got.Background != DarkPalette().Background {
		t.Fatalf("Auto: after system dark, background = %v, want dark %v", got.Background, DarkPalette().Background)
	}
	tt.SetDark(false)
	if got.Background != LightPalette().Background {
		t.Fatalf("Auto: after system light, background = %v, want light %v", got.Background, LightPalette().Background)
	}
}
