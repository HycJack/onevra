package gallery

import (
	"sync"
	"time"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

// terminalDemo loads the optional native library once when its page is entered.
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

func newShellTerminal() (*terminal.Terminal, error) {
	return terminal.New(terminal.Options{})
}

func exampleTerminal(c *ui.Context, s *demoState) {
	term, loading, err := s.terminal.snapshot()
	if loading {
		components.Spinner(c, components.SpinnerSpokes).Label(tr(c, "gallery.loading_terminal"))
		c.After(50 * time.Millisecond)
		return
	}
	if err != nil {
		components.Text(c, tr(c, "gallery.terminal_unavailable")).Bold()
		components.Text(c, err.Error()).TextColor(c.Theme().Danger).Wrap()
		return
	}
	if term != nil {
		terminal.View(c, term).Key("terminal-session").Fill().AutoFocus()
	}
}

// Close releases gallery-owned native resources. Call on the UI thread when
// its window closes, or after a headless view is no longer used.
func (g *Gallery) Close() {
	g.notification.close()
	g.dock.close()
	g.tray.close()
	for _, s := range g.mygoStates {
		s.terminal.close()
	}
}
