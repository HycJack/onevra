package gallery

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func controlSidebar(c *ui.Context, s *controlState, variant string) {
	switch variant {
	case "File tree":
		controlSidebarFileTree(c, s)
		return
	case "Collapsible groups":
		controlSidebarCollapsible(c, s)
		return
	case "Custom leading":
		controlSidebarLeading(c, s)
		return
	}
	components.Sidebar(c).Children(func() {
		components.SidebarSection(c, tr(c, "locations"), func() {
			if components.SidebarItem(c, tr(c, "home"), lucide.House(), s.index == 0).Key("home").Clicked() {
				s.index = 0
			}
			if components.SidebarItem(c, tr(c, "files"), lucide.Folder(), s.index == 1).Key("files").Clicked() {
				s.index = 1
			}
		})
		components.SidebarSection(c, tr(c, "account"), func() {
			if components.SidebarItem(c, tr(c, "preferences"), nil, s.index == 2).Key("preferences").Clicked() {
				s.index = 2
			}
			// A missing avatar leaves the item without an icon.
			avatar, _ := galleryImage()
			var image ui.ImageSource
			if avatar != nil {
				image = avatar
			}
			if components.SidebarImageItem(c, tr(c, "profile"), image, s.index == 3).Key("profile").Clicked() {
				s.index = 3
			}
		})
	})
}

// sidebarFile is a node of the file tree example; folders have children.
type sidebarFile struct {
	name     string
	children []sidebarFile
}

var sidebarFiles = []sidebarFile{
	{"source", []sidebarFile{
		{"components", []sidebarFile{{"button.go", nil}, {"sidebar.go", nil}, {"tree.go", nil}}},
		{"main.go", nil},
		{"go.mod", nil},
	}},
	{"docs", []sidebarFile{{"build.md", nil}}},
	{"README.md", nil},
}

// The file tree places a Tree inside a SidebarSection, so folders expand in
// place like an editor's project explorer.
func controlSidebarFileTree(c *ui.Context, s *controlState) {
	folder := 0
	var files func([]sidebarFile)
	files = func(nodes []sidebarFile) {
		for _, node := range nodes {
			var open *bool
			var children func()
			if node.children != nil {
				open = &s.folders[folder]
				folder++
				children = func() { files(node.children) }
			}
			if components.TreeItem(c, node.name, open, s.choice == node.name, children).Clicked() {
				s.choice = node.name
			}
		}
	}
	components.Sidebar(c).Children(func() {
		components.SidebarSection(c, tr(c, "explorer"), func() {
			components.Tree(c, func() { files(sidebarFiles) }).Label(tr(c, "explorer"))
		})
	})
}

// Collapsible groups let a long sidebar hide sections, with an icon before
// and a count beside each heading.
func controlSidebarCollapsible(c *ui.Context, s *controlState) {
	group := func(label string, icon *ui.SVG, open *bool, items []string, icons []*ui.SVG) {
		components.Collapsible(c, tr(c, label), open, components.CollapsibleOptions{
			ArrowPosition: components.ArrowLeft,
			Leading:       func() { components.Icon(c, icon).Size(14, 14).Shrink(0) },
			Trailing:      func() { components.Badge(c, strconv.Itoa(len(items))) },
		}, func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(c.Theme().Space(1)).Children(func() {
				for i, item := range items {
					if components.SidebarItem(c, tr(c, item), icons[i], s.choice == item).Key(item).Clicked() {
						s.choice = item
					}
				}
			})
		}).Root.FillWidth().MinWidth(0)
	}
	components.Sidebar(c).Children(func() {
		group("favorites", lucide.Star(), &s.folders[0], []string{"home", "recent", "downloads"}, []*ui.SVG{lucide.House(), lucide.Clock(), lucide.Download()})
		group("locations", lucide.MapPin(), &s.folders[1], []string{"documents", "shared"}, []*ui.SVG{lucide.FileText(), lucide.Users()})
	})
}

// Custom leading content replaces the icon with application-built marks,
// such as workspace monograms or tag colors, and trailing content adds
// counts after the label.
func controlSidebarLeading(c *ui.Context, s *controlState) {
	t := c.Theme()
	item := func(key string, leading func(c *ui.Context, size float32)) {
		if components.SidebarItemWithLeading(c, tr(c, key), s.choice == key, leading).Key(key).Clicked() {
			s.choice = key
		}
	}
	counted := func(key string, count int, leading func(c *ui.Context, size float32)) {
		trailing := func(c *ui.Context, _ float32) {
			components.Text(c, strconv.Itoa(count)).FontSize(t.Rem(.8)).TextColor(t.TextMuted).Shrink(0)
		}
		if components.SidebarItemWithSlots(c, tr(c, key), s.choice == key, leading, trailing).Key(key).Clicked() {
			s.choice = key
		}
	}
	// A monogram shows the first letter of its localized name.
	monogram := func(color ui.Color, key string) func(*ui.Context, float32) {
		return func(c *ui.Context, size float32) {
			letter := string([]rune(tr(c, key))[:1])
			ui.Box(c).Size(size, size).Shrink(0).Radius(size * .25).Background(color).
				Justify(ui.Center).AlignItems(ui.Center).Children(func() {
				components.Text(c, letter).FontSize(size * .6).Bold().TextColor(ui.RGB(255, 255, 255))
			})
		}
	}
	dot := func(color ui.Color) func(*ui.Context, float32) {
		return func(c *ui.Context, size float32) {
			ui.Box(c).Size(size, size).Shrink(0).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(size*.5, size*.5).Radius(size * .25).Background(color)
			})
		}
	}
	components.Sidebar(c).Children(func() {
		components.SidebarSection(c, tr(c, "workspaces"), func() {
			item("personal2", monogram(t.Accent, "personal2"))
			item("workspaceTeam", monogram(ui.RGB(232, 121, 52), "workspaceTeam"))
			item("archive", monogram(ui.RGB(120, 120, 128), "archive"))
		})
		components.SidebarSection(c, tr(c, "tags"), func() {
			counted("tagImportant", 3, dot(ui.RGB(235, 77, 75)))
			counted("tagInProgress", 12, dot(ui.RGB(242, 178, 51)))
			counted("tagDone", 48, dot(ui.RGB(72, 187, 120)))
		})
	})
}
