package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlTooltip(c *ui.Context, _ *controlState, _ string) {
	anchor := components.Button(c, tr(c, "hoverOrFocusMe"), components.ButtonSecondary)
	components.Tooltip(c, anchor, tr(c, "aHelpfulDescriptionOfThisAction"))
}
