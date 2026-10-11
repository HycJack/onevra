package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// releaseNames includes the uncompressed app only for an exact host match, so
// cross-builds cannot replace it with another architecture's executable.
func (c *config) releaseNames(t target) []string {
	var names []string
	for _, suffix := range c.archiveSuffixes(t) {
		names = append(names, c.archiveName(t)+suffix)
	}
	if c.KeepUncompressed && t == hostTarget() {
		name := c.Name
		if t.OS == "darwin" {
			name += ".app"
		}
		names = append(names, name)
	}
	return names
}

// Publish the named archives and optional app directory as one set, leaving other
// targets untouched. Keep backups until all artifacts are installed for rollback.
func publishArtifacts(src, dst string, names []string) error {
	for _, name := range names {
		info, err := os.Lstat(filepath.Join(src, name))
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() && !info.IsDir() {
			return fmt.Errorf("artifact is not a regular file or directory: %s", name)
		}
		kind := info.Mode().Type()
		info, err = os.Lstat(filepath.Join(dst, name))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err == nil && info.Mode().Type() != kind {
			return fmt.Errorf("refusing to replace artifact of a different type: %s", name)
		}
	}
	backup, err := os.MkdirTemp(dst, ".packer-artifacts-")
	if err != nil {
		return err
	}
	var saved, installed []string
	rollback := func(cause error) error {
		var recovery error
		for _, name := range installed {
			if err := os.RemoveAll(filepath.Join(dst, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
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
