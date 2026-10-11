package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

// macOSInfo builds a fresh dictionary for both development and release bundles.
// With no icon metadata it also performs configuration-only validation, before
// compilation or removal of previous release outputs.
func (c *config) macOSInfo(icon map[string]any) (map[string]any, error) {
	info := make(map[string]any, len(icon))
	for key, value := range icon {
		info[key] = value
	}
	for key, value := range c.bundleInfo() {
		info[key] = value
	}
	node := &c.MacOS.InfoPlist
	if node.IsZero() {
		return info, nil
	}
	const path = "macos.infoPlist"
	if node.Kind != yaml.MappingNode || node.Tag != "!!map" {
		return nil, plistError(node, path, "must be a string-keyed dictionary")
	}
	value, err := plistValue(node, path)
	if err != nil {
		return nil, err
	}
	extensions := value.(map[string]any)
	// Walk source order so errors are stable and point to the offending key.
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		_, reserved := info[key.Value]
		switch key.Value {
		case "CFBundleIconFile", "CFBundleIconFiles", "CFBundleIconName", "CFBundleIcons":
			reserved = true // Reserved even when no icon is selected.
		}
		if reserved {
			return nil, plistError(key, plistKeyPath(path, key.Value), "conflicts with Packer-managed bundle or icon metadata")
		}
		info[key.Value] = extensions[key.Value]
	}
	return info, nil
}

// plistValue keeps the supported YAML types through the JSON-to-plist bridge.
// Inspect nodes instead of decoding arbitrary values: the latter can coerce
// dictionary keys and silently turn timestamps or binary data into strings.
func plistValue(node *yaml.Node, path string) (any, error) {
	switch node.Kind {
	case yaml.MappingNode:
		if node.Tag != "!!map" {
			break
		}
		values := make(map[string]any, len(node.Content)/2)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
				return nil, plistError(key, path, "dictionary keys must be strings; YAML merge keys are not supported")
			}
			childPath := plistKeyPath(path, key.Value)
			if _, err := plistValue(key, childPath); err != nil {
				return nil, err
			}
			if _, exists := values[key.Value]; exists {
				return nil, plistError(key, childPath, "duplicate dictionary key")
			}
			value, err := plistValue(node.Content[i+1], childPath)
			if err != nil {
				return nil, err
			}
			values[key.Value] = value
		}
		return values, nil
	case yaml.SequenceNode:
		if node.Tag != "!!seq" {
			break
		}
		values := make([]any, len(node.Content))
		for i, child := range node.Content {
			value, err := plistValue(child, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			values[i] = value
		}
		return values, nil
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			for _, r := range node.Value {
				if r != '\t' && r != '\n' && r != '\r' && (r < 0x20 || r == 0xfffe || r == 0xffff) {
					return nil, plistError(node, path, "character U+%04X is not valid in XML plist text", r)
				}
			}
			return node.Value, nil
		case "!!bool", "!!int", "!!float":
			var value any
			if err := node.Decode(&value); err != nil {
				return nil, plistError(node, path, "%v", err)
			}
			if real, ok := value.(float64); ok {
				if math.IsNaN(real) || math.IsInf(real, 0) {
					return nil, plistError(node, path, "real numbers must be finite")
				}
				// YAML resolves overflowing decimal integers as floats. Reject
				// those rather than rounding them without the user's knowledge.
				if node.Style&yaml.TaggedStyle == 0 && !strings.ContainsAny(node.Value, ".eE") {
					decimal := strings.ReplaceAll(node.Value, "_", "")
					if _, err := strconv.ParseInt(decimal, 10, 64); err != nil {
						if _, err := strconv.ParseUint(strings.TrimPrefix(decimal, "+"), 10, 64); err != nil {
							return nil, plistError(node, path, "integer is outside the 64-bit signed or unsigned range")
						}
					}
				}
				text := strconv.FormatFloat(real, 'g', -1, 64)
				if !strings.ContainsAny(text, ".eE") {
					text += ".0" // json.Marshal(float64(1)) would erase the real type.
				}
				return json.Number(text), nil
			}
			return value, nil
		}
	case yaml.AliasNode:
		return nil, plistError(node, path, "YAML aliases are not supported")
	}
	return nil, plistError(node, path, "unsupported plist value %q; use strings, booleans, integers, finite reals, arrays or string-keyed dictionaries", node.Tag)
}

func plistKeyPath(path, key string) string {
	if key == "" || strings.IndexFunc(key, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_')
	}) >= 0 {
		return fmt.Sprintf("%s[%q]", path, key)
	}
	return path + "." + key
}

func plistError(node *yaml.Node, path, format string, args ...any) error {
	return fmt.Errorf("%s (line %d): %s", path, node.Line, fmt.Sprintf(format, args...))
}
