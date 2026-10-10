package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, tr(c, "gallery.open_dialog"), components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	components.Dialog(c, &s.dialogOpen, components.DialogOptions{Label: tr(c, "gallery.dialog"), DisableBackdropClose: variant == "Protected"}, func(panel ui.Element) {
		components.Text(c, tr(c, "gallery.a_page_in_a_dialog")).Bold()
		components.Text(c, tr(c, "gallery.open_another_layer_or_click_the_backdrop_to_go_back_one_layer"))
		if components.Button(c, tr(c, "gallery.open_another_layer"), components.ButtonSecondary).Clicked() {
			s.nestedOpen = true
		}
		if components.Button(c, tr(c, "gallery.done"), components.ButtonPrimary).Clicked() {
			s.dialogOpen = false
		}
		components.AlertDialog(c, &s.nestedOpen, components.AlertDialogOptions{Title: tr(c, "gallery.another_layer"), Message: tr(c, "gallery.the_page_underneath_stays_open"), ConfirmLabel: tr(c, "gallery.ok"), HideCancel: true})
	})
}
func controlAlertDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, fmt.Sprintf(tr(c, "gallery.show_s"), tr(c, variant)), components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	result := components.AlertDialog(c, &s.dialogOpen, components.AlertDialogOptions{Title: tr(c, variant), Message: tr(c, "gallery.this_dialog_returns_your_choice"), ConfirmLabel: tr(c, "gallery.ok"), CancelLabel: tr(c, "gallery.cancel"), HideCancel: variant == "Notice"})
	switch result {
	case components.DialogConfirmed:
		s.text = "Confirmed"
	case components.DialogCancelled:
		s.text = "Cancelled"
	}
	components.Text(c, tr(c, s.text))
}
