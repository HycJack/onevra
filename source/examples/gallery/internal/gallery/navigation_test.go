package gallery

import (
	"bytes"
	"math"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func TestNavigationSpacing(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 480)
	first, ok := tt.Find("Overview")
	search, searchOK := tt.Find("Search components")
	if !ok || !searchOK || search.Y != 48 || first.Y != search.Y+search.H+8 {
		t.Fatalf("first group must follow the search field at the shell's safe content edge: %v", first)
	}
	item, _ := tt.Find("Welcome")
	next, _ := tt.Find(newGroupLabel)
	if gap := next.Y - item.Y - item.H; math.Abs(float64(gap-12)) > 0.001 {
		t.Fatalf("group gap = %v, want 12", gap)
	}
}

func TestNavigationGroupOrder(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	// Entries repeat in the New components index, so their order is checked in
	// the render order rather than by the position of a unique label.
	texts := tt.Texts()
	consume := func(label string) {
		t.Helper()
		i := slices.Index(texts, label)
		if i < 0 {
			t.Fatalf("missing or unordered entry: %s", label)
		}
		texts = texts[i+1:]
	}
	for _, group := range navigationCatalog() {
		consume(g.text(group.label))
		for _, entry := range group.entries {
			consume(g.text(entry.label))
		}
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

func TestNavigationLeavesMyGoMarkersOut(t *testing.T) {
	for _, group := range navigationCatalog() {
		for _, entry := range group.entries {
			if !entry.mygo && isMyGoPage(entry.page) {
				t.Fatalf("MyGo page %s is listed outside the MyGo entry", entry.label)
			}
		}
	}
	// Pages behind the MyGo entry need no per-item MyGo label or marker.
	g := &Gallery{page: notificationsPage}
	for _, tt := range []*ui.Tester{ui.NewTester(g.Left, 220, 6000), ui.NewTester(g.View, 900, 6000)} {
		texts := tt.Texts()
		if slices.Contains(texts, "•") {
			t.Fatal("unexpected pending marker")
		}
		if i := slices.Index(texts, "MyGo"); i >= 0 && (i+1 >= len(texts) || texts[i+1] != "Official components") {
			t.Fatalf("unexpected MyGo label after %q", texts[max(i-1, 0)])
		}
	}
}

func TestNavigationSelectionAndHover(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(map[bool]string{false: "light", true: "dark"}[dark], func(t *testing.T) {
			g := &Gallery{}
			tt := ui.NewTester(g.Left, 220, 1600)
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
	purpose := map[galleryPage]bool{}
	hubs := 0
	visit := func(entry navigationEntry) {
		t.Helper()
		if entry.mygo {
			hubs++
			return
		}
		if purpose[entry.page] {
			t.Fatalf("duplicate navigation destination: %s", entry.label)
		}
		purpose[entry.page] = true
	}
	// The MyGo entry hosts every MyGo page once in its in-page list.
	for _, group := range append(all, mygoNavigationGroups()...) {
		for _, entry := range group.entries {
			visit(entry)
		}
	}
	if hubs != 1 {
		t.Fatalf("got %d MyGo entries, want 1", hubs)
	}
	if len(purpose) != 3+len(controlCatalog)+len(mygoCatalog) {
		t.Fatalf("missing catalog entries: %d", len(purpose))
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
		if g.text(group.label) == group.label && group.label != mygoGroupLabel {
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
	if err := tt.Click("日期和时间选择器"); err != nil {
		t.Fatal(err)
	}
	example, ok := controlExampleFor(g.page)
	if !ok || example.name != "Date & time pickers" {
		t.Fatal("translated navigation selected wrong page")
	}
	g.chooseLanguage("en")
	if example, ok := controlExampleFor(g.page); !ok || example.name != "Date & time pickers" {
		t.Fatal("language switch changed route")
	}
}
