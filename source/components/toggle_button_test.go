package components

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestToggleButton(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, disabled := range []bool{false, true} {
			on, changes := false, 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Disabled(disabled).Children(func() {
					if ToggleButton(c, &on, "Bold").Changed() {
						changes++
					}
				})
			}, 300, 100)
			tt.Click("Bold")
			if on != !disabled {
				t.Fatalf("click: on=%v disabled=%v", on, disabled)
			}
			tt.Click("Bold")
			if on {
				t.Fatal("second click did not turn off")
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeySpace)
			if on != !disabled {
				t.Fatalf("Space: on=%v disabled=%v", on, disabled)
			}
			want := 3
			if disabled {
				want = 0
			}
			if changes != want {
				t.Fatalf("changes=%d want=%d", changes, want)
			}
		}
	}
}
