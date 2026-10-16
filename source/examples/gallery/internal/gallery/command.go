package gallery

import (
	"runtime"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/icons/lucide"
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

// Commands is an application-owned integration of the public command component.
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
	// Trailing content marks the application-level command among pages.
	items := []command.Item{{ID: "settings", Title: tr(c, "settings"), Group: tr(c, "actionsNavigation"), Icon: lucide.Settings(), Keywords: []string{"settings", "preferences"},
		Trailing: func(c *ui.Context) { components.Badge(c, tr(c, "application")).AlignSelf(ui.Center) }}}
	for _, group := range navigationCatalog() {
		for _, entry := range group.entries {
			if entry.mygo {
				// The MyGo entry stands for the pages of its in-page list.
				for _, purpose := range mygoNavigationGroups() {
					for _, page := range purpose.entries {
						items = append(items, command.Item{
							ID: pageKey(page.page), Title: tr(c, page.label), Group: tr(c, group.label),
							Icon: page.icon(), Keywords: []string{page.label, purpose.label, group.label},
						})
					}
				}
				continue
			}
			keywords := []string{entry.label, group.label}
			if example, ok := controlExampleFor(entry.page); ok {
				keywords = append(keywords, pageComponents(example)...)
			}
			items = append(items, command.Item{
				ID: pageKey(entry.page), Title: tr(c, entry.label), Group: tr(c, group.label),
				Icon: entry.icon(), Keywords: keywords,
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
			SearchLabel:    tr(c, "commandPalette"), Placeholder: tr(c, "searchCommands"),
			ResultsLabel: tr(c, "commands"), EmptyText: tr(c, "noMatchingCommands"),
			OnActivate: func(id string) {
				g.commands.open = false
				if id == "settings" {
					g.settings.open = true
				} else {
					g.page = pageFromKey(id)
				}
			},
		}, components.DialogOptions{Label: tr(c, "commandPalette")})
	})
}
