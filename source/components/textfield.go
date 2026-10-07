package components

import "github.com/egoist/mygo/ui"

// TextField creates a single-line text input editing *value. It returns the
// element for chaining, as .Label("Name").Placeholder("Ada").
func TextField(c *ui.Context, value *string) *ui.Element {
	return textField(c, value, false)
}

// TextArea creates a multi-line text input editing *value. It returns the
// element for chaining, as .Label("Notes").Placeholder("…").
func TextArea(c *ui.Context, value *string) *ui.Element {
	return textField(c, value, true)
}

func textField(c *ui.Context, value *string, multiline bool) *ui.Element {
	t := c.Theme()
	var e *ui.Element
	if multiline {
		e = ui.TextAreaBase(c, value)
		e.MinHeight(t.Space(20))
	} else {
		e = ui.TextInputBase(c, value)
	}
	border := t.Border
	switch {
	case e.Focused():
		border = t.Accent
	case e.Hovered():
		border = t.Border.Mix(t.Text, 0.25)
	}
	e.Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Background(t.Surface).Border(1, border)
	return dimmed(e)
}
