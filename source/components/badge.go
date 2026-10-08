package components

import "github.com/egoist/mygo/ui"

// Badge creates a compact, non-interactive status label. Applications own its
// text and may override the inherited foreground and background colors.
func Badge(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	return Text(c, text).SingleLine().Unselectable().AlignSelf(ui.Start).Shrink(0).
		FontSize(t.Rem(.85)).FontWeight(500).
		Padding(t.Space(.5), t.Space(2)).Radius(t.Space(3)).
		Background(t.Surface).TextColor(t.Text).Border(1, t.Border)
}
