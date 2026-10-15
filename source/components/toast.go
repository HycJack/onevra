package components

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// Toast types for ui.Toast.Type select a colored leading status icon. Other
// types, including the empty one, show no icon.
const (
	ToastInfo    = "info"
	ToastSuccess = "success"
	ToastWarning = "warning"
	ToastError   = "error"
)

// ToastContent is optional ui.Toast.Data for content slots. Leading replaces
// the type's status icon, for example with an image or a progress indicator; it
// runs synchronously while Toasts builds the toast.
type ToastContent struct {
	Leading func()
}

// ToastPlacement selects the window edge and side where toasts stack. Bottom
// placements put the newest toast last, nearest the bottom edge; top
// placements put it first, nearest the top edge, and stack older ones downward.
type ToastPlacement int

const (
	ToastBottom ToastPlacement = iota
	ToastBottomStart
	ToastBottomEnd
	ToastTop
	ToastTopStart
	ToastTopEnd
)

// Top reports whether p places toasts at the top edge.
func (p ToastPlacement) Top() bool { return p >= ToastTop && p <= ToastTopEnd }

func (p ToastPlacement) align() ui.Align {
	switch p {
	case ToastBottomStart, ToastTopStart:
		return ui.Start
	case ToastBottomEnd, ToastTopEnd:
		return ui.End
	}
	return ui.Center
}

// ToastOptions customizes where Toasts places the window's toasts. Inset adds
// DIPs at the placement's edge beyond the standard margin, for example to clear
// a title bar; Shell adds its header height to top placements itself.
type ToastOptions struct {
	Placement ToastPlacement
	Inset     float32
}

// Toasts builds the window's MyGo toasts in the Onevra style, replacing MyGo's
// default look, and returns the viewport. Applications show and close toasts
// with MyGo's Context.Toast, ToastAction, AddToast and CloseToast; MyGo owns the
// queue, timing, pausing, Escape dismissal and announcements.
//
// Shell calls Toasts after its content with ShellOptions.Toasts, so
// applications call it only in standalone hosts; the first options value
// applies. A view that builds its own ToastViewportBase earlier in the frame
// keeps its look. Toasts stack at the bottom center by default, each with an
// optional action and a close button named by the "onevra.closeToast" message.
func Toasts(c *ui.Context, options ...ToastOptions) ui.Element {
	var o ToastOptions
	if len(options) > 0 {
		o = options[0]
	}
	t := c.Theme()
	closeLabel := localized(c).CloseToast
	width, _ := c.Size()
	margin, inset := t.Space(4), max(0, o.Inset)
	enter := t.Space(3)
	if o.Placement.Top() {
		enter = -enter
	}
	motion := ui.ElementTransition{Enter: &ui.Motion{Y: enter}, Exit: &ui.Motion{Scale: .96}}
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.AlignItems(o.Placement.align()).Gap(t.Space(2))
		if o.Placement.Top() {
			viewport.Justify(ui.Start).Padding(margin+inset, margin, margin, margin)
		} else {
			viewport.Padding(margin, margin, margin+inset, margin)
		}
		for i := range toasts {
			toast := toasts[i]
			if o.Placement.Top() {
				toast = toasts[len(toasts)-1-i]
			}
			parts := ui.ToastBase(c, toast)
			align := ui.Center
			if toast.Description != "" {
				align = ui.Start
			}
			parts.Root.Row().AlignItems(align).Gap(t.Space(3)).
				MaxWidth(min(max(0, width-2*margin), t.Space(80))).
				Padding(t.Space(3), t.Space(2), t.Space(3), t.Space(4)).Radius(t.Radius*1.5).
				Material(glass.Glass{Tint: t.Background.Alpha(.92)}).TextColor(t.Text).Border(1, t.Border).
				Shadow(0, t.Space(1), t.Space(4), 0, ui.RGBA(0, 0, 0, .14)).
				Transition(motion)
			parts.Root.Children(func() {
				toastLeading(c, toast)
				ui.Column(c).Grow(1).Shrink(1).MinWidth(0).Gap(t.Space(.5)).Children(func() {
					Text(c, toast.Title).FontWeight(500)
					if toast.Description != "" {
						Text(c, toast.Description).FontSize(t.FontSize - 1).TextColor(t.TextMuted)
					}
				})
				if toast.Action != "" {
					styleButton(c, parts.ActionButton(), ButtonSecondary).Shrink(0).
						Label(toast.Action).Children(func() { Text(c, toast.Action).SingleLine() })
				}
				size := t.Space(6)
				styleButton(c, parts.CloseButton(), ButtonGhost).Label(closeLabel).
					Size(size, size).Padding(0).Shrink(0).TextColor(t.TextMuted).Children(func() {
					Icon(c, lucide.X).Size(t.Rem(.875), t.Rem(.875))
				})
			})
		}
	})
}

func toastLeading(c *ui.Context, toast ui.Toast) {
	if content, ok := toast.Data.(ToastContent); ok && content.Leading != nil {
		ui.Row(c).Shrink(0).AlignItems(ui.Center).Children(content.Leading)
		return
	}
	t := c.Theme()
	var icon func() *ui.SVG
	var color ui.Color
	switch toast.Type {
	case ToastInfo:
		icon, color = lucide.Info, t.Accent
	case ToastSuccess:
		icon, color = lucide.CircleCheck, t.Success
	case ToastWarning:
		icon, color = lucide.TriangleAlert, t.Warning
	case ToastError:
		icon, color = lucide.CircleX, t.Danger
	default:
		return
	}
	size := t.Rem(1.125)
	Icon(c, icon).Size(size, size).Shrink(0).TextColor(color).Role(ui.RoleNone)
}
