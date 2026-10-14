package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func TestCodeEditorExample(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans"} {
		g, s, tt := featureTester("Code editor", "Default", lang, 540)
		if s.code == nil || s.code.DiffBase == nil || s.code.Label != g.text("gallery.code_editor_label") {
			t.Fatalf("editor state: %+v", s.code)
		}
		if err := tt.Click(g.text("gallery.code_editor_label")); err != nil {
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
		want := g.text("gallery.code_editor_position")
		if !tt.HasText(fmt.Sprintf(want, pos.Line+1, pos.Col+1)) {
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
		if added == 0 || deleted == 0 || !tt.HasText(fmt.Sprintf(g.text("gallery.diff_viewer_stats"), added, deleted)) {
			t.Fatalf("%s stats +%d -%d: %q", variant, added, deleted, tt.Texts())
		}
		// Reverting every change restores the old text in the new buffer.
		for len(s.diff.Hunks()) > 0 {
			s.diff.Revert(0)
		}
		tt.Frame()
		if s.diff.New.String() != editorSample || !tt.HasText(fmt.Sprintf(g.text("gallery.diff_viewer_stats"), 0, 0)) {
			t.Fatalf("%s revert left %q", variant, tt.Texts())
		}
	}
}

func TestCanvasExample(t *testing.T) {
	g, s, tt := featureTester("Canvas", "Default", "en", 540)
	if s.value != 0 || s.canvas.smooth == nil || s.canvas.drawn != 0 {
		t.Fatalf("initial canvas state: value %v, drawn %d", s.value, s.canvas.drawn)
	}
	for _, label := range []string{"gallery.canvas_smooth", "gallery.canvas_pixelated", "gallery.canvas_layer"} {
		if _, ok := tt.Find(g.text(label)); !ok {
			t.Fatalf("missing %s", label)
		}
	}
	if err := tt.Click(g.text("gallery.canvas_next_frame")); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(g.text("gallery.canvas_mirror")); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	// The offscreen layer is redrawn for the new frame, mirrored.
	if s.index != 1 || !s.on || s.canvas.drawn != 3 {
		t.Fatalf("frame %d, mirrored %v, layer drawn for %d", s.index, s.on, s.canvas.drawn)
	}
	for range spriteFrames - 1 {
		if err := tt.Click(g.text("gallery.canvas_next_frame")); err != nil {
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
	for range 6 {
		time.Sleep(20 * time.Millisecond)
		tt.Frame()
	}
	tt.KeyUp(ui.KeyD)
	tt.KeyUp(ui.KeyS)
	if d.steps == 0 || d.x <= x || d.y <= y || d.left {
		t.Fatalf("after %d steps: (%v, %v) from (%v, %v), left %v", d.steps, d.x, d.y, x, y, d.left)
	}
	// The player stays inside the world.
	d.x, d.y = gameWorld, -gameWorld
	tt.KeyDown(ui.KeyA)
	time.Sleep(40 * time.Millisecond)
	tt.Frame()
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
	want := map[string][]string{"Code & diffs": {"Code", "Code editor", "Diff viewer"}, "Games": {"Canvas", "Game loop", "Lawn defense"}}
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

func TestNewestPagesHighlighted(t *testing.T) {
	for key := range newestPages {
		if pageKey(pageFromKey(key)) != key {
			t.Fatalf("newest page %q has no route", key)
		}
	}
	for _, lang := range []string{"en", "zh-Hans"} {
		g := &Gallery{language: newTranslator(lang, "en")}
		tt := ui.NewTester(g.Left, 220, 6000)
		badge := g.text("gallery.navigation_new")
		texts, got := tt.Texts(), 0
		for i, text := range texts {
			if text != badge {
				continue
			}
			got++
			// The badge follows its item's label, or the label's upstream marker.
			label := texts[i-1]
			if label == "•" {
				label = texts[i-2]
			}
			found := false
			for key := range newestPages {
				found = found || g.text(sourceMessageIDs[routeLabel(key)]) == label
			}
			if !found {
				t.Fatalf("%s badge follows %q", lang, label)
			}
		}
		if got != len(newestPages) {
			t.Fatalf("%s: %d badges, want %d", lang, got, len(newestPages))
		}
	}
}

func TestNewestBadgeCentered(t *testing.T) {
	for key := range newestPages {
		g := &Gallery{}
		page := pageFromKey(key)
		label := routeLabel(key)
		tt := ui.NewTester(func(c *ui.Context) {
			g.withLanguage(c, func() { g.navigationItem(c, label, nil, page) })
		}, 220, 80)
		text, _ := tt.Find(label)
		badge, ok := tt.Find(g.text("gallery.navigation_new"))
		if !ok {
			t.Fatalf("%s has no badge", key)
		}
		if d := (badge.Y + badge.H/2) - (text.Y + text.H/2); d > .5 || d < -.5 {
			t.Fatalf("%s badge is off center by %v: badge %+v, label %+v", key, d, badge, text)
		}
	}
}

// routeLabel is the source label of a route key such as "mygo/Canvas".
func routeLabel(key string) string {
	_, label, _ := strings.Cut(key, "/")
	return label
}
