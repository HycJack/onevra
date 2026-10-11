package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"debug/elf"
	"debug/pe"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/tc-hib/winres"
)

// Cross-build real executables, but use local synthetic native-library payloads.
// This proves archive structure and target selection, not Windows/Linux UI behavior.
func TestCrossPackaging(t *testing.T) {
	if os.Getenv("ONEVRA_PACKER_SMOKE") != "1" {
		t.Skip("set ONEVRA_PACKER_SMOKE=1 to cross-build Windows and Linux fixtures")
	}
	root := t.TempDir()
	put(t, filepath.Join(root, "go.mod"), "module example.com/crossfixture\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "main.go"), "package main\nfunc main() {}\n")
	put(t, filepath.Join(root, "resources", "data", "payload.txt"), "payload")
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	}))
	defer server.Close()
	oldClient := http.DefaultClient
	http.DefaultClient = server.Client()
	t.Cleanup(func() { http.DefaultClient = oldClient })
	targets := []target{{"windows", "amd64"}, {"windows", "arm64"}, {"linux", "amd64"}, {"linux", "arm64"}}
	files := map[string]nativeFile{}
	for _, target := range targets {
		name := "libfixture.so"
		if target.OS == "windows" {
			name = "fixture.dll"
		}
		body := "/" + target.key()
		files[target.key()] = nativeFile{Name: name, URL: server.URL + body, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(body)))}
		put(t, filepath.Join(root, "runtimes", target.OS, target.Arch, "local-"+name), "local-"+target.key())
	}
	put(t, filepath.Join(root, "runtimes", "darwin"), "other targets must not be scanned")
	manifest, err := json.Marshal(map[string]any{"libraries": []any{map[string]any{"name": "fixture", "files": files}}})
	if err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "mygo-plugin.json"), string(manifest))
	icon, err := winres.NewIconFromImages([]image.Image{image.NewNRGBA(image.Rect(0, 0, 16, 16))})
	if err != nil {
		t.Fatal(err)
	}
	var ico bytes.Buffer
	if err := icon.SaveICO(&ico); err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "resources", "icon.ico"), ico.String())
	put(t, filepath.Join(root, "resources", "icon.svg"), `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"/>`)
	file := filepath.Join(root, "packer.yaml")
	put(t, file, `name: "Cross App"
identifier: "me.oneo.crossfixture"
out: "dist"
openReleases: false
platforms: ["windows/amd64","windows/arm64","linux/amd64","linux/arm64"]`)
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if err := run(context.Background(), []string{"build"}); err != nil {
		t.Fatal(err)
	}
	for _, target := range targets {
		t.Run(target.key(), func(t *testing.T) {
			path := filepath.Join(c.releaseDir(), "cross-app-"+target.key()+"-0.0.0"+c.archiveSuffixes(target)[0])
			entries, modes := readPortableArchive(t, path, target.OS == "windows")
			prefix := c.Name + "/"
			if string(entries[prefix+"data/payload.txt"]) != "payload" || string(entries[prefix+files[target.key()].Name]) != "/"+target.key() {
				t.Fatalf("missing payload or wrong native target: %v", entries)
			}
			localName := "local-" + files[target.key()].Name
			if string(entries[prefix+localName]) != "local-"+target.key() {
				t.Fatalf("missing local library or wrong target: %q", entries[prefix+localName])
			}
			executable := filepath.Base(c.executable("unused", target))
			data := entries[prefix+executable]
			if target == hostTarget() {
				retained, err := os.ReadFile(c.executable(filepath.Join(c.releaseDir(), c.Name), target))
				if err != nil || !bytes.Equal(retained, data) {
					t.Fatalf("retained app differs from the host archive: %v", err)
				}
				payload, err := os.ReadFile(filepath.Join(c.releaseDir(), c.Name, files[target.key()].Name))
				if err != nil || string(payload) != "/"+target.key() {
					t.Fatalf("retained native payload = %q, %v", payload, err)
				}
			}
			if target.OS == "windows" {
				f, err := pe.NewFile(bytes.NewReader(data))
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				machine := uint16(pe.IMAGE_FILE_MACHINE_AMD64)
				if target.Arch == "arm64" {
					machine = pe.IMAGE_FILE_MACHINE_ARM64
				}
				if f.Machine != machine || f.OptionalHeader.(*pe.OptionalHeader64).Subsystem != pe.IMAGE_SUBSYSTEM_WINDOWS_GUI {
					t.Fatalf("wrong PE target or subsystem: %+v", f.FileHeader)
				}
				rs, err := winres.LoadFromEXE(bytes.NewReader(data))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := rs.GetIcon(winres.ID(1)); err != nil {
					t.Fatalf("missing embedded ICO: %v", err)
				}
				if len(rs.Get(winres.RT_MANIFEST, winres.ID(1), winres.LCIDDefault)) == 0 {
					t.Fatal("missing Windows manifest")
				}
				app, err := build(context.Background(), c, t.TempDir(), false, target)
				if err != nil {
					t.Fatal(err)
				}
				if data, err := os.ReadFile(filepath.Join(app, localName)); err != nil || string(data) != "local-"+target.key() {
					t.Fatalf("development local library: %q, %v", data, err)
				}
				dev, err := pe.Open(c.executable(app, target))
				if err != nil {
					t.Fatal(err)
				}
				defer dev.Close()
				if dev.Machine != machine || dev.OptionalHeader.(*pe.OptionalHeader64).Subsystem != pe.IMAGE_SUBSYSTEM_WINDOWS_GUI {
					t.Fatal("Windows development executable should use the selected target and GUI subsystem")
				}
			} else {
				f, err := elf.NewFile(bytes.NewReader(data))
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				machine := elf.EM_X86_64
				if target.Arch == "arm64" {
					machine = elf.EM_AARCH64
				}
				if f.Machine != machine || modes[prefix+executable]&0o111 == 0 {
					t.Fatal("wrong ELF target or missing executable mode")
				}
				if len(entries[prefix+"icon.svg"]) == 0 {
					t.Fatal("Linux icon missing")
				}
			}
		})
	}
	entries, err := os.ReadDir(c.releaseDir())
	wantCount := len(targets)
	for _, target := range targets {
		if target == hostTarget() {
			wantCount++
		}
	}
	if err != nil || len(entries) != wantCount {
		t.Fatalf("unexpected release artifacts: %v, %v; want %d", entries, err, wantCount)
	}
	for _, entry := range entries {
		if entry.IsDir() && (entry.Name() != c.Name || wantCount == len(targets)) {
			t.Fatalf("release retained an unexpected directory: %s", entry.Name())
		}
	}
	// Turning retention off also removes the previous host app on the next build.
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintln(f, "\nkeepUncompressed: false"); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := run(context.Background(), []string{"build"}); err != nil {
		t.Fatal(err)
	}
	entries, err = os.ReadDir(c.releaseDir())
	if err != nil || len(entries) != len(targets) {
		t.Fatalf("expected only archives with retention disabled: %v, %v", entries, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			t.Fatalf("retention disabled but directory remains: %s", entry.Name())
		}
	}
}

func readPortableArchive(t *testing.T, path string, windows bool) (map[string][]byte, map[string]int64) {
	t.Helper()
	files, modes := map[string][]byte{}, map[string]int64{}
	if windows {
		z, err := zip.OpenReader(path)
		if err != nil {
			t.Fatal(err)
		}
		defer z.Close()
		for _, f := range z.File {
			if f.FileInfo().IsDir() {
				continue
			}
			r, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			files[f.Name], err = io.ReadAll(r)
			r.Close()
			if err != nil {
				t.Fatal(err)
			}
			modes[f.Name] = int64(f.Mode())
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		gz, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		defer gz.Close()
		r := tar.NewReader(gz)
		for {
			h, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if h.Typeflag != tar.TypeReg {
				continue
			}
			files[h.Name], err = io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			modes[h.Name] = h.Mode
		}
	}
	return files, modes
}
