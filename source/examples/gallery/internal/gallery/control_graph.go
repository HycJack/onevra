package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/graph"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

type graphDemoState struct {
	view                  graph.State
	attachments, settings bool
	empty                 bool
	added                 int
	opened                string
	cache                 graphCache
}

// graphCache keeps the sample between builds, so thousands of added notes
// are not regenerated for every frame.
type graphCache struct {
	added             int
	attachments, dark bool
	newTitle          string
	data              graph.Data
	titles            map[string]string
}

// graphNote is a document of the sample library; its path's first segment
// is its directory, as in a Markdown notes library.
type graphNote struct {
	path, title string
	links       []string
	attachments []string
}

var graphNotes = []graphNote{
	{path: "index.md", title: "Home", links: []string{"projects/index.md", "areas/index.md", "journal/index.md", "resources/index.md"}},
	{path: "scratchpad.md", title: "Scratchpad"},
	{path: "projects/index.md", title: "Projects", links: []string{"projects/garden.md", "projects/website.md", "projects/move.md"}},
	{path: "projects/garden.md", title: "Community garden", links: []string{"areas/health.md", "resources/plants.md"}, attachments: []string{"plot-plan.png"}},
	{path: "projects/website.md", title: "Personal website redesign", links: []string{"resources/typography.md"}, attachments: []string{"wireframe.pdf", "palette.png"}},
	{path: "projects/move.md", title: "Moving checklist", links: []string{"areas/finance.md"}},
	{path: "areas/index.md", title: "Areas", links: []string{"areas/health.md", "areas/finance.md", "areas/home.md"}},
	{path: "areas/health.md", title: "Health"},
	{path: "areas/finance.md", title: "Finance", attachments: []string{"budget.xlsx"}},
	{path: "areas/home.md", title: "Home maintenance", links: []string{"projects/move.md"}},
	{path: "journal/index.md", title: "Journal", links: []string{"journal/monday.md", "journal/tuesday.md", "journal/weekly.md"}},
	{path: "journal/monday.md", title: "Monday", links: []string{"projects/garden.md"}},
	{path: "journal/tuesday.md", title: "Tuesday", links: []string{"projects/website.md"}, attachments: []string{"sketch.jpg"}},
	{path: "journal/weekly.md", title: "Weekly review", links: []string{"areas/finance.md", "areas/health.md"}},
	{path: "resources/index.md", title: "Resources", links: []string{"resources/plants.md", "resources/typography.md", "resources/reading.md"}},
	{path: "resources/plants.md", title: "Plant care", attachments: []string{"seasons.png"}},
	{path: "resources/typography.md", title: "Typography notes"},
	{path: "resources/reading.md", title: "Reading list", links: []string{"journal/weekly.md"}},
}

func graphDirectory(path string) string {
	directory, _, nested := strings.Cut(path, "/")
	if !nested {
		return ""
	}
	return directory
}

// graphSample maps the library onto graph data: directories become groups,
// directory indexes own their notes, wiki links are references and
// attachments sit beside their note.
func graphSample(c *ui.Context, s *graphDemoState) (graph.Data, map[string]string) {
	dark, newTitle := c.Theme().Dark, tr(c, "gallery.graph_new_note")
	k := &s.cache
	if k.titles == nil || k.added != s.added || k.attachments != s.attachments || k.dark != dark || k.newTitle != newTitle {
		*k = graphCache{added: s.added, attachments: s.attachments, dark: dark, newTitle: newTitle}
		k.data = graphLibrary(s, dark, newTitle)
		k.titles = make(map[string]string, len(k.data.Nodes))
		for _, n := range k.data.Nodes {
			k.titles[n.ID] = n.Label
		}
	}
	return k.data, k.titles
}

// graphLibrary builds the sample with newTitle formatting added notes' titles.
func graphLibrary(s *graphDemoState, dark bool, newTitle string) graph.Data {
	attachment, attachmentLabel := ui.Hex("#e0b35a"), ui.Hex("#8a6a2a")
	referenceOpacity, attachmentOpacity := float32(.16), float32(.13)
	if dark {
		attachment, attachmentLabel = ui.Hex("#d4b06a"), ui.Hex("#ead7a8")
		referenceOpacity, attachmentOpacity = .2, .17
	}
	hierarchy := graph.EdgeStyle{Distance: 160, HoverOpacity: .4}
	reference := graph.EdgeStyle{Distance: 200, Width: .85, Opacity: referenceOpacity, HoverOpacity: .35}
	attached := graph.EdgeStyle{Distance: 64, Width: .8, Opacity: attachmentOpacity, Color: attachment}

	notes := append(make([]graphNote, 0, len(graphNotes)+s.added), graphNotes...)
	base := len(graphNotes) - 2
	for i := range s.added {
		// The first notes join the sample's folders; later ones branch from
		// earlier added notes, with a few links back across the library.
		parent := notes[2+i%base]
		var links []string
		if i >= base {
			parent = notes[len(graphNotes)+(i-base)/3]
			if i%7 == 0 {
				links = append(links, graphNotes[2+i%base].path)
			}
		}
		notes = append(notes, graphNote{
			path:  fmt.Sprintf("%s/new-%d.md", graphDirectory(parent.path), i+1),
			title: fmt.Sprintf(newTitle, i+1),
			links: append([]string{parent.path}, links...),
		})
	}
	var data graph.Data
	for _, note := range notes {
		directory := graphDirectory(note.path)
		data.Nodes = append(data.Nodes, graph.Node{ID: note.path, Label: note.title, Group: directory, Openable: true})
		for _, target := range note.links {
			style := reference
			if strings.HasSuffix(note.path, "index.md") && (graphDirectory(target) == directory || directory == "") {
				style = hierarchy
			}
			data.Edges = append(data.Edges, graph.Edge{Source: note.path, Target: target, Style: style})
		}
		if !s.attachments {
			continue
		}
		for _, name := range note.attachments {
			id := note.path + "#" + name
			data.Nodes = append(data.Nodes, graph.Node{ID: id, Label: name, Group: directory, Size: 8, Color: attachment, LabelColor: attachmentLabel})
			data.Edges = append(data.Edges, graph.Edge{Source: note.path, Target: id, Style: attached})
		}
	}
	return data
}

// graphBulkCounts are the larger additions offered for trying the graph's
// performance.
var graphBulkCounts = []int{100, 1000}

func controlGraph(c *ui.Context, s *controlState, _ string) {
	if s.graph == nil {
		s.graph = &graphDemoState{attachments: true}
	}
	d := s.graph
	t := c.Theme()
	components.Text(c, tr(c, "gallery.graph_intro")).TextColor(t.TextMuted).Wrap()
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(2)).Children(func() {
		if components.Button(c, tr(c, "gallery.graph_add"), components.ButtonOutline).Clicked() {
			d.added++
		}
		for _, count := range graphBulkCounts {
			if components.Button(c, fmt.Sprintf(tr(c, "gallery.graph_add_many"), count), components.ButtonOutline).Clicked() {
				d.added += count
			}
		}
		if d.added > 0 && components.Button(c, tr(c, "gallery.graph_reset"), components.ButtonGhost).Clicked() {
			d.added = 0
		}
		if components.Button(c, tr(c, "gallery.graph_fit"), components.ButtonOutline).Clicked() {
			d.view.Fit()
		}
		label := "gallery.chart_empty"
		if d.empty {
			label = "gallery.chart_restore"
		}
		if components.Button(c, tr(c, label), components.ButtonGhost).Clicked() {
			d.empty = !d.empty
		}
	})
	var data graph.Data
	var titles map[string]string
	if !d.empty {
		data, titles = graphSample(c, d)
	}
	ui.Column(c).FillWidth().Grow(1).MinWidth(0).MinHeight(240).Radius(t.Radius*1.5).Background(t.Surface).Border(1, t.Border.Alpha(.5)).Clip().Children(func() {
		graph.View(c, &d.view, data, graph.Options{
			Label:     fmt.Sprintf(tr(c, "gallery.graph_label"), len(data.Nodes), len(data.Edges)),
			EmptyText: tr(c, "gallery.graph_empty"),
			Open:      func(id string) { d.opened = id },
			Overlay:   func(c *ui.Context) { graphSettings(c, d) },
		})
	})
	opened := tr(c, "gallery.graph_hint")
	if title, ok := titles[d.opened]; ok {
		opened = fmt.Sprintf(tr(c, "gallery.graph_opened"), title)
	}
	components.Text(c, opened).TextColor(t.TextMuted)
}

func graphSettings(c *ui.Context, d *graphDemoState) {
	t := c.Theme()
	label := tr(c, "gallery.graph_settings")
	anchor := components.Button(c, "", components.ButtonSecondary).Label(label).Padding(t.Space(1.5))
	anchor.Children(func() { components.Icon(c, lucide.Ellipsis).Size(16, 16) })
	if anchor.Clicked() {
		d.settings = !d.settings
	}
	components.Popover(c, anchor, &d.settings, func(panel ui.Element) {
		panel.Label(label).Width(288)
		ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(t.Space(3)).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				components.Text(c, tr(c, "gallery.graph_attachments"))
				components.Text(c, tr(c, "gallery.graph_attachments_hint")).FontSize(t.FontSize - 1).TextColor(t.TextMuted).Wrap()
			})
			components.Switch(c, &d.attachments, tr(c, "gallery.graph_attachments"), components.SwitchOptions{HideLabel: true})
		})
	})
}
