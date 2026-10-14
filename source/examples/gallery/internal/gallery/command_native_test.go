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
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

func TestNativeCommand(t *testing.T) {
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
	for _, dark := range []bool{false, true} {
		for _, width := range []int{820, 360} {
			g := &Gallery{}
			g.openCommands()
			var shell *sidebar.Shell
			mygo.RunOnMain(func() {
				if width == 360 {
					g.chooseLanguage("zh-Hans")
				}
				shell = sidebar.NewShell(sidebar.ShellOptions{Hidden: nativeBackground(), Title: "Command palette review", Width: width, Height: 600,
					MinWidth: 280, MinHeight: 400, LeftCollapsed: width == 360, Theme: pinnedTheme(dark), Left: g.Left, Content: g.View})
			})
			win := shell.Window()
			time.Sleep(300 * time.Millisecond)
			data, err := win.CapturePage()
			if width == 820 {
				mygo.RunOnMain(func() { g.commands.open = false; win.Invalidate() })
				time.Sleep(200 * time.Millisecond)
				preview, captureErr := win.CapturePage()
				if captureErr != nil {
					win.Destroy()
					t.Fatal(captureErr)
				}
				if writeErr := os.WriteFile(filepath.Join(dir, fmt.Sprintf("command-trigger-%v.png", dark)), preview, 0600); writeErr != nil {
					win.Destroy()
					t.Fatal(writeErr)
				}
			}
			mygo.RunOnMain(func() {
				g.commands.open = false
				g.settings.open, g.settings.view.Category = true, "general"
				win.Invalidate()
			})
			time.Sleep(200 * time.Millisecond)
			settingsPreview, settingsErr := win.CapturePage()
			if settingsErr != nil {
				win.Destroy()
				t.Fatal(settingsErr)
			}
			if writeErr := os.WriteFile(filepath.Join(dir, fmt.Sprintf("command-settings-%v-%d.png", dark, width)), settingsPreview, 0600); writeErr != nil {
				win.Destroy()
				t.Fatal(writeErr)
			}
			win.Destroy()
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() < width {
				t.Fatalf("capture: %v", err)
			}
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("command-%v-%d.png", dark, width)), data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Logf("Command captures: %s", dir)
}

func TestNativeCommandExpansion(t *testing.T) {
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
	for _, dark := range []bool{false, true} {
		state, open := command.State{}, true
		var shell *sidebar.Shell
		mygo.RunOnMain(func() {
			shell = sidebar.NewShell(sidebar.ShellOptions{Hidden: nativeBackground(), Title: "Command expansion review", Width: 820, Height: 600,
				LeftCollapsed: true, Theme: pinnedTheme(dark), Content: func(c *ui.Context) {
					command.CommandDialog(c, &open, &state, []command.Item{{ID: "settings", Title: "Settings"}},
						command.Options{ExpandOnSearch: true}, components.DialogOptions{})
				}})
		})
		win := shell.Window()
		for _, phase := range []string{"collapsed", "expanded", "cleared"} {
			mygo.RunOnMain(func() {
				state.Query = ""
				if phase == "expanded" {
					state.Query = "set"
				}
				win.Invalidate()
			})
			time.Sleep(200 * time.Millisecond)
			data, err := win.CapturePage()
			if err != nil {
				win.Destroy()
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("command-%s-%v.png", phase, dark)), data, 0600); err != nil {
				win.Destroy()
				t.Fatal(err)
			}
		}
		win.Destroy()
	}
	t.Logf("Command expansion captures: %s", dir)
}
