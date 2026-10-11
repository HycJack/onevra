package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Exercise a real application-built Swift dylib without patching the bundle or
// archive afterwards. The loader uses the relocated copy, not a developer cache.
func TestNativeLibraryPackaging(t *testing.T) {
	if runtime.GOOS != "darwin" || os.Getenv("ONEVRA_PACKER_SMOKE") != "1" {
		t.Skip("set ONEVRA_PACKER_SMOKE=1 on macOS with Xcode")
	}
	ctx := context.Background()
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/localfixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), `package main
import (
    "os"
    "path/filepath"
)
func main() {
    exe, err := os.Executable()
    if err != nil { panic(err) }
    if _, err := os.Stat(filepath.Join(filepath.Dir(exe), "..", "Frameworks", "libaudio.dylib")); err != nil { panic(err) }
}
`)
	put(t, filepath.Join(root, "native", "Audio.swift"), `@_cdecl("fixture_value")
public func fixtureValue() -> Int32 { 42 }
`)
	library := filepath.Join(root, "runtimes", "darwin", runtime.GOARCH, "libaudio.dylib")
	if err := os.MkdirAll(filepath.Dir(library), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := command(ctx, "xcrun", "swiftc", "-emit-library", "-O", "-module-name", "FixtureAudio",
		"-Xlinker", "-install_name", "-Xlinker", "@rpath/libaudio.dylib",
		filepath.Join(root, "native", "Audio.swift"), "-o", library); err != nil {
		t.Fatal(err)
	}
	if _, err := command(ctx, "codesign", "--remove-signature", library); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(library)
	if err != nil {
		t.Fatal(err)
	}
	obsolete := filepath.Join(filepath.Dir(library), "libobsolete.dylib")
	put(t, obsolete, string(original))
	put(t, filepath.Join(root, "native", "loader.c"), `#include <dlfcn.h>
#include <stdio.h>
int main(int argc, char **argv) {
    if (argc != 2) return 1;
    void *lib = dlopen(argv[1], RTLD_NOW | RTLD_LOCAL);
    if (!lib) { fprintf(stderr, "%s\n", dlerror()); return 2; }
    int (*value)(void) = (int (*)(void))dlsym(lib, "fixture_value");
    int result = value && value() == 42 ? 0 : 3;
    dlclose(lib);
    return result;
}
`)
	loader := filepath.Join(root, "loader")
	if _, err := command(ctx, "xcrun", "clang", filepath.Join(root, "native", "loader.c"), "-o", loader); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "packer.yaml")
	put(t, file, fmt.Sprintf(`name: Local Audio
identifier: me.oneo.localfixture
out: dist
openReleases: false
platforms: [darwin/%s/zip]
`, runtime.GOARCH))
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if err := run(ctx, []string{"build"}); err != nil {
		t.Fatal(err)
	}
	extracted := filepath.Join(root, "relocated output")
	if _, err := command(ctx, "ditto", "-x", "-k", filepath.Join(c.releaseDir(), c.archiveName(hostTarget())+".zip"), extracted); err != nil {
		t.Fatal(err)
	}
	if err := run(ctx, []string{"run"}); err != nil {
		t.Fatal(err)
	}
	devApp := filepath.Join(c.path(c.Out), c.Name+" Dev.app")
	for _, app := range []string{filepath.Join(c.releaseDir(), c.Name+".app"), filepath.Join(extracted, c.Name+".app"), devApp} {
		bundled := filepath.Join(app, "Contents", "Frameworks", "libaudio.dylib")
		for _, path := range []string{bundled, app} {
			if _, err := command(ctx, "codesign", "--verify", "--deep", "--strict", path); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := command(ctx, loader, bundled); err != nil {
			t.Fatalf("load bundled Swift library: %v", err)
		}
	}
	if data, err := os.ReadFile(library); err != nil || !bytes.Equal(data, original) {
		t.Fatalf("packaging changed the application-owned library: %v", err)
	}
	if _, err := os.Stat(filepath.Join(devApp, "Contents", "Frameworks", "libobsolete.dylib")); err != nil {
		t.Fatal("did not discover all runtime libraries:", err)
	}
	if err := os.Remove(obsolete); err != nil {
		t.Fatal(err)
	}
	if err := run(ctx, []string{"run"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(devApp, "Contents", "Frameworks", "libobsolete.dylib")); !os.IsNotExist(err) {
		t.Fatal("rebuild retained a stale runtime library")
	}
	// An invalid runtime entry must not replace the last usable development app.
	if err := os.Remove(library); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(library, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := run(ctx, []string{"run"}); err == nil {
		t.Fatal("accepted a directory as a native library")
	}
	if _, err := command(ctx, loader, filepath.Join(devApp, "Contents", "Frameworks", "libaudio.dylib")); err != nil {
		t.Fatalf("failed build damaged the previous development library: %v", err)
	}
}
