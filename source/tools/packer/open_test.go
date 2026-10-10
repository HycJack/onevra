package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDirectoryOpener(t *testing.T) {
	for host, want := range map[string]string{"darwin": "open", "windows": "explorer.exe", "linux": "xdg-open"} {
		got, err := directoryOpener(host)
		if err != nil || got != want {
			t.Errorf("%s: got %q, %v; want %q", host, got, err, want)
		}
	}
	if _, err := directoryOpener("unsupported"); err == nil {
		t.Fatal("accepted unsupported host")
	}
}

func TestBuildOpenReleases(t *testing.T) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("fake desktop opener uses a POSIX shell")
	}
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/openerfixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), "package main\nfunc main() {}\n")
	bin := filepath.Join(root, "bin")
	name, err := directoryOpener(runtime.GOOS)
	if err != nil {
		t.Fatal(err)
	}
	opener := filepath.Join(bin, name)
	put(t, opener, "#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$PACKER_OPEN_MARKER\"\nexit \"$PACKER_OPEN_EXIT\"\n")
	if err := os.Chmod(opener, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "opened")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("PACKER_OPEN_MARKER", marker)
	t.Chdir(root)
	for _, tc := range []struct {
		name, setting, main, exit string
		wantOpen, wantError       bool
	}{
		{"default", "", ".", "0", true, false},
		{"disabled", "openReleases: false\n", ".", "0", false, false},
		{"opener failure", "", ".", "1", true, false},
		{"build failure", "", "./missing", "0", false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			os.Remove(marker)
			t.Setenv("PACKER_OPEN_EXIT", tc.exit)
			put(t, filepath.Join(root, "packer.yaml"), "name: App\nidentifier: me.app\nout: release output\nplatforms: [linux/amd64, linux/arm64]\nmain: "+tc.main+"\n"+tc.setting)
			releaseDir := filepath.Join(root, "release output", "releases", "0.0.0")
			stale := filepath.Join(releaseDir, "old-target", "stale.zip")
			put(t, stale, "stale archive")
			retained := []string{
				filepath.Join(root, "release output", "releases", "1.2.3", "previous.zip"),
				filepath.Join(root, "release output", "App Dev", "keep"),
			}
			for _, path := range retained {
				put(t, path, "keep")
			}
			err := run(context.Background(), []string{"build"})
			if (err != nil) != tc.wantError {
				t.Fatalf("build error = %v, wantError = %v", err, tc.wantError)
			}
			if _, err := os.Stat(stale); !os.IsNotExist(err) {
				t.Fatalf("stale archive was not removed: %v", err)
			}
			for _, path := range retained {
				if data, err := os.ReadFile(path); err != nil || string(data) != "keep" {
					t.Fatalf("unrelated output changed: %s: %q, %v", path, data, err)
				}
			}
			entries, err := os.ReadDir(releaseDir)
			wantCount := 2
			if tc.wantError {
				wantCount = 0
			}
			if err != nil || len(entries) != wantCount {
				t.Fatalf("release directory entries = %v, %v; want %d archives", entries, err, wantCount)
			}
			data, err := os.ReadFile(marker)
			if !tc.wantOpen {
				if !os.IsNotExist(err) {
					t.Fatalf("unexpected opener invocation: %q, %v", data, err)
				}
				return
			}
			want := releaseDir + "\n"
			if err != nil || string(data) != want {
				t.Fatalf("opener args = %q, %v; want one invocation with %q", data, err, want)
			}
			for _, arch := range []string{"amd64", "arm64"} {
				archive := filepath.Join(strings.TrimSpace(want), "app-linux-"+arch+"-0.0.0.tar.gz")
				if _, err := os.Stat(archive); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
