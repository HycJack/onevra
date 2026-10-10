package gallery

import (
	"fmt"
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
	components.Dialog(c, &g.settings.open, components.DialogOptions{Label: tr(c, "gallery.settings"), Inset: inset}, func(panel ui.Element) {
		panel.Width(width).Height(min(t.Space(180), h*.85)).Padding(0).Gap(0).Clip()
		components.Settings(c, &g.settings.view, g.settingsCategories(), components.SettingsOptions{
			Narrow: width < t.Space(170),
			HeaderActions: func(c *ui.Context) {
				if sidebar.IconButton(c, lucide.X(), tr(c, "gallery.close_settings")).Clicked() {
					g.settings.open = false
				}
			},
		})
	})
}

func (g *Gallery) settingsCategories() []components.SettingsCategory {
	notifications := components.SettingsCategory{
		ID: "notifications", Title: g.text("gallery.notifications"), Section: g.text("gallery.personal"), Icon: lucide.Bell(),
		Description: g.text("gallery.template_preview_changes_stay_in_this_session_no_notifications_are_sen"),
		Groups: []components.SettingsGroup{{
			ID: "smart",
			Items: []components.SettingsItem{{
				ID: "smart", Title: g.text("gallery.smart_notifications"), Description: g.text("gallery.pause_mobile_alerts_while_you_re_active_on_desktop"),
				Build: func(c *ui.Context) {
					components.Switch(c, &g.settings.smart, g.text("gallery.toggle_smart_notifications"), components.SwitchOptions{HideLabel: true})
				},
			}},
		}},
	}
	for i, row := range notificationRows {
		if notifications.Groups[len(notifications.Groups)-1].ID != row.group {
			notifications.Groups = append(notifications.Groups, components.SettingsGroup{ID: row.group, Title: g.text(row.group)})
		}
		group := &notifications.Groups[len(notifications.Groups)-1]
		group.Items = append(group.Items, components.SettingsItem{
			ID: row.title, Title: g.text(row.title), Description: g.text(row.description),
			Build: func(c *ui.Context) {
				settingsSelect(c, &g.settings.channels[i], fmt.Sprintf(g.text("gallery.s_delivery"), g.text(row.title)), []string{"Email", "In-app", "Off"}).Shrink(0)
			},
		})
	}
	return []components.SettingsCategory{
		{ID: "general", Title: g.text("gallery.general"), Section: g.text("gallery.application"), Icon: lucide.Settings(),
			Groups: []components.SettingsGroup{{ID: "language", Items: []components.SettingsItem{{
				ID: "language", Title: g.text("gallery.language"), Description: g.text("gallery.applies_immediately_system_menus_keep_their_current_language"), Keywords: []string{"language", "English", "中文", g.text("gallery.follow_system")},
				Build: func(c *ui.Context) {
					choice := g.config.Language
					if settingsSelect(c, &choice, g.text("gallery.language_preference"), []string{"", "en", "zh-Hans"}).Changed() {
						g.chooseLanguage(choice)
						c.Invalidate()
					}
				},
			}}}, g.commandSettingsGroup()}},
		{ID: "appearance", Title: g.text("gallery.appearance"), Section: g.text("gallery.application"), Icon: lucide.Palette(),
			Groups: []components.SettingsGroup{{ID: "theme", Items: []components.SettingsItem{{
				ID: "theme", Title: g.text("gallery.theme"), Description: g.text("gallery.applies_immediately_and_follows_the_toolbar_theme_menu"), Keywords: []string{"System", "Light", "Dark", g.text("gallery.system"), g.text("gallery.light"), g.text("gallery.dark")},
				Build: func(c *ui.Context) {
					choice := string(mygo.Theme.Source())
					if settingsSelect(c, &choice, g.text("gallery.theme_preference"), []string{"system", "light", "dark"}).Changed() {
						g.chooseTheme(mygo.ThemeSource(choice))
					}
				},
			}, {
				ID: "pixel-grid", Title: g.text("gallery.pixel_grid_background"), Description: g.text("gallery.pixel_grid_background_description"),
				Build: func(c *ui.Context) {
					enabled := !g.config.DisablePixelGrid
					components.Switch(c, &enabled, g.text("gallery.toggle_pixel_grid_background"), components.SwitchOptions{HideLabel: true}).Changed()
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
		return tr(c, "gallery.follow_system")
	case "en":
		return "English"
	case "zh-Hans":
		return "简体中文"
	case "system":
		return tr(c, "gallery.system")
	case "light":
		return tr(c, "gallery.light")
	case "dark":
		return tr(c, "gallery.dark")
	default:
		return tr(c, choice)
	}
}

func (g *Gallery) commandSettingsGroup() components.SettingsGroup {
	return components.SettingsGroup{ID: "commands", Title: g.text("gallery.command_palette"), Items: []components.SettingsItem{
		{ID: "command-abc", Title: g.text("gallery.command_auto_abc"), Description: g.text("gallery.command_auto_abc_description"), Keywords: []string{"ABC", "input", "Raycast"}, Build: func(c *ui.Context) {
			enabled := !g.config.DisableCommandABC
			components.Switch(c, &enabled, g.text("gallery.toggle_command_auto_abc"), components.SwitchOptions{HideLabel: true}).Disabled(runtime.GOOS != "darwin").Changed()
			if enabled != !g.config.DisableCommandABC {
				g.config.DisableCommandABC = !enabled
				g.savePreferences()
			}
		}},
		{ID: "command-expanded", Title: g.text("gallery.command_default_expanded"), Description: g.text("gallery.command_default_expanded_description"), Keywords: []string{"search", "results", "Raycast"}, Build: func(c *ui.Context) {
			expanded := !g.config.CommandExpandOnSearch
			components.Switch(c, &expanded, g.text("gallery.toggle_command_default_expanded"), components.SwitchOptions{HideLabel: true}).Changed()
			if expanded != !g.config.CommandExpandOnSearch {
				g.config.CommandExpandOnSearch = !expanded
				g.savePreferences()
			}
		}},
	}}
}
