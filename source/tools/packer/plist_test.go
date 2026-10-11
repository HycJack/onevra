package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const testInfoPlistConfig = `
  infoPlist:
    NSMicrophoneUsageDescription: >-
      Record only when requested: 中文 & <voice>.
    FixtureTypes:
      Text: "中文 <&> \"quoted\"\nnext line"
      Enabled: true
      Integer: 9007199254740993
      MinInteger: -9223372036854775808
      MaxInteger: 18446744073709551615
      Real: 1.0
      Fraction: -1.25
      Exponent: 1.0e+20
      NegativeZero: -0.0
      TaggedReal: !!float 2
      HexInteger: 0xff
      EmptyArray: []
      EmptyDict: {}
      Values: [false, "", 7, 2.0, {Nested: "value"}]
`

func TestInfoPlistConfig(t *testing.T) {
	file := filepath.Join(t.TempDir(), "packer.yaml")
	for _, setting := range []string{"", "macos:\n  infoPlist: {}\n"} {
		put(t, file, "name: App\nidentifier: me.app\n"+setting)
		c, err := loadConfig(file)
		if err != nil {
			t.Fatal(err)
		}
		info, err := c.macOSInfo(nil)
		if err != nil || !reflect.DeepEqual(info, c.bundleInfo()) {
			t.Fatalf("empty extension changed metadata: %v, %v", info, err)
		}
	}
	put(t, file, "name: App\nidentifier: me.app\nmacos:"+testInfoPlistConfig)
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	info, err := c.macOSInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	if info["NSMicrophoneUsageDescription"] != "Record only when requested: 中文 & <voice>." {
		t.Fatalf("usage description = %q", info["NSMicrophoneUsageDescription"])
	}
	// Compare the JSON bridge's number spelling as well as its value. A real
	// with an integral value must still arrive at plutil as a real.
	data, err := json.Marshal(info["FixtureTypes"])
	if err != nil {
		t.Fatal(err)
	}
	want := `{"EmptyArray":[],"EmptyDict":{},"Enabled":true,"Exponent":1e+20,"Fraction":-1.25,"HexInteger":255,"Integer":9007199254740993,"MaxInteger":18446744073709551615,"MinInteger":-9223372036854775808,"NegativeZero":-0.0,"Real":1.0,"TaggedReal":2.0,"Text":"中文 \u003c\u0026\u003e \"quoted\"\nnext line","Values":[false,"",7,2.0,{"Nested":"value"}]}`
	if string(data) != want {
		t.Fatalf("typed metadata:\n%s\nwant:\n%s", data, want)
	}
}

func TestInfoPlistInvalidConfig(t *testing.T) {
	for _, tc := range []struct{ name, yaml, path, reason string }{
		{"null dictionary", "null", "macos.infoPlist", "dictionary"},
		{"empty value", "", "macos.infoPlist", "dictionary"},
		{"sequence", "[]", "macos.infoPlist", "dictionary"},
		{"string", "text", "macos.infoPlist", "dictionary"},
		{"null", "{Example: null}", "macos.infoPlist.Example", "!!null"},
		{"nested null", "{Example: {Items: [1, null]}}", "macos.infoPlist.Example.Items[1]", "!!null"},
		{"numeric key", "{1: value}", "macos.infoPlist", "keys must be strings"},
		{"nested boolean key", "{Example: {true: value}}", "macos.infoPlist.Example", "keys must be strings"},
		{"duplicate", "{Example: 1, Example: 2}", "macos.infoPlist.Example", "duplicate"},
		{"nested duplicate", "{Example: {Key: 1, Key: 2}}", "macos.infoPlist.Example.Key", "duplicate"},
		{"date", "{Example: 2026-10-11}", "macos.infoPlist.Example", "!!timestamp"},
		{"binary", "{Example: !!binary YQ==}", "macos.infoPlist.Example", "!!binary"},
		{"custom scalar", "{Example: !custom text}", "macos.infoPlist.Example", "!custom"},
		{"custom dictionary", "{Example: !custom {Key: value}}", "macos.infoPlist.Example", "!custom"},
		{"custom sequence", "{Example: !custom [1]}", "macos.infoPlist.Example", "!custom"},
		{"nan", "{Example: .nan}", "macos.infoPlist.Example", "finite"},
		{"infinity", "{Example: .inf}", "macos.infoPlist.Example", "finite"},
		{"negative infinity", "{Example: -.inf}", "macos.infoPlist.Example", "finite"},
		{"large integer", "{Example: 18446744073709551616}", "macos.infoPlist.Example", "64-bit"},
		{"negative integer", "{Example: -9223372036854775809}", "macos.infoPlist.Example", "64-bit"},
		{"invalid integer", "{Example: !!int text}", "macos.infoPlist.Example", "cannot decode"},
		{"control character", `{Example: "a\0b"}`, "macos.infoPlist.Example", "XML"},
		{"invalid key character", `{"a\0b": value}`, `macos.infoPlist["a\x00b"]`, "XML"},
		{"quoted path", `{"a.b": [null]}`, `macos.infoPlist["a.b"][0]`, "!!null"},
		{"alias", "{First: &value text, Example: *value}", "macos.infoPlist.Example", "aliases"},
		{"recursive alias", "{Example: &value [*value]}", "macos.infoPlist.Example[0]", "aliases"},
		{"merge", "{Example: {<<: {Key: value}}}", "macos.infoPlist.Example", "merge"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "packer.yaml")
			put(t, file, "name: App\nidentifier: me.app\nmacos:\n  infoPlist: "+tc.yaml+"\n")
			_, err := loadConfig(file)
			if err == nil || !strings.Contains(err.Error(), tc.path) || !strings.Contains(err.Error(), tc.reason) || !strings.Contains(err.Error(), "line ") {
				t.Fatalf("error = %v; want path %q and reason %q with line number", err, tc.path, tc.reason)
			}
		})
	}
}

func TestInfoPlistLeadingZeroNumber(t *testing.T) {
	file := filepath.Join(t.TempDir(), "packer.yaml")
	// The YAML decoder resolves 09 as a float, but it is not an overflowing
	// integer. Keep that finite value and the type assigned by the decoder.
	put(t, file, "name: App\nidentifier: me.app\nmacos:\n  infoPlist: {Example: 09}\n")
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	info, err := c.macOSInfo(nil)
	if err != nil || info["Example"] != json.Number("9.0") {
		t.Fatalf("leading-zero number: %v, %v", info, err)
	}
}

func TestInfoPlistManagedKeys(t *testing.T) {
	c := &config{}
	keys := c.bundleInfo()
	for _, key := range []string{"CFBundleIconFile", "CFBundleIconFiles", "CFBundleIconName", "CFBundleIcons"} {
		keys[key] = nil
	}
	for key := range keys {
		t.Run(key, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "packer.yaml")
			put(t, file, "name: App\nidentifier: me.app\nmacos:\n  infoPlist:\n    "+key+": App\n")
			if _, err := loadConfig(file); err == nil || !strings.Contains(err.Error(), "macos.infoPlist."+key) || !strings.Contains(err.Error(), "conflicts") {
				t.Fatalf("managed field accepted: %v", err)
			}
		})
	}
}

func TestInfoPlistStrictOuterConfig(t *testing.T) {
	for _, setting := range []string{
		"infoPlist: {}", "macos:\n  infoPlsit: {}", "macos:\n  infoPlist: {}\n  signingIdentity: any",
		"windows:\n  infoPlist: {}", "linux:\n  infoPlist: {}", "unknown: true",
	} {
		file := filepath.Join(t.TempDir(), "packer.yaml")
		put(t, file, "name: App\nidentifier: me.app\n"+setting+"\n")
		if _, err := loadConfig(file); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("unknown field accepted in %q: %v", setting, err)
		}
	}
}

func TestInfoPlistMerge(t *testing.T) {
	file := filepath.Join(t.TempDir(), "packer.yaml")
	put(t, file, "name: App\nidentifier: me.app\nmacos:\n  infoPlist: {Custom: value}\n")
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	icon := map[string]any{"CFBundleIconName": "Fixture", "CompilerKey": "generated"}
	info, err := c.macOSInfo(icon)
	if err != nil || info["Custom"] != "value" || info["CompilerKey"] != "generated" || info["CFBundleIconName"] != "Fixture" {
		t.Fatalf("merged metadata: %v, %v", info, err)
	}
	if len(icon) != 2 {
		t.Fatal("merge modified the icon input")
	}
	dev := *c
	dev.Name += " Dev"
	dev.Identifier += ".dev"
	devInfo, err := dev.macOSInfo(icon)
	if err != nil || devInfo["CFBundleName"] != "App Dev" || devInfo["CFBundleIdentifier"] != "me.app.dev" || info["CFBundleName"] != "App" {
		t.Fatalf("development metadata: %v, %v", devInfo, err)
	}
	// An unknown key returned by a future icon compiler is protected too, even
	// if the application supplied the same value.
	icon["Custom"] = "value"
	if _, err := c.macOSInfo(icon); err == nil || !strings.Contains(err.Error(), "macos.infoPlist.Custom") {
		t.Fatalf("icon compiler metadata collision accepted: %v", err)
	}
}

func TestInfoPlistFailsBeforeBuild(t *testing.T) {
	root := t.TempDir()
	put(t, filepath.Join(root, "packer.yaml"), "name: App\nidentifier: me.app\nout: dist\nmacos:\n  infoPlist: {CFBundleName: App}\n")
	previous := filepath.Join(root, "dist", "releases", "0.0.0", "previous")
	put(t, previous, "retained")
	t.Chdir(root)
	// There is deliberately no Go module or icon toolchain in this fixture.
	for _, mode := range []string{"build", "run"} {
		if err := run(context.Background(), []string{mode}); err == nil || !strings.Contains(err.Error(), "macos.infoPlist.CFBundleName") {
			t.Fatalf("%s did not fail during config validation: %v", mode, err)
		}
	}
	if data, err := os.ReadFile(previous); err != nil || string(data) != "retained" {
		t.Fatalf("invalid config changed previous outputs: %v", err)
	}
}

func TestNativePlistPackaging(t *testing.T) {
	if runtime.GOOS != "darwin" || os.Getenv("ONEVRA_PACKER_SMOKE") != "1" {
		t.Skip("set ONEVRA_PACKER_SMOKE=1 on macOS")
	}
	root := t.TempDir()
	file := filepath.Join(root, "packer.yaml")
	put(t, file, "name: App\nidentifier: me.app\nmacos:"+testInfoPlistConfig)
	c, err := loadConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	info, err := c.macOSInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "Info.plist")
	ctx := context.Background()
	if err := writePlist(ctx, path, info); err != nil {
		t.Fatal(err)
	}
	if _, err := command(ctx, "plutil", "-lint", path); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{
		"Text": "string", "Enabled": "bool", "Integer": "integer", "MinInteger": "integer", "MaxInteger": "integer",
		"Real": "float", "Fraction": "float", "Exponent": "float", "NegativeZero": "float", "TaggedReal": "float",
		"HexInteger": "integer", "EmptyArray": "array", "EmptyDict": "dictionary", "Values": "array",
		"Values.0": "bool", "Values.1": "string", "Values.2": "integer", "Values.3": "float", "Values.4": "dictionary",
	} {
		data, err := command(ctx, "plutil", "-type", "FixtureTypes."+key, path)
		if err != nil || strings.TrimSpace(string(data)) != want {
			t.Errorf("%s type = %s, %v; want %s", key, data, err, want)
		}
	}
	for key, want := range map[string]string{
		"Text": "中文 <&> \"quoted\"\nnext line", "Integer": "9007199254740993",
		"MinInteger": "-9223372036854775808", "MaxInteger": "18446744073709551615", "Values.4.Nested": "value",
	} {
		data, err := command(ctx, "plutil", "-extract", "FixtureTypes."+key, "raw", "-o", "-", path)
		if err != nil || strings.TrimSuffix(string(data), "\n") != want {
			t.Errorf("%s = %q, %v; want %q", key, data, err, want)
		}
	}
}

func assertUsageDescription(t *testing.T, app string) {
	t.Helper()
	data, err := command(context.Background(), "plutil", "-extract", "NSMicrophoneUsageDescription", "raw", "-o", "-", filepath.Join(app, "Contents", "Info.plist"))
	if err != nil || strings.TrimSuffix(string(data), "\n") != "Record only when requested: 中文 & <voice>." {
		t.Fatalf("%s usage description: %q, %v", app, data, err)
	}
}
