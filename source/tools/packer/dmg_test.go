package main

import (
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestDMGBackground(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".background", "install.png")
	if err := writeDMGBackground(path); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 720 || img.Bounds().Dy() != 440 {
		t.Fatalf("unexpected background bounds: %v", img.Bounds())
	}
	for _, x := range []int{180, 540} {
		if got := color.RGBAModel.Convert(img.At(x, 185)); got != (color.RGBA{255, 255, 255, 255}) {
			t.Fatalf("icon area should be empty: %v", got)
		}
	}
	if got := color.RGBAModel.Convert(img.At(360, 185)); got == (color.RGBA{255, 255, 255, 255}) {
		t.Fatal("missing direction cue")
	}
	if err := writeDMGBackground(filepath.Join(path, "invalid.png")); err == nil {
		t.Fatal("expected invalid destination failure")
	}
}
