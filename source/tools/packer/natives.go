package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type nativeFile struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// Discover manifests in the selected application's dependency graph without
// importing plugin packages into Packer or running their code.
func bundleNatives(ctx context.Context, c *config, destination string, t target) error {
	files, err := c.nativeFiles(ctx, t)
	if err != nil || len(files) == 0 {
		return err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	for _, f := range files {
		data, err := nativeBytes(ctx, cache, f)
		if err != nil {
			return err
		}
		dst := filepath.Join(destination, f.Name)
		file, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o755)
		if err != nil {
			return err
		}
		_, err = file.Write(data)
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			return err
		}
		// Sign macOS copies only; the shared cache retains its pinned hash.
		if t.OS == "darwin" {
			if _, err := command(ctx, "codesign", "--force", "--sign", "-", dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *config) nativeFiles(ctx context.Context, t target) ([]nativeFile, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-deps", "-f", "{{.Dir}}", c.Main)
	cmd.Dir, cmd.Env = c.root, buildEnv(t)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("list native dependencies: %w", err)
	}
	return readNativeManifests(strings.Split(strings.TrimSpace(string(out)), "\n"), t.key())
}

func readNativeManifests(dirs []string, target string) ([]nativeFile, error) {
	var files []nativeFile
	byLibrary, byFilename := map[string]nativeFile{}, map[string]nativeFile{}
	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			continue
		}
		path := filepath.Join(dir, "mygo-plugin.json")
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		var manifest struct {
			Libraries []struct {
				Name  string                `json:"name"`
				Files map[string]nativeFile `json:"files"`
			} `json:"libraries"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		for _, lib := range manifest.Libraries {
			f, ok := lib.Files[target]
			if !ok {
				return nil, fmt.Errorf("%s: %s has no build for %s", path, lib.Name, target)
			}
			if strings.TrimSpace(lib.Name) == "" || !f.valid() {
				return nil, fmt.Errorf("%s: invalid native library %q", path, lib.Name)
			}
			if prev, ok := byLibrary[lib.Name]; ok && prev != f {
				return nil, fmt.Errorf("%s: conflicting definitions of %s", path, lib.Name)
			}
			byLibrary[lib.Name] = f
			key := strings.ToLower(f.Name)
			if prev, ok := byFilename[key]; ok {
				if prev != f {
					return nil, fmt.Errorf("%s: conflicting native filename %s", path, f.Name)
				}
				continue
			}
			byFilename[key] = f
			files = append(files, f)
		}
	}
	return files, nil
}

func (f nativeFile) valid() bool {
	sum, err := hex.DecodeString(f.SHA256)
	return component(f.Name) && err == nil && len(sum) == sha256.Size &&
		f.SHA256 == strings.ToLower(f.SHA256) && strings.HasPrefix(f.URL, "https://")
}

func nativeBytes(ctx context.Context, cache string, f nativeFile) ([]byte, error) {
	if !f.valid() {
		return nil, fmt.Errorf("invalid native library: %+v", f)
	}
	path := filepath.Join(cache, "onevra-packer", "natives", f.SHA256, f.Name)
	valid := func(data []byte) bool { return fmt.Sprintf("%x", sha256.Sum256(data)) == f.SHA256 }
	if data, err := os.ReadFile(path); err == nil && valid(data) {
		return data, nil
	}
	fmt.Fprintln(os.Stderr, "Downloading", f.Name)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.URL, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: %s", f.Name, res.Status)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, 64<<20+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 64<<20 || !valid(data) {
		return nil, fmt.Errorf("download %s: size limit or SHA-256 mismatch", f.Name)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".download-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	_, err = tmp.Write(data)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return nil, err
	}
	return data, nil
}
