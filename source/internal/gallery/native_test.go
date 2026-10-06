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
	"github.com/oneo-me/oneo-go/components"
)

// Ordinary tests never open windows.
func TestMain(m *testing.M) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "oneo-go-smoke-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mygo.App.SetName("oneo-go smoke")
	mygo.App.SetPath(mygo.PathUserData, dir)
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
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	g := &Gallery{}
	var shell *components.Shell
	mygo.RunOnMain(func() { shell = g.New() })
	win := shell.Window()
	defer win.Destroy()
	if win.Page() != nil || win.Title() != "oneo-go Gallery" {
		t.Fatal("expected native content and caller title")
	}
	var blank *mygo.Window
	mygo.RunOnMain(func() { blank = components.NewShell(components.ShellOptions{}).Window() })
	if blank.Page() != nil {
		t.Error("nil content created a web page")
	}
	if width, height := blank.Size(); width != 800 || height != 600 {
		t.Errorf("default size = %dx%d", width, height)
	}
	blank.Destroy()
	dir := os.Getenv("ONEO_GO_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []struct {
		name            string
		theme           mygo.ThemeSource
		width           int
		right, settings bool
	}{
		{"light", mygo.ThemeLight, 820, false, false},
		{"dark", mygo.ThemeDark, 820, false, false},
		{"inspector", mygo.ThemeLight, 820, true, false},
		{"settings", mygo.ThemeLight, 820, true, true},
		{"restored", mygo.ThemeDark, 820, true, false},
		{"narrow", mygo.ThemeLight, 360, false, false},
	} {
		mygo.Theme.SetSource(mode.theme)
		mygo.RunOnMain(func() {
			shell.SetPage(nil)
			shell.SetRightExpanded(false)
		})
		time.Sleep(400 * time.Millisecond)
		var before mygo.Rectangle
		mygo.RunOnMain(func() {
			win.SetSize(mode.width, 900)
			before = win.Bounds()
			shell.SetRightExpanded(mode.right)
			if mode.settings {
				shell.SetPage(g.Settings)
			}
		})
		time.Sleep(400 * time.Millisecond)
		mygo.RunOnMain(func() {
			after := win.Bounds()
			wantWidth := mode.width
			if mode.right && !mode.settings {
				wantWidth += 280
			}
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
