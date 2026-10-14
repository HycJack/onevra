package gallery

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

func TestNativeMyGoFeatures(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []struct {
		name, language string
		dark           bool
		width          int
	}{
		{"light", "en", false, 540},
		{"dark", "en", true, 540},
		{"narrow", "zh-Hans", false, 280},
	} {
		for _, name := range []string{"Transforms", "Motion"} {
			for _, variant := range mygoVariants(name) {
				t.Run(mode.name+"/"+name+"/"+variant, func(t *testing.T) {
					var shell *sidebar.Shell
					mygo.RunOnMain(func() {
						g := &Gallery{language: newTranslator(mode.language, "en")}
						s := g.mygoState(name, variant, true)
						s.index = 1
						s.open = name == "Transforms" && variant == "Editable controls"
						example, _ := mygoExampleFor(pageFromKey("mygo/" + name))
						shell = sidebar.NewShell(sidebar.ShellOptions{
							Hidden: nativeBackground(), Title: name, Width: mode.width, Height: 480,
							MinWidth: 280, MinHeight: 400, LeftCollapsed: true, Theme: pinnedTheme(mode.dark),
							Content: func(c *ui.Context) {
								g.withLanguage(c, func() {
									ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c)+16, 16, 16, 16).Gap(12).
										Children(func() { example.build(c, s) })
								})
							},
						})
					})
					win := shell.Window()
					defer win.Destroy()
					time.Sleep(250 * time.Millisecond)
					data, err := win.CapturePage()
					if err != nil {
						t.Fatal(err)
					}
					img, err := png.Decode(bytes.NewReader(data))
					if err != nil || img.Bounds().Dx() < mode.width {
						t.Fatalf("capture: %v", err)
					}
					path := filepath.Join(dir, fmt.Sprintf("features-%s-%s-%s.png", mode.name, name, variant))
					if err := os.WriteFile(path, data, 0600); err != nil {
						t.Fatal(err)
					}
				})
			}
		}
	}
}
