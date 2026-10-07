//go:build windows && (amd64 || arm64)

package components

import (
	"context"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestMain(m *testing.M) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		os.Exit(m.Run())
	}
	dir, err := os.MkdirTemp("", "oneo-hybrid-")
	if err != nil {
		panic(err)
	}
	mygo.App.SetName("oneo-go hybrid review")
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
func windowsWait(t *testing.T, what string, check func() bool) {
	t.Helper()
	until := time.Now().Add(15 * time.Second)
	for time.Now().Before(until) {
		ok := false
		mygo.RunOnMain(func() { ok = check() })
		if ok {
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Fatalf("timed out: %s", what)
}
func windowsEval(t *testing.T, w *WebView, script string) any {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	v, err := w.Page().EvalContext(ctx, script)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func TestNativeHybridComposition(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires Windows desktop and WebView2 Runtime")
	}
	var host *mygo.Window
	var web *WebView
	var err error
	ready, overlay, clicked, modal, omit := false, false, false, false, false
	mygo.RunOnMain(func() {
		host = mygo.NewWindow(mygo.WindowOptions{Title: "Windows GPU hybrid review", Width: 800, Height: 600, Hidden: true, Content: ui.View(func(c *ui.Context) {
			if web != nil && !omit {
				web.View(c).Fill()
			}
			if overlay {
				ui.Overlay(c, func() {
					if ui.Button(c, "Native overlay").Absolute().Left(0).Top(0).Width(200).Height(100).Clicked() {
						clicked = true
					}
				})
			}
			if modal {
				ui.AlertDialog(c, &modal, "Native modal", "Live page remains visible", "Cancel", "OK")
			}
		})})
		web, err = NewWebView(host, mygo.PageOptions{})
		if err != nil {
			return
		}
		web.Page().OnDidFinishLoad(func() { ready = true })
		web.Page().LoadHTML(`<button id="b" style="margin:20px" onclick="window.count++">Web button</button><script>window.count=0</script>`, "")
		host.Show()
	})
	defer host.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	windowsWait(t, "browser ready", func() bool { return web.Err() != nil || ready && web.placed && !web.pending })
	mygo.RunOnMain(func() { err = web.Err() })
	if err != nil {
		t.Fatal(err)
	}
	pos := windowsEval(t, web, `const r=document.querySelector('#b').getBoundingClientRect();return [r.x+r.width/2,r.y+r.height/2]`).([]any)
	click := func() {
		mygo.RunOnMain(func() {
			a := web.view.(*windowsAttachment)
			dpi, _, _ := winDPI.Call(a.host)
			p := winPoint{int32(pos[0].(float64) * float64(dpi) / 96), int32(pos[1].(float64) * float64(dpi) / 96)}
			winUser.NewProc("ClientToScreen").Call(a.surface, uintptr(unsafe.Pointer(&p)))
			winUser.NewProc("SetCursorPos").Call(uintptr(p.X), uintptr(p.Y))
			winUser.NewProc("mouse_event").Call(2, 0, 0, 0, 0)
			winUser.NewProc("mouse_event").Call(4, 0, 0, 0, 0)
		})
	}
	click()
	time.Sleep(150 * time.Millisecond)
	if got := windowsEval(t, web, "return count"); got != float64(1) {
		t.Fatalf("native WebView2 click: %v", got)
	}
	mygo.RunOnMain(func() { overlay = true; host.Invalidate() })
	windowsWait(t, "overlay routing", func() bool { return !web.route(40, 40, false) })
	click()
	windowsWait(t, "native overlay click", func() bool { return clicked })
	if got := windowsEval(t, web, "return count"); got != float64(1) {
		t.Fatalf("overlay leaked click: %v", got)
	}
	mygo.RunOnMain(func() { overlay = false; modal = true; host.Invalidate() })
	windowsWait(t, "modal barrier", func() bool { return !web.pending && !web.route(40, 40, false) && web.view.(*windowsAttachment).shown })
	mygo.RunOnMain(func() { modal = false; omit = true; host.Invalidate() })
	windowsWait(t, "automatic unmount", func() bool { return !web.view.(*windowsAttachment).shown })
	mygo.RunOnMain(func() { omit = false; host.Invalidate() })
	windowsWait(t, "remount", func() bool { return web.view.(*windowsAttachment).shown })
	if got := windowsEval(t, web, "return count"); got != float64(1) {
		t.Fatal("document state lost")
	}
	mygo.RunOnMain(func() { web.Close(); web.Close() })
	if host.IsDestroyed() || !web.donor.IsDestroyed() {
		t.Fatal("cleanup ownership")
	}
}
