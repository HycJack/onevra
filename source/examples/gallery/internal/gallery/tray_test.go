package gallery

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestTrayIcon(t *testing.T) {
	for _, template := range []bool{true, false} {
		img, err := png.Decode(bytes.NewReader(trayIcon(template)))
		if err != nil {
			t.Fatal(err)
		}
		if img.Bounds().Dx() != 32 || img.Bounds().Dy() != 32 {
			t.Fatalf("unexpected icon bounds: %v", img.Bounds())
		}
		if _, _, _, a := img.At(0, 0).RGBA(); a != 0 {
			t.Fatal("icon background must be transparent")
		}
		// The tall bar of the artwork sits left of center, the tiles right.
		r, g, b, a := img.At(8, 16).RGBA()
		if a == 0 {
			t.Fatal("icon must contain visible ink")
		}
		if template != (r == 0 && g == 0 && b == 0) {
			t.Fatalf("template=%v drew color %d,%d,%d", template, r, g, b)
		}
		if _, _, _, a := img.At(12, 16).RGBA(); a != 0 {
			t.Fatal("gap between the bar and the tiles must stay transparent")
		}
	}
}

func TestTrayPanel(t *testing.T) {
	shown := 0
	a := &trayDemo{showOwner: func() { shown++ }}
	tt := ui.NewTester(a.view, 300, 124)
	if !tt.HasText("Onevra Gallery") || !tt.HasText("More actions are in the tray menu.") {
		t.Fatal("missing panel text", tt.Texts())
	}
	if err := tt.Click("Show Gallery"); err != nil || shown != 1 {
		t.Fatalf("Show Gallery did not reach its owner: %v", err)
	}
}

func TestDesktopIntegrationsNeedShell(t *testing.T) {
	g := &Gallery{}
	ui.NewTester(g.View, 600, 500)
	g.EnableDesktopIntegrations()
	if g.tray.tray != nil || g.dock.menu != nil {
		t.Fatal("headless gallery must not create native resources")
	}
	for _, key := range []string{"desktop/tray", "desktop/dock-menu"} {
		if pageKey(pageFromKey(key)) == key {
			t.Fatalf("removed route %s must fall back", key)
		}
	}
	g.Close()
	g.Close()
}

func TestNativeTray(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_FOREGROUND=1")
	}
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	mygo.RunOnMain(func() {
		g := &Gallery{}
		g.New()
		defer g.Shell.Window().Destroy()
		defer g.Close()
		a := &g.tray
		if err := a.start(g.showMainWindow); err != nil {
			t.Error(err)
			return
		}
		original := a.tray
		if err := a.start(g.showMainWindow); err != nil || a.tray != original {
			t.Error("start should be idempotent")
		}
		if a.panel != nil {
			t.Error("panel should not open with the tray")
		}
		menu := a.menu
		g.chooseLanguage("zh-Hans")
		if a.menu != menu || menu.ItemByID("panel").Label != g.text("openTrayPanel") || menu.ItemByID("quit").Label != "退出" {
			t.Error("language switch must update tray menu without recreating its resources")
		}
		g.chooseLanguage("en")
		if menu.ItemByID("panel").Label != "Open tray panel" {
			t.Error("English tray menu was not restored")
		}
		a.toggle()
		if a.panel == nil {
			t.Error("toggle should open the panel")
		}
		a.toggle()
		if a.panel != nil {
			t.Error("toggle should close the panel")
		}
		g.Shell.Window().Hide()
		a.show()
		item := menu.ItemByID("gallery")
		item.Click(item, nil)
		if a.panel != nil || !g.Shell.Window().IsVisible() {
			t.Error("Show Gallery should close the panel and show the main window")
		}
		a.show()
		old := a.panel
		a.close()
		a.close()
		if !old.IsDestroyed() || a.tray != nil || a.panel != nil || a.offQuit != nil {
			t.Error("close must release native resources and quit subscription")
		}
		if err := a.start(g.showMainWindow); err != nil {
			t.Error(err)
		}
		g.Close()
		if a.tray != nil || a.panel != nil {
			t.Error("gallery cleanup left tray resources alive")
		}
	})
}
