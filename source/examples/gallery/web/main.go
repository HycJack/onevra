// The website's component preview: one page of the Gallery's demonstrations,
// built for browsers by website/scripts/preview-wasm.ts. It runs natively
// too, for checking a page without the website:
//
//	go run ./examples/gallery/web -page button
//	go run ./examples/gallery/web -manifest
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery"
)

func main() {
	manifest := flag.Bool("manifest", false, "print the component page catalog as JSON")
	page := flag.String("page", "", "component page slug, such as button")
	language := flag.String("lang", "en", "interface language: en or zh-Hans")
	flag.Parse()

	if *manifest {
		data, err := json.MarshalIndent(gallery.PreviewEntries(), "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
		return
	}

	preview := gallery.NewPreview(*page, *language)
	mygo.App.SetName("Onevra components")
	mygo.App.WhenReady(func() {
		if *language == "zh-Hans" {
			mygo.App.SetLocale("zh-Hans")
		} else {
			mygo.App.SetLocale("en")
		}
		window := mygo.NewWindow(mygo.WindowOptions{
			Title:   "Onevra components",
			Width:   960,
			Height:  720,
			Content: ui.View(preview.View),
		})
		installBridge(preview, window)
		window.OnClosed(func() { mygo.App.Quit() })
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
