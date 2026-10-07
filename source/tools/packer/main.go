// Command packer builds and runs desktop Go applications using packer.json.
// It reads plugin manifests as data; no application framework is imported.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type config struct {
	Name       string   `json:"name"`
	Identifier string   `json:"identifier"`
	Version    string   `json:"version"`
	Main       string   `json:"main"`
	Out        string   `json:"out"`
	Platforms  []string `json:"platforms"` // Release targets; omitted means host OS/architecture.
	MacOS      struct {
		Icon                 string `json:"icon"` // ICNS file or an Icon Composer .icon directory.
		MinimumSystemVersion string `json:"minimumSystemVersion"`
	} `json:"macos"`
	Windows struct {
		Icon string `json:"icon"` // ICO embedded into the executable.
	} `json:"windows"`
	Linux struct {
		Icon string `json:"icon"` // PNG or SVG included with the executable.
	} `json:"linux"`
	root string
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "packer:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) != 1 || (args[0] != "run" && args[0] != "build") {
		return fmt.Errorf("usage: go run ./tools/packer {run|build}; all settings belong in packer.json")
	}
	mode := args[0]
	c, err := loadConfig("packer.json")
	if err != nil {
		return err
	}
	targets, err := c.targets(mode)
	if err != nil {
		return err
	}
	if err := checkHost(targets, runtime.GOOS); err != nil {
		return err
	}
	out := c.path(c.Out)
	if mode == "build" {
		out = c.releaseDir()
		for _, t := range targets {
			if _, err := build(ctx, c, out, true, t); err != nil {
				return fmt.Errorf("%s: %w", t.key(), err)
			}
			for _, suffix := range t.archiveSuffixes() {
				fmt.Fprintln(os.Stderr, filepath.Join(out, c.archiveName(t)+suffix))
			}
		}
		return nil
	}
	dev := *c
	dev.Name += " Dev"
	dev.Identifier += ".dev"
	c = &dev
	app, err := build(ctx, c, out, false, targets[0])
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, app)
	cmd := exec.CommandContext(ctx, c.executable(app, targets[0]))
	cmd.Dir = c.root
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); ctx.Err() != nil {
		return ctx.Err()
	} else {
		return err
	}
}

func loadConfig(path string) (*config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := &config{Version: "0.0.0", Main: ".", Out: "dist", root: filepath.Dir(path)}
	c.MacOS.MinimumSystemVersion = "13.0"
	d := json.NewDecoder(f)
	d.DisallowUnknownFields() // Never silently ignore signing or resource settings.
	if err := d.Decode(c); err != nil {
		return nil, fmt.Errorf("%s (supported fields are defined in tools/packer): %w", path, err)
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("%s: expected one JSON object", path)
	}
	if !component(c.Name) {
		return nil, fmt.Errorf("name must be a non-empty filename, not a path")
	}
	for name, value := range map[string]string{"name": c.Name, "identifier": c.Identifier, "version": c.Version} {
		if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\x00\r\n") {
			return nil, fmt.Errorf("invalid %s: must be non-empty and cannot contain newlines or NUL", name)
		}
	}
	if !component(c.Version) {
		return nil, fmt.Errorf("version must be a non-empty filename component, not a path")
	}
	if c.Main == "" || strings.HasPrefix(c.Main, "-") || c.Out == "" || c.MacOS.MinimumSystemVersion == "" {
		return nil, fmt.Errorf("main, out and macos.minimumSystemVersion must be non-empty; main cannot be a flag")
	}
	if _, err := parseTargets(c.Platforms); err != nil {
		return nil, err
	}
	if err := c.resolveIcons(); err != nil {
		return nil, err
	}
	if err := c.validateIcons(); err != nil {
		return nil, err
	}
	return c, nil
}

// Default icon sources live beside packer.json, never in resources.
func (c *config) resolveIcons() error {
	for _, platform := range []struct {
		icon  *string
		names []string
	}{
		{&c.MacOS.Icon, []string{"icon.icon", "icon.icns"}},
		{&c.Windows.Icon, []string{"icon.ico"}},
		{&c.Linux.Icon, []string{"icon.png", "icon.svg"}},
	} {
		if *platform.icon != "" {
			continue
		}
		for _, name := range platform.names {
			info, err := os.Stat(c.path(name))
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return err
			}
			if (filepath.Ext(name) == ".icon" && !info.IsDir()) ||
				(filepath.Ext(name) != ".icon" && !info.Mode().IsRegular()) {
				return fmt.Errorf("invalid icon source %s", c.path(name))
			}
			*platform.icon = name
			break
		}
	}
	return nil
}

// Check declared formats without requiring other platforms' assets on this host.
// The selected platform's compiler validates the actual icon contents.
func (c *config) validateIcons() error {
	for _, icon := range []struct{ field, path, formats string }{
		{"macos.icon", c.MacOS.Icon, ".icon .icns"},
		{"windows.icon", c.Windows.Icon, ".ico"},
		{"linux.icon", c.Linux.Icon, ".png .svg"},
	} {
		if icon.path == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(icon.path))
		valid := false
		for _, format := range strings.Fields(icon.formats) {
			valid = valid || ext == format
		}
		if !valid {
			return fmt.Errorf("%s must use one of %s, got %q", icon.field, icon.formats, icon.path)
		}
	}
	return nil
}

func component(s string) bool {
	return strings.TrimSpace(s) != "" && s != "." && s != ".." && !strings.ContainsAny(s, "/\\:\x00\r\n")
}

func (c *config) path(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(c.root, p)
}

func (c *config) releaseDir() string {
	return filepath.Join(c.path(c.Out), "releases", c.Version)
}

func (c *config) buildArgs(exe string, release bool, t target) []string {
	args := []string{"build", "-o", exe}
	if release {
		ldflags := "-s -w"
		if t.OS == "windows" {
			ldflags += " -H=windowsgui"
		}
		args = append(args, "-trimpath", "-ldflags", ldflags)
	}
	return append(args, c.Main)
}

func command(ctx context.Context, name string, args ...string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %w\n%s", name, err, out)
	}
	return out, nil
}
