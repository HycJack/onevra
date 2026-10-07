package gallery

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

// terminalDemo loads the optional native library only after an explicit action.
// Loading never blocks a frame. A stopped load cannot resurrect its terminal.
type terminalDemo struct {
	mu               sync.Mutex
	term             *terminal.Terminal
	loading, stopped bool
	err              error
}

func (d *terminalDemo) start(create func() (*terminal.Terminal, error)) {
	d.mu.Lock()
	if d.loading || d.term != nil {
		d.mu.Unlock()
		return
	}
	d.loading, d.stopped, d.err = true, false, nil
	d.mu.Unlock()
	go func() {
		term, err := create()
		d.mu.Lock()
		defer d.mu.Unlock()
		d.loading = false
		if d.stopped {
			if term != nil {
				_ = term.Close()
			}
			return
		}
		d.term, d.err = term, err
	}()
}

func (d *terminalDemo) snapshot() (*terminal.Terminal, bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.term, d.loading, d.err
}

func (d *terminalDemo) close() {
	d.mu.Lock()
	d.stopped = true
	term := d.term
	d.term, d.err = nil, nil
	d.mu.Unlock()
	if term != nil {
		_ = term.Close()
	}
}

// echoConnection is a byte-stream fixture, not a shell or input system.
// MyGo owns keyboard encoding, selection, rendering and scrolling.
type echoConnection struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func newEchoConnection() *echoConnection {
	r, w := io.Pipe()
	return &echoConnection{r, w}
}
func (c *echoConnection) Read(p []byte) (int, error) { return c.reader.Read(p) }
func (c *echoConnection) Write(p []byte) (int, error) {
	data := bytes.ReplaceAll(p, []byte("\r"), []byte("\r\n"))
	data = bytes.ReplaceAll(data, []byte{127}, []byte("\b \b"))
	_, err := c.writer.Write(data)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func (c *echoConnection) Close() error {
	_ = c.reader.Close()
	return c.writer.Close()
}

func newPreviewTerminal() (*terminal.Terminal, error) {
	conn := newEchoConnection()
	term, err := terminal.New(terminal.Options{Conn: conn, Scrollback: 1 << 20})
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	term.Feed([]byte(terminalSample))
	return term, nil
}

const terminalSample = "\x1b[1mMyGo Terminal\x1b[0m\r\n" +
	"\x1b[31mRed  \x1b[32mGreen  \x1b[33mYellow  \x1b[34mBlue  \x1b[35mMagenta  \x1b[36mCyan\x1b[0m\r\n" +
	"\x1b[1mBold\x1b[0m  \x1b[3mItalic\x1b[0m  \x1b[4mUnderline\x1b[0m  \x1b[7mInverse\x1b[0m\r\n" +
	"\x1b[38;2;255;128;64mTrue color\x1b[0m  中文 日本語 🌍\r\n\r\n"

func exampleTerminal(c *ui.Context, s *demoState) {
	term, loading, err := s.terminal.snapshot()
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, "Start terminal", components.ButtonSecondary).Disabled(loading || term != nil).Clicked() {
			s.terminal.start(newPreviewTerminal)
		}
		if components.Button(c, "Stop", components.ButtonSecondary).Disabled(!loading && term == nil).Clicked() {
			s.terminal.close()
		}
	})
	// Re-read after actions so a closed terminal is never built in this frame.
	term, loading, err = s.terminal.snapshot()
	if loading {
		ui.Spinner(c).Label("Loading terminal")
		c.After(50 * time.Millisecond)
		return
	}
	if err != nil {
		ui.Text(c, "Terminal unavailable").Bold()
		ui.Text(c, err.Error()).TextColor(c.Theme().Danger).Wrap()
		return
	}
	if term == nil {
		ui.Text(c, "Stopped").TextColor(c.Theme().TextMuted)
		return
	}
	ui.Text(c, "Echo · no shell").TextColor(c.Theme().TextMuted)
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, "Colors", components.ButtonSecondary).Clicked() {
			term.Feed([]byte(terminalSample))
		}
		if components.Button(c, "Scrollback", components.ButtonSecondary).Clicked() {
			for i := range 80 {
				term.Feed(fmt.Appendf(nil, "\x1b[36m%02d\x1b[0m  Terminal output\r\n", i+1))
			}
		}
		if components.Button(c, "Clear", components.ButtonSecondary).Clicked() {
			term.Feed([]byte("\x1b[2J\x1b[3J\x1b[H"))
		}
	})
	terminal.View(c, term).Label("Terminal preview").FillWidth().Height(360)
}

// Close releases gallery-owned native resources. Call on the UI thread when
// its window closes, or after a headless view is no longer used.
func (g *Gallery) Close() {
	for _, s := range g.mygoStates {
		s.terminal.close()
	}
}
