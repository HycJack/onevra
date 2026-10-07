package ui

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestShellAppearance(t *testing.T) {
	for _, tc := range []struct {
		mode theme.Mode
		want ui.Color
	}{
		{theme.Auto, theme.LightPalette().Background}, // headless default is light
		{theme.Light, theme.LightPalette().Background},
		{theme.Dark, theme.DarkPalette().Background},
	} {
		s := newShell(ShellOptions{Appearance: tc.mode})
		var bg ui.Color
		ui.NewTester(func(c *ui.Context) {
			s.view(c)
			bg = c.Theme().Background
		}, 200, 100)
		if bg != tc.want {
			t.Errorf("Appearance %s: background = %v, want %v", tc.mode, bg, tc.want)
		}
	}
}
