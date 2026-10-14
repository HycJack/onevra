package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

var sample = Data{Labels: []string{"Mon", "Tue", "Wed", "Thu"}, Series: []Series{
	{ID: "a", Label: "Alpha", Values: []float64{10, 40, -20, 30}},
	{ID: "b", Label: "Beta", Values: []float64{20, 10, -10, 40}},
}}

var cartesianViews = map[string]func(*ui.Context, *State, Data, Options) ui.Element{
	"Bar": Bar, "Line": Line, "Area": Area, "StackedBar": StackedBar,
}

func TestCartesianInteraction(t *testing.T) {
	for name, view := range cartesianViews {
		for _, scale := range []float32{.75, 1, 1.5} {
			t.Run(fmt.Sprintf("%s/%g", name, scale), func(t *testing.T) {
				state := State{Index: -1}
				var selected Selection
				opts := Options{Label: "Plot", Description: "Arrow keys", EmptyText: "Empty", Detail: func(_ *ui.Context, s Selection) { selected = s }}
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Column(c).Width(280).Scale(scale, scale).TransformOrigin(0, 0).Children(func() { view(c, &state, sample, opts) })
				}, 500, 600)
				for _, dark := range []bool{false, true} {
					tt.SetDark(dark)
					r, ok := tt.Find("Plot")
					if !ok || r.W <= 0 || r.H <= 0 {
						t.Fatalf("invalid plot: %+v", r)
					}
					tt.Move(r.X+r.W*.92, r.Y+r.H/2)
					if state.Index != 3 || selected.Label != "Thu" || len(selected.Values) != 2 {
						t.Fatalf("pointer selection: %+v", selected)
					}
					tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
					tt.Key(0, ui.KeyHome)
					tt.Key(0, ui.KeyRight)
					if state.Index != 1 || selected.Values[0].Value != 40 {
						t.Fatalf("keyboard selection: %+v", selected)
					}
					tt.Key(0, ui.KeyEnd)
					if state.Index != 3 {
						t.Fatal("End did not select last category")
					}
					tt.Key(0, ui.KeyEscape)
					if state.Index != -1 {
						t.Fatal("Escape did not clear selection")
					}
				}
			})
		}
	}
}

func TestLegendEmptyAndDataChanges(t *testing.T) {
	for name, view := range cartesianViews {
		t.Run(name, func(t *testing.T) {
			data := sample
			state := State{Index: 3}
			var selected Selection
			tt := ui.NewTester(func(c *ui.Context) {
				view(c, &state, data, Options{Label: "Plot", EmptyText: "Empty", Detail: func(_ *ui.Context, s Selection) { selected = s }})
			}, 280, 500)
			if err := tt.Click("Beta"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if !state.Hidden["b"] || len(selected.Values) != 1 || selected.Values[0].ID != "a" {
				t.Fatalf("legend did not hide series: %+v", selected)
			}
			data.Series = []Series{sample.Series[1], sample.Series[0]}
			tt.Frame()
			if len(selected.Values) != 1 || selected.Values[0].ID != "a" {
				t.Fatal("hidden state followed series position instead of identity")
			}
			if err := tt.Click("Alpha"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if !tt.HasText("Empty") || state.Index != -1 {
				t.Fatal("all-hidden state is not empty")
			}
			if err := tt.Click("Beta"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if tt.HasText("Empty") {
				t.Fatal("legend cannot restore an empty chart")
			}
			state.Index, data.Labels = 3, data.Labels[:1]
			tt.Frame()
			if state.Index != -1 {
				t.Fatal("data shrink left a stale selection")
			}
		})
	}
}

func TestDisabledCharts(t *testing.T) {
	views := map[string]func(*ui.Context, *State){
		"Donut": func(c *ui.Context, state *State) {
			Donut(c, state, []Slice{{ID: "a", Label: "Alpha", Value: 2}, {ID: "b", Label: "Beta", Value: 3}}, Options{Label: "Plot"})
		},
	}
	for name, view := range cartesianViews {
		views[name] = func(c *ui.Context, state *State) { view(c, state, sample, Options{Label: "Plot"}) }
	}
	for name, view := range views {
		t.Run(name, func(t *testing.T) {
			state := State{Index: -1}
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Disabled(true).Children(func() { view(c, &state) })
			}, 280, 500)
			if err := tt.Click("Plot"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEnd)
			if err := tt.Click("Alpha"); err != nil {
				t.Fatal(err)
			}
			if state.Index != -1 || len(state.Hidden) != 0 {
				t.Fatalf("disabled chart changed state: %+v", state)
			}
		})
	}
}

func TestDomainAndMissingValues(t *testing.T) {
	data := Data{Labels: []string{"A", "B", "C", "D"}, Series: []Series{
		{ID: "a", Values: []float64{10, -20, math.NaN(), math.Inf(1), 999}},
		{ID: "b", Values: []float64{30, -40}},
	}}
	p := prepare(data, &State{}, stacks)
	if !p.valid || p.y.Lo > -60 || p.y.Hi < 40 || p.stack[1][0].Hi != 40 || p.stack[1][1].Lo != -60 {
		t.Fatalf("incorrect stacked domain: %+v", p)
	}
	if !math.IsNaN(p.values[0][2]) || !math.IsNaN(p.values[0][3]) || !math.IsNaN(p.values[1][3]) {
		t.Fatal("invalid and missing values must remain gaps")
	}
	if p := prepare(data, &State{Hidden: map[string]bool{"b": true}}, bars); p.y.Hi != 10 || p.y.Lo != -20 {
		t.Fatalf("hidden or excess values changed the domain: %+v", p.y)
	}
	for _, values := range [][]float64{nil, {math.NaN()}, {math.Inf(1)}, {-math.MaxFloat64, math.MaxFloat64}} {
		bad := Data{Labels: []string{"A", "B"}, Series: []Series{{ID: "a", Values: values}}}
		if prepare(bad, &State{}, bars).valid {
			t.Fatalf("invalid data has a scale: %v", values)
		}
	}
	zero := Data{Labels: []string{"A"}, Series: []Series{{ID: "a", Values: []float64{0}}}}
	if p := prepare(zero, &State{}, areas); !p.valid || p.y.Lo != 0 || p.y.Hi != 1 {
		t.Fatal("zero data should use a visible nondegenerate scale")
	}
}

func TestDonutSelectionAndFiltering(t *testing.T) {
	data := []Slice{{ID: "a", Label: "Alpha", Value: 3}, {ID: "zero", Label: "Zero"}, {ID: "b", Label: "Beta", Value: 1}, {ID: "bad", Label: "Invalid", Value: math.Inf(1)}}
	state := State{Index: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Donut(c, &state, data, Options{Label: "Plot", EmptyText: "Empty", FormatValue: func(v float64) string { return fmt.Sprintf("%.0f h", v) }})
	}, 280, 500)
	if tt.HasText("Zero") || tt.HasText("Invalid") || !tt.HasText("4 h") {
		t.Fatal("donut did not filter invalid values or format the total")
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if state.Index != 2 || !tt.HasText("Beta  1 h") {
		t.Fatal("legend selection lost input index or formatting")
	}
	if err := tt.Click("Plot"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	if state.Index != 2 {
		t.Fatal("keyboard selection did not skip zero values")
	}
	tt.Key(0, ui.KeyEscape)
	if state.Index != -1 {
		t.Fatal("Escape did not clear donut selection")
	}
	r, _ := tt.Find("Plot")
	tt.Move(r.X+r.W/2+80, r.Y+r.H/2)
	if state.Index != 0 {
		t.Fatal("ring hover did not select slice")
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	if state.Index != -1 {
		t.Fatal("the ring hole should not select a slice")
	}
	data = []Slice{{ID: "a", Value: math.MaxFloat64}, {ID: "b", Value: math.MaxFloat64}}
	tt.Frame()
	tt.SetDark(true)
	if tt.Image() == nil || tt.HasText("Empty") {
		t.Fatal("large positive slices should still render")
	}
	data = nil
	tt.Frame()
	if state.Index != -1 || !tt.HasText("Empty") {
		t.Fatal("empty donut kept a selection")
	}
}

func TestLocalStateAndSlots(t *testing.T) {
	selected := [2]Selection{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			for i := range 2 {
				ui.Box(c).Key(i).FillWidth().Children(func() {
					Line(c, nil, sample, Options{Label: fmt.Sprint(i), PlotHeight: 100, HideLegend: true,
						Detail: func(_ *ui.Context, s Selection) { selected[i] = s }})
				})
			}
			Bar(c, nil, Data{}, Options{PlotHeight: 60, Empty: func(c *ui.Context) { ui.Text(c, "Custom empty") }})
		})
	}, 280, 450)
	if err := tt.Click("0"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	if selected[0].Index != 0 || selected[1].Index != -1 || !tt.HasText("Custom empty") {
		t.Fatalf("element-local state leaked or slots failed: %+v", selected)
	}
}
