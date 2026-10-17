package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chat"
	"github.com/oneo-me/onevra/source/components/shaders"
)

// shaderBackground is one background the page switches between; label is its
// English source text.
type shaderBackground struct {
	id, label string
	build     func(*ui.Context, shaders.Options) ui.Element
}

var shaderBackgrounds = []shaderBackground{
	{"signal", "Signal field", shaders.SignalField},
}

func shaderBackgroundFor(id string) shaderBackground {
	for _, b := range shaderBackgrounds {
		if b.id == id {
			return b
		}
	}
	return shaderBackgrounds[0]
}

// shaderDemoState holds the chosen background and its tuning.
type shaderDemoState struct {
	background       string
	intensity, scale float64
	still, inert     bool
}

// controlShaderBackgrounds shows a background behind an empty conversation,
// with the backgrounds and the chosen one's settings at the top right and the
// conversation composer at the bottom.
func controlShaderBackgrounds(c *ui.Context, s *controlState, _ string) {
	if s.shader == nil {
		s.shader = &shaderDemoState{background: shaderBackgrounds[0].id, intensity: 1, scale: 1}
	}
	if s.chat == nil {
		s.chat = newChatDemoState(c.Now())
	}
	d := s.shader
	t := c.Theme()
	options := shaders.Options{Intensity: float32(d.intensity), Scale: float32(d.scale), Still: d.still, Inert: d.inert}
	shaderBackgroundFor(d.background).build(c, options).Column().Children(func() {
		ui.Row(c).FillWidth().Justify(ui.End).Padding(t.Space(4), t.Space(4), 0).Children(func() {
			shaderPanel(c, d)
		})
		ui.Spacer(c)
		ui.Column(c).FillWidth().MaxWidth(720).AlignSelf(ui.Center).MinWidth(0).Padding(0, t.Space(4), t.Space(4)).Children(func() {
			shaderComposer(c, s.chat)
		})
	})
}

// shaderPanel lists the backgrounds, with the chosen one's settings below.
func shaderPanel(c *ui.Context, d *shaderDemoState) {
	t := c.Theme()
	ui.Column(c).Width(240).Padding(t.Space(2)).Gap(t.Space(3)).
		Background(t.Surface.Alpha(.88)).Border(1, t.Border).Radius(12).Children(func() {
		components.SidebarSection(c, tr(c, "shaderBackground"), func() {
			for _, b := range shaderBackgrounds {
				if components.SidebarItem(c, tr(c, b.label), nil, b.id == d.background).Key(b.id).Clicked() {
					d.background = b.id
				}
			}
		})
		ui.Column(c).FillWidth().Padding(0, t.Space(2), t.Space(1)).Gap(t.Space(3)).Children(func() {
			for _, slider := range []struct {
				key    string
				value  *float64
				lo, hi float64
			}{{"shaderIntensity", &d.intensity, .2, 2}, {"shaderScale", &d.scale, .5, 2}} {
				ui.Column(c).Key(slider.key).FillWidth().Gap(t.Space(1)).Children(func() {
					components.Text(c, tr(c, slider.key)).TextColor(t.TextMuted)
					components.Slider(c, slider.value, slider.lo, slider.hi).Label(tr(c, slider.key)).FillWidth()
				})
			}
			components.Switch(c, &d.still, tr(c, "shaderStill"))
			components.Switch(c, &d.inert, tr(c, "shaderInert"))
		})
	})
}

// shaderComposer is the chat demonstration's composer and its toolbar; Send
// only clears the draft, as the page has no conversation.
func shaderComposer(c *ui.Context, s *chatDemoState) {
	result := chat.GlassComposer(c, nil, &s.draft, chatDemoComposerOptions(c, s))
	s.dropped(result.Dropped)
	if result.Send {
		s.draft, s.files = "", nil
	}
}
