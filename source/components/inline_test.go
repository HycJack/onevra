package components

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/theme"
)

func TestRichTextComposition(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, width := range []int{160, 480} {
			build := func(rich func(*ui.Context, ...ui.Span) ui.Element) *ui.Tester {
				return ui.NewTester(func(c *ui.Context) {
					c.SetTheme(theme.Palette(dark))
					rich(c, ui.Span{Text: "Hello ", Weight: 600}, ui.Span{Text: "世界 ", Color: c.Theme().Accent, Underline: true},
						ui.Span{Text: "italic ", Italic: true, Background: c.Theme().Surface}).
						FontSize(18).Children(func() {
						Text(c, "with ")
						Button(c, "a link", ButtonLink)
						Text(c, " and wrapping text.")
					})
				}, width, 200)
			}
			got, want := build(RichText), build(ui.RichText)
			if !bytes.Equal(got.Image().Pix, want.Image().Pix) {
				t.Fatal("inline composition differs from the text engine")
			}
		}
	}
}

func TestBadge(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for tone := BadgeNeutral; tone <= BadgeDanger; tone++ {
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				Badge(c, "New", BadgeOptions{Tone: tone, Icon: lucide.Check()}).Label("status")
			}, 200, 100)
			r, ok := tt.Find("status")
			if !ok || r.W <= 0 || r.H <= 0 || r.W >= 200 || !tt.HasText("New") {
				t.Fatalf("missing badge: rect=%v found=%v texts=%v", r, ok, tt.Texts())
			}
		}
	}
}

func TestBadgeMorph(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		text := "9"
		tt := ui.NewTester(func(c *ui.Context) {
			Badge(c, text).Label("count")
		}, 300, 100)
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce})
		tt.Frame()
		narrow, _ := tt.Find("count")
		text = "1,234 items"
		tt.Frame()
		during, _ := tt.Find("count")
		tt.Advance(time.Second)
		tt.Frame()
		wide, _ := tt.Find("count")
		if wide.W <= narrow.W || !tt.HasText(text) || tt.HasText("9") {
			t.Fatalf("settled badge: %v → %v, texts=%v", narrow, wide, tt.Texts())
		}
		if moving := during.W < wide.W; moving == reduce {
			t.Fatalf("reduce motion %v: width %v during, %v settled", reduce, during.W, wide.W)
		}
	}
}

func TestBreadcrumbsActivation(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		chosen, changes := -1, 0
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Disabled(disabled).Children(func() {
				b := Breadcrumbs(c, []string{"Home", "Library", "Current"}, &chosen)
				if b.Changed() {
					changes++
				}
			})
		}, 300, 100)
		for range 2 {
			if err := tt.Click("Home"); err != nil {
				t.Fatal(err)
			}
		}
		if disabled {
			if chosen != -1 || changes != 0 {
				t.Fatal("disabled navigation")
			}
		} else if chosen != 0 || changes != 2 {
			t.Fatalf("activation: %d/%d", chosen, changes)
		}
		if err := tt.Click("Current"); err != nil {
			t.Fatal(err)
		}
		if changes > 2 {
			t.Fatal("current item navigated")
		}
	}
}

func TestBreadcrumbsLayoutStable(t *testing.T) {
	var heights, widths []float32
	for _, items := range [][]string{{"Home"}, {"Home", "Current"}} {
		chosen := -1
		var h float32
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).AlignItems(ui.Start).Children(func() {
				Breadcrumbs(c, items, &chosen).DrawOver(func(_ *ui.Painter, r ui.Rect) { h = r.H })
			})
		}, 300, 100)
		home, ok := tt.Find("Home")
		if !ok {
			t.Fatal("missing Home")
		}
		heights, widths = append(heights, h), append(widths, home.W)
	}
	// Matching padding places current and linked Home text identically; the
	// current item's heavier weight alone may widen its glyphs slightly.
	if d := widths[0] - widths[1]; heights[0] == 0 || heights[0] != heights[1] || d < -1 || d > 1 {
		t.Fatalf("heights %v, Home widths %v", heights, widths)
	}
}

func TestRatingInput(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		value, changes := 0, 0
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Disabled(disabled).Children(func() {
				r := Rating(c, &value, 5)
				r.Label("rating")
				if r.Changed() {
					changes++
				}
			})
		}, 300, 100)
		r, _ := tt.Find("rating")
		tt.ClickAt(r.X+10, r.Y+r.H/2)
		if !disabled && (value != 1 || changes != 1) {
			t.Fatalf("click: %d/%d", value, changes)
		}
		tt.ClickAt(r.X+10, r.Y+r.H/2)
		if value != 0 {
			t.Fatal("selected star did not clear")
		}
		tt.Key(0, ui.KeyEnd)
		if !disabled && value != 5 {
			t.Fatal("End")
		}
		tt.Key(0, ui.KeyRight)
		if !disabled && changes != 3 {
			t.Fatal("clamped change")
		}
		tt.Key(0, ui.KeyLeft)
		if !disabled && value != 4 {
			t.Fatal("Left")
		}
		tt.Key(0, ui.KeyHome)
		if value != 0 || disabled && changes != 0 {
			t.Fatal("Home or disabled input")
		}
	}
}

func TestRatingBounds(t *testing.T) {
	for _, maximum := range []int{-1, 0, 5} {
		value := 10
		ui.NewTester(func(c *ui.Context) {
			r := Rating(c, &value, maximum)
			if r.Changed() {
				t.Fatal("clamping reported user input")
			}
		}, 200, 100)
		if value != max(0, maximum) {
			t.Fatal("not clamped")
		}
	}
}

// Compare against independently composed drawing commands so omitting the
// selected outline or drawing it in the unselected color is a regression.
func TestRatingSelectedStroke(t *testing.T) {
	for _, dark := range []bool{false, true} {
		value := 1
		got := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			Rating(c, &value, 1)
		}, 100, 100)
		want := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			t := c.Theme()
			ui.Box(c).Size(t.Rem(1.5), t.Rem(1.5)).Draw(func(p *ui.Painter, r ui.Rect) {
				path := ratingStar(r)
				p.FillPath(path, t.Warning)
				p.StrokePath(path, 1.2, t.Warning)
			})
		}, 100, 100)
		got.Move(90, 90)
		want.Move(90, 90)
		if !bytes.Equal(got.Image().Pix, want.Image().Pix) {
			t.Fatal("selected stroke must match fill")
		}
	}
}
