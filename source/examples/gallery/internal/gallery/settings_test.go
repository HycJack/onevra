package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSettingsDialog(t *testing.T) {
	g := &Gallery{page: shellPage}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			g.LeftToolbar(c)
			g.View(c)
		})
	}, 820, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	if !g.settings.open || !tt.HasText("Smart notifications") {
		t.Fatal("settings button did not open the dialog", tt.Texts())
	}
	if err := tt.Click("Toggle smart notifications"); err != nil {
		t.Fatal(err)
	}
	if g.settings.smart {
		t.Fatal("switch did not update template state")
	}
	g.settings.view.Query = "New sign-in"
	tt.Frame()
	if err := tt.Click("New sign-in delivery"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Off"); err != nil {
		t.Fatal(err)
	}
	if g.settings.channels[1] != "Off" || !g.settings.open {
		t.Fatal("delivery selection did not stay inside the dialog")
	}
	g.settings.view.Query = ""
	tt.Frame()
	if err := tt.Click("Appearance"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Theme") {
		t.Fatal("appearance page missing")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if g.settings.open || g.page != shellPage || !tt.HasText("Clicks: 0") {
		t.Fatal("dismissal changed the underlying page")
	}
	for _, close := range []string{"Close settings", "backdrop"} {
		if err := tt.Click("Settings"); err != nil {
			t.Fatal(err)
		}
		if close == "backdrop" {
			tt.ClickAt(1, 100)
		} else if err := tt.Click(close); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if g.settings.open || g.settings.smart {
			t.Fatal("close or session-state retention failed")
		}
	}
}

func TestSettingsWindowSizes(t *testing.T) {
	for _, size := range [][2]int{{360, 480}, {820, 600}, {820, 900}, {1200, 900}} {
		g := &Gallery{}
		g.settings.open = true
		tt := ui.NewTester(g.Settings, size[0], size[1])
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		tt.Frame()
		panel, ok := tt.Find("Settings")
		w, h := float32(size[0]), float32(size[1])
		if !ok || panel.X < 47.5 || panel.Y < 47.5 || panel.X+panel.W > w-47.5 || panel.Y+panel.H > h-47.5 {
			t.Fatalf("%v: settings escaped viewport clearance: %+v", size, panel)
		}
		if panel.W > min(800, w*.88)+1 {
			t.Fatalf("%v: settings width exceeded compact limit: %+v", size, panel)
		}
		if panel.H < min(720, h*.8, h-96)-1 || panel.Y+panel.H/2 > h/2+.5 {
			t.Fatalf("%v: settings too short or not biased upward: %+v", size, panel)
		}
		close, ok := tt.Find("Close settings")
		if !ok || close.Y < panel.Y || close.Y+close.H > panel.Y+60 || close.X < panel.X+panel.W-60 || close.X+close.W > panel.X+panel.W {
			t.Fatalf("%v: close button not inside the top-right corner: %+v", size, close)
		}
	}
}

func TestSettingsSearchAndGeometry(t *testing.T) {
	for _, width := range []int{360, 820} {
		g := &Gallery{}
		g.settings.open = true
		tt := ui.NewTester(g.Settings, width, 600)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		g.settings.view.Query = "NEW SIGN-IN"
		tt.Frame()
		if !tt.HasText("New sign-in") || tt.HasText("Password changed") {
			t.Fatal("search did not filter settings", tt.Texts())
		}
		if width == 820 {
			search, _ := tt.Find("Search settings")
			nav, _ := tt.Find("Notifications")
			if search.X+search.W > nav.X+nav.W+1 {
				t.Fatalf("search extends beyond sidebar: %+v, nav %+v", search, nav)
			}
		}
		for _, label := range []string{"Close settings", "Search settings", "New sign-in delivery"} {
			r, ok := tt.Find(label)
			if !ok || r.W <= 0 || r.X < 0 || r.X+r.W > float32(width) || r.Y < 0 || r.Y+r.H > 600 {
				t.Fatalf("%d: %s bounds = %+v, found=%v", width, label, r, ok)
			}
		}
		g.settings.view.Query = "no such setting"
		tt.Frame()
		if !tt.HasText("No matching settings. Try another search.") {
			t.Fatal("empty search feedback missing")
		}
	}
}
