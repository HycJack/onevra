package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/plugins/terminal"
)

func TestTerminalLibraryLocation(t *testing.T) {
	root := t.TempDir()
	contents := filepath.Join(root, "Test App.app", "Contents")
	exe := filepath.Join(contents, "MacOS", "Test App")
	want := filepath.Join(contents, "Frameworks", "libghostty-vt.dylib")
	t.Setenv("MYGO_GHOSTTY_VT", "")
	if err := useBundledTerminalLibrary(exe); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("MYGO_GHOSTTY_VT"); got != want {
		t.Fatalf("bundle path = %q, want %q", got, want)
	}
	// Missing bundled files must not fall back to an unrelated cached version.
	if got, err := terminal.LibraryPath(); err != nil || got != want {
		t.Fatalf("terminal lookup = %q, %v", got, err)
	}
	t.Setenv("MYGO_GHOSTTY_VT", "/explicit/library.dylib")
	if err := useBundledTerminalLibrary(exe); err != nil || os.Getenv("MYGO_GHOSTTY_VT") != "/explicit/library.dylib" {
		t.Fatalf("explicit override was not preserved: %v", err)
	}
	for _, path := range []string{
		filepath.Join(root, "go-build", "gallery"),
		filepath.Join(root, "not-a-bundle", "Contents", "MacOS", "gallery"),
		filepath.Join(contents, "gallery"),
	} {
		t.Setenv("MYGO_GHOSTTY_VT", "")
		if err := useBundledTerminalLibrary(path); err != nil || os.Getenv("MYGO_GHOSTTY_VT") != "" {
			t.Fatalf("unpackaged executable changed terminal lookup: %s, %v", path, err)
		}
	}
}

func TestBundledTerminalLoad(t *testing.T) {
	bundle := os.Getenv("ONEO_GALLERY_BUNDLE")
	if runtime.GOOS != "darwin" || bundle == "" {
		t.Skip("set ONEO_GALLERY_BUNDLE to a Packer-built gallery .app on macOS")
	}
	// Relocate the library so this cannot pass because of the build path or cache.
	data, err := os.ReadFile(filepath.Join(bundle, "Contents", "Frameworks", "libghostty-vt.dylib"))
	if err != nil {
		t.Fatal(err)
	}
	contents := filepath.Join(t.TempDir(), "Relocated Gallery.app", "Contents")
	frameworks := filepath.Join(contents, "Frameworks")
	if err := os.MkdirAll(frameworks, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(frameworks, "libghostty-vt.dylib")
	if err := os.WriteFile(path, data, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MYGO_GHOSTTY_VT", "")
	if err := useBundledTerminalLibrary(filepath.Join(contents, "MacOS", "Gallery")); err != nil {
		t.Fatal(err)
	}
	if got, err := terminal.LibraryPath(); err != nil || got != path || !strings.Contains(got, "Frameworks") {
		t.Fatalf("terminal selected %q: %v", got, err)
	}
	if err := terminal.Load(); err != nil {
		t.Fatalf("load relocated bundled library: %v", err)
	}
}
