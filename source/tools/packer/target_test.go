package main

import (
	"context"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestTargets(t *testing.T) {
	got, err := parseTargets([]string{"darwin/arm64", " windows/amd64 ", "darwin/arm64", "linux/arm64"})
	want := []target{{"darwin", "arm64"}, {"windows", "amd64"}, {"linux", "arm64"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("targets = %v, %v", got, err)
	}
	for _, bad := range []string{"", "darwin", "darwin/universal", "windows/386", "linux/arm", "plan9/amd64", "darwin/amd64/extra"} {
		if _, err := parseTargets([]string{bad}); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
	c := &config{Platforms: []string{"windows/amd64", "linux/arm64"}}
	for _, tc := range []struct {
		mode string
		want []target
	}{
		{"build", []target{{"windows", "amd64"}, {"linux", "arm64"}}},
		{"run", []target{hostTarget()}},
	} {
		got, err := c.targets(tc.mode)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("selection %+v: %v, %v", tc, got, err)
		}
	}
	if err := checkHost(want, "linux"); err == nil {
		t.Fatal("accepted macOS packaging on Linux")
	}
	if err := checkHost(want[1:], "windows"); err != nil {
		t.Fatal(err)
	}
}

func TestCommandsRejectOptions(t *testing.T) {
	for _, args := range [][]string{
		nil, {"unknown"}, {"build", "-platform", "darwin/arm64"},
		{"run", "-icon", "icon.png"}, {"build", "-config", "other.json"},
		{"run", "--", "-theme", "dark"},
	} {
		if err := run(context.Background(), args); err == nil || !strings.Contains(err.Error(), "all settings belong in packer.yaml") {
			t.Fatalf("unexpected command handling for %v: %v", args, err)
		}
	}
}

func TestTargetEnvironmentAndFlags(t *testing.T) {
	t.Setenv("GOOS", "plan9")
	t.Setenv("GOARCH", "386")
	t.Setenv("CGO_ENABLED", "1")
	target := target{"windows", "arm64"}
	cmd := exec.Command("go")
	cmd.Env = buildEnv(target)
	env := strings.Join(cmd.Environ(), "\n")
	for _, want := range []string{"GOOS=windows", "GOARCH=arm64", "CGO_ENABLED=0"} {
		if !strings.Contains(env, want) {
			t.Fatalf("missing %s in child environment", want)
		}
	}
	for _, wrong := range []string{"GOOS=plan9", "GOARCH=386", "CGO_ENABLED=1"} {
		if strings.Contains(env, wrong) {
			t.Fatalf("inherited %s overrides selected target", wrong)
		}
	}
	c := &config{Main: "."}
	if !strings.Contains(strings.Join(c.buildArgs("app.exe", true, target), " "), "-H=windowsgui") {
		t.Fatal("Windows release should use the GUI subsystem")
	}
	want := []string{"build", "-o", "app.exe", "-ldflags", "-H=windowsgui", "."}
	if got := c.buildArgs("app.exe", false, target); !reflect.DeepEqual(got, want) {
		t.Fatalf("Windows development should use the GUI subsystem without stripping debug information: %v", got)
	}
}

func TestArchiveFormatSelection(t *testing.T) {
	for _, tc := range []struct {
		name      string
		platforms []string
		target    target
		want      []string
	}{
		{"defaults", nil, target{"darwin", "arm64"}, []string{".dmg", ".zip"}},
		{"dmg", []string{"darwin/arm64/dmg"}, target{"darwin", "arm64"}, []string{".dmg"}},
		{"zip", []string{"darwin/arm64/zip"}, target{"darwin", "arm64"}, []string{".zip"}},
		{"merge", []string{"darwin/arm64/dmg", "darwin/arm64/zip", "darwin/arm64/dmg"}, target{"darwin", "arm64"}, []string{".dmg", ".zip"}},
		{"mixed", []string{"darwin/arm64/zip", "darwin/arm64"}, target{"darwin", "arm64"}, []string{".zip", ".dmg"}},
		{"separate-architectures", []string{"darwin/arm64/dmg", "darwin/amd64/zip"}, target{"darwin", "amd64"}, []string{".zip"}},
		{"windows", []string{"windows/amd64/zip"}, target{"windows", "amd64"}, []string{".zip"}},
		{"linux", []string{"linux/arm64/tar.gz"}, target{"linux", "arm64"}, []string{".tar.gz"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := &config{Platforms: tc.platforms}
			if _, err := c.targets("build"); err != nil {
				t.Fatal(err)
			}
			if got := c.archiveSuffixes(tc.target); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("formats = %v, want %v", got, tc.want)
			}
		})
	}
	got, err := parseTargets([]string{"darwin/arm64/dmg", "darwin/arm64/zip", "windows/amd64"})
	want := []target{{"darwin", "arm64"}, {"windows", "amd64"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("merged targets = %v, %v", got, err)
	}
	for _, bad := range []string{"darwin/arm64/", "darwin/arm64/dmg/zip", "windows/amd64/dmg", "linux/arm64/zip", "darwin/arm64/tar.gz"} {
		if _, err := parseTargets([]string{bad}); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}
