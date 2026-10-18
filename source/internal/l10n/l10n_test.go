package l10n

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/egoist/mygo/i18n"
)

func messages(t *testing.T, tag string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(tag + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestTranslationsAreComplete(t *testing.T) {
	en, zh := messages(t, "en"), messages(t, "zh-Hans")
	for key := range en {
		if strings.HasPrefix(key, "@") {
			continue
		}
		if _, ok := zh[key]; !ok {
			t.Errorf("zh-Hans lacks %s", key)
		}
		// Each Chinese message resolves through the catalog rather than falling back.
		if _, args := en["@"+key]; !args && In("zh-Hans").Text(key) == en[key] {
			t.Errorf("zh-Hans %s falls back to English", key)
		}
	}
	for key := range zh {
		if _, ok := en[key]; !ok && !strings.HasPrefix(key, "@") {
			t.Errorf("zh-Hans has unknown key %s", key)
		}
	}
}

func TestMessages(t *testing.T) {
	if got := In("").Text("toggleSidebar"); got != "Toggle sidebar" {
		t.Fatalf("English: %q", got)
	}
	if got := In("zh-Hans-CN").Text("closeWindow"); got != "关闭窗口" {
		t.Fatalf("zh-Hans-CN falls back to zh-Hans: %q", got)
	}
	if got := In("en").Label("numberAdjust", ""); got != "Adjust" {
		t.Fatalf("empty label leaves a trimmed name: %q", got)
	}
	if got := In("zh-Hans").Label("numberIncrease", "数量"); got != "增加数量" {
		t.Fatalf("Chinese label: %q", got)
	}
	if got := In("en").Format("slideOf", map[string]any{"index": 2, "count": 5}); got != "2 of 5" {
		t.Fatalf("formatted message: %q", got)
	}
}

func TestApplicationOverrides(t *testing.T) {
	app, err := i18n.NewCatalog(fstest.MapFS{"en.json": {Data: []byte(`{"onevra.commands": "Actions"}`)}}, "en")
	if err != nil {
		t.Fatal(err)
	}
	i18n.Register(app)
	if got := In("en").Text("commands"); got != "Actions" {
		t.Fatalf("an app's onevra.commands: %q", got)
	}
}
