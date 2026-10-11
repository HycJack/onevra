package main

import (
	"os"
	"path/filepath"
	"slices"
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

func TestReleaseNames(t *testing.T) {
	c := &config{Name: "Test App", Version: "1.2.3", Platforms: []string{"darwin/arm64/zip", "darwin/arm64/dmg", "darwin/arm64/zip"}}
	for _, keep := range []bool{true, false} {
		c.KeepUncompressed = keep
		for _, goos := range []string{"darwin", "windows", "linux"} {
			for _, arch := range []string{"arm64", "amd64"} {
				target := target{goos, arch}
				var want []string
				for _, suffix := range c.archiveSuffixes(target) {
					want = append(want, c.archiveName(target)+suffix)
				}
				if keep && target == hostTarget() {
					app := c.Name
					if goos == "darwin" {
						app += ".app"
					}
					want = append(want, app)
				}
				if got := c.releaseNames(target); !slices.Equal(got, want) {
					t.Errorf("%s, keep=%v: names = %v, want %v", target.key(), keep, got, want)
				}
			}
		}
	}
}

func TestPublishArtifactsRollback(t *testing.T) {
	src := t.TempDir()
	// Moving App into its own descendant fails after the preceding artifacts
	// have been installed, exercising rollback for both files and directories.
	dst := filepath.Join(src, "App", "dist")
	names := []string{"Other.app", "app.zip", "App"}
	put(t, filepath.Join(src, "Other.app", "new"), "new app")
	put(t, filepath.Join(src, "app.zip"), "new archive")
	put(t, filepath.Join(dst, "Other.app", "old"), "previous app")
	put(t, filepath.Join(dst, "app.zip"), "previous archive")
	if err := publishArtifacts(src, dst, names); err == nil {
		t.Fatal("expected failed publication")
	}
	for name, want := range map[string]string{"Other.app/old": "previous app", "app.zip": "previous archive"} {
		if data, err := os.ReadFile(filepath.Join(dst, name)); err != nil || string(data) != want {
			t.Fatalf("previous artifact was not restored: %s: %q, %v", name, data, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, "Other.app", "new")); !os.IsNotExist(err) {
		t.Fatal("new app was not rolled back")
	}
	if entries, err := os.ReadDir(dst); err != nil || len(entries) != 2 {
		t.Fatalf("unexpected rollback leftovers: %v, %v", entries, err)
	}
}

func TestPublishArtifactsRejectsTypeMismatch(t *testing.T) {
	for _, directorySource := range []bool{false, true} {
		src, dst := t.TempDir(), t.TempDir()
		sourcePath, destPath := filepath.Join(src, "app"), filepath.Join(dst, "app")
		if directorySource {
			sourcePath = filepath.Join(sourcePath, "payload")
		} else {
			destPath = filepath.Join(destPath, "payload")
		}
		put(t, sourcePath, "new")
		put(t, destPath, "previous")
		if err := publishArtifacts(src, dst, []string{"app"}); err == nil {
			t.Fatal("accepted mismatched artifact types")
		}
		if data, err := os.ReadFile(destPath); err != nil || string(data) != "previous" {
			t.Fatalf("previous artifact changed: %q, %v", data, err)
		}
	}
}

func TestPublishArtifactsRejectsSymlinks(t *testing.T) {
	for _, sourceLink := range []bool{true, false} {
		src, dst, external := t.TempDir(), t.TempDir(), t.TempDir()
		put(t, filepath.Join(external, "keep"), "unrelated")
		link, other := filepath.Join(src, "App"), filepath.Join(dst, "App")
		if !sourceLink {
			link, other = other, link
		}
		if err := os.Symlink(external, link); err != nil {
			t.Skipf("cannot create symlinks: %v", err)
		}
		put(t, filepath.Join(other, "payload"), "app")
		if err := publishArtifacts(src, dst, []string{"App"}); err == nil {
			t.Fatal("accepted a symlink artifact")
		}
		if data, err := os.ReadFile(filepath.Join(external, "keep")); err != nil || string(data) != "unrelated" {
			t.Fatalf("symlink target changed: %q, %v", data, err)
		}
	}
}

func TestPublishArtifacts(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "stage"), filepath.Join(root, "dist")
	names := []string{"app.dmg", "app.zip", "App.app"}
	files := []string{"app.dmg", "app.zip", "App.app/Contents/MacOS/App"}
	for _, name := range files {
		put(t, filepath.Join(dst, name), "previous")
	}
	put(t, filepath.Join(dst, "App.app", "stale"), "old payload")
	put(t, filepath.Join(dst, "App Dev.app", "keep"), "development app")
	put(t, filepath.Join(dst, "other-target.zip"), "other target")
	put(t, filepath.Join(src, names[0]), "new")
	if err := publishArtifacts(src, dst, names); err == nil {
		t.Fatal("incomplete artifact set must fail")
	}
	for _, name := range files {
		if data, _ := os.ReadFile(filepath.Join(dst, name)); string(data) != "previous" {
			t.Fatalf("previous archive %s lost", name)
		}
		put(t, filepath.Join(src, name), "new")
	}
	if err := publishArtifacts(src, dst, names); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dst, "App.app", "stale")); !os.IsNotExist(err) {
		t.Fatal("stale app payload survived publication")
	}
	for _, name := range files {
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
	if err != nil || len(entries) != 5 {
		t.Fatalf("unexpected leftovers: %v, %v", entries, err)
	}
}
