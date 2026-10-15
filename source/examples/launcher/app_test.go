package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestFilter(t *testing.T) {
	a := &application{config: defaultConfig()}
	for _, name := range []string{"Calendar", "Visual Studio Code", "Activity Monitor", "Code Editor", "Xcode", "Calculator"} {
		a.apps = append(a.apps, mygo.InstalledApp{Name: name, Path: "/" + name})
	}
	for query, want := range map[string][]string{
		"":     {"Calendar", "Visual Studio Code", "Activity Monitor", "Code Editor", "Xcode", "Calculator"},
		"cal":  {"Calendar", "Calculator"},
		"code": {"Code Editor", "Visual Studio Code", "Xcode"},
		"vsc":  {"Visual Studio Code"},
		"ACT":  {"Activity Monitor", "Calculator"}, // its letters in order, last
		"zzz":  nil,
	} {
		a.query = query
		a.filter()
		var got []string
		for _, app := range a.results {
			got = append(got, app.Name)
		}
		if !slices.Equal(got, want) {
			t.Errorf("%q: %v, want %v", query, got, want)
		}
	}
}

func TestPreferences(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app", "config.json")
	got, err := loadConfig(path)
	if err != nil || got != defaultConfig() {
		t.Fatalf("missing config: %+v, %v", got, err)
	}
	want := config{Language: "zh-Hans", Appearance: "dark", Shortcut: "cmdShiftSpace"}
	if err := saveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	got, err = loadConfig(path)
	if err != nil || got != want {
		t.Fatalf("round trip: %+v, %v", got, err)
	}
	if err := os.WriteFile(path, []byte(`{"language": "future", "appearance": "unknown", "shortcut": "Hyper+Q"}`), 0600); err != nil {
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

func TestShortcuts(t *testing.T) {
	if acceleratorFor(defaultShortcut) == "" {
		t.Fatalf("default shortcut %q has no accelerator", defaultShortcut)
	}
	for _, option := range shortcutOptions {
		if acceleratorFor(option.id) != option.accelerator {
			t.Errorf("%s: %q", option.id, acceleratorFor(option.id))
		}
	}
	if acceleratorFor("unknown") != "" {
		t.Fatal("unknown shortcut resolved")
	}
	if got := keys("CmdOrCtrl+Alt+Shift+Space"); !strings.Contains(got, "Space") {
		t.Fatalf("keys() lost the key: %q", got)
	}
}

func TestPanel(t *testing.T) {
	a := &application{config: defaultConfig(), accelerator: "Alt+Shift+Space"}
	a.apps = []mygo.InstalledApp{{Name: "Calculator", ID: "com.apple.calculator", Path: "/Applications/Calculator.app"}}
	a.query = "calc"
	a.filter()
	tt := ui.NewTester(a.panelView, width, height)
	for _, want := range []string{"Calculator", a.text("keyHints")} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q: %v", want, tt.Texts())
		}
	}
	if !tt.HasText(a.format("showOrHide", map[string]any{"shortcut": keys("Alt+Shift+Space")})) {
		t.Fatalf("missing footer: %v", tt.Texts())
	}
	a.query = "zzz"
	a.filter()
	tt = ui.NewTester(a.panelView, width, height)
	if !tt.HasText(a.text("noApps")) {
		t.Fatalf("missing empty state: %v", tt.Texts())
	}
}

func TestLanguages(t *testing.T) {
	a := &application{config: defaultConfig()}
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
	for key := range catalogs["en"] {
		if strings.HasPrefix(key, "@") {
			continue
		}
		if catalogs["zh-Hans"][key] == nil {
			t.Errorf("missing zh-Hans/%s", key)
		}
	}
	for _, value := range []string{"system", "en", "zh-Hans", "light", "dark"} {
		if _, ok := catalogs["en"][optionKey(value)]; !ok {
			t.Errorf("no message for setting value %s", value)
		}
	}
	a.config.Language = "zh-Hans"
	if a.text("noApps") == catalogs["en"]["noApps"] {
		t.Fatal("translation was not applied")
	}
}

func TestInitial(t *testing.T) {
	for in, want := range map[string]string{"Calculator": "C", "  zoom": "Z", "": "?", "1Password": "1"} {
		if got := initial(in); got != want {
			t.Errorf("initial(%q) = %q, want %q", in, got, want)
		}
	}
}
