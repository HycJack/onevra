package chat

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

func TestFloatingConversationScrollbarClearsComposer(t *testing.T) {
	for _, width := range []int{900, 360} {
		for _, dark := range []bool{false, true} {
			t.Run(fmt.Sprintf("width=%d/dark=%v", width, dark), func(t *testing.T) {
				var layout FloatingConversationState
				var state ConversationState
				draft, showError := "", false
				theme := ui.LightTheme()
				if dark {
					theme = ui.DarkTheme()
				}
				// Isolate the scrollbar pixels from the content and glass.
				theme.Scrollbar = ui.RGB(240, 30, 70)
				tt := ui.NewTester(func(c *ui.Context) {
					c.SetTheme(theme)
					FloatingConversation(c, &layout, &state, "Latest", func(c *ui.Context) {
						ui.Box(c).Height(1800).Shrink(0)
						components.Text(c, "Last message").Height(20)
					}, func(c *ui.Context) {
						ui.Column(c).Label("Composer").FillWidth().Gap(6).Children(func() {
							if showError {
								GlassPanel(c, func(c *ui.Context) { components.Text(c, "Retrying") }).Height(44)
							}
							GlassComposer(c, &layout, &draft, ComposerOptions{
								Label:  "Reply",
								Footer: func(c *ui.Context) { components.Text(c, "Project / branch") },
							})
						})
					}).Label("Viewport")
				}, width, 600)
				for _, tc := range []struct {
					draft string
					err   bool
				}{{"", false}, {strings.Repeat("More lines\n", 6), false}, {"", true}, {"", false}} {
					draft, showError = tc.draft, tc.err
					state.JumpToLatest()
					tt.Frame()
					tt.Frame()
					x := float32(width - 6)
					tt.Move(x, 40)
					composer, ok := tt.Find("Composer")
					if !ok {
						t.Fatal("missing composer")
					}
					bottom := -1
					img := tt.Image()
					for y := 0; y < 600; y++ {
						pixel := img.RGBAAt(int(x), y)
						if pixel.R == 240 && pixel.G == 30 && pixel.B == 70 {
							bottom = y
						}
					}
					trackBottom := composer.Y - 8
					if bottom < int(trackBottom)-4 || bottom >= int(trackBottom) {
						t.Fatalf("scrollbar bottom=%d does not leave an 8-pixel gap above composer=%+v", bottom, composer)
					}
					viewport, _ := tt.Find("Viewport")
					last, _ := tt.Find("Last message")
					if viewport.H != 600 || !state.AtLatest() || last.Y+last.H > composer.Y {
						t.Fatalf("scrollbar inset changed viewport or latest visibility: viewport=%+v last=%+v composer=%+v", viewport, last, composer)
					}
					// The cleared region is not an invisible paging target.
					state.Scroll.Y = 0
					tt.Frame()
					for _, y := range []float32{composer.Y - 4, 590} {
						tt.ClickAt(x, y)
						if state.Scroll.Y != 0 {
							t.Fatalf("click below the scrollbar at %v paged to %v", y, state.Scroll.Y)
						}
					}
					// The shortened track still reaches the end of the transcript.
					tt.Press(x, 8)
					tt.Move(x, trackBottom-2)
					tt.Release(x, trackBottom-2)
					if !state.AtLatest() {
						t.Fatalf("dragging to the track end did not reach latest: %+v", state.Scroll)
					}
				}
			})
		}
	}
}
