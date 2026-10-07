package gallery

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/components"
)

//go:embed hybrid.html
var hybridHTML string

const hybridPage galleryPage = -1

// hybridDemo owns only demonstration state. The reusable component owns the
// native attachment and the hidden browser's lifetime.
type hybridDemo struct {
	web                          *components.WebView
	host                         *mygo.Window
	visible, ready, busy, closed bool
	count                        int
	status, input                string
	alertOpen, alertHideWeb      bool
	alertDeadline                time.Time
}

func (g *Gallery) setHybridActive(active bool) {
	if g.hybridActive == active {
		return
	}
	g.hybridActive = active
	if g.hybrid != nil {
		if !active {
			g.hybrid.alertOpen = false
		}
		g.hybrid.web.SetVisible(active && g.hybrid.webVisible())
	}
}

func (g *Gallery) Hybrid(c *ui.Context) {
	if g.hybrid != nil {
		g.hybrid.build(c)
		return
	}
	ui.Column(c).Fill().Padding(20).Gap(12).Children(func() {
		ui.Text(c, "Native + WebView").Bold()
		if g.hybridErr != nil {
			ui.Text(c, g.hybridErr.Error()).TextColor(c.Theme().TextMuted).Wrap()
			return
		}
		if g.Shell == nil {
			ui.Text(c, "Open the gallery in a native window to load the WebView.").Wrap()
			return
		}
		ui.Text(c, "Loading WebView…")
	})
	if g.closed || g.Shell == nil || g.hybridErr != nil || g.hybridStarting {
		return
	}
	g.hybridStarting = true
	win := g.Shell.Window()
	// Create native resources only once and outside the build transaction.
	win.Update(func() {
		defer func() { g.hybridStarting = false }()
		if g.closed || !g.hybridActive || g.page != hybridPage {
			return
		}
		web, err := components.NewWebView(win, mygo.PageOptions{})
		if err != nil {
			g.hybridErr = err
			return
		}
		d := &hybridDemo{web: web, host: win, visible: true, status: "Loading local HTML…", input: "Native text input"}
		g.hybrid = d
		web.Page().OnDidFinishLoad(func() {
			if !d.closed {
				d.ready = true
				d.status = "Ready · real WKWebView, not a screenshot"
				win.Invalidate()
			}
		})
		web.Page().OnDidFailLoad(func(e *mygo.LoadError) {
			if !d.closed {
				d.status = fmt.Sprintf("Web load failed: %v", e)
				win.Invalidate()
			}
		})
		// Only local fixtures; no remote resources or application services.
		web.Page().LoadHTML(hybridHTML, "")
	})
}

func (d *hybridDemo) build(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().MinWidth(0).Padding(16).Gap(10).Children(func() {
		ui.Row(c).FillWidth().Shrink(0).Wrap().Gap(8).Children(func() {
			if ui.Button(c, "Send count to web").Disabled(!d.ready || d.busy).Clicked() {
				d.sendCount()
			}
			if ui.Button(c, "Show / hide web").Clicked() {
				d.visible = !d.visible
			}
			if ui.Button(c, "Toggle appearance").Clicked() {
				next := mygo.ThemeDark
				if t.Dark {
					next = mygo.ThemeLight
				}
				mygo.Theme.SetSource(next)
			}
		})
		ui.Row(c).FillWidth().Shrink(0).Wrap().Gap(8).Children(func() {
			if ui.Button(c, "Alert over WebView (8 s)").Clicked() {
				d.openAlert(false)
			}
			if ui.Button(c, "Alert with WebView hidden").Clicked() {
				d.openAlert(true)
			}
		})
		ui.TextInput(c, &d.input).Label("Native text input").FillWidth().Shrink(0)
		ui.Text(c, d.status).TextColor(t.TextMuted).Wrap().Shrink(0)
		d.web.View(c).Fill().MinHeight(0).Background(t.Surface).Children(func() {
			if !d.visible {
				ui.Text(c, "WebView hidden").Padding(16).Wrap()
			}
		})
		ui.Text(c, "Experimental · macOS ARM64 · system-view composition, not a shared GPU scene.").TextColor(t.TextMuted).Wrap().Shrink(0)
	})
	d.buildAlert(c)
	d.web.SetVisible(d.webVisible())
}

func (d *hybridDemo) webVisible() bool {
	return d.visible && !(d.alertOpen && d.alertHideWeb)
}

func (d *hybridDemo) openAlert(hideWeb bool) {
	d.alertOpen, d.alertHideWeb = true, hideWeb
	d.alertDeadline = time.Now().Add(8 * time.Second)
	d.status = "Occlusion probe · WebView stays visible · auto-closes in 8 seconds"
	if hideWeb {
		d.status = "Alert comparison · WebView explicitly hidden until dismissal"
	}
}

func (d *hybridDemo) buildAlert(c *ui.Context) {
	if !d.alertOpen {
		return
	}
	message := "The WebView is temporarily hidden. This does not prove shared composition."
	if !d.alertHideWeb {
		message = "The WebView stays visible and may cover this MyGo alert. Auto-closes after 8 seconds; Escape also cancels while native UI has focus."
		if !time.Now().Before(d.alertDeadline) {
			d.alertOpen = false
			d.status = "Overlay test finished · WebView was kept visible"
			return
		}
		c.After(100 * time.Millisecond)
	}
	if choice := ui.AlertDialog(c, &d.alertOpen, "Native alert / WebView test", message, "Cancel", "OK"); choice >= 0 {
		d.status = "Native alert closed"
	}
}

func (d *hybridDemo) sendCount() {
	d.count++
	count := d.count
	d.busy = true
	// Eval waits for WebKit, so do not block MyGo's view builder.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := d.web.Page().EvalContext(ctx, fmt.Sprintf("return window.receiveNative(%d)", count))
		d.host.Update(func() {
			if d.closed {
				return
			}
			d.busy = false
			d.status = fmt.Sprintf("Web acknowledged native count %d", count)
			if err != nil {
				d.status = "JavaScript: " + err.Error()
			}
		})
	}()
}
