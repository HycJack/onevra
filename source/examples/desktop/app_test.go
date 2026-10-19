package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPreferences(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app", "config.json")
	got, err := loadConfig(path)
	if err != nil || got != defaultConfig() {
		t.Fatalf("missing config: %+v, %v", got, err)
	}
	want := defaultConfig()
	want.Language, want.Appearance, want.Panel, want.Welcomed = "zh-Hans", "dark", false, true
	if err := saveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	got, err = loadConfig(path)
	if err != nil || got != want {
		t.Fatalf("round trip: %+v, %v", got, err)
	}
	if err := os.WriteFile(path, []byte(`{"language": "future", "appearance": "unknown", "closeAction": "invisible"}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err = loadConfig(path)
	if err != nil || got != defaultConfig() {
		t.Fatalf("invalid preferences: %+v, %v", got, err)
	}
	if err := os.WriteFile(path, []byte("[broken"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err = loadConfig(path)
	if err == nil || got != defaultConfig() {
		t.Fatalf("corrupt config: %+v, %v", got, err)
	}
	// A save error must remain visible; it must not reset the user's live choice.
	a := &application{config: want, path: filepath.Join(path, "impossible.json")}
	a.changed()
	if a.errorKey != "saveError" || a.config != want {
		t.Fatal("save failure was swallowed or reset preferences")
	}
}

func TestClosePolicy(t *testing.T) {
	for _, platform := range []string{"darwin", "windows", "linux"} {
		for _, action := range []string{"dock", "tray", "both", "quit"} {
			for _, tray := range []bool{false, true} {
				c := defaultConfig()
				c.CloseAction = action
				want := false
				switch appKind {
				case "desktop":
					want = platform == "darwin" && action == "dock"
				case "menubar":
					want = tray
				}
				if got := keepRunning(c, platform, tray); got != want {
					t.Errorf("%s/%s/tray=%v: %v", platform, action, tray, got)
				}
			}
		}
	}
}

func TestMainTabs(t *testing.T) {
	if appKind == "menubar" {
		t.Skip("the menu bar example has no main window")
	}
	a := &application{config: defaultConfig(), note: "Draft"}
	tt := ui.NewTester(a.headerLeading, 400, 60)
	if !tt.HasText(a.text("note")) || tt.HasText(a.text("description")) {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click(a.text("about")); err != nil {
		t.Fatal(err)
	}
	if a.mainTab != 1 {
		t.Fatalf("tab=%d", a.mainTab)
	}
	tt = ui.NewTester(a.mainView, 760, 560)
	if !tt.HasText(a.text("description")) || a.note != "Draft" {
		t.Fatalf("note=%q %v", a.note, tt.Texts())
	}
}

func TestPanelToolbar(t *testing.T) {
	if appKind == "menubar" {
		t.Skip("the menu bar example tests its panel separately")
	}
	a := &application{config: defaultConfig(), note: "Draft"}
	tt := ui.NewTester(a.panelView, 360, 300)
	for _, label := range []string{a.text("copy"), a.text("settings")} {
		if _, ok := tt.Find(label); !ok {
			t.Fatalf("missing %s action: %v", label, tt.Texts())
		}
	}
	if err := tt.Click(a.text("about")); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if a.panelTab != 1 || !tt.HasText(a.text("title")) || a.note != "Draft" {
		t.Fatalf("tab=%d note=%q", a.panelTab, a.note)
	}
	if err := tt.Click(a.text("copy")); err != nil {
		t.Fatal(err)
	}
	if got := tt.Clipboard(); got != "Draft" {
		t.Fatalf("clipboard=%q", got)
	}
}

func TestLanguages(t *testing.T) {
	a := &application{config: defaultConfig(), note: "Keep this draft"}
	for _, tc := range []struct{ choice, system, want string }{
		{"system", "zh-CN", "zh-Hans"}, {"system", "zh-SG", "zh-Hans"}, {"system", "zh-TW", "zh-Hant"},
		{"system", "ja-JP", "ja"}, {"system", "ko-KR", "ko"}, {"system", "de-DE", "en"}, {"en", "zh-CN", "en"}, {"zh-Hans", "en-US", "zh-Hans"},
	} {
		a.config.Language, a.preferred = tc.choice, []string{tc.system}
		if got := a.language(); got != tc.want {
			t.Errorf("%+v: %s", tc, got)
		}
	}
	a.config.Language, a.preferred = "system", []string{"de-DE", "ja-JP"}
	if got := a.language(); got != "ja" {
		t.Errorf("a later preferred language: %s", got)
	}
	catalogs := map[string]map[string]any{}
	for _, code := range languages {
		data, err := os.ReadFile(filepath.Join("locales", code+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var catalog map[string]any
		if err := json.Unmarshal(data, &catalog); err != nil {
			t.Fatal(err)
		}
		catalogs[code] = catalog
	}
	for _, code := range languages {
		if len(catalogs[code]) != len(catalogs["en"]) {
			t.Fatalf("%s catalog size differs", code)
		}
	}
	a.preferred = []string{"en-US"}
	for key := range catalogs["en"] {
		if strings.HasPrefix(key, "@") {
			continue
		}
		for _, code := range languages {
			a.config.Language = code
			if catalogs[code][key] == nil || a.text(key) != catalogs[code][key] {
				t.Errorf("missing %s/%s", code, key)
			}
		}
	}
	for _, value := range append([]string{"system", "light", "dark", "dock", "tray", "both", "quit"}, languages...) {
		if _, ok := catalogs["en"][optionKey(value)]; !ok {
			t.Errorf("no message for setting value %s", value)
		}
	}
	before := a.note
	if appKind != "menubar" {
		a.config.Language = "en"
		tt := ui.NewTester(a.mainView, 760, 560)
		if !tt.HasText(a.text("note")) {
			t.Fatal(tt.Texts())
		}
		a.config.Language = "zh-Hans"
		tt = ui.NewTester(a.mainView, 760, 560)
		if !tt.HasText(a.text("note")) {
			t.Fatal(tt.Texts())
		}
	}
	if a.note != before {
		t.Fatal("language change lost note")
	}
	for _, code := range []string{"en", "zh-Hans"} {
		a.config.Language = code
		tt := ui.NewTester(func(c *ui.Context) { a.scope(c, func() { a.settingsView(c, 360) }) }, 360, 540)
		if !tt.HasText(a.text("language")) || !tt.HasText(a.text("appearance")) {
			t.Fatal(tt.Texts())
		}
	}
	if a.settingsState.Query != "" {
		t.Fatal("render mutated search")
	}
}
