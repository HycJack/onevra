package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlBadge(c *ui.Context, _ *controlState, variant string) {
	b := components.Badge(c, "New")
	if variant == "Accent" {
		b.Background(c.Theme().Accent).TextColor(c.Theme().AccentText).Border(1, c.Theme().Accent)
	}
}

func controlRichText(c *ui.Context, s *controlState, _ string) {
	components.RichText(c, ui.Span{Text: "Styled ", Weight: 600}).Children(func() {
		components.Text(c, "text and an ")
		if components.Button(c, "inline action", components.ButtonLink).Clicked() {
			s.index++
		}
		ui.Textf(c, " wrap as one paragraph. Clicks: %d", s.index)
	}).MaxWidth(320)
}

func controlBreadcrumbs(c *ui.Context, s *controlState, _ string) {
	path := []string{"Home", "Components", "Breadcrumbs"}
	b := components.Breadcrumbs(c, path, &s.index)
	b.Label("Path")
	if b.Changed() {
		s.text = path[s.index]
	}
	components.Text(c, s.text).TextColor(c.Theme().TextMuted)
}

func controlRating(c *ui.Context, s *controlState, _ string) {
	components.Rating(c, &s.index, 5).Label("Rating")
	ui.Textf(c, "%d / 5", s.index).TextColor(c.Theme().TextMuted)
}
