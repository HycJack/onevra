package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tc-hib/winres"
)

// Patch only the staged executable; do not write .syso files into the application's
// source tree. Preserve existing resources, including an application-owned manifest.
func windowsResources(c *config, exe string) error {
	src, err := os.Open(exe)
	if err != nil {
		return err
	}
	defer src.Close()
	rs, err := winres.LoadFromEXE(src)
	if err != nil && !errors.Is(err, winres.ErrNoResources) {
		return err
	}
	hasManifest := false
	rs.WalkType(winres.RT_MANIFEST, func(_ winres.Identifier, _ uint16, _ []byte) bool {
		hasManifest = true
		return false
	})
	if !hasManifest {
		rs.SetManifest(winres.AppManifest{
			Description: c.Name, Compatibility: winres.Win10AndAbove,
			DPIAwareness: winres.DPIPerMonitorV2, LongPathAware: true,
		})
	}
	if c.Windows.Icon != "" {
		f, err := os.Open(c.path(c.Windows.Icon))
		if err != nil {
			return err
		}
		icon, err := winres.LoadICO(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("Windows icon: %w", err)
		}
		if err := rs.SetIcon(winres.ID(1), icon); err != nil {
			return err
		}
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}
	path := exe + ".resources"
	defer os.Remove(path)
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	err = rs.WriteToEXE(dst, src)
	if closeErr := dst.Close(); err == nil {
		err = closeErr
	}
	if closeErr := src.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(path, exe)
}
