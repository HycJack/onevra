package components

import (
	"time"

	"github.com/egoist/mygo/ui"
)

const tabsIndicatorDuration = 180 * time.Millisecond

// TabsVariant selects the look of a tab list.
type TabsVariant int

const (
	// TabsUnderline marks the chosen tab with a shared accent underline that
	// animates to it. This is the default.
	TabsUnderline TabsVariant = iota
	// TabsSegmented draws the tabs as the segments of a rounded track, with
	// the chosen one raised on a shadowed, lighter thumb, as macOS draws a
	// segmented control.
	TabsSegmented
)

// Tabs creates a row of tabs showing labels, of which *selected is the index
// of the one chosen. A click chooses a tab, as do the arrows, Home and End
// while one has the keyboard focus, which follows the choice. Reduced Motion
// skips transitions. List's Changed reports a new choice.
//
// TabsUnderline animates a shared accent underline to the chosen tab and
// keeps every label's weight, so selection never moves the labels.
// TabsSegmented raises the chosen tab on a shadowed thumb over a capsule
// track and shows its label in a heavier weight, like a macOS segmented
// control.
func Tabs(c *ui.Context, selected *int, variant TabsVariant, labels ...string) ui.TabsParts {
	tabs := ui.TabsBase(c, selected, len(labels))
	if variant == TabsSegmented {
		return segmentedTabs(c, tabs, selected, labels)
	}
	return underlinedTabs(c, tabs, selected, labels)
}

// underlinedTabs puts the accent underline under the chosen tab.
func underlinedTabs(c *ui.Context, tabs ui.TabsParts, selected *int, labels []string) ui.TabsParts {
	t := c.Theme()
	tabs.List.Gap(t.Space(1))
	var active ui.Element
	tabs.List.Children(func() {
		for i, label := range labels {
			tab := dimmed(tabs.Tab(i).Padding(t.Space(2), t.Space(3)).FocusRing(false))
			on := i == *selected
			fg := t.Text
			if !on {
				fg = fg.Alpha(0.6)
			}
			tab.TextColor(fg)
			if on {
				active = tab
			}
			tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if tab.FocusVisible() {
					p.FocusRing(r, [4]float32{t.Radius, t.Radius, t.Radius, t.Radius})
				}
			})
			tab.Children(func() { Text(c, label).SingleLine() })
		}
	})
	if len(labels) > 0 {
		tabs.List.DrawOver(func(p *ui.Painter, r ui.Rect) {
			// Read current layout during painting, and animate list-relative
			// geometry so moving the entire control does not move its underline.
			bounds := active.LayoutBounds()
			in, h := t.Space(1.5), t.Space(0.5)
			x := tabs.List.Animate("indicator-x", bounds.X-r.X+in, tabsIndicatorDuration)
			w := tabs.List.Animate("indicator-width", max(0, bounds.W-2*in), tabsIndicatorDuration)
			if c.Preferences().ReduceMotion {
				x, w = bounds.X-r.X+in, max(0, bounds.W-2*in)
			}
			p.Fill(ui.Rect{X: r.X + x, Y: bounds.Y + bounds.H - h, W: w, H: h}, t.Accent, h/2)
		})
	}
	return tabs
}

// segmentedTabs puts the chosen tab on a raised thumb inside a recessed track,
// as macOS draws a segmented control.
func segmentedTabs(c *ui.Context, tabs ui.TabsParts, selected *int, labels []string) ui.TabsParts {
	t := c.Theme()
	// The track is recessed from the page and the thumb raised from it, with a
	// cast shadow: the same raised face a toolbar button uses.
	track := toolbarTrack(t)
	thumb, _, _, shadow := toolbarFace(t)
	// FitRadii clamps any radius above half the height, so this makes the track
	// a capsule whatever its height and the theme's spacing.
	pill := t.Space(20)
	tabs.List.Padding(t.Space(0.5)).Gap(t.Space(0.5)).Radius(pill).Background(track).AlignSelf(ui.Start)
	var active ui.Element
	tabs.List.Children(func() {
		for i, label := range labels {
			tab := dimmed(tabs.Tab(i).Padding(t.Space(1.5), t.Space(3)).FocusRing(false))
			on := i == *selected
			fg := t.Text
			if !on {
				fg = fg.Alpha(0.6)
			}
			tab.TextColor(fg)
			if on {
				active = tab
			}
			tab.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if tab.FocusVisible() {
					p.FocusRing(r, [4]float32{pill, pill, pill, pill})
				}
			})
			tab.Children(func() {
				text := Text(c, label).SingleLine()
				if on {
					text.FontWeight(600)
				}
			})
		}
	})
	if len(labels) > 0 {
		// The thumb belongs behind the labels, so draw it with the track rather
		// than over the tabs.
		tabs.List.Draw(func(p *ui.Painter, r ui.Rect) {
			bounds := active.LayoutBounds()
			x := tabs.List.Animate("thumb-x", bounds.X-r.X, tabsIndicatorDuration)
			w := tabs.List.Animate("thumb-width", bounds.W, tabsIndicatorDuration)
			if c.Preferences().ReduceMotion {
				x, w = bounds.X-r.X, bounds.W
			}
			thumbRect := ui.Rect{X: r.X + x, Y: bounds.Y, W: w, H: bounds.H}
			radius := bounds.H / 2
			p.Shadow(thumbRect, radius, 0, 1, 2, 0, shadow)
			p.Fill(thumbRect, thumb, radius)
		})
	}
	return tabs
}
