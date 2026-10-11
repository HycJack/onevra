package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"unicode"
)

func build(ctx context.Context, c *config, out string, production bool, t target) (string, error) {
	if err := checkHost([]target{t}, runtime.GOOS); err != nil {
		return "", err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return "", err
	}
	stage, err := os.MkdirTemp(out, ".packer-stage-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(stage)
	if t.OS != "darwin" {
		return buildPortable(ctx, c, out, stage, production, t)
	}
	artifacts := filepath.Join(stage, "artifacts")
	app := filepath.Join(artifacts, c.Name+".app")
	contents := filepath.Join(app, "Contents")
	resources := filepath.Join(contents, "Resources")
	for _, p := range []string{filepath.Join(contents, "MacOS"), resources} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return "", err
		}
	}
	// Compile icons first: a missing Xcode toolchain should fail before Go builds.
	info, err := compileIcon(ctx, c, stage, resources)
	if err != nil {
		return "", err
	}
	info, err = c.macOSInfo(info)
	if err != nil {
		return "", err
	}
	if err := copyResources(c.path("resources"), resources, c.iconSources()...); err != nil {
		return "", err
	}
	if err := writePlist(ctx, filepath.Join(contents, "Info.plist"), info); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(contents, "PkgInfo"), []byte("APPL????"), 0o644); err != nil {
		return "", err
	}
	exe := filepath.Join(contents, "MacOS", c.Name)
	if err := compileGo(ctx, c, exe, production, t); err != nil {
		return "", err
	}
	if err := bundleNatives(ctx, c, filepath.Join(contents, "Frameworks"), t); err != nil {
		return "", err
	}
	if _, err := command(ctx, "codesign", "--force", "--sign", "-", app); err != nil {
		return "", err
	}
	if _, err := command(ctx, "codesign", "--verify", "--deep", "--strict", app); err != nil {
		return "", err
	}
	if production {
		if err := archive(ctx, c, stage, app, t); err != nil {
			return "", err
		}
		names := c.releaseNames(t)
		if err := publishArtifacts(artifacts, out, names); err != nil {
			return "", err
		}
		return filepath.Join(out, names[0]), nil
	}
	final := filepath.Join(out, filepath.Base(app))
	if err := publish(app, final); err != nil {
		return "", err
	}
	return final, nil
}

func (c *config) bundleInfo() map[string]any {
	return map[string]any{
		"CFBundleName":                         c.Name,
		"CFBundleDisplayName":                  c.Name,
		"CFBundleIdentifier":                   c.Identifier,
		"CFBundleVersion":                      c.Version,
		"CFBundleShortVersionString":           c.Version,
		"CFBundleExecutable":                   c.Name,
		"CFBundlePackageType":                  "APPL",
		"CFBundleInfoDictionaryVersion":        "6.0",
		"LSMinimumSystemVersion":               c.MacOS.MinimumSystemVersion,
		"NSPrincipalClass":                     "NSApplication",
		"NSHighResolutionCapable":              true,
		"NSSupportsAutomaticGraphicsSwitching": true,
	}
}

func writePlist(ctx context.Context, path string, info map[string]any) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	_, err = command(ctx, "plutil", "-convert", "xml1", path)
	return err
}

func (c *config) slug() string {
	words := strings.FieldsFunc(strings.ToLower(c.Name), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	name := strings.Join(words, "-")
	if name == "" {
		name = "app"
	}
	return name
}

func (c *config) archiveName(t target) string { return c.slug() + "-" + t.key() + "-" + c.Version }

func archive(ctx context.Context, c *config, stage, app string, t target) error {
	// ditto preserves macOS bundle metadata and includes the .app parent.
	base := filepath.Join(filepath.Dir(app), c.archiveName(t))
	if slices.Contains(c.archiveSuffixes(t), ".zip") {
		if _, err := command(ctx, "ditto", "-c", "-k", "--sequesterRsrc", "--keepParent", app, base+".zip"); err != nil {
			return err
		}
	}
	if !slices.Contains(c.archiveSuffixes(t), ".dmg") {
		return nil
	}
	// Isolate the image input, so the ZIP and image never include one another.
	imageRoot := filepath.Join(stage, "image")
	if err := os.Mkdir(imageRoot, 0o755); err != nil {
		return err
	}
	if _, err := command(ctx, "ditto", app, filepath.Join(imageRoot, filepath.Base(app))); err != nil {
		return err
	}
	if err := os.Symlink("/Applications", filepath.Join(imageRoot, "Applications")); err != nil {
		return err
	}
	return createDMG(ctx, stage, imageRoot, c.Name, filepath.Base(app), base+".dmg")
}

// publish replaces a development app only after every build step succeeds.
// A separate backup directory preserves the previous build if rollback fails.
func publish(src, dst string) error {
	info, err := os.Lstat(dst)
	if errors.Is(err, os.ErrNotExist) {
		return os.Rename(src, dst)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to replace non-directory %s", dst)
	}
	backup, err := os.MkdirTemp(filepath.Dir(dst), ".packer-previous-")
	if err != nil {
		return err
	}
	old := filepath.Join(backup, "previous")
	if err := os.Rename(dst, old); err != nil {
		os.Remove(backup)
		return err
	}
	if err := os.Rename(src, dst); err != nil {
		if restore := os.Rename(old, dst); restore != nil {
			return fmt.Errorf("publish: %w; rollback: %v; previous build retained at %s", err, restore, old)
		}
		os.Remove(backup)
		return err
	}
	return os.RemoveAll(backup)
}
