package gallery

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// webViewDemo owns one web view of the gallery window. It is created once on
// first entry and kept while other pages show; the view only places it.
// loaded is the theme and language its page was last rendered for. covered
// tells that a modal layer of the gallery is open over the window.
type webViewDemo struct {
	view    *mygo.WebView
	err     error
	loaded  string
	menu    bool
	covered bool
}

func (d *webViewDemo) start(create func() (*mygo.WebView, error)) {
	if d.view != nil || d.err != nil {
		return
	}
	d.view, d.err = create()
}

func (d *webViewDemo) close() {
	if d.view != nil {
		d.view.Destroy()
	}
	d.view, d.err, d.loaded = nil, nil, ""
}

func (g *Gallery) newWebView() (*mygo.WebView, error) {
	return g.Shell.Window().NewWebView(mygo.WebViewOptions{})
}

var webViewPage = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="{{.Lang}}"><head><meta charset="utf-8"><style>
body { margin: 0; padding: 24px; font: 14px/1.5 -apple-system, "Segoe UI", system-ui, sans-serif;
	background: {{.Background}}; color: {{.Text}}; }
h1 { margin: 0 0 8px; font-size: 20px; }
p { margin: 0 0 16px; color: {{.Muted}}; }
.row { display: flex; gap: 8px; align-items: center; }
input, button { font: inherit; padding: 6px 10px; border-radius: 6px; border: 1px solid {{.Border}};
	background: transparent; color: inherit; }
button { background: {{.Accent}}; color: {{.AccentText}}; border-color: transparent; }
</style></head><body>
<h1>{{.Title}}</h1>
<p>{{.Body}}</p>
<div class="row">
	<input placeholder="{{.Placeholder}}" aria-label="{{.Placeholder}}">
	<button onclick="count.value = Number(count.value) + 1">{{.Count}}</button>
	<output id="count">0</output>
</div>
</body></html>`))

func cssColor(c ui.Color) template.CSS {
	return template.CSS(fmt.Sprintf("rgba(%d, %d, %d, %.3f)", c.R, c.G, c.B, float32(c.A)/255))
}

// webViewHTML renders the demonstration page in the gallery's current theme
// and language, so the page reads as part of the native content around it.
func webViewHTML(c *ui.Context) string {
	t := c.Theme()
	lang := "en"
	if l := *languageSlot(c); l != nil {
		lang = l.code
	}
	var b bytes.Buffer
	_ = webViewPage.Execute(&b, map[string]any{
		"Lang":       lang,
		"Background": cssColor(t.Background), "Text": cssColor(t.Text), "Muted": cssColor(t.TextMuted),
		"Border": cssColor(t.Border), "Accent": cssColor(t.Accent), "AccentText": cssColor(t.AccentText),
		"Title":       tr(c, "aWebPageInsideNativeUi"),
		"Body":        tr(c, "theSystemWebEngineRendersThisPageNativePopoversDialogsAnd"),
		"Placeholder": tr(c, "typeIntoThePage"),
		"Count":       tr(c, "count"),
	})
	return b.String()
}

func exampleWebView(c *ui.Context, s *demoState) {
	t := c.Theme()
	d := &s.webView
	if d.err != nil {
		components.Text(c, tr(c, "webViewUnavailable")).Bold()
		components.Text(c, d.err.Error()).TextColor(t.Danger).Wrap()
		return
	}
	if d.view == nil {
		components.Text(c, tr(c, "webViewsNeedANativeWindow")).TextColor(t.TextMuted).Wrap()
		return
	}
	if html := webViewHTML(c); d.loaded != html {
		d.loaded = html
		// Fill with the new theme until the reloaded page paints, not white.
		if err := d.view.SetBackgroundColor(string(cssColor(t.Background))); err != nil {
			log.Printf("Set web view background: %v", err)
		}
		d.view.Page().LoadHTML(html, "")
	}
	ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		if components.Button(c, tr(c, "reloadPage"), components.ButtonSecondary).Clicked() {
			d.view.Page().LoadHTML(d.loaded, "")
		}
		anchor := components.Button(c, tr(c, "overThePage"), components.ButtonOutline)
		if anchor.Clicked() {
			d.menu = !d.menu
		}
		components.Popover(c, anchor, &d.menu, func(panel ui.Element) {
			components.Text(c, tr(c, "nativeUiPaintedAfterTheWebViewShowsOverItsPage")).Wrap().MaxWidth(t.Space(60))
		})
	})
	var frame ui.Element
	if d.covered {
		// Glass materials cannot see a web page under them, so a modal layer
		// over the window would show a hole. Hide the web view, keeping its
		// page, behind a box of the page's background meanwhile.
		frame = ui.Box(c).Background(t.Background)
	} else {
		frame = ui.WebView(c, d.view)
	}
	frame.FillWidth().Height(320).Radius(t.Radius*1.5).Border(1, t.Border).Children(func() {
		components.Badge(c, "HTML").Absolute().Top(t.Space(3)).Right(t.Space(3))
	})
}
