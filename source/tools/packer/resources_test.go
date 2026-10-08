package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResources(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "resources")
	dst := filepath.Join(root, "App.app", "Contents", "Resources")
	if err := copyResources(src, dst); err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(src, "nested", "message.txt"), "payload")
	for _, name := range []string{"icon.icon/icon.json", "icon.icon/Assets/layer.png", "icon.icns", "icon.ico", "icon.png", "icon.svg"} {
		put(t, filepath.Join(src, name), "icon source, not payload")
	}
	put(t, filepath.Join(src, "images", "icon.png"), "image payload")
	if err := copyResources(src, dst); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dst, "nested", "message.txt")); err != nil || string(data) != "payload" {
		t.Fatalf("copied resource: %q, %v", data, err)
	}
	for _, name := range []string{"icon.icon", "icon.icns", "icon.ico", "icon.png", "icon.svg"} {
		if _, err := os.Stat(filepath.Join(dst, name)); !os.IsNotExist(err) {
			t.Fatalf("copied icon source %s as payload", name)
		}
	}
	if data, err := os.ReadFile(filepath.Join(dst, "images", "icon.png")); err != nil || string(data) != "image payload" {
		t.Fatalf("nested image payload: %q, %v", data, err)
	}
	if err := copyResources(src, filepath.Join(src, "output")); err == nil {
		t.Fatal("accepted output inside the resource tree")
	}
	// Compiled icons already in the destination cannot be overwritten by payloads.
	put(t, filepath.Join(dst, "AppIcon.icns"), "compiled")
	other := filepath.Join(root, "collision")
	put(t, filepath.Join(other, "AppIcon.icns"), "payload")
	if err := copyResources(other, dst); err == nil {
		t.Fatal("overwrote generated icon")
	}
	if data, _ := os.ReadFile(filepath.Join(dst, "AppIcon.icns")); string(data) != "compiled" {
		t.Fatal("generated icon changed after collision")
	}
}

func TestConfiguredIconResources(t *testing.T) {
	root := t.TempDir()
	c := &config{root: root}
	c.MacOS.Icon = "resources/branding/custom.icon"
	c.Windows.Icon = filepath.Join(root, "resources", "branding", "custom.ico")
	c.Linux.Icon = "resources/branding/custom.svg"
	for _, name := range []string{"custom.icon/icon.json", "custom.icon/Assets/layer.png", "custom.ico", "custom.svg", "keep.png"} {
		put(t, filepath.Join(root, "resources", "branding", name), name)
	}
	dst := filepath.Join(root, "out")
	if err := copyResources(c.path("resources"), dst, c.iconSources()...); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Join(dst, "branding"))
	if err != nil || len(entries) != 1 || entries[0].Name() != "keep.png" {
		t.Fatalf("payload entries = %v, %v", entries, err)
	}
}

func TestResourceSymlinks(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "resources")
	put(t, filepath.Join(src, "payload"), "data")
	if err := os.Symlink("payload", filepath.Join(src, "link")); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	if err := copyResources(src, filepath.Join(root, "out")); err == nil {
		t.Fatal("accepted symlink within resources")
	}
	linked := filepath.Join(root, "linked")
	if err := os.Symlink(src, linked); err != nil {
		t.Fatal(err)
	}
	if err := copyResources(linked, filepath.Join(root, "out")); err == nil {
		t.Fatal("accepted symlinked resource root")
	}
}
