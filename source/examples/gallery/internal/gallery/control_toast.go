package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// Shell builds the Onevra toast viewport, so examples only add toasts.
func controlToast(c *ui.Context, s *controlState, variant string) {
	switch variant {
	case "Actions":
		// index counts deleted items; Undo restores them.
		items := 3 - s.index
		components.Text(c, fmt.Sprintf(tr(c, "gallery.toast_items"), items)).TextColor(c.Theme().TextMuted)
		if components.Button(c, tr(c, "gallery.delete"), components.ButtonSecondary).Disabled(items == 0).Clicked() {
			s.index++
			c.AddToast(ui.Toast{Title: tr(c, "gallery.toast_deleted"), Action: tr(c, "gallery.undo"),
				Timeout: 8 * time.Second, OnAction: func() { s.index = max(0, s.index-1) }})
		}
	case "Progress":
		ui.Row(c).Gap(8).Wrap().Children(func() {
			if components.Button(c, tr(c, "gallery.toast_start_upload"), components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{ID: "gallery-upload", Title: tr(c, "gallery.toast_uploading"), Timeout: -1,
					Data: components.ToastContent{Leading: func() {
						size := c.Theme().Rem(1.125)
						components.Spinner(c, components.SpinnerRing).Size(size, size)
					}}})
			}
			if components.Button(c, tr(c, "gallery.toast_finish_upload"), components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{ID: "gallery-upload", Title: tr(c, "gallery.toast_uploaded"), Type: components.ToastSuccess})
			}
		})
	default:
		ui.Row(c).Gap(8).Wrap().Children(func() {
			for _, kind := range []string{components.ToastInfo, components.ToastSuccess, components.ToastWarning, components.ToastError} {
				if components.Button(c, tr(c, "gallery.toast_"+kind), components.ButtonSecondary).Clicked() {
					c.AddToast(ui.Toast{Title: tr(c, "gallery.toast_"+kind+"_title"),
						Description: tr(c, "gallery.toast_status_description"), Type: kind})
				}
			}
		})
	}
}

var toastPlacements = []struct {
	placement components.ToastPlacement
	label     string
}{
	{components.ToastTop, "gallery.toast_top"}, {components.ToastTopStart, "gallery.toast_top_start"},
	{components.ToastTopEnd, "gallery.toast_top_end"}, {components.ToastBottom, "gallery.toast_bottom"},
	{components.ToastBottomStart, "gallery.toast_bottom_start"}, {components.ToastBottomEnd, "gallery.toast_bottom_end"},
}

// toastPlacement applies the window's toast placement through Shell, which
// owns the viewport; the choice stays in memory for this window.
func (g *Gallery) toastPlacement(c *ui.Context) {
	ui.Row(c).Gap(16).Wrap().Children(func() {
		for _, option := range toastPlacements {
			if components.Radio(c, &g.toasts.Placement, option.placement, tr(c, option.label)).Changed() && g.Shell != nil {
				g.Shell.SetToasts(g.toasts)
			}
		}
	})
}
