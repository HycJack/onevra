package components

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type webViewPlacement struct {
	rect    ui.Rect
	visible bool
}
type fakeWebViewAttachment struct {
	calls  []webViewPlacement
	closes int
}

func (a *fakeWebViewAttachment) place(r ui.Rect, visible bool) {
	a.calls = append(a.calls, webViewPlacement{r, visible})
}
func (a *fakeWebViewAttachment) close() { a.closes++ }

func TestWebViewDeferredLayoutAndVisibility(t *testing.T) {
	a := &fakeWebViewAttachment{}
	var pending []func()
	w := &WebView{view: a, visible: true, post: func(f func()) { pending = append(pending, f) }}
	flush := func() {
		queue := pending
		pending = nil
		for _, f := range queue {
			f()
		}
	}
	tt := ui.NewTester(func(c *ui.Context) { w.View(c).Fill() }, 640, 480)
	if len(a.calls) != 0 || len(pending) != 1 {
		t.Fatal("native work happened during paint, or queue was not coalesced")
	}
	tt.SetSize(320, 240)
	flush()
	if got := a.calls[0]; got.rect.W != 320 || got.rect.H != 240 || !got.visible {
		t.Fatalf("stale placement: %+v", got)
	}
	w.SetVisible(false)
	w.SetVisible(false)
	flush()
	if a.calls[len(a.calls)-1].visible {
		t.Fatal("hide did not hide")
	}
	w.SetVisible(true)
	flush()
	if a.calls[len(a.calls)-1].visible {
		t.Fatal("show used stale bounds before paint")
	}
	tt.Frame()
	flush()
	if !a.calls[len(a.calls)-1].visible {
		t.Fatal("fresh paint did not restore visibility")
	}
	w.SetVisible(false)
	before := len(a.calls)
	w.Close()
	w.Close()
	flush()
	tt.SetSize(400, 300)
	w.SetVisible(true)
	flush()
	if a.closes != 1 || len(a.calls) != before {
		t.Fatal("closed component still manipulates native resources")
	}
}

func TestWebViewRejectsMissingHost(t *testing.T) {
	w, err := NewWebView(nil, mygo.PageOptions{})
	if w != nil || err == nil {
		t.Fatal("nil host accepted")
	}
}
