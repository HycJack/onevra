package gallery

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestPanelPosition(t *testing.T) {
	panel := mygo.Rectangle{Width: 380, Height: 320}
	for _, tt := range []struct {
		name         string
		anchor, area mygo.Rectangle
		x, y         int
	}{
		{"top", mygo.Rectangle{X: 500, Width: 24, Height: 24}, mygo.Rectangle{Y: 24, Width: 1440, Height: 876}, 322, 30},
		{"bottom", mygo.Rectangle{X: 500, Y: 876, Width: 24, Height: 24}, mygo.Rectangle{Width: 1440, Height: 876}, 322, 550},
		{"left edge", mygo.Rectangle{Width: 24, Height: 24}, mygo.Rectangle{Y: 24, Width: 1440, Height: 876}, 0, 30},
		{"right edge", mygo.Rectangle{X: 1416, Width: 24, Height: 24}, mygo.Rectangle{Y: 24, Width: 1440, Height: 876}, 1060, 30},
		{"secondary display", mygo.Rectangle{X: -30, Y: -900, Width: 24, Height: 24}, mygo.Rectangle{X: -1440, Y: -876, Width: 1440, Height: 876}, -380, -870},
		{"small work area", mygo.Rectangle{Width: 24, Height: 24}, mygo.Rectangle{Y: 24, Width: 200, Height: 200}, 0, 24},
	} {
		t.Run(tt.name, func(t *testing.T) {
			x, y := panelPosition(tt.anchor, tt.area, panel)
			if x != tt.x || y != tt.y {
				t.Fatalf("position = (%d, %d), want (%d, %d)", x, y, tt.x, tt.y)
			}
		})
	}
}

func TestTrayIcon(t *testing.T) {
	img, err := png.Decode(bytes.NewReader(trayIcon()))
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 32 || img.Bounds().Dy() != 32 {
		t.Fatalf("unexpected icon bounds: %v", img.Bounds())
	}
	if _, _, _, a := img.At(0, 0).RGBA(); a != 0 {
		t.Fatal("icon background must be transparent")
	}
	if _, _, _, a := img.At(6, 4).RGBA(); a == 0 {
		t.Fatal("icon must contain visible ink")
	}
}

func TestTrayPage(t *testing.T) {
	g := &Gallery{page: trayPage}
	tt := ui.NewTester(g.View, 600, 500)
	if !tt.HasText("Tray") || !tt.HasText("Enable tray") {
		t.Fatal("missing tray page", tt.Texts())
	}
	if g.tray.tray != nil || g.tray.win != nil {
		t.Fatal("rendering must not create native resources")
	}
	if pageFromKey(pageKey(trayPage)) != trayPage {
		t.Fatal("tray route does not round trip")
	}
	g.Close()
	g.Close()
}

func TestNativeTray(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	mygo.RunOnMain(func() {
		g := &Gallery{page: trayPage}
		g.New()
		defer g.Shell.Window().Destroy()
		defer g.Close()
		a := &g.tray
		if err := a.start(g.Shell.Window()); err != nil {
			t.Error(err)
			return
		}
		original := a.tray
		if err := a.start(g.Shell.Window()); err != nil || a.tray != original {
			t.Error("enable should be idempotent")
		}
		if a.win.IsVisible() {
			t.Error("panel should start hidden")
		}
		a.note = "Preserved while hidden"
		a.toggle()
		if !a.win.IsVisible() {
			t.Error("toggle should show the panel")
		}
		a.toggle()
		if a.win.IsVisible() {
			t.Error("toggle should hide the panel")
		}
		a.show()
		a.win.Close()
		if a.win.IsVisible() || a.note != "Preserved while hidden" {
			t.Error("close should hide the panel without discarding state")
		}
		old := a.win
		a.close()
		a.close()
		if !old.IsDestroyed() || a.tray != nil || a.win != nil || a.offQuit != nil {
			t.Error("disable must release native resources and quit subscription")
		}
		if err := a.start(g.Shell.Window()); err != nil {
			t.Error(err)
		}
		if a.note != "Preserved while hidden" {
			t.Error("disable should retain the in-memory note")
		}
		g.Close()
		if a.tray != nil || a.win != nil {
			t.Error("gallery cleanup left tray resources alive")
		}
	})
}
