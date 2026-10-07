// Command prepare materializes the pinned MyGo integration dependency.
// Run once after checkout and whenever its patch changes, from source/:
//
//	go run ./internal/cmd/prepare
//
// Gallery, tests and MyGo packaging then use ordinary Go commands, without tags.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oneo-me/oneo-go/internal/mygopatch"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 1 {
		return fmt.Errorf("usage: go run ./internal/cmd/prepare")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	mod, err := os.ReadFile(filepath.Join(cwd, "go.mod"))
	if err != nil {
		return fmt.Errorf("run from source/: %w", err)
	}
	if !strings.Contains(string(mod), mygopatch.Module+" "+mygopatch.Version) {
		return fmt.Errorf("patch requires %s %s", mygopatch.Module, mygopatch.Version)
	}
	if flags := os.Getenv("GOFLAGS"); strings.Contains(flags, "-overlay") || strings.Contains(flags, "-modfile") {
		return fmt.Errorf("prepare requires the repository module, without an overlay or modfile")
	}
	// An explicit version downloads the original module, bypassing go.mod's
	// local replacement. No MyGo imports are needed to bootstrap this command.
	cmd := exec.Command("go", "mod", "download", "-json", mygopatch.Module+"@"+mygopatch.Version)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	data, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("download pinned MyGo: %w\n%s", err, data)
	}
	var module struct{ Dir, Error string }
	if err = json.Unmarshal(data, &module); err != nil {
		return err
	}
	if module.Error != "" || module.Dir == "" {
		return fmt.Errorf("MyGo download: %s", module.Error)
	}
	if err = prepare(module.Dir, filepath.Join(cwd, ".mygo"), mygopatch.Patch); err != nil {
		return err
	}
	fmt.Println("MyGo prepared. Run the existing gallery with: go tool mygo dev")
	return nil
}

func prepare(base, cache string, patch []byte) error {
	if err := os.MkdirAll(cache, 0700); err != nil {
		return err
	}
	lock := filepath.Join(cache, "prepare.lock")
	if err := os.Mkdir(lock, 0700); err != nil {
		return fmt.Errorf("cannot lock preparation (if interrupted, remove %s): %w", lock, err)
	}
	defer os.Remove(lock)
	target := filepath.Join(cache, "mygo")
	stamp := fmt.Sprintf("%s@%s\n%x\n", mygopatch.Module, mygopatch.Version, sha256.Sum256(patch))
	if previous, err := os.ReadFile(filepath.Join(target, ".oneo-patch")); err == nil && string(previous) == stamp {
		return nil
	}
	work, err := os.MkdirTemp(cache, "prepare-")
	if err != nil {
		return err
	}
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(work)
		}
	}()
	next := filepath.Join(work, "next")
	if err = copyModule(base, next); err != nil {
		return err
	}
	if _, err = mygopatch.Apply(next, next, patch); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(next, ".oneo-patch"), []byte(stamp), 0600); err != nil {
		return err
	}
	// Publish only a complete patched tree. Preserve the previous dependency
	// if patch validation or installation fails. Do not prepare during a build.
	backup := filepath.Join(work, "previous")
	hadPrevious := false
	if _, err = os.Lstat(target); err == nil {
		if err = os.Rename(target, backup); err != nil {
			return err
		}
		hadPrevious = true
	} else if !os.IsNotExist(err) {
		return err
	}
	if err = os.Rename(next, target); err != nil {
		if hadPrevious {
			if restore := os.Rename(backup, target); restore != nil {
				// Keep the backup for recovery rather than deleting it on return.
				cleanup = false
				return fmt.Errorf("install: %v; restore: %v; previous dependency saved at %s", err, restore, backup)
			}
		}
		return err
	}
	return nil
}

func copyModule(from, to string) error {
	return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}
		target := filepath.Join(to, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0700)
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("unexpected module entry %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm()|0200)
	})
}
