package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTooltip(c *ui.Context, _ *controlState, _ string) {
	anchor := components.Button(c, "Hover or focus me", components.ButtonSecondary)
	components.Tooltip(c, anchor, "A helpful description of this action")
}
