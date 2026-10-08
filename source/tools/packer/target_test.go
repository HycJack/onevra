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
	if strings.Contains(strings.Join(c.buildArgs("app.exe", false, target), " "), "-H=windowsgui") {
		t.Fatal("development should keep console output")
	}
}
