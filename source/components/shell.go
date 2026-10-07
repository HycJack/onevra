package components

import (
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
	"github.com/oneo-me/oneo-go/theme"
)

// ShellOptions describes a window and its default sidebar-and-content page.
// Content builders run synchronously on the UI thread.
type ShellOptions struct {
	Title string
	// Window dimensions and constraints are outer DIPs. Zero uses MyGo defaults.
	Width, Height, MinWidth, MinHeight int
	// Left and Content build the default page. Nil leaves the region empty.
	Left, Content func(*ui.Context)
	// ContentTitle builds the leading main-content titlebar area. It may contain
	// a title, breadcrumbs or actions; nil leaves it blank. Shell clips it to
	// the room after the sidebar control and hides it on replacement pages.
	ContentTitle func(*ui.Context)
	// LeftToolbar builds a fixed bottom toolbar beneath Left; nil omits it.
	// Use ToolbarIconButton for actions and ui.Spacer for trailing alignment.
	LeftToolbar func(*ui.Context)
	// Default left sidebar width in DIPs; non-positive values use 220.
	// Clamped to LeftMinWidth/LeftMaxWidth at creation. Dragging or focused
	// arrow keys adjust the width; double-clicking the separator restores it.
	// Shell retains the chosen width across toggles and page changes.
	LeftWidth int
	// LeftMinWidth defaults to 120 DIPs when non-positive. A non-positive
	// LeftMaxWidth adds no fixed limit; a positive maximum below the minimum
	// is raised to the minimum. Visible width is also capped at 40% of the
	// viewport, even below the minimum, without discarding the chosen width.
	LeftMinWidth, LeftMaxWidth int
	// The left sidebar starts open.
	LeftCollapsed bool
	// Labels allow applications to localize the sidebar controls.
	LeftToggleLabel string
	// LeftResizeLabel names the draggable, keyboard-accessible separator.
	LeftResizeLabel string
	// Windows caption button labels; empty values use English defaults.
	MinimizeLabel, MaximizeLabel, RestoreLabel, CloseLabel string
	// Appearance selects the window theme: Auto follows the system
	// appearance, Light and Dark pin oneo-go's built-in palettes.
	// It applies before any content is built; the zero value is Auto.
	Appearance theme.Mode
}

// Shell owns window presentation, sidebar state and the active page.
// Applications retain ownership of business state and application lifetime.
// Call its methods on the UI thread, after MyGo is ready.
type Shell struct {
	window    *mygo.Window
	host      shellWindow
	options   ShellOptions
	page      func(*ui.Context)
	leftOpen  bool
	leftWidth float32
}

type shellWindow interface {
	Invalidate()
}

// NewShell creates a native window and installs the page container.
// It configures system controls and behind-window blur on macOS, and custom
// caption buttons on Windows. Other platforms retain their native controls.
// Windows and Linux use an opaque background. It does not persist state.
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

func newShell(options ShellOptions) *Shell {
	if options.LeftWidth <= 0 {
		options.LeftWidth = 220
	}
	if options.LeftMinWidth <= 0 {
		options.LeftMinWidth = 120
	}
	options.LeftWidth = max(options.LeftWidth, options.LeftMinWidth)
	if options.LeftMaxWidth > 0 {
		options.LeftMaxWidth = max(options.LeftMaxWidth, options.LeftMinWidth)
		options.LeftWidth = min(options.LeftWidth, options.LeftMaxWidth)
	}
	if options.LeftToggleLabel == "" {
		options.LeftToggleLabel = "Toggle sidebar"
	}
	if options.LeftResizeLabel == "" {
		options.LeftResizeLabel = "Resize sidebar"
	}
	return &Shell{options: options, leftOpen: !options.LeftCollapsed, leftWidth: float32(options.LeftWidth)}
}

func (s *Shell) nativeOptions() mygo.WindowOptions {
	o := s.options
	return mygo.WindowOptions{
		Title: o.Title, Width: o.Width, Height: o.Height,
		MinWidth: o.MinWidth, MinHeight: o.MinHeight,
		TitleBarStyle:        mygo.TitleBarHidden,
		Frameless:            runtime.GOOS == "windows",
		TrafficLightPosition: &mygo.Point{X: 16, Y: 16},
		Transparent:          runtime.GOOS == "darwin",
		Vibrancy:             shellVibrancy(),
		Content:              ui.View(s.view),
	}
}

func shellVibrancy() mygo.Vibrancy {
	if runtime.GOOS == "darwin" {
		return mygo.VibrancyTitlebar
	}
	return mygo.VibrancyNone
}

// Window exposes the original MyGo handle for events and lifecycle operations.
// Shell owns its content, title-bar configuration and sidebar layout.
func (s *Shell) Window() *mygo.Window { return s.window }

// SetPage replaces the entire default page, including sidebar controls.
// Nil restores the default page and its previous sidebar expansion state.
// Replacement pages receive an opaque background and the same window-control
// clearance and draggable header as the default page.
func (s *Shell) SetPage(page func(*ui.Context)) {
	s.page = page
	s.refresh()
}

// LeftExpanded reports the requested state, including during an animation.
func (s *Shell) LeftExpanded() bool { return s.leftOpen }

// SetLeftExpanded animates the interior layout without resizing the window.
func (s *Shell) SetLeftExpanded(open bool) {
	s.leftOpen = open
	s.refresh()
}

func (s *Shell) refresh() {
	if s.host != nil {
		s.host.Invalidate()
	}
}

const sidebarDuration = 220 * time.Millisecond

func openFraction(open bool) float32 {
	if open {
		return 1
	}
	return 0
}

func (s *Shell) view(c *ui.Context) {
	// The theme must be installed before any widget is created.
	theme.Apply(c, s.options.Appearance)
	if runtime.GOOS == "darwin" {
		c.Root().Background(ui.Transparent)
	}
	background := c.Theme().Background.Alpha(1)
	root := ui.Column(c).Fill().Clip()
	if runtime.GOOS != "darwin" {
		root.Background(background)
	}
	// Snapshot the page for this frame; callbacks can change it for the next.
	page := s.page
	leftTarget := openFraction(s.leftOpen)
	leftProgress := root.Animate("left", leftTarget, sidebarDuration)
	// Also snap an animation already running when Reduce Motion is enabled.
	if c.Preferences().ReduceMotion {
		leftProgress = leftTarget
	}
	bar := c.TitleBar()
	headerHeight := max(48, bar.Height, c.Theme().Rem(3))
	width, _ := c.Size()
	maxLeftWidth := width * 0.4
	if s.options.LeftMaxWidth > 0 {
		maxLeftWidth = min(maxLeftWidth, float32(s.options.LeftMaxWidth))
	}
	leftWidth := min(s.leftWidth, maxLeftWidth)
	root.Children(func() {
		ui.Column(c).Key(page != nil).Fill().Clip().Children(func() {
			if page != nil {
				ui.Column(c).Fill().Padding(headerHeight, 0, 0, 0).Background(background).Children(func() { page(c) })
				return
			}
			ui.Row(c).Fill().Children(func() {
				if leftProgress > 0 {
					// Only the visible slice changes width. Sidebar content stays put
					// while the opaque main region covers or uncovers it.
					sidebar := ui.Column(c).Key("left").Width(leftWidth * leftProgress).Shrink(0).FillHeight().Clip()
					if runtime.GOOS == "windows" {
						sidebar.Background(c.Theme().Surface.Alpha(1))
					}
					sidebar.Children(func() {
						ui.Column(c).Width(leftWidth).Shrink(0).FillHeight().Padding(headerHeight, 0, 0, 0).
							Disabled(!s.leftOpen).Children(func() {
							ui.Column(c).FillWidth().Grow(1).MinHeight(0).Clip().Children(func() {
								if s.options.Left != nil {
									s.options.Left(c)
								}
							})
							if s.options.LeftToolbar != nil {
								ui.Toolbar(c, func() { s.options.LeftToolbar(c) }).
									FillWidth().Shrink(0).Padding(8)
							}
						})
					})
				}
				ui.Column(c).Key("content").Grow(1).MinWidth(0).FillHeight().Background(background).Clip().Children(func() {
					ui.Column(c).Fill().Padding(headerHeight, 0, 0, 0).Children(func() {
						if s.options.Content != nil {
							s.options.Content(c)
						}
					})
					if leftProgress > 0 {
						shellSeparator(c)
					}
				})
			})
		})
		if page == nil && s.leftOpen && leftProgress == 1 {
			s.resizeHandle(c, leftWidth, maxLeftWidth, headerHeight)
		}
		// An unpainted header overlays the full-height columns. Only its controls
		// occupy the safe area; column content is kept below it by top padding.
		leftInset := max(8, bar.Left)
		if runtime.GOOS == "darwin" && bar.Left > 0 {
			leftInset += 16
		}
		rightInset := max(8, bar.Right)
		ui.Row(c).Key("header").Absolute().Top(0).Left(0).Right(0).Height(headerHeight).
			AlignItems(ui.Center).Padding(6, rightInset, 6, leftInset).
			DragWindow().Children(func() {
			if page == nil {
				if ToolbarIconButton(c, lucide.PanelLeft(), s.options.LeftToggleLabel).Clicked() {
					s.SetLeftExpanded(!s.leftOpen)
				}
				// Align the title with the main surface, but keep it beyond the
				// leading toggle when the sidebar is narrow or collapsed.
				titleInset := max(8, leftWidth*leftProgress+16-leftInset-max(32, c.Theme().Rem(2)))
				ui.Row(c).Grow(1).MinWidth(0).FillHeight().Clip().Padding(0, 8, 0, titleInset).Children(func() {
					if s.options.ContentTitle != nil {
						s.options.ContentTitle(c)
					}
				})
			}
			if runtime.GOOS == "windows" {
				if page != nil {
					ui.Spacer(c)
				}
				s.captionButtons(c)
			}
		})
	})
}

// resizeHandle uses MyGo's pointer capture and keyboard input, overlaid on
// the separator rather than inserting another layout column.
func (s *Shell) resizeHandle(c *ui.Context, width, maxWidth, top float32) {
	handle := ui.Box(c).Key("left-resize").Absolute().Left(width - 3).Top(top).Bottom(0).Width(6).
		Cursor(ui.CursorResizeEW).Focusable().Role(ui.RoleSplitter).Label(s.options.LeftResizeLabel)
	next := width
	if dx, _, held := handle.Dragged(); held {
		handle.Focus()
		next += dx
	}
	if handle.Shortcut(0, ui.KeyLeft) {
		next -= 10
	}
	if handle.Shortcut(0, ui.KeyRight) {
		next += 10
	}
	if handle.DoubleClicked() {
		s.leftWidth = float32(s.options.LeftWidth)
		c.Invalidate()
	} else if next != width {
		s.leftWidth = max(float32(s.options.LeftMinWidth), min(next, maxWidth))
		c.Invalidate()
	}
	handle.FocusRing(handle.FocusVisible())
}

// Separators overlay the surface; they never reserve layout space.
func shellSeparator(c *ui.Context) {
	ui.Box(c).Key("separator").Absolute().Left(0).Top(0).Bottom(0).Width(1).Background(c.Theme().Border)
}
