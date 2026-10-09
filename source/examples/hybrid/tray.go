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
	if a.panel != nil && a.panel.IsVisible() {
		a.hidePanel()
	} else {
		a.showPanel()
	}
}

func (a *application) hidePanel() {
	if a.panel != nil {
		a.panel.Hide()
	}
}

func (a *application) showPanel() {
	if a.tray == nil || !a.config.Panel {
		a.openSettings()
		return
	}
	if a.panel == nil {
		a.panel = mygo.NewWindow(mygo.WindowOptions{
			Title: a.text("panel_title"), Width: 360, Height: 300, Hidden: true,
			Frameless: true, AlwaysOnTop: true, SkipTaskbar: true,
			DisableResize: true, DisableMinimize: true, DisableMaximize: true, DisableFullScreen: true,
			Content: ui.View(a.panelView),
		})
		a.panel.OnBlur(a.hidePanel)
		a.panel.OnClose(func(e *mygo.CloseEvent) {
			if !a.quitting {
				e.PreventDefault()
				a.hidePanel()
			}
		})
	}
	anchor := a.tray.Bounds()
	if anchor.Width <= 0 || anchor.Height <= 0 {
		a.panel.Center()
	} else {
		area := mygo.Screen.DisplayNearestPoint(mygo.Point{X: anchor.X + anchor.Width/2, Y: anchor.Y + anchor.Height/2}).WorkArea
		x, y := panelPosition(anchor, area, a.panel.Bounds())
		a.panel.SetPosition(x, y)
	}
	show(a.panel)
}

func panelPosition(anchor, area, panel mygo.Rectangle) (int, int) {
	x := anchor.X + (anchor.Width-panel.Width)/2
	y := anchor.Y + anchor.Height + 6
	if y+panel.Height > area.Y+area.Height {
		y = anchor.Y - panel.Height - 6
	}
	return max(area.X, min(x, area.X+max(0, area.Width-panel.Width))), max(area.Y, min(y, area.Y+max(0, area.Height-panel.Height)))
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
