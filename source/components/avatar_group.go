package components

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/internal/l10n"
)

// AvatarGroupMember is one person in an AvatarGroup.
type AvatarGroupMember struct {
	// Name is the accessible name, hover hint and source of the initials.
	// Members keep their identity by name, and by order among equal names.
	Name   string
	Image  ui.ImageSource
	Status AvatarStatus
}

// AvatarGroupOptions customize an AvatarGroup.
type AvatarGroupOptions struct {
	// Max is how many avatars show before the rest collapse into one "+N"
	// chip: 4 when zero; a negative Max collapses everyone.
	Max int
	// Size is the avatar and chip diameter in DIPs; zero uses Avatar's
	// standard size.
	Size float32
	// Label names the group for assistive technology; empty uses the
	// localized "Team members".
	Label string
}

const avatarGroupDuration = 200 * time.Millisecond

// avatarGroupKey keys a member's slot by name and occurrence of that name.
type avatarGroupKey struct {
	name string
	n    int
}

// avatarGroupOverflowKey keys the overflow chip's slot.
type avatarGroupOverflowKey struct{}

// avatarGroupCount keeps the last overflow count and the way it moved.
type avatarGroupCount struct{ count, direction int }

// AvatarGroup creates an overlapping stack of members in order. Members past
// Max collapse into an overflow chip, so long lists build only what shows.
// Pointing at the stack loosens it and lifts the member under the pointer,
// which names itself in a tooltip; members joining or leaving open or close
// their slots and the overflow count rolls the way it moved. Reduced Motion
// keeps the stack still. Applications own members, photos and localized
// labels. Key a parent for identity.
func AvatarGroup(c *ui.Context, members []AvatarGroupMember, options AvatarGroupOptions) ui.Element {
	t := c.Theme()
	d := avatarDiameter(c, options.Size)
	limit := options.Max
	if limit == 0 {
		limit = 4
	}
	limit = min(max(0, limit), len(members))
	overflow := len(members) - limit
	label := options.Label
	if label == "" {
		label = l10n.In(c.Locale()).Text("avatarGroup")
	}
	overlap := d / 6
	reduce := c.Preferences().ReduceMotion
	slotMotion := ui.ElementTransition{
		Animation: &ui.Smooth,
		Enter:     &ui.Motion{Collapse: true, Scale: .9},
		Exit:      &ui.Motion{Collapse: true, Scale: .9},
	}
	group := ui.Row(c).Height(d).Shrink(0).AlignSelf(ui.Start).PaddingStart(overlap).
		Role(ui.RoleGroup).Label(label)
	fanned := group.Hovered() && !reduce
	group.Children(func() {
		var slots []ui.Element
		hovered := -1
		slot := func(key any) ui.Element {
			s := ui.Row(c).Key(key).Height(d).Shrink(0).MarginStart(-overlap).Transition(slotMotion)
			if s.Hovered() {
				hovered = len(slots)
			}
			slots = append(slots, s)
			return s
		}
		ring := func(e ui.Element, hover bool) {
			color := t.Border
			if hover {
				color = t.TextMuted
			}
			e.Shadow(0, 0, 0, 1, color)
		}
		seen := map[string]int{}
		for _, m := range members[:limit] {
			key := avatarGroupKey{m.Name, seen[m.Name]}
			seen[m.Name]++
			s := slot(key)
			hover := hovered == len(slots)-1
			s.Children(func() {
				a := avatar(c, m.Name, AvatarOptions{Image: m.Image, Status: m.Status, Size: d}, 2, t.Surface)
				ring(a, hover)
				// The tooltip repeats the accessible name, so it adds no description.
				Tooltip(c, a, m.Name)
				a.Description("")
			})
		}
		if overflow > 0 {
			s := slot(avatarGroupOverflowKey{})
			hover := hovered == len(slots)-1
			s.Children(func() {
				chip := ui.Box(c).Size(d, d).Shrink(0).Center().Radius(d/2).Clip().
					Background(t.SurfaceHover).TextColor(t.TextMuted).Border(2, t.Surface).
					Role(ui.RoleImage).Label(l10n.In(c.Locale()).Count("avatarGroupMore", overflow))
				ring(chip, hover)
				last := ui.Local(chip, "count", func() avatarGroupCount { return avatarGroupCount{overflow, 1} })
				if last.count != overflow {
					last.direction = 1
					if overflow < last.count {
						last.direction = -1
					}
					last.count = overflow
				}
				// More people rise in from below, fewer drop in from above.
				shift := d * .2 * float32(last.direction)
				chip.Children(func() {
					ui.Text(c, "+"+strconv.Itoa(overflow)).Key(overflow).SingleLine().Unselectable().
						FontSize(avatarFontSize(d)).FontWeight(500).
						Transition(ui.ElementTransition{
							Duration: avatarGroupDuration,
							Enter:    &ui.Motion{Y: shift},
							Exit:     &ui.Motion{Y: -shift},
						})
				})
			})
		}
		if reduce {
			return
		}
		// Hovering loosens the overlap from the center; the hovered member
		// rises and its neighbors ease aside.
		for i, s := range slots {
			var x, y float32
			if fanned {
				x = (float32(i) - float32(len(slots)-1)/2) * overlap * 2 / 3
			}
			switch {
			case hovered < 0:
			case i == hovered:
				y = -2
			case i > hovered:
				x += overlap / 2
			default:
				x -= overlap / 2
			}
			motion := ui.Animation{Duration: avatarGroupDuration, Ease: ui.EaseInOut}
			s.Translate(s.AnimateTo("fanX", x, motion), s.AnimateTo("fanY", y, motion))
		}
	})
	return group
}
