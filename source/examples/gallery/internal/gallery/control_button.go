package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// controlButtonMenu composes the ordinary Onevra button with MyGo's menu behavior.
func controlButtonMenu(c *ui.Context, s *controlState) {
	components.Button(c, "", components.ButtonSecondary).Label("Choose size").AlignSelf(ui.Start).Gap(c.Theme().Space(2)).
		Children(func() {
			components.Text(c, "Size: "+s.choice).SingleLine()
			ui.Icon(c, lucide.ChevronDown()).Size(16, 16)
		}).Menu(func(m *ui.Menu) {
		for _, size := range []string{"Small", "Medium", "Large"} {
			if m.Item(size).Checked(s.choice == size).Chosen() {
				s.choice = size
			}
		}
	})
}

type buttonExampleState struct {
	clicks [5]int
}

func controlButton(c *ui.Context, s *controlState, variant string) {
	if variant == "ToggleButton" {
		controlToggleButton(c, s)
		return
	}
	if variant == "MenuButton" {
		controlButtonMenu(c, s)
		return
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		for i, variant := range []struct {
			name  string
			value components.ButtonVariant
		}{
			{"Primary", components.ButtonPrimary}, {"Secondary", components.ButtonSecondary},
			{"Outline", components.ButtonOutline}, {"Ghost", components.ButtonGhost}, {"Link", components.ButtonLink},
		} {
			// Key the parent so changing the count never changes button identity.
			ui.Box(c).Key(variant.name).Children(func() {
				button := components.Button(c, "", variant.value).Label(variant.name)
				if button.Clicked() {
					s.button.clicks[i]++
				}
				button.Children(func() {
					components.Textf(c, "%s: %d", variant.name, s.button.clicks[i]).SingleLine()
				})
			})
		}
	})
}
