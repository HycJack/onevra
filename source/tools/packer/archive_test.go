package main

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestDisableDMG(t *testing.T) {
	for _, tc := range []struct {
		name, macos string
		disabled    bool
	}{
		{"default", `{}`, false},
		{"enabled", `{"disableDMG":false}`, false},
		{"disabled", `{"disableDMG":true}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "packer.yaml")
			put(t, file, `{"name":"App","identifier":"me.app","macos":`+tc.macos+`}`)
			c, err := loadConfig(file)
			if err != nil {
				t.Fatal(err)
			}
			if c.MacOS.DisableDMG != tc.disabled {
				t.Fatalf("unexpected setting: %+v", c.MacOS)
			}
			for _, goos := range []string{"darwin", "windows", "linux"} {
				want := []string{".zip"}
				if goos == "darwin" && !tc.disabled {
					want = []string{".dmg", ".zip"}
				}
				if goos == "linux" {
					want = []string{".tar.gz"}
				}
				if got := c.archiveSuffixes(target{goos, "arm64"}); !reflect.DeepEqual(got, want) {
					t.Fatalf("%s suffixes = %v, want %v", goos, got, want)
				}
			}
		})
	}
}

// Exercises the real ZIP tool without requiring Finder or creating a disk image.
func TestArchiveWithoutDMG(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("requires macOS ditto")
	}
	stage := t.TempDir()
	app := filepath.Join(stage, "artifacts", "App.app")
	put(t, filepath.Join(app, "Contents", "payload"), "payload")
	c := &config{Name: "App", Version: "1.0"}
	c.MacOS.DisableDMG = true
	target := hostTarget()
	if err := archive(context.Background(), c, stage, app, target); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(filepath.Dir(app), c.archiveName(target))
	z, err := zip.OpenReader(base + ".zip")
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	found := false
	for _, f := range z.File {
		if f.Name == "App.app/Contents/payload" {
			found = true
		}
	}
	if !found {
		t.Fatal("ZIP does not contain app payload")
	}
	for _, path := range []string{base + ".dmg", filepath.Join(stage, "image")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected DMG staging/output %s: %v", path, err)
		}
	}
}
