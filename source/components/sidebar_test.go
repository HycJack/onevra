package components

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/theme"
)

func sidebarFixture(mode theme.Mode, selected *int, disabled bool) func(*ui.Context) {
	return func(c *ui.Context) {
		theme.Apply(c, mode)
		Sidebar(c).Disabled(disabled).Children(func() {
			SidebarSection(c, "Locations", func() {
				if SidebarItem(c, "Home", lucide.House(), *selected == 0).Clicked() {
					*selected = 0
				}
				if SidebarItem(c, "Files", nil, *selected == 1).Clicked() {
					*selected = 1
				}
			})
			SidebarSection(c, "Account", nil)
		})
	}
}

func TestSidebarAppearance(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		selected := 0
		tt := ui.NewTester(sidebarFixture(mode, &selected, false), 220, 300)
		home, _ := tt.Find("Home")
		files, _ := tt.Find("Files")
		heading, _ := tt.Find("Locations")
		account, _ := tt.Find("Account")
		if heading.Y != 0 || home.X != 12 || home.W != 196 || files.Y-home.Y-home.H != 4 || account.Y-files.Y-files.H != 12 {
			t.Fatalf("unexpected sidebar spacing: %v %v %v %v", heading, home, files, account)
		}
		blank := ui.NewTester(func(c *ui.Context) { theme.Apply(c, mode) }, 220, 300)
		if tt.Image().RGBAAt(20, 299) != blank.Image().RGBAAt(20, 299) {
			t.Fatal("sidebar must preserve its host background")
		}
		for _, r := range []ui.Rect{home, files} {
			tt.Move(219, 299)
			before := bytes.Clone(tt.Image().Pix)
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			if !bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("hover changed sidebar appearance")
			}
		}
		if err := tt.Click("Files"); err != nil {
			t.Fatal(err)
		}
		if selected != 1 {
			t.Fatal("click did not select Files")
		}
		pixel := tt.Image().RGBAAt(int(files.X+files.W-12), int(files.Y+files.H/2))
		palette := theme.LightPalette()
		if mode == theme.Dark {
			palette = theme.DarkPalette()
		}
		want := palette.Text.Alpha(.1).Over(palette.Background)
		delta := func(a, b uint8) int { return max(int(a)-int(b), int(b)-int(a)) }
		if pixel.A != 255 || delta(pixel.R, want.R) > 1 || delta(pixel.G, want.G) > 1 || delta(pixel.B, want.B) > 1 {
			t.Fatalf("selection tint = %v, want %v", pixel, want)
		}
	}
}

func TestSidebarKeyboardAndDisabled(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		selected := -1
		tt := ui.NewTester(sidebarFixture(theme.Light, &selected, disabled), 220, 300)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyEnter)
		if !disabled && selected != 0 {
			t.Fatal("Enter did not activate first item")
		}
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		if !disabled && selected != 1 {
			t.Fatal("Space did not activate second item")
		}
		if err := tt.Click("Files"); err != nil {
			t.Fatal(err)
		}
		if disabled && selected != -1 {
			t.Fatal("disabled sidebar accepted input")
		}
	}
}
