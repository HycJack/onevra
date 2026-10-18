//go:build js

package main

import (
	"syscall/js"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery"
)

// installBridge lets the website change the preview without loading the
// WebAssembly again: the frame's page calls onevraPreviewShow on hash
// changes, and the window redraws the selected page.
func installBridge(preview *gallery.Preview, window *mygo.Window) {
	js.Global().Set("onevraPreviewShow", js.FuncOf(func(js.Value, []js.Value) any {
		slug := js.Global().Get("location").Get("hash").String()
		if len(slug) > 0 && slug[0] == '#' {
			slug = slug[1:]
		}
		window.Update(func() { preview.Show(slug) })
		return nil
	}))
}
