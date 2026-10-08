package gallery

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGlassPreview(t *testing.T) {
	for _, width := range []int{280, 540} {
		s := newDemoState()
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Fill().Gap(12).Children(func() { exampleGlass(c, s) })
		}, width, 850)
		for _, dark := range []bool{false, true} {
			s.value = 0
			tt.SetDark(dark)
			for _, label := range []string{"Regular", "Clear", "Tinted", "Background position"} {
				if _, ok := tt.Find(label); !ok {
					t.Fatalf("missing glass preview control: %s", label)
				}
			}
			before := bytes.Clone(tt.Image().Pix)
			previous := s.value
			if err := tt.Click("Background position"); err != nil {
				t.Fatal(err)
			}
			if s.value == previous || bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("background control did not update the preview")
			}
		}
	}
}
