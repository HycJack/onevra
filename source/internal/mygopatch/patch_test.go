package mygopatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	base, dst := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(base, "file.go"), []byte("one\ntwo\nthree\nfour\nfive\n"), 0600); err != nil {
		t.Fatal(err)
	}
	patch := []byte("--- a/file.go\n+++ b/file.go\n@@ -1,2 +1,2 @@\n one\n-two\n+changed\n@@ -4,2 +4,3 @@\n four\n+inserted\n five\n--- /dev/null\n+++ b/sub/new.go\n@@ -0,0 +1 @@\n+new\n")
	files, err := Apply(base, dst, patch)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("files: %v", files)
	}
	for name, want := range map[string]string{"file.go": "one\nchanged\nthree\nfour\ninserted\nfive\n", "sub/new.go": "new\n"} {
		got, err := os.ReadFile(filepath.Join(dst, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Errorf("%s: %q", name, got)
		}
	}
	got, _ := os.ReadFile(filepath.Join(base, "file.go"))
	if strings.Contains(string(got), "changed") {
		t.Fatal("modified dependency cache")
	}
	if _, err := Apply(base, t.TempDir(), []byte(strings.Replace(string(patch), "-two", "-mismatch", 1))); err == nil {
		t.Fatal("accepted context mismatch")
	}
}

func TestRejectMalformedPatches(t *testing.T) {
	for _, patch := range []string{
		"--- /dev/null\n+++ b/../escape\n@@ -0,0 +1 @@\n+bad\n",
		"--- /dev/null\n+++ b/C:\\escape\n@@ -0,0 +1 @@\n+bad\n",
		"--- /dev/null\n+++ b/new\n@@ -0,0 +1 @@\n",
		"--- /dev/null\n+++ b/new\n@@ -0,0 +1 @@\n?bad\n",
		"--- /dev/null\n+++ b/new\n@@ -0,0 +1 @@\n-context\n",
		"--- /dev/null\n+++ b/new\n@@ -1,2,3 +1 @@\n+bad\n",
		"--- /dev/null\n+++ b/new\n@@ -0,0 +1 @@\n+one\n--- /dev/null\n+++ b/new\n@@ -0,0 +1 @@\n+two\n",
	} {
		if _, err := Apply(t.TempDir(), t.TempDir(), []byte(patch)); err == nil {
			t.Errorf("accepted %q", patch)
		}
	}
}

func TestPinnedPatchShape(t *testing.T) {
	if len(Patch) == 0 || !strings.Contains(string(Patch), "+++ b/ui/native.go") || !strings.Contains(string(Patch), "+++ b/internal/gpu/d3d11/composition.go") {
		t.Fatal("missing pinned integration")
	}
}
