package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestShellContent(t *testing.T) {
	g := &Gallery{page: shellPage}
	tt := ui.NewTester(g.View, 360, 480)
	if tt.HasText("Shell") || !tt.HasText("Clicks: 0") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Count"); err != nil {
		t.Fatal(err)
	}
	if g.Clicks != 1 || !tt.HasText("Clicks: 1") {
		t.Fatal("caller-owned state was not updated")
	}
}

func TestDevBuildMarkers(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, dev := range []bool{false, true} {
			g := &Gallery{page: pageFromKey("controls/Popover"), dev: dev}
			left := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				g.Left(c)
			}, 220, 500)
			if left.HasText("Dev") != dev {
				t.Fatalf("dark=%v dev=%v: %v", dark, dev, left.Texts())
			}
			view := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				g.View(c)
			}, 900, 600)
			img := view.Image()
			background := img.RGBAAt(5, 575)
			// The grid spans the header inset and the empty space below examples.
			for _, y := range []int{5, 575} {
				for x := 22; x < 890; x += 22 {
					if (img.RGBAAt(x, y) != background) != dev || img.RGBAAt(x+5, y) != background {
						t.Fatalf("dark=%v dev=%v: unexpected grid at %d,%d", dark, dev, x, y)
					}
				}
			}
		}
	}
}

func TestDevGridSkipsDedicatedBackgrounds(t *testing.T) {
	if _, ok := mygoExampleFor(pageFromKey("mygo/Terminal")); !ok {
		t.Fatal("missing Terminal route")
	}
	for _, page := range []galleryPage{welcomePage, pageFromKey("mygo/Terminal")} {
		if (&Gallery{page: page, dev: true}).pageGrid() {
			t.Fatalf("page %v should keep its dedicated background", page)
		}
	}
}
