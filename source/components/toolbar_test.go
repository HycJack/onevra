package components

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

func TestToolbarIconButton(t *testing.T) {
	for _, dark := range []bool{false, true} {
		clicks := 0
		disabled := false
		var button *ui.Element
		var size float32
		tt := ui.NewTester(func(c *ui.Context) {
			c.Root().Background(ui.Transparent)
			size = max(32, c.Theme().Rem(2))
			ui.Row(c).Fill().DragWindow().Children(func() {
				button = ToolbarIconButton(c, lucide.Settings(), "Settings").Disabled(disabled)
				if button.Clicked() {
					clicks++
				}
			})
		}, 200, 80)
		tt.SetDark(dark)
		r, ok := tt.Find("Settings")
		if !ok || r.W != 32 || r.H != 32 {
			t.Fatal("missing accessible square button", r)
		}
		alpha := func() uint8 { return tt.Image().RGBAAt(5, 40).A }
		tt.Move(150, 50)
		if alpha() != 0 {
			t.Fatal("resting button has a face")
		}
		tt.Move(5, 40)
		hover := alpha()
		if hover == 0 || hover == 255 {
			t.Fatal("hover must be translucent", hover)
		}
		tt.Press(5, 40)
		if alpha() <= hover || alpha() == 255 {
			t.Fatal("press must be a stronger translucent tint", alpha())
		}
		tt.Release(5, 40)
		if clicks != 1 {
			t.Fatal("click lost to draggable titlebar", clicks)
		}
		tt.Key(0, ui.KeyEnter)
		tt.Key(0, ui.KeySpace)
		if clicks != 3 {
			t.Fatal("keyboard activation failed", clicks)
		}
		disabled = true
		tt.Frame()
		if err := tt.Click("Settings"); err != nil {
			t.Fatal(err)
		}
		tt.Key(0, ui.KeyEnter)
		if clicks != 3 || alpha() != 0 {
			t.Fatal("disabled button still reacts", clicks, alpha())
		}
		tt.SetPreferences(ui.Preferences{TextScale: 1.5})
		if r := button.Bounds(); r.W != size || r.H != size {
			t.Fatal("button does not scale", r)
		}
	}
}

func TestShellToolbarPlacement(t *testing.T) {
	s, _ := fixture()
	var body *ui.Element
	clicks := 0
	s.options.Left = func(c *ui.Context) { body = ui.Box(c).Fill() }
	s.options.LeftToolbar = func(c *ui.Context) {
		ui.Spacer(c)
		if ToolbarIconButton(c, lucide.Settings(), "Settings").Clicked() {
			clicks++
		}
	}
	tt := ui.NewTester(s.view, 820, 600)
	tt.SetTitleBar(ui.TitleBar{Height: 48, Left: 80})
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	left, _ := tt.Find("Toggle sidebar")
	if runtime.GOOS == "darwin" && left.X < 96 {
		t.Fatal("sidebar toggle too close to traffic lights", left)
	}
	bottom, ok := tt.Find("Settings")
	if !ok || bottom.Y+bottom.H != 592 || bottom.X+bottom.W != 212 || body.Bounds().Y+body.Bounds().H > bottom.Y {
		t.Fatal("toolbar not pinned below sidebar content", bottom, body.Bounds())
	}
	if err := tt.Click("Settings"); err != nil || clicks != 1 {
		t.Fatal("bottom toolbar action failed", err, clicks)
	}
	s.SetLeftExpanded(false)
	tt.Frame()
	if _, ok := tt.Find("Settings"); ok {
		t.Fatal("collapsed sidebar retains toolbar")
	}
	s.SetLeftExpanded(true)
	tt.Frame()
	s.SetPage(func(*ui.Context) {})
	tt.Frame()
	if _, ok := tt.Find("Settings"); ok {
		t.Fatal("replacement page retains sidebar toolbar")
	}
	s.SetPage(nil)
	tt.Frame()
	if _, ok := tt.Find("Settings"); !ok {
		t.Fatal("returning did not restore sidebar toolbar")
	}
}
