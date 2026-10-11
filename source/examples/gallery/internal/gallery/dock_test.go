package gallery

import (
	"os"
	"runtime"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestDockMenuPage(t *testing.T) {
	g := &Gallery{page: dockPage}
	tt := ui.NewTester(g.View, 600, 500)
	if !tt.HasText("Dock menu") {
		t.Fatal("missing Dock menu page", tt.Texts())
	}
	if runtime.GOOS == "darwin" && !tt.HasText("Enable Dock menu") {
		t.Fatal("missing enable action")
	}
	if runtime.GOOS != "darwin" && tt.HasText("Enable Dock menu") {
		t.Fatal("unsupported platform should show explanation, not an enable action")
	}
	if g.dock.menu != nil {
		t.Fatal("rendering must not install a native menu")
	}
	if pageFromKey(pageKey(dockPage)) != dockPage {
		t.Fatal("Dock route does not round trip")
	}
	g.enableDockMenu()
	if g.dock.menu != nil {
		t.Fatal("headless gallery must not install a native menu")
	}
	g.Close()
	g.Close()
}

func TestDockMenuActions(t *testing.T) {
	g := &Gallery{page: dockPage}
	menu := g.dockMenu()
	for _, id := range []string{"show-gallery", "settings", "welcome"} {
		item := menu.ItemByID(id)
		if item == nil || item.Click == nil {
			t.Fatalf("missing action %s", id)
		}
		// Dock callbacks may supply a nil or unrelated focused window. The
		// actions use their Gallery owner rather than the callback window.
		item.Click(item, nil)
		if id == "settings" && (!g.settings.open || g.page != dockPage) {
			t.Fatal("settings action should preserve navigation")
		}
		if id == "welcome" && (g.settings.open || g.page != welcomePage) {
			t.Fatal("welcome action should dismiss settings and select Welcome")
		}
	}
}

func TestNativeDockMenu(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_FOREGROUND=1")
	}
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || runtime.GOOS != "darwin" {
		t.Skip("requires a macOS desktop session")
	}
	mygo.RunOnMain(func() {
		original := mygo.App.Dock.Menu()
		defer mygo.App.Dock.SetMenu(original)
		previous := mygo.NewMenu([]*mygo.MenuItem{{Label: "Previous menu"}})
		mygo.App.Dock.SetMenu(previous)
		g := &Gallery{page: dockPage}
		g.New()
		defer g.Shell.Window().Destroy()
		defer g.Close()
		g.enableDockMenu()
		installed := g.dock.menu
		if installed == nil || mygo.App.Dock.Menu() != installed {
			t.Error("Dock menu was not installed")
			return
		}
		g.enableDockMenu()
		if g.dock.menu != installed {
			t.Error("enabling twice should preserve menu ownership")
		}
		g.chooseLanguage("zh-Hans")
		if mygo.App.Dock.Menu() != installed || installed.ItemByID("settings").Label != "打开设置" {
			t.Error("language switch must update the installed menu in place")
		}
		g.chooseLanguage("en")
		if installed.ItemByID("settings").Label != "Open Settings" {
			t.Error("English menu was not restored")
		}
		g.Shell.Window().Hide()
		item := installed.ItemByID("show-gallery")
		item.Click(item, nil)
		if !g.Shell.Window().IsVisible() {
			t.Error("Show Gallery did not show its owner window")
		}
		g.dock.close()
		g.dock.close()
		if mygo.App.Dock.Menu() != previous || g.dock.menu != nil {
			t.Error("disable did not restore the previous menu")
		}
		g.enableDockMenu()
		g.Close()
		if mygo.App.Dock.Menu() != previous {
			t.Error("gallery cleanup did not restore the previous menu")
		}
		g.enableDockMenu()
		other := mygo.NewMenu([]*mygo.MenuItem{{Label: "Another owner"}})
		mygo.App.Dock.SetMenu(other)
		g.chooseLanguage("zh-Hans")
		if mygo.App.Dock.Menu() != other {
			t.Error("language switch replaced another owner's menu")
		}
		g.Close()
		if mygo.App.Dock.Menu() != other {
			t.Error("cleanup replaced another owner's menu")
		}
	})
}
