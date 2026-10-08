package components

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// DialogOptions controls dismissal independently for each layer.
// The zero value allows backdrop clicks and Escape to dismiss the dialog.
type DialogOptions struct {
	DisableBackdropClose bool
	DisableEscapeClose   bool
	// Label is the accessible name. Applications supply localized text.
	Label string
}

// Dialog builds arbitrary content in a modal panel and returns its frame-local
// handle (absent when closed). Applications own open and call Dialog every frame.
// Build layers in back-to-front order, with stable keyed parents for independent
// dialogs. Nested calls are supported. MyGo owns focus trapping and restoration.
// Content may customize the panel and should provide scrolling for long pages.
func Dialog(c *ui.Context, open *bool, options DialogOptions, build func(panel ui.Element)) ui.Element {
	var panel ui.Element
	if !*open {
		return panel
	}
	t := c.Theme()
	w, h := c.Size()
	// Bias the panel toward the upper third while retaining viewport clearance.
	bottomInset := max(t.Space(4), h*.34)
	ui.Overlay(c, func() {
		back := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			AlignItems(ui.Center).Justify(ui.Center).Padding(t.Space(4), t.Space(4), bottomInset, t.Space(4)).
			Modal().Background(ui.RGBA(0, 0, 0, .18))
		back.Transition(ui.ElementTransition{Duration: 180 * time.Millisecond, Colors: true, Enter: &ui.Motion{Opacity: 0}})
		clicked := back.Clicked()
		escape := back.OverlayShortcut(0, ui.KeyEscape)
		if clicked && !options.DisableBackdropClose || escape && !options.DisableEscapeClose {
			*open = false
		}
		back.Children(func() {
			panel = ui.Column(c).Role(ui.RoleDialog).Label(options.Label).
				Width(t.Space(110)).MaxWidth(max(0, w-t.Space(8))).MaxHeight(max(0, h-bottomInset-t.Space(4))).
				Padding(t.Space(5)).Gap(t.Space(3)).Radius(t.Radius*2).
				Background(t.Background).TextColor(t.Text).
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
