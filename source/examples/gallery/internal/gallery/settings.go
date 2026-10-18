package gallery

import (
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// Template values are session-only; application preferences are persisted.
type settingsState struct {
	open, initialized       bool
	view                    components.SettingsState
	hideSearch, hideSidebar bool
	smart                   bool
	channels                []string
	quietStart, quietEnd    time.Time
	displayName, bio        string
	profileColor            ui.Color
	scale                   float64
	density                 string
	fontSize                float64
	tabWidth                int
	lineNumbers, wordWrap   bool
	minimap                 bool
	cacheGB                 float64
	downloads               string
	upToDate                bool
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
	s.quietStart = time.Date(2000, 1, 1, 22, 0, 0, 0, time.Local)
	s.quietEnd = time.Date(2000, 1, 1, 7, 0, 0, 0, time.Local)
	s.displayName, s.profileColor = "Alex Morgan", ui.RGB(99, 102, 241)
	s.scale, s.density = 100, "comfortable"
	s.fontSize, s.tabWidth, s.lineNumbers, s.wordWrap = 13, 4, true, true
	s.cacheGB, s.downloads = 2.4, "~/Downloads/Onevra"
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
			Narrow: width < t.Space(170), HideSearch: g.settings.hideSearch, HideSidebar: g.settings.hideSidebar,
			Close: func() { g.settings.open = false },
			// The personal section's heading adds a member from the Account
			// page; the integrations section's has nothing to add yet.
			SectionActions: func(c *ui.Context, section string) {
				switch section {
				case g.text("personal"):
					if components.IconButton(c, lucide.Plus(), g.text("addMember"), components.IconButtonCompact).Clicked() {
						g.settings.view = components.SettingsState{Category: "account"}
					}
				case g.text("integrations"):
					components.IconButton(c, lucide.Plus(), g.text("addIntegration"), components.IconButtonCompact).Disabled(true)
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
			}, {
				ID: "quiet-hours", Title: g.text("quietHours"), Description: g.text("quietHoursDescription"), Keywords: []string{"do not disturb", "勿扰"},
				Leading: func(c *ui.Context) { settingsGlyph(c, lucide.Moon()) },
				Build: func(c *ui.Context) {
					ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(c.Theme().Space(2)).Children(func() {
						components.TimePicker(c, &g.settings.quietStart).Label(g.text("quietHoursStart"))
						components.Text(c, "–").TextColor(c.Theme().TextMuted)
						components.TimePicker(c, &g.settings.quietEnd).Label(g.text("quietHoursEnd"))
					})
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
				ID: "language", Title: g.text("language"), Description: g.text("appliesImmediatelyToMenusAndControls"), Keywords: []string{"language", "English", "中文", "日本語", "한국어", g.text("followSystem")},
				Build: func(c *ui.Context) {
					choice := g.config.Language
					if settingsSelect(c, &choice, g.text("languagePreference"), append([]string{""}, languages...)).Changed() {
						g.chooseLanguage(choice)
						c.Invalidate()
					}
				},
			}}}, g.commandSettingsGroup(), g.layoutSettingsGroup()}},
		{ID: "appearance", Title: g.text("appearance"), Section: g.text("application"), Icon: lucide.Palette(),
			Groups: []components.SettingsGroup{{ID: "theme", Items: []components.SettingsItem{{
				ID: "theme", Title: g.text("theme"), Description: g.text("appliesImmediatelyAndFollowsTheToolbarThemeMenu"), Keywords: []string{"System", "Light", "Dark", g.text("system"), g.text("light"), g.text("dark")},
				Build: func(c *ui.Context) {
					choice := string(mygo.Theme.Source())
					if settingsSelect(c, &choice, g.text("themePreference"), []string{"system", "light", "dark"}).Changed() {
						g.chooseTheme(mygo.ThemeSource(choice))
					}
				},
			}}}},
		},
		// Template categories stay together under their own section.
		g.accountSettings(),
		notifications,
		g.layoutSettings(),
		g.editorSettings(),
		g.storageSettings(),
		g.aboutSettings(),
		// An empty section shows a note in the sidebar instead of a page.
		{ID: "integrations", Title: g.text("integrations"), Section: g.text("integrations"), Placeholder: func(c *ui.Context) {
			components.Text(c, g.text("noIntegrationsYet")).FontSize(c.Theme().Rem(.875)).TextColor(c.Theme().TextMuted).Wrap()
		}},
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
	case "en", "zh-Hans", "zh-Hant", "ja", "ko":
		// Each language is named in itself.
		return i18n.Autonym(choice)
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
