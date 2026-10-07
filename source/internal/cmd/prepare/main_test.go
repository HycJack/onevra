package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareDependency(t *testing.T) {
	base, cache := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "file.go"), []byte("original\n"), 0600); err != nil {
		t.Fatal(err)
	}
	patch := []byte("--- a/file.go\n+++ b/file.go\n@@ -1 +1 @@\n-original\n+patched\n")
	if err := prepare(base, cache, patch); err != nil {
		t.Fatal(err)
	}
	check := func(path, want string) {
		t.Helper()
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Fatalf("%s: %q, %v", path, got, err)
		}
	}
	target := filepath.Join(cache, "mygo", "file.go")
	check(target, "patched\n")
	check(filepath.Join(base, "file.go"), "original\n")
	before, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if err := prepare(base, cache, patch); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) {
		t.Fatal("unchanged patch rebuilt the dependency")
	}
	bad := []byte("--- a/file.go\n+++ b/file.go\n@@ -1 +1 @@\n-mismatch\n+bad\n")
	if err := prepare(base, cache, bad); err == nil {
		t.Fatal("accepted context mismatch")
	}
	check(target, "patched\n")
	next := []byte("--- a/file.go\n+++ b/file.go\n@@ -1 +1 @@\n-original\n+updated\n")
	if err := prepare(base, cache, next); err != nil {
		t.Fatal(err)
	}
	check(target, "updated\n")
	entries, err := os.ReadDir(cache)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "mygo" {
		t.Fatalf("temporary state leaked: %v", entries)
	}
}

func TestPrepareRejectsConcurrentWriter(t *testing.T) {
	cache := t.TempDir()
	if err := os.Mkdir(filepath.Join(cache, "prepare.lock"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := prepare(t.TempDir(), cache, nil); err == nil {
		t.Fatal("ignored writer lock")
	}
}
