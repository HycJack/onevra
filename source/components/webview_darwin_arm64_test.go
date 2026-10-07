//go:build darwin && arm64

package components

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestMain(m *testing.M) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "oneo-webview-")
	if err != nil {
		panic(err)
	}
	mygo.App.SetName("oneo-go WebView review")
	mygo.App.SetPath(mygo.PathUserData, dir)
	mygo.App.OnWindowAllClosed(func() {})
	result := make(chan int, 1)
	mygo.App.WhenReady(func() { go func() { result <- m.Run(); mygo.App.Quit() }() })
	if err := mygo.App.Run(); err != nil {
		panic(err)
	}
	code := <-result
	os.RemoveAll(dir)
	os.Exit(code)
}

func waitWebView(t *testing.T, what string, check func() bool) {
	t.Helper()
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		ok := false
		mygo.RunOnMain(func() { ok = check() })
		if ok {
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func evalWebView(t *testing.T, w *WebView, script string) any {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	v, err := w.Page().EvalContext(ctx, script)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

const webViewFixture = `<!doctype html><html><body style="padding:20px">
<button id="increment" onclick="this.textContent=++window.count">0</button>
<textarea></textarea><script>window.count=0; window.receiveNative=n=>window.nativeCount=n;</script>
</body></html>`

func TestNativeWebView(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a macOS desktop session")
	}
	initial := len(mygo.Windows())
	var host *mygo.Window
	var first, second *WebView
	var err error
	var ready1, ready2 bool
	omit := false
	mygo.RunOnMain(func() {
		host = mygo.NewWindow(mygo.WindowOptions{Title: "WebView component review", Width: 900, Height: 600, Hidden: true,
			Content: ui.View(func(c *ui.Context) {
				ui.Column(c).Fill().Padding(20).Gap(16).Children(func() {
					ui.Text(c, "Native MyGo header")
					ui.Row(c).Fill().Gap(20).Children(func() {
						if first != nil && !omit {
							first.View(c).Fill()
						} else {
							ui.Box(c).Fill()
						}
						if second != nil {
							second.View(c).Width(200).FillHeight()
						}
					})
				})
			}),
		})
		first, err = NewWebView(host, mygo.PageOptions{})
		if err != nil {
			return
		}
		second, err = NewWebView(host, mygo.PageOptions{})
		if err != nil {
			return
		}
		first.Page().OnDidFinishLoad(func() { ready1 = true })
		second.Page().OnDidFinishLoad(func() { ready2 = true })
		first.Page().LoadHTML(webViewFixture, "")
		second.Page().LoadHTML(webViewFixture, "")
		host.Show()
	})
	defer host.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	waitWebView(t, "both pages", func() bool {
		return ready1 && ready2 && first.placed && !first.pending && second.placed && !second.pending && first.layout.Visible && second.layout.Visible
	})
	if first.Page().Window() == host || first.donor.IsVisible() || second.donor.IsVisible() {
		t.Fatal("browser owner must be a distinct hidden window")
	}
	mygo.RunOnMain(func() {
		if invalid, err := NewWebView(first.donor, mygo.PageOptions{}); err == nil || invalid != nil {
			t.Error("accepted a web window as a native host")
		}
	})
	for _, width := range []int{900, 700, 1000} {
		host.SetSize(width, 600)
		time.Sleep(100 * time.Millisecond)
		waitWebView(t, "geometry", func() bool {
			a := first.view.(*macAttachment)
			f := objc.Send[rect](a.root, objc.RegisterName("frame"))
			return !first.pending && f.Size.W == float64(first.rect.W) && f.Size.H == float64(first.rect.H)
		})
		mygo.RunOnMain(func() {
			for _, w := range []*WebView{first, second} {
				a := w.view.(*macAttachment)
				if send(a.root, "window") != a.host {
					t.Error("browser did not attach to host")
				}
				f := objc.Send[rect](a.root, objc.RegisterName("frame"))
				if f.Origin.X != float64(w.rect.X) || math.Abs(f.Origin.Y-(bounds(a.parent).Size.H-float64(w.rect.Y)-float64(w.rect.H))) > 0.01 {
					t.Error("incorrect slot coordinates")
				}
				hit := send(a.parent, "hitTest:", point{f.Origin.X + 10, f.Origin.Y + 10})
				if hit == 0 || send(hit, "isDescendantOf:", a.root) == 0 {
					t.Error("slot does not hit its browser")
				}
			}
		})
	}
	clickWebButton(t, first)
	typeInWeb(t, first)
	if got := evalWebView(t, second, "return window.count"); got != float64(0) {
		t.Fatalf("documents share state: %v", got)
	}
	if got := evalWebView(t, first, "return window.receiveNative(7)"); got != float64(7) {
		t.Fatalf("Eval = %v", got)
	}

	mygo.RunOnMain(func() { first.SetVisible(false); omit = true; host.Invalidate() })
	waitWebView(t, "hidden browser", func() bool { return !first.pending && send(first.view.(*macAttachment).root, "isHidden") != 0 })
	mygo.RunOnMain(func() {
		a := first.view.(*macAttachment)
		if send(a.host, "firstResponder") != a.surface {
			t.Error("hide did not return focus to native Surface")
		}
		first.SetVisible(true)
	})
	// No paint while omitted: SetVisible(true) must not show stale geometry.
	time.Sleep(100 * time.Millisecond)
	mygo.RunOnMain(func() {
		if send(first.view.(*macAttachment).root, "isHidden") == 0 {
			t.Error("omitted view became visible")
		}
		omit = false
		host.Invalidate()
	})
	waitWebView(t, "restored browser", func() bool { return !first.pending && send(first.view.(*macAttachment).root, "isHidden") == 0 })
	if got := evalWebView(t, first, "return window.count"); got != float64(1) {
		t.Fatal("hide/show lost page state")
	}
	first.donor.Close()
	if first.donor.IsDestroyed() {
		t.Fatal("page can close its attached owner")
	}
	typeInWeb(t, second)
	mygo.RunOnMain(func() {
		second.Close()
		second.Close()
		a := first.view.(*macAttachment)
		if send(a.host, "firstResponder") != a.surface {
			t.Error("early Close did not return focus to native Surface")
		}
	})
	if host.IsDestroyed() || !second.donor.IsDestroyed() {
		t.Fatal("early Close destroyed host or retained owner")
	}
	// A fresh instance must work after closing one, alongside the remaining one.
	mygo.RunOnMain(func() { second, err = NewWebView(host, mygo.PageOptions{}); host.Invalidate() })
	if err != nil {
		t.Fatal(err)
	}
	mygo.RunOnMain(func() { first.SetVisible(false); host.Destroy(); first.Close() })
	waitWebView(t, "cleanup", func() bool { return first.closed && second.closed && len(mygo.Windows()) == initial })
}

func clickWebButton(t *testing.T, w *WebView) {
	t.Helper()
	mygo.RunOnMain(w.host.Focus)
	// A committed zero-size layout is not yet a visible browser. Injecting a
	// click then can land in window chrome rather than the document.
	waitWebView(t, "visible browser in key window", func() bool {
		if w.closed || w.host.IsDestroyed() || !w.placed || w.pending || !w.layout.Visible {
			return false
		}
		a := w.view.(*macAttachment)
		return a.shown && send(a.root, "isHidden") == 0 && send(a.host, "isKeyWindow") != 0
	})
	v := evalWebView(t, w, "const r=document.querySelector('#increment').getBoundingClientRect(); return {x:r.x+r.width/2,y:r.y+r.height/2}").(map[string]any)
	mygo.RunOnMain(func() {
		a := w.view.(*macAttachment)
		p := point{float64(w.rect.X) + v["x"].(float64), bounds(a.parent).Size.H - float64(w.rect.Y) - v["y"].(float64)}
		for _, kind := range []uint64{1, 2} {
			event := send(objc.ID(objc.GetClass("NSEvent")), "mouseEventWithType:location:modifierFlags:timestamp:windowNumber:context:eventNumber:clickCount:pressure:",
				kind, p, uint64(0), float64(0), send(a.host, "windowNumber"), uintptr(0), int64(1), int64(1), float32(1))
			send(a.host, "sendEvent:", event)
		}
	})
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if evalWebView(t, w, "return window.count") == float64(1) {
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Fatal("AppKit mouse events did not reach the web button")
}

func typeInWeb(t *testing.T, w *WebView) {
	t.Helper()
	evalWebView(t, w, "document.querySelector('textarea').focus(); return true")
	mygo.RunOnMain(func() {
		a := w.view.(*macAttachment)
		chars := send(objc.ID(objc.GetClass("NSString")), "stringWithUTF8String:", "q")
		// Apple ARM64 packs BOOL and unsigned short into the first stack slot;
		// purego's variadic call needs that slot packed explicitly.
		for _, kind := range []uint64{10, 11} {
			event := send(objc.ID(objc.GetClass("NSEvent")), "keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:",
				kind, point{}, uint64(0), float64(0), send(a.host, "windowNumber"), uintptr(0), chars, chars, uintptr(12<<16))
			if send(event, "keyCode") != 12 {
				t.Error("invalid keyboard fixture")
			}
			send(a.host, "sendEvent:", event)
		}
	})
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if evalWebView(t, w, "return document.querySelector('textarea').value") == "q" {
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Fatal("AppKit key events did not reach the web input")
}
