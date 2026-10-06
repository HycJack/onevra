package main

import (
	"flag"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/internal/gallery"
)

func main() {
	theme := flag.String("theme", "system", "Appearance: system, light, dark")
	flag.Parse()
	if *theme != "system" && *theme != "light" && *theme != "dark" {
		log.Fatal("theme must be system, light, or dark")
	}
	g := &gallery.Gallery{}
	mygo.App.SetName("oneo-go Gallery")
	mygo.Theme.SetSource(mygo.ThemeSource(*theme))
	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title: "oneo-go Gallery", Width: 820, Height: 900,
			MinWidth: 360, MinHeight: 480, Content: ui.View(g.View),
		})
		win.OnClosed(func() { mygo.App.Quit() })
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
