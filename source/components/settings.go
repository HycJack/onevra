package components

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
)

// SettingsState is application-owned presentation state, independent of setting values.
// An empty or unknown Category selects the first category. Searching preserves Category.
type SettingsState struct {
	Category string
	Query    string
}

// SettingsCategory describes a navigation destination. IDs must be nonempty and
// unique among siblings; keep them stable across localization and rebuilds.
type SettingsCategory struct {
	ID, Title, Section, Description string
	Icon                            *ui.SVG
	Groups                          []SettingsGroup
}

// SettingsGroup groups related items. Title may be empty for an unheaded group.
type SettingsGroup struct {
	ID, Title string
	Items     []SettingsItem
}

// SettingsItem describes a searchable row. Build runs synchronously on the UI
// thread only when visible. Applications own values, accessible control labels,
// persistence and effects; Build must not start background work during rendering.
// IDs must be nonempty and unique within the group.
type SettingsItem struct {
	ID, Title, Description string
	Keywords               []string
	Build                  func(c *ui.Context)
}

// SettingsOptions controls presentation without imposing a window or dialog.
// Empty text fields use English defaults. Narrow is supplied by the host from
// its allocated width (not the window width); it replaces the sidebar with a
// category selector and stacks each row's control below its text.
type SettingsOptions struct {
	Narrow                                                                        bool
	Title, SearchLabel, SearchPlaceholder, CategoryLabel, ResultsTitle, EmptyText string
	// HeaderActions builds optional trailing actions, such as a dialog close button.
	HeaderActions func(c *ui.Context)
}

// SettingsResult identifies an item by slice indices in the supplied catalog.
// Results are valid only for that catalog snapshot and preserve declaration order.
type SettingsResult struct{ Category, Group, Item int }

// SearchSettings searches every category using case-insensitive, whitespace-
// separated terms. Every term must occur in the combined category, group, title,
// description or keywords. An empty query returns all items. No builders run.
func SearchSettings(categories []SettingsCategory, query string) []SettingsResult {
	terms := strings.Fields(strings.ToLower(query))
	var results []SettingsResult
	for ci, category := range categories {
		for gi, group := range category.Groups {
			for ii, item := range group.Items {
				haystack := strings.ToLower(strings.Join([]string{category.Title, category.Section, category.Description, group.Title, item.Title, item.Description, strings.Join(item.Keywords, " ")}, " "))
				match := true
				for _, term := range terms {
					if !strings.Contains(haystack, term) {
						match = false
						break
					}
				}
				if match {
					results = append(results, SettingsResult{ci, gi, ii})
				}
			}
		}
	}
	return results
}

// Settings builds a bounded, scrollable settings surface using application-owned
// state and definitions. Place it in a parent with a finite height. It fills its
// parent, and can be embedded in Dialog, Shell content, or a standalone window.
// The host supplies the background or material; navigation adds a subtle tint.
// Use a keyed parent to distinguish multiple instances. state must not be nil.
// Navigation clears Query; clearing Query alone restores the selected category.
func Settings(c *ui.Context, state *SettingsState, categories []SettingsCategory, opts SettingsOptions) ui.Element {
	if state == nil {
		panic("components.Settings: nil state")
	}
	if opts.Title == "" {
		opts.Title = "Settings"
	}
	if opts.SearchLabel == "" {
		opts.SearchLabel = "Search settings"
	}
	if opts.SearchPlaceholder == "" {
		opts.SearchPlaceholder = "Search settings…"
	}
	if opts.CategoryLabel == "" {
		opts.CategoryLabel = "Settings category"
	}
	if opts.ResultsTitle == "" {
		opts.ResultsTitle = "Search results"
	}
	if opts.EmptyText == "" {
		opts.EmptyText = "No matching settings. Try another search."
	}
	selected := -1
	for i, category := range categories {
		if category.ID == state.Category {
			selected = i
			break
		}
	}
	if selected < 0 && len(categories) > 0 {
		selected = 0
		state.Category = categories[0].ID
	}
	t := c.Theme()
	search := func() {
		TextBox(c, &state.Query, TextBoxOptions{LeadingIcon: lucide.Search()}).Label(opts.SearchLabel).Placeholder(opts.SearchPlaceholder).FillWidth().MinWidth(0)
	}
	heading := func(title string, actions bool) {
		ui.Row(c).FillWidth().MinHeight(t.Space(14)).Shrink(0).Padding(t.Space(3), t.Space(4)).Gap(t.Space(2)).Children(func() {
			Text(c, title).Bold()
			ui.Spacer(c)
			if actions && opts.HeaderActions != nil {
				opts.HeaderActions(c)
			}
		})
	}
	return ui.Row(c).Fill().MinWidth(0).Gap(0).Children(func() {
		if !opts.Narrow {
			ui.Column(c).Width(t.Space(54)).FillHeight().Shrink(0).Background(t.Text.Alpha(.035)).Children(func() {
				heading(opts.Title, false)
				ui.Column(c).FillWidth().Padding(0, t.Space(3), t.Space(4), t.Space(3)).Children(search)
				ScrollView(c, ScrollViewOptions{}).Fill().Children(func() {
					Sidebar(c).Children(func() {
						section := ""
						for i, category := range categories {
							if category.Section != "" && category.Section != section {
								SidebarSection(c, category.Section, nil)
							}
							section = category.Section
							if SidebarItem(c, category.Title, category.Icon, selected == i).Key(category.ID).Clicked() {
								state.Category = category.ID
								selected = i
								state.Query = ""
							}
						}
					})
				})
			})
			ui.Box(c).Width(1).FillHeight().Shrink(0).Background(t.Border)
		}
		ui.Column(c).Fill().MinWidth(0).Children(func() {
			title := opts.Title
			if !opts.Narrow && selected >= 0 {
				title = categories[selected].Title
			}
			if strings.TrimSpace(state.Query) != "" {
				title = opts.ResultsTitle
			}
			heading(title, true)
			if opts.Narrow {
				ui.Column(c).FillWidth().Padding(0, t.Space(3), t.Space(3), t.Space(3)).Gap(t.Space(2)).Children(func() {
					if selected >= 0 {
						sel := Select(c, &state.Category)
						sel.Trigger.Label(opts.CategoryLabel).FillWidth().Children(func() { Text(c, categories[selected].Title).SingleLine() })
						sel.Popup(func(panel ui.Element) {
							for _, category := range categories {
								sel.Item(category.ID).Label(category.Title).Children(func() { Text(c, category.Title) })
							}
						})
						if sel.Trigger.Changed() {
							state.Query = ""
							for i := range categories {
								if categories[i].ID == state.Category {
									selected = i
								}
							}
						}
					}
					search()
				})
			}
			searching := strings.TrimSpace(state.Query) != ""
			results := SearchSettings(categories, state.Query)
			ScrollView(c, ScrollViewOptions{}).Fill().Children(func() {
				ui.Column(c).Key(state.Category).FillWidth().MinWidth(0).Padding(t.Space(5)).Gap(t.Space(5)).Children(func() {
					count, lastCategory, lastGroup := 0, -1, -1
					for _, result := range results {
						if !searching && result.Category != selected {
							continue
						}
						category := categories[result.Category]
						group := category.Groups[result.Group]
						item := group.Items[result.Item]
						if lastCategory != result.Category {
							if searching {
								Text(c, category.Title).Bold()
							}
							if category.Description != "" {
								Text(c, category.Description).TextColor(t.TextMuted).Wrap()
							}
						}
						if (lastCategory != result.Category || lastGroup != result.Group) && group.Title != "" {
							ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() { Text(c, group.Title).Bold(); ui.Box(c).FillWidth().Height(1).Background(t.Border) })
						}
						lastCategory, lastGroup = result.Category, result.Group
						count++
						var row ui.Element
						if opts.Narrow {
							row = ui.Column(c)
						} else {
							row = ui.Row(c)
						}
						row.Key(struct{ Category, Group, Item string }{category.ID, group.ID, item.ID}).FillWidth().MinWidth(0).Gap(t.Space(3)).Children(func() {
							ui.Column(c).FillWidth().Shrink(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
								Text(c, item.Title).Bold().Wrap()
								if item.Description != "" {
									Text(c, item.Description).TextColor(t.TextMuted).FontSize(t.Rem(.875)).Wrap()
								}
							})
							if item.Build != nil {
								item.Build(c)
							}
						})
					}
					if count == 0 {
						Text(c, opts.EmptyText).TextColor(t.TextMuted).Wrap()
					}
				})
			})
		})
	})
}
