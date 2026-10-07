package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestControlPages(t *testing.T) {
	if len(controlCatalog) != 9 {
		t.Fatal("missing control examples")
	}
	for i, example := range controlCatalog {
		t.Run(example.name, func(t *testing.T) {
			page := controlPageStart() + galleryPage(i)
			if got, ok := controlExampleFor(page); !ok || got.name != example.name {
				t.Fatal("incorrect route")
			}
			if _, ok := mygoExampleFor(page); ok {
				t.Fatal("route overlaps MyGo")
			}
			for _, width := range []int{540, 280} {
				for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
					g := &Gallery{page: page}
					tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c, mode); g.View(c) }, width, 1800)
					if tt.Image() == nil || len(g.controlStates) != len(example.variants) {
						t.Fatal("missing page or variant")
					}
					if tt.HasText("Disabled") {
						t.Fatal("unexpected disabled example")
					}
					for _, variant := range example.variants {
						if !tt.HasText(variant) {
							t.Fatalf("missing variant %s", variant)
						}
					}
				}
			}
		})
	}
	for _, page := range []galleryPage{shellPage, lucidePage, mygoPageStart, controlPageStart() - 1, controlPageStart() + galleryPage(len(controlCatalog))} {
		if _, ok := controlExampleFor(page); ok {
			t.Fatalf("unexpected route %d", page)
		}
	}
}

func TestControlNavigationAndState(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 1600)
	for i, example := range controlCatalog {
		if err := tt.Click(example.label); err != nil {
			t.Fatal(err)
		}
		if g.page != controlPageStart()+galleryPage(i) {
			t.Fatal("incorrect navigation")
		}
		title := ui.NewTester(g.Titlebar, 540, 60)
		if !title.HasText("Components / " + example.name) {
			t.Fatal("incorrect title")
		}
	}
	off, on := g.controlState("Checkbox", "Off"), g.controlState("Checkbox", "On")
	if off.on || !on.on {
		t.Fatal("incorrect initial states")
	}
	off.on = true
	if g.controlState("Checkbox", "Off") != off || g.controlState("Switch", "Off").on {
		t.Fatal("state lost or shared between pages")
	}
	on.on = false
	if !off.on {
		t.Fatal("state shared between variants")
	}
	if g.mygoState("Checkbox", "Off", true).on {
		t.Fatal("state shared with MyGo")
	}
}

func TestControlExampleInteractions(t *testing.T) {
	for _, example := range controlCatalog {
		for _, disabled := range []bool{false, true} {
			t.Run(example.name+map[bool]string{false: "/enabled", true: "/disabled"}[disabled], func(t *testing.T) {
				s := (&Gallery{}).controlState(example.name, "Empty")
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c, theme.Light)
					ui.Column(c).Fill().Gap(12).Disabled(disabled).Children(func() { example.build(c, s, "Default") })
				}, 540, 900)
				click := func(label string) {
					t.Helper()
					if err := tt.Click(label); err != nil {
						t.Fatal(err)
					}
				}
				changed := false
				switch example.name {
				case "Button":
					click("Primary")
					changed = s.clicks == 1
				case "Checkbox":
					click("Notifications")
					changed = s.on
				case "Switch":
					click("Sound")
					changed = s.on
				case "Radio":
					click("Large")
					changed = s.choice == "Large"
				case "TextField", "TextArea":
					label := "Name"
					if example.name == "TextArea" {
						label = "Notes"
					}
					click(label)
					tt.Type("Hello")
					changed = s.text == "Hello"
				case "Select":
					click("Size")
					if !disabled {
						click("Large")
					}
					changed = s.choice == "Large"
				case "Slider":
					click("Volume")
					tt.Key(0, ui.KeyRight)
					changed = s.value != 40
				case "Tabs":
					click("Style")
					changed = s.index == 1
				}
				if changed == disabled {
					t.Fatalf("changed=%t, disabled=%t", changed, disabled)
				}
			})
		}
	}
}

func TestControlButtonWrap(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{540, 200} {
			s := &controlState{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Fill().Gap(12).Children(func() { controlButton(c, s, "Default") })
			}, width, 300)
			var first, previous ui.Rect
			wrapped := false
			for i, label := range []string{"Primary", "Secondary", "Outline", "Ghost", "Link"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing button %s", label)
				}
				if r.X < 0 || r.X+r.W > float32(width)+1 {
					t.Fatalf("%s overflows width %d: %+v", label, width, r)
				}
				if i == 0 {
					first = r
				} else {
					if r.Y > previous.Y+1 {
						wrapped = true
					} else if r.X < previous.X+previous.W {
						t.Fatalf("%s overlaps the previous button", label)
					}
					if width == 540 && r.Y+r.H/2 != first.Y+first.H/2 {
						t.Fatal("wide layout should keep buttons on one row")
					}
					if i == 1 && r.Y+r.H/2 != first.Y+first.H/2 {
						t.Fatal("the first two buttons should fit side by side")
					}
				}
				previous = r
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			if width == 200 && !wrapped {
				t.Fatal("narrow layout should wrap buttons")
			}
			if s.clicks != 5 {
				t.Fatalf("got %d clicks, want 5", s.clicks)
			}
		}
	}
}

func TestControlButtonPointerStyles(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, disabled := range []bool{false, true} {
			s := &controlState{}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				ui.Column(c).Fill().Disabled(disabled).Children(func() { controlButton(c, s, "Default") })
			}, 540, 300)
			for _, label := range []string{"Primary", "Secondary", "Outline", "Ghost"} {
				r, ok := tt.Find(label)
				if !ok {
					t.Fatalf("missing button %s", label)
				}
				x, y := r.X+6, r.Y+r.H/2
				tt.Move(530, 290)
				rest := tt.Image().RGBAAt(int(x), int(y))
				tt.Move(x, y)
				hover := tt.Image().RGBAAt(int(x), int(y))
				tt.Press(x, y)
				pressed := tt.Image().RGBAAt(int(x), int(y))
				tt.Release(x, y)
				tt.Move(530, 290)
				reset := tt.Image().RGBAAt(int(x), int(y))
				if disabled {
					if hover != rest || pressed != rest {
						t.Fatalf("%v/%s: disabled button changed on pointer input", mode, label)
					}
				} else if hover == rest || pressed == hover {
					t.Fatalf("%v/%s: missing hover or pressed feedback", mode, label)
				}
				if reset != rest {
					t.Fatalf("%v/%s: pointer exit did not restore resting style", mode, label)
				}
			}
		}
	}
}

func TestControlReadOnlyExamples(t *testing.T) {
	for _, example := range controlCatalog {
		if example.name != "TextField" && example.name != "TextArea" {
			continue
		}
		s := (&Gallery{}).controlState(example.name, "Read only")
		before := s.text
		tt := ui.NewTester(func(c *ui.Context) { example.build(c, s, "Read only") }, 540, 300)
		label := "Name"
		if example.name == "TextArea" {
			label = "Notes"
		}
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
		tt.Type("Changed")
		if s.text != before {
			t.Fatal("read-only example accepted input")
		}
	}
}
