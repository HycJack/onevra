package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSystemPagesRenderHeadless(t *testing.T) {
	pages := map[string]string{
		"Recent documents": "recentEmpty",
		"File watching":    "watchEmpty",
		"Global shortcuts": "shortcutRegistered",
	}
	for name, key := range pages {
		for _, lang := range []string{"en", "zh-Hans"} {
			g, s, tt := featureTester(name, "Default", lang, 540)
			if s.system.native {
				t.Fatalf("%s: headless views must not call system services", name)
			}
			if !tt.HasText(g.text(key)) {
				t.Errorf("%s/%s: missing %q in %q", name, lang, g.text(key), tt.Texts())
			}
			s.system.close()
		}
	}
	tabs := []string{"appearanceReduceMotion", "fontsSample", "languagesShown", "networkOffline", "secretsSave"}
	for _, lang := range []string{"en", "zh-Hans"} {
		for tab, key := range tabs {
			g, s, tt := featureTester("System information", "Default", lang, 540)
			if s.system.native {
				t.Fatal("System information: headless views must not call system services")
			}
			s.system.systemTab = tab
			tt.Frame()
			if !tt.HasText(g.text(key)) {
				t.Errorf("System information/%s tab %d: missing %q in %q", lang, tab, g.text(key), tt.Texts())
			}
			s.system.close()
		}
	}
}

func TestSystemAppearanceFollowsPreferences(t *testing.T) {
	g, _, tt := featureTester("System information", "Default", "en", 540)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1.25, Accent: ui.RGB(37, 99, 235)})
	if !tt.HasText("125%") || !tt.HasText(g.text("systemYes")) {
		t.Fatalf("preferences not shown: %q", tt.Texts())
	}
	if _, ok := tt.Find(g.text("appearanceAccent")); !ok {
		t.Fatal("accent swatch missing")
	}
}

func TestLanguagesFollowWindowLocale(t *testing.T) {
	_, s, tt := featureTester("System information", "Default", "zh-Hans", 540)
	s.system.systemTab = 2
	tt.SetLocale("zh-Hans")
	if !tt.HasText("简体中文") || !tt.HasText("1,234,567.891") {
		t.Fatalf("zh-Hans formats: %q", tt.Texts())
	}
}
