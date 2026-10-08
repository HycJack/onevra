package gallery

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestButtonPreview(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 280} {
			g := &Gallery{}
			example := controlCatalog[1]
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				g.Controls(c, example)
			}, width, 600)
			if _, ok := tt.Find("Button"); !ok || tt.HasText("Default") {
				t.Fatal("button section must use the component name")
			}
			for _, label := range []string{"Primary", "Secondary", "Outline", "Ghost", "Link"} {
				r, ok := tt.Find(label)
				if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
					t.Fatalf("%s overflows preview: %+v", label, r)
				}
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			if g.controlState("Button", "Button").button.clicks != [5]int{1, 1, 1, 1, 1} {
				t.Fatal("preview lost button interactions")
			}
		}
	}
}

func TestButtonMenu(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 280} {
			g := &Gallery{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				g.Controls(c, controlCatalog[1])
			}, width, 600)
			if !tt.HasText("MenuButton") {
				t.Fatal("missing menu button section")
			}
			for _, size := range []string{"Large", "Medium", "Small"} {
				if err := tt.Click("Choose size"); err != nil {
					t.Fatal(err)
				}
				if !slices.Equal(tt.Menu(), []string{"Small", "Medium", "Large"}) {
					t.Fatalf("unexpected menu: %v", tt.Menu())
				}
				if err := tt.ChooseMenuItem(size); err != nil {
					t.Fatal(err)
				}
				if !tt.HasText("Size: "+size) || g.controlState("Button", "MenuButton").choice != size {
					t.Fatal("menu selection did not update button")
				}
				r, ok := tt.Find("Choose size")
				if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
					t.Fatalf("menu button overflows preview: %+v", r)
				}
				first, found := tt.Find("Primary")
				if !found || r.X != first.X || r.W >= float32(width)/2 {
					t.Fatalf("menu button must be content-sized and left-aligned with Button: menu=%+v, first=%+v", r, first)
				}
			}
			if g.controlState("Button", "Button").button.clicks != [5]int{} {
				t.Fatal("menu changed ordinary button counts")
			}
			// Selection leaves focus on the trigger; Down reopens its menu.
			tt.Key(0, ui.KeyDown)
			if len(tt.Menu()) != 3 {
				t.Fatal("keyboard did not reopen menu")
			}
			tt.CloseMenu()
		}
	}
}

func TestButtonMenuDisabled(t *testing.T) {
	s := &controlState{choice: "Small"}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() { controlButton(c, s, "MenuButton") })
	}, 280, 200)
	if err := tt.Click("Choose size"); err != nil {
		t.Fatal(err)
	}
	if len(tt.Menu()) != 0 || s.choice != "Small" {
		t.Fatal("disabled menu button responded to input")
	}
}

func TestControlButtonWrap(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 200} {
			s := &controlState{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Fill().Gap(12).Children(func() { controlButton(c, s, "Default") })
			}, width, 300)
			var first, previous ui.Rect
			wrapped := false
			for i, label := range []string{"Primary", "Secondary", "Outline", "Ghost", "Link"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing button %s", label)
				}
				if r.X < 0 || r.X+r.W > float32(width)+1 {
					t.Fatalf("%s overflows width %d: %+v", label, width, r)
				}
				if i == 0 {
					first = r
				} else {
					if r.Y > previous.Y+1 {
						wrapped = true
					} else if r.X < previous.X+previous.W {
						t.Fatalf("%s overlaps the previous button", label)
					}
					if width == 540 && r.Y+r.H/2 != first.Y+first.H/2 {
						t.Fatal("wide layout should keep buttons on one row")
					}
				}
				previous = r
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			if width == 200 && !wrapped {
				t.Fatal("narrow layout should wrap buttons")
			}
			if s.button.clicks != [5]int{1, 1, 1, 1, 1} {
				t.Fatalf("unexpected per-button clicks: %v", s.button.clicks)
			}
		}
	}
}

func TestControlButtonPointerStyles(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, disabled := range []bool{false, true} {
			s := &controlState{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Fill().Disabled(disabled).Children(func() { controlButton(c, s, "Default") })
			}, 540, 300)
			for _, label := range []string{"Primary", "Secondary", "Outline", "Ghost"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing button %s", label)
				}
				x, y := r.X+6, r.Y+r.H/2
				tt.Move(530, 290)
				rest := tt.Image().RGBAAt(int(x), int(y))
				tt.Move(x, y)
				hover := tt.Image().RGBAAt(int(x), int(y))
				tt.Press(x, y)
				pressed := tt.Image().RGBAAt(int(x), int(y))
				tt.Release(x, y)
				tt.Move(530, 290)
				reset := tt.Image().RGBAAt(int(x), int(y))
				if disabled {
					if hover != rest || pressed != rest {
						t.Fatalf("%v/%s: disabled button changed on pointer input", mode, label)
					}
				} else if hover == rest || pressed == hover {
					t.Fatalf("%v/%s: missing hover or pressed feedback", mode, label)
				}
				if reset != rest {
					t.Fatalf("%v/%s: pointer exit did not restore resting style", mode, label)
				}
			}
		}
	}
}

func TestButtonIndependentCounts(t *testing.T) {
	s := &controlState{}
	tt := ui.NewTester(func(c *ui.Context) {
		controlButton(c, s, "Button")
	}, 540, 300)
	labels := []string{"Primary", "Secondary", "Outline", "Ghost", "Link"}
	for i, label := range labels {
		for range i + 1 {
			if err := tt.Click(label); err != nil {
				t.Fatal(err)
			}
		}
		for j, name := range labels {
			want := 0
			if j <= i {
				want = j + 1
			}
			if s.button.clicks[j] != want {
				t.Fatalf("%s count = %d, want %d", name, s.button.clicks[j], want)
			}
			// MyGo exposes Link's inline text through its accessible label,
			// not as a separate text rectangle.
			if name == "Link" {
				continue
			}
			content, ok := tt.Find(fmt.Sprintf("%s: %d", name, want))
			button, found := tt.Find(name)
			if !ok || !found || content.X < button.X || content.Y < button.Y || content.X+content.W > button.X+button.W || content.Y+content.H > button.Y+button.H {
				t.Fatalf("%s count is not inside its button: content=%+v (%t), button=%+v (%t)", name, content, ok, button, found)
			}
		}
	}
	if tt.HasText("Clicks:") {
		t.Fatal("unexpected shared click counter")
	}
}
