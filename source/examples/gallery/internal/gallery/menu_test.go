package gallery

import (
	"os"
	"runtime"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestApplicationMenu(t *testing.T) {
	g := &Gallery{page: shellPage}
	menu := g.ApplicationMenu()
	if runtime.GOOS != "darwin" {
		if menu != nil {
			t.Fatal("non-macOS menu must remain unchanged")
		}
		return
	}
	items := menu.Items()
	roles := []mygo.MenuRole{mygo.RoleAppMenu, mygo.RoleFileMenu, mygo.RoleEditMenu, mygo.RoleViewMenu, mygo.RoleWindowMenu}
	if len(items) != len(roles) {
		t.Fatal("missing standard menus")
	}
	for i, role := range roles {
		if items[i].Role != role {
			t.Fatalf("menu %d: %s", i, items[i].Role)
		}
	}
	settings := menu.ItemByID("settings")
	if settings == nil || settings.Label != "Settings…" || settings.Accelerator != "Cmd+," || settings.Click == nil {
		t.Fatal("missing Settings menu or accelerator")
	}
	if items[0].Submenu[2] != settings {
		t.Fatal("Settings must follow About and its separator")
	}
	for _, role := range []mygo.MenuRole{mygo.RoleAbout, mygo.RoleServices, mygo.RoleHide, mygo.RoleHideOthers, mygo.RoleUnhide, mygo.RoleQuit} {
		found := false
		for _, item := range items[0].Submenu {
			found = found || item.Role == role
		}
		if !found {
			t.Fatalf("missing standard application role %s", role)
		}
	}
	g.settings.view.Query = "Theme"
	for range 2 {
		settings.Click(settings, nil)
	}
	if !g.settings.open || g.page != shellPage || g.settings.view.Query != "Theme" {
		t.Fatal("settings action reset navigation or state")
	}
	tt := ui.NewTester(g.View, 820, 600)
	tt.Frame()
	if !tt.HasText("Theme") {
		t.Fatal("settings dialog did not render", tt.Texts())
	}
}

func TestNativeApplicationMenu(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires a macOS desktop session")
	}
	mygo.RunOnMain(func() {
		original := mygo.App.Menu()
		defer mygo.App.SetMenu(original)
		g := &Gallery{page: shellPage}
		g.New()
		defer g.Shell.Window().Destroy()
		defer g.Close()
		menu := g.ApplicationMenu()
		mygo.App.SetMenu(menu)
		if mygo.App.Menu() != menu {
			t.Error("menu installation failed")
			return
		}
		g.Shell.Window().Hide()
		settings := menu.ItemByID("settings")
		settings.Click(settings, nil)
		if !g.settings.open || !g.Shell.Window().IsVisible() || g.page != shellPage {
			t.Error("settings did not show its owner while preserving navigation")
		}
	})
}
