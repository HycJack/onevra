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
	selection chart.State
	empty     bool
}

func controlBarChart(c *ui.Context, s *controlState, _ string) {
	controlChart(c, s, "Bar chart", chart.Bar)
}

func controlLineChart(c *ui.Context, s *controlState, _ string) {
	controlChart(c, s, "Line chart", chart.Line)
}

func controlAreaChart(c *ui.Context, s *controlState, _ string) {
	controlChart(c, s, "Area chart", chart.Area)
}

func controlStackedBarChart(c *ui.Context, s *controlState, _ string) {
	controlChart(c, s, "Stacked bar chart", chart.StackedBar)
}

func controlDonutChart(c *ui.Context, s *controlState, _ string) {
	controlChart(c, s, "Donut chart", nil)
}

func chartDemo(s *controlState) *chartDemoState {
	if s.chart == nil {
		s.chart = &chartDemoState{selection: chart.State{Index: -1}}
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

func controlChart(c *ui.Context, s *controlState, name string, render func(*ui.Context, *chart.State, chart.Data, chart.Options) ui.Element) {
	chartDemo(s)
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
		if name == "Donut chart" {
			components.Text(c, tr(c, "chartDistribution")).Bold()
			labels := []string{"chartDesign", "chartDevelopment", "chartResearch", "chartPlanning"}
			var slices []chart.Slice
			if !s.chart.empty {
				for i, label := range labels {
					slices = append(slices, chart.Slice{ID: label, Label: tr(c, label), Value: values[i]})
				}
			}
			chart.Donut(c, &s.chart.selection, slices, opts)
			return
		}
		components.Text(c, tr(c, "chartDaily")).Bold()
		var data chart.Data
		if !s.chart.empty {
			for _, day := range []string{"chartMon", "chartTue", "chartWed", "chartThu", "chartFri", "chartSat"} {
				data.Labels = append(data.Labels, tr(c, day))
			}
			data.Series = []chart.Series{
				{ID: "current", Label: tr(c, "chartCurrent"), Values: values},
				{ID: "previous", Label: tr(c, "chartPrevious"), Values: []float64{26, 34, 28, 48, 43, 57}},
			}
		}
		render(c, &s.chart.selection, data, opts)
	})
}
