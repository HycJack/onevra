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

func TestFlyoutNoteAndHeadlessLifecycle(t *testing.T) {
	g := &Gallery{}
	d := &g.mygoState("Flyout", "Default", true).flyout
	d.show(ui.Rect{}, false)
	if d.window != nil {
		t.Fatal("headless example created a native window")
	}
	tt := ui.NewTester(d.view, 300, 180)
	tt.Type("Retained note")
	g.Close()
	if d.note != "Retained note" {
		t.Fatalf("note lost: %q", d.note)
	}
	ui.NewTester(d.view, 300, 180)
	if d.note != "Retained note" {
		t.Fatal("rebuilding reset note")
	}
}

func TestNativeFlyout(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	g := &Gallery{}
	mygo.RunOnMain(func() { g.New(); g.enterGallery() })
	parent := g.Shell.Window()
	defer parent.Destroy()
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	var d *flyoutDemo
	mygo.RunOnMain(func() { d = &g.mygoState("Flyout", "Default", true).flyout; d.note = "Retained note" })
	for _, dark := range []bool{false, true} {
		source := mygo.ThemeLight
		if dark {
			source = mygo.ThemeDark
		}
		mygo.Theme.SetSource(source)
		for _, popover := range []bool{false, true} {
			var w *mygo.Window
			mygo.RunOnMain(func() { d.show(ui.Rect{X: 280, Y: 80, W: 100, H: 30}, popover); w = d.window })
			if w == nil || !w.IsFlyout() {
				t.Fatal("expected native flyout")
			}
			time.Sleep(350 * time.Millisecond)
			data, err := w.CapturePage()
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() < 300 {
				t.Fatalf("invalid capture: %v", err)
			}
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("flyout-dark-%v-popover-%v.png", dark, popover)), data, 0600); err != nil {
				t.Fatal(err)
			}
			mygo.RunOnMain(func() { w.Close() })
			mygo.RunOnMain(func() {
				if d.window != nil || d.note != "Retained note" {
					t.Error("close did not release window and retain note")
				}
			})
		}
	}
	mygo.RunOnMain(func() {
		d.show(ui.Rect{X: 280, Y: 80, W: 100, H: 30}, false)
		parent.Destroy()
		if d.window != nil {
			t.Error("parent close retained flyout")
		}
	})
}
