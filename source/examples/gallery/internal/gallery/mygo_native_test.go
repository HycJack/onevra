package gallery

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	oneoui "github.com/oneo-me/oneo-go/ui"
)

func TestNativeMyGo(t *testing.T) {
	captureCatalog(t, false)
}

func TestNativeControls(t *testing.T) {
	captureCatalog(t, true)
}

func captureCatalog(t *testing.T, controls bool) {
	t.Helper()
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	catalog, start, prefix := mygoCatalog, mygoPageStart, "mygo-"
	if controls {
		catalog, start, prefix = nil, controlPageStart(), "controls-"
		for _, example := range controlCatalog {
			catalog = append(catalog, mygoExample{name: example.name})
		}
	}
	g := &Gallery{}
	var shell *oneoui.Shell
	mygo.RunOnMain(func() { shell = g.New() })
	win := shell.Window()
	defer win.Destroy()
	dir := os.Getenv("ONEO_GO_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []struct {
		name  string
		theme mygo.ThemeSource
		width int
	}{
		{"light", mygo.ThemeLight, 820}, {"dark", mygo.ThemeDark, 820}, {"narrow", mygo.ThemeLight, 360},
	} {
		mygo.Theme.SetSource(mode.theme)
		mygo.RunOnMain(func() { win.SetSize(mode.width, 900); shell.SetLeftExpanded(mode.width > 360) })
		time.Sleep(400 * time.Millisecond)
		for i, example := range catalog {
			mygo.RunOnMain(func() { g.page = start + galleryPage(i); win.Invalidate() })
			if example.name == "Terminal" {
				var state *demoState
				mygo.RunOnMain(func() { state = g.mygoState("Terminal", "Default", true) })
				state.terminal.start(newPreviewTerminal)
				if _, err := waitTerminal(t, &state.terminal); err != nil {
					t.Fatal(err)
				}
				win.Invalidate()
			}
			time.Sleep(60 * time.Millisecond)
			data, err := win.CapturePage()
			if err != nil {
				t.Fatalf("%s/%s: %v", mode.name, example.name, err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() < mode.width {
				t.Fatalf("invalid capture for %s/%s: %v", mode.name, example.name, err)
			}
			path := filepath.Join(dir, prefix+mode.name+"-"+example.name+".png")
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Logf("Captured %d pages in light, dark and narrow layouts in %s", len(catalog), dir)
}
