package gallery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"go.yaml.in/yaml/v3"
)

func TestLanguageResolution(t *testing.T) {
	for _, tc := range []struct{ choice, system, want string }{
		{"", "zh-CN", "zh-Hans"}, {"", "zh-SG", "zh-Hans"}, {"", "zh-Hans", "zh-Hans"},
		{"", "zh-TW", "en"}, {"", "zh-Hant-HK", "en"}, {"", "de-DE", "en"},
		{"", "C", "en"}, {"", "garbage!", "en"}, {"en", "zh-CN", "en"},
		{"zh-Hans", "en-US", "zh-Hans"}, {"unknown", "zh-CN", "zh-Hans"},
	} {
		got := newTranslator(tc.choice, tc.system)
		if got.code != tc.want {
			t.Errorf("%+v: %s", tc, got.code)
		}
		if got.code == "zh-Hans" && (got.text("gallery.clear_glass_style") != "清透" || got.text("Clear") != "清空") {
			t.Fatal("contextual translation")
		}
		if got.text("untranslated fixture") != "untranslated fixture" {
			t.Fatal("missing fallback")
		}
	}
}

func TestLanguageCatalogs(t *testing.T) {
	type message struct {
		Other string `yaml:"other"`
	}
	catalogs := map[string]map[string]message{}
	for _, code := range []string{"en", "zh-Hans"} {
		data, err := localeFiles.ReadFile("locales/" + code + ".yaml")
		if err != nil {
			t.Fatal(err)
		}
		var entries map[string]message
		if err := yaml.Unmarshal(data, &entries); err != nil {
			t.Fatal(err)
		}
		catalogs[code] = entries
	}
	formats := regexp.MustCompile(`%(?:\[[0-9]+\])?[-+ #0]*[0-9]*(?:\.[0-9]+)?[a-zA-Z%]`)
	for id, en := range catalogs["en"] {
		zh, ok := catalogs["zh-Hans"][id]
		if !ok || zh.Other == "" {
			t.Errorf("missing Chinese: %q", id)
			continue
		}
		a, b := formats.FindAllString(en.Other, -1), formats.FindAllString(zh.Other, -1)
		slices.Sort(a)
		slices.Sort(b)
		if !slices.Equal(a, b) {
			t.Errorf("format mismatch: %q", id)
		}
	}
	if len(catalogs["en"]) != len(catalogs["zh-Hans"]) {
		t.Fatal("catalog IDs differ")
	}
	// Every explicit source message must exist in both catalogs. This includes
	// component defaults supplied by the application adapter.
	files, _ := filepath.Glob("*.go")
	for _, file := range files {
		if filepath.Ext(file) != ".go" || len(file) > 8 && file[len(file)-8:] == "_test.go" {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			index := -1
			if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "tr" {
				index = 1
			}
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "text" {
				index = 0
			}
			if index < 0 || index >= len(call.Args) {
				return true
			}
			lit, ok := call.Args[index].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			id, _ := strconv.Unquote(lit.Value)
			if _, ok := catalogs["en"][id]; !ok {
				t.Errorf("%s: uncatalogued %q", file, id)
			}
			return true
		})
	}
	for _, group := range navigationGroups {
		if _, ok := catalogs["en"][sourceMessageIDs[group]]; !ok {
			t.Errorf("navigation: %s", group)
		}
	}
	for _, row := range notificationRows {
		for _, id := range []string{row.group, row.title, row.description, row.channel} {
			if _, ok := catalogs["en"][sourceMessageIDs[id]]; !ok {
				t.Errorf("settings: %s", id)
			}
		}
	}
}

func TestLanguageSwitchRetainsState(t *testing.T) {
	g := &Gallery{page: welcomePage, config: Config{Theme: mygo.ThemeDark}, configPath: filepath.Join(t.TempDir(), "config.yaml")}
	state := g.controlState("Select", "Selected")
	tt := ui.NewTester(g.Welcome, 820, 900)
	if !tt.HasText("Small components. Real possibilities.") {
		t.Fatal(tt.Texts())
	}
	g.chooseLanguage("zh-Hans")
	tt.Frame()
	if !tt.HasText("小组件，无限可能。") {
		t.Fatal(tt.Texts())
	}
	if state != g.controlState("Select", "Selected") || state.choice != "Large" {
		t.Fatal("state changed")
	}
	saved, err := LoadConfig(g.configPath)
	if err != nil || saved.Language != "zh-Hans" || saved.Theme != mygo.ThemeDark {
		t.Fatalf("config: %+v %v", saved, err)
	}
	other := &Gallery{}
	english := ui.NewTester(other.Welcome, 820, 900)
	if !english.HasText("Small components. Real possibilities.") {
		t.Fatal("language leaked across galleries")
	}
	g.chooseLanguage("en")
	tt.Frame()
	if !tt.HasText("Small components. Real possibilities.") {
		t.Fatal("switch back")
	}
}

func TestLanguageSettings(t *testing.T) {
	g := &Gallery{config: Config{Theme: mygo.ThemeSystem}}
	g.settings.open = true
	g.settings.view.Category = "general"
	tt := ui.NewTester(g.Settings, 820, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("Language preference"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("简体中文"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if g.config.Language != "zh-Hans" || !tt.HasText("语言") {
		t.Fatal(tt.Texts())
	}
	g.settings.view.Query = "密码"
	tt.Frame()
	if !tt.HasText("密码已更改") {
		t.Fatal("localized metadata not searchable", tt.Texts())
	}
	g.settings.view.Query = ""
	g.settings.view.Category = "general"
	tt.Frame()
	if err := tt.Click("语言偏好"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("English"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if g.config.Language != "en" || !tt.HasText("Language") {
		t.Fatal(tt.Texts())
	}
}

func TestChineseControls(t *testing.T) {
	g := &Gallery{}
	g.chooseLanguage("zh-Hans")
	for _, example := range controlCatalog {
		t.Run(example.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { g.Controls(c, example) }, 1000, 900)
			if len(tt.Texts()) == 0 {
				t.Fatal("empty localized example")
			}
		})
	}
}

func TestNativeLanguages(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires desktop session")
	}
	g := &Gallery{}
	mygo.RunOnMain(func() { g.newWindow(nativeBackground()) })
	win := g.Shell.Window()
	defer win.Destroy()
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, language string
		theme          mygo.ThemeSource
		width          int
		settings       bool
	}{
		{"zh-welcome-light", "zh-Hans", mygo.ThemeLight, 820, false},
		{"zh-settings-light", "zh-Hans", mygo.ThemeLight, 820, true},
		{"zh-notifications-light", "zh-Hans", mygo.ThemeLight, 820, true},
		{"zh-settings-dark", "zh-Hans", mygo.ThemeDark, 820, true},
		{"zh-settings-narrow", "zh-Hans", mygo.ThemeLight, 360, true},
		{"en-settings-restored", "en", mygo.ThemeLight, 820, true},
	} {
		mygo.RunOnMain(func() {
			mygo.Theme.SetSource(tc.theme)
			win.SetMinimumSize(360, 480)
			win.SetSize(tc.width, 600)
			g.chooseLanguage(tc.language)
			if tc.settings {
				g.settings.open = true
				g.settings.view.Category = "general"
				if tc.name == "zh-notifications-light" {
					g.settings.view.Category = "notifications"
				}
			}
			win.Invalidate()
		})
		time.Sleep(300 * time.Millisecond)
		data, err := win.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, tc.name+".png")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		t.Log(path)
	}
	for _, name := range []string{"Calendar", "ColorView"} {
		mygo.RunOnMain(func() {
			g.chooseLanguage("zh-Hans")
			g.settings.open = false
			for i, example := range controlCatalog {
				if example.name == name {
					g.page = controlPageStart() + galleryPage(i)
				}
			}
			win.SetSize(820, 600)
			win.Invalidate()
		})
		time.Sleep(300 * time.Millisecond)
		data, err := win.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, "zh-"+name+".png")
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		t.Log(path)
	}
	mygo.RunOnMain(func() { g.chooseLanguage("en") })
	mygo.RunOnMain(func() {
		if win.Title() != "Onevra Gallery" {
			t.Error(fmt.Sprint("title: ", win.Title()))
		}
	})
}

func TestLanguageSwitchPreservesEditorFocus(t *testing.T) {
	g := &Gallery{}
	s := g.controlState("TextBox", "Filled")
	tt := ui.NewTester(func(c *ui.Context) { g.withLanguage(c, func() { controlTextBox(c, s, "Filled") }) }, 500, 150)
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("A personal draft 中文")
	before := s.text
	g.chooseLanguage("zh-Hans")
	tt.Frame()
	if !tt.Focused("名称") || s.text != before {
		t.Fatal("switch lost input or focus", s.text)
	}
	tt.Type("!")
	if s.text == before {
		t.Fatal("editor no longer receives input")
	}
}

func TestLanguageConfigValidation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("theme: dark\nlanguage: unknown\n"), 0600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(path)
	if err == nil || config.Language != "" || config.Theme != mygo.ThemeDark {
		t.Fatalf("fallback: %+v, %v", config, err)
	}
	if err := SaveConfig(path, Config{Theme: mygo.ThemeDark, Language: "unknown"}); err == nil {
		t.Fatal("accepted invalid language")
	}
}
