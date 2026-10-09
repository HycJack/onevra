package gallery

import (
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/shells/sidebar"
)

// Template values are session-only; appearance uses the real gallery preference.
type settingsState struct {
	open, initialized bool
	page, query       string
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
	s.initialized, s.smart, s.page = true, true, "Notifications"
	for _, row := range notificationRows {
		s.channels = append(s.channels, row.channel)
	}
}

func settingsMatch(query string, values ...string) bool {
	return strings.Contains(strings.ToLower(strings.Join(values, " ")), strings.ToLower(strings.TrimSpace(query)))
}

// Settings builds a modal template over the current gallery page, never a Shell replacement page.
func (g *Gallery) Settings(c *ui.Context) {
	if !g.settings.open {
		return
	}
	g.settings.initialize()
	s, t := &g.settings, c.Theme()
	w, h := c.Size()
	inset := sidebar.HeaderHeight(c)
	width := min(t.Space(200), w*.88, max(0, w-2*inset))
	narrow := width < t.Space(170)
	components.Dialog(c, &s.open, components.DialogOptions{Label: "Settings", Inset: inset}, func(panel ui.Element) {
		panel.Width(width).Height(min(t.Space(180), h*.85)).Padding(0).Gap(0).Clip()
		ui.Row(c).Fill().Gap(0).Children(func() {
			if !narrow {
				ui.Column(c).Width(t.Space(54)).FillHeight().Shrink(0).Background(t.Background.Mix(t.Text, .035)).Children(func() {
					g.settingsHeading(c, "Settings", false)
					ui.Column(c).FillWidth().Padding(0, t.Space(3), t.Space(4), t.Space(3)).Children(func() {
						components.TextBox(c, &s.query, components.TextBoxOptions{LeadingIcon: lucide.Search()}).
							Label("Search settings").Placeholder("Search settings…").FillWidth().MinWidth(0)
					})
					components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
						components.Sidebar(c).Children(func() {
							components.SidebarSection(c, "PERSONAL", func() {
								if components.SidebarItem(c, "Notifications", lucide.Bell(), s.page == "Notifications").Clicked() {
									s.page = "Notifications"
								}
							})
							components.SidebarSection(c, "APPLICATION", func() {
								if components.SidebarItem(c, "Appearance", lucide.Palette(), s.page == "Appearance").Clicked() {
									s.page = "Appearance"
								}
							})
						})
					})
				})
				ui.Box(c).Width(1).FillHeight().Shrink(0).Background(t.Border)
			}
			ui.Column(c).Fill().MinWidth(0).Background(t.Background).Children(func() {
				if narrow {
					g.settingsHeading(c, "Settings", true)
					ui.Column(c).FillWidth().Padding(0, t.Space(3), t.Space(3), t.Space(3)).Gap(t.Space(2)).Children(func() {
						settingsSelect(c, &s.page, "Settings category", []string{"Notifications", "Appearance"}).FillWidth()
						components.TextBox(c, &s.query, components.TextBoxOptions{LeadingIcon: lucide.Search()}).Label("Search settings").Placeholder("Search settings…").FillWidth()
					})
				} else {
					g.settingsHeading(c, s.page, true)
				}
				components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
					ui.Column(c).Key(s.page).FillWidth().MinWidth(0).Padding(t.Space(5)).Gap(t.Space(5)).Children(func() {
						if s.page == "Appearance" {
							g.appearanceSettings(c, narrow)
						} else {
							g.notificationSettings(c, narrow)
						}
					})
				})
			})
		})
	})
}

func (g *Gallery) settingsHeading(c *ui.Context, title string, closeButton bool) {
	t := c.Theme()
	ui.Row(c).FillWidth().MinHeight(t.Space(14)).Shrink(0).Padding(t.Space(3), t.Space(4)).Gap(t.Space(2)).Children(func() {
		components.Text(c, title).Bold()
		ui.Spacer(c)
		if closeButton && sidebar.IconButton(c, lucide.X(), "Close settings").Clicked() {
			g.settings.open = false
		}
	})
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

func settingsRow(c *ui.Context, title, description string, narrow bool, control func()) {
	var container ui.Element
	if narrow {
		container = ui.Column(c)
	} else {
		container = ui.Row(c)
	}
	container.Key(title).FillWidth().MinWidth(0).Gap(c.Theme().Space(3)).Children(func() {
		text := ui.Column(c).MinWidth(0).Gap(c.Theme().Space(1))
		text.FillWidth().Shrink(1)
		text.Children(func() {
			components.Text(c, title).Bold().Wrap()
			components.Text(c, description).TextColor(c.Theme().TextMuted).FontSize(c.Theme().Rem(.875)).Wrap()
		})
		control()
	})
}

func (g *Gallery) notificationSettings(c *ui.Context, narrow bool) {
	s, t := &g.settings, c.Theme()
	components.Text(c, "Choose which updates you want to receive and where they should appear.").TextColor(t.TextMuted).Wrap()
	components.Text(c, "Template preview · Changes stay in this session; no notifications are sent.").FontSize(t.Rem(.8)).TextColor(t.TextMuted).Wrap()
	matches := 0
	if settingsMatch(s.query, "Notifications", "Smart notifications", "Pause mobile alerts while you’re active on desktop") {
		matches++
		settingsRow(c, "Smart notifications", "Pause mobile alerts while you’re active on desktop.", narrow, func() {
			components.Switch(c, &s.smart, "Enabled")
		})
	}
	group := ""
	for i, row := range notificationRows {
		if !settingsMatch(s.query, "Notifications", row.group, row.title, row.description) {
			continue
		}
		matches++
		if group != row.group {
			group = row.group
			ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
				components.Text(c, group).Bold()
				ui.Box(c).FillWidth().Height(1).Background(t.Border)
			})
		}
		settingsRow(c, row.title, row.description, narrow, func() {
			settingsSelect(c, &s.channels[i], row.title+" delivery", []string{"Email", "In-app", "Off"}).Shrink(0)
		})
	}
	if matches == 0 {
		components.Text(c, "No matching settings. Try another search.").TextColor(t.TextMuted).Wrap()
	}
}

func (g *Gallery) appearanceSettings(c *ui.Context, narrow bool) {
	if !settingsMatch(g.settings.query, "Appearance", "Theme", "System Light Dark") {
		components.Text(c, "No matching settings. Try another search.").Wrap()
		return
	}
	components.Text(c, "Make the gallery feel at home.").TextColor(c.Theme().TextMuted).Wrap()
	settingsRow(c, "Theme", "Applies immediately and follows the toolbar theme menu.", narrow, func() {
		choice := string(mygo.Theme.Source())
		if settingsSelect(c, &choice, "Theme preference", []string{"system", "light", "dark"}).Changed() {
			g.chooseTheme(mygo.ThemeSource(choice))
		}
	})
}
