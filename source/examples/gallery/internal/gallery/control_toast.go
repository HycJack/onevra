package gallery

import (
	"strings"
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
		components.Text(c, msg(c).ToastItems(items)).TextColor(c.Theme().TextMuted)
		if components.Button(c, tr(c, "delete"), components.ButtonSecondary).Disabled(items == 0).Clicked() {
			s.index++
			c.AddToast(ui.Toast{Title: tr(c, "toastDeleted"), Action: tr(c, "undo"),
				Timeout: 8 * time.Second, OnAction: func() { s.index = max(0, s.index-1) }})
		}
	case "Progress":
		ui.Row(c).Gap(8).Wrap().Children(func() {
			if components.Button(c, tr(c, "toastStartUpload"), components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{ID: "gallery-upload", Title: tr(c, "toastUploading"), Timeout: -1,
					Data: components.ToastContent{Leading: func() {
						size := c.Theme().Rem(1.125)
						components.Spinner(c, components.SpinnerRing).Size(size, size)
					}}})
			}
			if components.Button(c, tr(c, "toastFinishUpload"), components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{ID: "gallery-upload", Title: tr(c, "toastUploaded"), Type: components.ToastSuccess})
			}
		})
	default:
		ui.Row(c).Gap(8).Wrap().Children(func() {
			for _, kind := range []string{components.ToastInfo, components.ToastSuccess, components.ToastWarning, components.ToastError} {
				if components.Button(c, tr(c, "toast"+strings.ToUpper(kind[:1])+kind[1:]), components.ButtonSecondary).Clicked() {
					c.AddToast(ui.Toast{Title: tr(c, "toast"+strings.ToUpper(kind[:1])+kind[1:]+"Title"),
						Description: tr(c, "toastStatusDescription"), Type: kind})
				}
			}
		})
	}
}

var toastPlacements = []struct {
	placement components.ToastPlacement
	label     string
}{
	{components.ToastTop, "toastTop"}, {components.ToastTopStart, "toastTopStart"},
	{components.ToastTopEnd, "toastTopEnd"}, {components.ToastBottom, "toastBottom"},
	{components.ToastBottomStart, "toastBottomStart"}, {components.ToastBottomEnd, "toastBottomEnd"},
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
