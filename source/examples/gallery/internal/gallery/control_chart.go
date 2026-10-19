package gallery

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// chartTile is one chart's sample state; each tile has its own tools.
type chartTile struct {
	selection       chart.State
	round           int // data updates
	count           int // bubble categories shown
	empty, noLegend bool
}

type chartDemoState struct {
	tiles map[string]*chartTile
}

// tile returns a chart's sample state, starting with no selection.
func (d *chartDemoState) tile(name string) *chartTile {
	if d.tiles == nil {
		d.tiles = make(map[string]*chartTile)
	}
	if d.tiles[name] == nil {
		d.tiles[name] = &chartTile{selection: chart.State{Index: -1}, count: bubbleInitialCount}
	}
	return d.tiles[name]
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

// chartTools holds a tile's sample actions beside its title.
func chartTools(c *ui.Context, name string, d *chartTile) {
	tool := func(icon *ui.SVG, label string, disabled bool) bool {
		return components.IconButton(c, icon, tr(c, label), components.IconButtonCompact).Disabled(disabled).Clicked()
	}
	if tool(lucide.RefreshCw(), "chartUpdate", d.empty) {
		d.round++
	}
	if name == "Bubble chart" {
		if tool(lucide.Plus(), "bubbleAdd", d.empty || d.count == len(bubbleLabels)) {
			d.count++
		}
		if tool(lucide.Minus(), "bubbleRemove", d.empty || d.count == 1) {
			d.count--
		}
	}
	icon, label := lucide.EyeOff(), "chartHideLegend"
	if d.noLegend {
		icon, label = lucide.Eye(), "chartShowLegend"
	}
	if tool(icon, label, false) {
		d.noLegend = !d.noLegend
	}
	icon, label = lucide.Eraser(), "chartEmpty"
	if d.empty {
		icon, label = lucide.Undo2(), "chartRestore"
	}
	if tool(icon, label, false) {
		d.empty = !d.empty
	}
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

// chartTileWidth and chartTileHeight size the tile each chart takes.
const chartTileWidth, chartTileHeight = 480, 360

// controlChart shows the chart a variant names in a fixed tile, titled
// with the variant and its tools; selection details appear in the chart's
// tooltip.
func controlChart(c *ui.Context, s *controlState, name string) {
	d := chartDemo(s).tile(name)
	t := c.Theme()
	ui.Column(c).Width(chartTileWidth).Height(chartTileHeight).MaxWidthPercent(100).MinWidth(0).Padding(t.Space(4)).Gap(t.Space(4)).
		Radius(t.Radius*1.5).Background(t.Surface).Border(1, t.Border.Alpha(.5)).Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(t.Space(1)).Children(func() {
			components.Text(c, tr(c, name)).Bold().Grow(1).MinWidth(0).NoWrap().Ellipsis("…")
			chartTools(c, name, d)
		})
		opts := chart.Options{
			Label: tr(c, "chartSelection"), Description: tr(c, "chartKeys"),
			EmptyText: tr(c, "chartNoData"), PlotHeight: 220, HideLegend: d.noLegend,
			FormatValue: func(v float64) string { return strconv.FormatFloat(v, 'f', 0, 64) },
		}
		// Updates switch between two sets of sample data.
		updated := d.round%2 == 1
		values := []float64{38, 52, 46, 68, 61, 82}
		if updated {
			values = []float64{62, 43, 71, 54, 78, 65}
		}
		selection := &d.selection
		if name == "Split bar" {
			labels := []string{"chartOrganic", "chartDirect", "chartPaid", "chartSocial", "chartReferral", "chartOther"}
			sessions := []float64{41260, 24610, 15830, 8420, 5310, 3790}
			if updated {
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
				HideLegend: d.noLegend,
			})
			return
		}
		if name == "Bubble chart" {
			var slices []chart.Slice
			if !d.empty {
				values := bubbleValues(d.round)
				for i, label := range bubbleLabels[:d.count] {
					slices = append(slices, chart.Slice{ID: label, Label: tr(c, label), Value: values[i]})
				}
			}
			chart.Bubble(c, selection, slices, opts)
			return
		}
		if name == "Donut chart" {
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
