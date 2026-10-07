package components_test

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/components"
)

func ExampleNewWebView() {
	mygo.App.WhenReady(func() {
		var web *components.WebView
		show := true
		win := mygo.NewWindow(mygo.WindowOptions{
			Title: "Native + WebView", Width: 900, Height: 700, Hidden: true,
			Content: ui.View(func(c *ui.Context) {
				ui.Column(c).Fill().Gap(12).Padding(16).Children(func() {
					if web == nil {
						return
					}
					if ui.Button(c, "Show / hide browser").Clicked() {
						show = !show
						// Hide before omitting the slot; omission cannot unmount
						// a system view. Returning reuses the same document.
						web.SetVisible(show)
					}
					if show {
						web.View(c).Fill()
					}
				})
			}),
		})
		win.OnClosed(func() { mygo.App.Quit() })
		var err error
		web, err = components.NewWebView(win, mygo.PageOptions{})
		if err != nil {
			log.Print(err) // ErrWebViewUnsupported on other platforms.
			win.Destroy()
			return
		}
		web.Page().LoadHTML("<h1>Hello from the WebView</h1>", "")
		win.Show()
		// Host shutdown closes the browser automatically. To release it sooner,
		// call web.Close() on the UI thread, without closing the host.
	})
	if err := mygo.App.Run(); err != nil {
		log.Print(err)
	}
}
