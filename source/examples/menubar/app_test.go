package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"
	"go.yaml.in/yaml/v3"
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
	if a.errorKey != "save_error" || a.config != want {
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

func TestLanguages(t *testing.T) {
	for _, tc := range []struct{ choice, system, want string }{
		{"system", "zh-CN", "zh-Hans"}, {"system", "zh-SG", "zh-Hans"}, {"system", "zh-TW", "en"},
		{"system", "de-DE", "en"}, {"en", "zh-CN", "en"}, {"zh-Hans", "en-US", "zh-Hans"},
	} {
		if got := languageCode(tc.choice, tc.system); got != tc.want {
			t.Errorf("%+v: %s", tc, got)
		}
	}
	catalogs := map[string]map[string]struct {
		Other string `yaml:"other"`
	}{}
	for _, code := range []string{"en", "zh-Hans"} {
		data, err := localeFiles.ReadFile("locales/" + code + ".yaml")
		if err != nil {
			t.Fatal(err)
		}
		var catalog map[string]struct {
			Other string `yaml:"other"`
		}
		if err := yaml.Unmarshal(data, &catalog); err != nil {
			t.Fatal(err)
		}
		catalogs[code] = catalog
	}
	if len(catalogs["en"]) != len(catalogs["zh-Hans"]) {
		t.Fatal("catalog sizes differ")
	}
	a := &application{config: defaultConfig(), locale: "en-US", note: "Keep this draft"}
	for key := range catalogs["en"] {
		for _, code := range []string{"en", "zh-Hans"} {
			a.config.Language = code
			if catalogs[code][key].Other == "" || a.text(key) == key {
				t.Errorf("missing %s/%s", code, key)
			}
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
