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
	"github.com/oneo-me/onevra/source/plugins/chart"
	"github.com/oneo-me/onevra/source/shells/sidebar"
	"github.com/oneo-me/onevra/source/theme"
)

func TestNativeCharts(t *testing.T) {
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
		theme          theme.Mode
		width          int
	}{
		{"light", "en", theme.Light, 540},
		{"dark", "en", theme.Dark, 540},
		{"narrow", "zh-Hans", theme.Light, 320},
	} {
		for _, example := range controlCatalog {
			if example.group != "Charts" {
				continue
			}
			t.Run(mode.name+"/"+example.name, func(t *testing.T) {
				var shell *sidebar.Shell
				var state *controlState
				mygo.RunOnMain(func() {
					g := &Gallery{page: pageFromKey("controls/" + example.name), language: newTranslator(mode.language, "en")}
					state = g.controlState(example.name, "Default")
					state.chart = &chartDemoState{selection: chart.State{Index: -1}}
					shell = sidebar.NewShell(sidebar.ShellOptions{
						Hidden: nativeBackground(), Title: example.name, Width: mode.width, Height: 640,
						MinWidth: 280, MinHeight: 400, LeftCollapsed: true, Appearance: mode.theme,
						Content: g.View,
					})
				})
				win := shell.Window()
				defer win.Destroy()
				capture := func(suffix string) {
					t.Helper()
					data, err := win.CapturePage()
					if err != nil {
						t.Fatal(err)
					}
					img, err := png.Decode(bytes.NewReader(data))
					if err != nil || img.Bounds().Dx() < mode.width {
						t.Fatalf("capture: %v", err)
					}
					path := filepath.Join(dir, fmt.Sprintf("chart-%s-%s-%s.png", mode.name, example.name, suffix))
					if err := os.WriteFile(path, data, 0600); err != nil {
						t.Fatal(err)
					}
				}
				for _, index := range []int{-1, 1} {
					mygo.RunOnMain(func() { state.chart.selection.Index = index; win.Invalidate() })
					time.Sleep(180 * time.Millisecond)
					capture(fmt.Sprint(index))
				}
				mygo.RunOnMain(func() { state.on = true; win.Invalidate() })
				time.Sleep(90 * time.Millisecond)
				capture("updating")
				time.Sleep(300 * time.Millisecond)
				capture("updated")
			})
		}
	}
}
