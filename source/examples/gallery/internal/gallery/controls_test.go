package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/theme"
)

func TestControlPages(t *testing.T) {
	if len(controlCatalog) != 38 {
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
					if example.name == "TextBox" || example.name == "TextBoxBuffer" {
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

func TestPixelGridStandaloneExample(t *testing.T) {
	page := pageFromKey("controls/PixelGrid")
	example, ok := controlExampleFor(page)
	if !ok || example.name != "PixelGrid" {
		t.Fatal("missing pixel grid route")
	}
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, variant := range example.variants {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				controlPixelGrid(c, nil, variant)
			}, 240, 160)
			img := tt.Image()
			minor, major := img.RGBAAt(16, 8), img.RGBAAt(64, 8)
			if minor == img.RGBAAt(8, 8) || (minor != major) != (variant == "Major lines") {
				t.Fatalf("mode=%v variant=%s: incorrect grid emphasis", mode, variant)
			}
		}
	}
}

func TestControlNavigationAndState(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.Left, 220, 6000)
	for i, example := range controlCatalog {
		if err := tt.Click(example.label); err != nil {
			t.Fatal(err)
		}
		if g.page != controlPageStart()+galleryPage(i) {
			t.Fatal("incorrect navigation")
		}
		title := ui.NewTester(g.Titlebar, 540, 60)
		if !title.HasText(example.group + " / " + example.name) {
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
				case "Calendar":
					before := s.day
					tt.Key(0, ui.KeyTab)
					tt.Key(0, ui.KeyRight)
					tt.Key(0, ui.KeyEnter)
					changed = s.day != before
				case "TimePicker":
					click("Time")
					_, changed = tt.Find("Hours")
				case "CalendarPicker":
					click("Date")
					_, changed = tt.Find("Next month")
				case "ColorView":
					before := s.color
					click("Saturation and brightness")
					tt.Key(0, ui.KeyRight)
					changed = s.color != before
				case "ColorPicker":
					click("Default color")
					_, changed = tt.Find("Hue")
				case "Tooltip":
					tt.Key(0, ui.KeyTab)
					changed = tt.HasText("A helpful description of this action")
				case "Popover":
					click("Show popover")
					changed = s.dialogOpen
				case "Dialog":
					click("Open dialog")
					changed = s.dialogOpen
				case "AlertDialog":
					click("Show Default")
					changed = s.dialogOpen
				case "Code":
					click("Copy code")
					changed = tt.Clipboard() == galleryCode()[0].source
				case "Markdown":
					click("Increment")
					changed = s.index == 1
				case "RichText":
					click("inline action")
					changed = s.index == 1
				case "Breadcrumbs":
					click("Components")
					changed = s.index == 1
				case "Rating":
					tt.Key(0, ui.KeyTab)
					tt.Key(0, ui.KeyRight)
					changed = s.index == 1
				case "Badge", "Icon", "Image", "Spinner", "PixelGrid":
					return // Display-only components have no interaction state.
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
				case "TextBoxBuffer":
					click("Name")
					tt.Type("Hello")
					changed = s.buffer.String() == "Hello"
				case "Select":
					click("Size")
					if !disabled {
						click("Large")
					}
					changed = s.choice == "Large"
				case "Progress":
					click("Default value")
					tt.Key(0, ui.KeyRight)
					changed = s.value != 40
				case "Split":
					click("Default divider")
					tt.Key(0, ui.KeyRight)
					changed = s.split != 100
				case "NumberInput":
					click("Quantity")
					tt.Key(0, ui.KeyUp)
					changed = s.value == 41
				case "Slider":
					click("Volume")
					tt.Key(0, ui.KeyRight)
					changed = s.value != 40
				case "Tree":
					click("source")
					click("source")
					changed = s.on
				case "List", "GridView":
					click("Item 2")
					changed = s.index == 1
				case "Sidebar":
					click("Files")
					changed = s.index == 1
				case "Accordion":
					click("General")
					changed = s.toggles[0]
				case "Collapsible":
					click("Details")
					changed = s.on
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

func TestControlPageGrid(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{280, 900} {
			g := &Gallery{page: pageFromKey("controls/Popover")}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c, mode)
				g.View(c)
			}, width, 600)
			img := tt.Image()
			background := img.RGBAAt(5, 575)
			// The grid covers the header, outer margins and empty space far below
			// the examples, with one continuous origin rather than local canvases.
			for _, y := range []int{5, 575} {
				for x := 22; x < width-10; x += 22 {
					if img.RGBAAt(x, y) == background || img.RGBAAt(x+5, y) != background {
						t.Fatalf("mode=%v width=%d: missing page grid at %d,%d", mode, width, x, y)
					}
				}
			}
			if img.RGBAAt(5, 572) == background {
				t.Fatal("missing horizontal grid near the page bottom")
			}
			if img.RGBAAt(176, 575) != img.RGBAAt(22, 575) {
				t.Fatal("page grid must not emphasize major lines")
			}
		}
	}
}

func TestControlVariantLayout(t *testing.T) {
	for _, name := range []string{"Text", "TextBox", "TextBox (multiline)", "Select", "Tabs", "Collapsible", "Image", "Slider"} {
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
				var previous ui.Rect
				for i, variant := range example.variants {
					label := variant
					if name == "Tabs" {
						label = "General"
						if example.label != "Tabs" || len(example.variants) != 1 || tt.HasText("Default") || tt.HasText("Selected") {
							t.Fatal("Tabs should have one named demonstration without redundant state groups")
						}
					}
					r, ok := tt.Find(label)
					if !ok || r.X < 20 || r.X+r.W > float32(width)-20 {
						t.Fatalf("%s/%s outside preview: %+v", name, variant, r)
					}
					if i > 0 {
						vertical := name == "Text" || name == "Collapsible"
						if vertical && r.Y < previous.Y+previous.H {
							t.Fatalf("%s variants are not vertical", name)
						}
						if !vertical && width == 900 && i == 1 && r.Y != previous.Y {
							t.Fatalf("%s variants are not horizontal", name)
						}
						if r.Y < previous.Y+previous.H && r.X < previous.X+previous.W {
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
