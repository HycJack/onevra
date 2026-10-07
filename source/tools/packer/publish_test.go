package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReleaseDir(t *testing.T) {
	root := t.TempDir()
	for _, out := range []string{"dist", "../dist", filepath.Join(root, "custom")} {
		for _, version := range []string{"0.2026.1", "1.2.3", "2.0.0-rc.1+build.42"} {
			c := &config{root: root, Out: out, Version: version}
			base := out
			if !filepath.IsAbs(base) {
				base = filepath.Join(root, base)
			}
			if got, want := c.releaseDir(), filepath.Join(base, "releases", version); got != want {
				t.Errorf("release directory = %q, want %q", got, want)
			}
		}
	}
}

func TestArchiveName(t *testing.T) {
	for _, version := range []string{"0.0.0", "1.2.3", "2.0.0-rc.1+build.42"} {
		c := &config{Name: "Onevra Gallery", Version: version}
		for _, goos := range []string{"darwin", "windows", "linux"} {
			for _, arch := range []string{"arm64", "amd64"} {
				if got, want := c.archiveName(target{goos, arch}), "onevra-gallery-"+goos+"-"+arch+"-"+version; got != want {
					t.Errorf("archive name = %q, want %q", got, want)
				}
			}
		}
	}
}

func TestPublishArchives(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "stage"), filepath.Join(root, "dist")
	names := []string{"app.dmg", "app.zip"}
	for _, name := range names {
		put(t, filepath.Join(dst, name), "previous")
	}
	put(t, filepath.Join(dst, "App Dev.app", "keep"), "development app")
	put(t, filepath.Join(dst, "other-target.zip"), "other target")
	put(t, filepath.Join(src, names[0]), "new")
	if err := publishArchives(src, dst, names); err == nil {
		t.Fatal("incomplete archive set must fail")
	}
	for _, name := range names {
		if data, _ := os.ReadFile(filepath.Join(dst, name)); string(data) != "previous" {
			t.Fatalf("previous archive %s lost", name)
		}
		put(t, filepath.Join(src, name), "new")
	}
	if err := publishArchives(src, dst, names); err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if data, _ := os.ReadFile(filepath.Join(dst, name)); string(data) != "new" {
			t.Fatalf("archive %s was not replaced", name)
		}
	}
	for path, want := range map[string]string{"App Dev.app/keep": "development app", "other-target.zip": "other target"} {
		if data, _ := os.ReadFile(filepath.Join(dst, path)); string(data) != want {
			t.Fatalf("unrelated output changed: %s", path)
		}
	}
	entries, err := os.ReadDir(dst)
	if err != nil || len(entries) != 4 {
		t.Fatalf("unexpected leftovers: %v, %v", entries, err)
	}
}
