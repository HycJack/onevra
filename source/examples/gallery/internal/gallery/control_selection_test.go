package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestSelectionNavigation(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 1800)
	for _, old := range []string{"Check boxes", "Switches", "Radio buttons"} {
		if tt.HasText(old) {
			t.Fatalf("old sidebar entry remains: %s", old)
		}
	}
	if err := tt.Click("Selection"); err != nil {
		t.Fatal(err)
	}
	if pageKey(g.page) != "controls/Selection" {
		t.Fatal("incorrect selection route")
	}
	for _, name := range []string{"Checkbox", "Switch", "Radio", "Selection"} {
		if pageFromKey("controls/"+name) != g.page {
			t.Fatalf("saved route did not restore Selection: %s", name)
		}
	}
}

func TestSelectionInteractions(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{820, 280} {
			g := &Gallery{page: pageFromKey("controls/Selection")}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				g.View(c)
			}, width, 1000)
			for _, name := range []string{"Checkbox", "Switch", "Radio"} {
				if !tt.HasText(name) {
					t.Fatalf("missing section: %s", name)
				}
			}
			for _, text := range []string{"Enabled:", "Sound:", "Selected:", "Disabled"} {
				if tt.HasText(text) {
					t.Fatalf("unexpected status text: %s", text)
				}
			}
			initial := [2]bool{false, true}
			for _, section := range []struct {
				name   string
				labels []string
			}{
				{"Checkbox", []string{"Unchecked", "Checked"}},
				{"Switch", []string{"Off", "On"}},
			} {
				s := g.controlState("Selection", section.name)
				if s.toggles != initial {
					t.Fatalf("%s initial states = %v", section.name, s.toggles)
				}
				var first ui.Rect
				for i, label := range section.labels {
					r, ok := tt.Find(label)
					if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
						t.Fatalf("%s outside preview: %+v", label, r)
					}
					if i == 0 {
						first = r
					} else if width == 820 && (r.X <= first.X || r.Y != first.Y) {
						t.Fatalf("%s is not horizontally arranged", label)
					}
					if err := tt.Click(label); err != nil {
						t.Fatal(err)
					}
					want := initial
					for j := 0; j <= i && j < 2; j++ {
						want[j] = !want[j]
					}
					if s.toggles != want {
						t.Fatalf("%s: got %v, want %v", label, s.toggles, want)
					}
				}
			}
			radio := g.controlState("Selection", "Radio")
			if radio.choice != "Small" {
				t.Fatal("toggle interactions changed radio")
			}
			small, _ := tt.Find("Small")
			large, _ := tt.Find("Large")
			if width == 820 && (large.X <= small.X || large.Y != small.Y) {
				t.Fatal("radio options are not horizontally arranged")
			}
			if err := tt.Click("Large"); err != nil {
				t.Fatal(err)
			}
			if radio.choice != "Large" {
				t.Fatal("radio selection did not update")
			}
			g.page = pageFromKey("controls/Button")
			tt.Frame()
			g.page = pageFromKey("controls/Selection")
			tt.Frame()
			for _, name := range []string{"Checkbox", "Switch"} {
				if g.controlState("Selection", name).toggles != [2]bool{true, false} {
					t.Fatalf("%s state lost after navigation", name)
				}
			}
			if radio.choice != "Large" {
				t.Fatal("radio state lost after navigation")
			}
		}
	}
}

func TestSelectionDisabled(t *testing.T) {
	for _, name := range []string{"Checkbox", "Switch", "Radio"} {
		s := (&Gallery{}).controlState("Selection", name)
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Fill().Disabled(true).Children(func() { controlSelection(c, s, name) })
		}, 820, 200)
		label := map[string]string{"Checkbox": "Unchecked", "Switch": "Off", "Radio": "Large"}[name]
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
		if s.toggles != [2]bool{false, true} || s.choice != "Small" {
			t.Fatalf("%s responded while disabled", name)
		}
	}
}
