//go:build darwin && arm64

package components

import (
	"os"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestNativeHybridComposition(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires macOS desktop")
	}
	var host *mygo.Window
	var web *WebView
	var err error
	ready, overlay, clicked, modal, omit, clipped := false, false, false, false, false, false
	mygo.RunOnMain(func() {
		host = mygo.NewWindow(mygo.WindowOptions{Title: "GPU hybrid composition", Width: 700, Height: 600, Hidden: true, Content: ui.View(func(c *ui.Context) {
			ui.Column(c).Fill().Padding(20).Children(func() {
				if web != nil && !omit {
					if clipped {
						ui.Box(c).Width(200).Height(120).Clip().Children(func() { web.View(c).Width(400).Height(300).Shrink(0) })
					} else {
						web.View(c).Fill()
					}
				}
			})
			if overlay && web != nil {
				ui.Overlay(c, func() {
					if ui.Button(c, "Native overlay").Absolute().Left(web.rect.X).Top(web.rect.Y).Width(200).Height(100).Clicked() {
						clicked = true
					}
				})
			}
			if modal {
				ui.AlertDialog(c, &modal, "Native modal", "Live browser is blocked", "Cancel", "OK")
			}
		})})
		web, err = NewWebView(host, mygo.PageOptions{})
		if err != nil {
			return
		}
		web.Page().OnDidFinishLoad(func() { ready = true })
		web.Page().LoadHTML(webViewFixture+`<script>window.ticks=0;setInterval(()=>window.ticks++,30)</script>`, "")
		host.Show()
	})
	defer host.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	waitWebView(t, "browser ready in key window", func() bool {
		return ready && web.placed && !web.pending && send(web.view.(*macAttachment).host, "isKeyWindow") != 0
	})
	clickWebButton(t, web)
	mygo.RunOnMain(func() { overlay = true; host.Invalidate() })
	waitWebView(t, "overlay hit barrier", func() bool { return !web.route(float64(web.rect.X+40), float64(web.rect.Y+40), false) })
	clickWebButton(t, web)
	waitWebView(t, "native overlay receives click", func() bool { return clicked })
	if got := evalWebView(t, web, "return window.count"); got != float64(1) {
		t.Fatalf("overlay click reached browser: %v", got)
	}
	mygo.RunOnMain(func() { overlay = false; host.Invalidate() })
	waitWebView(t, "web input restored", func() bool { return web.route(float64(web.rect.X+10), float64(web.rect.Y+10), false) })
	typeInWeb(t, web)
	before := evalWebView(t, web, "return window.ticks").(float64)
	mygo.RunOnMain(func() { modal = true; host.Invalidate() })
	waitWebView(t, "modal takes native focus", func() bool {
		a := web.view.(*macAttachment)
		return !web.pending && send(a.host, "firstResponder") == a.surface && send(a.root, "isHidden") == 0 && !web.route(float64(web.rect.X+10), float64(web.rect.Y+10), false)
	})
	time.Sleep(120 * time.Millisecond)
	if evalWebView(t, web, "return window.ticks").(float64) <= before {
		t.Fatal("browser stopped under modal")
	}
	mygo.RunOnMain(func() { hybridEscape(web.view.(*macAttachment).host) })
	waitWebView(t, "Escape closes modal", func() bool { return !modal })
	mygo.RunOnMain(func() { omit = true; host.Invalidate() })
	waitWebView(t, "omission hides browser", func() bool { return send(web.view.(*macAttachment).root, "isHidden") != 0 })
	mygo.RunOnMain(func() { omit = false; clipped = true; host.Invalidate() })
	waitWebView(t, "ancestor clipping", func() bool { return !web.pending && web.layout.Clip.W == 200 && web.layout.Clip.H == 120 })
	mygo.RunOnMain(func() {
		a := web.view.(*macAttachment)
		f := nativeRect(a.root, objc.RegisterName("frame"))
		if f.Size.W != 200 || f.Size.H != 120 {
			t.Errorf("native clip frame: %+v", f)
		}
		if web.route(float64(web.rect.X+300), float64(web.rect.Y+10), false) {
			t.Error("input escaped clip")
		}
	})
	if got := evalWebView(t, web, "return [innerWidth,innerHeight]"); len(got.([]any)) != 2 || got.([]any)[0] != float64(400) || got.([]any)[1] != float64(300) {
		t.Fatalf("clipping resized document: %v", got)
	}
	mygo.RunOnMain(func() { web.Close(); web.Close() })
	if host.IsDestroyed() || !web.donor.IsDestroyed() {
		t.Fatal("close ownership")
	}
}

func hybridEscape(host objc.ID) {
	chars := send(objc.ID(objc.GetClass("NSString")), "stringWithUTF8String:", "\u001b")
	for _, kind := range []uint64{10, 11} {
		e := send(objc.ID(objc.GetClass("NSEvent")), "keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:", kind, point{}, uint64(0), float64(0), send(host, "windowNumber"), uintptr(0), chars, chars, uintptr(53<<16))
		send(host, "sendEvent:", e)
	}
}
