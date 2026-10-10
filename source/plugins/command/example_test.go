package command_test

import (
	"runtime"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/plugins/command"
)

func ExampleCommandDialog() {
	var state command.State
	var open bool
	var destination string
	inputSource := ""
	if runtime.GOOS == "darwin" {
		inputSource = command.InputSourceABC
	}
	items := []command.Item{{ID: "settings", Title: "Settings", Keywords: []string{"preferences"}}}
	view := func(c *ui.Context) {
		if c.Shortcut(ui.Cmd, ui.KeyE) {
			open = true
		}
		command.CommandDialog(c, &open, &state, items, command.Options{
			// Omit this to show results immediately, even with an empty query.
			InputSource:    inputSource,
			ExpandOnSearch: true,
			OnActivate:     func(id string) { destination = id; open = false },
		}, components.DialogOptions{})
	}
	// Also call state.ReleaseInputSource on the UI thread before destroying
	// the host window. Normal dialog dismissal restores the source automatically.
	_ = view // Supply to a MyGo window or Onevra Shell content builder.
	_ = destination
}
