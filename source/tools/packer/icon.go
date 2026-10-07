package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func compileIcon(ctx context.Context, c *config, work, resources string) (map[string]any, error) {
	info := map[string]any{}
	if c.MacOS.Icon == "" {
		return info, nil
	}
	path := c.path(c.MacOS.Icon)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".icns":
		if err := copyICNS(path, filepath.Join(resources, "AppIcon.icns")); err != nil {
			return nil, err
		}
		info["CFBundleIconFile"] = "AppIcon.icns"
	case ".icon":
		if stat, err := os.Stat(path); err != nil {
			return nil, err
		} else if !stat.IsDir() {
			return nil, fmt.Errorf("Icon Composer input must be a .icon directory: %s", path)
		}
		partial := filepath.Join(work, "icon.plist")
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if _, err := command(ctx, "xcrun", "actool", "--compile", resources,
			"--platform", "macosx", "--minimum-deployment-target", c.MacOS.MinimumSystemVersion,
			"--app-icon", name, "--output-partial-info-plist", partial, path); err != nil {
			return nil, fmt.Errorf("compiling .icon requires Xcode with Icon Composer support: %w", err)
		}
		data, err := command(ctx, "plutil", "-convert", "json", "-o", "-", partial)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &info); err != nil {
			return nil, fmt.Errorf("actool metadata: %w", err)
		}
		if info["CFBundleIconName"] != name {
			return nil, fmt.Errorf("actool did not produce CFBundleIconName for %s", name)
		}
		if _, err := os.Stat(filepath.Join(resources, "Assets.car")); err != nil {
			return nil, fmt.Errorf("actool did not produce Assets.car: %w", err)
		}
	default:
		return nil, fmt.Errorf("macos.icon must be an ICNS file or an Icon Composer .icon directory: %s", path)
	}
	return info, nil
}

// Check the container framing, then copy the prebuilt icon without re-encoding it.
func copyICNS(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if len(data) < 16 || string(data[:4]) != "icns" || uint64(binary.BigEndian.Uint32(data[4:8])) != uint64(len(data)) {
		return fmt.Errorf("invalid ICNS container: %s", src)
	}
	for offset := 8; offset < len(data); {
		if len(data)-offset < 8 {
			return fmt.Errorf("truncated ICNS element: %s", src)
		}
		size := uint64(binary.BigEndian.Uint32(data[offset+4 : offset+8]))
		if size < 8 || size > uint64(len(data)-offset) {
			return fmt.Errorf("invalid ICNS element size: %s", src)
		}
		offset += int(size)
	}
	return os.WriteFile(dst, data, 0o644)
}
