package gallery

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

// traySVGData is the Onevra application icon's artwork without its tile.
//
//go:embed tray.svg
var traySVGData []byte

var traySVG = ui.MustParseSVG(traySVGData)

// The application owns the tray, its panel and their lifetime. All callbacks
// run on MyGo's UI thread.
type trayDemo struct {
	localize  func(*ui.Context, func())
	translate func(string) string
	showOwner func()
	menu      *mygo.Menu
	tray      *mygo.Tray
	panel     *mygo.Window
	offQuit   func()
}

func (a *trayDemo) view(c *ui.Context) {
	if a.localize != nil {
		a.localize(c, func() { a.localizedView(c) })
		return
	}
	a.localizedView(c)
}

func (a *trayDemo) localizedView(c *ui.Context) {
	theme.Apply(c)
	if c.Shortcut(0, ui.KeyEscape) || c.Shortcut(ui.Cmd, ui.KeyW) {
		a.hide()
	}
	// The popover material shows on macOS; elsewhere the flyout is
	// transparent and the panel draws its own surface.
	c.Root().Background(ui.Transparent)
	panel := ui.Column(c).Fill().Padding(16).Gap(12)
	if !c.Vibrancy() {
		panel.Background(c.Theme().Surface)
		switch runtime.GOOS {
		case "darwin":
			panel.Radius(10)
		case "windows":
		default:
			panel.Radius(10).Border(1, c.Theme().Border)
		}
	}
	panel.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(12).Children(func() {
			ui.Image(c, traySVG).Size(32, 32)
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				components.Text(c, tr(c, "onevraGallery")).Bold()
				components.Text(c, tr(c, "trayHint")).Wrap().FontSize(12).TextColor(c.Theme().TextMuted)
			})
		})
		components.Button(c, tr(c, "showGallery"), components.ButtonPrimary).FillWidth().OnClick(a.showGallery)
	})
}

// show opens the panel as a system popover anchored to the tray icon. The
// platform owns placement and light dismissal.
func (a *trayDemo) show() {
	if a.tray == nil || a.panel != nil {
		return
	}
	panel := mygo.NewFlyout(mygo.FlyoutOptions{
		Tray: a.tray, Placement: mygo.PlacementBottom, Gap: 4,
		Width: 300, Height: 124, Focusable: true, Shadow: true, Popover: true,
		Content: ui.View(a.view),
	})
	a.panel = panel
	panel.OnClosed(func() {
		if a.panel == panel {
			a.panel = nil
		}
	})
}

func (a *trayDemo) hide() {
	if a.panel != nil {
		a.panel.Close()
	}
}

func (a *trayDemo) toggle() {
	if a.panel != nil {
		a.hide()
	} else {
		a.show()
	}
}

func (a *trayDemo) showGallery() {
	a.hide()
	if a.showOwner != nil {
		a.showOwner()
	}
}

// start creates the tray icon; showOwner brings the main Gallery window back.
func (a *trayDemo) start(showOwner func()) error {
	if a.tray != nil {
		return nil
	}
	a.showOwner = showOwner
	template := runtime.GOOS == "darwin"
	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{
		Icon: trayIcon(template), IconIsTemplate: template, ToolTip: a.text("trayTooltip"),
	})
	if err != nil {
		return err
	}
	a.menu = a.buildMenu()
	if runtime.GOOS == "linux" {
		// AppIndicator reports no click events; a menu action opens the panel.
		a.tray.SetMenu(a.menu)
	} else {
		a.tray.OnClick(a.toggle)
		a.tray.OnRightClick(func() {
			a.hide()
			// A permanent menu would consume primary clicks on macOS.
			a.tray.SetMenu(a.menu)
			defer a.tray.SetMenu(nil)
			a.tray.PopUpMenu()
		})
	}
	a.offQuit = mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.close() })
	return nil
}

// trayIcon renders the rounded rectangles of tray.svg into a 32px PNG, in
// black as a macOS template or in the icon's own colors for other taskbars.
func trayIcon(template bool) []byte {
	var doc struct {
		ViewBox string `xml:"viewBox,attr"`
		Rects   []struct {
			X      string `xml:"x,attr"`
			Y      string `xml:"y,attr"`
			Width  string `xml:"width,attr"`
			Height string `xml:"height,attr"`
			RX     string `xml:"rx,attr"`
			Fill   string `xml:"fill,attr"`
		} `xml:"rect"`
	}
	if err := xml.Unmarshal(traySVGData, &doc); err != nil {
		panic(err)
	}
	var box [4]float32
	parseFloats(doc.ViewBox, box[:])
	const size = 32
	scale := size / box[2]
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for _, r := range doc.Rects {
		var v [5]float32
		parseFloats(r.X+" "+r.Y+" "+r.Width+" "+r.Height+" "+r.RX, v[:])
		x, y := (v[0]-box[0])*scale, (v[1]-box[1])*scale
		w, h, rad := v[2]*scale, v[3]*scale, v[4]*scale
		ink := color.NRGBA{A: 255}
		if !template {
			n, _ := strconv.ParseUint(r.Fill[1:], 16, 32)
			ink = color.NRGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}
		}
		fillRoundedRect(img, x, y, w, h, rad, ink)
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		panic(err) // bytes.Buffer cannot fail to write.
	}
	return b.Bytes()
}

// fillRoundedRect blends a rounded rectangle over img, with 4×4 samples per
// pixel for antialiased edges.
func fillRoundedRect(img *image.NRGBA, x, y, w, h, r float32, ink color.NRGBA) {
	const n = 4
	b := img.Bounds()
	for py := b.Min.Y; py < b.Max.Y; py++ {
		for px := b.Min.X; px < b.Max.X; px++ {
			hits := 0
			for sy := range n {
				for sx := range n {
					// Distance from the sample to the rectangle shrunk by r.
					dx := max(x+r-(float32(px)+(float32(sx)+0.5)/n), (float32(px)+(float32(sx)+0.5)/n)-(x+w-r), 0)
					dy := max(y+r-(float32(py)+(float32(sy)+0.5)/n), (float32(py)+(float32(sy)+0.5)/n)-(y+h-r), 0)
					if dx*dx+dy*dy <= r*r {
						hits++
					}
				}
			}
			if hits > 0 {
				c := ink
				c.A = uint8(math.Round(float64(ink.A) * float64(hits) / (n * n)))
				img.Set(px, py, blend(img.NRGBAAt(px, py), c))
			}
		}
	}
}

// blend composites src over dst in non-premultiplied color.
func blend(dst, src color.NRGBA) color.NRGBA {
	sa, da := float64(src.A)/255, float64(dst.A)/255
	a := sa + da*(1-sa)
	if a == 0 {
		return color.NRGBA{}
	}
	mix := func(s, d uint8) uint8 {
		return uint8(math.Round((float64(s)*sa + float64(d)*da*(1-sa)) / a))
	}
	return color.NRGBA{mix(src.R, dst.R), mix(src.G, dst.G), mix(src.B, dst.B), uint8(math.Round(a * 255))}
}

func parseFloats(s string, dst []float32) {
	i := 0
	for _, f := range bytes.Fields([]byte(s)) {
		if i < len(dst) {
			v, _ := strconv.ParseFloat(string(f), 32)
			dst[i] = float32(v)
			i++
		}
	}
}

func (a *trayDemo) close() {
	if a.offQuit != nil {
		a.offQuit()
		a.offQuit = nil
	}
	if a.panel != nil {
		a.panel.Destroy()
		a.panel = nil
	}
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
	a.showOwner = nil
	a.menu = nil
}

// EnableDesktopIntegrations installs the tray icon and, on macOS, the Dock
// menu. The standalone Gallery calls it at startup; disk-free galleries and
// headless tests never create these native resources.
func (g *Gallery) EnableDesktopIntegrations() {
	if g.Shell == nil {
		return
	}
	if err := g.tray.start(g.showMainWindow); err != nil {
		log.Printf("Enable gallery tray: %v", err)
	}
	g.enableDockMenu()
	g.updateLanguageWindows()
}

func (a *trayDemo) text(key string) string {
	if a.translate != nil {
		return a.translate(key)
	}
	return (*translator)(nil).text(key)
}

func (a *trayDemo) buildMenu() *mygo.Menu {
	return mygo.NewMenu([]*mygo.MenuItem{
		{ID: "panel", Label: a.text("openTrayPanel"), Click: func(*mygo.MenuItem, *mygo.Window) { a.show() }},
		{ID: "gallery", Label: a.text("showGallery"), Click: func(*mygo.MenuItem, *mygo.Window) { a.showGallery() }},
		mygo.Separator(),
		{ID: "quit", Role: mygo.RoleQuit, Label: a.text("quit")},
	})
}
