package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// floatingGroupCopied is how long Share reads Copied after a click.
const floatingGroupCopied = 1500 * time.Millisecond

// floatingGroupState is a sample board's history, link and presentation.
type floatingGroupState struct {
	edits, undone  int
	presenting     bool
	copiedAt       time.Time
	zoom           int
	locked, hidden bool
}

func controlFloatingButtonGroup(c *ui.Context, s *controlState, variant string) {
	st := &s.floating
	if variant == "Vertical" {
		floatingToolsRail(c, st)
		return
	}
	group := components.FloatingButtonGroupMuted
	if variant == "Raised" {
		group = components.FloatingButtonGroupRaised
	}
	share := tr(c, "share")
	if !st.copiedAt.IsZero() {
		if left := floatingGroupCopied - c.Now().Sub(st.copiedAt); left > 0 {
			share = tr(c, "copied")
			c.After(left)
		} else {
			st.copiedAt = time.Time{}
		}
	}
	present := tr(c, "present")
	if st.presenting {
		present = tr(c, "stopPresenting")
	}
	items := []components.FloatingButtonGroupItem{
		{ID: "undo", Label: tr(c, "undo"), Icon: lucide.Undo2, IconOnly: true, Shortcut: "⌘Z", Disabled: st.edits == 0},
		{ID: "redo", Label: tr(c, "redo"), Icon: lucide.Redo2, IconOnly: true, Shortcut: "⇧⌘Z", Disabled: st.undone == 0},
		{Separator: true},
		{ID: "comment", Label: tr(c, "comment"), Icon: lucide.MessageCircle, Shortcut: "C"},
		{ID: "share", Label: share, Icon: lucide.Link2, ReserveLabels: []string{tr(c, "share"), tr(c, "copied")}},
		{Separator: true},
		{ID: "present", Label: present, Icon: lucide.Play, Toggle: true, On: st.presenting,
			ReserveLabels: []string{tr(c, "present"), tr(c, "stopPresenting")}},
	}
	switch components.FloatingButtonGroup(c, tr(c, "boardActions"), group, ui.Horizontal, items).Clicked() {
	case "undo":
		st.edits, st.undone = st.edits-1, st.undone+1
	case "redo":
		st.edits, st.undone = st.edits+1, st.undone-1
	case "comment":
		st.edits, st.undone = st.edits+1, 0
	case "share":
		st.copiedAt = c.Now()
	case "present":
		st.presenting = !st.presenting
	}
}

// floatingToolsRail is a column of zoom and layer tools, as along a canvas
// edge.
func floatingToolsRail(c *ui.Context, st *floatingGroupState) {
	lock, hide := lucide.LockOpen, lucide.Eye
	if st.locked {
		lock = lucide.Lock
	}
	if st.hidden {
		hide = lucide.EyeOff
	}
	items := []components.FloatingButtonGroupItem{
		{ID: "in", Label: tr(c, "zoomIn"), Icon: lucide.ZoomIn, IconOnly: true, Shortcut: "⌘+", Disabled: st.zoom >= 3},
		{ID: "out", Label: tr(c, "zoomOut"), Icon: lucide.ZoomOut, IconOnly: true, Shortcut: "⌘−", Disabled: st.zoom <= -3},
		{ID: "fit", Label: tr(c, "zoomToFit"), Icon: lucide.Scan, IconOnly: true, Shortcut: "⇧1", Disabled: st.zoom == 0},
		{Separator: true},
		{ID: "lock", Label: tr(c, "lockLayer"), Icon: lock, IconOnly: true, Toggle: true, On: st.locked},
		{ID: "hide", Label: tr(c, "hideLayer"), Icon: hide, IconOnly: true, Toggle: true, On: st.hidden},
	}
	switch components.FloatingButtonGroup(c, tr(c, "canvasTools"), components.FloatingButtonGroupRaised, ui.Vertical, items).Clicked() {
	case "in":
		st.zoom++
	case "out":
		st.zoom--
	case "fit":
		st.zoom = 0
	case "lock":
		st.locked = !st.locked
	case "hide":
		st.hidden = !st.hidden
	}
}
