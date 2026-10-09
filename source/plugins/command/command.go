// Package command provides application-owned command palettes and dialogs.
// Hosts register opening shortcuts, retain State, and execute activated IDs.
package command

import (
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// Item describes a command. IDs must be nonempty and unique, independent of
// localized titles. Shortcut is display text; the host registers actual keys.
// Items in the same Group are displayed together in first-appearance order.
type Item struct {
	ID, Title, Description, Group, Shortcut string
	Keywords                                []string
	Icon                                    *ui.SVG
	Disabled                                bool
}

// State is retained by the host for one palette; do not copy it after use.
// Query and ActiveID may be updated by the host. The zero value is ready to use.
// Closing a dialog preserves the query; hosts choose when to reset it.
type State struct {
	// InputSourceError records the last selection or restoration failure.
	// A failure never prevents searching. It is cleared on the next attempt.
	InputSourceError      error
	inputSource           inputSourceSession
	Query                 string
	ActiveID              string
	scroll                components.ListState
	lastQuery, lastActive string
	lastIndex             int
	composing             bool
	refocus               bool
}

// Options supplies application text, results and activation behavior.
type Options struct {
	// InputSource optionally selects an enabled macOS input source once when
	// search receives focus, e.g. InputSourceABC. Empty keeps the current source.
	// Selection waits for any marked-text composition to finish. Unsupported
	// platforms and headless views report State.InputSourceError without switching.
	InputSource string
	// KeepInputSource skips restoration on close. By default restoration only
	// occurs if the user has not manually selected a different source.
	KeepInputSource bool
	// ExpandOnSearch hides the results area until Query contains non-whitespace
	// committed text, and collapses it again when cleared. Hidden results cannot
	// be activated; loading, error and empty states are also hidden until search.
	// The zero value shows the results area immediately.
	ExpandOnSearch bool
	// PreFiltered uses host-supplied matches instead of applying Filter.
	// For asynchronous search, the host owns cancellation and stale-result checks.
	PreFiltered bool
	// Loading and Error replace results and suppress activation. Error wins.
	Loading bool
	Error   string
	// Empty labels use English defaults.
	SearchLabel, Placeholder, ResultsLabel, EmptyText, LoadingText string
	// OnActivate runs on the UI thread for Enter or a single click. The host
	// owns execution, background work and whether to close the dialog.
	OnActivate func(id string)
}

// Filter returns matching items without modifying the input. All whitespace-
// separated query terms must match the title or keywords, case-insensitively.
// Exact titles rank before title prefixes, title substrings, then keywords;
// ties retain input order. An empty query returns a copy in original order.
func Filter(items []Item, query string) []Item {
	terms := strings.Fields(strings.ToLower(query))
	type match struct {
		item  Item
		score int
	}
	matches := make([]match, 0, len(items))
	for _, item := range items {
		title := strings.ToLower(item.Title)
		keywords := strings.ToLower(strings.Join(item.Keywords, " "))
		score, ok := 0, true
		for _, term := range terms {
			switch {
			case title == term:
			case strings.HasPrefix(title, term):
				score += 1
			case strings.Contains(title, term):
				score += 2
			case strings.Contains(keywords, term):
				score += 3
			default:
				ok = false
			}
		}
		if ok {
			matches = append(matches, match{item, score})
		}
	}
	slices.SortStableFunc(matches, func(a, b match) int { return a.score - b.score })
	out := make([]Item, len(matches))
	for i, m := range matches {
		out[i] = m.item
	}
	return out
}

func grouped(items []Item) []Item {
	var groups []string
	buckets := make(map[string][]Item)
	for _, item := range items {
		if _, ok := buckets[item.Group]; !ok {
			groups = append(groups, item.Group)
		}
		buckets[item.Group] = append(buckets[item.Group], item)
	}
	out := make([]Item, 0, len(items))
	for _, group := range groups {
		out = append(out, buckets[group]...)
	}
	return out
}

func (s *State) sync(items []Item) int {
	index := -1
	for i, item := range items {
		if !item.Disabled && item.ID != "" && (index < 0 || item.ID == s.ActiveID) {
			index = i
			if item.ID == s.ActiveID {
				break
			}
		}
	}
	s.ActiveID = ""
	if index >= 0 {
		s.ActiveID = items[index].ID
	}
	if s.Query != s.lastQuery || s.ActiveID != s.lastActive || index != s.lastIndex {
		if index >= 0 {
			s.scroll.ScrollIntoView(index)
		}
	}
	s.lastQuery, s.lastActive = s.Query, s.ActiveID
	s.lastIndex = index
	return index
}

func (s *State) move(items []Item, delta int) {
	i := s.sync(items)
	for next := i + delta; next >= 0 && next < len(items); next += delta {
		if !items[next].Disabled && items[next].ID != "" {
			s.ActiveID = items[next].ID
			s.scroll.ScrollIntoView(next)
			return
		}
	}
}

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// CommandPalette builds search and grouped virtual results in a
// bounded parent. Expanded results fill that parent; collapsed search uses only
// the input's height. The host supplies the available size. Up/Down keep
// focus in the search field and skip disabled commands, stopping at either end.
// MyGo owns text editing and IME handling. Tab follows normal focus navigation.
// Stable keyed parents distinguish independently built palettes.
func CommandPalette(c *ui.Context, s *State, items []Item, options Options) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().MinWidth(0).MinHeight(0)
	root.Children(func() {
		input := components.TextBox(c, &s.Query, components.TextBoxOptions{}).
			Label(fallback(options.SearchLabel, "Search commands")).Placeholder(fallback(options.Placeholder, "Search commands…")).
			FillWidth().Shrink(0).FontSize(t.Rem(1.15)).
			Padding(t.Space(4.5), t.Space(4.5), t.Space(4.5), t.Space(12)).
			Background(ui.Transparent).Border(0, ui.Transparent).Radius(0).FocusRing(false)
		input.DrawOver(func(p *ui.Painter, r ui.Rect) {
			size := t.Space(5)
			p.Icon(lucide.Search(), ui.Rect{X: r.X + t.Space(4.5), Y: r.Y + (r.H-size)/2, W: size, H: size}, t.TextMuted)
		})
		// Resolve bound input before deriving results for this frame.
		input.Changed()
		s.composing = input.Composing()
		if s.refocus {
			input.Focus()
			s.refocus = false
		}
		if !s.inputSource.attempted && options.InputSource != "" {
			if !s.inputSource.waited {
				// Modal focus and the native input context are committed after
				// this build. Select on a subsequent frame, not during creation.
				s.inputSource.waited = true
				c.After(time.Millisecond)
			} else if input.Focused() && !s.composing {
				s.InputSourceError = s.inputSource.begin(options.InputSource, options.KeepInputSource, selectInputSource)
			}
		}
		if !resultsExpanded(s, options) {
			s.sync(nil)
			return
		}
		root.FillHeight()
		results := items
		if !options.PreFiltered {
			results = Filter(items, s.Query)
		}
		results = grouped(results)
		if options.Loading || options.Error != "" {
			results = nil
		}
		s.sync(results)
		activate := func(id string) {
			for _, item := range results {
				if item.ID == id && id != "" && !item.Disabled {
					s.ActiveID = id
					if options.OnActivate != nil {
						options.OnActivate(id)
					}
					return
				}
			}
		}
		input.OnShortcut(0, ui.KeyDown, func() { s.move(results, 1) })
		input.OnShortcut(0, ui.KeyUp, func() { s.move(results, -1) })
		input.OnSubmit(func() { activate(s.ActiveID) })
		if len(results) == 0 {
			message := fallback(options.EmptyText, "No matching commands")
			if options.Loading {
				message = fallback(options.LoadingText, "Loading…")
			}
			if options.Error != "" {
				message = options.Error
			}
			ui.Box(c).Fill().MinHeight(0).Padding(t.Space(4)).Center().Children(func() {
				components.Text(c, message).Role(ui.RoleStatus).TextColor(t.TextMuted).Wrap()
			})
		} else {
			s.scroll.Key = func(i int) any { return results[i].ID }
			components.List(c, &s.scroll, len(results), func(i int) {
				item := results[i]
				if item.Group != "" && (i == 0 || results[i-1].Group != item.Group) {
					components.Text(c, item.Group).TextColor(t.TextMuted).FontSize(t.Rem(.8)).Padding(t.Space(2), t.Space(2.5)).SingleLine()
				}
				row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(2.5)).
					Padding(t.Space(2.5), t.Space(2.5)).Radius(t.Radius * 1.5).Role(ui.RoleButton).
					Label(item.Title).Description(item.Description).Disabled(item.Disabled || item.ID == "").Selected(item.ID == s.ActiveID).TextColor(t.Text)
				if item.ID == s.ActiveID {
					row.Background(t.Text.Alpha(.1))
					input.ActiveDescendant(row)
				} else if row.Hovered() && !item.Disabled {
					row.Background(t.Text.Alpha(.04))
				}
				row.OnClick(func() { s.refocus = true; activate(item.ID) })
				row.Children(func() {
					if item.Icon != nil {
						components.Icon(c, item.Icon).Size(t.Space(5), t.Space(5)).Shrink(0)
					}
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						color := t.Text
						if item.Disabled {
							color = t.TextMuted
						}
						components.Text(c, item.Title).TextColor(color).SingleLine()
						if item.Description != "" {
							components.Text(c, item.Description).TextColor(t.TextMuted).FontSize(t.Rem(.85)).SingleLine()
						}
					})
					if item.Shortcut != "" {
						components.Text(c, item.Shortcut).TextColor(t.TextMuted).FontSize(t.Rem(.8)).Shrink(0).SingleLine()
					}
				})
			}).Label(fallback(options.ResultsLabel, "Commands")).Fill().MinHeight(0).Padding(0, t.Space(2), t.Space(3), t.Space(2))
		}
	})
	return root
}

// CommandDialog hosts CommandPalette in Onevra's modal surface. Call every
// frame with retained open and state. The host owns opening shortcuts, query
// resets and closing on activation. Dialog dismissal and focus restoration
// follow components.Dialog. Its returned panel may be resized in this frame.
// With ExpandOnSearch, its height fits the input until a query is committed.
func CommandDialog(c *ui.Context, open *bool, s *State, items []Item, options Options, dialog components.DialogOptions) ui.Element {
	if !*open {
		s.ReleaseInputSource()
		s.composing = false
		s.refocus = false
	}
	// Dismissal is evaluated before the input is built. Preserve the previous
	// frame's marked-text state so Escape ending composition cannot also close
	// the surrounding dialog. Text composition remains owned by MyGo.
	dialog.DisableEscapeClose = dialog.DisableEscapeClose || s.composing
	if dialog.Label == "" {
		dialog.Label = fallback(options.SearchLabel, "Commands")
	}
	panel := components.Dialog(c, open, dialog, func(panel ui.Element) {
		panel.Width(640).Padding(0).Gap(0).Radius(c.Theme().Radius * 3).Clip()
		CommandPalette(c, s, items, options)
		// Read after the palette resolves bound text input so expansion and
		// collapse use this frame's query, without rebuilding the search field.
		if resultsExpanded(s, options) {
			panel.Height(440)
		}
	})
	if !*open {
		s.ReleaseInputSource()
	}
	return panel
}

func resultsExpanded(s *State, options Options) bool {
	return !options.ExpandOnSearch || strings.TrimSpace(s.Query) != ""
}
