package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileGo(ctx context.Context, c *config, exe string, release bool, t target) error {
	cmd := exec.CommandContext(ctx, "go", c.buildArgs(exe, release, t)...)
	cmd.Dir, cmd.Env = c.root, buildEnv(t)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build %s: %w", t.key(), err)
	}
	return nil
}

// Windows and Linux ship a relocatable directory inside one archive. Libraries
// and payloads live beside the executable, where the native loaders expect them.
func buildPortable(ctx context.Context, c *config, out, stage string, release bool, t target) (string, error) {
	if t.OS == "windows" {
		base, _, _ := strings.Cut(strings.ToUpper(c.Name), ".")
		reserved := base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" ||
			(len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && base[3] >= '1' && base[3] <= '9')
		if reserved || strings.ContainsAny(c.Name, `<>"|?*`) || strings.HasSuffix(c.Name, ".") || strings.HasSuffix(c.Name, " ") {
			return "", fmt.Errorf("name is not a valid Windows filename: %q", c.Name)
		}
	}
	artifacts := filepath.Join(stage, "artifacts")
	app := filepath.Join(artifacts, c.Name)
	if err := os.MkdirAll(app, 0o755); err != nil {
		return "", err
	}
	exe := c.executable(app, t)
	if err := compileGo(ctx, c, exe, release, t); err != nil {
		return "", err
	}
	if t.OS == "windows" {
		if err := windowsResources(c, exe); err != nil {
			return "", err
		}
	} else if c.Linux.Icon != "" {
		data, err := os.ReadFile(c.path(c.Linux.Icon))
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(app, "icon"+strings.ToLower(filepath.Ext(c.Linux.Icon))), data, 0o644); err != nil {
			return "", err
		}
	}
	if err := bundleNatives(ctx, c, app, t); err != nil {
		return "", err
	}
	if err := copyResources(c.path("resources"), app, c.iconSources()...); err != nil {
		return "", err
	}
	if release {
		name := c.archiveName(t) + c.archiveSuffixes(t)[0]
		if err := portableArchive(ctx, app, filepath.Join(artifacts, name), t.OS == "windows"); err != nil {
			return "", err
		}
		if err := publishArtifacts(artifacts, out, c.releaseNames(t)); err != nil {
			return "", err
		}
		return filepath.Join(out, name), nil
	}
	final := filepath.Join(out, c.Name)
	if err := publish(app, final); err != nil {
		return "", err
	}
	return final, nil
}
