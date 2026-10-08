package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, "Open dialog", components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	components.Dialog(c, &s.dialogOpen, components.DialogOptions{Label: "Dialog", DisableBackdropClose: variant == "Protected"}, func(panel ui.Element) {
		components.Text(c, "A page in a dialog").Bold()
		components.Text(c, "Open another layer, or click the backdrop to go back one layer.")
		if components.Button(c, "Open another layer", components.ButtonSecondary).Clicked() {
			s.nestedOpen = true
		}
		if components.Button(c, "Done", components.ButtonPrimary).Clicked() {
			s.dialogOpen = false
		}
		components.AlertDialog(c, &s.nestedOpen, components.AlertDialogOptions{Title: "Another layer", Message: "The page underneath stays open.", ConfirmLabel: "OK", HideCancel: true})
	})
}
func controlAlertDialog(c *ui.Context, s *controlState, variant string) {
	if components.Button(c, "Show "+variant, components.ButtonSecondary).Clicked() {
		s.dialogOpen = true
	}
	result := components.AlertDialog(c, &s.dialogOpen, components.AlertDialogOptions{Title: variant, Message: "This dialog returns your choice.", ConfirmLabel: "OK", CancelLabel: "Cancel", HideCancel: variant == "Notice"})
	switch result {
	case components.DialogConfirmed:
		s.text = "Confirmed"
	case components.DialogCancelled:
		s.text = "Cancelled"
	}
	components.Text(c, s.text)
}
