package flow

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// CardRadius rounds Card's corners; the selection ring follows it.
const CardRadius = 10

// CardOptions describes a Card. Width defaults to 240 DIPs. Title and
// Description each keep to one line. Header builds a small row above them,
// such as a status; Leading builds content before them, such as an icon
// tile; Footer builds a small muted row below, such as counts. A nonzero
// Border replaces the theme's border color, as for a status.
type CardOptions struct {
	Width                   float32
	Title, Description      string
	Header, Leading, Footer func(*ui.Context)
	Border                  ui.Color
}

// Card builds the default node face for Node.Content: a rounded surface
// with optional header, leading, title, description and footer content.
func Card(c *ui.Context, opts CardOptions) ui.Element {
	t := c.Theme()
	width := opts.Width
	if width <= 0 {
		width = 240
	}
	border := opts.Border
	if border == (ui.Color{}) {
		border = t.Border
	}
	small := t.Rem(.85)
	return ui.Column(c).Width(width).Padding(t.Space(3)).Gap(t.Space(2.5)).
		Radius(CardRadius).Background(t.Surface).Border(1, border).Children(func() {
		if opts.Header != nil {
			ui.Row(c).Gap(t.Space(1.5)).FontSize(small).Children(func() { opts.Header(c) })
		}
		if opts.Leading != nil || opts.Title != "" || opts.Description != "" {
			ui.Row(c).Gap(t.Space(2.5)).Children(func() {
				if opts.Leading != nil {
					opts.Leading(c)
				}
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(.5)).Children(func() {
					if opts.Title != "" {
						components.Text(c, opts.Title).SingleLine().FontWeight(600)
					}
					if opts.Description != "" {
						components.Text(c, opts.Description).SingleLine().FontSize(small).TextColor(t.TextMuted)
					}
				})
			})
		}
		if opts.Footer != nil {
			ui.Row(c).Gap(t.Space(3)).FontSize(small).TextColor(t.TextMuted).Children(func() { opts.Footer(c) })
		}
	})
}
