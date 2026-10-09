package components

import (
	"time"

	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// DialogOptions controls dismissal independently for each layer.
// The zero value allows backdrop clicks and Escape to dismiss the dialog.
type DialogOptions struct {
	DisableBackdropClose bool
	DisableEscapeClose   bool
	// Label is the accessible name. Applications supply localized text.
	Label string
	// Inset is the minimum viewport clearance in DIPs on every edge.
	// Non-positive values use 20. WindowDragHeight can increase this clearance.
	Inset float32
	// WindowDragHeight reserves a full-width native window drag strip at the
	// viewport top, inside this modal layer. Non-positive values disable it.
	// Set this and Inset to the host's header height to retain window dragging
	// without enabling underlying header actions. MyGo owns drag and focus behavior.
	WindowDragHeight float32
}

// Dialog builds arbitrary content in a modal panel and returns its frame-local
// handle (absent when closed). Applications own open and call Dialog every frame.
// Build layers in back-to-front order, with stable keyed parents for independent
// dialogs. Nested calls are supported. MyGo owns focus trapping and restoration.
// The panel is horizontally centered with spare vertical space biased below it.
// Its default maximum size leaves 20 logical pixels clear on every viewport edge;
// Inset overrides the clearance; WindowDragHeight can increase it.
// Content may customize the panel and should provide scrolling for long pages.
func Dialog(c *ui.Context, open *bool, options DialogOptions, build func(panel ui.Element)) ui.Element {
	var panel ui.Element
	if !*open {
		return panel
	}
	t := c.Theme()
	w, h := c.Size()
	// Distribute spare vertical space 1:2 above/below the panel.
	inset := options.Inset
	if inset <= 0 {
		inset = 20
	}
	dragHeight := min(max(0, options.WindowDragHeight), max(0, h))
	inset = max(inset, dragHeight)
	ui.Overlay(c, func() {
		back := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			AlignItems(ui.Center).Padding(inset).
			Modal().Background(ui.RGBA(0, 0, 0, .18))
		back.Transition(ui.ElementTransition{Duration: 180 * time.Millisecond, Colors: true, Enter: &ui.Motion{Opacity: 0}})
		clicked := back.Clicked()
		escape := back.OverlayShortcut(0, ui.KeyEscape)
		if clicked && !options.DisableBackdropClose || escape && !options.DisableEscapeClose {
			*open = false
		}
		back.Children(func() {
			ui.Box(c).Grow(1).MinHeight(0)
			panel = ui.Column(c).Role(ui.RoleDialog).Label(options.Label).
				Width(t.Space(110)).MaxWidth(max(0, w-2*inset)).MaxHeight(max(0, h-2*inset)).Shrink(0).
				Padding(t.Space(5)).Gap(t.Space(3)).Radius(t.Radius*2).
				Material(glass.Glass{Tint: t.Background}).TextColor(t.Text).
				Shadow(0, t.Space(2), t.Space(6), 0, ui.RGBA(0, 0, 0, .12))
			if t.Dark {
				panel.Border(1, t.Border)
			}
			// Register a hit target so blank panel space never dismisses the backdrop.
			panel.Clicked()
			panel.Children(func() {
				if build != nil {
					build(panel)
				}
			})
			ui.Box(c).Grow(2).MinHeight(0)
			if dragHeight > 0 {
				// A dedicated hit target, not a clickable backdrop or a forwarded
				// event: MyGo starts the native drag without dismissing this layer.
				ui.Box(c).Absolute().Left(0).Top(0).Right(0).Height(dragHeight).DragWindow()
			}
		})
	})
	return panel
}

// DialogResult is emitted only in the frame of a user action. Store it in
// application state if needed later; this UI-thread API never blocks.
type DialogResult uint8

const (
	DialogNone DialogResult = iota
	DialogConfirmed
	DialogCancelled
)

// AlertDialogOptions supplies localized text and optional cancellation.
// ConfirmLabel must be supplied. HideCancel produces a one-button notice;
// backdrop and Escape dismissal still return DialogCancelled unless disabled.
type AlertDialogOptions struct {
	DialogOptions
	Title, Message            string
	ConfirmLabel, CancelLabel string
	HideCancel                bool
}

// AlertDialog is a preset Dialog returning confirmation or cancellation.
// Programmatically closing open does not emit a result.
func AlertDialog(c *ui.Context, open *bool, options AlertDialogOptions) DialogResult {
	result := DialogNone
	wasOpen := *open
	opts := options.DialogOptions
	opts.Label = options.Title
	Dialog(c, open, opts, func(panel ui.Element) {
		panel.Role(ui.RoleAlertDialog)
		if options.Title != "" {
			Text(c, options.Title).Bold()
		}
		if options.Message != "" {
			Text(c, options.Message).TextColor(c.Theme().TextMuted)
		}
		ui.Row(c).FillWidth().Wrap().Justify(ui.End).Gap(c.Theme().Space(2)).Children(func() {
			if !options.HideCancel && Button(c, options.CancelLabel, ButtonSecondary).Clicked() {
				result = DialogCancelled
				*open = false
			}
			if Button(c, options.ConfirmLabel, ButtonPrimary).Clicked() {
				result = DialogConfirmed
				*open = false
			}
		})
	})
	if wasOpen && !*open && result == DialogNone {
		result = DialogCancelled
	}
	return result
}
