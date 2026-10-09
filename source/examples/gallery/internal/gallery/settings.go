package gallery

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/shells/sidebar"
)

// Template values are session-only; appearance uses the real gallery preference.
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
	if !g.settings.open {
		return
	}
	g.settings.initialize()
	t := c.Theme()
	w, h := c.Size()
	inset := sidebar.HeaderHeight(c)
	width := min(t.Space(200), w*.88, max(0, w-2*inset))
	components.Dialog(c, &g.settings.open, components.DialogOptions{Label: "Settings", Inset: inset}, func(panel ui.Element) {
		panel.Width(width).Height(min(t.Space(180), h*.85)).Padding(0).Gap(0).Clip()
		components.Settings(c, &g.settings.view, g.settingsCategories(), components.SettingsOptions{
			Narrow: width < t.Space(170),
			HeaderActions: func(c *ui.Context) {
				if sidebar.IconButton(c, lucide.X(), "Close settings").Clicked() {
					g.settings.open = false
				}
			},
		})
	})
}

func (g *Gallery) settingsCategories() []components.SettingsCategory {
	notifications := components.SettingsCategory{
		ID: "notifications", Title: "Notifications", Section: "PERSONAL", Icon: lucide.Bell(),
		Description: "Template preview · Changes stay in this session; no notifications are sent.",
		Groups: []components.SettingsGroup{{
			ID: "smart",
			Items: []components.SettingsItem{{
				ID: "smart", Title: "Smart notifications", Description: "Pause mobile alerts while you’re active on desktop.",
				Build: func(c *ui.Context) { components.Switch(c, &g.settings.smart, "Enabled") },
			}},
		}},
	}
	for i, row := range notificationRows {
		if notifications.Groups[len(notifications.Groups)-1].Title != row.group {
			notifications.Groups = append(notifications.Groups, components.SettingsGroup{ID: row.group, Title: row.group})
		}
		group := &notifications.Groups[len(notifications.Groups)-1]
		group.Items = append(group.Items, components.SettingsItem{
			ID: row.title, Title: row.title, Description: row.description,
			Build: func(c *ui.Context) {
				settingsSelect(c, &g.settings.channels[i], row.title+" delivery", []string{"Email", "In-app", "Off"}).Shrink(0)
			},
		})
	}
	return []components.SettingsCategory{
		notifications,
		{ID: "appearance", Title: "Appearance", Section: "APPLICATION", Icon: lucide.Palette(),
			Groups: []components.SettingsGroup{{ID: "theme", Items: []components.SettingsItem{{
				ID: "theme", Title: "Theme", Description: "Applies immediately and follows the toolbar theme menu.", Keywords: []string{"System", "Light", "Dark"},
				Build: func(c *ui.Context) {
					choice := string(mygo.Theme.Source())
					if settingsSelect(c, &choice, "Theme preference", []string{"system", "light", "dark"}).Changed() {
						g.chooseTheme(mygo.ThemeSource(choice))
					}
				},
			}}}},
		},
	}
}

func settingsSelect(c *ui.Context, value *string, label string, choices []string) ui.Element {
	sel := components.Select(c, value)
	sel.Trigger.Label(label).Width(c.Theme().Space(32)).Children(func() { components.Text(c, *value).SingleLine() })
	sel.Popup(func(panel ui.Element) {
		for _, choice := range choices {
			sel.Item(choice).Label(choice).Children(func() { components.Text(c, choice) })
		}
	})
	return sel.Trigger
}
