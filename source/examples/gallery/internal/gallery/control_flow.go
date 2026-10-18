package gallery

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/flow"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

type flowDemoState struct {
	view             flow.State
	selected, opened string
	added            []flowAddition
	leftToRight      bool
	empty            bool
}

// flowAddition is a step added after another, to watch the layout adapt.
type flowAddition struct{ id, after string }

type flowKind uint8

const (
	flowTask flowKind = iota
	flowExperiment
	flowFinding
)

type flowStatus uint8

const (
	flowNoStatus flowStatus = iota
	flowDone
	flowRunning
	flowPlanned
	flowWaiting
)

// flowStep is a step of the sample plan; title and detail are message keys.
type flowStep struct {
	id            string
	kind          flowKind
	status        flowStatus
	title, detail string
	experiment    int
	notes, files  int
	muted         bool
}

type flowLink struct {
	from, to string
	dashed   bool
	label    string
}

var flowSteps = []flowStep{
	{id: "goal", status: flowRunning, title: "flowGoal", detail: "flowGoalDetail", notes: 3},
	{id: "audit", status: flowDone, title: "flowAudit", detail: "flowAuditDetail", notes: 2, files: 1},
	{id: "content", status: flowDone, title: "flowContent", detail: "flowContentDetail", notes: 1},
	{id: "navigation", kind: flowExperiment, status: flowRunning, title: "flowNavigation", detail: "flowNavigationDetail", experiment: 1},
	{id: "speed", status: flowDone, title: "flowSpeed", detail: "flowSpeedDetail", files: 2},
	{id: "menus", kind: flowFinding, title: "flowMenus", detail: "flowMenusDetail"},
	{id: "retest", status: flowWaiting, title: "flowRetest", detail: "flowRetestDetail", muted: true},
	{id: "images", kind: flowFinding, title: "flowImages", detail: "flowImagesDetail"},
	{id: "review", status: flowRunning, title: "flowReview", detail: "flowReviewDetail", notes: 1},
	{id: "ready", kind: flowFinding, title: "flowReady", detail: "flowReadyDetail"},
	{id: "notes", status: flowPlanned, title: "flowNotes", detail: "flowNotesDetail", muted: true},
}

var flowLinks = []flowLink{
	{from: "goal", to: "audit"},
	{from: "goal", to: "content"},
	{from: "goal", to: "speed"},
	{from: "goal", to: "review"},
	{from: "audit", to: "navigation", dashed: true, label: "flowInforms"},
	{from: "audit", to: "speed", label: "flowDependsOn"},
	{from: "navigation", to: "menus"},
	{from: "navigation", to: "retest", dashed: true},
	{from: "speed", to: "images"},
	{from: "content", to: "review"},
	{from: "menus", to: "review", dashed: true},
	{from: "images", to: "review", dashed: true, label: "flowInforms"},
	{from: "review", to: "ready"},
	{from: "review", to: "notes", dashed: true},
}

func flowDemo(s *controlState) *flowDemoState {
	if s.flow == nil {
		s.flow = &flowDemoState{}
	}
	return s.flow
}

// flowSample maps the plan, and the steps added to it, onto diagram data
// with the titles of its nodes.
func flowSample(c *ui.Context, d *flowDemoState) (flow.Data, map[string]string) {
	var data flow.Data
	titles := map[string]string{}
	add := func(step flowStep, title, detail string) {
		titles[step.id] = title
		data.Nodes = append(data.Nodes, flow.Node{
			ID: step.id, Label: title, Muted: step.muted,
			Content: func(c *ui.Context) { flowCard(c, step, title, detail) },
		})
	}
	for _, step := range flowSteps {
		add(step, tr(c, step.title), tr(c, step.detail))
	}
	for _, l := range flowLinks {
		e := flow.Edge{Source: l.from, Target: l.to, Dashed: l.dashed}
		if l.label != "" {
			e.Label = tr(c, l.label)
		}
		data.Edges = append(data.Edges, e)
	}
	for i, a := range d.added {
		add(flowStep{id: a.id, status: flowPlanned}, msg(c).FlowNewStep(i+1), tr(c, "flowNewStepDetail"))
		data.Edges = append(data.Edges, flow.Edge{Source: a.after, Target: a.id})
	}
	return data, titles
}

// flowCard draws a step as a card with its status, kind and activity.
func flowCard(c *ui.Context, step flowStep, title, detail string) {
	t := c.Theme()
	tile, icon := t.Accent, lucide.ListChecks
	switch step.kind {
	case flowExperiment:
		tile, icon = t.Warning, lucide.FlaskConical
	case flowFinding:
		tile, icon = t.Success, lucide.Lightbulb
	}
	if step.id == "goal" {
		icon = lucide.Target
	}
	opts := flow.CardOptions{
		Title: title, Description: detail,
		Leading: func(c *ui.Context) {
			ui.Box(c).Size(30, 30).Radius(8).Background(tile.Alpha(.16)).TextColor(tile).Center().Children(func() {
				components.Icon(c, icon).Size(16, 16)
			})
		},
	}
	if step.status == flowRunning {
		opts.Border = t.Accent.Alpha(.6)
	}
	if step.status != flowNoStatus {
		opts.Header = func(c *ui.Context) {
			color, icon, label := t.Success, lucide.CircleCheck, "flowDone"
			switch step.status {
			case flowRunning:
				color, icon, label = t.Accent, lucide.LoaderCircle, "flowRunning"
			case flowPlanned:
				color, icon, label = t.TextMuted, lucide.CircleDashed, "flowPlanned"
			case flowWaiting:
				color, icon, label = t.Warning, lucide.Clock, "flowWaiting"
			}
			ui.Row(c).Gap(t.Space(1)).TextColor(color).Children(func() {
				components.Icon(c, icon).Size(12, 12)
				components.Text(c, tr(c, label)).FontWeight(500)
			})
			ui.Spacer(c)
			if step.experiment > 0 {
				components.Text(c, msg(c).FlowExperiment(step.experiment)).TextColor(t.TextMuted)
			}
		}
	}
	if step.notes > 0 || step.files > 0 {
		opts.Footer = func(c *ui.Context) {
			count := func(icon func() *ui.SVG, n int) {
				if n == 0 {
					return
				}
				ui.Row(c).Gap(t.Space(1)).Children(func() {
					components.Icon(c, icon).Size(12, 12)
					components.Text(c, strconv.Itoa(n))
				})
			}
			count(lucide.MessageSquare, step.notes)
			count(lucide.Paperclip, step.files)
		}
	}
	flow.Card(c, opts)
}

func controlFlow(c *ui.Context, s *controlState, _ string) {
	d := flowDemo(s)
	t := c.Theme()
	var data flow.Data
	var titles map[string]string
	if !d.empty {
		data, titles = flowSample(c, d)
	}
	direction := flow.TopToBottom
	if d.leftToRight {
		direction = flow.LeftToRight
	}
	flow.View(c, &d.view, data, flow.Options{
		Label:     msg(c).FlowLabel(len(data.Nodes), len(data.Edges)),
		EmptyText: tr(c, "flowEmpty"),
		Direction: direction,
		Selected:  d.selected,
		Select:    func(id string) { d.selected = id },
		Open:      func(id string) { d.opened = id },
	}).Grow(1).MinHeight(240)
	ui.Row(c).Gap(t.Space(4)).TextColor(t.TextMuted).Children(func() {
		if title, ok := titles[d.selected]; ok {
			components.Text(c, msg(c).FlowSelected(title))
		}
		if title, ok := titles[d.opened]; ok {
			components.Text(c, msg(c).GraphOpened(title))
		}
	})
}

// flowToolbar holds the diagram's sample actions in the page titlebar.
func flowToolbar(c *ui.Context, s *controlState) {
	d := flowDemo(s)
	if sidebar.IconButton(c, lucide.Plus(), tr(c, "flowAdd")).Disabled(d.empty).Clicked() {
		// A step follows the selected one, or the previous addition.
		after := "review"
		if n := len(d.added); n > 0 {
			after = d.added[n-1].id
		}
		if d.selected != "" {
			after = d.selected
		}
		id := fmt.Sprintf("added-%d", len(d.added)+1)
		d.added = append(d.added, flowAddition{id: id, after: after})
		d.selected = id
	}
	if len(d.added) > 0 && sidebar.IconButton(c, lucide.RotateCcw(), tr(c, "flowReset")).Clicked() {
		d.added = nil
		d.selected = ""
	}
	icon, label := lucide.ArrowRightLeft(), "flowLeftToRight"
	if d.leftToRight {
		icon, label = lucide.ArrowDownUp(), "flowTopToBottom"
	}
	if sidebar.IconButton(c, icon, tr(c, label)).Clicked() {
		d.leftToRight = !d.leftToRight
	}
	if sidebar.IconButton(c, lucide.ZoomOut(), tr(c, "flowZoomOut")).Disabled(d.empty).Clicked() {
		d.view.ZoomOut()
	}
	if sidebar.IconButton(c, lucide.ZoomIn(), tr(c, "flowZoomIn")).Disabled(d.empty).Clicked() {
		d.view.ZoomIn()
	}
	if sidebar.IconButton(c, lucide.Scan(), tr(c, "graphFit")).Disabled(d.empty).Clicked() {
		d.view.Fit()
	}
	emptyStateToggle(c, &d.empty)
}
