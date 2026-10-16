package chat

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

func TestFloatingConversationClearanceAndJump(t *testing.T) {
	var layout FloatingConversationState
	var state ConversationState
	var latest, composer, viewport ui.Rect
	draft := "Reply"
	tt := ui.NewTester(func(c *ui.Context) {
		view := FloatingConversation(c, &layout, &state, "Latest", func(c *ui.Context) {
			for i := range 20 {
				row := components.Text(c, fmt.Sprintf("Message %d", i)).Height(40).Shrink(0)
				if i == 19 {
					latest = row.Bounds()
				}
			}
		}, func(c *ui.Context) {
			result := GlassComposer(c, &layout, &draft, ComposerOptions{Label: "Reply", SendLabel: "Send"})
			composer = result.Root.Bounds()
		})
		viewport = view.Bounds()
	}, 900, 600)
	tt.Frame()
	if !state.AtLatest() || latest.Y+latest.H > composer.Y {
		t.Fatalf("latest obscured: scroll=%+v latest=%+v composer=%+v", state.Scroll, latest, composer)
	}
	// The viewport extends behind the floating card, not just to its top.
	if viewport.H < composer.Y+composer.H-2 {
		t.Fatalf("viewport clipped above composer: %+v, %+v", state.Scroll, composer)
	}
	state.Scroll.Y = 0
	tt.Frame()
	jump, ok := tt.Find("Latest")
	if !ok || jump.Y+jump.H > composer.Y {
		t.Fatalf("jump hidden by composer: %+v, %+v", jump, composer)
	}
	if err := tt.Click("Latest"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !state.AtLatest() {
		t.Fatal("jump did not resume following")
	}
	oldHeight := layout.height
	draft = strings.Repeat("More lines\n", 6)
	tt.Frame()
	tt.Frame()
	if layout.height <= oldHeight || latest.Y+latest.H > composer.Y {
		t.Fatalf("multiline clearance did not update: height=%v old=%v latest=%+v composer=%+v", layout.height, oldHeight, latest, composer)
	}
}

func TestFloatingConversationThemedGlassJump(t *testing.T) {
	for _, palette := range []*ui.Theme{theme.LightPalette(), theme.DarkPalette()} {
		t.Run(fmt.Sprintf("dark=%v", palette.Dark), func(t *testing.T) {
			var layout FloatingConversationState
			var state ConversationState
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(palette)
				FloatingConversation(c, &layout, &state, "Latest", func(c *ui.Context) {
					ui.Box(c).Height(2000).FillWidth().Shrink(0).Background(ui.RGB(0, 0, 0))
				}, nil)
			}, 900, 600)
			tt.Frame()
			state.Scroll.Y = 0
			tt.Frame()
			jump, ok := tt.Find("Latest")
			if !ok {
				t.Fatal("missing jump control")
			}
			for _, hovered := range []bool{false, true} {
				if hovered {
					tt.Move(jump.X+jump.W/2, jump.Y+jump.H/2)
				}
				// Sample clear of the rim and arrow. The translucent glass
				// should stay close to the theme background over the transcript.
				pixel := tt.Image().RGBAAt(int(jump.X+7), int(jump.Y+jump.H/2))
				bg := palette.Background
				if math.Abs(float64(pixel.R)-float64(bg.R)) > 60 ||
					math.Abs(float64(pixel.G)-float64(bg.G)) > 60 ||
					math.Abs(float64(pixel.B)-float64(bg.B)) > 60 {
					t.Fatalf("hovered=%v: jump surface %v does not follow background %v", hovered, pixel, bg)
				}
			}
			if err := tt.Click("Latest"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if !state.AtLatest() {
				t.Fatal("glass jump did not resume following")
			}
		})
	}
}

func TestFloatingConversationFooterFade(t *testing.T) {
	for _, theme := range []*ui.Theme{ui.LightTheme(), ui.DarkTheme()} {
		var layout FloatingConversationState
		var state ConversationState
		var card ui.Rect
		draft := "Reply"
		backdrop := ui.RGB(240, 30, 70)
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme)
			FloatingConversation(c, &layout, &state, "Latest", func(c *ui.Context) {
				ui.Box(c).Width(900).Height(2000).Margin(0, -90, 0, -90).Shrink(0).Background(backdrop)
			}, func(c *ui.Context) {
				result := GlassComposer(c, &layout, &draft, ComposerOptions{
					Label: "Reply", Footer: func(c *ui.Context) { components.Text(c, "Project / branch") },
				})
				card = result.Card.Bounds()
			})
		}, 900, 600)
		for _, text := range []string{"Reply", strings.Repeat("More lines\n", 6)} {
			draft = text
			state.Scroll.Y = 200
			tt.Frame()
			tt.Frame()
			img := tt.Image()
			// Sample beside the card, where only the full-width fade overlays
			// the transcript; the glass itself does not affect these pixels.
			x := 30
			above := img.RGBAAt(x, int(card.Y+card.H/2)-6)
			middle := img.RGBAAt(x, int(card.Y+card.H*.75))
			footer := img.RGBAAt(x, int(card.Y+card.H)+14)
			bg := theme.Background
			want := color.RGBA{R: bg.R, G: bg.G, B: bg.B, A: 255}
			if above.R < 230 || middle == above || middle == want || footer != want {
				t.Fatalf("incorrect fade for card %+v: above=%v middle=%v footer=%v want=%v", card, above, middle, footer, want)
			}
		}
	}
}

func TestGlassComposerPreservesThemeAndRootBackground(t *testing.T) {
	for _, theme := range []*ui.Theme{ui.LightTheme(), ui.DarkTheme()} {
		draft := ""
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme)
			c.Root().Background(ui.Transparent)
			GlassComposer(c, nil, &draft, ComposerOptions{
				Attachments: func(c *ui.Context) {
					if c.Theme() != theme {
						t.Error("attachments theme changed")
					}
				},
				Accessory: func(c *ui.Context) {
					if c.Theme().Surface != theme.Surface {
						t.Error("accessory theme changed")
					}
				},
				Footer: func(c *ui.Context) {
					if c.Theme().Surface != theme.Surface {
						t.Error("footer theme changed")
					}
				},
			})
			if c.Theme() != theme {
				t.Error("composer changed frame theme")
			}
		}, 900, 400)
		tt.Frame()
		// Below the composer, the host's transparent root must remain exposed.
		if got := tt.Image().RGBAAt(450, 390); got.A != 0 {
			t.Fatalf("composer made the window root opaque: %v", got)
		}
	}
}
