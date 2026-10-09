package main

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Exercises the real ZIP tool without requiring Finder or creating a disk image.
func TestArchiveWithoutDMG(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("requires macOS ditto")
	}
	stage := t.TempDir()
	app := filepath.Join(stage, "artifacts", "App.app")
	put(t, filepath.Join(app, "Contents", "payload"), "payload")
	c := &config{Name: "App", Version: "1.0"}
	c.Platforms = []string{"darwin/" + runtime.GOARCH + "/zip"}
	target := hostTarget()
	if err := archive(context.Background(), c, stage, app, target); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(filepath.Dir(app), c.archiveName(target))
	z, err := zip.OpenReader(base + ".zip")
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	found := false
	for _, f := range z.File {
		if f.Name == "App.app/Contents/payload" {
			found = true
		}
	}
	if !found {
		t.Fatal("ZIP does not contain app payload")
	}
	for _, path := range []string{base + ".dmg", filepath.Join(stage, "image")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected DMG staging/output %s: %v", path, err)
		}
	}
}
