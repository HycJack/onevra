package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func (a *application) createTray() error {
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{Icon: trayIcon(), IconIsTemplate: true, ToolTip: a.text("title")})
	if err != nil {
		return err
	}
	if runtime.GOOS != "linux" {
		a.tray.OnClick(a.primaryClick)
		a.tray.OnRightClick(func() {
			a.hidePanel()
			a.tray.SetMenu(a.trayMenu)
			defer a.tray.SetMenu(nil)
			a.tray.PopUpMenu()
		})
	}
	return nil
}

// A product without a panel can replace the settings action here with its own
// quick action. No placeholder main window is required.
func (a *application) primaryClick() {
	if !a.config.Panel {
		a.openSettings()
		return
	}
	if a.panel != nil {
		a.hidePanel()
	} else {
		a.showPanel()
	}
}

func (a *application) hidePanel() {
	if a.panel != nil {
		a.panel.Close()
	}
}

// showPanel opens the quick-note panel as a system popover anchored to the
// tray icon. The platform owns placement, light dismissal and, on macOS, the
// popover material; panelView draws its own surface elsewhere.
func (a *application) showPanel() {
	if a.tray == nil || !a.config.Panel {
		a.openSettings()
		return
	}
	if a.panel == nil {
		panel := mygo.NewFlyout(mygo.FlyoutOptions{
			Tray: a.tray, Placement: mygo.PlacementBottom, Gap: 4,
			Width: 360, Height: 300, Focusable: true, Shadow: true, Popover: true,
			Content: ui.View(a.panelView),
		})
		a.panel = panel
		panel.OnClosed(func() {
			if a.panel == panel {
				a.panel = nil
			}
		})
	}
}

func (a *application) destroyTray() {
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
}

func trayIcon() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	ink := color.NRGBA{R: 79, G: 124, B: 255, A: 255}
	if runtime.GOOS == "darwin" {
		ink = color.NRGBA{A: 255}
	}
	for _, rect := range []image.Rectangle{image.Rect(6, 4, 26, 7), image.Rect(6, 25, 26, 28), image.Rect(6, 7, 9, 25), image.Rect(23, 7, 26, 25), image.Rect(12, 11, 20, 14), image.Rect(12, 18, 20, 21)} {
		draw.Draw(img, rect, image.NewUniform(ink), image.Point{}, draw.Src)
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		panic(err)
	}
	return b.Bytes()
}
