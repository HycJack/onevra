package gallery

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

func TestNativeMyGo(t *testing.T) {
	captureCatalog(t, false)
}

func TestNativeControls(t *testing.T) {
	captureCatalog(t, true)
}

func captureCatalog(t *testing.T, controls bool) {
	t.Helper()
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
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
	var appShell *sidebar.Shell
	mygo.RunOnMain(func() { appShell = g.newWindow(nativeBackground()); g.enterGallery() })
	win := appShell.Window()
	defer win.Destroy()
	// Only native captures relax the default minimum to exercise narrow layouts.
	mygo.RunOnMain(func() { win.SetMinimumSize(360, 480) })
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
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
		mygo.RunOnMain(func() { win.SetSize(mode.width, 900); appShell.SetLeftExpanded(mode.width > 360) })
		time.Sleep(400 * time.Millisecond)
		for i, example := range catalog {
			mygo.RunOnMain(func() {
				g.page = start + galleryPage(i)
				if controls && example.name == "Popover" {
					g.controlState("Popover", "Default").dialogOpen = true
				}
				if controls && example.name == "Dialog" {
					g.controlState("Dialog", "Default").dialogOpen = true
				}
				if controls && example.name == "AlertDialog" {
					g.controlState("AlertDialog", "Confirmation").dialogOpen = true
				}
				if controls && example.name == "Tree" {
					s := g.controlState("Tree", "Tree")
					s.on, s.choice = true, "main.go"
					s = g.controlState("Tree", "TreeVirtualized")
					s.tree.Open.Add("source")
					s.index = 1
				}
				if controls && example.name == "Table" {
					g.controlState("Table", "Table").choice = "main.go"
					g.controlState("Table", "TableVirtualized").index = 1
				}
				win.Invalidate()
			})
			if example.name == "Terminal" {
				var state *demoState
				mygo.RunOnMain(func() { state = g.mygoState("Terminal", "Default", true) })
				if _, err := waitTerminal(t, &state.terminal); err != nil {
					t.Fatal(err)
				}
				win.Invalidate()
			}
			time.Sleep(240 * time.Millisecond)
			data, err := win.CapturePage()
			if err != nil {
				t.Fatalf("%s/%s: %v", mode.name, example.name, err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() < mode.width {
				t.Fatalf("invalid capture for %s/%s: %v", mode.name, example.name, err)
			}
			if controls && runtime.GOOS == "darwin" && mode.name != "narrow" && (example.name == "List" || example.name == "Tree" || example.name == "Table") {
				// Content captures omit the native material, but their alpha
				// must still expose it below an otherwise empty sidebar pixel.
				_, _, _, alpha := img.At(4, img.Bounds().Dy()/2).RGBA()
				if alpha == 0xffff {
					t.Fatalf("%s/%s made the sidebar opaque", mode.name, example.name)
				}
			}
			path := filepath.Join(dir, prefix+mode.name+"-"+example.name+".png")
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			if example.name == "WebView" {
				// Window captures leave a hole where the system web view shows.
				captureWebViewPage(t, g, filepath.Join(dir, prefix+mode.name+"-WebView-page.png"))
			}
		}
	}
	t.Logf("Captured %d pages in light, dark and narrow layouts in %s", len(catalog), dir)
}

// captureWebViewPage waits for the demonstration page to paint more than its
// background, then writes the web view's own capture to path.
func captureWebViewPage(t *testing.T, g *Gallery, path string) {
	t.Helper()
	var d *webViewDemo
	mygo.RunOnMain(func() { d = &g.mygoState("WebView", "Default", true).webView })
	if d.err != nil || d.view == nil {
		t.Fatalf("web view not created: %v", d.err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		data, err := d.view.CapturePage()
		if err == nil {
			if img, err := png.Decode(bytes.NewReader(data)); err == nil && !uniformImage(img) {
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("web view page did not paint: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func uniformImage(img image.Image) bool {
	b := img.Bounds()
	first := img.At(b.Min.X, b.Min.Y)
	for y := b.Min.Y; y < b.Max.Y; y += 4 {
		for x := b.Min.X; x < b.Max.X; x += 4 {
			if img.At(x, y) != first {
				return false
			}
		}
	}
	return true
}

// TestNativeWebViewUnderModal checks that gallery modal layers hide the web
// view: glass cannot see its page, so it would otherwise leave a hole.
func TestNativeWebViewUnderModal(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	g := &Gallery{}
	var appShell *sidebar.Shell
	mygo.RunOnMain(func() { appShell = g.newWindow(nativeBackground()); g.enterGallery() })
	win := appShell.Window()
	defer win.Destroy()
	mygo.Theme.SetSource(mygo.ThemeLight)
	mygo.RunOnMain(func() {
		win.SetSize(820, 600)
		for i, example := range mygoCatalog {
			if example.name == "WebView" {
				g.page = mygoPageStart + galleryPage(i)
			}
		}
		win.Invalidate()
	})
	time.Sleep(800 * time.Millisecond)
	for _, layer := range []string{"settings", "commands"} {
		mygo.RunOnMain(func() {
			g.settings.open, g.commands.open = layer == "settings", layer == "commands"
			win.Invalidate()
		})
		time.Sleep(600 * time.Millisecond)
		data, err := win.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		// The content area right of the sidebar is otherwise opaque.
		b := img.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y += 4 {
			for x := b.Min.X + b.Dx()*3/10; x < b.Max.X; x += 4 {
				if _, _, _, a := img.At(x, y).RGBA(); a != 0xffff {
					t.Fatalf("%s: transparent content at %d,%d under the modal", layer, x, y)
				}
			}
		}
	}
	mygo.RunOnMain(func() { g.settings.open, g.commands.open = false, false })
}
