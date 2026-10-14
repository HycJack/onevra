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

func TestNativeCode(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	g := &Gallery{page: pageFromKey("controls/Code")}
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, dark := range []bool{false, true} {
		for _, width := range []int{640, 320} {
			var shell *sidebar.Shell
			mygo.RunOnMain(func() {
				shell = sidebar.NewShell(sidebar.ShellOptions{Hidden: nativeBackground(),
					Title: "Code review", Width: width, Height: 820, MinWidth: 280, MinHeight: 400,
					LeftCollapsed: true, Theme: pinnedTheme(dark),
					Content: func(c *ui.Context) {
						ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() { g.body(c) })
					},
				})
			})
			win := shell.Window()
			time.Sleep(300 * time.Millisecond)
			data, err := win.CapturePage()
			win.Destroy()
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() < width {
				t.Fatalf("capture: %v", err)
			}
			path := filepath.Join(dir, fmt.Sprintf("code-%v-%d.png", dark, width))
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Logf("Code captures: %s", dir)
}
