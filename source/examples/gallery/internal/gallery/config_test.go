package gallery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preferences", "config.yaml")
	config, err := LoadConfig(path)
	if err != nil || config.Theme != mygo.ThemeSystem {
		t.Fatalf("missing: %+v, %v", config, err)
	}
	for _, source := range []mygo.ThemeSource{mygo.ThemeSystem, mygo.ThemeLight, mygo.ThemeDark} {
		want := Config{Theme: source}
		if err := SaveConfig(path, want); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil || string(data) != "theme: "+string(source)+"\n" {
			t.Fatalf("expected YAML config: %s, %v", data, err)
		}
		got, err := LoadConfig(path)
		if err != nil || got != want {
			t.Fatalf("round trip: %+v, %v", got, err)
		}
	}
	for _, data := range []string{`{}`, `theme: "unknown"`, `theme: 123`, `{`} {
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		got, err := LoadConfig(path)
		if got.Theme != mygo.ThemeSystem || (data != `{}` && err == nil) {
			t.Fatalf("%s: %+v, %v", data, got, err)
		}
	}
	if err := SaveConfig(path, Config{Theme: "invalid"}); err == nil {
		t.Fatal("accepted invalid theme")
	}
	if err := SaveConfig(filepath.Join(path, "config.yaml"), Config{Theme: mygo.ThemeLight}); err == nil {
		t.Fatal("expected write error")
	}
	if _, err := LoadConfig(filepath.Dir(path)); err == nil {
		t.Fatal("expected read error")
	}
}

func TestNativeThemeMenu(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop event loop for theme updates")
	}
	previous := mygo.Theme.Source()
	defer mygo.Theme.SetSource(previous)
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := SaveConfig(path, Config{Theme: mygo.ThemeDark, Language: "en"}); err != nil {
		t.Fatal(err)
	}
	g := &Gallery{}
	g.Configure(path, mygo.ThemeLight)
	if mygo.Theme.Source() != mygo.ThemeLight {
		t.Fatal("launch override ignored")
	}
	saved, err := LoadConfig(path)
	if err != nil || saved.Theme != mygo.ThemeDark {
		t.Fatal("override rewrote preference")
	}
	g.Configure(path, "")
	if mygo.Theme.Source() != mygo.ThemeDark {
		t.Fatal("saved preference not restored")
	}
	tt := ui.NewTester(func(c *ui.Context) { ui.Toolbar(c, func() { g.LeftToolbar(c) }).FillWidth() }, 220, 240)
	theme, ok := tt.Find("Theme")
	settings, settingsOK := tt.Find("Settings")
	if !ok || !settingsOK || settings.X >= theme.X || theme.W != settings.W {
		t.Fatalf("toolbar layout: %v, %v", theme, settings)
	}
	for _, option := range []struct {
		label  string
		source mygo.ThemeSource
		icon   *ui.SVG
	}{{"System", mygo.ThemeSystem, lucide.Monitor()}, {"Light", mygo.ThemeLight, lucide.Sun()}, {"Dark", mygo.ThemeDark, lucide.Moon()}} {
		if err := tt.Click("Theme"); err != nil {
			t.Fatal(err)
		}
		if err := tt.ChooseMenuItem(option.label); err != nil {
			t.Fatal(err)
		}
		if mygo.Theme.Source() != option.source {
			t.Fatalf("%s not applied", option.label)
		}
		if themeIcon() != option.icon {
			t.Fatalf("%s has the wrong toolbar icon", option.label)
		}
		saved, err := LoadConfig(path)
		if err != nil || saved.Theme != option.source {
			t.Fatalf("%s not persisted: %+v, %v", option.label, saved, err)
		}
	}
	g.settings.initialize()
	g.settings.open, g.settings.view.Category = true, "appearance"
	settingsUI := ui.NewTester(g.Settings, 820, 600)
	settingsUI.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := settingsUI.Click("Theme preference"); err != nil {
		t.Fatal(err)
	}
	if err := settingsUI.Click("Light"); err != nil {
		t.Fatal(err)
	}
	saved, err = LoadConfig(path)
	if err != nil || saved.Theme != mygo.ThemeLight || mygo.Theme.Source() != mygo.ThemeLight {
		t.Fatalf("settings theme did not apply and persist: %+v, %v", saved, err)
	}
}
