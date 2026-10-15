package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSystemPagesRenderHeadless(t *testing.T) {
	pages := map[string]string{
		"Network": "networkOffline", "Secrets": "secretsSave", "Recent documents": "recentEmpty",
		"File watching": "watchEmpty", "Installed apps": "appsHint", "Fonts": "fontsMonospaced",
		"System appearance": "appearanceReduceMotion", "Languages & formats": "languagesShown",
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
}

func TestSystemAppearanceFollowsPreferences(t *testing.T) {
	g, _, tt := featureTester("System appearance", "Default", "en", 540)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1.25, Accent: ui.RGB(37, 99, 235)})
	if !tt.HasText("125%") || !tt.HasText(g.text("systemYes")) {
		t.Fatalf("preferences not shown: %q", tt.Texts())
	}
	if _, ok := tt.Find(g.text("appearanceAccent")); !ok {
		t.Fatal("accent swatch missing")
	}
}

func TestLanguagesFollowWindowLocale(t *testing.T) {
	_, _, tt := featureTester("Languages & formats", "Default", "zh-Hans", 540)
	tt.SetLocale("zh-Hans")
	if !tt.HasText("简体中文") || !tt.HasText("1,234,567.891") {
		t.Fatalf("zh-Hans formats: %q", tt.Texts())
	}
}
