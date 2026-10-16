package gallery

import (
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

// flyoutDemo retains application state separately from the transient window.
type flyoutDemo struct {
	gallery *Gallery
	window  *mygo.Window
	note    string
}

func (d *flyoutDemo) close() {
	if d.window != nil {
		d.window.Destroy()
	}
	d.window = nil
}

func (d *flyoutDemo) show(bounds ui.Rect, popover bool) {
	if d.gallery == nil || d.gallery.Shell == nil {
		return
	}
	d.close()
	parent := d.gallery.Shell.Window()
	w := mygo.NewFlyout(mygo.FlyoutOptions{
		Parent:    parent,
		Anchor:    mygo.Rectangle{X: int(bounds.X), Y: int(bounds.Y), Width: int(bounds.W), Height: int(bounds.H)},
		Placement: mygo.PlacementBottomStart, Gap: 6,
		Width: 300, Height: 180, Focusable: true, Shadow: true, Popover: popover,
		Content: ui.View(d.view),
	})
	d.window = w
	w.OnClosed(func() {
		if d.window == w {
			d.window = nil
		}
		parent.Invalidate()
	})
}

func (d *flyoutDemo) view(c *ui.Context) {
	theme.Apply(c)
	c.Root().Background(ui.Transparent)
	d.gallery.withLanguage(c, func() {
		t := c.Theme()
		panel := ui.Column(c).Fill().Padding(16).Gap(12)
		if !c.Vibrancy() {
			panel.Background(t.Surface)
			switch runtime.GOOS {
			case "darwin":
				panel.Radius(10)
			case "windows":
			default:
				panel.Radius(10).Border(1, t.Border)
			}
		}
		panel.Children(func() {
			components.TextBox(c, &d.note, components.TextBoxOptions{}).Label(tr(c, "notes")).FillWidth().AutoFocus()
			components.Button(c, tr(c, "close"), components.ButtonOutline).OnClick(func() {
				if d.window != nil {
					d.window.Close()
				}
			})
		})
	})
}

func exampleFlyout(c *ui.Context, s *demoState) {
	d := &s.flyout
	components.Text(c, tr(c, "flyoutDescription")).Wrap()
	for i, label := range []string{"flyoutEditor", "flyoutPopover"} {
		b := components.Button(c, tr(c, label), components.ButtonOutline).Disabled(d.gallery == nil || d.gallery.Shell == nil)
		if b.Clicked() {
			d.show(b.Bounds(), i == 1)
		}
	}
	components.Text(c, d.note).Wrap()
}
