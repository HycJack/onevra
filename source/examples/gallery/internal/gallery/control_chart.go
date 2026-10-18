package gallery

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

type chartDemoState struct {
	// Each chart keeps its own selection; the data actions apply to all.
	selections map[string]*chart.State
	empty      bool
}

// selection returns a chart's selection, starting with none.
func (d *chartDemoState) selection(name string) *chart.State {
	if d.selections == nil {
		d.selections = make(map[string]*chart.State)
	}
	if d.selections[name] == nil {
		d.selections[name] = &chart.State{Index: -1}
	}
	return d.selections[name]
}

var chartViews = map[string]func(*ui.Context, *chart.State, chart.Data, chart.Options) ui.Element{
	"Bar chart": chart.Bar, "Line chart": chart.Line, "Area chart": chart.Area, "Stacked bar chart": chart.StackedBar,
}

func chartDemo(s *controlState) *chartDemoState {
	if s.chart == nil {
		s.chart = &chartDemoState{}
	}
	return s.chart
}

// chartToolbar holds a chart page's data actions in the page titlebar.
func chartToolbar(c *ui.Context, s *controlState) {
	d := chartDemo(s)
	if sidebar.IconButton(c, lucide.RefreshCw(), tr(c, "chartUpdate")).Disabled(d.empty).Clicked() {
		s.on = !s.on
	}
	emptyStateToggle(c, &d.empty)
}

// emptyStateToggle switches a sample between its data and its empty state.
func emptyStateToggle(c *ui.Context, empty *bool) {
	icon, label := lucide.Eraser(), "chartEmpty"
	if *empty {
		icon, label = lucide.Undo2(), "chartRestore"
	}
	if sidebar.IconButton(c, icon, tr(c, label)).Clicked() {
		*empty = !*empty
	}
}

// controlChart shows the chart a variant names.
func controlChart(c *ui.Context, s *controlState, name string) {
	d := chartDemo(s)
	t := c.Theme()
	ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(4)).Gap(t.Space(4)).
		Radius(t.Radius*1.5).Background(t.Surface).Border(1, t.Border.Alpha(.5)).Children(func() {
		opts := chart.Options{
			Label: tr(c, "chartSelection"), Description: tr(c, "chartKeys"),
			EmptyText: tr(c, "chartNoData"), PlotHeight: 220,
			FormatValue: func(v float64) string { return strconv.FormatFloat(v, 'f', 0, 64) },
		}
		values := []float64{38, 52, 46, 68, 61, 82}
		if s.on {
			values = []float64{62, 43, 71, 54, 78, 65}
		}
		selection := d.selection(name)
		if name == "Split bar" {
			// The data action switches between two weeks of traffic.
			labels := []string{"chartOrganic", "chartDirect", "chartPaid", "chartSocial", "chartReferral", "chartOther"}
			sessions := []float64{41260, 24610, 15830, 8420, 5310, 3790}
			if s.on {
				sessions = []float64{36120, 27940, 11250, 13870, 4980, 2610}
			}
			var slices []chart.Slice
			if !d.empty {
				for i, label := range labels {
					slices = append(slices, chart.Slice{ID: label, Label: tr(c, label), Value: sessions[i]})
				}
			}
			chart.SplitBar(c, selection, slices, chart.Options{
				Label: tr(c, "chartTraffic"), Description: tr(c, "chartTrafficSessions"), EmptyText: tr(c, "chartNoData"),
			})
			return
		}
		if name == "Donut chart" {
			components.Text(c, tr(c, "chartDistribution")).Bold()
			labels := []string{"chartDesign", "chartDevelopment", "chartResearch", "chartPlanning"}
			var slices []chart.Slice
			if !d.empty {
				for i, label := range labels {
					slices = append(slices, chart.Slice{ID: label, Label: tr(c, label), Value: values[i]})
				}
			}
			chart.Donut(c, selection, slices, opts)
			return
		}
		components.Text(c, tr(c, "chartDaily")).Bold()
		var data chart.Data
		if !d.empty {
			for _, day := range []string{"chartMon", "chartTue", "chartWed", "chartThu", "chartFri", "chartSat"} {
				data.Labels = append(data.Labels, tr(c, day))
			}
			data.Series = []chart.Series{
				{ID: "current", Label: tr(c, "chartCurrent"), Values: values},
				{ID: "previous", Label: tr(c, "chartPrevious"), Values: []float64{26, 34, 28, 48, 43, 57}},
			}
		}
		chartViews[name](c, selection, data, opts)
	})
}
