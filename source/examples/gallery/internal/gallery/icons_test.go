package gallery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestIconCatalogCoverage(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "../../../../icons/lucide/icons_gen.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	pattern := regexp.MustCompile(`returns the Lucide ([a-z0-9-]+) icon\.`)
	var names []string
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			match := pattern.FindStringSubmatch(f.Doc.Text())
			if match == nil {
				t.Fatalf("missing name: %s", f.Name.Name)
			}
			names = append(names, match[1])
		}
	}
	var catalogNames []string
	for _, icon := range iconCatalog {
		catalogNames = append(catalogNames, icon.name)
		if icon.get == nil {
			t.Fatalf("nil accessor: %s", icon.name)
		}
	}
	if !reflect.DeepEqual(names, catalogNames) {
		t.Fatal("catalog is stale: rerun tools/lucidegen with the pinned upstream checkout")
	}
}

func TestIconBrowserLayout(t *testing.T) {
	for _, width := range []int{220, 600} {
		g := &Gallery{page: lucidePage}
		tt := ui.NewTester(g.Icons, width, 480)
		if _, ok := tt.Find("Search icons"); ok {
			t.Fatal("search field remains")
		}
		grid, ok := tt.Find("Icon grid")
		if !ok || grid.X != 0 || grid.Y != 0 || grid.W != float32(width) || grid.Y+grid.H != 480 {
			t.Fatalf("grid must fill the content: %v", grid)
		}
		first, ok := tt.Find(iconCatalog[0].name)
		if !ok || first.X < grid.X+16 || first.Y < grid.Y+16 || first.X+first.W > grid.X+grid.W-16 {
			t.Fatalf("grid content must retain internal padding: grid=%v, first=%v", grid, first)
		}
		if first.W != first.H {
			t.Fatalf("icon button must be square: %v", first)
		}
		for _, icon := range iconCatalog[:2] {
			if err := tt.Click(icon.name); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() != icon.name {
				t.Fatalf("copied %q, want %q", tt.Clipboard(), icon.name)
			}
		}
		tt.SetClipboard("")
		tt.Key(0, ui.KeySpace)
		if tt.Clipboard() != iconCatalog[1].name {
			t.Fatal("keyboard activation must copy the focused icon name")
		}
		g.grid.ScrollTo(len(iconCatalog)-1, ui.End)
		tt.Frame()
		last := iconCatalog[len(iconCatalog)-1]
		if err := tt.Click(last.name); err != nil {
			t.Fatal(err)
		}
		if tt.Clipboard() != last.name {
			t.Fatal("scrolled icon copied the wrong name")
		}
	}
}

func TestIconButtonFeedback(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		g := &Gallery{}
		var palette *ui.Theme
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			palette = c.Theme()
			g.Icons(c)
		}, 600, 480)
		r, ok := tt.Find(iconCatalog[0].name)
		if !ok {
			t.Fatal("missing first icon button")
		}
		x, y := r.X+r.W/2, r.Y+r.H/2
		check := func(state string, want ui.Color) {
			t.Helper()
			pixel := tt.Image().RGBAAt(int(r.X+8), int(y))
			if pixel.R != want.R || pixel.G != want.G || pixel.B != want.B {
				t.Fatalf("mode=%v %s: pixel=%v, want=%v", mode, state, pixel, want)
			}
		}
		tt.Move(590, 470)
		check("rest", palette.Background)
		tt.Move(x, y)
		check("hover", palette.SurfaceHover)
		tt.Press(x, y)
		check("pressed", palette.SurfacePressed)
		tt.Release(x, y)
		check("released", palette.SurfaceHover)
		if tt.Clipboard() != iconCatalog[0].name {
			t.Fatal("release must copy the icon name")
		}
		tt.Move(590, 470)
		check("pointer exit", palette.Background)
	}
}

func TestGalleryTitlebar(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Titlebar, 400, 48)
	if !tt.HasText("Icon library / Lucide") {
		t.Fatal(tt.Texts())
	}
	g.page = shellPage
	tt.Frame()
	if !tt.HasText("Shell") {
		t.Fatal(tt.Texts())
	}
}

func TestIconNavigation(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() {
			ui.Column(c).Width(200).FillHeight().Children(func() { g.Left(c) })
			ui.Column(c).Grow(1).FillHeight().Children(func() { g.View(c) })
		})
	}, 820, 700)
	if g.page != lucidePage {
		t.Fatal("default page is not the icon library")
	}
	if err := tt.Click("Lucide"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("Icon grid"); !ok || g.page != lucidePage {
		t.Fatal("Lucide navigation did not show the icon grid")
	}
	if tt.HasText("Shell") || tt.HasText("Application") {
		t.Fatal("removed Shell navigation is still visible")
	}
}
