package gallery

import (
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// Template values are session-only; application preferences are persisted.
type settingsState struct {
	open, initialized bool
	view              components.SettingsState
	smart             bool
	channels          []string
}

var notificationRows = []struct{ group, title, description, channel string }{
	{"Account & security", "Security alert", "Unusual activity or important security updates.", "Email"},
	{"Account & security", "New sign-in", "Your account is accessed from a new device.", "Email"},
	{"Account & security", "Password changed", "Your password is updated.", "Email"},
	{"Workspace updates", "Role or permission changed", "Your access level changes.", "In-app"},
	{"Workspace updates", "New team member joined", "Someone joins the workspace.", "In-app"},
	{"Workspace updates", "Workspace settings changed", "Important workspace settings are updated.", "In-app"},
	{"Workspace updates", "Member removed", "Someone is removed from the workspace.", "In-app"},
	{"Messages & notes", "Mentions", "Someone mentions you in a conversation.", "In-app"},
	{"Messages & notes", "Shared notes", "A note is shared with you.", "Email"},
}

func (s *settingsState) initialize() {
	if s.initialized {
		return
	}
	s.initialized, s.smart = true, true
	for _, row := range notificationRows {
		s.channels = append(s.channels, row.channel)
	}
}

// Settings hosts the public settings surface in an application-owned dialog.
func (g *Gallery) Settings(c *ui.Context) {
	g.withLanguage(c, func() { g.settingsLocalized(c) })
}

func (g *Gallery) settingsLocalized(c *ui.Context) {
	if !g.settings.open {
		return
	}
	g.settings.initialize()
	t := c.Theme()
	w, h := c.Size()
	inset := sidebar.HeaderHeight(c)
	width := min(t.Space(200), w*.88, max(0, w-2*inset))
	components.Dialog(c, &g.settings.open, components.DialogOptions{Label: tr(c, "settings"), Inset: inset}, func(panel ui.Element) {
		panel.Width(width).Height(min(t.Space(180), h*.85)).Padding(0).Gap(0).Clip()
		components.Settings(c, &g.settings.view, g.settingsCategories(), components.SettingsOptions{
			Narrow: width < t.Space(170),
			HeaderActions: func(c *ui.Context) {
				if sidebar.IconButton(c, lucide.X(), tr(c, "closeSettings")).Clicked() {
					g.settings.open = false
				}
			},
		})
	})
}

func (g *Gallery) settingsCategories() []components.SettingsCategory {
	notifications := components.SettingsCategory{
		ID: "notifications", Title: g.text("notifications"), Section: g.text("personal"), Icon: lucide.Bell(),
		Description: g.text("templatePreviewChangesStayInThisSessionNoNotificationsAreSen"),
		Groups: []components.SettingsGroup{{
			ID: "smart",
			Items: []components.SettingsItem{{
				ID: "smart", Title: g.text("smartNotifications"), Description: g.text("pauseMobileAlertsWhileYouReActiveOnDesktop"),
				// The description ends with an inline link that toggles the switch.
				DescriptionContent: func(c *ui.Context) {
					ui.Text(c, g.text("pauseMobileAlertsWhileYouReActiveOnDesktop")+" ")
					action := "turnOn"
					if g.settings.smart {
						action = "turnOff"
					}
					if components.Button(c, g.text(action), components.ButtonLink).Clicked() {
						g.settings.smart = !g.settings.smart
					}
				},
				Build: func(c *ui.Context) {
					components.Switch(c, &g.settings.smart, g.text("toggleSmartNotifications"), components.SwitchOptions{HideLabel: true})
				},
			}},
		}},
	}
	for i, row := range notificationRows {
		if notifications.Groups[len(notifications.Groups)-1].ID != row.group {
			group := row.group
			notifications.Groups = append(notifications.Groups, components.SettingsGroup{ID: group, Title: g.text(group), Actions: func(c *ui.Context) {
				if components.Button(c, g.text("resetDefaults"), components.ButtonGhost).Label(g.messages().ResetS(g.text(group))).Shrink(0).Clicked() {
					for i, row := range notificationRows {
						if row.group == group {
							g.settings.channels[i] = row.channel
						}
					}
				}
			}})
		}
		group := &notifications.Groups[len(notifications.Groups)-1]
		group.Items = append(group.Items, components.SettingsItem{
			ID: row.title, Title: g.text(row.title), Description: g.text(row.description),
			Build: func(c *ui.Context) {
				settingsSelect(c, &g.settings.channels[i], g.messages().DeliveryOf(g.text(row.title)), []string{"Email", "In-app", "Off"}).Shrink(0)
			},
		})
	}
	return []components.SettingsCategory{
		{ID: "general", Title: g.text("general"), Section: g.text("application"), Icon: lucide.Settings(),
			Groups: []components.SettingsGroup{{ID: "language", Items: []components.SettingsItem{{
				ID: "language", Title: g.text("language"), Description: g.text("appliesImmediatelyToMenusAndControls"), Keywords: []string{"language", "English", "中文", g.text("followSystem")},
				Build: func(c *ui.Context) {
					choice := g.config.Language
					if settingsSelect(c, &choice, g.text("languagePreference"), []string{"", "en", "zh-Hans"}).Changed() {
						g.chooseLanguage(choice)
						c.Invalidate()
					}
				},
			}}}, g.commandSettingsGroup()}},
		{ID: "appearance", Title: g.text("appearance"), Section: g.text("application"), Icon: lucide.Palette(),
			Groups: []components.SettingsGroup{{ID: "theme", Items: []components.SettingsItem{{
				ID: "theme", Title: g.text("theme"), Description: g.text("appliesImmediatelyAndFollowsTheToolbarThemeMenu"), Keywords: []string{"System", "Light", "Dark", g.text("system"), g.text("light"), g.text("dark")},
				Build: func(c *ui.Context) {
					choice := string(mygo.Theme.Source())
					if settingsSelect(c, &choice, g.text("themePreference"), []string{"system", "light", "dark"}).Changed() {
						g.chooseTheme(mygo.ThemeSource(choice))
					}
				},
			}, {
				ID: "pixel-grid", Title: g.text("pixelGridBackground"), Description: g.text("pixelGridBackgroundDescription"),
				Build: func(c *ui.Context) {
					enabled := !g.config.DisablePixelGrid
					components.Switch(c, &enabled, g.text("togglePixelGridBackground"), components.SwitchOptions{HideLabel: true}).Changed()
					if enabled != !g.config.DisablePixelGrid {
						g.config.DisablePixelGrid = !enabled
						g.savePreferences()
						c.Invalidate()
					}
				},
			}}}},
		},
		notifications,
	}
}

func settingsSelect(c *ui.Context, value *string, label string, choices []string) ui.Element {
	sel := components.Select(c, value)
	sel.Trigger.Label(label).MinWidth(c.Theme().Space(32)).MaxWidthPercent(100).Children(func() { components.Text(c, settingChoiceLabel(c, *value)).SingleLine() })
	sel.Popup(func(panel ui.Element) {
		for _, choice := range choices {
			sel.Item(choice).Label(settingChoiceLabel(c, choice)).Children(func() { components.Text(c, settingChoiceLabel(c, choice)) })
		}
	})
	return sel.Trigger
}

func settingChoiceLabel(c *ui.Context, choice string) string {
	switch choice {
	case "":
		return tr(c, "followSystem")
	case "en":
		return "English"
	case "zh-Hans":
		return "简体中文"
	case "system":
		return tr(c, "system")
	case "light":
		return tr(c, "light")
	case "dark":
		return tr(c, "dark")
	default:
		return tr(c, choice)
	}
}

func (g *Gallery) commandSettingsGroup() components.SettingsGroup {
	return components.SettingsGroup{ID: "commands", Title: g.text("commandPalette"), Items: []components.SettingsItem{
		{ID: "command-abc", Title: g.text("commandAutoAbc"), Description: g.text("commandAutoAbcDescription"), Keywords: []string{"ABC", "input", "Raycast"}, Build: func(c *ui.Context) {
			enabled := !g.config.DisableCommandABC
			components.Switch(c, &enabled, g.text("toggleCommandAutoAbc"), components.SwitchOptions{HideLabel: true}).Disabled(runtime.GOOS != "darwin").Changed()
			if enabled != !g.config.DisableCommandABC {
				g.config.DisableCommandABC = !enabled
				g.savePreferences()
			}
		}},
		{ID: "command-expanded", Title: g.text("commandDefaultExpanded"), Description: g.text("commandDefaultExpandedDescription"), Keywords: []string{"search", "results", "Raycast"}, Build: func(c *ui.Context) {
			expanded := !g.config.CommandExpandOnSearch
			components.Switch(c, &expanded, g.text("toggleCommandDefaultExpanded"), components.SwitchOptions{HideLabel: true}).Changed()
			if expanded != !g.config.CommandExpandOnSearch {
				g.config.CommandExpandOnSearch = !expanded
				g.savePreferences()
			}
		}},
	}}
}
