package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

func putManifest(t *testing.T, dir, name string, f nativeFile) {
	t.Helper()
	data, err := json.Marshal(map[string]any{"libraries": []any{
		map[string]any{"name": name, "source": "upstream provenance", "files": map[string]nativeFile{
			"darwin-arm64": f, "darwin-amd64": f,
		}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(dir, "mygo-plugin.json"), string(data))
}

func TestNativeManifests(t *testing.T) {
	root := t.TempDir()
	a, b := filepath.Join(root, "a"), filepath.Join(root, "b")
	f := nativeFile{Name: "libfixture.dylib", URL: "https://example.com/native", SHA256: strings.Repeat("a", 64)}
	if files, err := readNativeManifests([]string{root, ""}, "darwin-arm64"); err != nil || len(files) != 0 {
		t.Fatalf("absent manifest: %v, %v", files, err)
	}
	putManifest(t, a, "fixture", f)
	putManifest(t, b, "fixture", f)
	files, err := readNativeManifests([]string{a, b, a}, "darwin-arm64")
	if err != nil || len(files) != 1 || files[0] != f {
		t.Fatalf("selection and deduplication: %v, %v", files, err)
	}
	if _, err := readNativeManifests([]string{a}, "linux-arm64"); err == nil {
		t.Fatal("missing platform silently omitted a library")
	}
	for _, tc := range []struct {
		name string
		file nativeFile
	}{
		{"", f},
		{"fixture", nativeFile{Name: "../escape", URL: f.URL, SHA256: f.SHA256}},
		{"fixture", nativeFile{Name: f.Name, URL: f.URL, SHA256: "bad"}},
		{"fixture", nativeFile{Name: "different.dylib", URL: f.URL, SHA256: f.SHA256}},
		{"other", nativeFile{Name: f.Name, URL: f.URL, SHA256: strings.Repeat("b", 64)}},
		{"other", nativeFile{Name: "LibFixture.dylib", URL: f.URL, SHA256: f.SHA256}},
	} {
		putManifest(t, b, tc.name, tc.file)
		if _, err := readNativeManifests([]string{a, b}, "darwin-arm64"); err == nil {
			t.Fatalf("accepted invalid or conflicting manifest: %+v", tc)
		}
	}
	put(t, filepath.Join(b, "mygo-plugin.json"), "invalid JSON")
	if _, err := readNativeManifests([]string{b}, "darwin-arm64"); err == nil {
		t.Fatal("accepted malformed manifest")
	}
}

func TestNativeDependencyDiscovery(t *testing.T) {
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/fixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), "package main\nimport _ \"example.com/fixture/plugin\"\nfunc main() {}\n")
	plugin := filepath.Join(root, "plugin")
	put(t, filepath.Join(plugin, "plugin.go"), "package plugin\n")
	f := nativeFile{Name: "libfixture.dylib", URL: "https://example.com/native", SHA256: strings.Repeat("a", 64)}
	putManifest(t, plugin, "fixture", f)
	// An unrelated manifest must not be discovered by a recursive directory scan.
	put(t, filepath.Join(root, "unused", "mygo-plugin.json"), "invalid JSON")
	c := &config{root: root, Main: "."}
	files, err := c.nativeFiles(context.Background(), target{"darwin", "arm64"})
	if err != nil || len(files) != 1 || files[0] != f {
		t.Fatalf("dependency discovery: %v, %v", files, err)
	}
	c.Main = "./missing"
	if _, err := c.nativeFiles(context.Background(), target{"darwin", "arm64"}); err == nil {
		t.Fatal("go list failure must not produce an incomplete bundle")
	}
}

func TestNativeCache(t *testing.T) {
	body := "native library fixture"
	var requests atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, body)
	}))
	defer server.Close()
	client := http.DefaultClient
	http.DefaultClient = server.Client()
	t.Cleanup(func() { http.DefaultClient = client })
	cache := t.TempDir()
	f := nativeFile{Name: "libfixture.dylib", URL: server.URL, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(body)))}
	for range 2 {
		data, err := nativeBytes(context.Background(), cache, f)
		if err != nil || string(data) != body {
			t.Fatalf("nativeBytes: %q, %v", data, err)
		}
	}
	if requests.Load() != 1 {
		t.Fatalf("cache missed: %d requests", requests.Load())
	}
	path := filepath.Join(cache, "onevra-packer", "natives", f.SHA256, f.Name)
	put(t, path, "corrupted")
	if _, err := nativeBytes(context.Background(), cache, f); err != nil || requests.Load() != 2 {
		t.Fatalf("corrupt cache not repaired: %v, %d requests", err, requests.Load())
	}
	f.SHA256 = strings.Repeat("0", 64)
	if _, err := nativeBytes(context.Background(), cache, f); err == nil {
		t.Fatal("accepted incorrect hash")
	}
	if _, err := os.Stat(filepath.Join(cache, "onevra-packer", "natives", f.SHA256, f.Name)); !os.IsNotExist(err) {
		t.Fatal("cached unverified bytes")
	}
	for _, name := range []string{"../escape", "", "..", "/absolute"} {
		f.Name = name
		if _, err := nativeBytes(context.Background(), cache, f); err == nil {
			t.Errorf("accepted native filename %q", name)
		}
	}
}
