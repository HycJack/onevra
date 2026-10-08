package gallery

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestButtonPreview(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 280} {
			g := &Gallery{}
			example, _ := controlExampleFor(pageFromKey("controls/Button"))
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

func TestButtonTogglePreview(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 280} {
			g := &Gallery{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				example, _ := controlExampleFor(pageFromKey("controls/Button"))
				g.Controls(c, example)
			}, width, 1000)
			if !tt.HasText("ToggleButton") || tt.HasText("Disabled") {
				t.Fatal("missing toggle section or unexpected disabled preview")
			}
			s := g.controlState("Button", "ToggleButton")
			if s.on || !s.toggles[1] {
				t.Fatal("preview must demonstrate off and on states")
			}
			for _, label := range []string{"Bold", "Italic"} {
				r, ok := tt.Find(label)
				if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
					t.Fatalf("%s overflows preview: %+v", label, r)
				}
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			if !s.on || s.toggles[1] {
				t.Fatal("toggle preview did not update independent values")
			}
			if g.controlState("Button", "Button").button.clicks != [5]int{} {
				t.Fatal("toggle state leaked into other button examples")
			}
			if pageFromKey("controls/ToggleButton") != lucidePage {
				t.Fatal("standalone toggle route remains")
			}
		}
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
			var previous ui.Rect
			for i, label := range []string{"Primary", "Secondary", "Outline", "Ghost", "Link"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing button %s", label)
				}
				if r.X < 0 || r.X+r.W > float32(width)+1 {
					t.Fatalf("%s overflows width %d: %+v", label, width, r)
				}
				if i > 0 {
					if width == 540 && r.Y+r.H/2 != previous.Y+previous.H/2 {
						t.Fatalf("%s is not on the same row", label)
					}
					if r.Y < previous.Y+previous.H && r.X < previous.X+previous.W {
						t.Fatalf("%s overlaps the previous button", label)
					}
				}
				previous = r
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
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
