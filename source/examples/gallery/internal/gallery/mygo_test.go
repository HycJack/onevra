package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
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
	groups := map[string]int{"Plugins": 0}
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

func TestMyGoInteractions(t *testing.T) {
	for _, name := range []string{"Checkbox", "Tabs", "Modal", "AlertDialog", "Popover"} {
		t.Run(name, func(t *testing.T) {
			var entry mygoExample
			for _, e := range mygoCatalog {
				if e.name == name {
					entry = e
					break
				}
			}
			s := newDemoState()
			tt := ui.NewTester(func(c *ui.Context) { ui.Column(c).Fill().Children(func() { entry.build(c, s) }) }, 540, 600)
			click := func(label string) {
				t.Helper()
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			switch name {
			case "Checkbox":
				click("Enable notifications")
				if !s.on {
					t.Fatal("checkbox did not toggle")
				}
			case "Tabs":
				click("Appearance")
				if s.index != 1 {
					t.Fatal("tab did not change")
				}
			case "Modal":
				click("Open dialog")
				if !s.open {
					t.Fatal("dialog did not open")
				}
				click("Done")
				if s.open {
					t.Fatal("dialog did not close")
				}
			case "AlertDialog":
				click("Ask for confirmation")
				click("Apply")
				if s.count != 1 || s.open {
					t.Fatal("confirmation did not apply")
				}
			case "Popover":
				click("Show popover")
				if !s.open {
					t.Fatal("popover did not open")
				}
				click("Close")
				if s.open {
					t.Fatal("popover did not close")
				}
			}
		})
	}
}

func TestMyGoGroupedNavigation(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 5000)
	lastY := float32(0)
	for _, group := range mygoGroups {
		heading, ok := tt.Find(group)
		if !ok || heading.Y <= lastY {
			t.Fatalf("missing or unordered group: %s", group)
		}
		lastY = heading.Y
		for _, entry := range mygoCatalog {
			if entry.group != group {
				continue
			}
			item, ok := tt.Find(entry.name)
			if !ok || item.Y <= lastY {
				t.Fatalf("misplaced entry: %s", entry.name)
			}
			lastY = item.Y
		}
	}
	if err := tt.Click("Divider"); err != nil {
		t.Fatal(err)
	}
	if entry, ok := mygoExampleFor(g.page); !ok || entry.name != "Divider" {
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
		ui.Column(c).Disabled(true).Children(func() { exampleCheckbox(c, disabled) })
	}, 400, 200)
	tt.Click("Enable notifications")
	if disabled.on {
		t.Fatal("disabled control accepted input")
	}
}

func TestMyGoNavigationAndState(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 900)
	if err := tt.Click("Button"); err != nil {
		t.Fatal(err)
	}
	if example, ok := mygoExampleFor(g.page); !ok || example.name != "Button" {
		t.Fatal("button route not selected")
	}
	content := ui.NewTester(g.View, 540, 850)
	if err := content.Click("Count"); err != nil {
		t.Fatal(err)
	}
	if g.mygoStates["Button"].count != 1 {
		t.Fatal("button did not update state")
	}
	g.page = mygoPageStart + 1
	ui.NewTester(g.View, 540, 850)
	if g.mygoStates["PrimaryButton"].count != 0 {
		t.Fatal("state leaked between pages")
	}
	g.page = mygoPageStart
	ui.NewTester(g.View, 540, 850)
	if g.mygoStates["Button"].count != 1 {
		t.Fatal("state lost on navigation")
	}
}
