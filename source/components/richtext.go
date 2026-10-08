package components

import "github.com/egoist/mygo/ui"

// RichText composes styled spans and inline Text or ButtonLink children into a
// single wrapping paragraph. Zero span fields inherit the surrounding style.
// MyGo's unstyled text primitive owns shaping, selection, inline interaction
// and accessibility; only text elements belong in Children.
func RichText(c *ui.Context, spans ...ui.Span) ui.Element {
	return ui.Text(c, "").Children(func() {
		for _, span := range spans {
			e := Text(c, span.Text).Element
			if span.Font != "" {
				e.Font(span.Font)
			}
			if span.Size != 0 {
				e.FontSize(span.Size)
			}
			if span.Weight != 0 {
				e.FontWeight(span.Weight)
			}
			if span.Color != (ui.Color{}) {
				e.TextColor(span.Color)
			}
			if span.Background != (ui.Color{}) {
				e.Background(span.Background)
			}
			if span.Italic {
				e.Italic()
			}
			if span.Underline {
				e.Underline()
			}
			if span.WavyUnderline {
				e.WavyUnderline()
			}
			if span.Strikethrough {
				e.Strikethrough()
			}
			if span.DecorationColor != (ui.Color{}) {
				e.DecorationColor(span.DecorationColor)
			}
			if span.DecorationThickness != 0 {
				e.DecorationThickness(span.DecorationThickness)
			}
			if span.LetterSpacing != 0 {
				e.LetterSpacing(span.LetterSpacing)
			}
			if span.Features != "" {
				e.FontFeatures(span.Features)
			}
		}
	})
}
