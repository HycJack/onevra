package components

import (
	"math"
	"reflect"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestGridPositions(t *testing.T) {
	for _, tc := range []struct {
		name           string
		length, offset float32
		step           float64
		want           []float32
	}{
		{"origin", 40, 0, 16, []float32{16, 32}},
		{"positive", 40, 5, 16, []float32{5, 21, 37}},
		{"negative", 40, -5, 16, []float32{11, 27}},
		{"whole period", 40, -32, 16, []float32{16, 32}},
		{"major phase", 80, -17, 64, []float32{47}},
		{"large pan", 40, 16000005, 16, []float32{5, 21, 37}},
		{"fractional", 12, 0.5, 4.5, []float32{0.5, 5, 9.5}},
		{"invalid offset", 40, float32(math.NaN()), 16, []float32{16, 32}},
		{"empty", 0, 0, 16, nil},
		{"invalid length", float32(math.Inf(1)), 0, 16, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []float32
			gridPositions(tc.length, tc.step, tc.offset, func(x float32) { got = append(got, x) })
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPixelGridRendering(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, spacing := range []float32{0, -1, 16, float32(math.NaN()), float32(math.Inf(1))} {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				PixelGrid(c, PixelGridOptions{Spacing: spacing, MajorEvery: 4}).Width(80).Height(80)
			}, 100, 100)
			img := tt.Image()
			background := img.RGBAAt(8, 8)
			if img.RGBAAt(0, 8) != background || img.RGBAAt(8, 0) != background || img.RGBAAt(0, 0) != background {
				t.Fatalf("mode=%v spacing=%v: unexpected leading border", mode, spacing)
			}
			minor := img.RGBAAt(16, 8)
			major := img.RGBAAt(64, 8)
			if background == minor || minor == major || major == background {
				t.Fatalf("mode=%v spacing=%v: missing grid contrast", mode, spacing)
			}
		}
	}
}

func TestPixelGridClippingAndChildren(t *testing.T) {
	clicked := false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Background(ui.RGB(255, 0, 0)).Padding(10).Children(func() {
			PixelGrid(c, PixelGridOptions{Spacing: 0.1, MajorEvery: -1, OffsetX: -1}).Width(80).Height(80).Children(func() {
				if Button(c, "Content", ButtonPrimary).Clicked() {
					clicked = true
				}
			})
		})
	}, 120, 120)
	img := tt.Image()
	if img.RGBAAt(9, 70) != img.RGBAAt(90, 70) || img.RGBAAt(90, 70) != img.RGBAAt(110, 70) {
		t.Fatal("grid paints outside its bounds")
	}
	if img.RGBAAt(13, 70) == img.RGBAAt(14, 70) {
		t.Fatal("minimum spacing or negative offset not applied")
	}
	if img.RGBAAt(13, 70) != img.RGBAAt(17, 70) {
		t.Fatal("major emphasis was not disabled")
	}
	if err := tt.Click("Content"); err != nil {
		t.Fatal(err)
	}
	if !clicked {
		t.Fatal("grid blocked child interaction")
	}
}
