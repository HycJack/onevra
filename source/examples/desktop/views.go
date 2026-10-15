package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
	"runtime"
)

// themeSource selects light or dark for the whole app, including native
// window chrome; config values match MyGo's theme sources.
func (a *application) themeSource() mygo.ThemeSource { return mygo.ThemeSource(a.config.Appearance) }

func (a *application) scope(c *ui.Context, build func()) {
	c.OnShortcut(ui.Cmd, ui.KeyQ, mygo.App.Quit)
	c.OnShortcut(ui.Cmd, ui.KeyComma, a.openSettings)
	build()
}

func (a *application) feedback(c *ui.Context) {
	if a.errorKey != "" {
		components.Text(c, a.text(a.errorKey)).Wrap().TextColor(c.Theme().Danger)
		components.Text(c, a.errorDetail).Wrap().FontSize(11).TextColor(c.Theme().TextMuted)
	}
}

func (a *application) mainView(c *ui.Context) {
	a.scope(c, func() {
		c.OnShortcut(ui.Cmd, ui.KeyW, a.main.Close)
		t := c.Theme()
		ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
			switch a.mainTab {
			case 1:
				components.Text(c, a.text("title")).FontSize(26).Bold()
				components.Text(c, a.text("description")).Wrap().TextColor(t.TextMuted)
			default:
				components.TextBox(c, &a.note, components.TextBoxOptions{Multiline: true}).Grow(1).Label(a.text("note")).Placeholder(a.text("placeholder"))
				components.Text(c, a.text("sessionNote")).Wrap().FontSize(12).TextColor(t.TextMuted)
			}
			a.feedback(c)
		})
		components.Dialog(c, &a.settingsOpen, components.DialogOptions{Label: a.text("settings")}, func(panel ui.Element) {
			w, h := c.Size()
			inset := plain.HeaderHeight(c)
			width := min(float32(680), w-2*inset)
			panel.Width(width).Height(min(float32(480), h-2*inset)).Padding(0).Gap(0)
			a.settingsView(c, width)
		})
	})
}

// headerLeading puts the view switcher at the left of the title bar.
func (a *application) headerLeading(c *ui.Context) {
	components.Tabs(c, &a.mainTab, components.TabsSegmented, a.text("note"), a.text("about"))
}

// headerActions puts the window's settings action at the right of the title bar.
func (a *application) headerActions(c *ui.Context) {
	if components.IconButton(c, lucide.Settings(), a.text("settings"), components.IconButtonToolbar).Clicked() {
		a.openSettings()
	}
}

func (a *application) panelView(c *ui.Context) {
	theme.Apply(c)
	a.scope(c, func() {
		if c.Shortcut(0, ui.KeyEscape) || c.Shortcut(ui.Cmd, ui.KeyW) {
			a.hidePanel()
		}
		// The popover material shows on macOS; elsewhere the flyout is
		// transparent and the panel draws its own surface.
		c.Root().Background(ui.Transparent)
		panel := ui.Column(c).Fill().Padding(16).Gap(12)
		if !c.Vibrancy() {
			panel.Background(c.Theme().Surface)
			switch runtime.GOOS {
			case "darwin":
				panel.Radius(10)
			case "windows":
			default:
				panel.Radius(10).Border(1, c.Theme().Border)
			}
		}
		panel.Children(func() {
			a.panelToolbar(c)
			switch a.panelTab {
			case 1:
				components.Text(c, a.text("title")).FontSize(18).Bold()
				components.Text(c, a.text("description")).Wrap().FontSize(12).TextColor(c.Theme().TextMuted)
			default:
				components.TextBox(c, &a.note, components.TextBoxOptions{Multiline: true}).Grow(1).Label(a.text("note")).Placeholder(a.text("placeholder")).AutoFocus()
				components.Text(c, a.text("sessionNote")).Wrap().FontSize(11)
			}
		})
	})
}

// panelToolbar groups a segmented view switcher on the left and icon actions on
// the right of one toolbar row.
func (a *application) panelToolbar(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
		components.Tabs(c, &a.panelTab, components.TabsSegmented, a.text("note"), a.text("about"))
		ui.Spacer(c)
		if components.IconButton(c, lucide.Copy(), a.text("copy"), components.IconButtonToolbar).Clicked() {
			c.WriteClipboard(a.note)
		}
		if components.IconButton(c, lucide.Settings(), a.text("settings"), components.IconButtonToolbar).Clicked() {
			a.openSettings()
		}
	})
}

func (a *application) welcomeView(c *ui.Context) {
	a.scope(c, func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
			ui.Column(c).FillWidth().Padding(24).Gap(16).Children(func() {
				components.Text(c, a.text("welcome")).FontSize(24).Bold()
				key := "welcomeBody"
				if runtime.GOOS == "linux" {
					key = "welcomeLinux"
				}
				components.Text(c, a.text(key)).Wrap()
				a.feedback(c)
				components.Button(c, a.text("continue"), components.ButtonPrimary).OnClick(a.finishWelcome)
			})
		})
	})
}

func (a *application) selectSetting(c *ui.Context, value *string, label string, choices []string) {
	before := *value
	selectBox := components.Select(c, value)
	selectBox.Trigger.Label(a.text(label)).MinWidth(160).Children(func() { components.Text(c, a.text(optionKey(*value))).SingleLine() })
	selectBox.Popup(func(ui.Element) {
		for _, choice := range choices {
			selectBox.Item(choice).Label(a.text(optionKey(choice))).Children(func() { components.Text(c, a.text(optionKey(choice))) })
		}
	})
	if before != *value {
		a.changed()
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
	}
	if appKind != "menubar" {
		choices := []string{"quit"}
		if runtime.GOOS == "darwin" {
			choices = []string{"dock", "quit"}
		}
		if len(choices) > 1 {
			items = append(items, components.SettingsItem{ID: "close", Title: a.text("closeAction"), Build: func(c *ui.Context) { a.selectSetting(c, &a.config.CloseAction, "closeAction", choices) }})
		}
	}
	if appKind != "desktop" {
		items = append(items, components.SettingsItem{ID: "panel", Title: a.text("panel"), Description: a.text("panelHint"), Build: func(c *ui.Context) {
			if components.Switch(c, &a.config.Panel, a.text("panel"), components.SwitchOptions{HideLabel: true}).Changed() {
				a.changed()
			}
		}})
	}
	ui.Column(c).Fill().Children(func() {
		components.Settings(c, &a.settingsState, []components.SettingsCategory{{ID: "general", Title: a.text("general"), Groups: []components.SettingsGroup{{ID: "preferences", Items: items}}}}, components.SettingsOptions{
			Narrow: width < 620, HideSearch: true, HideSidebar: true,
			HeaderActions: func(c *ui.Context) {
				if components.IconButton(c, lucide.X(), a.text("closeSettings"), components.IconButtonPlain).Clicked() {
					if appKind == "menubar" {
						a.settings.Close()
					} else {
						a.settingsOpen = false
					}
				}
			},
		})
		a.feedback(c)
	})
}
