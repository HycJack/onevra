package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func generateFixture(t *testing.T) *config {
	t.Helper()
	root := t.TempDir()
	c := &config{Name: "Fixture", Identifier: "me.test.fixture", Version: "1.0", Main: ".", Out: "dist", Platforms: []string{"linux/amd64", "linux/arm64"}, Generate: []string{"./assets"}, root: root}
	c.MacOS.MinimumSystemVersion = "13.0"
	data, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	put(t, filepath.Join(root, "packer.yaml"), string(data))
	put(t, filepath.Join(root, "go.mod"), "module example.com/generate\n\ngo 1.27.1\n")
	put(t, filepath.Join(root, "assets", "generate.go"), "package assets\n//go:generate go run ../generator\n")
	put(t, filepath.Join(root, "generator", "main.go"), `package main
import ("os"; "runtime")
func main() {
 if os.Getenv("GOOS") != runtime.GOOS || os.Getenv("GOARCH") != runtime.GOARCH { panic("wrong host") }
 if _, err := os.Stat("generate.go"); err != nil { panic("wrong working directory") }
 f, err := os.OpenFile("../calls", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
 if err != nil { panic(err) }
 if _, err := f.WriteString("called\n"); err != nil { panic(err) }
 if err := f.Close(); err != nil { panic(err) }
 if os.Getenv("PACKER_GENERATE_FAIL") == "1" { os.Exit(1) }
 if err := os.WriteFile("../main.go", []byte("package main\nfunc main() {}\n"), 0600); err != nil { panic(err) }
 if err := os.MkdirAll("../resources", 0755); err != nil { panic(err) }
 if err := os.WriteFile("../resources/icon.svg", []byte("<svg xmlns=\"http://www.w3.org/2000/svg\"/>"), 0600); err != nil { panic(err) }
}
`)
	// Unselected packages must not run, even though they belong to the same module.
	put(t, filepath.Join(root, "maintenance", "generate.go"), "package maintenance\n//go:generate nonexistent-maintenance-command\n")
	t.Chdir(root)
	return c
}

func TestGenerateBeforeMultiTargetBuild(t *testing.T) {
	c := generateFixture(t)
	if err := run(context.Background(), []string{"build"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(c.root, "calls"))
	if err != nil || string(data) != "called\n" {
		t.Fatalf("generate calls = %q, %v", data, err)
	}
	for _, arch := range []string{"amd64", "arm64"} {
		path := filepath.Join(c.releaseDir(), c.archiveName(target{"linux", arch})+".tar.gz")
		entries, _ := readPortableArchive(t, path, false)
		if len(entries[c.Name+"/icon.svg"]) == 0 {
			t.Fatal("generated icon was not included in archive")
		}
	}
}

func TestGenerateFailureStopsRunAndBuild(t *testing.T) {
	for _, mode := range []string{"run", "build"} {
		t.Run(mode, func(t *testing.T) {
			c := generateFixture(t)
			t.Setenv("PACKER_GENERATE_FAIL", "1")
			put(t, filepath.Join(c.root, "dist", "previous"), "keep")
			err := run(context.Background(), []string{mode})
			if err == nil || !strings.Contains(err.Error(), "go generate") {
				t.Fatalf("expected generate failure, got %v", err)
			}
			entries, err := os.ReadDir(filepath.Join(c.root, "dist"))
			if err != nil || len(entries) != 1 || entries[0].Name() != "previous" {
				t.Fatalf("outputs changed: %v, %v", entries, err)
			}
		})
	}
}

func TestGenerateExecution(t *testing.T) {
	c := generateFixture(t)
	t.Setenv("GOOS", "not-the-host")
	t.Setenv("GOARCH", "not-the-host")
	t.Chdir(t.TempDir())
	if err := c.generate(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.generate(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation = %v", err)
	}
	c.Generate = []string{"./missing-package"}
	if err := c.generate(context.Background()); err == nil {
		t.Fatal("missing package succeeded")
	}
	c.Generate = nil
	if err := c.generate(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateConfig(t *testing.T) {
	file := filepath.Join(t.TempDir(), "packer.yaml")
	for _, value := range []string{`"./assets"`, `[""]`, `["  "]`, `[" ./assets"]`, `["-run=anything"]`, `["./assets", "-n"]`, `["\u0000"]`, `["./assets\n"]`, `[123]`} {
		put(t, file, `{"name":"App","identifier":"me.app","generate":`+value+`}`)
		if _, err := loadConfig(file); err == nil {
			t.Errorf("accepted generate %s", value)
		}
	}
	for _, value := range []string{`[]`, `["./assets"]`, `["./assets/...", "./other"]`} {
		put(t, file, `{"name":"App","identifier":"me.app","generate":`+value+`}`)
		if _, err := loadConfig(file); err != nil {
			t.Errorf("rejected generate %s: %v", value, err)
		}
	}
	put(t, file, `name: "App"
identifier: "me.app"
preprocess: ["go","generate"]`)
	if _, err := loadConfig(file); err == nil {
		t.Fatal("accepted obsolete preprocess configuration")
	}
}
