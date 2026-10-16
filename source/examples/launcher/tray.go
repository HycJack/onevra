package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"runtime"

	"github.com/egoist/mygo"
)

func (a *application) createTray() error {
	if a.tray != nil {
		return nil
	}
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: trayIcon(), IconIsTemplate: true, ToolTip: a.text("title")})
	if err != nil {
		return err
	}
	if runtime.GOOS != "linux" {
		a.tray.OnClick(a.toggle)
		a.tray.OnRightClick(func() {
			a.hide()
			a.tray.SetMenu(a.trayMenu)
			defer a.tray.SetMenu(nil)
			a.tray.PopUpMenu()
		})
	}
	return nil
}

func (a *application) destroyTray() {
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
}

// trayIcon draws a magnifying glass, black on transparent for macOS to tint,
// blue elsewhere to show on light and dark taskbars alike.
func trayIcon() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	ink := color.NRGBA{R: 79, G: 124, B: 255, A: 255}
	if runtime.GOOS == "darwin" {
		ink = color.NRGBA{A: 255}
	}
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			dx, dy := float64(x)-13, float64(y)-13
			if r := dx*dx + dy*dy; r >= 7*7 && r <= 10*10 {
				img.Set(x, y, ink)
			}
		}
	}
	for i := 0; i < 8; i++ { // the handle
		draw.Draw(img, image.Rect(20+i, 20+i, 23+i, 23+i), image.NewUniform(ink), image.Point{}, draw.Src)
	}
	var b bytes.Buffer
	_ = png.Encode(&b, img)
	return b.Bytes()
}
