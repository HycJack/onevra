package gallery

// Component sections describe purpose, never implementation ownership.
var navigationGroups = []string{
	"Plugins", "Actions & navigation", "Text input", "Selection & values",
	"Date & time", "Color", "Content display", "Code & diffs", "Charts", "Layout & containers",
	"Feedback & overlays", "Games",
}

// newPages are the routes added in the latest published version or since; the
// top New components group lists them in place of their purpose groups. Add
// new pages alongside the existing entries. When
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

	"controls/Split panes":  true,
	"mygo/Gradients & blur": true,
	"mygo/Shaders":          true,
	"mygo/Live chart":       true,
}

type navigationEntry struct {
	label string
	page  galleryPage
	// mygo marks the single sidebar entry that hosts every MyGo page behind
	// its own in-page list; page is then only the first of them.
	mygo bool
}

// mygoGroupLabel names the sidebar group of the MyGo entry, and the command
// palette group of the pages it hosts.
const mygoGroupLabel = "MyGo"

// newGroupLabel names the group of the pages in newPages.
const newGroupLabel = "New components"

type navigationGroup struct {
	label   string
	entries []navigationEntry
}

func navigationCatalog() []navigationGroup {
	var groups []navigationGroup
	add := func(label string, entries []navigationEntry) {
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label: label, entries: entries})
		}
	}
	add("Overview", []navigationEntry{{"Welcome", welcomePage, false}})
	if entries := newNavigationEntries(false); len(entries) > 0 {
		groups = append(groups, navigationGroup{label: newGroupLabel, entries: entries})
	}
	add("Icon library", []navigationEntry{{"Lucide", lucidePage, false}})
	add(mygoGroupLabel, []navigationEntry{{"Official components", trayPage, true}})
	for _, group := range navigationGroups {
		var entries []navigationEntry
		for i, example := range controlCatalog {
			if example.group == group && !newPages["controls/"+example.name] {
				entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i), false})
			}
		}
		add(group, entries)
	}
	return groups
}

// mygoNavigationGroups lists the pages of the MyGo entry for its in-page
// list: its own New components group, then purposes in the sidebar's order.
func mygoNavigationGroups() []navigationGroup {
	var groups []navigationGroup
	if entries := newNavigationEntries(true); len(entries) > 0 {
		groups = append(groups, navigationGroup{label: newGroupLabel, entries: entries})
	}
	for _, group := range append([]string{"Desktop"}, navigationGroups...) {
		var entries []navigationEntry
		if group == "Desktop" {
			entries = []navigationEntry{{"Tray", trayPage, false}, {"Dock menu", dockPage, false}, {"System notifications", notificationsPage, false}}
		}
		for i, example := range mygoCatalog {
			if example.group == group && !newPages["mygo/"+example.name] {
				entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i), false})
			}
		}
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label: group, entries: entries})
		}
	}
	return groups
}

// isMyGoPage reports pages behind the MyGo entry: desktop integrations and
// the MyGo catalog, which directly follow one another.
func isMyGoPage(page galleryPage) bool {
	return page >= trayPage && page < controlPageStart()
}

// newNavigationEntries collects the pages in newPages in catalog order, either
// Onevra's for the sidebar or MyGo's for the MyGo entry's in-page list.
func newNavigationEntries(mygo bool) []navigationEntry {
	var entries []navigationEntry
	if !mygo {
		for i, example := range controlCatalog {
			if newPages["controls/"+example.name] {
				entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i), false})
			}
		}
		return entries
	}
	for i, example := range mygoCatalog {
		if newPages["mygo/"+example.name] {
			entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i), false})
		}
	}
	return entries
}
