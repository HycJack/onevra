package gallery

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// Positions belong to this demo page and survive navigation, not app restarts.
type glassDemo struct {
	x, y float32
}

func exampleGlass(c *ui.Context, s *demoState) {
	bitmap, err := galleryImage()
	if err != nil {
		components.Text(c, tr(c, "gallery.unable_to_load_avatar_jpeg")).TextColor(c.Theme().TextMuted)
		return
	}
	glassPreview(c, s, bitmap)
}

func glassPreview(c *ui.Context, s *demoState, bitmap *ui.Bitmap) {
	ui.Row(c).Gap(8).Wrap().Children(func() {
		for i, name := range []string{"Regular", "Clear", "Tinted"} {
			style := components.ButtonSecondary
			if s.index == i {
				style = components.ButtonPrimary
			}
			if components.Button(c, glassStyleName(c, name), style).Clicked() {
				s.index = i
			}
		}
		if components.Button(c, tr(c, "gallery.reset_position"), components.ButtonOutline).Clicked() {
			s.glass = glassDemo{}
		}
	})
	components.Text(c, tr(c, "gallery.drag_the_glass_over_the_avatar")).TextColor(c.Theme().TextMuted)
	canvas := ui.Box(c).FillWidth().Height(420).Radius(20).Clip().Label(tr(c, "gallery.glass_canvas"))
	previousWidth := canvas.Bounds().W
	// Bounds is the previous layout; settle initial layout and window resizing.
	canvas.Draw(func(_ *ui.Painter, r ui.Rect) {
		if r.W != previousWidth {
			c.Invalidate()
		}
	})
	canvas.Children(func() {
		// Paint the avatar first: the material samples content beneath it.
		components.Image(c, bitmap).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Fit(ui.Cover).Label(tr(c, "gallery.author_avatar"))
		width := min(float32(220), max(0, canvas.Bounds().W-32))
		maxX := max(0, canvas.Bounds().W-width-32)
		material := glass.Glass{}
		color := c.Theme().Text
		switch s.index {
		case 1:
			material.Style = glass.Clear
			color = ui.RGB(15, 23, 42)
		case 2:
			material.Tint = c.Theme().Accent
			color = c.Theme().AccentText
		}
		pane := ui.Column(c).Absolute().Width(width).Height(120).
			Radius(30).Material(material).Center().Gap(8).Label(tr(c, "gallery.draggable_glass"))
		if dx, dy, ok := pane.Dragged(); ok {
			s.glass.x += dx
			s.glass.y += dy
		}
		// Keep the entire drag target reachable after dragging or resizing.
		s.glass.x = min(maxX, max(0, s.glass.x))
		s.glass.y = min(float32(268), max(0, s.glass.y))
		pane.Left(16 + s.glass.x).Top(16 + s.glass.y).Children(func() {
			components.Text(c, tr(c, "gallery.liquid_glass")).Bold().TextColor(color)
			components.Text(c, tr(c, "gallery.onevra_mygo")).TextColor(color)
		})
	})
}

func glassStyleName(c *ui.Context, name string) string {
	if name == "Clear" {
		return tr(c, "gallery.clear_glass_style")
	}
	return tr(c, name)
}
