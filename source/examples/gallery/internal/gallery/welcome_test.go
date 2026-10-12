package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestWelcomeStartupPage(t *testing.T) {
	for _, width := range []int{360, 820, 1100} {
		for _, dark := range []bool{false, true} {
			g := &Gallery{page: componentStatusPage}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, theme.Auto)
				g.WelcomePage(c)
			}, width, 600)
			tt.SetDark(dark)
			if g.page != componentStatusPage || !tt.HasText("Welcome to Onevra") {
				t.Fatal("startup must show Welcome without replacing the saved main route")
			}
			if _, ok := tt.Find("Search components"); ok {
				t.Fatal("startup must not show main navigation")
			}
			r, ok := tt.Find("Enter gallery")
			if !ok || r.Y < 48 || r.X < 0 || r.X+r.W > float32(width) || r.Y+r.H > 600 {
				t.Fatalf("entry outside the safe visible area: %v", r)
			}
			if err := tt.Click("Enter gallery"); err != nil {
				t.Fatal(err)
			}
			if g.page != componentStatusPage {
				t.Fatal("entry must retain the saved main route")
			}
		}
	}
}

func TestWelcome(t *testing.T) {
	for _, width := range []int{360, 600, 1100} {
		for _, dark := range []bool{false, true} {
			g := &Gallery{page: welcomePage}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, theme.Auto)
				g.View(c)
			}, width, 1600)
			tt.SetDark(dark)
			for _, label := range []string{"Onevra", "Small components. Real possibilities.", repositoryLabel, "Search everything", "Version " + galleryVersion} {
				r, ok := tt.Find(label)
				if !ok || r.W <= 0 || r.X < 0 || r.X+r.W > float32(width)+1 {
					t.Fatalf("width %d, dark %v: %s missing or outside content: %v", width, dark, label, r)
				}
			}
			for _, removed := range []string{"Your next idea", "Browse the collection", "Explore Lucide", "Quiet mode"} {
				if tt.HasText(removed) {
					t.Fatalf("welcome must stay focused; found %q", removed)
				}
			}
		}
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
}

func TestWelcomeLinks(t *testing.T) {
	var opened []string
	g := &Gallery{page: welcomePage, openURL: func(url string) error {
		opened = append(opened, url)
		return nil
	}}
	tt := ui.NewTester(g.View, 600, 1200)
	if err := tt.Click(repositoryLabel); err != nil {
		t.Fatal(err)
	}
	if len(opened) != 1 || opened[0] != repositoryURL {
		t.Fatalf("repository opened %v", opened)
	}
	if err := tt.Click("Search everything"); err != nil {
		t.Fatal(err)
	}
	if !g.commands.open {
		t.Fatal("search did not open the command palette")
	}
	headless := &Gallery{page: welcomePage}
	if err := ui.NewTester(headless.View, 600, 1200).Click(repositoryLabel); err != nil {
		t.Fatal(err)
	}
}
