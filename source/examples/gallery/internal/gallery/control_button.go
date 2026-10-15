package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

type buttonExampleState struct {
	clicks [6]int
}

func controlButton(c *ui.Context, s *controlState, variant string) {
	if variant == "ToggleButton" {
		controlToggleButton(c, s)
		return
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		for i, variant := range []struct {
			name  string
			value components.ButtonVariant
		}{
			{"Primary", components.ButtonPrimary}, {"Secondary", components.ButtonSecondary},
			{"Outline", components.ButtonOutline}, {"Ghost", components.ButtonGhost}, {"Link", components.ButtonLink},
			{"Toolbar", components.ButtonToolbar},
		} {
			// Key the parent so changing the count never changes button identity.
			ui.Box(c).Key(variant.name).Children(func() {
				button := components.Button(c, "", variant.value).Label(tr(c, variant.name))
				if button.Clicked() {
					s.button.clicks[i]++
				}
				button.Children(func() {
					ui.Text(c, msg(c).LabelCount(tr(c, variant.name), s.button.clicks[i])).SingleLine()
				})
			})
		}
	})
}
