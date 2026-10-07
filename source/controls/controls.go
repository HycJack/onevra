package controls

import "github.com/egoist/mygo/ui"

// dimmed applies the design style's disabled treatment. The bases own the
// disabled semantics; without their styleFn hook, opacity is the faithful
// visual equivalent.
func dimmed(e *ui.Element) *ui.Element {
	if e.IsDisabled() {
		e.Opacity(0.45)
	}
	return e
}
