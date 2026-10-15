package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// These categories demonstrate row kinds beyond a trailing select or switch:
// leading content, stacked controls and rows built entirely by the application.

// settingsGlyph is a row's leading icon on a tinted tile.
func settingsGlyph(c *ui.Context, icon *ui.SVG) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(8), t.Space(8)).Radius(t.Radius).Background(t.Text.Alpha(.06)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		components.Icon(c, icon).Size(t.Space(4), t.Space(4)).TextColor(t.TextMuted)
	})
}

// layoutSettingsGroup toggles the surface's own search field and sidebar.
func (g *Gallery) layoutSettingsGroup() components.SettingsGroup {
	toggle := func(id, title, description, label string, hidden *bool) components.SettingsItem {
		return components.SettingsItem{ID: id, Title: g.text(title), Description: g.text(description), Build: func(c *ui.Context) {
			shown := !*hidden
			if components.Switch(c, &shown, g.text(label), components.SwitchOptions{HideLabel: true}).Changed() {
				*hidden = !shown
			}
		}}
	}
	return components.SettingsGroup{ID: "layout", Title: g.text("settingsLayout"), Items: []components.SettingsItem{
		toggle("search", "showSearch", "showSearchDescription", "toggleSettingsSearch", &g.settings.hideSearch),
		toggle("sidebar", "showSidebar", "showSidebarDescription", "toggleSettingsSidebar", &g.settings.hideSidebar),
	}}
}

func (g *Gallery) layoutSettings() components.SettingsCategory {
	return components.SettingsCategory{ID: "layout", Title: g.text("layout"), Section: g.text("personal"), Icon: lucide.LayoutGrid(),
		Groups: []components.SettingsGroup{{ID: "layout", Items: []components.SettingsItem{{
			ID: "scale", Title: g.text("interfaceScale"), Description: g.text("interfaceScaleDescription"), Keywords: []string{"zoom", "缩放"}, Stacked: true,
			Build: func(c *ui.Context) {
				ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(c.Theme().Space(3)).Children(func() {
					components.Slider(c, &g.settings.scale, 75, 150, components.SliderOptions{Step: 25}).Label(g.text("interfaceScale")).Grow(1).MinWidth(0)
					components.Text(c, fmt.Sprintf("%.0f%%", g.settings.scale)).SingleLine().Width(c.Theme().Space(12)).Shrink(0)
				})
			},
		}, {
			ID: "density", Title: g.text("density"), Description: g.text("densityDescription"), Stacked: true,
			Build: func(c *ui.Context) {
				ui.RadioGroup(c, func() {
					for _, choice := range []string{"compact", "comfortable", "spacious"} {
						components.Radio(c, &g.settings.density, choice, g.text(choice))
					}
				}).Label(g.text("density")).Row().FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4))
			},
		}}}},
	}
}

func (g *Gallery) editorSettings() components.SettingsCategory {
	return components.SettingsCategory{ID: "editor", Title: g.text("editor"), Section: g.text("personal"), Icon: lucide.Code(),
		Groups: []components.SettingsGroup{{ID: "text", Items: []components.SettingsItem{{
			ID: "font-size", Title: g.text("fontSize"), Description: g.text("fontSizeDescription"), Keywords: []string{"font", "字体"},
			Leading: func(c *ui.Context) { settingsGlyph(c, lucide.Type()) },
			Build: func(c *ui.Context) {
				components.NumberInput(c, &g.settings.fontSize, 8, 32, 1).Label(g.text("fontSize")).Width(c.Theme().Space(28)).Shrink(0)
			},
		}, {
			ID: "tab-width", Title: g.text("tabWidth"), Description: g.text("tabWidthDescription"), Keywords: []string{"indent", "缩进"},
			Leading: func(c *ui.Context) { settingsGlyph(c, lucide.Keyboard()) },
			Build: func(c *ui.Context) {
				ui.RadioGroup(c, func() {
					for _, width := range []int{2, 4, 8} {
						components.Radio(c, &g.settings.tabWidth, width, fmt.Sprint(width))
					}
				}).Label(g.text("tabWidth")).Row().Shrink(0).Gap(c.Theme().Space(4))
			},
		}}}, {ID: "display", Items: []components.SettingsItem{{
			ID: "gutter", Title: g.text("editorDisplay"), Description: g.text("editorDisplayDescription"), Stacked: true,
			Keywords: []string{g.text("lineNumbers"), g.text("wordWrap"), g.text("minimap")},
			Build: func(c *ui.Context) {
				ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4)).Children(func() {
					components.Checkbox(c, &g.settings.lineNumbers, g.text("lineNumbers"))
					components.Checkbox(c, &g.settings.wordWrap, g.text("wordWrap"))
					components.Checkbox(c, &g.settings.minimap, g.text("minimap"))
				})
			},
		}}}},
	}
}

func (g *Gallery) storageSettings() components.SettingsCategory {
	const capacity = 10.0
	return components.SettingsCategory{ID: "storage", Title: g.text("storage"), Section: g.text("personal"), Icon: lucide.HardDrive(),
		Groups: []components.SettingsGroup{{ID: "usage", Items: []components.SettingsItem{{
			// The whole row is application content; its title remains searchable.
			ID: "usage", Title: g.text("storageUsage"), Keywords: []string{g.text("storageCache"), g.text("documents")},
			Content: func(c *ui.Context) {
				t := c.Theme()
				documents := 3.1
				used := documents + g.settings.cacheGB
				parts := []struct {
					key   string
					color ui.Color
					size  float64
				}{{"documents", t.Accent, documents}, {"storageCache", t.Text.Alpha(.4), g.settings.cacheGB}, {"storageOther", t.Border, capacity - used}}
				ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(3)).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border).Children(func() {
					ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
						components.Text(c, g.text("storageUsage")).Bold().Grow(1).MinWidth(0)
						components.Text(c, fmt.Sprintf("%.1f / %.0f GB", used, capacity)).TextColor(t.TextMuted).SingleLine()
					})
					// A segmented bar shows what each kind of data occupies.
					ui.Row(c).FillWidth().Height(t.Space(2)).Radius(t.Space(1)).Clip().Gap(0).Role(ui.RoleProgress).Range(0, 1, used/capacity).Label(g.text("storageUsage")).Children(func() {
						for _, part := range parts {
							ui.Box(c).WidthPercent(float32(part.size / capacity * 100)).FillHeight().Background(part.color)
						}
					})
					ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(4)).Children(func() {
						for _, part := range parts {
							ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
								ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(t.Space(1)).Background(part.color)
								components.Text(c, fmt.Sprintf("%s %.1f GB", g.text(part.key), part.size)).FontSize(t.Rem(.875)).TextColor(t.TextMuted)
							})
						}
					})
				})
			},
		}, {
			ID: "clear-cache", Title: g.text("clearCache"), Description: g.text("clearCacheDescription"),
			Leading: func(c *ui.Context) { settingsGlyph(c, lucide.Trash()) },
			Build: func(c *ui.Context) {
				ui.Row(c).Shrink(0).Children(func() {
					if components.Button(c, g.text("clear"), components.ButtonSecondary).Label(g.text("clearCache")).Clicked() {
						g.settings.cacheGB = 0
					}
				}).Disabled(g.settings.cacheGB == 0)
			},
		}, {
			ID: "downloads", Title: g.text("downloadLocation"), Description: g.text("downloadLocationDescription"), Stacked: true,
			Build: func(c *ui.Context) {
				components.TextBox(c, &g.settings.downloads, components.TextBoxOptions{LeadingIcon: lucide.Folder()}).Label(g.text("downloadLocation")).FillWidth().MinWidth(0)
			},
		}}}},
	}
}

func (g *Gallery) accountSettings() components.SettingsCategory {
	return components.SettingsCategory{ID: "account", Title: g.text("account"), Section: g.text("personal"), Icon: lucide.User(),
		Groups: []components.SettingsGroup{{ID: "profile", Items: []components.SettingsItem{{
			// A profile card that follows the fields below as they are edited.
			ID: "profile", Title: g.text("profile"), Keywords: []string{"avatar", "头像"},
			Content: func(c *ui.Context) {
				t := c.Theme()
				ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(4)).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border).Children(func() {
					initials := ""
					for _, word := range strings.Fields(g.settings.displayName) {
						initials += strings.ToUpper(string([]rune(word)[:1]))
					}
					ui.Box(c).Size(t.Space(12), t.Space(12)).Shrink(0).Radius(t.Space(6)).Background(g.settings.profileColor).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
						components.Text(c, string([]rune(initials)[:min(2, len([]rune(initials)))])).Bold().TextColor(ui.RGB(255, 255, 255)).SingleLine()
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						components.Text(c, g.settings.displayName).Bold().SingleLine()
						components.Text(c, "alex@example.com").TextColor(t.TextMuted).SingleLine()
					})
				})
			},
		}}}, {ID: "identity", Title: g.text("identity"), Items: []components.SettingsItem{{
			ID: "name", Title: g.text("displayName"), Description: g.text("displayNameDescription"), Stacked: true,
			Build: func(c *ui.Context) {
				components.TextBox(c, &g.settings.displayName, components.TextBoxOptions{}).Label(g.text("displayName")).FillWidth().MinWidth(0)
			},
		}, {
			ID: "bio", Title: g.text("bio"), Description: g.text("bioDescription"), Stacked: true,
			Build: func(c *ui.Context) {
				components.TextBox(c, &g.settings.bio, components.TextBoxOptions{Multiline: true}).Label(g.text("bio")).Lines(3, 6).FillWidth().MinWidth(0)
			},
		}, {
			ID: "color", Title: g.text("profileColor"), Description: g.text("profileColorDescription"),
			Build: func(c *ui.Context) {
				components.ColorPicker(c, &g.settings.profileColor).Label(g.text("profileColor")).Shrink(0)
			},
		}}}},
	}
}

func (g *Gallery) aboutSettings() components.SettingsCategory {
	return components.SettingsCategory{ID: "about", Title: g.text("about"), Section: g.text("personal"), Icon: lucide.Info(),
		Groups: []components.SettingsGroup{{ID: "about", Items: []components.SettingsItem{{
			ID: "about", Title: g.text("onevraGallery"), Description: g.text("aboutGallery"), Keywords: []string{"version", "update", "版本", "更新"},
			Content: func(c *ui.Context) {
				t := c.Theme()
				ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(3)).Padding(t.Space(6), t.Space(4)).Children(func() {
					ui.Box(c).Size(t.Space(16), t.Space(16)).Radius(t.Space(4)).Background(t.Accent).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
						components.Icon(c, lucide.Hexagon()).Size(t.Space(8), t.Space(8)).TextColor(t.AccentText)
					})
					components.Text(c, g.text("onevraGallery")).Bold().FontSize(t.Rem(1.25))
					components.Badge(c, "MyGo").AlignSelf(ui.Center)
					components.Text(c, g.text("aboutGallery")).TextColor(t.TextMuted).Wrap()
					if g.settings.upToDate {
						components.Text(c, g.text("upToDate")).TextColor(t.Success)
					} else if components.Button(c, g.text("checkForUpdates"), components.ButtonOutline).Clicked() {
						g.settings.upToDate = true
					}
				})
			},
		}}}},
	}
}
