package components

import (
	"math"
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

// ShellOptions describes a window and its default three-column page.
// Content builders run synchronously on the UI thread.
type ShellOptions struct {
	Title string
	// Window dimensions and constraints are outer DIPs. Zero uses MyGo defaults.
	Width, Height, MinWidth, MinHeight int
	// Left and Content build the default page. Nil leaves the region empty.
	Left, Content func(*ui.Context)
	// ContentTitle builds the leading main-content titlebar area. It may contain
	// a title, breadcrumbs or actions; nil leaves it blank. Shell clips it to
	// the room between the sidebar controls and hides it on replacement pages.
	ContentTitle func(*ui.Context)
	// LeftToolbar builds a fixed bottom toolbar beneath Left; nil omits it.
	// Use ToolbarIconButton for actions and ui.Spacer for trailing alignment.
	LeftToolbar func(*ui.Context)
	// Right enables the optional inspector and its toggle; nil disables both.
	Right func(*ui.Context)
	// Initial sidebar widths in DIPs; non-positive values use 220 and 280.
	// The left width can be changed by dragging its separator or using arrow
	// keys while the separator is focused. Shell retains it across page changes.
	LeftWidth, RightWidth int
	// The left sidebar starts open; the right sidebar starts closed.
	LeftCollapsed bool
	// Labels allow applications to localize the sidebar controls.
	LeftToggleLabel, RightToggleLabel string
	// LeftResizeLabel names the draggable, keyboard-accessible separator.
	LeftResizeLabel string
}

// Shell owns window presentation, sidebar state and the active page.
// Applications retain ownership of business state and application lifetime.
// Call its methods on the UI thread, after MyGo is ready.
type Shell struct {
	window              *mygo.Window
	host                shellWindow
	options             ShellOptions
	page                func(*ui.Context)
	leftOpen, rightOpen bool
	extension           int
	requestedExtension  int
	leftWidth           float32
	// Native resizing must run outside MyGo's build/layout/paint transaction.
	post             func(func())
	resizeQueued     bool
	pendingExtension int
	resizing         bool
	resizeBaseWidth  int
}

type shellWindow interface {
	Bounds() mygo.Rectangle
	SetBounds(mygo.Rectangle)
	IsMaximized() bool
	IsFullScreen() bool
	Invalidate()
}

// NewShell creates a native window and installs the page container.
// It configures system controls and behind-window blur on macOS; other
// platforms retain their native window background. It does not persist state.
func NewShell(options ShellOptions) *Shell {
	s := newShell(options)
	s.window = mygo.NewWindow(s.nativeOptions())
	s.host = s.window
	// RunOnMain is synchronous on the UI thread. Dispatch from a goroutine so
	// AppKit can redraw at the new size without reentering an unfinished frame.
	s.post = func(fn func()) { go mygo.RunOnMain(fn) }
	return s
}

func newShell(options ShellOptions) *Shell {
	if options.LeftWidth <= 0 {
		options.LeftWidth = 220
	}
	if options.RightWidth <= 0 {
		options.RightWidth = 280
	}
	if options.LeftToggleLabel == "" {
		options.LeftToggleLabel = "Toggle sidebar"
	}
	if options.RightToggleLabel == "" {
		options.RightToggleLabel = "Toggle inspector"
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
// Shell owns its content, title-bar configuration and sidebar width adjustments.
func (s *Shell) Window() *mygo.Window { return s.window }

// SetPage replaces the entire default page, including sidebar controls.
// Nil restores the default page and its previous sidebar expansion state.
// Replacement pages receive an opaque background and the same window-control
// clearance and draggable header as the default page. The inspector width
// adjustment animates across subsequent frames.
func (s *Shell) SetPage(page func(*ui.Context)) {
	s.page = page
	s.refresh()
}

// LeftExpanded reports the requested state, including during an animation.
func (s *Shell) LeftExpanded() bool { return s.leftOpen }

// RightExpanded reports the requested state, retained on replacement pages.
func (s *Shell) RightExpanded() bool { return s.rightOpen }

// SetLeftExpanded animates the interior layout without resizing the window.
func (s *Shell) SetLeftExpanded(open bool) {
	s.leftOpen = open
	s.refresh()
}

// SetRightExpanded animates the window extension to the right, leaving its origin and
// height unchanged. Disabled inspectors ignore this method. Screen constraints
// may limit expansion. In maximized/full-screen windows, resizing is deferred
// until the window returns to its normal state.
func (s *Shell) SetRightExpanded(open bool) {
	s.rightOpen = open && s.options.Right != nil
	s.refresh()
}

func (s *Shell) refresh() {
	if s.host != nil {
		s.host.Invalidate()
	}
}

// requestWidth coalesces resize requests from MyGo's build passes. A queued
// callback reads the latest target, including reversals and page changes.
func (s *Shell) requestWidth(requested int) {
	if s.post == nil {
		// Headless fixtures can apply the step without a native event loop.
		s.syncWidth(requested)
		return
	}
	s.pendingExtension = requested
	if s.resizeQueued || requested == s.requestedExtension || s.host.IsMaximized() || s.host.IsFullScreen() {
		return
	}
	s.resizeQueued = true
	s.post(func() {
		s.syncWidth(s.pendingExtension)
		s.resizeQueued = false
		s.refresh()
	})
}

// syncWidth applies one animation step outside rendering, retaining only the
// width the OS actually granted. User resizing remains part of the base width.
func (s *Shell) syncWidth(requested int) {
	if s.host == nil || s.host.IsMaximized() || s.host.IsFullScreen() {
		return
	}
	delta := requested - s.requestedExtension
	if delta == 0 {
		return
	}
	before := s.host.Bounds()
	after := before
	if delta < 0 {
		delta = -min(-delta, s.extension)
	}
	after.Width += delta
	// SetBounds can synchronously ask MyGo to draw. Publish the layout target
	// first; those frames measure the granted extension from the native bounds.
	s.requestedExtension = requested
	s.resizeBaseWidth = before.Width - s.extension
	s.resizing = true
	if after.Width != before.Width {
		s.host.SetBounds(after)
	}
	s.extension = max(0, s.host.Bounds().Width-s.resizeBaseWidth)
	s.resizing = false
}

const sidebarDuration = 220 * time.Millisecond

func openFraction(open bool) float32 {
	if open {
		return 1
	}
	return 0
}

func (s *Shell) view(c *ui.Context) {
	if runtime.GOOS == "darwin" {
		c.Root().Background(ui.Transparent)
	}
	background := c.Theme().Background.Alpha(1)
	inspectorBackground := background
	if runtime.GOOS == "darwin" {
		inspectorBackground = background.Alpha(0.8)
	}
	root := ui.Column(c).Fill().Clip()
	if runtime.GOOS != "darwin" {
		root.Background(background)
	}
	// Snapshot the page for this frame; callbacks can change it for the next.
	page := s.page
	leftTarget := openFraction(s.leftOpen)
	rightTarget := openFraction(s.rightOpen && page == nil)
	leftProgress := root.Animate("left", leftTarget, sidebarDuration)
	rightProgress := root.Animate("right", rightTarget, sidebarDuration)
	// Also snap an animation already running when Reduce Motion is enabled.
	if c.Preferences().ReduceMotion {
		leftProgress, rightProgress = leftTarget, rightTarget
	}
	requested := int(math.Round(float64(rightProgress * float32(s.options.RightWidth))))
	s.requestWidth(requested)
	bar := c.TitleBar()
	headerHeight := max(48, bar.Height, c.Theme().Rem(3))
	width, _ := c.Size()
	baseWidth := width
	visibleRight := float32(requested)
	if s.host != nil && !s.host.IsMaximized() && !s.host.IsFullScreen() {
		extension := s.extension
		if s.resizing {
			extension = max(0, s.host.Bounds().Width-s.resizeBaseWidth)
		}
		baseWidth = max(0, width-float32(extension))
		// Use the applied request, never the animation's next native size.
		visibleRight = float32(s.requestedExtension)
	}
	leftWidth := min(s.leftWidth, baseWidth*0.4)
	rightWidth := float32(s.options.RightWidth)
	// If the OS cannot grant enough room, reveal inside the existing window
	// while keeping some space for the main region. Inspector content stays fixed.
	visibleRight = min(visibleRight, max(0, width-leftWidth*leftProgress-min(160, baseWidth*0.6)))
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
					ui.Column(c).Key("left").Width(leftWidth * leftProgress).Shrink(0).FillHeight().Clip().Children(func() {
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
				if visibleRight > 0 && s.options.Right != nil {
					// Reveal a full-width inspector through the growing window edge,
					// rather than reflowing its contents at each animation step.
					ui.Column(c).Key("right").Width(visibleRight).Shrink(0).FillHeight().Background(inspectorBackground).Clip().Children(func() {
						ui.Column(c).Width(rightWidth).Shrink(0).FillHeight().Padding(headerHeight, 0, 0, 0).
							Disabled(!s.rightOpen).Children(func() { s.options.Right(c) })
						shellSeparator(c)
					})
				}
			})
		})
		if page == nil && s.leftOpen && leftProgress == 1 {
			s.resizeHandle(c, leftWidth, baseWidth*0.4, headerHeight)
		}
		// An unpainted header overlays the full-height columns. Only its controls
		// occupy the safe area; column content is kept below it by top padding.
		leftInset := max(8, bar.Left)
		if runtime.GOOS == "darwin" && bar.Left > 0 {
			leftInset += 16
		}
		rightInset := max(8, bar.Right)
		if page == nil {
			// Anchor the inspector control to the main content's right edge.
			// In tight windows, keep both controls inside the native safe area.
			rightInset = max(rightInset, min(visibleRight+8, width-leftInset-2*max(32, c.Theme().Rem(2))-8))
		}
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
				if s.options.Right != nil {
					ui.Toolbar(c, func() {
						if ToolbarIconButton(c, lucide.PanelRight(), s.options.RightToggleLabel).Clicked() {
							s.SetRightExpanded(!s.rightOpen)
						}
					}).Padding(0).Width(max(32, c.Theme().Rem(2))).Shrink(0)
				}
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
	if next != width {
		s.leftWidth = max(min(120, maxWidth), min(next, maxWidth))
		c.Invalidate()
	}
	handle.FocusRing(handle.FocusVisible())
}

// Separators overlay the surface; they never reserve layout space.
func shellSeparator(c *ui.Context) {
	ui.Box(c).Key("separator").Absolute().Left(0).Top(0).Bottom(0).Width(1).Background(c.Theme().Border)
}
