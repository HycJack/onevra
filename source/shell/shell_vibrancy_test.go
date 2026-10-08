package shell

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

// Exercise effective material changes, not host OS assumptions. This includes
// unsupported systems, compositor failure and Window.SetVibrancy(None).
func TestShellEffectiveVibrancy(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		t.Run(mode.String(), func(t *testing.T) {
			s := newShell(ShellOptions{Appearance: mode, LeftToolbar: func(c *ui.Context) { ui.Spacer(c) }})
			tt := ui.NewTester(s.view, 820, 600)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true})
			palette := theme.LightPalette()
			if mode == theme.Dark {
				palette = theme.DarkPalette()
			}
			for _, vibrant := range []bool{false, true, false, true} {
				s.SetPage(nil)
				s.SetLeftExpanded(true)
				tt.SetVibrancy(vibrant)
				wantAlpha := uint8(255)
				if vibrant {
					wantAlpha = 0
					if mode == theme.Dark && runtime.GOOS != "windows" {
						wantAlpha = palette.Surface.Alpha(0.65).A
					}
				}
				for _, y := range []int{2, 400, 599} {
					pixel := tt.Image().RGBAAt(100, y)
					if pixel.A != wantAlpha {
						t.Fatalf("material=%v y=%d: sidebar=%v, want alpha %d", vibrant, y, pixel, wantAlpha)
					}
					if !vibrant && (pixel.R != palette.Surface.R || pixel.G != palette.Surface.G || pixel.B != palette.Surface.B) {
						t.Fatalf("fallback does not use theme surface: %v", pixel)
					}
					if tt.Image().RGBAAt(400, y).A != 255 {
						t.Fatal("main content must stay opaque")
					}
				}
				s.SetLeftExpanded(false)
				tt.Frame()
				if tt.Image().RGBAAt(100, 400).A != 255 {
					t.Fatal("collapsed sidebar left transparency")
				}
				s.SetLeftExpanded(true)
				s.SetPage(func(*ui.Context) {})
				tt.Frame()
				if tt.Image().RGBAAt(100, 400).A != 255 {
					t.Fatal("replacement page must stay opaque")
				}
			}
		})
	}
}
