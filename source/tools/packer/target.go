package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type target struct{ OS, Arch string }

func hostTarget() target     { return target{runtime.GOOS, runtime.GOARCH} }
func (t target) key() string { return t.OS + "-" + t.Arch }

func parseTargets(values []string) ([]target, error) {
	if len(values) == 0 {
		values = []string{runtime.GOOS + "/" + runtime.GOARCH}
	}
	var targets []target
	seen := map[target]bool{}
	for _, value := range values {
		goos, arch, ok := strings.Cut(strings.TrimSpace(value), "/")
		if !ok || (goos != "darwin" && goos != "windows" && goos != "linux") || (arch != "arm64" && arch != "amd64") {
			return nil, fmt.Errorf("invalid platform %q: use darwin, windows or linux with /arm64 or /amd64", value)
		}
		t := target{goos, arch}
		if !seen[t] {
			seen[t] = true
			targets = append(targets, t)
		}
	}
	return targets, nil
}

func (c *config) targets(mode string) ([]target, error) {
	if mode == "run" {
		return parseTargets(nil)
	}
	return parseTargets(c.Platforms)
}

func checkHost(targets []target, hostOS string) error {
	for _, t := range targets {
		if t.OS == "darwin" && hostOS != "darwin" {
			return fmt.Errorf("%s requires macOS for icon compilation and signing", t.key())
		}
	}
	return nil
}

func buildEnv(t target) []string {
	return append(os.Environ(), "CGO_ENABLED=0", "GOOS="+t.OS, "GOARCH="+t.Arch)
}

func (c *config) archiveSuffixes(t target) []string {
	switch t.OS {
	case "darwin":
		if c.MacOS.DisableDMG {
			return []string{".zip"}
		}
		return []string{".dmg", ".zip"}
	case "windows":
		return []string{".zip"}
	default:
		return []string{".tar.gz"}
	}
}

func (c *config) executable(app string, t target) string {
	switch t.OS {
	case "darwin":
		return filepath.Join(app, "Contents", "MacOS", c.Name)
	case "windows":
		return filepath.Join(app, c.Name+".exe")
	default:
		return filepath.Join(app, c.slug())
	}
}
