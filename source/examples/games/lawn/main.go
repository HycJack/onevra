// Lawn Defense is a small tower-defense game example on MyGo's game
// plugin, offscreen canvas and audio plugin, framed by Onevra's plain
// Shell and components. MyGo owns the window and the fixed-step loop; the
// game synthesizes its pixel art, sound effects and looping tune in
// memory, so the example ships no assets.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/games/lawn/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

// app owns the round and the window that shows it.
type app struct{ game lawnDemo }

func (a *app) view(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(3)).Children(func() { a.game.view(c) })
}

func main() {
	mygo.App.SetName("Onevra Lawn Defense")
	mygo.App.WhenReady(func() {
		a := &app{}
		// Native windows play the game's effects and music; headless tests keep it silent.
		a.game.output = true
		win := plain.NewShell(plain.ShellOptions{
			Title: locales.In(mygo.App.Locale()).Title(),
			Width: 460, Height: 560, MinWidth: 360, MinHeight: 460,
			StateKey: "lawn",
			Theme:    theme.Onevra,
			Content:  a.view,
		}).Window()
		win.OnClosed(func() {
			a.game.closeMusic()
			mygo.App.Quit()
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
