package components

import (
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/icons/lucide"
)

// NumberInputControl exposes layout on Element and editing on the text field.
// Keep this handle to query Changed after applying Element layout methods.
type NumberInputControl struct {
	ui.Element
	input, increase, decrease ui.Element
	changed                   bool
}

// Changed reports an input-driven value change, not initial clamping.
func (n *NumberInputControl) Changed() bool { return n.changed }

// Label names the text field and its two step buttons.
func (n *NumberInputControl) Label(label string) *NumberInputControl {
	n.input.Label(label)
	n.increase.Label(strings.TrimSpace(label + " Increase"))
	n.decrease.Label(strings.TrimSpace(label + " Decrease"))
	return n
}

// AutoFocus focuses the editable text field.
func (n *NumberInputControl) AutoFocus() *NumberInputControl {
	n.input.AutoFocus()
	return n
}

type numberInputState struct {
	text  string
	value float64
}

// NumberInput creates a text field with an integrated vertical stepper on the
// right. Bounds and value must be finite, hi >= lo, and step finite and positive.
// Valid in-range text updates value immediately without snapping to step. Partial
// or invalid text is retained while editing and restored on Enter or focus loss.
// Up/Down and the buttons add/subtract step and clamp to the bounds. Application
// changes replace even a focused draft. Disable a parent to disable the control.
// Key a parent for stable identity; construction resolves interaction state.
func NumberInput(c *ui.Context, value *float64, lo, hi, step float64) *NumberInputControl {
	finite := func(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
	if value == nil || !finite(*value) || !finite(lo) || !finite(hi) || hi < lo || !finite(step) || step <= 0 {
		panic("components.NumberInput: finite value and ordered bounds, and positive finite step required")
	}
	*value = max(lo, min(hi, *value))
	t := c.Theme()
	n := &NumberInputControl{Element: ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Stretch).Width(t.Space(32)).Role(ui.RoleGroup)}
	format := func() string { return strconv.FormatFloat(*value, 'f', -1, 64) }
	state := ui.Local(n.Element, "number-input", func() numberInputState {
		return numberInputState{text: format(), value: *value}
	})
	if state.value != *value {
		state.text = format()
	}
	set := func(v float64) {
		v = max(lo, min(hi, v))
		if *value != v {
			*value = v
			n.changed = true
		}
	}
	advance := func(delta float64) {
		// Decimal arithmetic through shortest decimal representations avoids
		// accumulating binary tails for ordinary fractional steps.
		v := *value + delta
		precision := func(x float64) int {
			s := strconv.FormatFloat(x, 'f', -1, 64)
			if i := strings.IndexByte(s, '.'); i >= 0 {
				return len(s) - i - 1
			}
			return 0
		}
		p := math.Pow10(max(precision(*value), precision(step)))
		if finite(v*p) && finite(p) {
			v = math.Round(v*p) / p
		}
		set(v)
		state.text = format()
	}
	n.Element.Children(func() {
		n.input = TextBox(c, &state.text, TextBoxOptions{}).FillWidth().MinWidth(0)
		if n.input.Changed() {
			if v, err := strconv.ParseFloat(strings.TrimSpace(state.text), 64); err == nil && finite(v) && v >= lo && v <= hi {
				set(v)
			}
		}
		if n.input.Shortcut(0, ui.KeyUp) {
			advance(step)
		}
		if n.input.Shortcut(0, ui.KeyDown) {
			advance(-step)
		}
		if !n.input.Focused() || n.input.Submitted() {
			state.text = format()
		}
		// The text field alone determines intrinsic height. The step buttons
		// fill that height without contributing their own preferred sizes.
		ui.Box(c).Width(t.Space(4)).Shrink(0).Children(func() {
			ui.Column(c).Absolute().Top(0).Bottom(0).Left(0).Right(0).Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Children(func() {
				button := func(up bool) ui.Element {
					b := ui.ButtonBase(c).FillWidth().Height(0).Grow(1).Padding(0).Radius(t.Radius / 2)
					b.Disabled((up && *value >= hi) || (!up && *value <= lo))
					if !b.IsDisabled() {
						if b.Pressed() {
							b.Background(t.SurfacePressed)
						} else if b.Hovered() {
							b.Background(t.SurfaceHover)
						}
					}
					b.DrawOver(func(p *ui.Painter, r ui.Rect) {
						icon := lucide.ChevronDown()
						if up {
							icon = lucide.ChevronUp()
						}
						size := t.Space(2.5)
						p.Icon(icon, ui.Rect{X: r.X + (r.W-size)/2, Y: r.Y + (r.H-size)/2, W: size, H: size}, t.Text)
					})
					return dimmed(b)
				}
				n.increase = button(true)
				ui.Box(c).FillWidth().Height(1).Background(t.Border)
				n.decrease = button(false)
				if n.increase.Clicked() {
					advance(step)
				}
				if n.decrease.Clicked() {
					advance(-step)
				}
			})
		})
		n.input.Range(lo, hi, *value).Step(step)
	})
	state.value = *value
	n.Label("")
	return n
}
