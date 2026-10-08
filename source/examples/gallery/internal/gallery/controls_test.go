package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestControlPages(t *testing.T) {
	if len(controlCatalog) != 12 {
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
					wantStates := len(example.variants)
					if example.name == "TextBox" {
						wantStates = 7
					}
					if tt.Image() == nil || len(g.controlStates) != wantStates {
						t.Fatal("missing page or variant")
					}
					if tt.HasText("Disabled") {
						t.Fatal("unexpected disabled example")
					}
					if example.name != "Selection" && !tt.HasText(example.name) {
						t.Fatalf("missing component heading %s", example.name)
					}
					if len(example.variants) > 1 {
						for _, variant := range example.variants {
							if !tt.HasText(variant) {
								t.Fatalf("missing variant %s", variant)
							}
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

func TestPixelGridHasNoStandaloneExample(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 1800)
	if tt.HasText("Pixel grid") || tt.HasText("PixelGrid") {
		t.Fatal("pixel grid sidebar entry remains")
	}
	if pageFromKey("controls/PixelGrid") != lucidePage {
		t.Fatal("removed pixel grid route must use the unknown-route fallback")
	}
	for _, example := range controlCatalog {
		if example.name == "PixelGrid" {
			t.Fatal("pixel grid standalone example remains")
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

func TestEditableTextExample(t *testing.T) {
	g := &Gallery{}
	s := g.controlState("Text", "Editable")
	tt := ui.NewTester(func(c *ui.Context) { controlText(c, s, "Editable") }, 540, 100)
	tt.Key(0, ui.KeyTab)
	tt.Type("Edited name")
	if s.text != "Edited name" {
		t.Fatalf("text binding not updated: %q", s.text)
	}
	if g.controlState("Text", "Default").text != "Hello Onevra" {
		t.Fatal("text variants share state")
	}
	if mygoHandled["EditableText"].Replacement != "components.Text" {
		t.Fatal("EditableText is not merged into Text")
	}
}

func TestControlExampleInteractions(t *testing.T) {
	for _, example := range append(append([]controlExample{}, controlCatalog...), textBoxControlCatalog[1:]...) {
		if example.name == "Selection" {
			continue // Grouped selection interactions are covered separately.
		}
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
				case "Icon", "Image":
					return // Media components have no interaction state.
				case "ScrollView":
					return // ScrollView interactions are covered in the component tests.
				case "Text":
					if !tt.HasText("Count: 42") || !tt.HasText("Inline bold") {
						t.Fatal("missing display text")
					}
					return // Display text has no interaction state.
				case "Menu":
					click("Actions")
					if !disabled {
						if err := tt.ChooseMenuItem("Open"); err != nil {
							t.Fatal(err)
						}
					}
					changed = s.index == 1
				case "Button":
					click("Primary")
					changed = s.button.clicks[0] == 1
				case "TextBox", "TextBox (multiline)":
					label := "Name"
					if example.name == "TextBox (multiline)" {
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
				case "Sidebar":
					click("Files")
					changed = s.index == 1
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

func TestControlPreviewGeometry(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			controlPreview(c).Height(120).Children(func() {
				ui.Box(c).Size(10, 10).Label("content")
			})
		}, 200, 120)
		r, ok := tt.Find("content")
		// The one-pixel border sits outside the 20-pixel padding.
		if !ok || r.X != 21 || r.Y != 21 {
			t.Fatalf("preview padding: %+v", r)
		}
		img := tt.Image()
		background := img.RGBAAt(50, 50)
		for _, x := range []int{22, 44, 66, 88} {
			if img.RGBAAt(x, 50) == background || img.RGBAAt(x+5, 50) != background {
				t.Fatalf("missing 22-pixel grid at x=%d", x)
			}
		}
		if img.RGBAAt(100, 0) == background {
			t.Fatal("missing preview border")
		}
		if img.RGBAAt(0, 0) == img.RGBAAt(100, 0) {
			t.Fatal("preview corner is not rounded")
		}
	}
}

func TestControlVariantLayout(t *testing.T) {
	for _, name := range []string{"TextBox", "TextBox (multiline)", "Select", "Tabs"} {
		for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
			for _, width := range []int{900, 280} {
				g := &Gallery{page: pageFromKey("controls/" + name)}
				example, _ := controlExampleFor(g.page)
				for _, textExample := range textBoxControlCatalog {
					if textExample.name == name {
						example = textExample
					}
				}
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c, mode)
					ui.Column(c).Fill().Padding(24).Children(func() { g.controlVariants(c, example) })
				}, width, 1200)
				var first, previous ui.Rect
				for i, variant := range example.variants {
					r, ok := tt.Find(variant)
					if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
						t.Fatalf("%s/%s outside preview: %+v", name, variant, r)
					}
					if i == 0 {
						first = r
					} else {
						if width == 900 && i < 3 && (r.Y != first.Y || r.X <= previous.X) {
							t.Fatalf("%s variants are not arranged horizontally", name)
						}
						if r.Y == previous.Y && r.X < previous.X+previous.W {
							t.Fatalf("%s variants overlap", name)
						}
					}
					previous = r
				}
				for _, text := range []string{"Selected:", "characters", "General settings", "Style settings", "Grid spacing:"} {
					if tt.HasText(text) {
						t.Fatalf("%s retains extra status text %q", name, text)
					}
				}
				if name == "TextBox" || name == "TextBox (multiline)" {
					label := "Name"
					if name == "TextBox (multiline)" {
						label = "Notes"
					}
					if err := tt.Click(label); err != nil {
						t.Fatal(err)
					}
					tt.Type("Edited")
					if g.controlState(name, "Empty").text != "Edited" || g.controlState(name, "Filled").text == "Edited" || g.controlState(name, "Read only").text == "Edited" {
						t.Fatalf("%s variant values are not independent", name)
					}
				}
			}
		}
	}
}

func TestControlReadOnlyExamples(t *testing.T) {
	for _, example := range textBoxControlCatalog {
		if example.name != "TextBox" && example.name != "TextBox (multiline)" {
			continue
		}
		s := (&Gallery{}).controlState(example.name, "Read only")
		before := s.text
		tt := ui.NewTester(func(c *ui.Context) { example.build(c, s, "Read only") }, 540, 300)
		label := "Name"
		if example.name == "TextBox (multiline)" {
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
