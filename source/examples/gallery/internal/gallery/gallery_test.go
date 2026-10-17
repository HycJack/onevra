package gallery

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/sidebar"
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
	terminal := pageFromKey("mygo/Terminal")
	if _, ok := mygoExampleFor(terminal); !ok {
		t.Fatal("missing Terminal route")
	}
	g := &Gallery{page: welcomePage, dev: true}
	if g.pageGrid() {
		t.Fatal("Welcome should keep its dedicated background")
	}
	g.page = pageFromKey("controls/Shader backgrounds")
	if g.pageGrid() {
		t.Fatal("Shader backgrounds should keep their dedicated background")
	}
	g.page = terminal
	if g.pageGrid() || g.exampleGrid() {
		t.Fatal("Terminal should keep its dedicated background")
	}
}

func TestDevGridBehindSelectedMyGoExample(t *testing.T) {
	page := notificationsPage
	render := func(dev bool) (*image.RGBA, float32) {
		var header float32
		g := &Gallery{page: page, dev: dev}
		view := ui.NewTester(func(c *ui.Context) {
			header = sidebar.HeaderHeight(c)
			g.View(c)
		}, 1200, 600)
		return view.Image(), header
	}
	plain, header := render(false)
	grid, _ := render(true)
	top := int(header)
	// The top border separates the host from the floating header.
	if plain.RGBAAt(600, top) == plain.RGBAAt(600, top-2) {
		t.Fatal("missing top border below the header")
	}
	same := func(x0, y0, x1, y1 int) bool {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				if plain.RGBAAt(x, y) != grid.RGBAAt(x, y) {
					return false
				}
			}
		}
		return true
	}
	if !same(0, 0, 1200, top) {
		t.Fatal("grid drawn behind the header")
	}
	if !same(0, top, 150, 600) {
		t.Fatal("grid drawn behind the example list")
	}
	if same(400, 560, 1200, 600) {
		t.Fatal("missing grid behind the selected example")
	}
}

func TestMyGoExampleTitlebar(t *testing.T) {
	g := &Gallery{page: pageFromKey("mygo/Rich text editor")}
	if _, ok := mygoExampleFor(g.page); !ok {
		t.Fatal("missing Rich text editor route")
	}
	title := ui.NewTester(g.Titlebar, 600, 40)
	if !title.HasText("MyGo / Official components") || title.HasText("Text input / Rich text editor") {
		t.Fatalf("shell titlebar: %v", title.Texts())
	}
	// The example pane names the page and carries its actions.
	view := ui.NewTester(g.View, 1200, 600)
	for _, text := range []string{"Text input / Rich text editor", "Bold", "Italic", "Code"} {
		if !view.HasText(text) {
			t.Fatalf("missing %q: %v", text, view.Texts())
		}
	}
}
