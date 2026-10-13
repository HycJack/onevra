package components

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestToastsActionAndClose(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		var undone, closed int
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			if Button(c, "Delete", ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{ID: "deleted", Title: "Note deleted", Description: "Moved to trash", Type: ToastSuccess,
					Timeout: -1, Action: "Undo", OnAction: func() { undone++ }, OnClose: func() { closed++ }})
			}
			Toasts(c)
		}, 480, 360)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		if err := tt.Click("Delete"); err != nil {
			t.Fatal(err)
		}
		// Onevra's look replaces MyGo's, which would show the title again.
		if n := len(slices.DeleteFunc(tt.Texts(), func(s string) bool { return s != "Note deleted" })); n != 1 {
			t.Fatalf("title shown %d times: %q", n, tt.Texts())
		}
		if !tt.HasText("Moved to trash") {
			t.Fatalf("missing description: %q", tt.Texts())
		}
		r, ok := tt.Find("Note deleted")
		if !ok || r.Y < 280 || r.X < 60 || r.X+r.W > 420 {
			t.Fatalf("toast is not at the bottom center: %+v", r)
		}
		if err := tt.Click("Undo"); err != nil {
			t.Fatal(err)
		}
		if undone != 1 || closed != 1 {
			t.Fatalf("undo: undone %d, closed %d", undone, closed)
		}
		if err := tt.Click("Delete"); err != nil {
			t.Fatal(err)
		}
		if err := tt.Click("Close"); err != nil {
			t.Fatal(err)
		}
		if undone != 1 || closed != 2 {
			t.Fatalf("close: undone %d, closed %d", undone, closed)
		}
	}
}

func TestToastsLocalizedCloseAndLeadingSlot(t *testing.T) {
	leading := false
	var closed int
	tt := ui.NewTester(func(c *ui.Context) {
		if Button(c, "Upload", ButtonSecondary).Clicked() {
			c.AddToast(ui.Toast{Title: "Uploading", Timeout: -1, OnClose: func() { closed++ },
				Data: ToastContent{Leading: func() {
					leading = true
					Spinner(c, SpinnerRing)
				}}})
		}
		WithLocalization(c, Localization{CloseToast: "关闭"}, func() { Toasts(c) })
	}, 480, 360)
	if err := tt.Click("Upload"); err != nil {
		t.Fatal(err)
	}
	if !leading || !tt.HasText("Uploading") {
		t.Fatalf("leading slot %v, texts %q", leading, tt.Texts())
	}
	if err := tt.Click("关闭"); err != nil {
		t.Fatal(err)
	}
	if closed != 1 {
		t.Fatalf("localized close did not close: %d", closed)
	}
}

// A view's own toast viewport, built before Toasts, keeps its look.
func TestToastsYieldToEarlierViewport(t *testing.T) {
	custom := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Button(c, "Show", ButtonSecondary).Clicked() {
			c.Toast("Saved")
		}
		ui.ToastViewportBase(c, func(_ ui.Element, toasts []ui.Toast) {
			for _, toast := range toasts {
				custom++
				ui.ToastBase(c, toast).Root.Children(func() { ui.Text(c, toast.Title) })
			}
		})
		Toasts(c)
	}, 480, 360)
	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	if custom == 0 || tt.HasText("Close") {
		t.Fatalf("custom viewport %d, texts %q", custom, tt.Texts())
	}
	if _, ok := tt.Find("Close"); ok {
		t.Fatal("Onevra close button built over the custom viewport")
	}
}

func TestToastsPlacement(t *testing.T) {
	for _, tc := range []struct {
		placement ToastPlacement
		top       bool
		align     ui.Align
	}{
		{ToastBottom, false, ui.Center}, {ToastBottomStart, false, ui.Start}, {ToastBottomEnd, false, ui.End},
		{ToastTop, true, ui.Center}, {ToastTopStart, true, ui.Start}, {ToastTopEnd, true, ui.End},
	} {
		added := false
		tt := ui.NewTester(func(c *ui.Context) {
			if !added {
				added = true
				c.AddToast(ui.Toast{Title: "Older", Timeout: -1})
				c.AddToast(ui.Toast{Title: "Newer", Timeout: -1})
			}
			Toasts(c, ToastOptions{Placement: tc.placement, Inset: 40})
		}, 800, 600)
		older, ok1 := tt.Find("Older")
		newer, ok2 := tt.Find("Newer")
		if !ok1 || !ok2 {
			t.Fatalf("%v: missing toasts %q", tc.placement, tt.Texts())
		}
		if tc.top {
			// Top toasts clear the inset and stack downward, newest first.
			if newer.Y < 40 || newer.Y > 120 || older.Y <= newer.Y {
				t.Fatalf("%v: newer %+v, older %+v", tc.placement, newer, older)
			}
		} else if older.Y > newer.Y || newer.Y > 600-40 || newer.Y < 400 {
			t.Fatalf("%v: newer %+v, older %+v", tc.placement, newer, older)
		}
		switch x := newer.X; tc.align {
		case ui.Start:
			if x > 100 {
				t.Fatalf("%v: not at the start: %v", tc.placement, x)
			}
		case ui.End:
			if x < 400 {
				t.Fatalf("%v: not at the end: %v", tc.placement, x)
			}
		default:
			if x < 200 || x > 400 {
				t.Fatalf("%v: not centered: %v", tc.placement, x)
			}
		}
	}
}
