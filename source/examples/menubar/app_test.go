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
	path := filepath.Join(t.TempDir(), "app", "config.yaml")
	got, err := loadConfig(path)
	if err != nil || got != defaultConfig() {
		t.Fatalf("missing config: %+v, %v", got, err)
	}
	want := defaultConfig()
	want.Language, want.Appearance, want.Panel, want.Welcomed, want.Silent = "zh-Hans", "dark", false, true, true
	if err := saveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	got, err = loadConfig(path)
	if err != nil || got != want {
		t.Fatalf("round trip: %+v, %v", got, err)
	}
	if err := os.WriteFile(path, []byte("language: future\nappearance: unknown\ncloseAction: invisible\n"), 0600); err != nil {
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
	a := &application{config: want, path: filepath.Join(path, "impossible.yaml")}
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
				case "hybrid":
					want = action != "quit" && tray
				}
				if got := keepRunning(c, platform, tray); got != want {
					t.Errorf("%s/%s/tray=%v: %v", platform, action, tray, got)
				}
			}
		}
	}
}

func TestPanelToolbar(t *testing.T) {
	a := &application{config: defaultConfig(), note: "Draft"}
	tt := ui.NewTester(a.panelView, 360, 300)
	if !tt.HasText(a.text("note")) {
		t.Fatal(tt.Texts())
	}
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
		{"system", "zh-CN", "zh-Hans"}, {"system", "zh-SG", "zh-Hans"}, {"system", "zh-TW", "en"},
		{"system", "de-DE", "en"}, {"en", "zh-CN", "en"}, {"zh-Hans", "en-US", "zh-Hans"},
	} {
		a.config.Language, a.locale = tc.choice, tc.system
		if got := a.language(); got != tc.want {
			t.Errorf("%+v: %s", tc, got)
		}
	}
	catalogs := map[string]map[string]any{}
	for _, code := range []string{"en", "zh-Hans"} {
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
	if len(catalogs["en"]) != len(catalogs["zh-Hans"]) {
		t.Fatal("catalog sizes differ")
	}
	a.locale = "en-US"
	for key := range catalogs["en"] {
		if strings.HasPrefix(key, "@") {
			continue
		}
		a.config.Language = "zh-Hans"
		if catalogs["zh-Hans"][key] == nil || a.text(key) != catalogs["zh-Hans"][key] {
			t.Errorf("missing zh-Hans/%s", key)
		}
	}
	for _, value := range []string{"system", "en", "zh-Hans", "light", "dark", "dock", "tray", "both", "quit"} {
		if _, ok := catalogs["en"][optionKey(value)]; !ok {
			t.Errorf("no message for setting value %s", value)
		}
	}
	before := a.note
	if appKind != "menubar" {
		a.config.Language = "en"
		tt := ui.NewTester(a.mainView, 760, 560)
		if !tt.HasText(a.text("title")) {
			t.Fatal(tt.Texts())
		}
		a.config.Language = "zh-Hans"
		tt = ui.NewTester(a.mainView, 760, 560)
		if !tt.HasText(a.text("title")) {
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
