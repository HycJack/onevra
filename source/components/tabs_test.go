package components

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestTabsClick(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := false
		if dark {
			mode = true
		}
		selected := 0
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(mode))
			Tabs(c, &selected, "General", "Appearance", "Advanced")
		}, 400, 100)
		if err := tt.Click("Appearance"); err != nil {
			t.Fatal(err)
		}
		if selected != 1 {
			t.Errorf("dark=%v: selected=%d, want 1", dark, selected)
		}
		if err := tt.Click("Advanced"); err != nil {
			t.Fatal(err)
		}
		if selected != 2 {
			t.Errorf("dark=%v: selected=%d, want 2", dark, selected)
		}
	}
}

func TestTabsIndicatorAnimation(t *testing.T) {
	for _, mode := range []bool{false, true} {
		selected := 0
		var list ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(mode))
			parts := Tabs(c, &selected, "General", "Appearance", "Advanced")
			parts.List.Draw(func(_ *ui.Painter, r ui.Rect) { list = r })
		}, 400, 100)
		underline := func() (int, int) {
			t.Helper()
			img := tt.Image()
			y := int(list.Y + list.H - 1)
			bg := img.RGBAAt(399, y)
			left, right := -1, -1
			for x := 0; x < 400; x++ {
				if img.RGBAAt(x, y) != bg {
					if left < 0 {
						left = x
					}
					right = x
				}
			}
			if left < 0 {
				t.Fatal("missing underline")
			}
			return left, right
		}
		start, _ := underline()
		before, _ := tt.Find("Appearance")
		selected = 1
		tt.Frame()
		if x, _ := underline(); x < start || x >= int(before.X) {
			t.Fatal("underline jumped on selection")
		}
		time.Sleep(60 * time.Millisecond)
		tt.Frame()
		middle, _ := underline()
		if middle <= start || middle >= int(before.X) {
			t.Fatalf("not moving between tabs: %d", middle)
		}
		time.Sleep(tabsIndicatorDuration)
		tt.Frame()
		end, _ := underline()
		if end <= middle {
			t.Fatal("underline did not reach selected tab")
		}
		after, _ := tt.Find("Appearance")
		if before != after {
			t.Fatalf("selection changed text geometry: %v -> %v", before, after)
		}
		selected = 0
		tt.Frame()
		time.Sleep(40 * time.Millisecond)
		tt.Frame()
		selected = 1
		tt.Frame()
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		if x, _ := underline(); x != end {
			t.Fatal("reduced motion did not finish transition")
		}
		selected = 0
		tt.Frame()
		if x, _ := underline(); x != start {
			t.Fatal("reduced motion animated selection")
		}
	}
}

func TestTabsHoverDoesNotChangeAppearance(t *testing.T) {
	for _, mode := range []bool{false, true} {
		selected := 0
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(mode))
			Tabs(c, &selected, "General", "Appearance")
		}, 400, 100)
		before := bytes.Clone(tt.Image().Pix)
		for _, label := range []string{"General", "Appearance"} {
			r, ok := tt.Find(label)
			if !ok {
				t.Fatalf("missing %s", label)
			}
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			if !bytes.Equal(before, tt.Image().Pix) {
				t.Fatalf("hover changed %s", label)
			}
		}
	}
}

func TestTabsEmpty(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) { Tabs(c, &selected) }, 400, 100)
	tt.Frame()
}

func TestTabsKeyboard(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.LightPalette())
		Tabs(c, &selected, "General", "Appearance", "Advanced")
	}, 400, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if selected != 1 {
		t.Errorf("Right chose %d, want 1", selected)
	}
	tt.Key(0, ui.KeyLeft)
	if selected != 0 {
		t.Errorf("Left chose %d, want 0", selected)
	}
}
