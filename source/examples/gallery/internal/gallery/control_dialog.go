package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, tr(c, "openDialog"), components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	components.Dialog(c, &s.dialogOpen, components.DialogOptions{Label: tr(c, "dialog"), DisableBackdropClose: variant == "Protected"}, func(panel ui.Element) {
		components.Text(c, tr(c, "aPageInADialog")).Bold()
		components.Text(c, tr(c, "openAnotherLayerOrClickTheBackdropToGoBackOneLayer"))
		if components.Button(c, tr(c, "openAnotherLayer"), components.ButtonSecondary).Clicked() {
			s.nestedOpen = true
		}
		if components.Button(c, tr(c, "done"), components.ButtonPrimary).Clicked() {
			s.dialogOpen = false
		}
		components.AlertDialog(c, &s.nestedOpen, components.AlertDialogOptions{Title: tr(c, "anotherLayer"), Message: tr(c, "thePageUnderneathStaysOpen"), ConfirmLabel: tr(c, "ok"), HideCancel: true})
	})
}
func controlAlertDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, msg(c).ShowS(tr(c, variant)), components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	result := components.AlertDialog(c, &s.dialogOpen, components.AlertDialogOptions{Title: tr(c, variant), Message: tr(c, "thisDialogReturnsYourChoice"), ConfirmLabel: tr(c, "ok"), CancelLabel: tr(c, "cancel"), HideCancel: variant == "Notice"})
	switch result {
	case components.DialogConfirmed:
		s.text = "Confirmed"
	case components.DialogCancelled:
		s.text = "Cancelled"
	}
	components.Text(c, tr(c, s.text))
}
