package gallery

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
)

func waitTerminal(t *testing.T, d *terminalDemo) (*terminal.Terminal, error) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		term, loading, err := d.snapshot()
		if !loading {
			return term, err
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("terminal load did not finish")
	return nil, nil
}

func TestTerminalHeadlessDoesNotStartShell(t *testing.T) {
	g := &Gallery{}
	defer g.Close()
	for i, entry := range mygoCatalog {
		if entry.name == "Terminal" {
			g.page = mygoPageStart + galleryPage(i)
		}
	}
	tt := ui.NewTester(g.View, 540, 850)
	if _, ok := tt.Find("Start terminal"); ok {
		t.Fatal("unexpected demo controls")
	}
	term, loading, err := g.mygoStates["Terminal"].terminal.snapshot()
	if term != nil || loading || err != nil {
		t.Fatal("rendering started a terminal")
	}
}

func TestTerminalLoadLifecycle(t *testing.T) {
	var d terminalDemo
	defer d.close()
	started, release := make(chan struct{}), make(chan struct{})
	failure := errors.New("native library unavailable")
	d.start(func() (*terminal.Terminal, error) { close(started); <-release; return nil, failure })
	<-started
	d.start(func() (*terminal.Terminal, error) { t.Error("duplicate load"); return nil, nil })
	d.close()
	close(release)
	if term, err := waitTerminal(t, &d); term != nil || err != nil {
		t.Fatal("stopped load published a result")
	}
	d.start(func() (*terminal.Terminal, error) { return nil, failure })
	if _, err := waitTerminal(t, &d); !errors.Is(err, failure) {
		t.Fatalf("error = %v", err)
	}
	s := newDemoState()
	s.terminal.err = failure
	tt := ui.NewTester(func(c *ui.Context) { exampleTerminal(c, s) }, 540, 600)
	if _, ok := tt.Find("Terminal unavailable"); !ok {
		t.Fatal("missing error state")
	}
}

func TestTerminalEchoConnection(t *testing.T) {
	conn := newEchoConnection()
	defer conn.Close()
	done := make(chan error, 1)
	go func() { _, err := conn.Write([]byte("hello\r\x7f")); done <- err }()
	buf := make([]byte, len("hello\r\n\b \b"))
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "hello\r\n\b \b" {
		t.Fatalf("echo = %q", buf)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Read(buf); err == nil {
		t.Fatal("closed reader did not unblock")
	}
}

func TestNativeTerminalEcho(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires opt-in native library loading")
	}
	var d terminalDemo
	defer d.close()
	d.start(newPreviewTerminal)
	term, err := waitTerminal(t, &d)
	if err != nil {
		t.Fatal(err)
	}
	tt := ui.NewTester(func(c *ui.Context) { terminal.View(c, term).Fill().AutoFocus() }, 600, 360)
	if !strings.Contains(term.Text(), "MyGo Terminal") {
		t.Fatal("missing ANSI fixture")
	}
	tt.Type("hello-terminal")
	tt.Key(0, ui.KeyEnter)
	deadline := time.Now().Add(2 * time.Second)
	for !strings.Contains(term.Text(), "hello-terminal") && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		tt.Frame()
	}
	if !strings.Contains(term.Text(), "hello-terminal") {
		t.Fatal("keyboard input was not echoed")
	}
	d.close()
	select {
	case <-term.Done():
	case <-time.After(time.Second):
		t.Fatal("terminal connection did not close")
	}
	d.start(newPreviewTerminal)
	restarted, err := waitTerminal(t, &d)
	if err != nil || restarted == nil || restarted == term {
		t.Fatalf("restart failed: %v", err)
	}
}
