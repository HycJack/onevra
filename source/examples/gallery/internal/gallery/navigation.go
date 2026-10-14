package gallery

// Component sections describe purpose, never implementation ownership.
var navigationGroups = []string{
	"Plugins", "Actions & navigation", "Text input", "Selection & values",
	"Date & time", "Color", "Content display", "Code & diffs", "Charts", "Layout & containers",
	"Feedback & overlays", "Games",
}

// newPages are the routes added in the latest published version or since; the
// sidebar highlights them. Add new pages alongside the existing entries. When
// publishing a version, keep the pages it ships and remove those from earlier
// versions.
var newPages = map[string]bool{
	"controls/Relationship graph": true,
}

type navigationEntry struct {
	label string
	page  galleryPage
}

type navigationGroup struct {
	label   string
	entries []navigationEntry
}

func navigationCatalog() []navigationGroup {
	var groups []navigationGroup
	add := func(label string, entries []navigationEntry) {
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label, entries})
		}
	}
	add("Overview", []navigationEntry{{"Welcome", welcomePage}, {"Component status", componentStatusPage}})
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
