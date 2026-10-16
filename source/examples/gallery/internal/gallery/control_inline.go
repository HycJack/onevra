package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlBadge(c *ui.Context, _ *controlState, variant string) {
	b := components.Badge(c, tr(c, "new"))
	if variant == "Accent" {
		b.Background(c.Theme().Accent).TextColor(c.Theme().AccentText).Border(1, c.Theme().Accent)
	}
}

func controlRichText(c *ui.Context, s *controlState, _ string) {
	components.RichText(c, ui.Span{Text: tr(c, "styled"), Weight: 600}).Children(func() {
		components.Text(c, tr(c, "textAndAn"))
		if components.Button(c, tr(c, "inlineAction"), components.ButtonLink).Clicked() {
			s.index++
		}
		ui.Text(c, msg(c).WrapAsOneParagraphClicksD(s.index))
	}).MaxWidth(320)
}

func controlBreadcrumbs(c *ui.Context, s *controlState, _ string) {
	path := []string{tr(c, "home"), tr(c, "components"), "Breadcrumbs"}
	b := components.Breadcrumbs(c, path, &s.index)
	b.Label(tr(c, "path"))
	if b.Changed() {
		s.text = []string{"Home", "Components", "Breadcrumbs"}[s.index]
	}
	if s.text != "" {
		components.Text(c, tr(c, s.text)).TextColor(c.Theme().TextMuted)
	}
}

func controlRating(c *ui.Context, s *controlState, _ string) {
	components.Rating(c, &s.index, 5).Label(tr(c, "rating"))
	ui.Text(c, msg(c).RatingOutOfFive(s.index)).TextColor(c.Theme().TextMuted)
}
