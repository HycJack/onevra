package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestControlPopover(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	if err := tt.Click("Popover"); err != nil {
		t.Fatal(err)
	}
	if example, ok := controlExampleFor(g.page); !ok || example.name != "Popover" {
		t.Fatal("popover is not a Onevra page")
	}
	content := ui.NewTester(g.View, 540, 850)
	for _, label := range []string{"Show popover", "Enable notifications", "Close"} {
		if err := content.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	s := g.controlState("Popover", "Default")
	if s.dialogOpen || !s.on {
		t.Fatal("popover actions did not update state")
	}
	for _, e := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if e.Name == "Popover" && (e.Status != mygoReplaced || e.Replacement != "components.Popover") {
			t.Fatal("popover replacement is not recorded")
		}
	}
}
