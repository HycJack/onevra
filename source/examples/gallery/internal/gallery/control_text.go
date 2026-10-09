package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlText(c *ui.Context, s *controlState, variant string) {
	if variant == "Editable" {
		components.Text(c, s.text, components.TextOptions{Editable: &s.text}).Label(tr(c, "gallery.editable_text")).FontSize(c.Theme().Rem(1.5)).Bold()
		return
	}
	ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(c.Theme().Space(4)).Children(func() {
		components.Text(c, tr(c, "gallery.a_heading")).FontSize(c.Theme().Rem(1.5)).Bold()
		components.Text(c, tr(c, "gallery.hello_world"))
		components.Text(c, tr(c, "gallery.secondary_text")).TextColor(c.Theme().TextMuted)
		components.Text(c, tr(c, "gallery.long_text_wraps_within_the_available_width_including_in_narrow_windows")).FillWidth()
		components.Text(c, tr(c, "gallery.a_single_line_label")).SingleLine()
		ui.Textf(c, tr(c, "gallery.s_d"), tr(c, "gallery.count"), 42)
		components.RichText(c).Children(func() {
			components.Text(c, tr(c, "gallery.inline"))
			components.Text(c, tr(c, "gallery.bold_text")).Bold()
		})
	})
}
