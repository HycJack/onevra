package main

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRuntimeConfig(t *testing.T) {
	file := filepath.Join(t.TempDir(), "packer.yaml")
	put(t, file, "name: App\nidentifier: me.app\n")
	if _, err := loadConfig(file); err != nil {
		t.Fatal(err)
	}
	// Runtime discovery is a directory convention, not another YAML setting.
	for _, setting := range []string{
		"nativeLibraries: {}",
		"nativeLibraries:\n  darwin/arm64: [libaudio.dylib]",
		"runtimes: runtimes",
	} {
		put(t, file, "name: App\nidentifier: me.app\n"+setting+"\n")
		if _, err := loadConfig(file); err == nil {
			t.Errorf("accepted unsupported setting: %s", setting)
		}
	}
}

func TestRuntimeFiles(t *testing.T) {
	root := t.TempDir()
	c := &config{root: root}
	platform := target{"darwin", "arm64"}
	check := func(want []string) {
		t.Helper()
		files, err := c.runtimeFiles(platform)
		if err != nil || !reflect.DeepEqual(files, want) {
			t.Fatalf("runtime files = %v, %v; want %v", files, err, want)
		}
	}
	check(nil)
	// Other targets, including malformed directories, are not read.
	put(t, filepath.Join(root, "runtimes", "linux"), "not a directory")
	put(t, filepath.Join(root, "runtimes", "darwin", "amd64", "other.dylib"), "other architecture")
	check(nil)
	dir := filepath.Join(root, "runtimes", "darwin", "arm64")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	check(nil)
	audio := filepath.Join(dir, "libaudio.dylib")
	other := filepath.Join(dir, "libother.dylib")
	put(t, other, "other")
	put(t, audio, "audio")
	// Resolve from the configuration directory, not the process working directory.
	t.Chdir(t.TempDir())
	check([]string{audio, other})
	if files, err := c.runtimeFiles(target{"windows", "amd64"}); err != nil || len(files) != 0 {
		t.Fatalf("absent target should add no libraries: %v, %v", files, err)
	}
}

func TestRuntimeInvalidSources(t *testing.T) {
	for _, name := range []string{"runtimes", "runtimes/darwin", "runtimes/darwin/arm64"} {
		t.Run("file/"+name, func(t *testing.T) {
			root := t.TempDir()
			put(t, filepath.Join(root, name), "not a directory")
			if _, err := (&config{root: root}).runtimeFiles(target{"darwin", "arm64"}); err == nil {
				t.Fatal("accepted a file in place of a runtime directory")
			}
		})
	}
	for _, name := range []string{"runtimes", "runtimes/darwin", "runtimes/darwin/arm64", "runtimes/darwin/arm64/libaudio.dylib"} {
		t.Run("symlink/"+name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, name)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			for _, source := range []string{t.TempDir(), filepath.Join(root, "missing")} {
				if err := os.Symlink(source, path); err != nil {
					t.Skipf("cannot create symlink: %v", err)
				}
				if _, err := (&config{root: root}).runtimeFiles(target{"darwin", "arm64"}); err == nil {
					t.Fatal("accepted a runtime symlink")
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
	root := t.TempDir()
	put(t, filepath.Join(root, "runtimes", "darwin", "arm64", "nested", "libaudio.dylib"), "audio")
	if _, err := (&config{root: root}).runtimeFiles(target{"darwin", "arm64"}); err == nil {
		t.Fatal("accepted a nested runtime directory")
	}
}

func TestNativeNameCollisions(t *testing.T) {
	for _, names := range [][]string{
		{"libaudio.so", "libaudio.so"},
		{"libaudio.so", "LibAudio.so"},
		{"fixture"},
		{"FIXTURE"},
		{"assets"},
		{"../escape"},
	} {
		dst := t.TempDir()
		put(t, filepath.Join(dst, "fixture"), "executable")
		if err := os.Mkdir(filepath.Join(dst, "assets"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := checkNativeNames(dst, names); err == nil {
			t.Errorf("accepted conflicting or unsafe names: %v", names)
		}
		if data, _ := os.ReadFile(filepath.Join(dst, "fixture")); string(data) != "executable" {
			t.Fatal("collision modified existing output")
		}
	}
	if err := checkNativeNames(filepath.Join(t.TempDir(), "new"), []string{"libaudio.so", "libother.so"}); err != nil {
		t.Fatal(err)
	}
}

func TestPortableNativeCollisions(t *testing.T) {
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), "package main\nfunc main() {}\n")
	put(t, filepath.Join(root, "resources", "payload.so"), "payload")
	put(t, filepath.Join(root, "icon.svg"), "icon")
	c := &config{Name: "Fixture", root: root, Main: "."}
	c.Linux.Icon = "icon.svg"
	out := filepath.Join(root, "out")
	marker := filepath.Join(out, c.Name, "previous")
	put(t, marker, "previous build")
	for _, name := range []string{"fixture", "FIXTURE", "payload.so", "PAYLOAD.SO", "icon.svg"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(root, "runtimes", "linux", "amd64", name)
			put(t, path, "local library")
			t.Cleanup(func() { os.Remove(path) })
			if _, err := build(context.Background(), c, out, false, target{"linux", "amd64"}); err == nil || !strings.Contains(err.Error(), "conflicting native filename") {
				t.Fatalf("expected a collision with the executable, payload or icon: %v", err)
			}
			if data, err := os.ReadFile(marker); err != nil || string(data) != "previous build" {
				t.Fatalf("failed build replaced previous output: %q, %v", data, err)
			}
			entries, err := os.ReadDir(out)
			if err != nil || len(entries) != 1 || entries[0].Name() != c.Name {
				t.Fatalf("failed build left staging output: %v, %v", entries, err)
			}
		})
	}
}

func TestBundleRuntimes(t *testing.T) {
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), "package main\nfunc main() {}\n")
	lib := filepath.Join(root, "runtimes", "linux", "amd64", "libaudio.so")
	put(t, lib, "local audio")
	if err := os.Chmod(lib, 0o755); err != nil {
		t.Fatal(err)
	}
	c := &config{root: root, Main: "."}
	dst := filepath.Join(root, "app")
	ctx := context.Background()
	platform := target{"linux", "amd64"}
	if err := bundleNatives(ctx, c, dst, platform); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dst, "libaudio.so")); err != nil || string(data) != "local audio" {
		t.Fatalf("copied library: %q, %v", data, err)
	}
	if info, err := os.Stat(filepath.Join(dst, "libaudio.so")); err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("library lost executable permissions: %v, %v", info, err)
	}
	if err := bundleNatives(ctx, c, dst, platform); err == nil {
		t.Fatal("replaced an existing library")
	}
	// Reject collisions with dependency libraries before any download is attempted.
	put(t, filepath.Join(root, "mygo-plugin.json"), `{"libraries":[{"name":"upstream","files":{"linux-amd64":{"name":"LibAudio.so","url":"https://invalid.example/native","sha256":"`+strings.Repeat("a", 64)+`"}}}]}`)
	if err := bundleNatives(ctx, c, filepath.Join(root, "collision"), platform); err == nil || !strings.Contains(err.Error(), "conflicting native filename") {
		t.Fatalf("expected a runtime/manifest collision, got %v", err)
	}
}
