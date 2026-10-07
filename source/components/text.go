package components

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// Text creates display text using MyGo's text primitive. It inherits the frame
// theme and parent text styles without overriding font, size or color. MyGo
// owns shaping, wrapping and rendering. The returned element supports chaining,
// such as .Bold(), .TextColor(...) and .SingleLine(), and RichText composition.
func Text(c *ui.Context, value string) *ui.Element {
	return ui.Text(c, value)
}

// Textf creates display text formatted with fmt.Sprintf, with the same behavior
// as Text. Callers own the format string and localized content.
func Textf(c *ui.Context, format string, args ...any) *ui.Element {
	return Text(c, fmt.Sprintf(format, args...))
}
