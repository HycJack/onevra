package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (c *config) iconSources() []string {
	var paths []string
	for _, icon := range []string{c.MacOS.Icon, c.Windows.Icon, c.Linux.Icon} {
		if icon != "" {
			paths = append(paths, c.path(icon))
		}
	}
	return paths
}

// Copy payloads without icon sources or overwriting generated files.
// Payload symlinks and special files cannot introduce external dependencies.
func copyResources(src, dst string, icons ...string) error {
	info, err := os.Lstat(src)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("resources must be a directory: %s", src)
	}
	rel, err := filepath.Rel(src, dst)
	if err != nil {
		return err
	}
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("resource destination cannot be inside its source: %s", dst)
	}
	excluded := make(map[string]bool)
	for _, name := range []string{"icon.icon", "icon.icns", "icon.ico", "icon.png", "icon.svg"} {
		excluded[filepath.Join(src, name)] = true
	}
	for _, icon := range icons {
		excluded[filepath.Clean(icon)] = true
	}
	if err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != src && excluded[path] {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() && !entry.Type().IsRegular() {
			return fmt.Errorf("resource must be a regular file or directory: %s", path)
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyResourceFile(path, target)
	}); err != nil {
		return fmt.Errorf("copy resources: %w", err)
	}
	return nil
}

func copyResourceFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666|info.Mode().Perm()&0o111)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	return errors.Join(copyErr, out.Close())
}
