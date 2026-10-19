package components

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// ScrollEdgeStyle selects what separates a ScrollEdgeBar from the content
// scrolling under it.
type ScrollEdgeStyle int

const (
	// ScrollEdgeGlass sets the whole bar on one capsule of Liquid Glass and
	// leaves the content beside it untouched. This is the default.
	ScrollEdgeGlass ScrollEdgeStyle = iota
	// ScrollEdgeSoft fades the content into the background toward the edge,
	// without a blur, as macOS's soft scroll edge effect.
	ScrollEdgeSoft
	// ScrollEdgeHard frosts the content evenly under the bar and draws a
	// hairline where it ends, as macOS's hard scroll edge effect.
	ScrollEdgeHard
	// ScrollEdgeBlur blurs the content progressively, the more the nearer
	// the edge.
	ScrollEdgeBlur
)

// ScrollEdgeBarOptions configures a ScrollEdgeBar.
type ScrollEdgeBarOptions struct {
	Style ScrollEdgeStyle
	// Bottom docks the bar at the bottom edge, with its effect rising from
	// it, rather than at the top.
	Bottom bool
	// Height is the bar's height in DIPs; zero is 52, which fits 40-DIP
	// buttons.
	Height float32
	// Inset is the bar's distance from its edge and from the sides in DIPs;
	// zero is three theme spacing units.
	Inset float32
	// Glass is the bar's material in ScrollEdgeGlass, such as glass.Clear
	// or a tint; Interactive is ignored. The zero value is regular glass.
	Glass glass.Glass
	// Background is what the content scrolls over, which ScrollEdgeSoft and
	// ScrollEdgeHard fade toward; unset, the theme's Background.
	Background ui.Color
}

// Extra room the effects take beyond the bar, as MyGo's glass gallery and
// macOS size them: a soft edge fades out 10 DIPs past the bar, a
// progressive blur 24 DIPs past it.
const (
	scrollEdgeSoftExtra = 10
	scrollEdgeBlurExtra = 24
	scrollEdgeBlur      = 6
)

func (o ScrollEdgeBarOptions) resolved(c *ui.Context) ScrollEdgeBarOptions {
	if o.Height <= 0 {
		o.Height = 52
	}
	if o.Inset <= 0 {
		o.Inset = c.Theme().Space(3)
	}
	return o
}

// ScrollEdgeBarExtent returns how far from its edge a ScrollEdgeBar with
// options reaches, in DIPs: the scroll content's padding on that edge, plus
// any spacing, and its ScrollbarInsets keep content and scroll bar clear of
// the bar.
func ScrollEdgeBarExtent(c *ui.Context, options ScrollEdgeBarOptions) float32 {
	o := options.resolved(c)
	return o.Inset + o.Height
}

// ScrollEdgeBar floats a bar of controls over the top or bottom edge of
// scrolling content, such as a title bar, toolbar or bottom bar, and
// separates it from the content in the chosen style. Build it after the
// content, as a child of a bounded, clipping container that the content
// fills: materials sample what was painted beneath them, and the bar and
// its effect position themselves absolutely within that container.
//
//	ui.Box(c).Fill().Clip().Children(func() {
//		inset := components.ScrollEdgeBarExtent(c, options)
//		components.ScrollView(c, components.ScrollViewOptions{}).Fill().
//			Padding(inset+12, 16, 16).ScrollbarInsets(inset, 0, 0).Children(content)
//		components.ScrollEdgeBar(c, options, toolbar).Label("Toolbar")
//	})
//
// The returned element is the bar, a centered row that content builds the
// children of; the application owns them, their actions and the bar's
// accessible label. In styles other than ScrollEdgeGlass, the bar is
// transparent, so the controls should carry their own surfaces, such as
// interactive glass. The effect layer never takes input, so content under
// it beside the bar stays reachable.
func ScrollEdgeBar(c *ui.Context, options ScrollEdgeBarOptions, content func()) ui.Element {
	o := options.resolved(c)
	extent := o.Inset + o.Height
	if o.Style != ScrollEdgeGlass {
		edge := ui.Box(c).Absolute().Left(0).Right(0).PassThrough()
		if o.Bottom {
			edge.Bottom(0)
		} else {
			edge.Top(0)
		}
		switch o.Style {
		case ScrollEdgeHard:
			edge.Height(extent).Material(glass.ScrollEdge{Hard: true, Bottom: o.Bottom, Background: o.Background})
		case ScrollEdgeBlur:
			mask := &ui.LinearGradient{From: ui.RGB(0, 0, 0), To: ui.Transparent, Angle: 180, Start: 0.3, End: 1}
			if o.Bottom {
				mask.Angle = 0
			}
			edge.Height(extent + scrollEdgeBlurExtra).Material(glass.Blur{Radius: scrollEdgeBlur, Mask: mask})
		default:
			edge.Height(extent + scrollEdgeSoftExtra).Material(glass.ScrollEdge{Bottom: o.Bottom, Background: o.Background})
		}
	}
	bar := ui.Row(c).Absolute().Left(o.Inset).Right(o.Inset).Height(o.Height).
		Padding(0, 6).Gap(6).AlignItems(ui.Center).Radius(o.Height / 2)
	if o.Bottom {
		bar.Bottom(o.Inset)
	} else {
		bar.Top(o.Inset)
	}
	if o.Style == ScrollEdgeGlass {
		material := o.Glass
		material.Interactive = false
		bar.Material(material)
	} else {
		bar.PassThrough()
	}
	if content != nil {
		bar.Children(content)
	}
	return bar
}
