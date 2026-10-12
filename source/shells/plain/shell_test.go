package plain

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

type testWindow struct {
	minimized, closed, frames int
	maximized                 bool
}

func (w *testWindow) Invalidate()       { w.frames++ }
func (w *testWindow) Minimize()         { w.minimized++ }
func (w *testWindow) Close()            { w.closed++ }
func (w *testWindow) ToggleMaximize()   { w.maximized = !w.maximized }
func (w *testWindow) IsMaximized() bool { return w.maximized }

func TestNativeConfiguration(t *testing.T) {
	s := newShell(ShellOptions{Hidden: true, Width: 200, Height: 100})
	o := s.nativeOptions()
	if !o.Hidden || o.TitleBarStyle != mygo.TitleBarHidden || o.Frameless != (runtime.GOOS == "windows") {
		t.Fatalf("unexpected native chrome: %+v", o)
	}
	if o.Width < o.MinWidth || o.Height < o.MinHeight {
		t.Fatal("initial window violates minimum")
	}
	if o.TrafficLightPosition == nil || *o.TrafficLightPosition != (mygo.Point{X: 16, Y: 16}) {
		t.Fatal("missing traffic light placement")
	}
	if o.StateKey != "" {
		t.Fatalf("unrequested state key: %q", o.StateKey)
	}
	if o := newShell(ShellOptions{StateKey: "settings"}).nativeOptions(); o.StateKey != "settings" {
		t.Fatalf("state key: %q", o.StateKey)
	}
}

func TestContentAndDialogClearance(t *testing.T) {
	for _, scale := range []float32{1, 1.5, 2} {
		open := false
		var height float32
		s := newShell(ShellOptions{Content: func(c *ui.Context) {
			height = HeaderHeight(c)
			components.Button(c, "Body action", components.ButtonPrimary)
			components.Dialog(c, &open, components.DialogOptions{Label: "Settings"}, func(p ui.Element) { p.Width(1000).Height(1000) })
		}})
		tt := ui.NewTester(s.view, 760, 560)
		tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
		r, ok := tt.Find("Body action")
		if !ok || r.Y < height {
			t.Fatalf("content overlaps header: %+v / %v", r, height)
		}
		open = true
		tt.Move(0, 0)
		r, ok = tt.Find("Settings")
		if !ok || r.Y < height {
			t.Fatalf("dialog overlaps header: %+v / %v", r, height)
		}
		tt.ClickAt(120, height/2)
		if !open {
			t.Fatal("dragging strip dismissed modal")
		}
	}
}

func TestCaptionAndLiveUpdates(t *testing.T) {
	s := newShell(ShellOptions{})
	w := &testWindow{}
	s.host = w
	tt := ui.NewTester(s.captionButtons, 400, 60)
	for _, name := range []string{"Minimize window", "Maximize window", "Restore window", "Close window"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	if w.minimized != 1 || w.closed != 1 || w.maximized {
		t.Fatalf("caption actions: %+v", w)
	}
	s.SetLabels(Labels{Close: "关闭窗口"})
	s.SetAppearance(theme.Dark)
	tt.Frame()
	if err := tt.Click("关闭窗口"); err != nil {
		t.Fatal(err)
	}
	if w.closed != 2 || w.frames != 2 {
		t.Fatalf("live updates: %+v", w)
	}
	var dark bool
	s.options.Content = func(c *ui.Context) { dark = c.Theme().Dark }
	ui.NewTester(s.view, 760, 560)
	if !dark {
		t.Fatal("theme not installed before content")
	}
}
