package gallery

import (
	"sync"
	"time"

	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// terminalDemo loads the optional native library once when its page is entered.
// Loading never blocks a frame. A stopped load cannot resurrect its terminal.
type terminalDemo struct {
	mu                        sync.Mutex
	term                      *terminal.Terminal
	started, loading, stopped bool
	err                       error
}

func (d *terminalDemo) start(create func() (*terminal.Terminal, error)) {
	d.mu.Lock()
	if d.started {
		d.mu.Unlock()
		return
	}
	d.started, d.loading, d.stopped, d.err = true, true, false, nil
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
	d.started, d.stopped = false, true
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
		components.Spinner(c, components.SpinnerSpokes).Label(tr(c, "loadingTerminal"))
		c.After(50 * time.Millisecond)
		return
	}
	if err != nil {
		components.Text(c, tr(c, "terminalUnavailable")).Bold()
		components.Text(c, err.Error()).TextColor(c.Theme().Danger).Wrap()
		return
	}
	if term != nil {
		terminal.View(c, term).Key("terminal-session").Fill().AutoFocus()
	}
}

// terminalToolbar switches copy on select from the titlebar once the session
// is running.
func terminalToolbar(c *ui.Context, s *demoState) {
	term, _, _ := s.terminal.snapshot()
	if term == nil {
		return
	}
	if components.ToggleButton(c, &s.on, tr(c, "terminalCopyOnSelect")).KeepFocus().Changed() {
		mode := terminal.CopyOnSelectNone
		if s.on {
			mode = terminal.CopyOnSelectClipboard
		}
		term.SetCopyOnSelect(mode)
	}
}

// Close releases gallery-owned native resources. Call on the UI thread when
// its window closes, or after a headless view is no longer used.
func (g *Gallery) Close() {
	g.commands.palette.ReleaseInputSource()
	g.notification.close()
	g.dock.close()
	g.tray.close()
	for _, s := range g.mygoStates {
		s.terminal.close()
		s.webView.close()
		s.audio.close()
		s.system.close()
		s.flyout.close()
	}
}
