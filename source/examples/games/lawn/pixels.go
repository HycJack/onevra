package main

import (
	"image"
	"image/color"
)

// pixelSheet paints square frames of size pixels side by side; each byte
// of a row picks a palette color, and others are transparent.
func pixelSheet(size int, palette map[byte]color.NRGBA, frames ...[]string) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, size*len(frames), size))
	for f, rows := range frames {
		for y, row := range rows {
			for x := range size {
				img.SetNRGBA(f*size+x, y, palette[row[x]])
			}
		}
	}
	return img
}
