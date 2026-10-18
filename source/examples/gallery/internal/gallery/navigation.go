package gallery

// Component sections describe purpose, never implementation ownership.
var navigationGroups = []string{
	"Shaders", "Conversation", "Plugins", "Actions & navigation", "Text input", "Selection & values",
	"Date & time", "Color", "Content display", "Code & diffs", "Charts", "Layout & containers",
	"Feedback & overlays", "Games",
}

// newPages are the routes added in the latest published version or since; the
// top New components group lists them in place of their purpose groups. Add
// new pages alongside the existing entries. When
// publishing a version, keep the pages it ships and remove those from earlier
// versions.
var newPages = map[string]bool{
	"controls/Carousel":           true,
	"controls/Split panes":        true,
	"controls/Shader backgrounds": true,
	"controls/Bubble chart":       true,
	"controls/TagInput":           true,
	"controls/HoldToConfirm":      true,
	"controls/Skeleton":           true,
	"controls/Avatars":            true,
	"controls/Treemap":            true,
	"mygo/Global shortcuts":       true,
	"mygo/Gradients & blur":       true,
	"mygo/Shaders":                true,
	"mygo/Live chart":             true,
}

// updatedPages are the routes whose components changed in the latest published
// version or since; the Updated components group lists them in place of their
// purpose groups. Keep this list disjoint from newPages. When publishing a
// version, keep the pages it changes and remove those from earlier versions.
var updatedPages = map[string]bool{
	"controls/Badge":      true,
	"controls/Chat":       true,
	"controls/Dialogs":    true,
	"controls/Disclosure": true,
	"controls/Honeycomb":  true,
	"controls/Sidebar":    true,
	"controls/Toast":      true,
	"controls/Tree":       true,
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

// newGroupLabel names the group of the pages in newPages, and
// updatedGroupLabel the group of the pages in updatedPages.
const (
	newGroupLabel     = "New components"
	updatedGroupLabel = "Updated components"
)

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
	add(newGroupLabel, newNavigationEntries(false))
	add(updatedGroupLabel, updatedNavigationEntries(false))
	add("Icon library", []navigationEntry{{"Lucide", lucidePage, false}})
	add(mygoGroupLabel, []navigationEntry{{"Official components", notificationsPage, true}})
	for _, group := range navigationGroups {
		var entries []navigationEntry
		for i, example := range controlCatalog {
			if example.group == group && !newPages["controls/"+example.name] && !updatedPages["controls/"+example.name] {
				entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i), false})
			}
		}
		add(group, entries)
	}
	return groups
}

// mygoNavigationGroups lists the pages of the MyGo entry for its in-page
// list: its own New and Updated components groups, then purposes in the
// sidebar's order.
func mygoNavigationGroups() []navigationGroup {
	var groups []navigationGroup
	add := func(label string, entries []navigationEntry) {
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label: label, entries: entries})
		}
	}
	add(newGroupLabel, newNavigationEntries(true))
	add(updatedGroupLabel, updatedNavigationEntries(true))
	for _, group := range append([]string{"Desktop"}, navigationGroups...) {
		var entries []navigationEntry
		if group == "Desktop" {
			entries = []navigationEntry{{"System notifications", notificationsPage, false}}
		}
		for i, example := range mygoCatalog {
			if example.group == group && !newPages["mygo/"+example.name] && !updatedPages["mygo/"+example.name] {
				entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i), false})
			}
		}
		if len(entries) > 0 {
			groups = append(groups, navigationGroup{label: group, entries: entries})
		}
	}
	return groups
}

// isMyGoPage reports pages behind the MyGo entry: System notifications and
// the MyGo catalog, which directly follow one another.
func isMyGoPage(page galleryPage) bool {
	return page >= notificationsPage && page < controlPageStart()
}

// newNavigationEntries collects the pages in newPages in catalog order, either
// Onevra's for the sidebar or MyGo's for the MyGo entry's in-page list.
func newNavigationEntries(mygo bool) []navigationEntry {
	return navigationEntries(newPages, mygo)
}

// updatedNavigationEntries collects the pages in updatedPages in catalog order.
func updatedNavigationEntries(mygo bool) []navigationEntry {
	return navigationEntries(updatedPages, mygo)
}

func navigationEntries(pages map[string]bool, mygo bool) []navigationEntry {
	var entries []navigationEntry
	if !mygo {
		for i, example := range controlCatalog {
			if pages["controls/"+example.name] {
				entries = append(entries, navigationEntry{example.label, controlPageStart() + galleryPage(i), false})
			}
		}
		return entries
	}
	for i, example := range mygoCatalog {
		if pages["mygo/"+example.name] {
			entries = append(entries, navigationEntry{example.name, mygoPageStart + galleryPage(i), false})
		}
	}
	return entries
}
