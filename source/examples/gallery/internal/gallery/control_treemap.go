package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/treemap"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

type treemapDemoState struct {
	view treemap.State
	// updated shows the sample after a cleanup, with sizes changed.
	updated, empty bool
	opened         string
}

func treemapDemo(s *controlState) *treemapDemoState {
	if s.treemap == nil {
		s.treemap = &treemapDemoState{}
	}
	return s.treemap
}

// treemapSample is a simulated home folder's disk usage in megabytes, with
// days since each file was last opened as the shading measure. File and
// folder names stay as they are on disk, so only the captions are translated.
// After a cleanup, the old installers and caches are gone and new photos have
// arrived.
func treemapSample(c *ui.Context, updated bool) []treemap.Node {
	reclaimable, clutter := tr(c, "treemapReclaimable"), ui.Hex("#e8a33d")
	file := func(path string, mb, days float64) treemap.Node {
		return treemap.Node{ID: path, Label: path[strings.LastIndex(path, "/")+1:], Value: mb, Measure: days}
	}
	// Reclaimable files stand out with a dashed outline over a flat fill.
	stale := func(n treemap.Node) treemap.Node {
		n.Caption, n.Color = reclaimable, clutter
		n.Style = treemap.Style{Background: clutter.Alpha(.16), Border: clutter, BorderWidth: 1.5, BorderStyle: ui.BorderDashed}
		n.HoverStyle = treemap.Style{Background: clutter.Alpha(.3), BorderWidth: 2}
		return n
	}
	// A dashed outline marks a whole top-level folder as one to review.
	outlined := func(n treemap.Node) treemap.Node {
		n.Style = treemap.Style{Border: c.Theme().TextMuted.Alpha(.7), BorderWidth: 1.5, BorderStyle: ui.BorderDashed}
		n.HoverStyle = treemap.Style{Border: c.Theme().Text, BorderWidth: 2}
		return n
	}
	folder := func(path string, children ...treemap.Node) treemap.Node {
		return treemap.Node{ID: path, Label: path[strings.LastIndex(path, "/")+1:], Children: children}
	}
	installers, caches, trips := 1840.0, 920.0, 0.0
	if updated {
		installers, caches, trips = 0, 210, 1360
	}
	return []treemap.Node{
		folder("Videos",
			file("Videos/Family reunion.mov", 2140, 210),
			file("Videos/Screen recordings", 640, 34),
			folder("Videos/Clips", file("Videos/Clips/Birthday.mp4", 380, 300), file("Videos/Clips/Concert.mp4", 290, 150))),
		folder("Pictures",
			folder("Pictures/2025", file("Pictures/2025/Summer", 860, 90), file("Pictures/2025/Winter", 540, 280)),
			folder("Pictures/2026", file("Pictures/2026/Spring", 720, 12), file("Pictures/2026/Trip", trips, 1)),
			file("Pictures/Wallpapers", 180, 45)),
		outlined(folder("Downloads",
			stale(file("Downloads/Installers", installers, 160)),
			file("Downloads/Archive.zip", 640, 120),
			file("Downloads/Manual.pdf", 48, 60), file("Downloads/Slides.key", 96, 8))),
		folder("Documents",
			folder("Documents/Projects", file("Documents/Projects/Website", 420, 3), file("Documents/Projects/Garden", 160, 20),
				file("Documents/Projects/Budget.xlsx", 12, 6)),
			folder("Documents/Writing", file("Documents/Writing/Drafts", 140, 30), file("Documents/Writing/Notes", 36, 2)),
			file("Documents/Scans", 260, 200)),
		folder("Music",
			folder("Music/Albums", file("Music/Albums/Jazz", 680, 40), file("Music/Albums/Classical", 520, 90), file("Music/Albums/Folk", 210, 25)),
			file("Music/Podcasts", 330, 14)),
		folder("Library",
			stale(file("Library/Caches", caches, 70)),
			file("Library/Mail", 610, 1), file("Library/Fonts", 120, 180)),
	}
}

func formatMegabytes(mb float64) string {
	if mb >= 1000 {
		return fmt.Sprintf("%.1f GB", mb/1000)
	}
	return fmt.Sprintf("%.0f MB", mb)
}

func controlTreemap(c *ui.Context, s *controlState, _ string) {
	d := treemapDemo(s)
	t := c.Theme()
	var data []treemap.Node
	if !d.empty {
		data = treemapSample(c, d.updated)
	}
	m := msg(c)
	treemap.View(c, &d.view, data, treemap.Options{
		Label:         tr(c, "treemapHome"),
		Description:   tr(c, "treemapKeys"),
		EmptyText:     tr(c, "treemapEmpty"),
		FormatValue:   formatMegabytes,
		MeasureLabel:  tr(c, "treemapUnused"),
		FormatMeasure: func(days float64) string { return m.TreemapDays(int(days + .5)) },
		Open:          func(id string) { d.opened = id },
		// A folder's menu can also drill into it, and a file's opens it.
		Menu: func(id string) func(m *ui.Menu) {
			path := treemap.Path(data, id)
			if path == nil {
				return nil
			}
			return func(m *ui.Menu) {
				if m.Item(tr(c, "treemapOpen")).Chosen() {
					d.opened = id
				}
				if len(path[len(path)-1].Children) > 0 && m.Item(tr(c, "treemapShowContents")).Chosen() {
					d.view.Root, d.view.Active = id, ""
				}
			}
		},
	}).Grow(1).MinHeight(240)
	status := tr(c, "treemapKeys")
	if opened := treemap.Path(data, d.opened); opened != nil {
		status = m.TreemapOpened(opened[len(opened)-1].Label)
	}
	components.Text(c, status).TextColor(t.TextMuted)
}

// treemapToolbar holds the treemap's sample actions in the page titlebar.
func treemapToolbar(c *ui.Context, s *controlState) {
	d := treemapDemo(s)
	if sidebar.IconButton(c, lucide.RefreshCw(), tr(c, "chartUpdate")).Disabled(d.empty).Clicked() {
		d.updated = !d.updated
	}
	emptyStateToggle(c, &d.empty)
}
