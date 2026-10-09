package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
)

// Welcome samples are local compositions, independent of the component pages.
type welcomeState struct {
	initialized bool
	name, kind  string
	include     bool
	created     bool
	mode        int
	intensity   [2]float64
	quiet       bool
}

func (s *welcomeState) initialize() {
	if s.initialized {
		return
	}
	*s = welcomeState{initialized: true, name: "A new beginning", kind: "Personal", include: true, quiet: true, intensity: [2]float64{72, 36}}
}

func (g *Gallery) Welcome(c *ui.Context) {
	g.welcome.initialize()
	t := c.Theme()
	components.ScrollView(c, components.ScrollViewOptions{}).Key("welcome").Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).MaxWidth(t.Space(280)).AlignSelf(ui.Center).
			Padding(t.Space(6)).Gap(t.Space(7)).Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(4)).Children(func() {
				ui.Column(c).Gap(t.Space(1.5)).Children(func() {
					components.Text(c, "Small components. Real possibilities.").FontSize(t.Rem(1.35)).Bold().Wrap()
					components.Text(c, "A few familiar pieces, brought together. Go ahead, make them yours.").TextColor(t.TextMuted).Wrap()
				})
				ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Stretch).Gap(t.Space(4)).Children(func() {
					welcomeCard(c, "idea").Children(func() { g.welcomeIdea(c) })
					welcomeCard(c, "balance").Children(func() { g.welcomeBalance(c) })
				})
			})
			g.welcomeIcons(c)
			components.Text(c, "Built with Onevra · Powered by MyGo · Written in Go").
				FontSize(t.Rem(.8)).TextColor(t.TextMuted).Wrap()
		})
	})
}

// Paint behind the titlebar inset as well as the scrolling page, not inside a card.
func welcomeBackground(c *ui.Context, content ui.Element) {
	t := c.Theme()
	content.Gradient(t.Background.Mix(t.Accent, .10), t.Background, 155).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.Clip(r, 0, func() {
				for _, radius := range []float32{24, 40, 56, 72, 88} {
					var orbit ui.Path
					orbit.Circle(r.X+r.W-t.Space(4), r.Y+t.Space(14), t.Space(radius))
					p.StrokePath(&orbit, 1, t.Accent.Alpha(.09))
				}
			})
		})
}

func welcomeCard(c *ui.Context, key string) ui.Element {
	t := c.Theme()
	return ui.Column(c).Key(key).Basis(t.Space(60)).Grow(1).MinWidth(0).MaxWidthPercent(100).
		Padding(t.Space(5)).Gap(t.Space(4)).Radius(t.Space(4)).
		Background(t.Surface.Alpha(.94)).Border(1, t.Border).
		Shadow(0, t.Space(1), t.Space(4), -t.Space(2), t.Text.Alpha(.08))
}

func welcomeCardHeading(c *ui.Context, icon func() *ui.SVG, title, subtitle string, color ui.Color) {
	t := c.Theme()
	ui.Row(c).Wrap().Gap(t.Space(3)).Children(func() {
		ui.Box(c).Size(t.Space(10), t.Space(10)).Radius(t.Space(3)).Center().
			Background(color.Alpha(.10)).TextColor(color).Children(func() {
			components.Icon(c, icon).Size(t.Space(5), t.Space(5))
		})
		ui.Column(c).MinWidth(0).Gap(t.Space(1)).Children(func() {
			components.Text(c, title).Bold().Wrap()
			components.Text(c, subtitle).FontSize(t.Rem(.8)).TextColor(t.TextMuted).Wrap()
		})
	})
}

func (g *Gallery) welcomeIdea(c *ui.Context) {
	s, t := &g.welcome, c.Theme()
	welcomeCardHeading(c, lucide.Sparkles, "Your next idea", "A fresh start, in a few small steps.", t.Accent)
	ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(2)).Children(func() {
		components.Text(c, "Give it a name").FontSize(t.Rem(.85)).TextColor(t.TextMuted)
		if components.TextBox(c, &s.name, components.TextBoxOptions{}).Label("Preview name").Placeholder("Something wonderful…").FillWidth().MinWidth(0).Changed() {
			s.created = false
		}
	})
	ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(2)).Children(func() {
		components.Text(c, "Make space for").FontSize(t.Rem(.85)).TextColor(t.TextMuted)
		sel := components.Select(c, &s.kind)
		sel.Trigger.Label("Preview category").FillWidth().MinWidth(0).Children(func() {
			components.Text(c, s.kind).SingleLine()
		})
		sel.Popup(func(panel ui.Element) {
			for _, kind := range []string{"Personal", "Work", "Play"} {
				sel.Item(kind).Label(kind).Children(func() { components.Text(c, kind) })
			}
		})
		if sel.Trigger.Changed() {
			s.created = false
		}
	})
	if components.Checkbox(c, &s.include, "Include a little inspiration").Changed() {
		s.created = false
	}
	ui.Spacer(c)
	if welcomeAction(c, "Create preview", components.ButtonPrimary, lucide.ArrowUpRight).FillWidth().Clicked() {
		s.created = true
	}
	if s.created {
		ui.Textf(c, "Ready to explore: %s", s.name).FontSize(t.Rem(.8)).TextColor(t.Success).Wrap()
	} else {
		components.Text(c, "Just a demo. Your ideas stay here.").FontSize(t.Rem(.8)).TextColor(t.TextMuted).Wrap()
	}
}

func (g *Gallery) welcomeBalance(c *ui.Context) {
	s, t := &g.welcome, c.Theme()
	welcomeCardHeading(c, lucide.SlidersHorizontal, "Find your balance", "Less noise. A little more flow.", t.Success)
	components.Tabs(c, &s.mode, "Focus", "Relax")
	value := &s.intensity[s.mode]
	ui.Row(c).FillWidth().Wrap().Gap(t.Space(4)).Padding(t.Space(3)).Radius(t.Space(3)).
		Gradient(t.Background.Mix(t.Success, .07), t.Background, 110).Children(func() {
		components.ProgressRing(c, *value/100).Size(t.Space(16), t.Space(16)).Label("Focus intensity")
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			ui.Textf(c, "%.0f%%", *value).FontSize(t.Rem(2)).Bold()
			caption := "In the zone"
			if s.mode == 1 {
				caption = "Room to breathe"
			}
			components.Text(c, caption).FontSize(t.Rem(.8)).TextColor(t.TextMuted)
		})
	})
	ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
		components.Text(c, "Intensity").FontSize(t.Rem(.85)).TextColor(t.TextMuted)
		components.Slider(c, value, 0, 100).Label("Preview intensity").FillWidth().MinWidth(0)
	})
	ui.Spacer(c)
	components.Switch(c, &s.quiet, "Quiet mode")
	caption := "A quieter corner, just for you."
	if !s.quiet {
		caption = "Back in the mix."
	}
	components.Text(c, caption).FontSize(t.Rem(.8)).TextColor(t.TextMuted).Wrap()
}

func (g *Gallery) welcomeIcons(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(5)).Gap(t.Space(4)).Radius(t.Space(4)).
		Background(t.Background.Mix(t.Accent, .035)).Border(1, t.Border).Children(func() {
		ui.Row(c).Wrap().Gap(t.Space(3)).TextColor(t.TextMuted).Children(func() {
			for _, icon := range []func() *ui.SVG{lucide.Command, lucide.Sun, lucide.CloudRain, lucide.Headphones, lucide.Heart, lucide.Zap, lucide.Palette, lucide.FingerprintPattern} {
				ui.Box(c).Size(t.Space(9), t.Space(9)).Center().Background(t.Surface).Radius(t.Space(2.5)).Children(func() {
					components.Icon(c, icon).Size(t.Space(5), t.Space(5))
				})
			}
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(4)).Children(func() {
			ui.Column(c).Basis(t.Space(50)).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				components.Text(c, "A small mark. A big personality.").Bold().Wrap()
				components.Text(c, "Meet the Lucide icon collection.").TextColor(t.TextMuted).FontSize(t.Rem(.85)).Wrap()
			})
			if welcomeAction(c, "Explore Lucide", components.ButtonOutline, lucide.ArrowUpRight).Clicked() {
				g.page = lucidePage
				g.enterGallery()
			}
		})
	})
}

func welcomeAction(c *ui.Context, label string, variant components.ButtonVariant, icon func() *ui.SVG) ui.Element {
	b := components.Button(c, "", variant).Label(label)
	b.Children(func() {
		components.Text(c, label).Wrap()
		components.Icon(c, icon).Size(c.Theme().Space(4), c.Theme().Space(4)).Shrink(0)
	})
	return b
}
