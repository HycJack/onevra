package gallery

import "testing"

func TestPreviewEntries(t *testing.T) {
	entries := PreviewEntries()
	if len(entries) != len(controlCatalog)+4 {
		t.Fatalf("got %d previews, want %d", len(entries), len(controlCatalog)+4)
	}
	seen := map[string]bool{}
	components := map[string]bool{}
	for _, entry := range entries {
		if entry.Slug == "" || seen[entry.Slug] {
			t.Errorf("preview slug %q is empty or repeated", entry.Slug)
		}
		seen[entry.Slug] = true
		if entry.Key == "" || entry.Name == "" || entry.NameZh == "" || entry.Group == "" || entry.GroupZh == "" {
			t.Errorf("%q has incomplete metadata: %+v", entry.Slug, entry)
		}
		if len(entry.Components) == 0 {
			t.Errorf("%q demonstrates no component", entry.Slug)
		}
		for _, component := range entry.Components {
			components[previewSlug(component)] = true
		}
	}
	for _, name := range []string{
		"button", "icon-button", "toggle-button", "select", "native-select",
		"calendar", "text-box", "text-box-buffer", "table", "chat",
		"settings", "search-settings", "command-view", "command-dialog", "icon", "pixel-grid",
	} {
		if !components[name] {
			t.Errorf("no preview demonstrates %s", name)
		}
	}
}

func TestPreviewSelection(t *testing.T) {
	p := NewPreview("calendar", "en")
	if got := p.Title(); got != "Date & time pickers" {
		t.Errorf("calendar title = %q", got)
	}
	if !p.Show("Slider") || p.Title() != "Slider" {
		t.Errorf("slider title = %q", p.Title())
	}
	if p.Show("missing") {
		t.Error("missing slug reported as shown")
	}
	zh := NewPreview("button", "zh-Hans")
	if got := zh.Title(); got != "按钮" {
		t.Errorf("button title in Chinese = %q", got)
	}
	if got := NewPreview("settings", "zh-Hans").Title(); got != "设置" {
		t.Errorf("settings title in Chinese = %q", got)
	}
	if got := NewPreview("command-palette", "en").Title(); got != "Command palette" {
		t.Errorf("command palette title = %q", got)
	}
	fallback := NewPreview("missing", "en")
	if got := fallback.Title(); got != PreviewEntries()[0].Name {
		t.Errorf("fallback title = %q", got)
	}
}
