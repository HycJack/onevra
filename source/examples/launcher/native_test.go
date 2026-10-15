package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
)

func TestMain(m *testing.M) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "onevra-launcher-")
	if err != nil {
		panic(err)
	}
	mygo.App.SetName(appName + " Test")
	mygo.App.SetPath(mygo.PathUserData, dir)
	// Language choices made by tests apply without being saved for later launches.
	os.Setenv("MYGO_LOCALE", "en")
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

func TestNativeLauncher(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" || os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_SMOKE=1 and ONEVRA_NATIVE_FOREGROUND=1")
	}
	a := &application{config: defaultConfig(), path: filepath.Join(t.TempDir(), "config.yaml")}
	mygo.RunOnMain(func() {
		a.start()
		if a.panel == nil {
			t.Error("startup did not create the panel")
			return
		}
		a.show()
		if !a.panel.IsVisible() {
			t.Error("show did not open the panel")
		}
		a.toggle()
		if a.panel.IsVisible() {
			t.Error("toggle did not hide the panel")
		}
		a.toggle()
		if !a.panel.IsVisible() {
			t.Error("toggle did not show the panel")
		}
		a.apps = []mygo.InstalledApp{{Name: "Calculator", ID: "com.apple.calculator", Path: "/Applications/Calculator.app"}}
		a.query, a.chosen = "calc", 0
		a.filter()
		if len(a.results) != 1 || a.results[0].Name != "Calculator" {
			t.Errorf("filter: %v", a.results)
		}
		a.openSettings()
		if a.settings == nil || !a.settings.IsVisible() {
			t.Error("settings did not open")
		}
	})
	defer mygo.RunOnMain(func() {
		a.quitting = true
		for _, win := range []*mygo.Window{a.panel, a.settings} {
			if win != nil {
				win.Close()
			}
		}
		a.destroyTray()
		mygo.GlobalShortcut.UnregisterAll()
		if len(mygo.Windows()) != 0 {
			t.Error("quit path retained windows")
		}
	})
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"light", "dark"} {
		mygo.RunOnMain(func() {
			a.config.Appearance = mode
			a.changed()
			a.show()
		})
		time.Sleep(350 * time.Millisecond)
		for name, capture := range map[string]*mygo.Window{"launcher": a.panel, "settings": a.settings} {
			if capture == nil {
				continue
			}
			data, err := capture.CapturePage()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, name+"-"+mode+".png"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
}
