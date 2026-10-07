package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func put(t *testing.T, path, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestConfig(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "packer.json")
	put(t, file, `{"name":"Test & App","identifier":"me.test.app","macos":{"icon":"icon.icns"}}`)
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	if c.root != root || c.MacOS.Icon != "icon.icns" || c.Version != "0.0.0" || c.MacOS.MinimumSystemVersion != "13.0" {
		t.Fatalf("unexpected config: %+v", c)
	}
	if c.path(c.MacOS.Icon) != filepath.Join(root, "icon.icns") || c.path(file) != file {
		t.Fatal("paths are not relative to the configuration file")
	}
	for _, input := range []string{
		`{"name":"../App","identifier":"me.app"}`,
		`{"name":"App","identifier":""}`,
		`{"name":"App","identifier":"me.app","resources":["files"]}`,
		`{"name":"App","identifier":"me.app","nativeLibraries":{}}`,
		`{"name":"App","identifier":"me.app","macos":{"signingIdentity":"Developer ID"}}`,
		`{"name":"App","identifier":"me.app","main":"-bad"}`,
		`{"name":"App","identifier":"me.app","version":"bad\nversion"}`,
		`{"name":"App","identifier":"me.app","version":"../1.0"}`,
		`{"name":"App","identifier":"me.app","version":"1\\2"}`,
		`{"name":"App","identifier":"me.app","version":"1:2"}`,
		`{"name":"App","identifier":"me.app"} {}`,
		`null`,
	} {
		put(t, file, input)
		if _, err := loadConfig(file); err == nil {
			t.Errorf("accepted %s", input)
		}
	}
}

func TestBuildModes(t *testing.T) {
	c := &config{Name: "Test App", Identifier: "me.test", Version: "1.0", Main: "./cmd/app"}
	fast := strings.Join(c.buildArgs("app", false, target{"darwin", "arm64"}), " ")
	small := strings.Join(c.buildArgs("app", true, target{"darwin", "arm64"}), " ")
	for _, flag := range []string{"-trimpath", "-s -w"} {
		if strings.Contains(fast, flag) || !strings.Contains(small, flag) {
			t.Errorf("incorrect mode flags: run=%s build=%s", fast, small)
		}
	}
	for _, flag := range []string{"-race", "-gcflags", "-a "} {
		if strings.Contains(fast, flag) {
			t.Errorf("run should use ordinary cached Go compilation: %s", fast)
		}
	}
	if fast != "build -o app ./cmd/app" || small != "build -o app -trimpath -ldflags -s -w ./cmd/app" {
		t.Fatalf("unexpected framework-specific build flags: run=%s build=%s", fast, small)
	}
}

func TestPublish(t *testing.T) {
	root := t.TempDir()
	dst := filepath.Join(root, "output")
	put(t, filepath.Join(dst, "old"), "previous")
	if err := publish(filepath.Join(root, "missing"), dst); err == nil {
		t.Fatal("expected failed publication")
	}
	if data, _ := os.ReadFile(filepath.Join(dst, "old")); string(data) != "previous" {
		t.Fatal("failed publication lost previous build")
	}
	src := filepath.Join(root, "new")
	put(t, filepath.Join(src, "app"), "new")
	if err := publish(src, dst); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dst, "old")); !os.IsNotExist(err) {
		t.Fatal("stale artifact survived publication")
	}
	put(t, src, "file")
	if err := publish(dst, src); err == nil {
		t.Fatal("replaced non-directory destination")
	}
}

// This opt-in test runs real Apple tools and briefly opens Finder to save the
// DMG layout. A tiny Go app verifies packaging without downloading native libraries.
func TestNativePackaging(t *testing.T) {
	if runtime.GOOS != "darwin" || os.Getenv("ONEO_PACKER_SMOKE") != "1" {
		t.Skip("set ONEO_PACKER_SMOKE=1 on macOS with Xcode")
	}
	ctx := context.Background()
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), `package main
import "os"
func main() { if path := os.Getenv("PACKER_TEST_MARKER"); path != "" { if err := os.WriteFile(path, []byte("ran"), 0600); err != nil { panic(err) } } }
`)
	file := filepath.Join(root, "packer.json")
	put(t, file, `{"name":"Test & App","identifier":"me.oneo.fixture","out":"dist","macos":{"icon":"icon.icns"}}`)
	put(t, filepath.Join(root, "icon.icns"), string(testICNS(t)))
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "resources", "nested", "payload.txt"), "bundled payload")
	dmg, err := build(ctx, c, c.releaseDir(), true, hostTarget())
	if err != nil {
		t.Fatal(err)
	}
	if dmg != filepath.Join(c.releaseDir(), "test-app-darwin-"+runtime.GOARCH+"-0.0.0.dmg") {
		t.Fatalf("unexpected archive path: %s", dmg)
	}
	published, err := os.ReadDir(c.releaseDir())
	if err != nil || len(published) != 2 {
		t.Fatalf("release should retain only two archives: %v, %v", published, err)
	}
	base := strings.TrimSuffix(dmg, ".dmg")
	extracted := filepath.Join(root, "extracted")
	if _, err := command(ctx, "ditto", "-x", "-k", base+".zip", extracted); err != nil {
		t.Fatal(err)
	}
	app := filepath.Join(extracted, c.Name+".app")
	if data, err := os.ReadFile(filepath.Join(app, "Contents", "Resources", "nested", "payload.txt")); err != nil || string(data) != "bundled payload" {
		t.Fatalf("archive payload: %q, %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(app, "Contents", "Resources", "icon.icns")); !os.IsNotExist(err) {
		t.Fatal("icon source was copied into the release")
	}
	data, err := command(ctx, "plutil", "-convert", "json", "-o", "-", filepath.Join(app, "Contents", "Info.plist"))
	if err != nil {
		t.Fatal(err)
	}
	var info map[string]any
	if err := json.Unmarshal(data, &info); err != nil || info["CFBundleName"] != c.Name || info["CFBundleIconFile"] != "AppIcon.icns" {
		t.Fatalf("plist: %s, %v", data, err)
	}
	z, err := zip.OpenReader(base + ".zip")
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	found := false
	for _, f := range z.File {
		if f.Name == c.Name+".app/Contents/MacOS/"+c.Name {
			found = true
		}
		if strings.HasSuffix(f.Name, ".zip") || strings.HasSuffix(f.Name, ".dmg") {
			t.Fatalf("archive includes another archive: %s", f.Name)
		}
	}
	if !found {
		t.Fatal("ZIP is missing the app executable")
	}
	if _, err := command(ctx, "hdiutil", "verify", base+".dmg"); err != nil {
		t.Fatal(err)
	}
	mount := filepath.Join(root, "mount")
	if _, err := command(ctx, "hdiutil", "attach", "-quiet", "-readonly", "-nobrowse", "-mountpoint", mount, base+".dmg"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := command(ctx, "hdiutil", "detach", "-quiet", mount); err != nil {
			t.Error(err)
		}
	})
	entries, err := os.ReadDir(mount)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".zip") || strings.HasSuffix(entry.Name(), ".dmg") {
			t.Fatalf("DMG contains another archive: %s", entry.Name())
		}
	}
	for _, name := range []string{".DS_Store", ".background/install.png"} {
		if info, err := os.Stat(filepath.Join(mount, name)); err != nil || info.Size() == 0 {
			t.Fatalf("missing DMG presentation asset %s: %v", name, err)
		}
	}
	layout, err := command(ctx, "osascript", "-e", `on run argv
	set imageFolder to POSIX file (item 1 of argv) as alias
	tell application "Finder"
		set imageFolder to folder imageFolder
		open imageFolder
		delay 1
		tell imageFolder
			set options to icon view options of container window
			set validLayout to (current view of container window is icon view) and (icon size of options is 160) and (arrangement of options is (not arranged)) and (position of item (item 2 of argv) is {180, 185}) and (position of item "Applications" is {540, 185})
		end tell
		close container window of imageFolder
	end tell
	return validLayout
end run`, mount, c.Name+".app")
	if err != nil || strings.TrimSpace(string(layout)) != "true" {
		t.Fatalf("compressed DMG did not retain Finder layout: %s, %v", layout, err)
	}
	if target, err := os.Readlink(filepath.Join(mount, "Applications")); err != nil || target != "/Applications" {
		t.Fatalf("Applications shortcut: %q, %v", target, err)
	}
	if _, err := command(ctx, "codesign", "--verify", "--deep", "--strict", filepath.Join(mount, c.Name+".app")); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "launch marker")
	t.Setenv("PACKER_TEST_MARKER", marker)
	t.Chdir(root)
	if err := run(ctx, []string{"run"}); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(marker); string(data) != "ran" {
		t.Fatal("run did not launch the configured application")
	}
	devApp := filepath.Join(c.path(c.Out), c.Name+" Dev.app")
	data, err = command(ctx, "plutil", "-convert", "json", "-o", "-", filepath.Join(devApp, "Contents", "Info.plist"))
	if err != nil {
		t.Fatal(err)
	}
	var devInfo map[string]any
	if err := json.Unmarshal(data, &devInfo); err != nil || devInfo["CFBundleName"] != c.Name+" Dev" || devInfo["CFBundleIdentifier"] != c.Identifier+".dev" {
		t.Fatalf("development app metadata: %s, %v", data, err)
	}
	// Synthetic Icon Composer fixture, independent of any consumer's artwork.
	icon := filepath.Join(root, "Fixture.icon")
	put(t, filepath.Join(icon, "icon.json"), `{"fill":{"solid":"srgb:0.2,0.4,0.8,1"},"groups":[],"supported-platforms":{"squares":"shared"}}`)
	c.MacOS.Icon = icon
	composerApp, err := build(ctx, c, filepath.Join(root, "composer-output"), false, hostTarget())
	if err != nil {
		t.Fatal(err)
	}
	data, err = command(ctx, "plutil", "-convert", "json", "-o", "-", filepath.Join(composerApp, "Contents", "Info.plist"))
	if err != nil {
		t.Fatal(err)
	}
	var meta map[string]any
	if err := json.Unmarshal(data, &meta); err != nil || meta["CFBundleIconName"] != "Fixture" || meta["CFBundleIdentifier"] != c.Identifier {
		t.Fatalf("missing merged Icon Composer metadata: %s, %v", data, err)
	}
	debug, err := os.Stat(filepath.Join(composerApp, "Contents", "MacOS", c.Name))
	if err != nil {
		t.Fatal(err)
	}
	stripped, err := os.Stat(filepath.Join(app, "Contents", "MacOS", c.Name))
	if err != nil {
		t.Fatal(err)
	}
	if stripped.Size() >= debug.Size() {
		t.Fatalf("build binary (%d) is not smaller than run binary (%d)", stripped.Size(), debug.Size())
	}
	// A failed icon build must not remove the previously published artifacts.
	c.MacOS.Icon = "missing.icon"
	if _, err := build(ctx, c, c.releaseDir(), true, hostTarget()); err == nil {
		t.Fatal("expected missing-icon build failure")
	}
	if _, err := command(ctx, "hdiutil", "verify", dmg); err != nil {
		t.Fatalf("previous archive no longer works: %v", err)
	}
	if err := exec.Command(filepath.Join(devApp, "Contents", "MacOS", c.Name+" Dev")).Run(); err != nil {
		t.Fatalf("development app was not preserved: %v", err)
	}
}
