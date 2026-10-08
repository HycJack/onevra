package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestWelcome(t *testing.T) {
	for _, width := range []int{360, 600, 1100} {
		for _, dark := range []bool{false, true} {
			g := &Gallery{page: welcomePage}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, theme.Auto)
				g.View(c)
			}, width, 2200)
			tt.SetDark(dark)
			if tt.HasText("Native UI.") || tt.HasText("Made to feel right.") {
				t.Fatal("welcome hero must not return")
			}
			for _, label := range []string{"Small components. Real possibilities.", "Your next idea", "Find your balance", "Preview name", "Preview category", "Preview intensity", "Create preview", "Explore Lucide"} {
				r, ok := tt.Find(label)
				if !ok || r.W <= 0 || r.X < 0 || r.X+r.W > float32(width)+1 {
					t.Fatalf("width %d, dark %v: %s missing or outside content: %v", width, dark, label, r)
				}
			}
			if err := tt.Click("Create preview"); err != nil {
				t.Fatal(err)
			}
			if !g.welcome.created || !tt.HasText("Ready to explore: A new beginning") {
				t.Fatal("preview did not respond")
			}
			if err := tt.Click("Relax"); err != nil {
				t.Fatal(err)
			}
			if g.welcome.mode != 1 || !tt.HasText("36%") || !tt.HasText("Room to breathe") {
				t.Fatal("mode did not switch to its independent value")
			}
			if err := tt.Click("Quiet mode"); err != nil {
				t.Fatal(err)
			}
			if g.welcome.quiet || !tt.HasText("Back in the mix.") {
				t.Fatal("quiet mode did not respond")
			}
			if g.controlStates != nil {
				t.Fatal("welcome samples must not use individual component demo state")
			}
		}
	}
}

func TestWelcomeEditing(t *testing.T) {
	g := &Gallery{page: welcomePage}
	tt := ui.NewTester(g.View, 600, 1800)
	if err := tt.Click("Preview name"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("A thoughtful app")
	if err := tt.Click("Create preview"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Ready to explore: A thoughtful app") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Include a little inspiration"); err != nil {
		t.Fatal(err)
	}
	if g.welcome.include || g.welcome.created {
		t.Fatal("editing a preview should invalidate the previous result")
	}
	if err := tt.Click("Preview intensity"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	if g.welcome.intensity[0] != 100 || !tt.HasText("100%") {
		t.Fatal("slider did not update the live indicator")
	}
}

func TestWelcomeNavigation(t *testing.T) {
	g := &Gallery{page: componentStatusPage}
	tt := ui.NewTester(g.Left, 220, 600)
	if err := tt.Click("Welcome"); err != nil {
		t.Fatal(err)
	}
	if g.page != welcomePage || pageFromKey(pageKey(g.page)) != welcomePage {
		t.Fatal("welcome navigation did not round trip")
	}
	view := ui.NewTester(g.View, 600, 1800)
	if err := view.Click("Explore Lucide"); err != nil {
		t.Fatal(err)
	}
	if g.page != lucidePage {
		t.Fatalf("Explore Lucide navigated to %s", pageKey(g.page))
	}
}
