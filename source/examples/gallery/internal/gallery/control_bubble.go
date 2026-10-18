package gallery

import (
	"math/rand/v2"
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// bubbleLabels name the sample's categories, of which the first count show.
var bubbleLabels = []string{
	"chartDesign", "chartDevelopment", "chartResearch", "chartPlanning", "chartReview",
	"chartSupport", "chartTesting", "chartOperations", "chartDocumentation", "chartMarketing",
}

const bubbleInitialCount = 7

type bubbleDemoState struct {
	selection chart.State
	count     int // categories shown
	round     int // data updates
	empty     bool
}

func bubbleDemo(s *controlState) *bubbleDemoState {
	if s.bubble == nil {
		s.bubble = &bubbleDemoState{selection: chart.State{Index: -1}, count: bubbleInitialCount}
	}
	return s.bubble
}

// bubbleValues returns the sample values of an update round, the same for
// every category added later, so updates and additions are independent.
func bubbleValues(round int) []float64 {
	values := []float64{38, 82, 46, 20, 27, 14, 61, 9, 33, 52}
	if round == 0 {
		return values
	}
	r := rand.New(rand.NewPCG(uint64(round), 1))
	for i := range values {
		values[i] = float64(5 + r.IntN(95))
	}
	return values
}

// bubbleToolbar holds the bubble sample's data actions in the page titlebar.
func bubbleToolbar(c *ui.Context, s *controlState) {
	d := bubbleDemo(s)
	if sidebar.IconButton(c, lucide.RefreshCw(), tr(c, "chartUpdate")).Disabled(d.empty).Clicked() {
		d.round++
	}
	if sidebar.IconButton(c, lucide.Plus(), tr(c, "bubbleAdd")).Disabled(d.empty || d.count == len(bubbleLabels)).Clicked() {
		d.count++
	}
	if sidebar.IconButton(c, lucide.Minus(), tr(c, "bubbleRemove")).Disabled(d.empty || d.count == 1).Clicked() {
		d.count--
	}
	emptyStateToggle(c, &d.empty)
}

func controlBubbleChart(c *ui.Context, s *controlState, _ string) {
	d := bubbleDemo(s)
	t := c.Theme()
	ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(4)).Gap(t.Space(4)).
		Radius(t.Radius*1.5).Background(t.Surface).Border(1, t.Border.Alpha(.5)).Children(func() {
		components.Text(c, tr(c, "chartDistribution")).Bold()
		var slices []chart.Slice
		if !d.empty {
			values := bubbleValues(d.round)
			for i, label := range bubbleLabels[:d.count] {
				slices = append(slices, chart.Slice{ID: label, Label: tr(c, label), Value: values[i]})
			}
		}
		chart.Bubble(c, &d.selection, slices, chart.Options{
			Label: tr(c, "chartSelection"), Description: tr(c, "chartKeys"),
			EmptyText: tr(c, "chartNoData"), PlotHeight: 320,
			FormatValue: func(v float64) string { return strconv.FormatFloat(v, 'f', 0, 64) },
		})
	})
}
