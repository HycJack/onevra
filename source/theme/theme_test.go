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
	for _, tc := range []struct {
		name    string
		palette *ui.Theme
		alpha   float32
	}{
		{"light", LightPalette(), 0.30},
		{"dark", DarkPalette(), 0.35},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Verify the chosen text tint, not an unrelated RGB-distance threshold.
			want := tc.palette.Text.Alpha(tc.alpha)
			if got := tc.palette.Scrollbar; got != want {
				t.Errorf("scrollbar = %v, want text tint at %.0f%% opacity: %v", got, tc.alpha*100, want)
			}
		})
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
