// Command packer builds and runs desktop Go applications using packer.yaml.
// It reads plugin manifests as data; no application framework is imported.
package main

import (
	"context"
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

	"go.yaml.in/yaml/v3"
)

type config struct {
	Name       string          `yaml:"name"`
	Identifier string          `yaml:"identifier"`
	Version    string          `yaml:"version"`
	Main       string          `yaml:"main"`
	Out        string          `yaml:"out"`
	Platforms  []string        `yaml:"platforms"` // Release targets; omitted means host OS/architecture.
	Generate   packagePatterns `yaml:"generate"`  // Optional package patterns passed to go generate once on the host.
	MacOS      struct {
		DisableDMG           bool   `yaml:"disableDMG"` // Skip release disk images; ZIP output is always retained.
		Icon                 string `yaml:"icon"`       // ICNS file or an Icon Composer .icon directory.
		MinimumSystemVersion string `yaml:"minimumSystemVersion"`
	} `yaml:"macos"`
	Windows struct {
		Icon string `yaml:"icon"` // ICO embedded into the executable.
	} `yaml:"windows"`
	Linux struct {
		Icon string `yaml:"icon"` // PNG or SVG included with the executable.
	} `yaml:"linux"`
	root string
}

// packagePatterns rejects YAML's implicit number-to-string coercion for Go packages.
type packagePatterns []string

func (p *packagePatterns) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("generate must be a sequence of package patterns")
	}
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			return fmt.Errorf("generate entries must be strings")
		}
	}
	var values []string
	if err := node.Decode(&values); err != nil {
		return err
	}
	*p = values
	return nil
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
		return fmt.Errorf("usage: go run ./tools/packer {run|build}; all settings belong in packer.yaml")
	}
	mode := args[0]
	c, err := loadConfig("packer.yaml")
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
	if err := c.generate(ctx); err != nil {
		return err
	}
	if err := c.resolveIcons(); err != nil {
		return err
	}
	out := c.path(c.Out)
	if mode == "build" {
		out = c.releaseDir()
		for _, t := range targets {
			if _, err := build(ctx, c, out, true, t); err != nil {
				return fmt.Errorf("%s: %w", t.key(), err)
			}
			for _, suffix := range c.archiveSuffixes(t) {
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
	d := yaml.NewDecoder(f)
	d.KnownFields(true) // Never silently ignore signing or resource settings.
	if err := d.Decode(c); err != nil {
		return nil, fmt.Errorf("%s (supported fields are defined in tools/packer): %w", path, err)
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("%s: expected one YAML document", path)
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
	for _, pkg := range c.Generate {
		if strings.TrimSpace(pkg) != pkg || pkg == "" || strings.HasPrefix(pkg, "-") || strings.ContainsAny(pkg, "\x00\r\n") {
			return nil, fmt.Errorf("generate entries must be non-empty package patterns, not flags: %q", pkg)
		}
	}
	if err := c.validateIcons(); err != nil {
		return nil, err
	}
	return c, nil
}

// Applications own go:generate directives; Packer only schedules their execution
// before dependency and icon discovery, using the host rather than release target.
func (c *config) generate(ctx context.Context) error {
	if len(c.Generate) == 0 {
		return nil
	}
	args := append([]string{"generate"}, c.Generate...)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = c.root
	cmd.Env = append(os.Environ(), "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); ctx.Err() != nil {
		return fmt.Errorf("go generate: %w", ctx.Err())
	} else if err != nil {
		return fmt.Errorf("go generate: %w", err)
	}
	return nil
}

// Default icon sources live beside packer.yaml, never in resources.
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
