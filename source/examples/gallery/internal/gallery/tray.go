package gallery

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/theme"
)

const trayPage galleryPage = componentStatusPage + 1

// The application owns the panel, in-memory note and lifetime. All callbacks
// run on MyGo's UI thread; hiding the panel does not discard the note.
type trayDemo struct {
	tray    *mygo.Tray
	win     *mygo.Window
	note    string
	owner   *mygo.Window
	offQuit func()
	err     error
}

func (a *trayDemo) view(c *ui.Context) {
	theme.Apply(c, theme.Auto)
	if c.Shortcut(0, ui.KeyEscape) || c.Shortcut(ui.Cmd, ui.KeyW) {
		a.win.Hide()
	}
	c.OnShortcut(ui.Cmd, ui.KeyQ, mygo.App.Quit)
	ui.Column(c).Fill().Padding(20).Gap(12).Children(func() {
		ui.Row(c).Gap(12).Children(func() {
			components.Text(c, "Tray notes").FontSize(22).Bold().Grow(1)
			components.Button(c, "Hide", components.ButtonGhost).OnClick(a.win.Hide)
		})
		components.Text(c, "A small window, opened from the tray.").TextColor(c.Theme().TextMuted)
		components.TextBox(c, &a.note, components.TextBoxOptions{Multiline: true}).
			Grow(1).Label("Note").Placeholder("Write a quick note…").AutoFocus()
		ui.Row(c).Gap(8).Children(func() {
			components.Button(c, "Clear", components.ButtonOutline).
				Disabled(a.note == "").OnClick(func() { a.note = "" })
			ui.Spacer(c)
			components.Button(c, "Quit", components.ButtonSecondary).OnClick(mygo.App.Quit)
		})
		components.Text(c, "Kept in memory until you quit. Escape hides this window.").
			FontSize(12).TextColor(c.Theme().TextMuted)
	})
}

// panelPosition prefers below the icon, flips above a bottom taskbar, and
// clamps to the display work area (which can have negative coordinates).
func panelPosition(anchor, area, panel mygo.Rectangle) (x, y int) {
	x = anchor.X + (anchor.Width-panel.Width)/2
	y = anchor.Y + anchor.Height + 6
	if y+panel.Height > area.Y+area.Height {
		y = anchor.Y - panel.Height - 6
	}
	x = max(area.X, min(x, area.X+max(0, area.Width-panel.Width)))
	y = max(area.Y, min(y, area.Y+max(0, area.Height-panel.Height)))
	return
}

func (a *trayDemo) show() {
	anchor := a.tray.Bounds()
	if anchor.Width <= 0 || anchor.Height <= 0 {
		// AppIndicator and icons in the Windows overflow area may have no bounds.
		a.win.Center()
	} else {
		area := mygo.Screen.DisplayNearestPoint(mygo.Point{
			X: anchor.X + anchor.Width/2, Y: anchor.Y + anchor.Height/2,
		}).WorkArea
		x, y := panelPosition(anchor, area, a.win.Bounds())
		a.win.SetPosition(x, y)
	}
	a.win.Show()
}

func (a *trayDemo) toggle() {
	if a.win.IsVisible() {
		a.win.Hide()
	} else {
		a.show()
	}
}

func (a *trayDemo) start(owner *mygo.Window) error {
	if a.tray != nil {
		return nil
	}
	a.owner = owner
	a.err = nil
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{
		Icon: trayIcon(), IconIsTemplate: true, ToolTip: "Onevra Gallery — Tray demo",
	})
	if err != nil {
		a.err = err
		return err
	}
	a.win = mygo.NewWindow(mygo.WindowOptions{
		Title: "Gallery Tray notes", Width: 380, Height: 320,
		Hidden: true, Frameless: true, AlwaysOnTop: true, SkipTaskbar: true,
		DisableResize: true, DisableMinimize: true,
		DisableMaximize: true, DisableFullScreen: true,
		Content: ui.View(a.view),
	})
	a.win.OnBlur(a.win.Hide)
	a.win.OnClose(func(e *mygo.CloseEvent) {
		// Application quit destroys this auxiliary window before closing windows.
		e.PreventDefault()
		a.win.Hide()
	})
	menu := mygo.NewMenu([]*mygo.MenuItem{
		{Label: "Open notes", Click: func(*mygo.MenuItem, *mygo.Window) { a.show() }},
		{Label: "Show Gallery", Click: func(*mygo.MenuItem, *mygo.Window) {
			a.win.Hide()
			a.owner.Show()
		}},
		mygo.Separator(),
		{Role: mygo.RoleQuit},
	})
	if runtime.GOOS == "linux" {
		// AppIndicator reports no click events; a menu action opens the panel.
		a.tray.SetMenu(menu)
	} else {
		a.tray.OnClick(a.toggle)
		a.tray.OnRightClick(func() {
			a.win.Hide()
			// A permanent menu would consume primary clicks on macOS.
			a.tray.SetMenu(menu)
			defer a.tray.SetMenu(nil)
			a.tray.PopUpMenu()
		})
	}
	a.offQuit = mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.close() })
	return nil
}

// A self-contained 32px PNG avoids resource-path and packaging requirements.
// macOS tints the template; blue remains visible on other platforms' taskbars.
func trayIcon() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	ink := color.NRGBA{R: 79, G: 124, B: 255, A: 255}
	if runtime.GOOS == "darwin" {
		ink = color.NRGBA{A: 255}
	}
	for _, r := range []image.Rectangle{
		image.Rect(6, 4, 26, 7), image.Rect(6, 25, 26, 28),
		image.Rect(6, 7, 9, 25), image.Rect(23, 7, 26, 25),
		image.Rect(12, 11, 20, 14), image.Rect(12, 18, 20, 21),
	} {
		draw.Draw(img, r, image.NewUniform(ink), image.Point{}, draw.Src)
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		panic(err) // bytes.Buffer cannot fail to write.
	}
	return b.Bytes()
}

func (a *trayDemo) close() {
	if a.offQuit != nil {
		a.offQuit()
		a.offQuit = nil
	}
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
	if a.win != nil {
		a.win.Destroy()
		a.win = nil
	}
	a.owner = nil
}

// Tray creates native resources only in response to an explicit action. Merely
// navigating here, including in a headless test, never creates a status item.
func (g *Gallery) Tray(c *ui.Context) {
	d := &g.tray
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(c.Theme().Space(6)).Gap(16).Children(func() {
			components.Text(c, "Tray").FontSize(22).Bold()
			instructions := "Left-click the tray icon to toggle the notes window. Right-click for its native menu."
			if runtime.GOOS == "linux" {
				instructions = "AppIndicator supports menus only. Choose Open notes from the tray menu."
			}
			components.Text(c, instructions).Wrap()
			components.Text(c, "Escape, closing the notes window or clicking away hides it. Notes stay in memory; closing Gallery removes the tray.").Wrap().TextColor(c.Theme().TextMuted)
			ui.Row(c).Gap(8).Wrap().Children(func() {
				if d.tray == nil {
					if components.Button(c, "Enable tray", components.ButtonPrimary).Disabled(g.Shell == nil).Clicked() {
						_ = d.start(g.Shell.Window())
					}
				} else {
					components.Button(c, "Open notes", components.ButtonPrimary).OnClick(d.show)
					components.Button(c, "Disable tray", components.ButtonOutline).OnClick(d.close)
				}
			})
			if d.err != nil {
				components.Text(c, "Tray unavailable: "+d.err.Error()).Wrap().TextColor(c.Theme().Danger)
			}
		})
	})
}
