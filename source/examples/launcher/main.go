package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
)

const appName = "Onevra Launcher"

func main() {
	mygo.App.SetName(appName)
	if !mygo.App.RequestSingleInstanceLock() {
		return // the running launcher shows itself (OnSecondInstance)
	}
	base, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(base, "onevra-examples", "launcher")
	mygo.App.SetPath(mygo.PathUserData, dir)
	a := &application{path: filepath.Join(dir, "config.json")}
	a.config, err = loadConfig(a.path)
	if err != nil {
		a.report("loadError", err)
	}
	// Start in the saved appearance rather than switching after launch.
	mygo.Theme.SetSource(a.themeSource())
	// A launcher lives in the background: no Dock icon or app switcher entry.
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.OnWindowAllClosed(func() {}) // keep running without windows
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.destroyTray() })
	mygo.App.OnActivate(func(bool) {
		if a.ready {
			a.show()
		}
	})
	mygo.App.OnSecondInstance(func([]string, string) {
		if a.ready {
			a.show()
		} else {
			a.pendingShow = true
		}
	})
	mygo.App.WhenReady(a.start)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
