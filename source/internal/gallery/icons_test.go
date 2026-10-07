package gallery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestIconCatalogCoverage(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "../../icons/lucide/icons_gen.go", nil, parser.ParseComments)
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
		t.Fatal("catalog is stale: run go generate ./internal/gallery")
	}
}

func TestFilterIcons(t *testing.T) {
	for _, query := range []string{"", "  "} {
		if len(filterIcons(query)) != len(iconCatalog) {
			t.Fatal("empty search must show all icons")
		}
	}
	for _, query := range []string{"cloud-rain", " CLOUD RAIN ", "rain cloud"} {
		found := false
		for _, icon := range filterIcons(query) {
			if icon.name == "cloud-rain" {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing cloud-rain for %q", query)
		}
	}
	if len(filterIcons("no-such-lucide-icon-xyz")) != 0 {
		t.Fatal("expected no matches")
	}
}

func TestIconBrowserLayout(t *testing.T) {
	for _, width := range []int{220, 600} {
		g := &Gallery{page: lucidePage}
		tt := ui.NewTester(g.Icons, width, 480)
		if tt.HasText("Lucide") || tt.HasText(" of ") {
			t.Fatal("content repeats the title or icon count", tt.Texts())
		}
		grid, ok := tt.Find("Icon grid")
		if !ok || grid.X != 0 || grid.W != float32(width) || grid.Y+grid.H != 480 {
			t.Fatalf("grid must reach the side and bottom edges: %v", grid)
		}
		first, ok := tt.Find(iconCatalog[0].name)
		if !ok || first.X-grid.X != 16 || first.Y-grid.Y < 16 || first.X+first.W > grid.X+grid.W-16 {
			t.Fatalf("grid content must retain internal padding: grid=%v, first=%v", grid, first)
		}
		input, ok := tt.Find("Search icons")
		if !ok || input.H <= 0 || input.Y+input.H > grid.Y {
			t.Fatalf("search field must remain above the grid: %v", input)
		}
		// The decorative icon must not intercept clicks on the input.
		tt.ClickAt(input.X+19, input.Y+input.H/2)
		tt.Type("cloud")
		if g.query != "cloud" {
			t.Fatal("search icon intercepted input focus")
		}
	}
}

func TestGalleryTitlebar(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Titlebar, 400, 48)
	if !tt.HasText("Shell") {
		t.Fatal(tt.Texts())
	}
	g.page = lucidePage
	tt.Frame()
	if !tt.HasText("Icon library / Lucide") {
		t.Fatal(tt.Texts())
	}
}

func TestIconNavigationAndSearch(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() {
			ui.Column(c).Width(200).FillHeight().Children(func() { g.Left(c) })
			ui.Column(c).Grow(1).FillHeight().Children(func() { g.View(c) })
		})
	}, 820, 700)
	if err := tt.Click("Lucide"); err != nil {
		t.Fatal(err)
	}
	if g.page != lucidePage || !tt.HasText("Search icons") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Search icons"); err != nil {
		t.Fatal(err)
	}
	tt.Type("no-such-lucide-icon-xyz")
	if !tt.HasText("No icons found. Try another name.") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if g.query != "" || tt.HasText("No icons found. Try another name.") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Shell"); err != nil {
		t.Fatal(err)
	}
	if g.page != shellPage || !tt.HasText("Clicks: 0") {
		t.Fatal(tt.Texts())
	}
}
