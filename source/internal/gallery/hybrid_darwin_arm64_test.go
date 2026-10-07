//go:build darwin && arm64

package gallery

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
)

func waitHybrid(t *testing.T, what string, check func() bool) {
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

func hybridSend(id objc.ID, name string, args ...any) objc.ID {
	return id.Send(objc.RegisterName(name), args...)
}

func hybridNativeRoot(win *mygo.Window) objc.ID {
	children := hybridSend(hybridSend(objc.ID(win.NativeHandle()), "contentView"), "subviews")
	for i := uintptr(0); i < uintptr(hybridSend(children, "count")); i++ {
		root := hybridSend(children, "objectAtIndex:", i)
		views := hybridSend(root, "subviews")
		for j := uintptr(0); j < uintptr(hybridSend(views, "count")); j++ {
			if hybridSend(hybridSend(views, "objectAtIndex:", j), "isKindOfClass:", objc.GetClass("WKWebView")) != 0 {
				return root
			}
		}
	}
	return 0
}

func TestNativeHybrid(t *testing.T) {
	if os.Getenv("ONEO_GO_NATIVE_SMOKE") != "1" {
		t.Skip("requires a macOS desktop session")
	}
	initial := len(mygo.Windows())
	g := &Gallery{}
	mygo.RunOnMain(func() { g.New() })
	win := g.Shell.Window()
	defer win.Destroy()
	if len(mygo.Windows()) != initial+1 {
		t.Fatal("WebView was eagerly created on the Shell page")
	}
	mygo.RunOnMain(func() { g.page = hybridPage; win.Invalidate() })
	waitHybrid(t, "gallery browser", func() bool { return g.hybridErr != nil || g.hybrid != nil && g.hybrid.ready })
	var root objc.ID
	mygo.RunOnMain(func() {
		if g.hybridErr != nil {
			t.Error(g.hybridErr)
			return
		}
		root = hybridNativeRoot(win)
	})
	if root == 0 {
		t.Fatal("no WebView in gallery host")
	}
	d := g.hybrid
	page := d.web.Page()
	visible := func() bool { return hybridSend(root, "isHidden") == 0 }
	for _, mode := range []struct {
		name    string
		theme   mygo.ThemeSource
		width   int
		sidebar bool
	}{
		{"light", mygo.ThemeLight, 1120, true},
		{"dark", mygo.ThemeDark, 820, true},
		{"narrow", mygo.ThemeLight, 360, false},
	} {
		mygo.RunOnMain(func() {
			mygo.Theme.SetSource(mode.theme)
			win.SetSize(mode.width, 900)
			g.Shell.SetLeftExpanded(mode.sidebar)
		})
		time.Sleep(400 * time.Millisecond)
		waitHybrid(t, "visible browser", visible)
		mygo.RunOnMain(func() {
			if hybridSend(root, "window") != objc.ID(win.NativeHandle()) {
				t.Error("web is not in Gallery window")
			}
			type pair struct{ X, Y float64 }
			type rectangle struct{ Origin, Size pair }
			f := objc.Send[rectangle](root, objc.RegisterName("frame"))
			parent := hybridSend(objc.ID(win.NativeHandle()), "contentView")
			b := objc.Send[rectangle](parent, objc.RegisterName("bounds"))
			if f.Origin.X < 0 || f.Origin.Y < 0 || f.Origin.X+f.Size.X > b.Size.X+0.1 || f.Origin.Y+f.Size.Y > b.Size.Y-40 {
				t.Errorf("web escaped body/safe header: frame=%+v bounds=%+v", f, b)
			}
		})
		captureHybrid(t, win, mode.name)
	}
	mygo.RunOnMain(func() { win.SetSize(1120, 900); g.Shell.SetLeftExpanded(true) })
	time.Sleep(400 * time.Millisecond)
	mygo.RunOnMain(func() { d.openAlert(false); win.Invalidate() })
	waitHybrid(t, "occlusion probe", func() bool { return d.alertOpen && visible() })
	time.Sleep(200 * time.Millisecond)
	captureHybrid(t, win, "alert-over-web")
	waitHybrid(t, "automatic probe dismissal", func() bool { return !d.alertOpen && visible() })
	mygo.RunOnMain(func() { d.openAlert(true); win.Invalidate() })
	waitHybrid(t, "safe alert", func() bool { return d.alertOpen && !visible() })
	time.Sleep(100 * time.Millisecond)
	captureHybrid(t, win, "alert-web-hidden")
	mygo.RunOnMain(func() {
		host := objc.ID(win.NativeHandle())
		// Opening in test code does not perform the toolbar click that normally
		// gives MyGo's Surface keyboard focus.
		children := hybridSend(hybridSend(host, "contentView"), "subviews")
		for i := uintptr(0); i < uintptr(hybridSend(children, "count")); i++ {
			child := hybridSend(children, "objectAtIndex:", i)
			if hybridSend(child, "isKindOfClass:", objc.GetClass("MyGoSurfaceView")) != 0 {
				hybridSend(host, "makeFirstResponder:", child)
			}
		}
		chars := hybridSend(objc.ID(objc.GetClass("NSString")), "stringWithUTF8String:", "\u001b")
		// Apple ARM64 packs the trailing BOOL and unsigned short into one
		// stack slot. purego's variadic call needs that slot packed explicitly.
		for _, kind := range []uint64{10, 11} {
			event := hybridSend(objc.ID(objc.GetClass("NSEvent")), "keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:",
				kind, struct{ X, Y float64 }{}, uint64(0), float64(0), hybridSend(host, "windowNumber"), uintptr(0), chars, chars, uintptr(53<<16))
			if hybridSend(event, "keyCode") != 53 {
				t.Error("invalid Escape fixture")
			}
			hybridSend(host, "sendEvent:", event)
		}
	})
	waitHybrid(t, "Escape dismissal and web restoration", func() bool { return !d.alertOpen && visible() })
	mygo.RunOnMain(d.sendCount)
	waitHybrid(t, "native-to-web counter", func() bool { return !d.busy })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if value, err := page.EvalContext(ctx, "return document.querySelector('#native-count').textContent"); err != nil || value != "1" {
		t.Fatalf("counter: %v, %v", value, err)
	}
	for _, p := range []galleryPage{shellPage, lucidePage, mygoPageStart} {
		mygo.RunOnMain(func() { g.page = p; win.Invalidate() })
		waitHybrid(t, "browser hidden after navigation", func() bool { return !visible() })
		mygo.RunOnMain(func() { g.page = hybridPage; win.Invalidate() })
		waitHybrid(t, "same browser restored", visible)
		if g.hybrid != d || g.hybrid.web.Page() != page {
			t.Fatal("navigation recreated document")
		}
	}
	mygo.RunOnMain(func() { g.Shell.SetRightExpanded(true) })
	time.Sleep(400 * time.Millisecond)
	captureHybrid(t, win, "inspector")
	mygo.RunOnMain(func() { g.Shell.SetPage(g.Settings) })
	waitHybrid(t, "hidden in Settings", func() bool { return !visible() })
	time.Sleep(400 * time.Millisecond)
	captureHybrid(t, win, "settings")
	mygo.RunOnMain(func() { g.Shell.SetPage(nil) })
	waitHybrid(t, "return from Settings", visible)
	mygo.RunOnMain(func() { d.visible = false; d.web.SetVisible(false) })
	waitHybrid(t, "manual hide", func() bool { return !visible() })
	mygo.RunOnMain(func() { g.page = shellPage; win.Invalidate() })
	time.Sleep(100 * time.Millisecond)
	mygo.RunOnMain(func() { g.page = hybridPage; win.Invalidate() })
	time.Sleep(100 * time.Millisecond)
	mygo.RunOnMain(func() {
		if visible() {
			t.Error("navigation discarded manual hide")
		}
	})
	mygo.RunOnMain(func() { d.visible = true; d.web.SetVisible(true) })
	waitHybrid(t, "manual restore", visible)
	win.Destroy()
	waitHybrid(t, "gallery cleanup", func() bool { return d.closed && len(mygo.Windows()) == initial })
}

func captureHybrid(t *testing.T, win *mygo.Window, name string) {
	t.Helper()
	dir := os.Getenv("ONEO_GO_SCREENSHOTS")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	var number objc.ID
	mygo.RunOnMain(func() { number = hybridSend(objc.ID(win.NativeHandle()), "windowNumber") })
	path := filepath.Join(dir, "hybrid-"+name+".png")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "/usr/sbin/screencapture", "-x", "-o", "-l", strconv.FormatUint(uint64(number), 10), path).CombinedOutput()
	if err != nil {
		t.Fatalf("OS capture (Screen Recording permission required): %v: %s", err, out)
	}
	t.Log(path)
}
