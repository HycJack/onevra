package gallery

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestMenuPreview(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, variant := range []string{"Left-click", "Right-click"} {
			s := &controlState{}
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				ui.Column(c).Children(func() { controlMenu(c, s, variant) })
			}, 280, 240)
			var err error
			if variant == "Right-click" {
				tt.Click("Right-click here")
				if len(tt.Menu()) != 0 {
					t.Fatal("left click opened context menu")
				}
				err = tt.RightClick("Right-click here")
			} else {
				err = tt.Click("Actions")
			}
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Contains(tt.Menu(), "Duplicate") {
				t.Fatal(tt.Menu())
			}
			if err := tt.ChooseMenuItem("Duplicate"); err != nil {
				t.Fatal(err)
			}
			if s.index != 1 || s.text != "Duplicate" {
				t.Fatal("action not executed")
			}
			if variant == "Right-click" {
				tt.Key(ui.Shift, ui.KeyF10)
			} else {
				tt.Key(0, ui.KeyDown)
			}
			if err := tt.ChooseMenuItem("Delete"); err != nil {
				t.Fatal(err)
			}
			if s.index != 2 || s.text != "Delete" {
				t.Fatal("keyboard action not executed")
			}
		}
	}
}

func TestNativeSelectPreview(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, width := range []int{540, 280} {
			g := &Gallery{}
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				for _, example := range controlCatalog {
					if example.name == "Select" {
						g.Controls(c, example)
					}
				}
			}, width, 700)
			if !tt.HasText("NativeSelect") {
				t.Fatal("missing native select")
			}
			nativeHeading, _ := tt.Find("NativeSelect")
			for _, label := range []string{"Default", "Selected", "Searchable"} {
				variant, ok := tt.Find(label)
				if !ok || variant.Y+variant.H >= nativeHeading.Y {
					t.Fatalf("%s is not above the separate NativeSelect section", label)
				}
			}
			if err := tt.Click("Choose size"); err != nil {
				t.Fatal(err)
			}
			if err := tt.ChooseMenuItem("Large"); err != nil {
				t.Fatal(err)
			}
			if g.controlState("Select", "NativeSelect").choice != "Large" || g.controlState("Select", "Default").choice != "Small" {
				t.Fatal("selection state leaked")
			}
			r, ok := tt.Find("Choose size")
			if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
				t.Fatal("select overflows")
			}
		}
	}
}

func TestMenuLanguageSwitch(t *testing.T) {
	for _, variant := range []string{"Left-click", "Right-click"} {
		g := &Gallery{}
		g.chooseLanguage("zh-Hans")
		state := &controlState{}
		tt := ui.NewTester(func(c *ui.Context) {
			g.withLanguage(c, func() { ui.Column(c).Children(func() { controlMenu(c, state, variant) }) })
		}, 320, 240)
		var err error
		if variant == "Right-click" {
			err = tt.RightClick(g.text("rightClickHere"))
		} else {
			err = tt.Click(g.text("actions"))
		}
		if err != nil {
			t.Fatal(err)
		}
		if err := tt.ChooseMenuItem("创建副本"); err != nil {
			t.Fatal(err)
		}
		if state.text != "Duplicate" || state.index != 1 {
			t.Fatal("translated menu changed action identity")
		}
		if !tt.HasText("创建副本: 1") {
			t.Fatal(tt.Texts())
		}
		g.chooseLanguage("en")
		tt.Frame()
		if !tt.HasText("Duplicate: 1") {
			t.Fatal(tt.Texts())
		}
	}
}
