package main

import "github.com/egoist/mygo"

type shortcutOption struct {
	id, accelerator string
}

// shortcutOptions are tried in turn: the first one free is the app's. Not
// Alt+Space, which Raycast, PowerToys Run and others take: macOS gives a
// shortcut to every app that registers it, so both would open.
var shortcutOptions = []shortcutOption{
	{"altShiftSpace", "Alt+Shift+Space"},
	{"cmdAltShiftSpace", "CmdOrCtrl+Alt+Shift+Space"},
	{"cmdShiftSpace", "CmdOrCtrl+Shift+Space"},
}

const defaultShortcut = "altShiftSpace"

func acceleratorFor(id string) string {
	for _, option := range shortcutOptions {
		if option.id == id {
			return option.accelerator
		}
	}
	return ""
}

func shortcutLabel(id string) string { return keys(acceleratorFor(id)) }

// startShortcut registers the saved shortcut, or the first other one free when
// another app took it.
func (a *application) startShortcut() {
	order := []string{a.config.Shortcut}
	for _, option := range shortcutOptions {
		if option.id != a.config.Shortcut {
			order = append(order, option.id)
		}
	}
	for _, id := range order {
		a.applyShortcut(id)
		if a.shortcut == id {
			if id != a.config.Shortcut {
				a.config.Shortcut = id
				a.changed()
			}
			return
		}
	}
}

// applyShortcut moves the global shortcut to id, keeping the previous one when
// the system refuses it because another app took it.
func (a *application) applyShortcut(id string) {
	accelerator := acceleratorFor(id)
	if accelerator == "" || accelerator == a.accelerator {
		return
	}
	if a.accelerator != "" {
		mygo.GlobalShortcut.Unregister(a.accelerator)
	}
	if err := mygo.GlobalShortcut.Register(accelerator, a.toggle); err != nil {
		a.report("shortcutError", err)
		if a.accelerator != "" {
			_ = mygo.GlobalShortcut.Register(a.accelerator, a.toggle)
		}
		return
	}
	a.shortcut, a.accelerator = id, accelerator
	if a.errorKey == "shortcutError" {
		a.errorKey, a.errorDetail = "", ""
	}
}
