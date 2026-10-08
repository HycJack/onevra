package gallery

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func exampleGlass(c *ui.Context, s *demoState) {
	components.Slider(c, &s.value, 0, 100).Label("Background position").FillWidth()
	for _, sample := range []struct {
		name     string
		material glass.Glass
	}{
		{"Regular", glass.Glass{}},
		{"Clear", glass.Glass{Style: glass.Clear}},
		{"Tinted", glass.Glass{Tint: c.Theme().Accent}},
	} {
		components.Text(c, sample.name).Bold()
		ui.Box(c).FillWidth().Height(170).Radius(16).Clip().Children(func() {
			// Paint the backdrop first: glass samples content in the same window.
			ui.Image(c, demoImage).Absolute().Top(0).Bottom(0).
				Left(-float32(s.value)).Right(float32(s.value) - 100).Fit(ui.Cover)
			ui.Column(c).Absolute().Top(30).Left(24).Right(24).Height(110).
				Radius(28).Material(sample.material).Center().Gap(8).Children(func() {
				color := c.Theme().Text
				if sample.name == "Tinted" {
					color = c.Theme().AccentText
				} else if sample.name == "Clear" {
					color = ui.RGB(15, 23, 42)
				}
				components.Text(c, "Liquid Glass").Bold().TextColor(color)
				components.Text(c, "Onevra × MyGo").TextColor(color)
			})
		})
	}
}
