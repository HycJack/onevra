package gallery

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/shell"
)

// Ordinary tests never open windows.
func TestMain(m *testing.M) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "onevra-smoke-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mygo.App.SetName("Onevra smoke")
	mygo.App.SetPath(mygo.PathUserData, dir)
	// Tests own separate windows; only the end of m.Run should quit the app.
	mygo.App.OnWindowAllClosed(func() {})
	result := make(chan int, 1)
	mygo.App.WhenReady(func() { go func() { result <- m.Run(); mygo.App.Quit() }() })
	if err := mygo.App.Run(); err != nil {
		os.RemoveAll(dir)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := <-result
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestNativeGallery(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	g := &Gallery{}
	var appShell *shell.Shell
	mygo.RunOnMain(func() { appShell = g.New() })
	win := appShell.Window()
	defer win.Destroy()
	if win.Page() != nil || win.Title() != "Onevra Gallery" {
		t.Fatal("expected native content and caller title")
	}
	if width, height := win.Size(); width != 820 || height != 600 {
		t.Errorf("gallery default size = %dx%d, want 820x600", width, height)
	}
	// Only native captures relax the default minimum to exercise narrow layouts.
	mygo.RunOnMain(func() { win.SetMinimumSize(360, 480) })
	var blank *mygo.Window
	mygo.RunOnMain(func() { blank = shell.NewShell(shell.ShellOptions{}).Window() })
	if blank.Page() != nil {
		t.Error("nil content created a web page")
	}
	if width, height := blank.Size(); width != 820 || height != 600 {
		t.Errorf("default size = %dx%d", width, height)
	}
	blank.Destroy()
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []struct {
		name     string
		theme    mygo.ThemeSource
		width    int
		settings bool
		icons    bool
		review   bool
	}{
		{"light", mygo.ThemeLight, 820, false, false, false},
		{"dark", mygo.ThemeDark, 820, false, false, false},
		{"settings", mygo.ThemeLight, 820, true, false, false},
		{"restored", mygo.ThemeDark, 820, false, false, false},
		{"narrow", mygo.ThemeLight, 360, false, false, false},
		{"icons-light", mygo.ThemeLight, 820, false, true, false},
		{"icons-dark", mygo.ThemeDark, 820, false, true, false},
		{"icons-narrow", mygo.ThemeLight, 360, false, true, false},
		{"status-light", mygo.ThemeLight, 820, false, false, true},
		{"status-dark", mygo.ThemeDark, 820, false, false, true},
		{"status-narrow", mygo.ThemeLight, 360, false, false, true},
	} {
		mygo.Theme.SetSource(mode.theme)
		mygo.RunOnMain(func() {
			appShell.SetPage(nil)
		})
		time.Sleep(400 * time.Millisecond)
		var before mygo.Rectangle
		mygo.RunOnMain(func() {
			g.page = shellPage
			if mode.icons {
				g.page = lucidePage
			}
			if mode.review {
				g.page = componentStatusPage
			}
			win.SetSize(mode.width, 900)
			before = win.Bounds()
			if mode.settings {
				appShell.SetPage(g.Settings)
			}
		})
		time.Sleep(400 * time.Millisecond)
		mygo.RunOnMain(func() {
			after := win.Bounds()
			wantWidth := mode.width
			if after.Width != wantWidth || after.X != before.X || after.Y != before.Y || after.Height != before.Height {
				t.Errorf("%s: bounds = %+v, expected width %d and unchanged origin/height (%+v)", mode.name, after, wantWidth, before)
			}
		})
		data, err := win.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		if img.Bounds().Dx() < mode.width {
			t.Fatal("capture smaller than window")
		}
		path := filepath.Join(dir, mode.name+".png")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		t.Log(path)
	}
}
