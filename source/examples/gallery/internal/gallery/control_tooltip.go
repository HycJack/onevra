package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTooltip(c *ui.Context, _ *controlState, _ string) {
	anchor := components.Button(c, tr(c, "gallery.hover_or_focus_me"), components.ButtonSecondary)
	components.Tooltip(c, anchor, tr(c, "gallery.a_helpful_description_of_this_action"))
}
