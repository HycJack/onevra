package components

import (
	"reflect"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func settingsFixture() []SettingsCategory {
	return []SettingsCategory{
		{ID: "account", Title: "Account", Section: "Personal", Groups: []SettingsGroup{{ID: "security", Title: "Security", Items: []SettingsItem{{ID: "alerts", Title: "Alerts", Description: "New sign-in", Keywords: []string{"邮件"}}}}}},
		{ID: "appearance", Title: "Appearance", Groups: []SettingsGroup{{ID: "colors", Title: "Colors", Items: []SettingsItem{{ID: "theme", Title: "Theme", Description: "Choose a mode", Keywords: []string{"dark", "深色"}}}}}},
	}
}

func TestSearchSettings(t *testing.T) {
	catalog := settingsFixture()
	for _, tc := range []struct {
		query string
		want  []SettingsResult
	}{
		{"", []SettingsResult{{0, 0, 0}, {1, 0, 0}}},
		{" \t ", []SettingsResult{{0, 0, 0}, {1, 0, 0}}},
		{" DARK appearance ", []SettingsResult{{1, 0, 0}}},
		{"personal SECURITY sign-in", []SettingsResult{{0, 0, 0}}},
		{"深色", []SettingsResult{{1, 0, 0}}},
		{"邮件", []SettingsResult{{0, 0, 0}}},
		{"dark alerts", nil},
		{"missing", nil},
	} {
		if got := SearchSettings(catalog, tc.query); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.query, got, tc.want)
		}
	}
	if got := SearchSettings(nil, ""); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestSettingsStateAndBuilders(t *testing.T) {
	for _, narrow := range []bool{false, true} {
		state := SettingsState{Category: "obsolete"}
		catalog := settingsFixture()
		enabled := false
		catalog[1].Groups[0].Items[0].Build = func(c *ui.Context) { Switch(c, &enabled, "Enable theme") }
		tt := ui.NewTester(func(c *ui.Context) { Settings(c, &state, catalog, SettingsOptions{Narrow: narrow}) }, 800, 600)
		tt.Frame()
		if state.Category != "account" || !tt.HasText("Alerts") || tt.HasText("Theme") {
			t.Fatal("initial selection", tt.Texts())
		}
		state.Query = "dark"
		tt.Frame()
		if !tt.HasText("Theme") || tt.HasText("Alerts") || state.Category != "account" {
			t.Fatal("global search", tt.Texts())
		}
		if err := tt.Click("Enable theme"); err != nil {
			t.Fatal(err)
		}
		if !enabled {
			t.Fatal("result control did not update app state")
		}
		state.Query = ""
		tt.Frame()
		if !tt.HasText("Alerts") || tt.HasText("Theme") {
			t.Fatal("search did not restore category")
		}
		state.Query = "unmatched"
		tt.Frame()
		if !tt.HasText("No matching settings. Try another search.") {
			t.Fatal("empty result")
		}
		if narrow {
			if err := tt.Click("Settings category"); err != nil {
				t.Fatal(err)
			}
		}
		if err := tt.Click("Appearance"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if state.Query != "" || state.Category != "appearance" || !tt.HasText("Theme") {
			t.Fatal("category navigation", state, tt.Texts())
		}
		catalog = nil
		tt.Frame()
		if !tt.HasText("No matching settings. Try another search.") {
			t.Fatal("empty catalog")
		}
	}
}

func TestSettingsSidebarSpacing(t *testing.T) {
	catalog := []SettingsCategory{
		{ID: "account", Title: "Account", Section: "Personal"},
		{ID: "notifications", Title: "Notifications", Section: "Personal"},
		{ID: "general", Title: "General", Section: "Application"},
		{ID: "appearance", Title: "Appearance", Section: "Application"},
		{ID: "about", Title: "About"},
		{ID: "help", Title: "Help"},
		{ID: "privacy", Title: "Privacy", Section: "Personal"},
	}
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		state := SettingsState{}
		var itemGap, sectionGap float32
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			itemGap, sectionGap = c.Theme().Space(1), c.Theme().Space(3)
			Settings(c, &state, catalog, SettingsOptions{})
		}, 800, 600)
		find := func(label string) ui.Rect {
			t.Helper()
			r, ok := tt.Find(label)
			if !ok {
				t.Fatalf("missing navigation element %q", label)
			}
			return r
		}
		for _, pair := range [][2]string{{"Account", "Notifications"}, {"General", "Appearance"}, {"About", "Help"}} {
			first, second := find(pair[0]), find(pair[1])
			if gap := second.Y - first.Y - first.H; gap != itemGap {
				t.Fatalf("%v: item gap = %v, want %v", pair, gap, itemGap)
			}
		}
		notifications, application := find("Notifications"), find("Application")
		appearance, about := find("Appearance"), find("About")
		for _, pair := range [][2]ui.Rect{{notifications, application}, {appearance, about}} {
			if gap := pair[1].Y - pair[0].Y - pair[0].H; gap != sectionGap {
				t.Fatalf("section gap = %v, want %v", gap, sectionGap)
			}
		}
		// Repeated and unheaded sections must keep every category actionable.
		for _, category := range catalog {
			state.Query = "unmatched"
			tt.Frame()
			if err := tt.Click(category.Title); err != nil {
				t.Fatal(err)
			}
			if state.Category != category.ID || state.Query != "" {
				t.Fatalf("navigation to %q: %+v", category.ID, state)
			}
		}
	}
}

func TestSettingsLocalizedAndIndependent(t *testing.T) {
	a, b := SettingsState{Query: "not found"}, SettingsState{}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() {
			ui.Column(c).Key("a").Fill().Children(func() {
				Settings(c, &a, settingsFixture(), SettingsOptions{Narrow: true, Title: "偏好设置", ResultsTitle: "搜索结果", SearchLabel: "搜索设置", EmptyText: "无结果"})
			})
			ui.Column(c).Key("b").Fill().Children(func() { Settings(c, &b, settingsFixture(), SettingsOptions{Narrow: true}) })
		})
	}, 900, 600)
	tt.Frame()
	if !tt.HasText("搜索结果") || !tt.HasText("无结果") || !tt.HasText("Alerts") || b.Query != "" {
		t.Fatal(tt.Texts())
	}
}

func TestSettingsGroupActions(t *testing.T) {
	for _, narrow := range []bool{false, true} {
		state := SettingsState{}
		catalog := settingsFixture()
		clicks := 0
		catalog[0].Groups[0].Actions = func(c *ui.Context) {
			if Button(c, "Refresh", ButtonSecondary).Clicked() {
				clicks++
			}
		}
		tt := ui.NewTester(func(c *ui.Context) { Settings(c, &state, catalog, SettingsOptions{Narrow: narrow}) }, 800, 600)
		tt.Frame()
		title, ok := tt.Find("Security")
		action, found := tt.Find("Refresh")
		if !ok || !found || action.X <= title.X || action.Y > title.Y+title.H || action.Y+action.H < title.Y {
			t.Fatalf("narrow %v: action not trailing the group title: %+v %+v", narrow, title, action)
		}
		if err := tt.Click("Refresh"); err != nil || clicks != 1 {
			t.Fatalf("narrow %v: action not clicked: %v %d", narrow, err, clicks)
		}
		// Actions follow their group into search results but are not searched.
		state.Query = "Refresh"
		tt.Frame()
		if tt.HasText("Refresh") || tt.HasText("Alerts") {
			t.Fatal("actions were searched", tt.Texts())
		}
		state.Query = "sign-in"
		tt.Frame()
		if !tt.HasText("Refresh") {
			t.Fatal("actions missing from search results", tt.Texts())
		}
		state.Query = ""
		state.Category = "appearance"
		tt.Frame()
		if tt.HasText("Refresh") {
			t.Fatal("actions shown outside their category", tt.Texts())
		}
	}
}
