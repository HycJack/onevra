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
	"github.com/egoist/mygo/ui"
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
	win := mygo.NewWindow(mygo.WindowOptions{Title: "oneo-go Gallery", Width: 820, Height: 1000, Content: ui.View(g.View)})
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
		{"light", mygo.ThemeLight, 820},
		{"dark", mygo.ThemeDark, 820},
		{"narrow", mygo.ThemeLight, 360},
	} {
		mygo.Theme.SetSource(mode.theme)
		win.SetSize(mode.width, 1000)
		time.Sleep(400 * time.Millisecond)
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
