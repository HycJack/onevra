package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/shaders"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// shaderDemoState holds a background page's tuning, kept apart per page.
type shaderDemoState struct {
	intensity, scale float64
	still, inert     bool
	prompt           string
}

func controlSignalField(c *ui.Context, s *controlState, _ string) {
	shaderWelcome(c, s, shaders.SignalField)
}

// shaderWelcome shows a background behind the welcome page of an Agent tool,
// with a panel tuning it.
func shaderWelcome(c *ui.Context, s *controlState, background func(*ui.Context, shaders.Options) ui.Element) {
	if s.shader == nil {
		s.shader = &shaderDemoState{intensity: 1, scale: 1}
	}
	d := s.shader
	t := c.Theme()
	background(c, shaders.Options{Intensity: float32(d.intensity), Scale: float32(d.scale), Still: d.still, Inert: d.inert}).Column().Children(func() {
		ui.Column(c).FillWidth().Grow(1).MinWidth(0).MinHeight(0).AlignItems(ui.Center).Justify(ui.Center).Padding(t.Space(6)).Gap(t.Space(5)).Children(func() {
			ui.Column(c).AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
				components.Icon(c, lucide.Sparkles).Size(32, 32).TextColor(t.Accent)
				components.Text(c, tr(c, "shaderWelcomeTitle")).FontSize(t.Rem(1.6)).Bold().TextAlign(ui.Center).Wrap()
				components.Text(c, tr(c, "shaderWelcomeSubtitle")).TextColor(t.TextMuted).TextAlign(ui.Center).Wrap()
			})
			ui.Row(c).FillWidth().MaxWidth(560).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				components.TextBox(c, &d.prompt, components.TextBoxOptions{LeadingIcon: lucide.MessageSquare()}).
					Placeholder(tr(c, "shaderWelcomePrompt")).Label(tr(c, "shaderWelcomePrompt")).Grow(1).MinWidth(0)
				components.Button(c, tr(c, "start"), components.ButtonPrimary)
			})
			ui.Row(c).Wrap().Justify(ui.Center).Gap(t.Space(2)).Children(func() {
				for _, key := range []string{"shaderWelcomeReview", "shaderWelcomeTests", "shaderWelcomeExplain"} {
					if components.Button(c, tr(c, key), components.ButtonOutline).Key(key).Clicked() {
						d.prompt = tr(c, key)
					}
				}
			})
		})
		ui.Row(c).FillWidth().Justify(ui.Center).Padding(0, t.Space(4), t.Space(4)).Children(func() {
			ui.Row(c).Wrap().MinWidth(0).AlignItems(ui.Center).Justify(ui.Center).GapX(t.Space(4)).GapY(t.Space(2)).
				Padding(t.Space(2), t.Space(3)).Background(t.Surface.Alpha(.88)).Border(1, t.Border).Radius(12).Children(func() {
				for _, slider := range []struct {
					key    string
					value  *float64
					lo, hi float64
				}{{"shaderIntensity", &d.intensity, .2, 2}, {"shaderScale", &d.scale, .5, 2}} {
					ui.Row(c).Key(slider.key).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
						components.Text(c, tr(c, slider.key)).TextColor(t.TextMuted)
						components.Slider(c, slider.value, slider.lo, slider.hi).Label(tr(c, slider.key)).Width(120)
					})
				}
				components.Switch(c, &d.still, tr(c, "shaderStill"))
				components.Switch(c, &d.inert, tr(c, "shaderInert"))
			})
		})
	})
}
