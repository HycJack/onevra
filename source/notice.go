package oneogo

import "github.com/egoist/mygo/ui"

// NoticeKind selects a semantic heading color, never the only status indicator.
type NoticeKind uint8

const (
	NoticeInfo NoticeKind = iota
	NoticeWarning
	NoticeError
	NoticeSuccess
)

// NoticeOptions supplies a caller-localized status title and optional details.
// Actions can build retry or dismiss controls; Notice never hides itself.
type NoticeOptions struct {
	Kind        NoticeKind
	Title       string
	Description string
	Actions     func()
}

// Notice renders persistent feedback using native text and caller-built controls.
// It is a labeled group, not an automatically announced live alert.
func Notice(c *ui.Context, options NoticeOptions) *ui.Element {
	t := c.Theme()
	color := t.Text
	switch options.Kind {
	case NoticeWarning:
		color = t.Warning
	case NoticeError:
		color = t.Danger
	case NoticeSuccess:
		color = t.Success
	}
	return ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(2)).Background(t.Surface).Radius(t.Radius).Role(ui.RoleGroup).Label(options.Title).Children(func() {
		ui.Text(c, options.Title).Bold().TextColor(color)
		if options.Description != "" {
			ui.Text(c, options.Description).Selectable()
		}
		if options.Actions != nil {
			options.Actions()
		}
	})
}
