package gallery

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// honeycombDemoItem is one cell of the gallery's honeycomb demonstration. The
// label is kept as a locale key and translated per build, so switching language
// updates existing items and their optional duplicate suffixes.
type honeycombDemoItem struct {
	id, labelKey string
	suffix       int
	icon         func() *ui.SVG
	ring         ui.Color
}

type honeycombDemoEntry struct {
	key  string
	icon func() *ui.SVG
	ring ui.Color
}

// honeycombDemoPool is the candidate set the Add button draws from at random.
var honeycombDemoPool = []honeycombDemoEntry{
	{"documents", lucide.FileText, ui.RGB(79, 124, 255)},
	{"images", lucide.Image, ui.RGB(236, 120, 82)},
	{"notes", lucide.StickyNote, ui.RGB(226, 178, 60)},
	{"audio", lucide.Music, ui.RGB(206, 96, 176)},
	{"files", lucide.Folder, ui.RGB(72, 168, 128)},
	{"home", lucide.House, ui.RGB(120, 140, 158)},
	{"profile", lucide.UserRound, ui.RGB(150, 118, 226)},
	{"settings", lucide.Settings, ui.RGB(110, 110, 122)},
	{"downloads", lucide.Download, ui.RGB(56, 148, 214)},
	{"favorites", lucide.Star, ui.RGB(224, 86, 92)},
	{"inbox", lucide.MessageSquare, ui.RGB(66, 176, 176)},
	{"mail", lucide.Mail, ui.RGB(120, 176, 92)},
	{"bookmarks", lucide.Bookmark, ui.RGB(196, 108, 148)},
	{"archive", lucide.Archive, ui.RGB(158, 138, 96)},
	{"cloud", lucide.Cloud, ui.RGB(122, 162, 214)},
	{"lock", lucide.Lock, ui.RGB(132, 112, 176)},
	{"camera", lucide.Camera, ui.RGB(196, 118, 198)},
	{"videos", lucide.Video, ui.RGB(88, 128, 196)},
}

func honeycombDefaults() []honeycombDemoItem {
	items := make([]honeycombDemoItem, len(honeycombDemoPool))
	for i, entry := range honeycombDemoPool {
		items[i] = honeycombDemoItem{id: entry.key, labelKey: entry.key, icon: entry.icon, ring: entry.ring}
	}
	return items
}

func honeycombLabel(c *ui.Context, item honeycombDemoItem) string {
	label := tr(c, item.labelKey)
	if item.suffix > 0 {
		label = fmt.Sprintf("%s %d", label, item.suffix)
	}
	return label
}

// honeycombAdd appends one cell chosen at random from the pool, numbering a
// label whose icon already appears.
func honeycombAdd(s *controlState) {
	entry := honeycombDemoPool[rand.IntN(len(honeycombDemoPool))]
	suffix := 0
	for _, item := range s.honeycombItems {
		if item.labelKey == entry.key {
			suffix++
		}
	}
	if suffix > 0 {
		suffix++
	}
	id := fmt.Sprintf("%s-%d", entry.key, s.honeycombNext)
	s.honeycombNext++
	s.honeycombItems = append(s.honeycombItems, honeycombDemoItem{id: id, labelKey: entry.key, suffix: suffix, icon: entry.icon, ring: entry.ring})
	// Center the new cell on the next build, once it is among the items.
	s.honeycomb = id
}

// honeycombRemove drops one cell at random. The selection follows on the next
// build, when the honeycomb notices its item is gone.
func honeycombRemove(s *controlState) {
	if len(s.honeycombItems) == 0 {
		return
	}
	i := rand.IntN(len(s.honeycombItems))
	s.honeycombItems = slices.Delete(s.honeycombItems, i, i+1)
}

// honeycombToolbar holds the sample's item actions in the page titlebar.
func honeycombToolbar(c *ui.Context, s *controlState) {
	if sidebar.IconButton(c, lucide.Plus(), tr(c, "honeycombAdd")).Clicked() {
		honeycombAdd(s)
	}
	if sidebar.IconButton(c, lucide.Minus(), tr(c, "honeycombRemove")).Disabled(len(s.honeycombItems) == 0).Clicked() {
		honeycombRemove(s)
	}
}

func controlHoneycomb(c *ui.Context, s *controlState, _ string) {
	if s.honeycombItems == nil {
		s.honeycombItems = honeycombDefaults()
		s.honeycomb = s.honeycombItems[0].id
	}
	if !slices.ContainsFunc(s.honeycombItems, func(item honeycombDemoItem) bool { return item.id == s.honeycomb }) {
		if len(s.honeycombItems) > 0 {
			s.honeycomb = s.honeycombItems[0].id
		} else {
			s.honeycomb = ""
		}
	}
	selected := s.honeycomb
	items := make([]components.HoneycombItem, len(s.honeycombItems))
	selectedLabel := ""
	for i, item := range s.honeycombItems {
		label := honeycombLabel(c, item)
		ring, icon := item.ring, item.icon
		items[i] = components.HoneycombItem{
			ID: item.id, Label: label, Ring: ring,
			Content: func(c *ui.Context, size float32) {
				ui.Box(c).Fill().Radius(size/2).Background(ring.Alpha(.16)).Border(1, ring.Alpha(.45)).Center().Children(func() {
					components.Icon(c, icon).Size(size*.42, size*.42).TextColor(ring)
				})
			},
		}
		// Favorites wear a badge on their edge, over the ring.
		if item.labelKey == "favorites" {
			items[i].Overlay = func(c *ui.Context, size float32) {
				d := size * .3
				ui.Box(c).Absolute().Right(0).Bottom(0).Size(d, d).Radius(d / 2).Center().Background(c.Theme().Background).Children(func() {
					components.Icon(c, lucide.Heart).Size(d*.6, d*.6).TextColor(ring)
				})
			}
		}
		if item.id == selected {
			selectedLabel = msg(c).SelectedS(label)
		}
	}
	ui.Column(c).Key("honeycomb-demo").FillWidth().Grow(1).MinWidth(0).MinHeight(0).Children(func() {
		ui.Box(c).FillWidth().Grow(1).Shrink(1).MinWidth(0).MinHeight(0).Children(func() {
			if chosen := components.Honeycomb(c, components.HoneycombOptions{Items: items, Selected: selected,
				SelectedLabel: selectedLabel}); chosen != "" {
				s.honeycomb = chosen
			}
		})
	})
}
