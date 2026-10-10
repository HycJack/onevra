package gallery

import (
	"runtime"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/plugins/command"
)

type commandState struct {
	open    bool
	palette command.State
}

func commandShortcut() (ui.Modifiers, string) {
	if runtime.GOOS == "darwin" {
		return ui.Cmd, "⌘E"
	}
	return ui.Ctrl, "Ctrl+E"
}

func (g *Gallery) openCommands() {
	g.commands.palette.Query = ""
	g.commands.palette.ActiveID = ""
	g.commands.open = true
}

// Commands is an application-owned integration of the public command plugin.
// Shortcuts are scoped to this window and wait behind other modal surfaces.
func (g *Gallery) Commands(c *ui.Context) {
	mods, _ := commandShortcut()
	if c.Shortcut(mods, ui.KeyE) {
		g.openCommands()
	}
	if !g.commands.open {
		command.CommandDialog(c, &g.commands.open, &g.commands.palette, nil, command.Options{}, components.DialogOptions{})
		return
	}
	items := []command.Item{{ID: "settings", Title: tr(c, "gallery.settings"), Group: tr(c, "gallery.actions_navigation"), Icon: lucide.Settings(), Keywords: []string{"settings", "preferences"}}}
	for _, group := range navigationCatalog() {
		for _, entry := range group.entries {
			items = append(items, command.Item{
				ID: pageKey(entry.page), Title: tr(c, entry.label), Group: tr(c, group.label),
				Icon: entry.icon(), Keywords: []string{entry.label, group.label},
			})
		}
	}
	ui.Box(c).Key("command-palette").Absolute().Children(func() {
		inputSource := ""
		if runtime.GOOS == "darwin" && g.Shell != nil && !g.config.DisableCommandABC {
			inputSource = command.InputSourceABC
		}
		command.CommandDialog(c, &g.commands.open, &g.commands.palette, items, command.Options{
			InputSource:    inputSource,
			ExpandOnSearch: g.config.CommandExpandOnSearch,
			SearchLabel:    tr(c, "gallery.command_palette"), Placeholder: tr(c, "gallery.search_commands"),
			ResultsLabel: tr(c, "gallery.commands"), EmptyText: tr(c, "gallery.no_matching_commands"),
			OnActivate: func(id string) {
				g.commands.open = false
				g.enterGallery()
				if id == "settings" {
					g.settings.open = true
				} else {
					g.page = pageFromKey(id)
				}
			},
		}, components.DialogOptions{Label: tr(c, "gallery.command_palette")})
	})
}
