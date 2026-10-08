package components

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestDialogLayers(t *testing.T) {
	for _, protected := range []bool{false, true} {
		outer, inner := true, true
		tt := ui.NewTester(func(c *ui.Context) {
			Dialog(c, &outer, DialogOptions{Label: "Outer"}, func(p ui.Element) {
				Button(c, "Outer action", ButtonPrimary)
				Dialog(c, &inner, DialogOptions{Label: "Inner", DisableBackdropClose: protected}, func(p ui.Element) {
					Button(c, "Inner action", ButtonPrimary)
				})
			})
		}, 600, 500)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		tt.ClickAt(2, 2)
		if !outer || inner != protected {
			t.Fatalf("backdrop: outer=%v inner=%v protected=%v", outer, inner, protected)
		}
		if protected {
			tt.Key(0, ui.KeyEscape)
		}
		if inner || !outer {
			t.Fatal("top layer dismissal")
		}
		tt.Frame()
		tt.ClickAt(2, 2)
		if outer {
			t.Fatal("outer layer did not dismiss")
		}
	}
}

func TestAlertDialogResults(t *testing.T) {
	for _, action := range []string{"OK", "Cancel", "backdrop", "escape"} {
		open := true
		var results []DialogResult
		tt := ui.NewTester(func(c *ui.Context) {
			r := AlertDialog(c, &open, AlertDialogOptions{Title: "Question", Message: "Continue?", ConfirmLabel: "OK", CancelLabel: "Cancel"})
			if r != DialogNone {
				results = append(results, r)
			}
		}, 600, 500)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		switch action {
		case "backdrop":
			tt.ClickAt(2, 2)
		case "escape":
			tt.Key(0, ui.KeyEscape)
		default:
			if err := tt.Click(action); err != nil {
				t.Fatal(err)
			}
		}
		tt.Frame()
		want := DialogCancelled
		if action == "OK" {
			want = DialogConfirmed
		}
		if open || len(results) != 1 || results[0] != want {
			t.Fatalf("%s: open=%v results=%v", action, open, results)
		}
	}
}

func TestDialogGeometryAndNotice(t *testing.T) {
	open := true
	var bounds ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Dialog(c, &open, DialogOptions{DisableEscapeClose: true, DisableBackdropClose: true}, func(p ui.Element) {
			p.DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
			Text(c, "Body")
		})
	}, 280, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	if bounds.X < 0 || bounds.X+bounds.W > 280 || bounds.Y < 0 || bounds.Y+bounds.H/2 > 480*.36 || bounds.Y+bounds.H/2 < 480*.30 {
		t.Fatalf("bounds: %+v", bounds)
	}
	tt.Key(0, ui.KeyEscape)
	tt.ClickAt(2, 2)
	if !open {
		t.Fatal("protected dialog closed")
	}
	notice := ui.NewTester(func(c *ui.Context) {
		AlertDialog(c, &open, AlertDialogOptions{Title: "Notice", ConfirmLabel: "OK", CancelLabel: "Cancel", HideCancel: true})
	}, 400, 400)
	if notice.HasText("Cancel") || !notice.HasText("OK") {
		t.Fatal("notice actions")
	}
}
