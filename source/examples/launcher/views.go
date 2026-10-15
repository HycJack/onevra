package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/theme"
)

// scope installs the shortcuts every Onevra window of the app shares.
func (a *application) scope(c *ui.Context, build func()) {
	c.OnShortcut(ui.Cmd, ui.KeyQ, mygo.App.Quit)
	c.OnShortcut(ui.Cmd, ui.KeyComma, a.openSettings)
	build()
}

// panelView is the launcher itself: a search field, the matching apps and a
// footer that names the keys. It runs on the hidden panel window and reads the
// application-owned state the shortcut, the tray and the goroutines change.
func (a *application) panelView(c *ui.Context) {
	theme.Apply(c)
	t := c.Theme()
	root := c.Root()
	if c.Vibrancy() {
		root.Background(ui.Transparent) // the material shows through
	} else {
		root.Background(t.Background)
	}
	// An Escape left to the panel closes it, which hides it: with a query,
	// Escape clears it first.
	if a.query != "" && c.Shortcut(0, ui.KeyEscape) {
		a.query = ""
		a.filter()
		a.input.Focus()
	}
	if c.Shortcut(0, ui.KeyDown) {
		a.choose(a.chosen + 1)
	}
	if c.Shortcut(0, ui.KeyUp) {
		a.choose(a.chosen - 1)
	}
	if c.Shortcut(0, ui.KeyEnter) {
		a.open()
	}
	c.OnShortcut(ui.Cmd, ui.KeyQ, mygo.App.Quit)
	a.prefetchIcons()

	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Padding(16, 18).Gap(12).AlignItems(ui.Center).Children(func() {
			input := components.TextBox(c, &a.query, components.TextBoxOptions{LeadingIcon: lucide.Search()})
			input.Bind(&a.input).Grow(1).FontSize(22).Label(a.text("search")).Placeholder(a.text("search"))
			if input.Changed() {
				a.filter()
			}
		})
		divider(c)
		a.list.Selected = &a.chosen
		a.list.Key = func(i int) any { return a.results[i].Path }
		a.list.Label = func(i int) string { return a.results[i].Name }
		a.list.OnSubmit = func(int) { a.open() }
		message := ""
		switch {
		case a.loading && len(a.results) == 0:
			message = a.text("loading")
		case len(a.results) == 0:
			message = a.text("noApps")
		}
		if message != "" {
			ui.Box(c).Grow(1).Center().Children(func() {
				components.Text(c, message).TextColor(t.TextMuted)
			})
		} else {
			list := components.List(c, &a.list, len(a.results), func(i int) { a.row(c, i) }).Grow(1).Padding(6)
			if list.Submitted() {
				a.open()
			}
		}
		divider(c)
		ui.Row(c).Padding(8, 16).Gap(16).Children(func() {
			components.Text(c, a.format("showOrHide", map[string]any{"shortcut": keys(a.accelerator)})).FontSize(11).TextColor(t.TextMuted).Grow(1).SingleLine()
			components.Text(c, a.text("keyHints")).FontSize(11).TextColor(t.TextMuted).SingleLine()
		})
	})
}

// row shows an app with its icon, its name and where it is. The List draws the
// chosen row's face; the row owns only its content.
func (a *application) row(c *ui.Context, i int) {
	app := a.results[i]
	t := c.Theme()
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(12).Padding(6, 10).Children(func() {
		a.appIcon(c, app)
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			components.Text(c, app.Name).SingleLine()
			components.Text(c, where(app)).FontSize(11).TextColor(t.TextMuted).SingleLine()
		})
	})
}

func (a *application) feedback(c *ui.Context) {
	if a.errorKey == "" {
		return
	}
	components.Text(c, a.text(a.errorKey)).Wrap().TextColor(c.Theme().Danger)
	if a.errorDetail != "" {
		components.Text(c, a.errorDetail).Wrap().FontSize(11).TextColor(c.Theme().TextMuted)
	}
}

func (a *application) settingsView(c *ui.Context, width float32) {
	items := []components.SettingsItem{
		{ID: "language", Title: a.text("language"), Keywords: []string{"English", "中文"}, Build: func(c *ui.Context) {
			a.selectSetting(c, &a.config.Language, "language", []string{"system", "en", "zh-Hans"})
		}},
		{ID: "appearance", Title: a.text("appearance"), Build: func(c *ui.Context) {
			a.selectSetting(c, &a.config.Appearance, "appearance", []string{"system", "light", "dark"})
		}},
		{ID: "shortcut", Title: a.text("shortcut"), Description: a.text("shortcutHint"), Build: a.selectShortcut},
		{ID: "login", Title: a.text("openAtLogin"), Description: a.text("openAtLoginHint"), Build: func(c *ui.Context) {
			if components.Switch(c, &a.login, a.text("openAtLogin"), components.SwitchOptions{HideLabel: true}).Changed() {
				if err := mygo.App.SetOpenAtLogin(a.login); err != nil {
					a.report("loginError", err)
				}
				a.refresh() // the tray menu shows the same choice
			}
		}},
	}
	ui.Column(c).Fill().Children(func() {
		components.Settings(c, &a.settingsState, []components.SettingsCategory{{ID: "general", Title: a.text("general"), Groups: []components.SettingsGroup{{ID: "preferences", Items: items}}}}, components.SettingsOptions{
			Narrow: width < 620,
			HeaderActions: func(c *ui.Context) {
				if components.IconButton(c, lucide.X(), a.text("closeSettings"), components.IconButtonPlain).Clicked() {
					a.settings.Hide()
				}
			},
		})
		a.feedback(c)
	})
}

// selectSetting binds a saved choice and persists it when the user picks
// another value.
func (a *application) selectSetting(c *ui.Context, value *string, label string, choices []string) {
	before := *value
	selectBox := components.Select(c, value)
	selectBox.Trigger.Label(a.text(label)).MinWidth(160).Children(func() {
		components.Text(c, a.text(optionKey(*value))).SingleLine()
	})
	selectBox.Popup(func(ui.Element) {
		for _, choice := range choices {
			selectBox.Item(choice).Label(a.text(optionKey(choice))).Children(func() {
				components.Text(c, a.text(optionKey(choice)))
			})
		}
	})
	if before != *value {
		a.changed()
	}
}

// selectShortcut moves the global shortcut and keeps the previous one when the
// system refuses the new one.
func (a *application) selectShortcut(c *ui.Context) {
	choice := a.shortcut
	selectBox := components.Select(c, &choice)
	selectBox.Trigger.Label(a.text("shortcut")).MinWidth(160).Children(func() {
		components.Text(c, shortcutLabel(choice)).SingleLine()
	})
	selectBox.Popup(func(ui.Element) {
		for _, option := range shortcutOptions {
			selectBox.Item(option.id).Label(shortcutLabel(option.id)).Children(func() {
				components.Text(c, shortcutLabel(option.id))
			})
		}
	})
	if choice != a.shortcut {
		a.applyShortcut(choice)
		if a.shortcut == choice {
			a.config.Shortcut = choice
			a.changed()
		}
	}
}
