package components

import (
	"reflect"
	"testing"

	"github.com/egoist/mygo/ui"
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
