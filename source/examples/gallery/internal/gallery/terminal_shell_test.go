package gallery

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/shell"
)

func terminalPage() galleryPage {
	for i, entry := range mygoCatalog {
		if entry.name == "Terminal" {
			return mygoPageStart + galleryPage(i)
		}
	}
	panic("Terminal missing from catalog")
}

func TestTerminalAutomaticLoadOnce(t *testing.T) {
	var calls atomic.Int32
	failure := errors.New("load failed")
	g := &Gallery{page: terminalPage(), createTerminal: func() (*terminal.Terminal, error) {
		calls.Add(1)
		return nil, failure
	}}
	defer g.Close()
	tt := ui.NewTester(g.View, 540, 850)
	if _, err := waitTerminal(t, &g.mygoStates["Terminal"].terminal); !errors.Is(err, failure) {
		t.Fatalf("load error = %v", err)
	}
	tt.Frame()
	g.page = shellPage
	tt.Frame()
	g.page = terminalPage()
	tt.Frame()
	if calls.Load() != 1 {
		t.Fatalf("loaded %d times", calls.Load())
	}
	if _, ok := tt.Find("Terminal unavailable"); !ok {
		t.Fatal("missing load failure")
	}
}

func TestNativeTerminal(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires opt-in native library loading")
	}
	g := &Gallery{page: terminalPage(), createTerminal: newShellTerminal}
	defer g.Close()
	var header float32
	tt := ui.NewTester(func(c *ui.Context) {
		header = shell.HeaderHeight(c)
		g.View(c)
	}, 600, 400)
	term, err := waitTerminal(t, &g.mygoStates["Terminal"].terminal)
	if err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	view, ok := tt.Find("Terminal")
	if !ok {
		t.Fatal("missing terminal view")
	}
	if r := view; r.X != 0 || r.Y != header || r.W != 600 || r.H != 400-header {
		t.Fatalf("terminal bounds = %+v, header = %v", r, header)
	}
	// Split the marker so command echo alone cannot satisfy the assertion.
	command := "echo onevra-shell-$((17+25))"
	if runtime.GOOS == "windows" {
		command = `Write-Output ("onevra-shell-" + (17+25))`
	}
	tt.Type(command)
	tt.Key(0, ui.KeyEnter)
	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(term.Text(), "onevra-shell-42") && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
		tt.Frame()
	}
	if !strings.Contains(term.Text(), "onevra-shell-42") {
		t.Fatalf("shell did not execute input: %s", term.Text())
	}
	g.Close()
	select {
	case <-term.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("shell did not close")
	}
}
