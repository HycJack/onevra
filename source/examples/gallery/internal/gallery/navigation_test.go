package gallery

import (
	"bytes"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestNavigationSpacing(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 480)
	first, ok := tt.Find("Overview")
	if !ok || first.Y != 48 {
		t.Fatalf("first group must start at the shell's safe content edge: %v", first)
	}
	item, _ := tt.Find("Component status")
	next, _ := tt.Find("Icon library")
	if gap := next.Y - item.Y - item.H; math.Abs(float64(gap-12)) > 0.001 {
		t.Fatalf("group gap = %v, want 12", gap)
	}
}

func TestNavigationGroupOrder(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	var bottom float32
	for _, label := range []string{"Overview", "Component status", "Icon library", "Lucide", "Plugins", "Terminal", "Actions", "Buttons", "Selection & values", "Text input", "Overlays"} {
		r, ok := tt.Find(label)
		if !ok || r.Y < bottom {
			t.Fatalf("%s missing or out of order: %v (previous bottom %v)", label, r, bottom)
		}
		bottom = r.Y + r.H
	}
	if tt.HasText("Onevra controls") || tt.HasText("Components") || tt.HasText("MyGo components") {
		t.Fatal("obsolete component group is still visible")
	}
	for _, group := range mygoGroups {
		if group == "Plugins" {
			t.Fatal("plugins must not also appear among built-in MyGo groups")
		}
	}
}

func TestNavigationMyGoMarkers(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	texts := tt.Texts()
	markers := 0
	for i, text := range texts {
		if text != "•" {
			continue
		}
		markers++
		if i == 0 {
			t.Fatal("marker has no preceding item")
		}
		found := false
		for _, example := range mygoCatalog {
			found = found || texts[i-1] == example.name
		}
		if !found {
			t.Fatalf("marker follows non-MyGo item: %q", texts[i-1])
		}
	}
	if markers != len(mygoCatalog) {
		t.Fatalf("got %d MyGo markers, want %d", markers, len(mygoCatalog))
	}
}

func TestNavigationSelectionAndHover(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(map[bool]string{false: "light", true: "dark"}[dark], func(t *testing.T) {
			g := &Gallery{}
			tt := ui.NewTester(g.Left, 220, 480)
			tt.SetDark(dark)
			for _, label := range []string{"Lucide", "Buttons"} {
				tt.Move(219, 479)
				before := bytes.Clone(tt.Image().Pix)
				r, _ := tt.Find(label)
				tt.Move(r.X+r.W/2, r.Y+r.H/2)
				if !bytes.Equal(before, tt.Image().Pix) {
					t.Fatalf("%s changed appearance on hover", label)
				}
			}
			for _, label := range []string{"Buttons", "Lucide"} {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				r, _ := tt.Find(label)
				// Sample the face away from the icon, text and rounded corners.
				pixel := tt.Image().RGBAAt(int(r.X+r.W-12), int(r.Y+r.H/2))
				if pixel.A == 0 || absDiff(pixel.R, pixel.G) > 4 || absDiff(pixel.G, pixel.B) > 4 {
					t.Fatalf("%s selection is not neutral gray: %v", label, pixel)
				}
			}
		})
	}
}

func absDiff(a, b uint8) int {
	return max(int(a)-int(b), int(b)-int(a))
}
