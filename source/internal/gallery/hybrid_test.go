package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHybridNavigationWithoutNativeResources(t *testing.T) {
	g := &Gallery{}
	nav := ui.NewTester(g.Left, 220, 900)
	if err := nav.Click("Native + WebView"); err != nil {
		t.Fatal(err)
	}
	if g.page != hybridPage {
		t.Fatal("experiment did not select its gallery page")
	}
	view := ui.NewTester(g.View, 600, 700)
	if !view.HasText("Open the gallery in a native window to load the WebView.") {
		t.Fatal(view.Texts())
	}
	if g.hybrid != nil || g.hybridStarting {
		t.Fatal("headless rendering created browser resources")
	}
	if err := nav.Click("Shell"); err != nil {
		t.Fatal(err)
	}
	view.Frame()
	if g.hybridActive {
		t.Fatal("leaving experiment did not deactivate it")
	}
	g.Close()
}

func TestHybridAlertModes(t *testing.T) {
	d := &hybridDemo{visible: true}
	tt := ui.NewTester(d.buildAlert, 800, 600)
	d.openAlert(false)
	tt.Frame()
	if !tt.HasText("Native alert / WebView test") || !d.webVisible() {
		t.Fatal("native alert did not keep the real web visible")
	}
	if !tt.HasText("This native alert covers the live browser and blocks its input. The page keeps running underneath. Auto-closes after 8 seconds; Escape cancels.") {
		t.Fatal("gallery is not demonstrating native overlay composition")
	}
	d.alertDeadline = time.Now().Add(-time.Second)
	tt.Frame()
	if d.alertOpen || !d.webVisible() {
		t.Fatal("native alert did not safely time out")
	}
	d.openAlert(true)
	tt.Frame()
	if !d.alertOpen || d.webVisible() {
		t.Fatal("visibility comparison must hide the web explicitly")
	}
	if err := tt.Click("OK"); err != nil {
		t.Fatal(err)
	}
	if d.alertOpen || !d.webVisible() {
		t.Fatal("dismiss did not restore visibility")
	}
	d.visible = false
	d.openAlert(true)
	tt.Frame()
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if d.webVisible() {
		t.Fatal("alert dismissal discarded manual hide")
	}
}
