package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMotionDraggableSnapsAndResets(t *testing.T) {
	_, s, tt := featureTester("Motion", "Draggable", "en", 540)
	// Rest at once so the snap is observable without a clock.
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	card, ok := tt.Find("Throw to either end")
	if !ok {
		t.Fatal("missing draggable card")
	}
	x, y := card.X+card.W/2, card.Y+card.H/2
	tt.Move(x, y)
	tt.Press(x, y)
	for dx := float32(10); dx <= 80; dx += 10 {
		tt.Move(x+dx, y+20)
	}
	if s.offset != (ui.Point{X: 80}) {
		t.Fatalf("dragged by 80 along the horizontal axis: offset %v", s.offset)
	}
	tt.Release(x+80, y+20)
	tt.Frame()
	if s.offset != (ui.Point{X: 120}) {
		t.Fatalf("let go past the middle: offset %v, want the far end", s.offset)
	}
	if err := tt.Click("Reset position"); err != nil {
		t.Fatal(err)
	}
	if s.offset != (ui.Point{}) {
		t.Fatalf("reset: offset %v", s.offset)
	}
}
