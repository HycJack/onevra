package gallery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGalleryImageErrors(t *testing.T) {
	dir := t.TempDir()
	if _, err := loadGalleryImage(dir); !os.IsNotExist(err) {
		t.Fatalf("missing resource error: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "images"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "images", "avatar.jpeg"), []byte("invalid JPEG"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadGalleryImage(dir); err == nil {
		t.Fatal("invalid image accepted")
	}
}

func TestGalleryImageAsset(t *testing.T) {
	bitmap, err := galleryImage()
	if err != nil {
		t.Fatal(err)
	}
	if w, h := bitmap.Size(); w != 1254 || h != 1254 {
		t.Fatalf("unexpected avatar dimensions: %dx%d", w, h)
	}
	again, err := galleryImage()
	if err != nil || again != bitmap {
		t.Fatal("decoded avatar is not cached")
	}
}
