package sidebar

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestLeftWidthPersistence(t *testing.T) {
	s := newShell(ShellOptions{LeftWidth: 220, LeftMinWidth: 120, LeftMaxWidth: 520})
	for _, tc := range []struct{ input, want float32 }{{315.5, 315.5}, {0, 120}, {900, 520}, {float32(math.NaN()), 520}, {float32(math.Inf(1)), 520}} {
		s.SetLeftWidth(tc.input)
		if got := s.LeftWidth(); got != tc.want {
			t.Fatalf("width=%v want=%v", got, tc.want)
		}
	}
	s.SetLeftWidth(300)
	s.SetLeftExpanded(false)
	s.SetPage(func(*ui.Context) {})
	s.SetPage(nil)
	if s.LeftWidth() != 300 || s.options.LeftWidth != 220 {
		t.Fatal("lost chosen width or changed reset default")
	}
}
