package chat

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// AttachmentPreviewOptions describes an application-prepared thumbnail.
// Label names the image for accessibility and its tooltip. An empty
// RemoveLabel makes it read-only. Loading and decoding belong to the host.
type AttachmentPreviewOptions struct {
	Label, RemoveLabel string
	Image              ui.ImageSource
}

// AttachmentPreviewResult reports a frame-local removal request.
type AttachmentPreviewResult struct {
	Root   ui.Element
	Remove bool
}

// AttachmentPreview shows a compact image tile with an optional remove
// action. A nil image shows a placeholder, without reading any file.
func AttachmentPreview(c *ui.Context, options AttachmentPreviewOptions) AttachmentPreviewResult {
	var result AttachmentPreviewResult
	t := c.Theme()
	result.Root = ui.Box(c).Size(64, 64).Shrink(0).Radius(8).Clip().
		Border(1, strongBorder(t)).Background(t.Surface).Children(func() {
		if options.Image != nil {
			ui.Image(c, options.Image).Label(options.Label).Size(62, 62).Fit(ui.Cover).Shrink(0)
		} else {
			ui.Box(c).Label(options.Label).Size(62, 62).Center().Children(func() {
				components.Icon(c, lucide.Image).Size(20, 20).TextColor(t.TextMuted)
			})
		}
		if options.RemoveLabel != "" {
			b := ui.ButtonBase(c).Label(options.RemoveLabel).Absolute().Top(3).Right(3).
				Size(18, 18).Radius(5).Center().Background(t.Background.Alpha(.9)).TextColor(t.Text)
			if b.Hovered() || b.Focused() {
				b.Background(t.Surface).Border(1, t.Accent)
			}
			b.Children(func() { components.Icon(c, lucide.X).Size(11, 11) })
			result.Remove = b.Clicked()
		}
	})
	components.Tooltip(c, result.Root, options.Label)
	return result
}
