package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

func TestToastExampleUndoAndUpload(t *testing.T) {
	g := &Gallery{}
	actions, progress := g.controlState("Toast", "Actions"), g.controlState("Toast", "Progress")
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.LightPalette())
		ui.Column(c).Gap(12).Children(func() {
			controlToast(c, actions, "Actions")
			controlToast(c, progress, "Progress")
		})
		components.Toasts(c)
	}, 540, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	click := func(label string) {
		t.Helper()
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	click("Delete")
	if !tt.HasText("Items: 2") || !tt.HasText("Item deleted") {
		t.Fatalf("delete: %q", tt.Texts())
	}
	click("Undo")
	if !tt.HasText("Items: 3") {
		t.Fatalf("undo: %q", tt.Texts())
	}
	click("Start upload")
	if !tt.HasText("Uploading report.pdf") {
		t.Fatalf("upload: %q", tt.Texts())
	}
	// Finishing updates the same toast rather than stacking another.
	click("Finish upload")
	if tt.HasText("Uploading report.pdf") || !tt.HasText("Upload complete") {
		t.Fatalf("finish: %q", tt.Texts())
	}
}

func TestToastExamplePlacement(t *testing.T) {
	g := &Gallery{page: pageFromKey("controls/Toast")}
	tt := ui.NewTester(func(c *ui.Context) { c.SetTheme(theme.LightPalette()); g.View(c) }, 900, 900)
	if !tt.HasText("Placement") {
		t.Fatalf("missing placement section: %q", tt.Texts())
	}
	if err := tt.Click("Top right"); err != nil {
		t.Fatal(err)
	}
	if g.toasts.Placement != components.ToastTopEnd {
		t.Fatalf("placement = %v", g.toasts.Placement)
	}
}
