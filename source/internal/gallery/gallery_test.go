package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestShellContent(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.View, 360, 480)
	if !tt.HasText("Shell") || !tt.HasText("Clicks: 0") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Count"); err != nil {
		t.Fatal(err)
	}
	if g.Clicks != 1 || !tt.HasText("Clicks: 1") {
		t.Fatal("caller-owned state was not updated")
	}
}
