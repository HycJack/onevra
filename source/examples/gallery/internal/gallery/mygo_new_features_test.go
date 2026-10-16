package gallery

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func TestCodeEditorExample(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		g, s, tt := featureTester("Code editor", "Editable", lang, 540)
		if s.code == nil || s.code.DiffBase == nil || s.code.Label != g.text("codeEditorLabel") {
			t.Fatalf("editor state: %+v", s.code)
		}
		if err := tt.Click(g.text("codeEditorLabel")); err != nil {
			t.Fatal(err)
		}
		tt.Type("// edited\n")
		tt.Frame()
		if !s.code.Buffer.Modified() || s.code.Buffer.String() == s.code.DiffBase.String() {
			t.Fatal("typing did not edit the buffer")
		}
		if s.code.DiffBase.String() != editorSample {
			t.Fatal("editing changed the diff base")
		}
		pos := s.code.CursorPos()
		if !tt.HasText(g.messages().CodeEditorPosition(pos.Line+1, pos.Col+1)) {
			t.Fatalf("missing caret position %v: %q", pos, tt.Texts())
		}
	}
}

func TestDiffViewerExample(t *testing.T) {
	for _, variant := range []string{"Side by side", "Unified"} {
		g, s, tt := featureTester("Diff viewer", variant, "zh-Hans", 540)
		if want := map[string]editor.DiffLayout{"Side by side": editor.DiffSplit, "Unified": editor.DiffUnified}[variant]; s.diff.Layout != want {
			t.Fatalf("%s layout = %v", variant, s.diff.Layout)
		}
		added, deleted := s.diff.Stats()
		if added == 0 || deleted == 0 || !tt.HasText(g.messages().DiffViewerStats(added, deleted)) {
			t.Fatalf("%s stats +%d -%d: %q", variant, added, deleted, tt.Texts())
		}
		// Reverting every change restores the old text in the new buffer.
		for len(s.diff.Hunks()) > 0 {
			s.diff.Revert(0)
		}
		tt.Frame()
		if s.diff.New.String() != editorSample || !tt.HasText(g.messages().DiffViewerStats(0, 0)) {
			t.Fatalf("%s revert left %q", variant, tt.Texts())
		}
	}
}

func TestCanvasExample(t *testing.T) {
	g, s, tt := featureTester("Canvas", "Default", "en", 540)
	if s.value != 0 || s.canvas.smooth == nil || s.canvas.drawn != 0 {
		t.Fatalf("initial canvas state: value %v, drawn %d", s.value, s.canvas.drawn)
	}
	for _, label := range []string{"canvasSmooth", "canvasPixelated", "canvasLayer"} {
		if _, ok := tt.Find(g.text(label)); !ok {
			t.Fatalf("missing %s", label)
		}
	}
	if err := tt.Click(g.text("canvasNextFrame")); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(g.text("canvasMirror")); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	// The offscreen layer is redrawn for the new frame, mirrored.
	if s.index != 1 || !s.on || s.canvas.drawn != 3 {
		t.Fatalf("frame %d, mirrored %v, layer drawn for %d", s.index, s.on, s.canvas.drawn)
	}
	for range spriteFrames - 1 {
		if err := tt.Click(g.text("canvasNextFrame")); err != nil {
			t.Fatal(err)
		}
	}
	if s.index != 0 {
		t.Fatalf("frames did not wrap: %d", s.index)
	}
}

func TestSpriteSheet(t *testing.T) {
	img := spriteSheet()
	if b := img.Bounds(); b.Dx() != spriteSize*spriteFrames || b.Dy() != spriteSize {
		t.Fatalf("sheet bounds %v", b)
	}
	for f, rows := range spriteRows {
		for y, row := range rows {
			if len(row) != spriteSize {
				t.Fatalf("frame %d row %d has %d pixels", f, y, len(row))
			}
			for _, p := range []byte(row) {
				if _, ok := spritePalette[p]; !ok && p != '.' {
					t.Fatalf("frame %d uses unknown pixel %q", f, p)
				}
			}
		}
	}
}

func TestGameExampleMovesWithHeldKeys(t *testing.T) {
	_, s, tt := featureTester("Game loop", "Default", "en", 540)
	d := &s.game
	if d.gamepads {
		t.Fatal("headless views must not read gamepads")
	}
	x, y := d.x, d.y
	tt.KeyDown(ui.KeyD)
	tt.KeyDown(ui.KeyS)
	tt.Advance(120 * time.Millisecond)
	tt.KeyUp(ui.KeyD)
	tt.KeyUp(ui.KeyS)
	if d.steps == 0 || d.x <= x || d.y <= y || d.left {
		t.Fatalf("after %d steps: (%v, %v) from (%v, %v), left %v", d.steps, d.x, d.y, x, y, d.left)
	}
	// The player stays inside the world.
	d.x, d.y = gameWorld, -gameWorld
	tt.KeyDown(ui.KeyA)
	tt.Advance(40 * time.Millisecond)
	tt.KeyUp(ui.KeyA)
	if d.x > gameWorld-spriteSize/2 || d.y < spriteSize/2 || !d.left {
		t.Fatalf("player left the world: (%v, %v), left %v", d.x, d.y, d.left)
	}
}

func TestNewMyGoReviewEntries(t *testing.T) {
	entries := map[string]mygoReviewEntry{}
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		entries[entry.Name] = entry
	}
	for _, name := range []string{"Code editor", "Diff viewer"} {
		entry := entries[name]
		if entry.Status != mygoPrimitive || entry.Source != "ui/editor" || !entry.HasExample || entry.Reason == "" {
			t.Fatalf("%s: %+v", name, entry)
		}
	}
	for _, name := range []string{"Canvas", "Game loop"} {
		if entry := entries[name]; entry.Status != mygoPrimitive || entry.Reason == "" {
			t.Fatalf("%s: %+v", name, entry)
		}
	}
}

func TestCodeAndGameGroups(t *testing.T) {
	want := map[string][]string{"Code & diffs": {"Code editor", "Diff viewer"}, "Games": {"Canvas", "Shaders", "Game loop", "3D renderer"}}
	for _, group := range navigationCatalog() {
		if entries, ok := want[group.label]; ok {
			var got []string
			for _, entry := range group.entries {
				got = append(got, entry.label)
			}
			if fmt.Sprint(got) != fmt.Sprint(entries) {
				t.Fatalf("%s = %v, want %v", group.label, got, entries)
			}
			delete(want, group.label)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing groups: %v", want)
	}
}

func TestNewComponentsGroup(t *testing.T) {
	groups := navigationCatalog()
	if len(groups) < 2 || groups[0].label != "Overview" || groups[1].label != newGroupLabel || !groups[1].index {
		t.Fatalf("New components group must follow Overview: %v", groups)
	}
	want := map[galleryPage]bool{}
	for key := range newPages {
		if pageKey(pageFromKey(key)) != key {
			t.Fatalf("new page %q has no route", key)
		}
		want[pageFromKey(key)] = true
	}
	if len(groups[1].entries) != len(want) {
		t.Fatalf("New components group has %d entries, want %d", len(groups[1].entries), len(want))
	}
	// The index repeats pages that stay in their purpose groups for browsing.
	purpose := map[galleryPage]bool{}
	for _, group := range groups {
		if group.index {
			continue
		}
		for _, entry := range group.entries {
			purpose[entry.page] = true
		}
	}
	for _, entry := range groups[1].entries {
		if !want[entry.page] {
			t.Fatalf("unexpected New components entry: %s", entry.label)
		}
		if !purpose[entry.page] {
			t.Fatalf("%s left its purpose group", entry.label)
		}
		delete(want, entry.page)
	}
	if len(want) != 0 {
		t.Fatalf("missing New components entries: %v", want)
	}
}

func TestNewComponentsGroupRenders(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		g := &Gallery{language: newTranslator(lang, "en")}
		tt := ui.NewTester(g.Left, 220, 6000)
		index, ok := tt.Find(g.text(newGroupLabel))
		if !ok {
			t.Fatalf("%s: missing New components group", lang)
		}
		lucide, ok := tt.Find(g.text("Lucide"))
		if !ok || index.Y >= lucide.Y {
			t.Fatalf("%s: New components group is not above the icon library", lang)
		}
		for _, entry := range newNavigationEntries() {
			if !tt.HasText(g.text(entry.label)) {
				t.Fatalf("%s: missing %s", lang, entry.label)
			}
		}
	}
}
