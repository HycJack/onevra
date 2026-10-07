package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Copy payload resources without overwriting compiled icons. Symlinks and
// special files are rejected so a bundle cannot depend on external paths.
func copyResources(src, dst string) error {
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
	if err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && !entry.Type().IsRegular() {
			return fmt.Errorf("resource must be a regular file or directory: %s", path)
		}
		return nil
	}); err != nil {
		return err
	}
	if err := os.CopyFS(dst, os.DirFS(src)); err != nil {
		return fmt.Errorf("copy resources: %w", err)
	}
	return nil
}
