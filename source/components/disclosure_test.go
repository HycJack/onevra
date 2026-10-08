package components

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestDisclosureAppearance(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, accordion := range []bool{false, true} {
			for _, open := range []bool{false, true} {
				build := func(upstream bool) *ui.Tester {
					return ui.NewTester(func(c *ui.Context) {
						theme.Apply(c, mode)
						body := func() { Text(c, "Content") }
						if accordion {
							if upstream {
								ui.Accordion(c, func() { ui.AccordionItem(c, "Details", &open, body); ui.AccordionItem(c, "Other", &open, body) })
							} else {
								Accordion(c, func(g *AccordionGroup) { g.Item("Details", &open, body); g.Item("Other", &open, body) })
							}
						} else if upstream {
							ui.Collapsible(c, "Details", &open, body)
						} else {
							Collapsible(c, "Details", &open, CollapsibleOptions{ArrowPosition: ArrowLeft}, body)
						}
					}, 320, 240)
				}
				ours, upstream := build(false), build(true)
				before := bytes.Clone(ours.Image().Pix)
				if !bytes.Equal(before, upstream.Image().Pix) {
					t.Fatalf("resting appearance differs: mode=%v accordion=%v open=%v", mode, accordion, open)
				}
				r, _ := ours.Find("Details")
				ours.Move(r.X+r.W/2, r.Y+r.H/2)
				if !bytes.Equal(before, ours.Image().Pix) {
					t.Fatal("hover changed appearance")
				}
			}
		}
	}
}

func TestCollapsibleInteractionAndPlacement(t *testing.T) {
	for _, position := range []ArrowPosition{ArrowRight, ArrowLeft} {
		open, disabled, changes := false, false, 0
		var trigger ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, theme.Light)
			ui.Column(c).Disabled(disabled).Children(func() {
				p := Collapsible(c, "Details", &open, CollapsibleOptions{ArrowPosition: position}, func() { Button(c, "Inside", ButtonSecondary) })
				p.Trigger.DrawOver(func(_ *ui.Painter, r ui.Rect) { trigger = r })
				if p.Changed() {
					changes++
				}
			})
		}, 320, 240)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true})
		label, _ := tt.Find("Details")
		left, right := label.X-trigger.X, trigger.X+trigger.W-label.X-label.W
		if position == ArrowRight && right <= left || position == ArrowLeft && left <= right {
			t.Fatalf("wrong arrow position: left=%v right=%v", left, right)
		}
		if tt.HasText("Inside") {
			t.Fatal("closed content built")
		}
		if err := tt.Click("Details"); err != nil {
			t.Fatal(err)
		}
		if !open || changes != 1 || !tt.HasText("Inside") {
			t.Fatalf("click: open=%v changes=%d", open, changes)
		}
		tt.Key(0, ui.KeyTab)
		if !tt.Focused("Inside") {
			t.Fatal("Tab skipped panel")
		}
		tt.Key(ui.Shift, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		if open || changes != 2 || tt.HasText("Inside") {
			t.Fatal("Space did not close")
		}
		tt.Key(0, ui.KeyEnter)
		if !open || changes != 3 {
			t.Fatal("Enter did not open")
		}
		disabled = true
		tt.Frame()
		tt.Click("Details")
		tt.Key(0, ui.KeySpace)
		if !open || changes != 3 {
			t.Fatal("disabled disclosure changed")
		}
	}
}

func TestAccordionKeyboard(t *testing.T) {
	states := [3]bool{}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Accordion(c, func(g *AccordionGroup) {
			for i, title := range []string{"General", "Privacy", "Advanced"} {
				if g.Item(title, &states[i], func() { Button(c, "Inside "+title, ButtonSecondary) }).Changed() {
					changes++
				}
			}
		})
	}, 320, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Key(0, ui.KeyTab)
	for _, step := range []struct {
		key  ui.Key
		want string
	}{
		{ui.KeyDown, "Privacy"}, {ui.KeyDown, "Advanced"}, {ui.KeyDown, "Advanced"},
		{ui.KeyHome, "General"}, {ui.KeyUp, "General"}, {ui.KeyEnd, "Advanced"}, {ui.KeyUp, "Privacy"},
	} {
		tt.Key(0, step.key)
		if !tt.Focused(step.want) {
			t.Fatalf("key %v: not focused %s", step.key, step.want)
		}
	}
	tt.Key(0, ui.KeyEnter)
	if !states[1] || changes != 1 {
		t.Fatal("Enter did not open")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Inside Privacy") {
		t.Fatal("Tab skipped panel")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeySpace)
	if !states[0] || !states[1] || changes != 2 {
		t.Fatal("sections not independent")
	}
}

func TestDisclosureAnimation(t *testing.T) {
	open := false
	var root ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		p := Collapsible(c, "Details", &open, CollapsibleOptions{}, func() { ui.Box(c).Height(100) })
		p.Root.Draw(func(_ *ui.Painter, r ui.Rect) { root = r })
	}, 320, 240)
	closed := root.H
	tt.Click("Details")
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	middle := root.H
	time.Sleep(220 * time.Millisecond)
	tt.Frame()
	if middle <= closed || middle >= root.H {
		t.Fatalf("no expansion animation: %v %v %v", closed, middle, root.H)
	}
	tt.Click("Details")
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if root.H <= closed {
		t.Fatal("closing snapped shut")
	}
	time.Sleep(220 * time.Millisecond)
	tt.Frame()
	if root.H != closed {
		t.Fatal("closing did not finish")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Click("Details")
	if root.H <= closed {
		t.Fatal("reduced motion did not open immediately")
	}
	tt.Click("Details")
	if root.H != closed {
		t.Fatal("reduced motion animated closing")
	}
}
