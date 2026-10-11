package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Discover only the current target's libraries beside packer.yaml. Each level
// must be a real directory; do not follow links outside the runtime tree.
func (c *config) runtimeFiles(t target) ([]string, error) {
	dir := c.root
	for _, part := range []string{"runtimes", t.OS, t.Arch} {
		dir = filepath.Join(dir, part)
		info, err := os.Lstat(dir)
		if os.IsNotExist(err) {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read runtimes: %w", err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("runtime path must be a directory, not a symlink or file: %s", dir)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read runtimes: %w", err)
	}
	var files []string
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if !entry.Type().IsRegular() {
			return nil, fmt.Errorf("runtime library must be a regular file, not a symlink or directory: %s", path)
		}
		if !component(entry.Name()) {
			return nil, fmt.Errorf("invalid runtime filename: %s", path)
		}
		files = append(files, path)
	}
	return files, nil
}

// Treat names case-insensitively even when cross-packaging on a case-sensitive
// filesystem. Neither local libraries nor manifests may replace staged files.
func checkNativeNames(destination string, names []string) error {
	entries, err := os.ReadDir(destination)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	seen := make(map[string]bool, len(entries)+len(names))
	for _, entry := range entries {
		seen[strings.ToLower(entry.Name())] = true
	}
	for _, name := range names {
		if !component(name) {
			return fmt.Errorf("invalid native filename %q", name)
		}
		key := strings.ToLower(name)
		if seen[key] {
			return fmt.Errorf("conflicting native filename %q in %s", name, destination)
		}
		seen[key] = true
	}
	return nil
}
