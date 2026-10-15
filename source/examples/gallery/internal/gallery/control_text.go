package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlText(c *ui.Context, s *controlState, variant string) {
	if variant == "Editable" {
		components.Text(c, s.text, components.TextOptions{Editable: &s.text}).Label(tr(c, "editableText")).FontSize(c.Theme().Rem(1.5)).Bold()
		return
	}
	ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(c.Theme().Space(4)).Children(func() {
		components.Text(c, tr(c, "aHeading")).FontSize(c.Theme().Rem(1.5)).Bold()
		components.Text(c, tr(c, "helloWorld"))
		components.Text(c, tr(c, "secondaryText")).TextColor(c.Theme().TextMuted)
		components.Text(c, tr(c, "longTextWrapsWithinTheAvailableWidthIncludingInNarrowWindows")).FillWidth()
		components.Text(c, tr(c, "aSingleLineLabel")).SingleLine()
		ui.Text(c, msg(c).LabelCount(tr(c, "count"), 42))
		components.RichText(c).Children(func() {
			components.Text(c, tr(c, "inline"))
			components.Text(c, tr(c, "boldText")).Bold()
		})
	})
}
