package components

import (
	"runtime"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type fakeWindow struct {
	bounds                 mygo.Rectangle
	maximized, fullscreen  bool
	limit                  int
	changes, invalidations int
}

func (w *fakeWindow) Bounds() mygo.Rectangle { return w.bounds }
func (w *fakeWindow) SetBounds(b mygo.Rectangle) {
	if w.limit > 0 {
		b.Width = min(b.Width, w.limit)
	}
	w.bounds = b
	w.changes++
}
func (w *fakeWindow) IsMaximized() bool  { return w.maximized }
func (w *fakeWindow) IsFullScreen() bool { return w.fullscreen }
func (w *fakeWindow) Invalidate()        { w.invalidations++ }

func fixture() (*Shell, *fakeWindow) {
	s := newShell(ShellOptions{
		Left:    func(c *ui.Context) { ui.Text(c, "Navigation") },
		Content: func(c *ui.Context) { ui.Text(c, "Document") },
		Right:   func(c *ui.Context) { ui.Text(c, "Inspector") },
	})
	w := &fakeWindow{bounds: mygo.Rectangle{X: 50, Y: 60, Width: 820, Height: 600}}
	s.host = w
	return s, w
}

func TestShellNativeOptions(t *testing.T) {
	s := newShell(ShellOptions{Title: "Example", Width: 820, Height: 600, MinWidth: 360, MinHeight: 320})
	o := s.nativeOptions()
	if o.Title != "Example" || o.Width != 820 || o.Height != 600 || o.MinWidth != 360 || o.MinHeight != 320 {
		t.Fatal(o)
	}
	if o.Content == nil || o.TitleBarStyle != mygo.TitleBarHidden || o.TrafficLightPosition == nil || *o.TrafficLightPosition != (mygo.Point{X: 16, Y: 16}) {
		t.Fatal(o)
	}
	material := mygo.VibrancyNone
	if runtime.GOOS == "darwin" {
		material = mygo.VibrancyUnderWindow
	}
	if o.Transparent != (runtime.GOOS == "darwin") || o.Vibrancy != material || o.Frameless || o.DisableResize || o.StateKey != "" {
		t.Fatal(o)
	}
}

// Bounds unit tests apply the final animation step explicitly.
func settleWidth(s *Shell) {
	width := 0
	if s.rightOpen && s.page == nil {
		width = s.options.RightWidth
	}
	s.syncWidth(width)
}

func TestShellWidthAndPages(t *testing.T) {
	s, w := fixture()
	original := w.bounds
	s.SetLeftExpanded(false)
	settleWidth(s)
	if w.bounds != original {
		t.Fatal("left toggle resized window")
	}
	s.SetRightExpanded(true)
	settleWidth(s)
	if w.bounds.Width != 1100 || w.bounds.X != original.X || w.bounds.Y != original.Y || w.bounds.Height != original.Height {
		t.Fatal(w.bounds)
	}
	s.SetRightExpanded(true)
	settleWidth(s)
	if w.changes != 1 {
		t.Fatal("repeated expansion resized window")
	}
	// Preserve a user's resize while the inspector is open.
	w.bounds.Width += 100
	s.SetPage(func(c *ui.Context) { ui.Text(c, "Settings") })
	settleWidth(s)
	if w.bounds.Width != 920 || !s.RightExpanded() || s.LeftExpanded() {
		t.Fatal("page replacement lost sidebar state or width")
	}
	s.SetPage(nil)
	settleWidth(s)
	if w.bounds.Width != 1200 {
		t.Fatal(w.bounds)
	}
	s.SetRightExpanded(false)
	settleWidth(s)
	if w.bounds.Width != 920 {
		t.Fatal(w.bounds)
	}
}

func TestShellConstrainedWidth(t *testing.T) {
	s, w := fixture()
	w.limit = 900
	s.SetRightExpanded(true)
	settleWidth(s)
	s.SetRightExpanded(false)
	settleWidth(s)
	if w.bounds.Width != 820 {
		t.Fatal("must undo only actual extension", w.bounds)
	}
	for _, fullscreen := range []bool{false, true} {
		w.maximized, w.fullscreen = !fullscreen, fullscreen
		s.SetRightExpanded(true)
		settleWidth(s)
		if w.bounds.Width != 820 {
			t.Fatal("resized maximized/fullscreen window")
		}
		w.maximized, w.fullscreen = false, false
		settleWidth(s)
		if w.bounds.Width != 900 {
			t.Fatal("deferred expansion missing")
		}
		w.maximized = true
		s.SetPage(func(*ui.Context) {})
		settleWidth(s)
		w.maximized = false
		settleWidth(s)
		if w.bounds.Width != 820 {
			t.Fatal("deferred contraction missing")
		}
		s.SetRightExpanded(false)
		settleWidth(s)
		s.SetPage(nil)
		settleWidth(s)
	}
}

func TestShellLayoutAndReplacement(t *testing.T) {
	s, w := fixture()
	tt := ui.NewTester(s.view, 820, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if !tt.HasText("Navigation") || !tt.HasText("Document") || tt.HasText("Inspector") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Toggle sidebar"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Navigation") || w.changes != 0 {
		t.Fatal("left sidebar did not collapse internally")
	}
	if err := tt.Click("Toggle inspector"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Inspector") || w.bounds.Width != 1100 {
		t.Fatal(tt.Texts(), w.bounds)
	}
	s.SetPage(func(c *ui.Context) {
		ui.Text(c, "Settings")
		if ui.Button(c, "Back").Clicked() {
			s.SetPage(nil)
		}
	})
	tt = ui.NewTester(s.view, 820, 600)
	_, hasLeftToggle := tt.Find("Toggle sidebar")
	_, hasRightToggle := tt.Find("Toggle inspector")
	if !tt.HasText("Settings") || tt.HasText("Document") || tt.HasText("Inspector") || hasLeftToggle || hasRightToggle {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(sidebarDuration + 30*time.Millisecond)
	tt.Frame()
	if !tt.HasText("Document") || !tt.HasText("Inspector") || tt.HasText("Navigation") {
		t.Fatal(tt.Texts())
	}
}

func TestShellControlClearanceAndStableContent(t *testing.T) {
	s, _ := fixture()
	tt := ui.NewTester(s.view, 820, 600)
	before, ok := tt.Find("Document")
	if !ok {
		t.Fatal("missing content")
	}
	s.SetRightExpanded(true)
	tt = ui.NewTester(s.view, 1100, 600)
	after, ok := tt.Find("Document")
	if !ok || before.X != after.X || before.W != after.W {
		t.Fatalf("content moved: %+v -> %+v", before, after)
	}
	for _, scale := range []float32{1, 1.5, 2} {
		tt = ui.NewTester(s.view, 360, 480)
		tt.SetTitleBar(ui.TitleBar{Height: 48, Left: 80, Right: 100})
		tt.SetPreferences(ui.Preferences{TextScale: scale})
		for _, label := range []string{"Toggle sidebar", "Toggle inspector"} {
			r, ok := tt.Find(label)
			if !ok || r.X < 80 || r.X+r.W > 260 || r.Y < 0 || r.Y+r.H > max(48, 48*scale) {
				t.Fatalf("scale %v: control %q overlaps native chrome: %+v", scale, label, r)
			}
		}
	}
}

func TestShellSidebarAnimations(t *testing.T) {
	s, w := fixture()
	tt := ui.NewTester(s.view, 820, 600)
	contentX := func() float32 {
		t.Helper()
		r, ok := tt.Find("Document")
		if !ok {
			t.Fatal("missing main content")
		}
		return r.X
	}
	step := func() {
		time.Sleep(55 * time.Millisecond)
		tt.Frame()
	}
	finish := func() {
		time.Sleep(sidebarDuration + 30*time.Millisecond)
		tt.Frame()
	}
	startX := contentX()
	s.SetLeftExpanded(false)
	tt.Frame()
	step()
	if x := contentX(); x <= 0 || x >= startX || w.bounds.Width != 820 {
		t.Fatalf("left collapse should be intermediate without window resizing: x=%v, bounds=%+v", x, w.bounds)
	}
	// Reverse before completion, without losing the current layout.
	s.SetLeftExpanded(true)
	tt.Frame()
	finish()
	if x := contentX(); x != startX {
		t.Fatalf("left reversal ended at %v", x)
	}
	s.SetLeftExpanded(false)
	tt.Frame()
	finish()
	if contentX() != 0 || tt.HasText("Navigation") {
		t.Fatal("left collapse did not finish")
	}

	s.SetRightExpanded(true)
	tt.Frame()
	step()
	if w.bounds.Width <= 820 || w.bounds.Width >= 1100 || w.bounds.X != 50 || w.bounds.Height != 600 {
		t.Fatalf("right expansion must resize incrementally: %+v", w.bounds)
	}
	s.SetRightExpanded(false)
	tt.Frame()
	finish()
	if w.bounds.Width != 820 || tt.HasText("Inspector") {
		t.Fatal("right reversal did not restore base width", w.bounds)
	}
	s.SetRightExpanded(true)
	tt.Frame()
	finish()
	if w.bounds.Width != 1100 {
		t.Fatal("right expansion did not finish", w.bounds)
	}
	// Replacing a page during contraction reverses cleanly on return.
	s.SetPage(func(c *ui.Context) { ui.Text(c, "Settings") })
	tt.Frame()
	step()
	if w.bounds.Width <= 820 || w.bounds.Width >= 1100 || !tt.HasText("Settings") {
		t.Fatal("replacement contraction missing", w.bounds)
	}
	s.SetPage(nil)
	tt.Frame()
	finish()
	if w.bounds.Width != 1100 || !tt.HasText("Inspector") {
		t.Fatal("replacement reversal did not restore inspector")
	}

	// Enabling reduced motion midway must snap both layout and native width.
	s.SetLeftExpanded(true)
	s.SetRightExpanded(false)
	tt.Frame()
	step()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if contentX() != startX || w.bounds.Width != 820 || tt.HasText("Inspector") {
		t.Fatal("Reduce Motion did not finish active animations")
	}
	s.SetRightExpanded(true)
	s.SetLeftExpanded(false)
	tt.Frame()
	if contentX() != 0 || w.bounds.Width != 1100 {
		t.Fatal("Reduce Motion animated a new toggle")
	}
}

func TestShellSidebarsRevealWithoutReflow(t *testing.T) {
	var left, right, main *ui.Element
	s, w := fixture()
	s.options.Left = func(c *ui.Context) { left = ui.Box(c).FillWidth().Height(32) }
	s.options.Right = func(c *ui.Context) { right = ui.Box(c).FillWidth().Height(32) }
	s.options.Content = func(c *ui.Context) { main = ui.Box(c).Fill() }
	tt := ui.NewTester(s.view, 820, 600)
	leftBefore := left.Bounds()
	s.SetLeftExpanded(false)
	tt.Frame()
	time.Sleep(55 * time.Millisecond)
	tt.Frame()
	if got := left.Bounds(); got != leftBefore {
		t.Fatalf("left content reflowed instead of being covered: %+v -> %+v", leftBefore, got)
	}
	if x := main.Bounds().X; x <= 0 || x >= leftBefore.W {
		t.Fatalf("main surface did not move over sidebar: %v", x)
	}
	time.Sleep(sidebarDuration + 30*time.Millisecond)
	tt.Frame()
	if r := main.Bounds(); r.X != 0 || r.W != 820 {
		t.Fatal("collapsed sidebar reserved space", r)
	}

	s.SetRightExpanded(true)
	tt.Frame()
	time.Sleep(55 * time.Millisecond)
	tt.Frame()
	// Follow the fake native window, as the desktop does after each resize.
	tt.SetSize(w.bounds.Width, 600)
	if r := right.Bounds(); r.W != 280 || r.H != 32 {
		t.Fatal("inspector reflowed while being revealed", r)
	}
	if w.bounds.Width <= 820 || w.bounds.Width >= 1100 {
		t.Fatal("missing intermediate window size")
	}
	if r := main.Bounds(); r.X != 0 || r.W <= 0 || right.Bounds().X != r.X+r.W {
		t.Fatal("inspector overlaps the main surface", r, right.Bounds())
	}
	time.Sleep(sidebarDuration + 30*time.Millisecond)
	tt.Frame()
	tt.SetSize(w.bounds.Width, 600)
	if r := right.Bounds(); r.W != 280 || r.X != 820 {
		t.Fatal("inspector ended at wrong geometry", r)
	}
}

func TestShellSidebarResize(t *testing.T) {
	s, w := fixture()
	tt := ui.NewTester(s.view, 820, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	separator := func() ui.Rect {
		t.Helper()
		r, ok := tt.Find("Resize sidebar")
		if !ok {
			t.Fatal("missing resize handle")
		}
		return r
	}
	drag := func(dx float32) {
		r := separator()
		x, y := r.X+r.W/2, r.Y+30
		tt.Move(x, y)
		if tt.Cursor() != ui.CursorResizeEW {
			t.Fatal("missing resize cursor")
		}
		tt.Press(x, y)
		tt.Move(x+dx, y)
		tt.Release(x+dx, y)
	}
	drag(50)
	if s.leftWidth != 270 || w.bounds.Width != 820 {
		t.Fatal("drag did not resize interior sidebar", s.leftWidth, w.bounds)
	}
	tt.Key(0, ui.KeyRight)
	if s.leftWidth != 280 {
		t.Fatal("keyboard resize did not move separator", s.leftWidth)
	}
	drag(1000)
	if s.leftWidth != 328 {
		t.Fatal("drag exceeded available width", s.leftWidth)
	}
	drag(-1000)
	if s.leftWidth != 120 {
		t.Fatal("drag exceeded minimum width", s.leftWidth)
	}
	drag(80)
	s.SetLeftExpanded(false)
	tt.Frame()
	if _, ok := tt.Find("Resize sidebar"); ok {
		t.Fatal("collapsed sidebar kept resize hit area")
	}
	s.SetLeftExpanded(true)
	tt.Frame()
	if s.leftWidth != 200 {
		t.Fatal("toggle lost resized width")
	}
	if r := separator(); r.X+r.W/2 != 200 || r.Y < 48 {
		t.Fatal("misplaced resize hit area", r)
	}
	s.SetPage(func(*ui.Context) {})
	tt.Frame()
	if _, ok := tt.Find("Resize sidebar"); ok {
		t.Fatal("replacement page kept resize hit area")
	}
}

func TestShellSeparatorsDoNotReserveSpace(t *testing.T) {
	var main *ui.Element
	s, _ := fixture()
	s.host = nil // Exercise layout independently of native resize acknowledgements.
	s.options.Content = func(c *ui.Context) { main = ui.Box(c).Fill() }
	s.SetRightExpanded(true)
	view := func(c *ui.Context) {
		theme := *c.Theme()
		theme.Background = ui.RGB(0, 0, 200)
		theme.Border = ui.RGB(200, 0, 0)
		c.SetTheme(&theme)
		s.view(c)
	}
	tt := ui.NewTester(view, 1100, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	assertPixel := func(x, y int, red bool) {
		t.Helper()
		p := tt.Image().RGBAAt(x, y)
		if (red && (p.R != 200 || p.B != 0)) || (!red && (p.R != 0 || p.B != 200)) {
			t.Fatalf("pixel (%d,%d) = %+v, separator=%v", x, y, p, red)
		}
	}
	if r := main.Bounds(); r.X != 220 || r.W != 600 {
		t.Fatal("separator reserved main layout space", r)
	}
	for _, y := range []int{2, 400, 599} {
		assertPixel(220, y, true)
		assertPixel(221, y, false)
		assertPixel(820, y, true)
		assertPixel(821, y, false)
	}
	s.SetLeftExpanded(false)
	tt.Frame()
	if r := main.Bounds(); r.X != 0 || r.W != 820 {
		t.Fatal("closed left separator reserved space", r)
	}
	assertPixel(0, 400, false)
	assertPixel(820, 400, true)
	s.SetRightExpanded(false)
	tt.Frame()
	if r := main.Bounds(); r.X != 0 || r.W != 1100 {
		t.Fatal("closed right separator reserved space", r)
	}
	assertPixel(820, 400, false)
	s.SetPage(func(*ui.Context) {})
	tt.Frame()
	assertPixel(220, 400, false)
	assertPixel(820, 400, false)
}

func TestShellOptionalInspector(t *testing.T) {
	s := newShell(ShellOptions{})
	s.SetRightExpanded(true)
	tt := ui.NewTester(s.view, 360, 480)
	_, hasToggle := tt.Find("Toggle inspector")
	if s.RightExpanded() || hasToggle {
		t.Fatal("disabled inspector exposed")
	}
}

func TestShellBackgroundsAndScaling(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s, _ := fixture()
		s.SetRightExpanded(true)
		tt := ui.NewTester(s.view, 1100, 600)
		tt.SetDark(dark)
		for _, p := range []struct {
			x     int
			alpha uint32
		}{{100, 0}, {400, 65535}, {1000, 65535}} {
			for _, y := range []int{2, 400, 599} {
				_, _, _, a := tt.Image().At(p.x, y).RGBA()
				want := p.alpha
				if runtime.GOOS != "darwin" {
					want = 65535
				}
				if a != want {
					t.Fatalf("dark=%v x=%d y=%d alpha=%d want=%d", dark, p.x, y, a, want)
				}
			}
		}
		for _, label := range []string{"Navigation", "Document", "Inspector"} {
			r, ok := tt.Find(label)
			if !ok || r.Y < 48 {
				t.Fatalf("%q overlaps header: %+v", label, r)
			}
		}
		s.SetLeftExpanded(false)
		tt = ui.NewTester(s.view, 1100, 600)
		_, _, _, a := tt.Image().At(100, 2).RGBA()
		if a != 65535 {
			t.Fatal("collapsed sidebar left a transparent header gap")
		}
		s.SetPage(func(*ui.Context) {})
		tt = ui.NewTester(s.view, 820, 600)
		for _, y := range []int{2, 400} {
			_, _, _, a := tt.Image().At(100, y).RGBA()
			if a != 65535 {
				t.Fatal("replacement page is not opaque")
			}
		}
		s.SetPage(nil)
		tt = ui.NewTester(s.view, 360, 480)
		tt.SetPreferences(ui.Preferences{TextScale: 1.5})
		if !tt.HasText("Document") {
			t.Fatal(tt.Texts())
		}
	}
}
