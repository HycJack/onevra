package gallery

import "strings"

// Component sections describe purpose, never implementation ownership.
var navigationGroups = []string{
	"Plugins", "Actions", "Selection & values", "Text input",
	"Collections", "Layout", "Containers", "Navigation", "Overlays",
	"Status", "Media & text",
}

type navigationEntry struct {
	label string
	page  galleryPage
}

type navigationGroup struct {
	label   string
	entries []navigationEntry
}

func filteredNavigation(query string) []navigationGroup {
	query = strings.ToLower(strings.TrimSpace(query))
	var groups []navigationGroup
	add := func(label string, entries []navigationEntry) {
		var matches []navigationEntry
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.label), query) {
				matches = append(matches, entry)
			}
		}
		if len(matches) > 0 {
			groups = append(groups, navigationGroup{label, matches})
		}
	}
	add("Overview", []navigationEntry{{"Welcome", welcomePage}, {"Component status", componentStatusPage}})
	add("Icon library", []navigationEntry{{"Lucide", lucidePage}})
	add("Desktop", []navigationEntry{{"Tray", trayPage}})
	for _, group := range navigationGroups {
		var entries []navigationEntry
		for i, example := range controlCatalog {
			if example.group == group {
				entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i)})
			}
		}
		for i, example := range mygoCatalog {
			if example.group == group {
				entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i)})
			}
		}
		add(group, entries)
	}
	return groups
}
