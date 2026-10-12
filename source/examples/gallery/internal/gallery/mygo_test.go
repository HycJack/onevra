package gallery

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func TestMyGoCatalog(t *testing.T) {
	seen := map[string]bool{}
	for i, example := range mygoCatalog {
		if seen[example.name] || example.group == "" || example.build == nil {
			t.Fatalf("invalid catalog entry: %s", example.name)
		}
		seen[example.name] = true
		page := mygoPageStart + galleryPage(i)
		got, ok := mygoExampleFor(page)
		if !ok || got.name != example.name {
			t.Fatalf("wrong route for %s", example.name)
		}
	}
	groups := map[string]int{"Plugins": 0, "Desktop": 0}
	for _, group := range mygoGroups {
		if _, exists := groups[group]; exists {
			t.Fatalf("duplicate group: %s", group)
		}
		groups[group] = 0
	}
	for _, entry := range mygoCatalog {
		if _, ok := groups[entry.group]; !ok {
			t.Fatalf("unknown group: %s", entry.group)
		}
		groups[entry.group]++
	}
	for group, count := range groups {
		if count == 0 {
			t.Fatalf("empty group: %s", group)
		}
	}
	for _, page := range []galleryPage{shellPage, lucidePage, -1, mygoPageStart + galleryPage(len(mygoCatalog))} {
		if _, ok := mygoExampleFor(page); ok {
			t.Fatalf("unexpected MyGo route: %d", page)
		}
	}
}

func TestMyGoPages(t *testing.T) {
	for i, example := range mygoCatalog {
		t.Run(example.name, func(t *testing.T) {
			for _, width := range []int{540, 280} {
				g := &Gallery{page: mygoPageStart + galleryPage(i)}
				tt := ui.NewTester(g.View, width, 850)
				for _, dark := range []bool{false, true} {
					tt.SetDark(dark)
					for _, label := range []string{"Live example", "Go usage", "Demo state and initialization"} {
						if _, ok := tt.Find(label); ok {
							t.Fatalf("unexpected explanatory content: %s", label)
						}
					}
					variants := mygoVariants(example.name)
					if len(g.mygoStates) != len(variants) {
						t.Fatal("missing variant state")
					}
					if len(variants) > 1 {
						for _, label := range variants {
							if _, ok := tt.Find(label); !ok {
								t.Fatalf("missing variant: %s", label)
							}
						}
					}
					if tt.Image() == nil {
						t.Fatal("missing rendered page")
					}
				}
			}
		})
	}
}

func TestMyGoGroupedNavigation(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 5000)
	texts := tt.Texts()
	consume := func(label string) {
		t.Helper()
		i := slices.Index(texts, label)
		if i < 0 {
			t.Fatalf("missing or unordered entry: %s", label)
		}
		texts = texts[i+1:]
	}
	for _, group := range navigationGroups {
		consume(group)
		for _, entry := range controlCatalog {
			if entry.group == group {
				consume(entry.label)
			}
		}
		for _, entry := range mygoCatalog {
			if entry.group == group {
				consume(entry.name)
			}
		}
	}
	if err := tt.Click("Tooltip"); err != nil {
		t.Fatal(err)
	}
	if entry, ok := controlExampleFor(g.page); !ok || entry.name != "Tooltip" {
		t.Fatal("grouping changed route identity")
	}
}

func TestMyGoVariantIsolation(t *testing.T) {
	g := &Gallery{}
	off := g.mygoState("Checkbox", "Off", true)
	on := g.mygoState("Checkbox", "On", false)
	disabled := g.mygoState("Checkbox", "Disabled", false)
	if off.on || !on.on || disabled.on {
		t.Fatal("wrong initial variants")
	}
	off.on = true
	on.on = false
	if g.mygoState("Checkbox", "On", false).on || disabled.on {
		t.Fatal("variant state leaked")
	}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() { components.Checkbox(c, &disabled.on, "Enable notifications") })
	}, 400, 200)
	tt.Click("Enable notifications")
	if disabled.on {
		t.Fatal("disabled control accepted input")
	}
}

func TestMyGoNavigationAndState(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 5000)
	if err := tt.Click("Collapsible"); err != nil {
		t.Fatal(err)
	}
	if example, ok := controlExampleFor(g.page); !ok || example.name != "Collapsible" {
		t.Fatal("collapsible route not selected")
	}
	detailsPage := g.page
	content := ui.NewTester(g.View, 540, 850)
	if err := content.Click("Details"); err != nil {
		t.Fatal(err)
	}
	if !g.controlState("Collapsible", "Right arrow").on {
		t.Fatal("collapsible did not update state")
	}
	if err := tt.Click("Popover"); err != nil {
		t.Fatal(err)
	}
	ui.NewTester(g.View, 540, 850)
	if g.controlState("Popover", "Default").dialogOpen {
		t.Fatal("state leaked between pages")
	}
	g.page = detailsPage
	ui.NewTester(g.View, 540, 850)
	if !g.controlState("Collapsible", "Right arrow").on {
		t.Fatal("state lost on navigation")
	}
}
