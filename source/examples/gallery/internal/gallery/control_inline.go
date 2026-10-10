package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func controlBadge(c *ui.Context, _ *controlState, variant string) {
	b := components.Badge(c, tr(c, "gallery.new"))
	if variant == "Accent" {
		b.Background(c.Theme().Accent).TextColor(c.Theme().AccentText).Border(1, c.Theme().Accent)
	}
}

func controlRichText(c *ui.Context, s *controlState, _ string) {
	components.RichText(c, ui.Span{Text: tr(c, "gallery.styled"), Weight: 600}).Children(func() {
		components.Text(c, tr(c, "gallery.text_and_an"))
		if components.Button(c, tr(c, "gallery.inline_action"), components.ButtonLink).Clicked() {
			s.index++
		}
		ui.Textf(c, tr(c, "gallery.wrap_as_one_paragraph_clicks_d"), s.index)
	}).MaxWidth(320)
}

func controlBreadcrumbs(c *ui.Context, s *controlState, _ string) {
	path := []string{tr(c, "gallery.home"), tr(c, "gallery.components"), "Breadcrumbs"}
	b := components.Breadcrumbs(c, path, &s.index)
	b.Label(tr(c, "gallery.path"))
	if b.Changed() {
		s.text = []string{"Home", "Components", "Breadcrumbs"}[s.index]
	}
	components.Text(c, tr(c, s.text)).TextColor(c.Theme().TextMuted)
}

func controlRating(c *ui.Context, s *controlState, _ string) {
	components.Rating(c, &s.index, 5).Label(tr(c, "gallery.rating"))
	ui.Textf(c, tr(c, "gallery.d_5"), s.index).TextColor(c.Theme().TextMuted)
}
