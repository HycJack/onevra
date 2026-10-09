package gallery

import (
	"math"
	"os"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGalleryMarkdownResource(t *testing.T) {
	if _, err := loadGalleryMarkdown(t.TempDir()); !os.IsNotExist(err) {
		t.Fatalf("missing Markdown error: %v", err)
	}
	document, err := galleryMarkdown()
	if err != nil {
		t.Fatal(err)
	}
	again, err := galleryMarkdown()
	if err != nil || document != again {
		t.Fatal("parsed resource is not cached")
	}
	tt := ui.NewTester(func(c *ui.Context) { controlMarkdown(c, &controlState{}, "Default") }, 320, 5000)
	image, ok := tt.Find("Author avatar")
	if !ok || image.W <= 0 || image.H <= 0 || image.W > 320 || math.Abs(float64(image.W-image.H)) > .01 {
		t.Fatalf("resource image: %v, found %v", image, ok)
	}
	for _, value := range []string{"Missing image — alternative text", "An unregistered language remains a selectable code block.", "<div>Raw HTML stays visible as text.</div>", "End of the document. Select All includes content outside the viewport."} {
		if !tt.HasText(value) {
			t.Errorf("missing fixture content: %q", value)
		}
	}
}

func TestMarkdownPreviewHasOnlyDocument(t *testing.T) {
	g := &Gallery{page: pageFromKey("controls/Markdown")}
	tt := ui.NewTester(g.View, 600, 1000)
	headings := 0
	for _, text := range tt.Texts() {
		if text == "Markdown" {
			headings++
		}
		if text == "Default" || text == "Plugins / Markdown" {
			t.Fatalf("extra gallery content: %q", text)
		}
	}
	if headings != 1 {
		t.Fatalf("expected only the document heading, got %d", headings)
	}
}
