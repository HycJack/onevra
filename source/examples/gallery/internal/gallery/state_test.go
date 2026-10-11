package gallery

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/egoist/mygo"
	"go.yaml.in/yaml/v3"
)

func TestGalleryVersionMatchesPackaging(t *testing.T) {
	data, err := os.ReadFile("../../../../packer.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		Version string `yaml:"version"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if config.Version != galleryVersion {
		t.Fatalf("update galleryVersion to match packer.yaml: %q != %q", galleryVersion, config.Version)
	}
}

func TestWelcomeVersionStorage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.yaml")
	state := State{SelectedItem: "controls/Button", SidebarWidth: 310, SidebarCollapsed: true}
	for _, step := range []struct {
		version string
		show    bool
	}{
		{"1.0", true}, {"1.0", false}, {"1.1", true}, {"1.1", false}, {"1.0", true},
	} {
		if got := recordWelcomeVersion(path, &state, step.version); got != step.show {
			t.Fatalf("version %s: show = %v", step.version, got)
		}
		loaded, err := LoadState(path)
		if err != nil || loaded != state || loaded.WelcomeVersion != step.version {
			t.Fatalf("immediate persistence: %+v, %v", loaded, err)
		}
		if loaded.SelectedItem != "controls/Button" || loaded.SidebarWidth != 310 || !loaded.SidebarCollapsed {
			t.Fatal("welcome changed the saved main interface")
		}
		state = loaded
	}
}

func TestStateStorage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "state.yaml")
	if state, err := LoadState(path); err != nil || state != (State{}) {
		t.Fatalf("missing: %+v %v", state, err)
	}
	want := State{Window: &WindowState{-900, 0, 1000, 700}, SidebarWidth: 312.5, SelectedItem: "controls/Button"}
	for _, width := range []float32{312.5, 400} {
		want.SidebarWidth = width
		want.SidebarCollapsed = width == 400
		if err := SaveState(path, want); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil || !strings.Contains(string(data), "sidebar_width:") {
			t.Fatalf("expected YAML state: %s, %v", data, err)
		}
		got, err := LoadState(path)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip: %+v %v", got, err)
		}
	}
	if err := os.WriteFile(path, []byte("sidebar_width: 300\nwindow: ["), 0600); err != nil {
		t.Fatal(err)
	}
	if state, err := LoadState(path); err == nil || state != (State{}) {
		t.Fatalf("corrupt: %+v %v", state, err)
	}
	if err := SaveState(filepath.Join(path, "state.yaml"), want); err == nil {
		t.Fatal("expected write failure")
	}
}

func TestStateSidebarDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.yaml")
	if err := os.WriteFile(path, []byte(`# Partial presentation state
sidebar_width: 300`), 0600); err != nil {
		t.Fatal(err)
	}
	state, err := LoadState(path)
	if err != nil || state.SidebarCollapsed || state.SidebarWidth != 300 {
		t.Fatalf("legacy state: %+v %v", state, err)
	}
}

func TestStateRoutes(t *testing.T) {
	keys := map[string]bool{}
	for page := welcomePage; page < controlPageStart()+galleryPage(len(controlCatalog)); page++ {
		key := pageKey(page)
		if keys[key] || pageFromKey(key) != page {
			t.Fatalf("invalid route %d: %s", page, key)
		}
		keys[key] = true
	}
	if pageFromKey("removed/item") != lucidePage || pageFromKey("") != lucidePage {
		t.Fatal("unknown route did not fall back")
	}
}

func TestNativeState(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_FOREGROUND=1")
	}
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	for _, collapsed := range []bool{true, false} {
		path := filepath.Join(t.TempDir(), "state.yaml")
		var want State
		mygo.RunOnMain(func() {
			g := &Gallery{}
			s := g.NewPersistent(path)
			s.Window().SetSize(940, 680)
			s.Window().SetPosition(100, 100)
			if !s.LeftExpanded() {
				t.Error("missing state should start expanded")
			}
			s.SetLeftWidth(315.5)
			s.SetLeftExpanded(!collapsed)
			g.page = controlPageStart()
			b := s.Window().Bounds()
			want = State{Window: &WindowState{b.X, b.Y, b.Width, b.Height}, SidebarWidth: 315.5, SidebarCollapsed: collapsed, SelectedItem: pageKey(g.page), WelcomeVersion: galleryVersion}
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
			if s.LeftExpanded() != !want.SidebarCollapsed {
				t.Errorf("restored expanded=%v, want %v", s.LeftExpanded(), !want.SidebarCollapsed)
			}
			b := s.Window().Bounds()
			if b.Width != want.Window.Width || b.Height != want.Window.Height || b.X != want.Window.X || b.Y != want.Window.Y || s.LeftWidth() != want.SidebarWidth || pageKey(g.page) != want.SelectedItem {
				t.Errorf("restore: bounds=%+v width=%v route=%s", b, s.LeftWidth(), pageKey(g.page))
			}
		})
	}
}
