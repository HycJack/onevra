package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlTagInput(c *ui.Context, s *controlState, _ string) {
	t := c.Theme()
	components.TagInput(c, &s.tags).Label(tr(c, "topics")).Placeholder(tr(c, "addATag")).FillWidth().MinWidth(0)
	components.Text(c, tr(c, "pressEnterOrCommaToAdd")).FontSize(t.Rem(.85)).TextColor(t.TextMuted)
}
