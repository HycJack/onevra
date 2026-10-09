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
		t, _, err := parsePlatform(value)
		if err != nil {
			return nil, err
		}
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

// A missing format selects the platform's default archive formats.
func parsePlatform(value string) (target, string, error) {
	parts := strings.Split(strings.TrimSpace(value), "/")
	if len(parts) < 2 || len(parts) > 3 {
		return target{}, "", fmt.Errorf("invalid platform %q: use OS/ARCH[/FORMAT]", value)
	}
	t := target{parts[0], parts[1]}
	if (t.OS != "darwin" && t.OS != "windows" && t.OS != "linux") || (t.Arch != "arm64" && t.Arch != "amd64") {
		return target{}, "", fmt.Errorf("invalid platform %q: use darwin, windows or linux with /arm64 or /amd64", value)
	}
	if len(parts) == 2 {
		return t, "", nil
	}
	for _, suffix := range defaultArchiveSuffixes(t) {
		if parts[2] == strings.TrimPrefix(suffix, ".") {
			return t, parts[2], nil
		}
	}
	return target{}, "", fmt.Errorf("unsupported archive format in %q", value)
}

func defaultArchiveSuffixes(t target) []string {
	switch t.OS {
	case "darwin":
		return []string{".dmg", ".zip"}
	case "windows":
		return []string{".zip"}
	default:
		return []string{".tar.gz"}
	}
}

func (c *config) archiveSuffixes(t target) []string {
	var suffixes []string
	seen := map[string]bool{}
	for _, value := range c.Platforms {
		selected, format, err := parsePlatform(value)
		if err != nil || selected != t {
			continue
		}
		formats := defaultArchiveSuffixes(t)
		if format != "" {
			formats = []string{"." + format}
		}
		for _, suffix := range formats {
			if !seen[suffix] {
				seen[suffix] = true
				suffixes = append(suffixes, suffix)
			}
		}
	}
	if len(suffixes) == 0 {
		return defaultArchiveSuffixes(t)
	}
	return suffixes
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
