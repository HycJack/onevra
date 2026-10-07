package gallery

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/egoist/mygo"
)

func TestStateStorage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "state.json")
	if state, err := LoadState(path); err != nil || state != (State{}) {
		t.Fatalf("missing: %+v %v", state, err)
	}
	want := State{Window: &WindowState{-900, 0, 1000, 700}, SidebarWidth: 312.5, SelectedItem: "controls/Button"}
	for _, width := range []float32{312.5, 400} {
		want.SidebarWidth = width
		if err := SaveState(path, want); err != nil {
			t.Fatal(err)
		}
		got, err := LoadState(path)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip: %+v %v", got, err)
		}
	}
	if err := os.WriteFile(path, []byte(`{"sidebar_width":300,"window":`), 0600); err != nil {
		t.Fatal(err)
	}
	if state, err := LoadState(path); err == nil || state != (State{}) {
		t.Fatalf("corrupt: %+v %v", state, err)
	}
	if err := SaveState(filepath.Join(path, "state.json"), want); err == nil {
		t.Fatal("expected write failure")
	}
}

func TestStateRoutes(t *testing.T) {
	keys := map[string]bool{}
	for page := lucidePage; page < controlPageStart()+galleryPage(len(controlCatalog)); page++ {
		key := pageKey(page)
		if keys[key] || pageFromKey(key) != page {
			t.Fatalf("invalid route %d: %s", page, key)
		}
		keys[key] = true
	}
	if pageFromKey("removed/item") != lucidePage {
		t.Fatal("unknown route did not fall back")
	}
}

func TestNativeState(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	path := filepath.Join(t.TempDir(), "state.json")
	var want State
	mygo.RunOnMain(func() {
		g := &Gallery{}
		s := g.NewPersistent(path)
		s.Window().SetSize(940, 680)
		s.Window().SetPosition(100, 100)
		s.SetLeftWidth(315.5)
		g.page = controlPageStart()
		b := s.Window().Bounds()
		want = State{Window: &WindowState{b.X, b.Y, b.Width, b.Height}, SidebarWidth: 315.5, SelectedItem: pageKey(g.page)}
		s.Window().Close()
	})
	got, err := LoadState(path)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("saved: %+v want %+v err %v", got, want, err)
	}
	mygo.RunOnMain(func() {
		g := &Gallery{}
		s := g.NewPersistent(path)
		defer s.Window().Destroy()
		b := s.Window().Bounds()
		if b.Width != want.Window.Width || b.Height != want.Window.Height || b.X != want.Window.X || b.Y != want.Window.Y || s.LeftWidth() != want.SidebarWidth || pageKey(g.page) != want.SelectedItem {
			t.Errorf("restore: bounds=%+v width=%v route=%s", b, s.LeftWidth(), pageKey(g.page))
		}
	})
}
