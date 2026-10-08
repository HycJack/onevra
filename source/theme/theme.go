package theme

import "github.com/egoist/mygo/ui"

// Mode selects which appearance a window uses.
type Mode int

const (
	// Auto follows the system's appearance.
	Auto Mode = iota
	// Light pins the built-in light palette.
	Light
	// Dark pins the built-in dark palette.
	Dark
)

func (m Mode) String() string {
	switch m {
	case Light:
		return "light"
	case Dark:
		return "dark"
	default:
		return "auto"
	}
}

// Palettes follow the Light and Shadow (light-shadow) design style from
// oneo-me/desgin-styles: a warm white light mode and a deep navy dark mode.
// Accent hover/pressed states derive the same way MyGo derives them, and the
// control radius follows the style; spacing and font size keep following MyGo.

// LightPalette returns Onevra's light palette.
func LightPalette() *ui.Theme {
	t := ui.LightTheme()
	accent := ui.Hex("#bd481b")
	black := ui.RGB(0, 0, 0)
	t.Background = ui.Hex("#fcfbf9")
	t.Surface = ui.Hex("#fefefd")
	t.SurfaceHover = ui.Hex("#eceae4")
	t.SurfacePressed = ui.Hex("#d9d7d2")
	t.Border = ui.Hex("#deded5")
	t.Text = ui.Hex("#33332f")
	t.TextMuted = ui.Hex("#67675f")
	t.Accent = accent
	t.AccentHover = accent.Mix(black, 0.12)
	t.AccentPressed = accent.Mix(black, 0.25)
	t.AccentText = ui.Hex("#ffffff")
	t.Danger = ui.Hex("#b4444c")
	// The style defines no warning token; amber matches its warm palette.
	t.Warning = ui.Hex("#b45309")
	t.Success = ui.Hex("#059669")
	t.Selection = ui.RGBA(230, 101, 43, 0.22)
	t.Focus = ui.RGBA(185, 71, 22, 0.55)
	t.Scrollbar = t.Text.Alpha(.30)
	t.Radius = 8
	return t
}

// DarkPalette returns Onevra's dark palette.
// MyGo widgets branch on the Dark flag, so it stays set.
func DarkPalette() *ui.Theme {
	t := ui.DarkTheme()
	t.Dark = true
	accent := ui.Hex("#9cbff9")
	black := ui.RGB(0, 0, 0)
	white := ui.RGB(255, 255, 255)
	t.Background = ui.Hex("#192332")
	t.Surface = ui.Hex("#243247")
	t.SurfaceHover = ui.Hex("#30435e")
	t.SurfacePressed = ui.Hex("#2b3c54")
	t.Border = ui.Hex("#2e3e55")
	t.Text = ui.Hex("#e7edf8")
	t.TextMuted = ui.Hex("#a4b4cc")
	t.Accent = accent
	t.AccentHover = accent.Mix(white, 0.2)
	t.AccentPressed = accent.Mix(black, 0.15)
	t.AccentText = ui.Hex("#122540")
	t.Danger = ui.Hex("#ef9aa1")
	// The style defines no warning token; amber matches its palette.
	t.Warning = ui.Hex("#fbbf24")
	t.Success = ui.Hex("#34d399")
	t.Selection = ui.RGBA(113, 162, 236, 0.25)
	t.Focus = ui.RGBA(131, 175, 244, 0.55)
	t.Scrollbar = t.Text.Alpha(.35)
	t.Radius = 8
	return t
}

// Apply installs the theme for the current frame. Call it before creating
// any widgets, typically at the start of the root content builder: SetTheme
// only affects the window root and the widgets created after the call.
//
// Auto reads the frame's default theme first, which still follows the
// system's appearance, then installs the matching built-in palette. MyGo
// reruns the frame when the system appearance changes, so Auto tracks it
// without polling. A pinned Light or Dark ignores appearance changes.
//
// A custom theme does not follow the desktop preferences' colors; the
// operating system still owns window chrome.
func Apply(c *ui.Context, m Mode) {
	dark := c.Theme().Dark
	switch m {
	case Light:
		c.SetTheme(LightPalette())
	case Dark:
		c.SetTheme(DarkPalette())
	default:
		if dark {
			c.SetTheme(DarkPalette())
		} else {
			c.SetTheme(LightPalette())
		}
	}
}
