package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func exampleMotion(c *ui.Context, s *demoState) {
	t := c.Theme()
	switch s.variant {
	case "Keyframes":
		if components.Button(c, tr(c, "gallery.motion_replay"), components.ButtonOutline).Clicked() {
			s.count++
		}
		ui.Box(c).FillWidth().Height(80).Padding(12).Children(func() {
			tile := ui.Box(c).Key("sequence").Size(44, 44).Radius(10).Background(t.Accent).Label(tr(c, "gallery.motion_tile"))
			x := tile.Keyframes("x", []ui.Keyframe{
				{Value: 0}, {At: 200 * time.Millisecond, Value: 110, Ease: ui.EaseOut},
				{At: 500 * time.Millisecond, Value: 55}, {At: 800 * time.Millisecond, Value: 0, Ease: ui.EaseInOut},
			}, ui.KeyframeOptions{Run: uint64(s.count), Repeat: 1, Alternate: true})
			tile.Rotate(x/3).Translate(x, 0)
		})
	case "Animated pose":
		components.Checkbox(c, &s.on, tr(c, "gallery.motion_lift"))
		ui.Box(c).FillWidth().Height(100).Padding(20).Children(func() {
			card := components.Button(c, tr(c, "gallery.motion_hover"), components.ButtonOutline)
			pose := ui.Pose{}
			if s.on || card.Hovered() {
				pose = ui.Pose{Y: -8, Scale: 1.04, Rotate: -4}
			}
			card.AnimatePose(pose, ui.Snappy)
		})
	default:
		if components.Button(c, tr(c, "gallery.motion_reverse"), components.ButtonOutline).Clicked() {
			s.on = !s.on
		}
		for i, animation := range []ui.Animation{ui.Smooth, ui.Snappy, ui.Bouncy} {
			ui.Column(c).Key(i).FillWidth().Gap(4).Children(func() {
				components.Text(c, tr(c, []string{"gallery.motion_smooth", "gallery.motion_snappy", "gallery.motion_bouncy"}[i])).TextColor(t.TextMuted)
				ui.Box(c).FillWidth().Height(56).Padding(6).Children(func() {
					tile := ui.Box(c).Size(36, 36).Radius(9).Background(t.Accent)
					target := float32(0)
					if s.on {
						target = 120
					}
					tile.Translate(tile.AnimateTo("x", target, animation), 0)
				})
			})
		}
	}
}

func exampleTransforms(c *ui.Context, s *demoState) {
	t := c.Theme()
	if s.variant == "Perspective" {
		if components.Button(c, tr(c, "gallery.transform_flip"), components.ButtonOutline).Clicked() {
			s.on = !s.on
		}
		ui.Box(c).FillWidth().Height(140).Padding(20).Children(func() {
			card := ui.Box(c).Size(150, 96)
			target := float32(0)
			if s.on {
				target = 180
			}
			angle := card.Spring("turn", target, ui.SpringOptions{Duration: 600 * time.Millisecond, Bounce: .1})
			card.Children(func() {
				for i, label := range []string{"gallery.transform_front", "gallery.transform_back"} {
					ui.Box(c).Absolute().Fill().Radius(12).Background(t.Surface).Border(1, t.Accent).
						RotateY(angle + float32(i)*180).Perspective(800).BackfaceHidden().
						AlignItems(ui.Center).Justify(ui.Center).Children(func() { components.Text(c, tr(c, label)) })
				}
			})
		})
		return
	}
	components.Slider(c, &s.value, -15, 15).Label(tr(c, "gallery.transform_angle")).FillWidth()
	components.Checkbox(c, &s.on, tr(c, "gallery.transform_scale"))
	ui.Box(c).FillWidth().Height(230).Padding(24, 12).Children(func() {
		scale := float32(1)
		if s.on {
			scale = .85
		}
		pose := ui.Scaling(scale, scale).Then(ui.Rotation(float32(s.value))).Then(ui.Translation(8, 8))
		ui.Column(c).Width(160).Padding(12).Gap(10).Radius(12).Background(t.Surface).Border(1, t.Border).
			Transform(pose).TransformOrigin(.5, .5).Children(func() {
			components.TextBox(c, &s.text, components.TextBoxOptions{}).Label(tr(c, "gallery.transform_note")).FillWidth()
			anchor := components.Button(c, tr(c, "gallery.show_popover"), components.ButtonOutline).FillWidth()
			if anchor.Clicked() {
				s.open = !s.open
			}
			components.Popover(c, anchor, &s.open, func(panel ui.Element) {
				panel.Label(tr(c, "gallery.popover_content"))
				components.Text(c, tr(c, "gallery.popover_content"))
				if components.Button(c, tr(c, "gallery.close"), components.ButtonSecondary).Clicked() {
					s.open = false
				}
			})
		})
	})
}
