package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
)

const appKind = "menubar"
const appName = "Onevra Menu Bar"

func main() {
	mygo.App.SetName(appName)
	if !mygo.App.RequestSingleInstanceLock() {
		return
	}
	base, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(base, "onevra-examples", appKind)
	mygo.App.SetPath(mygo.PathUserData, dir)
	a := &application{path: filepath.Join(dir, "config.yaml")}
	a.config, err = loadConfig(a.path)
	if err != nil {
		a.report("loadError", err)
	}
	// Start in the saved appearance rather than switching after launch.
	mygo.Theme.SetSource(a.themeSource())
	if appKind == "menubar" || appKind == "hybrid" && a.config.Silent && a.config.CloseAction != "both" {
		mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	}
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.destroyTray() })
	mygo.App.OnActivate(func(bool) {
		if a.ready {
			a.reopen()
		}
	})
	mygo.App.OnSecondInstance(func([]string, string) {
		if a.ready {
			a.reopen()
		} else {
			a.pendingOpen = true
		}
	})
	mygo.App.WhenReady(a.start)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
