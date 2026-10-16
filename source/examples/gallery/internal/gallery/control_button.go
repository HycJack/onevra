package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

type buttonExampleState struct {
	clicks [5]int
	// actions counts clicks on the icon, mixed and toolbar examples by their
	// message key.
	actions map[string]int
}

// buttonAction is a sample action shared by the icon examples.
type buttonAction struct {
	key  string
	icon func() *ui.SVG
}

func controlButton(c *ui.Context, s *controlState, variant string) {
	switch variant {
	case "ToggleButton":
		controlToggleButton(c, s)
	case "Icon and text":
		controlIconTextButton(c, s)
	case "IconButton":
		controlIconButton(c, s)
	case "Toolbar":
		controlToolbarButton(c, s)
	default:
		controlTextButton(c, s)
	}
}

func controlTextButton(c *ui.Context, s *controlState) {
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

// iconTextButton places a leading icon before a button's text.
func iconTextButton(c *ui.Context, s *controlState, variant components.ButtonVariant, action buttonAction) {
	t := c.Theme()
	label := tr(c, action.key)
	ui.Box(c).Key(action.key).Children(func() {
		button := components.Button(c, "", variant).Label(label).Gap(t.Space(1.5))
		if button.Clicked() {
			s.button.count(action.key)
		}
		button.Children(func() {
			components.Icon(c, action.icon()).Size(t.Rem(1), t.Rem(1))
			components.Text(c, label).SingleLine()
		})
	})
}

func controlIconTextButton(c *ui.Context, s *controlState) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		iconTextButton(c, s, components.ButtonPrimary, buttonAction{"new", lucide.Plus})
		iconTextButton(c, s, components.ButtonSecondary, buttonAction{"download", lucide.Download})
		iconTextButton(c, s, components.ButtonOutline, buttonAction{"share", lucide.Share})
		iconTextButton(c, s, components.ButtonGhost, buttonAction{"edit", lucide.Pencil})
	})
}

func controlIconButton(c *ui.Context, s *controlState) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(2)).Children(func() {
		for _, action := range []buttonAction{
			{"copy", lucide.Copy}, {"favorite", lucide.Star}, {"refresh", lucide.RefreshCw},
		} {
			if components.IconButton(c, action.icon(), tr(c, action.key), components.IconButtonPlain).Clicked() {
				s.button.count(action.key)
			}
		}
		components.IconButton(c, lucide.Trash(), tr(c, "delete"), components.IconButtonPlain).Disabled(true)
	})
}

// controlToolbarButton keeps toolbar faces apart from the standard variants:
// the raised Onevra faces first, then the Shell's own toolbar strip.
func controlToolbarButton(c *ui.Context, s *controlState) {
	t := c.Theme()
	ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(4)).Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(2)).Children(func() {
			ui.Box(c).Key("filter").Children(func() {
				if components.Button(c, tr(c, "filter"), components.ButtonToolbar).Clicked() {
					s.button.count("filter")
				}
			})
			iconTextButton(c, s, components.ButtonToolbar, buttonAction{"upload", lucide.Upload})
			for _, action := range []buttonAction{{"undo", lucide.Undo}, {"redo", lucide.Redo}} {
				if components.IconButton(c, action.icon(), tr(c, action.key), components.IconButtonToolbar).Clicked() {
					s.button.count(action.key)
				}
			}
		})
		ui.Column(c).Width(240).MaxWidthPercent(100).Gap(12).Children(func() {
			components.Text(c, tr(c, "shellToolbar")).TextColor(t.TextMuted)
			// The same strip Shell pins beneath the sidebar for LeftToolbar.
			ui.Toolbar(c, func() {
				for _, action := range []buttonAction{{"newFolder", lucide.FolderPlus}, {"search", lucide.Search}} {
					if sidebar.IconButton(c, action.icon(), tr(c, action.key)).Clicked() {
						s.button.count(action.key)
					}
				}
				ui.Spacer(c)
				if sidebar.IconButton(c, lucide.Settings(), tr(c, "settings")).Clicked() {
					s.button.count("settings")
				}
			}).FillWidth().Padding(8).Radius(t.Radius).Border(1, t.Border)
		})
	})
}

func (b *buttonExampleState) count(key string) {
	if b.actions == nil {
		b.actions = map[string]int{}
	}
	b.actions[key]++
}
