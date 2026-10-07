package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Publish only the named archives, leaving development apps and other targets
// untouched. Keep backups until all files are installed so failures can roll back.
func publishArchives(src, dst string, names []string) error {
	for _, name := range names {
		info, err := os.Lstat(filepath.Join(src, name))
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("archive is not a regular file: %s", name)
		}
		info, err = os.Lstat(filepath.Join(dst, name))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err == nil && !info.Mode().IsRegular() {
			return fmt.Errorf("refusing to replace non-file archive: %s", name)
		}
	}
	backup, err := os.MkdirTemp(dst, ".packer-archives-")
	if err != nil {
		return err
	}
	var saved, installed []string
	rollback := func(cause error) error {
		var recovery error
		for _, name := range installed {
			if err := os.Remove(filepath.Join(dst, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
				recovery = errors.Join(recovery, err)
			}
		}
		for _, name := range saved {
			if err := os.Rename(filepath.Join(backup, name), filepath.Join(dst, name)); err != nil {
				recovery = errors.Join(recovery, err)
			}
		}
		if recovery != nil {
			return fmt.Errorf("publish: %w; rollback: %v; backup retained at %s", cause, recovery, backup)
		}
		return errors.Join(cause, os.RemoveAll(backup))
	}
	for _, name := range names {
		path := filepath.Join(dst, name)
		if err := os.Rename(path, filepath.Join(backup, name)); err == nil {
			saved = append(saved, name)
		} else if !errors.Is(err, os.ErrNotExist) {
			return rollback(err)
		}
		if err := os.Rename(filepath.Join(src, name), path); err != nil {
			return rollback(err)
		}
		installed = append(installed, name)
	}
	return os.RemoveAll(backup)
}
