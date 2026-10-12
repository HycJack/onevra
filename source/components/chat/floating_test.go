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

func TestFloatingConversationHeaderClearance(t *testing.T) {
	var layout FloatingConversationState
	var state ConversationState
	var first, header, composer ui.Rect
	show, overlay := true, false
	tt := ui.NewTester(func(c *ui.Context) {
		options := FloatingConversationOptions{JumpLabel: "Latest", HeaderOverlay: overlay, Content: func(c *ui.Context) {
			first = components.Text(c, "First").Height(40).Shrink(0).Bounds()
		}, Composer: func(c *ui.Context) {
			composer = ui.Box(c).Height(60).FillWidth().Bounds()
		}}
		if show {
			options.Header = func(c *ui.Context) {
				header = GlassPanel(c, func(c *ui.Context) { components.Text(c, "Output") }).Bounds()
			}
		}
		FloatingConversationWith(c, &layout, &state, options)
	}, 900, 600)
	tt.Frame()
	tt.Frame()
	if header.H <= 0 || first.Y < header.Y+header.H {
		t.Fatalf("first turn under header: first=%+v header=%+v", first, header)
	}
	// The header shares the composer's width.
	if header.W != composer.W || header.X != composer.X {
		t.Fatalf("header width differs: header=%+v composer=%+v", header, composer)
	}
	top := first.Y
	show = false
	tt.Frame()
	tt.Frame()
	if layout.headerHeight != 0 || first.Y >= top {
		t.Fatalf("header space kept: height=%v first=%+v", layout.headerHeight, first)
	}
	withoutHeader := first
	show = true
	tt.Frame()
	tt.Frame()
	if first.Y != top {
		t.Fatalf("header clearance not restored: first=%+v top=%v", first, top)
	}
	overlay = true
	tt.Frame()
	tt.Frame()
	if layout.headerHeight != 0 || first != withoutHeader {
		t.Fatalf("overlay kept stale clearance: height=%v first=%+v want=%+v", layout.headerHeight, first, withoutHeader)
	}
}

func TestFloatingConversationHeaderTopInset(t *testing.T) {
	for _, overlay := range []bool{false, true} {
		t.Run(fmt.Sprintf("overlay=%v", overlay), func(t *testing.T) {
			var layout FloatingConversationState
			var state ConversationState
			var inset *float32
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Padding(30, 0, 0, 0).Children(func() {
					FloatingConversationWith(c, &layout, &state, FloatingConversationOptions{
						HeaderOverlay: overlay, HeaderTopInset: inset,
						Header:   func(c *ui.Context) { GlassPanel(c, nil).Label("Header").Height(32) },
						Content:  func(c *ui.Context) { components.Text(c, "First").Height(40) },
						Composer: func(c *ui.Context) { ui.Box(c).Label("Composer").Height(60).FillWidth() },
					}).Label("Viewport")
				})
			}, 900, 600)
			bounds := func(label string) ui.Rect {
				t.Helper()
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing %s", label)
				}
				return r
			}
			zero, custom, negative := float32(0), float32(24), float32(-8)
			for _, tc := range []struct {
				inset *float32
				want  float32
			}{{nil, 8}, {&zero, 0}, {&custom, 24}, {&negative, 0}, {nil, 8}} {
				inset = tc.inset
				tt.Frame()
				tt.Frame()
				header, viewport, composer := bounds("Header"), bounds("Viewport"), bounds("Composer")
				if header.Y != viewport.Y+tc.want || header.X != composer.X || header.W != composer.W {
					t.Fatalf("inset %v: header=%+v viewport=%+v composer=%+v", tc.want, header, viewport, composer)
				}
				clearance := float32(0)
				if !overlay {
					clearance = header.H + tc.want
				}
				first := bounds("First")
				if layout.headerHeight != clearance || first.Y != viewport.Y+22+clearance {
					t.Fatalf("inset %v: clearance=%v want=%v first=%+v", tc.want, layout.headerHeight, clearance, first)
				}
			}
		})
	}
}

func TestFloatingConversationHeaderOverlayKeepsTranscriptStill(t *testing.T) {
	for _, width := range []int{900, 360} {
		for _, position := range []string{"short", "latest", "reading"} {
			t.Run(fmt.Sprintf("width=%d/%s", width, position), func(t *testing.T) {
				var layout FloatingConversationState
				var state ConversationState
				var headerHeight float32
				tt := ui.NewTester(func(c *ui.Context) {
					options := FloatingConversationOptions{JumpLabel: "Latest", HeaderOverlay: true,
						Content: func(c *ui.Context) {
							components.Text(c, "First").Height(40).Shrink(0)
							if position != "short" {
								ui.Box(c).Height(1500).Shrink(0)
							}
							components.Text(c, "Last").Height(40).Shrink(0)
						}, Composer: func(c *ui.Context) {
							ui.Box(c).Label("Composer").Height(60).FillWidth()
						}}
					if headerHeight > 0 {
						options.Header = func(c *ui.Context) {
							GlassPanel(c, nil).Label("Header").Height(headerHeight)
						}
					}
					FloatingConversationWith(c, &layout, &state, options).Label("Viewport")
				}, width, 600)
				tt.Frame()
				if position == "reading" {
					state.Scroll.Y = 200
					tt.Frame()
				}
				bounds := func(label string) ui.Rect {
					t.Helper()
					r, ok := tt.Find(label)
					if !ok {
						t.Fatalf("missing %s", label)
					}
					return r
				}
				before := map[string]ui.Rect{}
				for _, label := range []string{"First", "Last", "Composer", "Viewport"} {
					before[label] = bounds(label)
				}
				y, maxY := state.Scroll.Y, state.Scroll.MaxY
				for _, height := range []float32{32, 108, 48, 0, 80, 0} {
					headerHeight = height
					for range 2 {
						tt.Frame()
						for label, want := range before {
							if got := bounds(label); got != want {
								t.Fatalf("header height %v moved %s: got=%+v want=%+v", height, label, got, want)
							}
						}
						if layout.headerHeight != 0 || state.Scroll.Y != y || state.Scroll.MaxY != maxY {
							t.Fatalf("overlay changed clearance or scroll: height=%v scroll=%+v want=%v/%v", layout.headerHeight, state.Scroll, y, maxY)
						}
					}
					if height == 0 {
						if _, ok := tt.Find("Header"); ok {
							t.Fatal("hidden header still present")
						}
						continue
					}
					header := bounds("Header")
					composer := before["Composer"]
					if header.Y != 8 || header.H != height || header.X != composer.X || header.W != composer.W {
						t.Fatalf("incorrect header geometry: header=%+v composer=%+v height=%v", header, composer, height)
					}
					if position == "short" && header.Y+header.H <= before["First"].Y {
						t.Fatal("header does not overlap the first message")
					}
				}
			})
		}
	}
}

func TestFloatingConversationTopInset(t *testing.T) {
	for _, header := range []string{"none", "overlay", "reserved"} {
		t.Run(header, func(t *testing.T) {
			var layout FloatingConversationState
			var state ConversationState
			theme := ui.LightTheme()
			// Isolate the scrollbar pixels from the content and glass.
			theme.Scrollbar = ui.RGB(240, 30, 70)
			options := FloatingConversationOptions{
				TopInset: 48, HeaderOverlay: header == "overlay", HeaderTopInset: new(float32),
				Content: func(c *ui.Context) {
					components.Text(c, "First").Height(40)
					ui.Box(c).Height(1800).Shrink(0)
				},
				Composer: func(c *ui.Context) { ui.Box(c).Label("Composer").Height(60).FillWidth() },
			}
			if header != "none" {
				options.Header = func(c *ui.Context) { GlassPanel(c, nil).Label("Header").Height(32) }
			}
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme)
				FloatingConversationWith(c, &layout, &state, options).Label("Viewport")
			}, 900, 600)
			tt.Frame()
			tt.Frame()
			// Read from the top rather than following the latest turn.
			state.Scroll.Y = 0
			tt.Frame()
			x := float32(900 - 6)
			tt.Move(x, 300)
			viewport, _ := tt.Find("Viewport")
			first, _ := tt.Find("First")
			want := viewport.Y + 22 + 48
			if header == "reserved" {
				want += 32
			}
			if viewport.Y != 0 || viewport.H != 600 || first.Y != want {
				t.Fatalf("viewport=%+v first=%+v want first at %v", viewport, first, want)
			}
			if header != "none" {
				if r, _ := tt.Find("Header"); r.Y != 48 {
					t.Fatalf("header=%+v does not start below the inset", r)
				}
			}
			top := -1
			img := tt.Image()
			for y := 0; y < 600; y++ {
				pixel := img.RGBAAt(int(x), y)
				if pixel.R == 240 && pixel.G == 30 && pixel.B == 70 {
					top = y
					break
				}
			}
			if top < 48 || top > 52 {
				t.Fatalf("scrollbar top=%d does not start below the inset", top)
			}
			// Scrolled content passes beneath the inset rather than being clipped.
			state.Scroll.Y = 100
			tt.Frame()
			if first, _ = tt.Find("First"); first.Y >= 48 {
				t.Fatalf("scrolled first=%+v did not move beneath the inset", first)
			}
		})
	}
}
