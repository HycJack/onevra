package gallery

import (
	"bytes"
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
	if err := os.WriteFile(filepath.Join(dir, "images", "icon.png"), []byte("invalid PNG"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadGalleryImage(dir); err == nil {
		t.Fatal("invalid image accepted")
	}
}

func TestGalleryImageAsset(t *testing.T) {
	original, err := os.ReadFile("../../../../icon.png")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := os.ReadFile("../../../../resources/images/icon.png")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, payload) {
		t.Fatal("refresh resources/images/icon.png from source/icon.png")
	}
	bitmap, err := galleryImage()
	if err != nil {
		t.Fatal(err)
	}
	if w, h := bitmap.Size(); w != 1024 || h != 1024 {
		t.Fatalf("unexpected icon dimensions: %dx%d", w, h)
	}
	again, err := galleryImage()
	if err != nil || again != bitmap {
		t.Fatal("decoded icon is not cached")
	}
}
