package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/egoist/mygo"
)

func TestMain(m *testing.M) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "onevra-"+appKind+"-")
	if err != nil {
		panic(err)
	}
	mygo.App.SetName(appName + " Test")
	mygo.App.SetPath(mygo.PathUserData, dir)
	mygo.App.OnWindowAllClosed(func() {})
	result := make(chan int, 1)
	mygo.App.WhenReady(func() { go func() { result <- m.Run(); mygo.App.Quit() }() })
	if err := mygo.App.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := <-result
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestNativeLifecycle(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_SMOKE=1 and ONEVRA_NATIVE_FOREGROUND=1")
	}
	a := &application{config: defaultConfig(), path: filepath.Join(t.TempDir(), "config.yaml")}
	var capture *mygo.Window
	mygo.RunOnMain(func() {
		a.start()
		if appKind == "menubar" {
			if a.main != nil || a.panel != nil || a.welcome == nil || !a.welcome.IsVisible() {
				t.Error("first launch must show only welcome")
			}
			a.finishWelcome()
			if !a.config.Welcomed || a.welcome.IsVisible() {
				t.Error("welcome did not complete")
			}
			a.reopen()
			if a.panel == nil || !a.panel.IsVisible() {
				t.Error("reopen did not show panel")
			}
		} else {
			if a.main == nil || !a.main.IsVisible() {
				t.Error("startup did not show main window")
			}
			if runtime.GOOS == "darwin" || appKind == "hybrid" {
				a.main.Close()
			} else {
				a.main.Hide()
			}
			if a.main.IsVisible() {
				t.Error("close did not hide main window")
			}
			a.reopen()
			if !a.main.IsVisible() {
				t.Error("reopen did not restore main window")
			}
		}
		a.note = "Retained draft"
		if appKind != "desktop" {
			if a.tray == nil {
				t.Error("tray unavailable")
				return
			}
			a.showPanel()
			a.panel.Close()
			if a.panel.IsVisible() {
				t.Error("panel close did not hide")
			}
			a.primaryClick()
			if !a.panel.IsVisible() {
				t.Error("left click did not open panel")
			}
			a.config.Panel = false
			a.changed()
			if a.panel.IsVisible() {
				t.Error("disabling panel did not hide it")
			}
			a.primaryClick()
		} else {
			a.openSettings()
		}
		if appKind == "menubar" {
			capture = a.settings
		} else {
			capture = a.main
			if !a.settingsOpen {
				t.Error("settings dialog not open")
			}
		}
		a.config.Language = "zh-Hans"
		a.changed()
		if a.note != "Retained draft" {
			t.Error("language change lost draft")
		}
		if appKind != "desktop" && a.trayMenu.ItemByID("settings").Label != "设置…" {
			t.Error("tray menu not translated")
		}
		// Quit must let every close pass through, including a visible settings dialog.
	})
	defer mygo.RunOnMain(func() {
		a.quitting = true
		for _, win := range []*mygo.Window{a.main, a.panel, a.settings, a.welcome} {
			if win != nil {
				win.Close()
			}
		}
		a.destroyTray()
		if len(mygo.Windows()) != 0 {
			t.Error("quit path retained windows")
		}
	})
	if capture == nil {
		t.Fatal("missing settings host")
	}
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"light", "dark"} {
		mygo.RunOnMain(func() { a.config.Appearance = mode; a.changed() })
		time.Sleep(350 * time.Millisecond)
		data, err := capture.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, appKind+"-settings-"+mode+".png"), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNativeBackgroundStartup(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_SMOKE=1 and ONEVRA_NATIVE_FOREGROUND=1")
	}
	if appKind == "desktop" {
		t.Skip("desktop always starts with its main window")
	}
	a := &application{config: defaultConfig()}
	a.config.Silent, a.config.Welcomed, a.config.Panel = true, true, false
	mygo.RunOnMain(func() {
		a.start()
		defer func() {
			a.quitting = true
			for _, win := range []*mygo.Window{a.main, a.panel, a.settings, a.welcome} {
				if win != nil {
					win.Close()
				}
			}
			a.destroyTray()
		}()
		if a.tray == nil {
			t.Error("tray unavailable")
			return
		}
		if a.main != nil || a.panel != nil || a.welcome != nil || a.settings != nil {
			t.Error("background startup created a window")
		}
		a.reopen()
		if appKind == "hybrid" {
			if a.main == nil || !a.main.IsVisible() {
				t.Error("reopen did not override silent startup")
			}
			a.config.CloseAction = "both"
			a.changed()
			a.main.Close()
			if a.main.IsVisible() || a.tray == nil {
				t.Error("close lost tray or kept window")
			}
			a.showMain()
			a.main.Minimize()
			a.reopen()
			if a.main.IsMinimized() || !a.main.IsVisible() {
				t.Error("reopen did not restore minimized window")
			}
		} else {
			if a.main != nil || a.panel != nil || a.settings == nil || !a.settings.IsVisible() {
				t.Error("panel-free reopen must open settings")
			}
			first := a.settings
			a.openSettings()
			if a.settings != first {
				t.Error("settings opened a duplicate window")
			}
		}
		// Losing the tray must leave an accessible fallback, never an invisible app.
		a.destroyTray()
		a.openSettings()
		if appKind == "hybrid" && (!a.main.IsVisible() || !a.settingsOpen) {
			t.Error("missing fallback dialog")
		}
		if appKind == "menubar" && !a.settings.IsVisible() {
			t.Error("missing fallback window")
		}
	})
}
