package ui

import (
	"bytes"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go/theme"
)

func TestButtonVariantsClick(t *testing.T) {
	variants := []ButtonVariant{ButtonPrimary, ButtonDefault, ButtonOutline, ButtonGhost, ButtonLink}
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, variant := range variants {
			clicked := false
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				if Button(c, "OK", variant).Clicked() {
					clicked = true
				}
			}, 300, 100)
			if err := tt.Click("OK"); err != nil {
				t.Fatal(err)
			}
			if !clicked {
				t.Errorf("mode=%v variant=%d: click not reported", mode, variant)
			}
		}
	}
}

func TestButtonKeyboard(t *testing.T) {
	clicked := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		if Button(c, "Save", ButtonPrimary).Clicked() {
			clicked = true
		}
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !clicked {
		t.Error("Space did not click the focused button")
	}
}

func TestButtonDisabledDims(t *testing.T) {
	sample := func(disabled bool) (r, g, b uint8) {
		var rect ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, theme.Light)
			ui.Column(c).Disabled(disabled).Fill().Children(func() {
				Button(c, "OK", ButtonPrimary)
			})
		}, 300, 100)
		var ok bool
		rect, ok = tt.Find("OK")
		if !ok {
			t.Fatal("button not found")
		}
		// Inside the left padding, clear of the label glyphs.
		p := tt.Image().RGBAAt(int(rect.X+6), int(rect.Y+rect.H/2))
		return p.R, p.G, p.B
	}
	er, eg, eb := sample(false)
	dr, dg, db := sample(true)
	if er == dr && eg == dg && eb == db {
		t.Fatal("disabled button renders identically to enabled")
	}
	// Dimming blends the accent toward the light background.
	if dr <= er || dg <= eg || db <= eb {
		t.Errorf("disabled pixel (%d,%d,%d) not lighter than enabled (%d,%d,%d)", dr, dg, db, er, eg, eb)
	}
}

func TestButtonLinkInline(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, disabled := range []bool{false, true} {
			clicks := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Fill().AlignItems(ui.Start).Disabled(disabled).Children(func() {
					ui.Text(c, "Details").Label("plain")
					ui.RichText(c).Children(func() {
						ui.Text(c, "Read ")
						if Button(c, "Details", ButtonLink).Clicked() {
							clicks++
						}
						ui.Text(c, " for more.")
					})
				})
			}, 300, 120)
			link, ok := tt.Find("Details")
			if !ok {
				t.Fatal("missing inline link")
			}
			plain, ok := tt.Find("plain")
			if !ok || link.H != plain.H || link.W != plain.W {
				t.Fatalf("link %+v must match plain text %+v", link, plain)
			}
			if !tt.HasText("Read Details for more.") || link.X <= 0 {
				t.Fatal("link is not composed inline")
			}
			tt.Move(290, 110)
			rest := slices.Clone(tt.Image().Pix)
			tt.Move(link.X+link.W/2, link.Y+link.H/2)
			hover := slices.Clone(tt.Image().Pix)
			if bytes.Equal(rest, hover) != disabled {
				t.Fatalf("disabled=%v: incorrect underline feedback", disabled)
			}
			tt.Move(290, 110)
			if !bytes.Equal(rest, tt.Image().Pix) {
				t.Fatal("underline remains after pointer exit")
			}
			if err := tt.Click("Details"); err != nil {
				t.Fatal(err)
			}
			want := 1
			if disabled {
				want = 0
			}
			if clicks != want || len(tt.OpenedURLs()) != 0 {
				t.Fatal("link must report clicks without opening a URL")
			}
		}
	}
}

func TestButtonLabel(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light)
		Button(c, "Delete", ButtonDefault)
	}, 300, 100)
	if !tt.HasText("Delete") {
		t.Error("button label not shown")
	}
}
