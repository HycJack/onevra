package gallery

import (
	"os"
	"testing"
)

func TestGalleryMarkdownResource(t *testing.T) {
	if _, err := loadGalleryMarkdown(t.TempDir()); !os.IsNotExist(err) {
		t.Fatalf("missing Markdown error: %v", err)
	}
	if _, err := galleryMarkdown(); err != nil {
		t.Fatal(err)
	}
	g, s, tt := featureTester("Markdown", "Default", "en", 320)
	for _, value := range []string{"Missing image — alternative text", "An unregistered language remains a selectable code block.", "<div>Raw HTML stays visible as text.</div>", "End of the document. Select All includes content outside the viewport."} {
		if !tt.HasText(value) {
			t.Errorf("missing fixture content: %q", value)
		}
	}
	if err := tt.Click(g.text("increment")); err != nil || s.index != 1 {
		t.Fatalf("custom code block: %v, count %d", err, s.index)
	}
}
