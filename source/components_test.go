package oneogo_test

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go"
)

func TestCompositionAndCallerOwnedAction(t *testing.T) {
	clicks := 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		oneogo.Section(c, "Section", "", func() {
			oneogo.EmptyState(c, oneogo.EmptyStateOptions{Title: "Empty", Description: "Caller content", Actions: func() {
				if ui.Button(c, "Create").Disabled(disabled).Clicked() {
					clicks++
				}
			}})
			oneogo.Notice(c, oneogo.NoticeOptions{Title: "Status"})
		})
	}, 400, 600)
	if !tt.HasText("Caller content") || !tt.HasText("Status") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Create"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatalf("clicks = %d", clicks)
	}
	disabled = true
	tt.Frame()
	if err := tt.Click("Create"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatal("disabled action activated")
	}
}

func TestThemesScalingAndWrapping(t *testing.T) {
	for _, dark := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			original := c.Theme()
			ui.Column(c).FillWidth().Children(func() {
				oneogo.Section(c, "Heading", "", nil)
				oneogo.EmptyState(c, oneogo.EmptyStateOptions{Title: "Empty"})
				for kind := oneogo.NoticeInfo; kind <= oneogo.NoticeSuccess; kind++ {
					oneogo.Notice(c, oneogo.NoticeOptions{Kind: kind, Title: "Status", Description: "A long description wraps across multiple lines in a narrow window. 中英文内容保持可读。"}).Key(kind)
				}
			})
			if c.Theme() != original {
				t.Error("component replaced caller theme")
			}
		}, 360, 1200)
		tt.SetDark(dark)
		before, ok := tt.Find("Heading")
		if !ok {
			t.Fatal("missing heading")
		}
		tt.SetPreferences(ui.Preferences{TextScale: 1.5, HighContrast: true, ReduceMotion: true})
		after, ok := tt.Find("Heading")
		if !ok || after.H <= before.H {
			t.Fatalf("text did not scale: %+v -> %+v", before, after)
		}
		for _, label := range []string{"Heading", "Empty", "Status"} {
			r, ok := tt.Find(label)
			if !ok || r.W <= 0 || r.X < 0 || r.X+r.W > 361 {
				t.Fatalf("invalid narrow layout for %s: %+v", label, r)
			}
		}
	}
}
