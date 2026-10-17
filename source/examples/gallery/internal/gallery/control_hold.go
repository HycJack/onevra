package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// controlHoldToConfirm keeps the confirmation in s.on; the hint beside the
// button follows the hold, and Reset returns the button to its label.
func controlHoldToConfirm(c *ui.Context, s *controlState, variant string) {
	t := c.Theme()
	label, done, options := "holdToArchive", "holdArchived", components.HoldToConfirmOptions{Icon: lucide.Archive()}
	switch variant {
	case "Danger":
		label, done, options = "holdToDeleteProject", "holdDeleted", components.HoldToConfirmOptions{Tone: components.HoldDanger}
	case "Neutral":
		label, done, options = "holdToSignOut", "holdSignedOut", components.HoldToConfirmOptions{Icon: lucide.LogOut(), Tone: components.HoldNeutral}
	}
	options.ConfirmedLabel = tr(c, done)
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
		h := components.HoldToConfirm(c.Key("hold"), &s.on, tr(c, label), options)
		switch {
		case h.Holding():
			components.Text(c, tr(c, "holdKeepHolding")).TextColor(t.TextMuted)
		case s.on:
			if components.Button(c, tr(c, "holdReset"), components.ButtonGhost).Clicked() {
				s.on = false
			}
		}
	})
}
