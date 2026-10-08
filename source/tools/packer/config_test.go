package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlatformIcons(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "packer.yaml")
	for _, config := range []string{
		`name: "App"
identifier: "me.app"
macos:
  icon: "custom.icon"
windows:
  icon: "custom.ico"
linux:
  icon: "custom.svg"`,
		`name: "App"
identifier: "me.app"
macos:
  icon: "custom.icns"
linux:
  icon: "custom.png"`,
	} {
		put(t, file, config)
		if _, err := loadConfig(file); err != nil {
			t.Fatalf("valid per-platform configuration: %v", err)
		}
	}
	for _, config := range []string{
		`name: "App"
identifier: "me.app"
icon: "icon.png"`,
		`name: "App"
identifier: "me.app"
minimumSystemVersion: "13.0"`,
		`name: "App"
identifier: "me.app"
macos:
  icon: "icon.ico"`,
		`name: "App"
identifier: "me.app"
macos:
  icon: "icon.png"`,
		`name: "App"
identifier: "me.app"
windows:
  icon: "icon.png"`,
		`name: "App"
identifier: "me.app"
linux:
  icon: "icon.icon"`,
	} {
		put(t, file, config)
		if _, err := loadConfig(file); err == nil {
			t.Fatalf("accepted invalid platform configuration: %s", config)
		}
	}
}

func TestIconDefaults(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "packer.yaml")
	put(t, file, `name: "App"
identifier: "me.app"`)
	put(t, filepath.Join(root, "resources", "images", "icon.png"), "payload, not the app icon")
	check := func(macos, windows, linux string) {
		t.Helper()
		for _, icon := range []*string{&macos, &windows, &linux} {
			if *icon != "" && *icon != "alternate.icns" {
				*icon = filepath.Join("resources", *icon)
			}
		}
		c, err := loadConfig(file)
		if err != nil {
			t.Fatal(err)
		}
		if err := c.resolveIcons(); err != nil {
			t.Fatal(err)
		}
		if c.MacOS.Icon != macos || c.Windows.Icon != windows || c.Linux.Icon != linux {
			t.Fatalf("icons = %q, %q, %q; want %q, %q, %q", c.MacOS.Icon, c.Windows.Icon, c.Linux.Icon, macos, windows, linux)
		}
	}
	check("", "", "")
	put(t, filepath.Join(root, "resources", "icon.svg"), "SVG source")
	check("", "", "icon.svg")
	put(t, filepath.Join(root, "resources", "icon.png"), "PNG source")
	put(t, filepath.Join(root, "resources", "icon.ico"), "ICO source")
	check("", "icon.ico", "icon.png")
	put(t, filepath.Join(root, "resources", "icon.icns"), "ICNS source")
	check("icon.icns", "icon.ico", "icon.png")
	if err := os.Mkdir(filepath.Join(root, "resources", "icon.icon"), 0o755); err != nil {
		t.Fatal(err)
	}
	check("icon.icon", "icon.ico", "icon.png")
	put(t, file, `name: "App"
identifier: "me.app"
macos:
  icon: "alternate.icns"`)
	check("alternate.icns", "icon.ico", "icon.png")
	if err := os.Remove(filepath.Join(root, "resources", "icon.icon")); err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "resources", "icon.icon"), "not a directory")
	put(t, file, `name: "App"
identifier: "me.app"`)
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.resolveIcons(); err == nil {
		t.Fatal("invalid preferred icon silently fell back to ICNS")
	}
}
