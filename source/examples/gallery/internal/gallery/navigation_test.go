package gallery

import (
	"bytes"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
)

func TestNavigationSpacing(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 480)
	first, ok := tt.Find("Overview")
	search, searchOK := tt.Find("Search components")
	if !ok || !searchOK || search.Y != 48 || first.Y != search.Y+search.H+8 {
		t.Fatalf("first group must follow the search field at the shell's safe content edge: %v", first)
	}
	item, _ := tt.Find("Component status")
	next, _ := tt.Find("Icon library")
	if gap := next.Y - item.Y - item.H; math.Abs(float64(gap-12)) > 0.001 {
		t.Fatalf("group gap = %v, want 12", gap)
	}
}

func TestNavigationGroupOrder(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	var bottom float32
	for _, label := range []string{"Overview", "Welcome", "Component status", "Icon library", "Lucide", "Desktop", "Tray", "Dock menu", "System notifications", "Plugins", "Markdown", "Terminal", "Actions & navigation", "Buttons", "Text input", "Selection & values", "Date & time", "Color", "Content display", "Layout & containers", "Feedback & overlays"} {
		r, ok := tt.Find(label)
		if !ok || r.Y < bottom {
			t.Fatalf("%s missing or out of order: %v (previous bottom %v)", label, r, bottom)
		}
		bottom = r.Y + r.H
	}
	if tt.HasText("Onevra controls") || tt.HasText("Components") || tt.HasText("MyGo components") {
		t.Fatal("obsolete component group is still visible")
	}
	for _, group := range mygoGroups {
		if group == "Plugins" {
			t.Fatal("plugins must not also appear among built-in MyGo groups")
		}
	}
}

func TestNavigationMyGoMarkers(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	texts := tt.Texts()
	markers := 0
	want := 0
	for _, example := range mygoCatalog {
		if example.group != "Plugins" {
			want++
		}
	}
	for i, text := range texts {
		if text != "•" {
			continue
		}
		markers++
		if i == 0 {
			t.Fatal("marker has no preceding item")
		}
		found := false
		for _, example := range mygoCatalog {
			found = found || (example.group != "Plugins" && texts[i-1] == example.name)
		}
		if !found {
			t.Fatalf("marker follows non-MyGo item: %q", texts[i-1])
		}
	}
	if markers != want {
		t.Fatalf("got %d MyGo markers, want %d", markers, want)
	}
}

func TestNavigationSelectionAndHover(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(map[bool]string{false: "light", true: "dark"}[dark], func(t *testing.T) {
			g := &Gallery{}
			tt := ui.NewTester(g.Left, 220, 800)
			tt.SetDark(dark)
			for _, label := range []string{"Lucide", "Buttons"} {
				tt.Move(219, 479)
				before := bytes.Clone(tt.Image().Pix)
				r, _ := tt.Find(label)
				tt.Move(r.X+r.W/2, r.Y+r.H/2)
				if !bytes.Equal(before, tt.Image().Pix) {
					t.Fatalf("%s changed appearance on hover", label)
				}
			}
			for _, label := range []string{"Buttons", "Lucide"} {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				r, _ := tt.Find(label)
				// Sample the face away from the icon, text and rounded corners.
				pixel := tt.Image().RGBAAt(int(r.X+r.W-12), int(r.Y+r.H/2))
				if pixel.A == 0 || absDiff(pixel.R, pixel.G) > 4 || absDiff(pixel.G, pixel.B) > 4 {
					t.Fatalf("%s selection is not neutral gray: %v", label, pixel)
				}
			}
		})
	}
}

func absDiff(a, b uint8) int {
	return max(int(a)-int(b), int(b)-int(a))
}

func TestNavigationIcons(t *testing.T) {
	for _, group := range navigationCatalog() {
		for _, entry := range group.entries {
			icon := entry.icon()
			if icon == nil || icon == lucide.Component() {
				t.Errorf("%s has no specific navigation icon", entry.label)
			}
			translated := entry
			translated.label = "本地化名称"
			if translated.icon() != icon {
				t.Errorf("%s icon depends on its display label", entry.label)
			}
		}
	}
	if icon := (navigationEntry{page: -999}).icon(); icon != lucide.Component() {
		t.Fatal("unknown destination must retain the generic fallback")
	}
}

func TestNavigationCatalog(t *testing.T) {
	all := navigationCatalog()
	count := 0
	seen := map[galleryPage]bool{}
	for _, group := range all {
		count += len(group.entries)
		for _, entry := range group.entries {
			if seen[entry.page] {
				t.Fatalf("duplicate navigation destination: %s", entry.label)
			}
			seen[entry.page] = true
		}
	}
	if count != 6+len(controlCatalog)+len(mygoCatalog) {
		t.Fatalf("missing catalog entries: %d", count)
	}
}

func TestNavigationSearchTrigger(t *testing.T) {
	for _, width := range []int{220, 520} {
		g := &Gallery{page: welcomePage}
		tt := ui.NewTester(g.Left, width, 600)
		before, ok := tt.Find("Search components")
		if !ok || before.X < 0 || before.X+before.W > float32(width) || before.Y < 48 {
			t.Fatalf("search bounds: %+v", before)
		}
		_, shortcut := commandShortcut()
		if !tt.HasText(shortcut) {
			t.Fatal("missing shortcut hint")
		}
		if err := tt.Click("Search components"); err != nil {
			t.Fatal(err)
		}
		tt.Type("button")
		if !g.commands.open || !tt.HasText("Welcome") || !tt.HasText("Icon library") || g.page != welcomePage {
			t.Fatal("trigger must open commands without filtering or navigating")
		}
	}
}

func TestChineseNavigationNames(t *testing.T) {
	g := &Gallery{}
	g.chooseLanguage("zh-Hans")
	for _, group := range navigationCatalog() {
		if g.text(group.label) == group.label {
			t.Errorf("untranslated navigation group: %s", group.label)
		}
		for _, entry := range group.entries {
			if entry.label == "Lucide" || entry.label == "Markdown" {
				continue
			}
			if g.text(entry.label) == entry.label {
				t.Errorf("untranslated navigation label: %s", entry.label)
			}
		}
	}
	tt := ui.NewTester(g.Left, 240, 6000)
	if err := tt.Click("时间选择器"); err != nil {
		t.Fatal(err)
	}
	example, ok := controlExampleFor(g.page)
	if !ok || example.name != "TimePicker" {
		t.Fatal("translated navigation selected wrong page")
	}
	g.chooseLanguage("en")
	if example, ok := controlExampleFor(g.page); !ok || example.name != "TimePicker" {
		t.Fatal("language switch changed route")
	}
}
