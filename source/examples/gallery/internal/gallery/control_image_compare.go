package gallery

import (
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// controlImageCompare compares the avatar in grayscale with its original
// colors, so both sides stay aligned without a second asset.
func controlImageCompare(c *ui.Context, s *controlState, variant string) {
	bitmap, err := galleryImage()
	if err != nil {
		components.Text(c, tr(c, "unableToLoadAvatarJpeg")).TextColor(c.Theme().TextMuted)
		return
	}
	options := components.ImageCompareOptions{BeforeLabel: tr(c, "imageCompareBefore"), AfterLabel: tr(c, "imageCompareAfter")}
	if variant == "Vertical" {
		options.Orientation = components.Vertical
	}
	side := func(gray bool) func() {
		return func() {
			image := components.Image(c, bitmap).Fill().Fit(ui.Cover)
			if gray {
				image.Grayscale()
			}
		}
	}
	components.ImageCompare(c, &s.compare, side(true), side(false), options).
		Label(msg(c).ImageCompareOf(tr(c, variant))).FillWidth().MaxWidth(720)
	components.Text(c, msg(c).ImageCompareShare(int(math.Round(float64(s.compare))))).TextColor(c.Theme().TextMuted)
}
