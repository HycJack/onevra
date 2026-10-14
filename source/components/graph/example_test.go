package graph_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/graph"
	"github.com/oneo-me/onevra/source/theme"
)

func ExampleView() {
	var state graph.State
	var opened string
	data := graph.Data{
		Nodes: []graph.Node{
			{ID: "index", Label: "Index", Openable: true},
			{ID: "ideas", Label: "Ideas", Group: "notes", Openable: true},
			{ID: "sketch", Label: "sketch.png", Group: "notes", Size: 8, Color: ui.Hex("#e0b35a")},
		},
		Edges: []graph.Edge{
			{Source: "index", Target: "ideas", Style: graph.EdgeStyle{Distance: 160}},
			{Source: "ideas", Target: "sketch", Style: graph.EdgeStyle{Distance: 64, Opacity: .13}},
		},
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Auto)
		graph.View(c, &state, data, graph.Options{
			Label: "Notes graph", Height: 360,
			Open: func(id string) { opened = id },
		})
	}
	// Supply view to a Onevra Shell, or run it without a window.
	tt := ui.NewTester(view, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	fmt.Println(state.Settled(), opened == "")
	// Output: true true
}
