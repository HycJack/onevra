package components

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestTextMatchesPrimitive(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{140, 480} {
			t.Run(fmt.Sprintf("%v/%d", mode, width), func(t *testing.T) {
				build := func(text func(*ui.Context, string) ui.Element) *ui.Tester {
					return ui.NewTester(func(c *ui.Context) {
						theme.Apply(c, mode)
						ui.Column(c).Fill().TextColor(c.Theme().TextMuted).FontSize(c.Theme().Rem(1.2)).Children(func() {
							text(c, "Hello, world. 你好，世界。こんにちは 🌍").Label("wrapped")
							text(c, "A long single-line label with inherited styles").SingleLine().Label("single")
							text(c, "Heading").Bold().FontSize(c.Theme().Rem(1.5)).TextColor(c.Theme().Text)
							text(c, "")
							ui.RichText(c).Children(func() {
								text(c, "Inline ")
								text(c, "bold").Bold()
							})
						})
					}, width, 400)
				}
				got, want := build(Text), build(ui.Text)
				if !bytes.Equal(got.Image().Pix, want.Image().Pix) {
					t.Fatal("Text changed primitive rendering or style inheritance")
				}
				for _, label := range []string{"wrapped", "single"} {
					g, gok := got.Find(label)
					w, wok := want.Find(label)
					if !gok || !wok || g != w {
						t.Fatalf("%s: bounds %v, want %v", label, g, w)
					}
				}
				if !got.HasText("Inline bold") {
					t.Fatal("Text does not compose in RichText")
				}
			})
		}
	}
}

func TestTextf(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			Textf(c, "%s: %d (%.1f%%)", "进度", 3, 12.5).Bold()
			Textf(c, "No arguments").SingleLine()
		})
	}, 300, 100)
	for _, want := range []string{"进度: 3 (12.5%)", "No arguments"} {
		if !tt.HasText(want) {
			t.Fatalf("missing formatted text %q", want)
		}
	}
}
