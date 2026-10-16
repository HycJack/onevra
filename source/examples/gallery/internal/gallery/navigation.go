package gallery

// Component sections describe purpose, never implementation ownership.
var navigationGroups = []string{
	"Plugins", "Actions & navigation", "Text input", "Selection & values",
	"Date & time", "Color", "Content display", "Code & diffs", "Charts", "Layout & containers",
	"Feedback & overlays", "Games",
}

// newPages are the routes added in the latest published version or since; the
// top New components group indexes them without removing them from their
// purpose groups. Add new pages alongside the existing entries. When
// publishing a version, keep the pages it ships and remove those from earlier
// versions.
var newPages = map[string]bool{
	"controls/Honeycomb":      true,
	"mygo/Rich text editor":   true,
	"mygo/Markdown":           true,
	"mygo/Markdown stream":    true,
	"mygo/System information": true,
	"mygo/Recent documents":   true,
	"mygo/File watching":      true,
	"mygo/Global shortcuts":   true,

	"controls/SplitPercent": true,
	"mygo/Gradients & blur": true,
	"mygo/Shaders":          true,
	"mygo/Live chart":       true,
}

type navigationEntry struct {
	label string
	page  galleryPage
}

// newGroupLabel names the index group of the pages in newPages.
const newGroupLabel = "New components"

type navigationGroup struct {
	label   string
	entries []navigationEntry
	// index marks a shortcut group that repeats pages already listed in
	// their purpose groups. Command search skips it to keep results unique.
	index bool
}

func navigationCatalog() []navigationGroup {
	var groups []navigationGroup
	add := func(label string, entries []navigationEntry) {
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label: label, entries: entries})
		}
	}
	add("Overview", []navigationEntry{{"Welcome", welcomePage}})
	if entries := newNavigationEntries(); len(entries) > 0 {
		groups = append(groups, navigationGroup{label: newGroupLabel, entries: entries, index: true})
	}
	add("Icon library", []navigationEntry{{"Lucide", lucidePage}})
	desktop := []navigationEntry{{"Tray", trayPage}, {"Dock menu", dockPage}, {"System notifications", notificationsPage}}
	for _, group := range append([]string{"Desktop"}, navigationGroups...) {
		var entries []navigationEntry
		if group == "Desktop" {
			entries = desktop
		}
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

// newNavigationEntries collects the pages in newPages in catalog order.
func newNavigationEntries() []navigationEntry {
	var entries []navigationEntry
	for i, example := range controlCatalog {
		if newPages["controls/"+example.name] {
			entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i)})
		}
	}
	for i, example := range mygoCatalog {
		if newPages["mygo/"+example.name] {
			entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i)})
		}
	}
	return entries
}
