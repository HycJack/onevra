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
	put(t, filepath.Join(root, "icon.png"), "icon source, not payload")
	if err := copyResources(src, dst); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dst, "nested", "message.txt")); err != nil || string(data) != "payload" {
		t.Fatalf("copied resource: %q, %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(dst, "icon.png")); !os.IsNotExist(err) {
		t.Fatal("copied icon source as payload")
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
