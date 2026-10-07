package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/examples/gallery/internal/gallery"
)

func main() {
	theme := flag.String("theme", "system", "Appearance: system, light, dark")
	flag.Parse()
	if *theme != "system" && *theme != "light" && *theme != "dark" {
		log.Fatal("theme must be system, light, or dark")
	}
	gallery.LogMyGoReview(log.Printf)
	if err := configureTerminalLibrary(); err != nil {
		log.Fatal(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(home, ".oneo", "onevra")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatal(err)
	}
	mygo.App.SetPath(mygo.PathUserData, configDir)
	g := &gallery.Gallery{}
	mygo.App.SetName("Onevra Gallery")
	var override mygo.ThemeSource
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "theme" {
			override = mygo.ThemeSource(*theme)
		}
	})
	g.Configure(filepath.Join(configDir, "config.json"), override)
	mygo.App.WhenReady(func() {
		win := g.NewPersistent(filepath.Join(configDir, "state.json")).Window()
		win.OnClosed(func() { mygo.App.Quit() })
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
