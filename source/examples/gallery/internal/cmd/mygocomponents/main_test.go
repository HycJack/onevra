package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"widgets.go": `package ui
func Button(c *Context, label string) *Element { return nil }
func Radio[T comparable](c *Context, value T) *Element { return nil }
func AlertDialog(c *Context) int { return 0 }
func Overlay(c *Context) {}
func ButtonBase(c *Context) *Element { return nil }
func private(c *Context) {}
func Helper(c Context) {}
func Other(c *Elsewhere) {}
func NewTester(view func(*Context)) {}
func (c *Context) Method() {}
`,
		"future_windows.go": "//go:build windows\n\npackage ui\nfunc Future(c *Context) {}\nfunc Button(c *Context) {}\n",
		"widgets_test.go":   "package ui\nfunc TestOnly(c *Context) {}\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range got {
		names = append(names, entry.Name)
		if entry.Source == "" {
			t.Fatal("missing source")
		}
	}
	want := []string{"AlertDialog", "Button", "Future", "Overlay", "Radio"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}
	again, err := scan(dir)
	if err != nil || !reflect.DeepEqual(got, again) {
		t.Fatalf("unstable scan: %v", err)
	}
}

func TestScanErrors(t *testing.T) {
	if _, err := scan(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing directory accepted")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.go"), []byte("package ui\nfunc"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := scan(dir); err == nil {
		t.Fatal("invalid source accepted")
	}
}
