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

func inspectorBool(t *testing.T, view objc.ID, script string) bool {
	t.Helper()
	type result struct{ value, failed bool }
	done := make(chan result, 1)
	mygo.RunOnMain(func() {
		block := objc.NewBlock(func(_ objc.Block, value, err objc.ID) {
			done <- result{err == 0 && value != 0 && send(value, "boolValue") != 0, err != 0}
		})
		defer block.Release()
		js := send(objc.ID(objc.GetClass("NSString")), "stringWithUTF8String:", script)
		send(view, "evaluateJavaScript:completionHandler:", js, block)
	})
	select {
	case r := <-done:
		if r.failed {
			t.Fatal("inspector JavaScript failed")
		}
		return r.value
	case <-time.After(5 * time.Second):
		t.Fatal("inspector JavaScript timed out")
		return false
	}
}

func TestNativeWebInspector(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires macOS desktop")
	}
	var host *mygo.Window
	var web *WebView
	var err error
	ready := false
	mygo.RunOnMain(func() {
		host = mygo.NewWindow(mygo.WindowOptions{Title: "WebView inspector review", Width: 1000, Height: 800, Hidden: true, Content: ui.View(func(c *ui.Context) {
			if web != nil {
				web.View(c).Fill()
			}
		})})
		web, err = NewWebView(host, mygo.PageOptions{DevTools: mygo.DevToolsEnabled})
		if err != nil {
			return
		}
		web.Page().OnDidFinishLoad(func() { ready = true })
		web.Page().LoadHTML(webViewFixture, "")
		host.Show()
	})
	defer host.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	waitWebView(t, "inspected page ready", func() bool { return ready && web.layout.Visible && !web.pending })
	var inspector objc.ID
	mygo.RunOnMain(func() {
		web.Page().OpenDevTools()
		// MyGo uses this same WebKit inspector SPI. Force docking so the
		// regression cannot accidentally pass with a separate inspector window.
		send(send(web.view.(*macAttachment).web, "_inspector"), "attach")
	})
	waitWebView(t, "docked inspector", func() bool {
		a := web.view.(*macAttachment)
		children := send(a.root, "subviews")
		for i := uintptr(0); i < uintptr(send(children, "count")); i++ {
			v := send(children, "objectAtIndex:", i)
			if v != a.web && send(v, "isKindOfClass:", objc.GetClass("WKWebView")) != 0 {
				inspector = v
				return bounds(v).Size.H > 0
			}
		}
		return false
	})
	until := time.Now().Add(10 * time.Second)
	for !inspectorBool(t, inspector, "document.readyState === 'complete' && !!globalThis.WI && !!WI.tabBrowser") {
		if time.Now().After(until) {
			t.Fatal("inspector document did not load")
		}
		time.Sleep(30 * time.Millisecond)
	}
	// A deterministic input inside the real inspector document verifies native
	// pointer and keyboard delivery without depending on WebKit's toolbar DOM.
	if !inspectorBool(t, inspector, `(()=>{const i=document.createElement('input');i.id='oneo-input';i.style.cssText='position:fixed;left:20px;top:40px;width:160px;height:30px;z-index:2147483647';document.body.append(i);return true})()`) {
		t.Fatal("input fixture")
	}
	click := func() {
		mygo.RunOnMain(func() {
			host.Focus()
			a := web.view.(*macAttachment)
			p := point{40, 55}
			if send(inspector, "isFlipped") == 0 {
				p.Y = bounds(inspector).Size.H - p.Y
			}
			p = nativeConvert(a.parent, objc.RegisterName("convertPoint:fromView:"), p, inspector)
			hit := nativeHit(a.parent, objc.RegisterName("hitTest:"), p)
			if hit == 0 || send(hit, "isDescendantOf:", inspector) == 0 {
				t.Error("inspector hit was intercepted by the native Surface")
			}
			p = nativeConvert(a.parent, objc.RegisterName("convertPoint:toView:"), p, 0)
			for _, kind := range []uint64{1, 2} {
				ev := send(objc.ID(objc.GetClass("NSEvent")), "mouseEventWithType:location:modifierFlags:timestamp:windowNumber:context:eventNumber:clickCount:pressure:", kind, p, uint64(0), float64(0), send(a.host, "windowNumber"), uintptr(0), int64(1), int64(1), float32(1))
				send(a.host, "sendEvent:", ev)
			}
		})
	}
	click()
	until = time.Now().Add(5 * time.Second)
	for !inspectorBool(t, inspector, "document.activeElement.id === 'oneo-input'") {
		if time.Now().After(until) {
			t.Fatal("inspector input did not receive pointer focus")
		}
		time.Sleep(30 * time.Millisecond)
	}
	mygo.RunOnMain(func() {
		a := web.view.(*macAttachment)
		chars := send(objc.ID(objc.GetClass("NSString")), "stringWithUTF8String:", "q")
		for _, kind := range []uint64{10, 11} {
			ev := send(objc.ID(objc.GetClass("NSEvent")), "keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:", kind, point{}, uint64(0), float64(0), send(a.host, "windowNumber"), uintptr(0), chars, chars, uintptr(12<<16))
			send(a.host, "sendEvent:", ev)
		}
	})
	until = time.Now().Add(5 * time.Second)
	for !inspectorBool(t, inspector, "document.querySelector('#oneo-input').value === 'q'") {
		if time.Now().After(until) {
			t.Fatal("inspector input did not receive keyboard input")
		}
		time.Sleep(30 * time.Millisecond)
	}
	host.SetSize(1100, 850)
	waitWebView(t, "resized browser", func() bool { return !web.pending && web.rect.W == 1100 })
	if !inspectorBool(t, inspector, "document.querySelector('#oneo-input').blur(); true") {
		t.Fatal("blur inspector input")
	}
	click()
	until = time.Now().Add(5 * time.Second)
	for !inspectorBool(t, inspector, "document.activeElement.id === 'oneo-input'") {
		if time.Now().After(until) {
			t.Fatal("inspector stopped accepting clicks after resize")
		}
		time.Sleep(30 * time.Millisecond)
	}
	mygo.RunOnMain(func() { web.Page().CloseDevTools() })
	waitWebView(t, "inspector closed", func() bool { return !web.Page().IsDevToolsOpened() })
	clickWebButton(t, web)
}
