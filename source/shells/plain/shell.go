// Package plain provides a single-content native window without the default
// title bar. Applications own content, localization, state and lifetime.
package plain

import (
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/theme"
)

// ShellOptions configures a single-content window. Builders and Shell methods
// run synchronously on the UI thread, after MyGo is ready.
type ShellOptions struct {
	Title string
	// Dimensions are outer DIPs. Defaults are 760x560 with a 360x300 minimum.
	Width, Height, MinWidth, MinHeight int
	// StateKey, when non-empty, has MyGo remember the window's bounds and
	// maximized or full-screen state under this key and restore them at
	// creation. MyGo stores them in its user-data directory, or in the store
	// set with mygo.App.SetWindowStateStore. Use a distinct key per window kind.
	StateKey string
	// Hidden allows the application to install lifecycle handlers before showing.
	Hidden bool
	// Content is automatically inset below the draggable window-control area.
	// Dialogs inherit its drag strip and remain relative to the whole window.
	Content    func(*ui.Context)
	Appearance theme.Mode
	Labels     Labels
}

// Labels names the Windows caption buttons. Empty values use English defaults.
type Labels struct{ Minimize, Maximize, Restore, Close string }

// Shell owns presentation only; Window exposes MyGo's original lifecycle API.
type Shell struct {
	window  *mygo.Window
	host    windowHost
	options ShellOptions
}

type windowHost interface {
	Invalidate()
	Minimize()
	ToggleMaximize()
	IsMaximized() bool
	Close()
}

// NewShell hides the standard title bar, retaining native macOS/Linux window
// controls and composing Onevra caption buttons on Windows. It does not show a
// title or add navigation, persistence, or application lifecycle policies;
// StateKey delegates window geometry to MyGo.
func NewShell(options ShellOptions) *Shell {
	s := newShell(options)
	s.window = mygo.NewWindow(s.nativeOptions())
	s.host = s.window
	if runtime.GOOS == "windows" {
		s.window.OnMaximize(s.refresh)
		s.window.OnUnmaximize(s.refresh)
	}
	return s
}

func newShell(o ShellOptions) *Shell {
	if o.Width <= 0 {
		o.Width = 760
	}
	if o.Height <= 0 {
		o.Height = 560
	}
	if o.MinWidth <= 0 {
		o.MinWidth = 360
	}
	if o.MinHeight <= 0 {
		o.MinHeight = 300
	}
	o.Width, o.Height = max(o.Width, o.MinWidth), max(o.Height, o.MinHeight)
	return &Shell{options: o}
}

func (s *Shell) nativeOptions() mygo.WindowOptions {
	o := s.options
	return mygo.WindowOptions{
		Title: o.Title, Width: o.Width, Height: o.Height, MinWidth: o.MinWidth, MinHeight: o.MinHeight, Hidden: o.Hidden, StateKey: o.StateKey,
		TitleBarStyle: mygo.TitleBarHidden, Frameless: runtime.GOOS == "windows",
		TrafficLightPosition: &mygo.Point{X: 16, Y: 16}, Content: ui.View(s.view),
	}
}

// Window exposes native events and operations. Do not replace its content.
func (s *Shell) Window() *mygo.Window { return s.window }

// SetAppearance changes the theme without recreating content or window state.
func (s *Shell) SetAppearance(mode theme.Mode) { s.options.Appearance = mode; s.refresh() }

// SetLabels updates caption accessibility names without changing their identity.
func (s *Shell) SetLabels(labels Labels) { s.options.Labels = labels; s.refresh() }

func (s *Shell) refresh() {
	if s.host != nil {
		s.host.Invalidate()
	}
}

// HeaderHeight is the frame's native-control clearance and modal drag height.
// Shell already insets ordinary content by this amount.
func HeaderHeight(c *ui.Context) float32 { return max(48, c.TitleBar().Height, c.Theme().Rem(3)) }

func (s *Shell) view(c *ui.Context) {
	theme.Apply(c, s.options.Appearance)
	height := HeaderHeight(c)
	components.WithDialogWindowDrag(c, height, func() {
		ui.Column(c).Fill().Background(c.Theme().Background).Clip().Children(func() {
			ui.Column(c).Key("content").Fill().Padding(height, 0, 0, 0).Children(func() {
				if s.options.Content != nil {
					s.options.Content(c)
				}
			})
			bar := c.TitleBar()
			ui.Row(c).Key("header").Absolute().Top(0).Left(0).Right(0).Height(height).
				Padding(6, max(8, bar.Right), 6, max(8, bar.Left)).AlignItems(ui.Center).DragWindow().Children(func() {
				ui.Spacer(c)
				if runtime.GOOS == "windows" {
					s.captionButtons(c)
				}
			})
		})
	})
}

func label(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func captionButton(c *ui.Context, icon *ui.SVG, name string) ui.Element {
	size := max(32, c.Theme().Rem(2))
	return components.Button(c, "", components.ButtonGhost).Label(name).Size(size, size).Padding(0).Children(func() {
		components.Icon(c, icon).Size(c.Theme().Rem(.875), c.Theme().Rem(.875))
	})
}

func (s *Shell) captionButtons(c *ui.Context) {
	o, w := s.options.Labels, s.host
	ui.Row(c).Key("window-controls").Shrink(0).Gap(4).Children(func() {
		if captionButton(c, lucide.Minus(), label(o.Minimize, "Minimize window")).Clicked() && w != nil {
			w.Minimize()
		}
		icon, name := lucide.Square(), label(o.Maximize, "Maximize window")
		if w != nil && w.IsMaximized() {
			icon, name = lucide.Copy(), label(o.Restore, "Restore window")
		}
		if captionButton(c, icon, name).Clicked() && w != nil {
			w.ToggleMaximize()
		}
		if captionButton(c, lucide.X(), label(o.Close, "Close window")).Clicked() && w != nil {
			w.Close()
		}
	})
}
