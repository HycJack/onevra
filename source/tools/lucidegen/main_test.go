package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIdentifier(t *testing.T) {
	for name, want := range map[string]string{"arrow-up-right": "ArrowUpRight", "circle-help": "CircleHelp", "3d-box": "Icon3dBox", "x": "X"} {
		got, err := identifier(name)
		if err != nil || got != want {
			t.Errorf("%s: %q, %v; want %q", name, got, err, want)
		}
	}
	for _, name := range []string{"", "-x", "a--b", "Foo", "a_b", "../x"} {
		if _, err := identifier(name); err == nil {
			t.Errorf("accepted %q", name)
		}
	}
}

func TestPinnedGeneration(t *testing.T) {
	upstream := t.TempDir()
	out := filepath.Join(t.TempDir(), "generated")
	mustWrite := func(path, data string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	mustGit := func(args ...string) string {
		t.Helper()
		result, err := git(upstream, args...)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	if err := os.Mkdir(filepath.Join(upstream, "icons"), 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(filepath.Join(upstream, "LICENSE"), "fixture license\n")
	mustWrite(filepath.Join(upstream, "icons", "search.svg"), `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8"/></svg>`)
	mustGit("init")
	mustGit("add", ".")
	mustGit("-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "commit", "-m", "fixture")
	commit := mustGit("rev-parse", "HEAD")
	if err := run(upstream, commit, out, false); err != nil {
		t.Fatal(err)
	}
	if err := run(upstream, commit, out, true); err != nil {
		t.Fatal(err)
	}
	mustWrite(filepath.Join(out, "icons_gen.go"), "outdated")
	if err := run(upstream, commit, out, true); err == nil || !strings.Contains(err.Error(), "not up to date") {
		t.Fatalf("drift: %v", err)
	}
	if err := run(upstream, strings.Repeat("0", 40), out, false); err == nil {
		t.Fatal("accepted wrong revision")
	}
	mustWrite(filepath.Join(upstream, "LICENSE"), "dirty")
	if err := run(upstream, commit, out, false); err == nil || !strings.Contains(err.Error(), "clean") {
		t.Fatalf("dirty checkout: %v", err)
	}
}

func TestRender(t *testing.T) {
	icons := []icon{{"search", "Search", "<svg viewBox=\"0 0 24 24\"/>\n"}}
	a, at, err := render(icons)
	if err != nil {
		t.Fatal(err)
	}
	b, bt, err := render(icons)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) || !bytes.Equal(at, bt) {
		t.Fatal("nondeterministic generation")
	}
	if !bytes.Contains(a, []byte("var cacheSearch lazyIcon")) || !bytes.Contains(at, []byte("{\"search\", Search}")) {
		t.Fatal("missing accessor or test entry")
	}
}
