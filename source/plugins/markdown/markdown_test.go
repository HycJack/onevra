package markdown

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/theme"
)

func TestMarkdownSelection(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{180, 600} {
			document := ParseMarkdown("# Heading\n\nHello **世界** and *friends*.\n\n```go\nfmt.Println(42)\n```\n\nLast paragraph.")
			tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c, mode); document.Render(c) }, width, 500)
			rect, ok := tt.Find("Heading")
			if !ok {
				t.Fatal(tt.Texts())
			}
			tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
			tt.Key(ui.Cmd, ui.KeyA)
			tt.Key(ui.Cmd, ui.KeyC)
			if got, want := tt.Clipboard(), "Heading\nHello 世界 and friends.\nfmt.Println(42)\nLast paragraph."; got != want {
				t.Fatalf("copy: %q, want %q", got, want)
			}
			last, ok := tt.Find("Last paragraph.")
			if !ok {
				t.Fatal("missing last paragraph")
			}
			tt.Press(rect.X+1, rect.Y+rect.H/2)
			tt.Move(last.X+last.W-1, last.Y+last.H/2)
			tt.Release(last.X+last.W-1, last.Y+last.H/2)
			tt.Key(ui.Cmd, ui.KeyC)
			if !strings.Contains(tt.Clipboard(), "Hello 世界 and friends.\nfmt.Println(42)") {
				t.Fatalf("cross-paragraph drag: %q", tt.Clipboard())
			}
		}
	}
}

func TestCodeHeaderCopyAndSelection(t *testing.T) {
	doc := ParseMarkdown("Before\n\n```go\nfmt.Println(42)\n```\n\nAfter")
	tt := ui.NewTester(func(c *ui.Context) { doc.Render(c, MarkdownOptions{CodeHeaders: true}) }, 400, 400)
	if !tt.HasText("go") {
		t.Fatal("missing language label")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Copy code") {
		t.Fatal("code copy requires mouse hover")
	}
	tt.Key(0, ui.KeyEnter)
	if tt.Clipboard() != "fmt.Println(42)\n" {
		t.Fatal("code copy changed source")
	}
	r, _ := tt.Find("Before")
	tt.ClickAt(r.X+2, r.Y+r.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "Before\nfmt.Println(42)\nAfter" {
		t.Fatalf("header contaminated selection: %q", tt.Clipboard())
	}
}

func TestMarkdownCodeViews(t *testing.T) {
	doc := ParseMarkdown("```unknown\nfallback\n```\n\n~~~counter title=demo\n世界\n~~~\n\n```counter\nsecond\n```\n\n    indented\n")
	var blocks []MarkdownCodeBlock
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		blocks = nil
		doc.Render(c, MarkdownOptions{CodeBlocks: map[string]func(*ui.Context, MarkdownCodeBlock){
			"counter": func(c *ui.Context, b MarkdownCodeBlock) {
				blocks = append(blocks, b)
				components.Text(c, "custom "+b.Code)
				if components.Button(c, "Increment", components.ButtonSecondary).Clicked() {
					clicks++
				}
			},
		}})
	}, 400, 500)
	if len(blocks) != 2 || blocks[0] != (MarkdownCodeBlock{Language: "counter", Info: "counter title=demo", Code: "世界\n", Index: 1}) || blocks[1].Index != 2 {
		t.Fatalf("blocks: %+v", blocks)
	}
	if !tt.HasText("fallback") || !tt.HasText("indented") || !tt.HasText("custom 世界\n") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("Increment"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatalf("clicks %d", clicks)
	}
	rect, _ := tt.Find("fallback")
	tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	copied := tt.Clipboard()
	if !strings.Contains(copied, "custom 世界") || !strings.Contains(copied, "custom second") || strings.Contains(copied, "Increment") {
		t.Fatalf("custom view selection: %q", copied)
	}
}

func TestMarkdownSyntax(t *testing.T) {
	const source = "Text &amp; \\*literal\\* `a &amp; b`  \nnext\nline\n\n> Quote\n\n3. Third\n   - Nested\n\n- [x] Done\n- [ ] Pending\n\n| Left | Right |\n| :--- | ---: |\n| one | ~~two~~ |\n\n![Alt](https://example.com/image.png)\n\n<div>raw</div>\n\n---\n"
	tt := ui.NewTester(func(c *ui.Context) { Markdown(c, source) }, 500, 900)
	for _, want := range []string{"Text & *literal* a &amp; b\nnext line", "Quote", "Third", "Nested", "Done", "Pending", "Left", "Right", "one", "two", "Alt", "<div>raw</div>"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %v", want, tt.Texts())
		}
	}
}

func TestMarkdownLinkAndReplacement(t *testing.T) {
	doc := ParseMarkdown("[Open **guide**](guide.md)\n\n<https://example.com>")
	var destination string
	tt := ui.NewTester(func(c *ui.Context) { doc.Render(c, MarkdownOptions{OnLink: func(url string) { destination = url }}) }, 400, 300)
	rect, ok := tt.Find("Open guide")
	if !ok {
		t.Fatal(tt.Texts())
	}
	tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
	if destination != "guide.md" {
		t.Fatal(destination)
	}
	rect, ok = tt.Find("https://example.com")
	if !ok {
		t.Fatal(tt.Texts())
	}
	tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
	if destination != "https://example.com" {
		t.Fatal(destination)
	}
	doc = ParseMarkdown("Replacement")
	tt.Key(0, ui.KeyTab)
	if !tt.HasText("Replacement") || tt.HasText("Open guide") {
		t.Fatal(tt.Texts())
	}
}

func TestMarkdownEmptyAndIndependent(t *testing.T) {
	var empty MarkdownDocument
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			empty.Render(c)
			Markdown(c, "First")
			Markdown(c, "Second")
		})
	}, 300, 200)
	rect, _ := tt.Find("First")
	tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != "First" {
		t.Fatal(tt.Clipboard())
	}
}

func TestMarkdownImages(t *testing.T) {
	bitmap := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 400, 200)))
	doc := ParseMarkdown("Before\n\n![Photo](photo.png)\n\n![Missing](missing.png)\n\nInline ![small](photo.png) after.\n\nAfter")
	for _, width := range []int{160, 600} {
		tt := ui.NewTester(func(c *ui.Context) {
			doc.Render(c, MarkdownOptions{Images: map[string]*ui.Bitmap{"photo.png": bitmap}})
		}, width, 500)
		rect, ok := tt.Find("Photo")
		if !ok || rect.W != float32(min(width, 400)) || rect.H != rect.W/2 {
			t.Fatalf("image size at %d: %v, found %v", width, rect, ok)
		}
		if !tt.HasText("Missing") || !tt.HasText("Inline small after.") {
			t.Fatal(tt.Texts())
		}
		before, _ := tt.Find("Before")
		tt.ClickAt(before.X+2, before.Y+before.H/2)
		tt.Key(ui.Cmd, ui.KeyA)
		tt.Key(ui.Cmd, ui.KeyC)
		if tt.Clipboard() != "Before\nMissing\nInline small after.\nAfter" {
			t.Fatalf("selection across image: %q", tt.Clipboard())
		}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		doc.Render(c, MarkdownOptions{Images: map[string]*ui.Bitmap{"photo.png": nil}})
	}, 300, 500)
	if !tt.HasText("Photo") {
		t.Fatal("nil bitmap did not fall back to alt text")
	}
}

func TestMarkdownTaskSelection(t *testing.T) {
	document := ParseMarkdown("Before\n\n- [x] **Done**\n  - [ ] Nested\n- Plain item\n- [ ] Pending\n\nAfter")
	tt := ui.NewTester(func(c *ui.Context) { document.Render(c) }, 320, 600)
	before, _ := tt.Find("Before")
	tt.ClickAt(before.X+2, before.Y+before.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	if want := "Before\nDone\nNested\nPlain item\nPending\nAfter"; tt.Clipboard() != want {
		t.Fatalf("task text must remain in document selection: %q", tt.Clipboard())
	}
	if tt.HasText("☑ Done") || tt.HasText("☐ Pending") {
		t.Fatal("legacy text task markers")
	}
}

func TestMarkdownReferenceDefinitionsDoNotReserveSpace(t *testing.T) {
	for _, width := range []int{180, 600} {
		render := func(source string) *ui.Tester {
			return ui.NewTester(func(c *ui.Context) { Markdown(c, source) }, width, 300)
		}
		inline := render("[Guide](guide.md)\n\n## Lists and quotes\n\n- Item")
		reference := render("[Guide][guide]\n\n[guide]: guide.md\n\n## Lists and quotes\n\n- Item")
		for _, label := range []string{"Guide", "Lists and quotes", "Item"} {
			want, wantOK := inline.Find(label)
			got, gotOK := reference.Find(label)
			if !wantOK || !gotOK || got != want {
				t.Fatalf("reference definition changed %q layout: got %v, want %v", label, got, want)
			}
		}
	}
}

func TestMarkdownTaskAppearanceAndInput(t *testing.T) {
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		checked := true
		reference := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			components.Checkbox(c, &checked, "")
		}, 100, 60)
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, mode)
			Markdown(c, "- [x] Done")
		}, 200, 100)
		rect, ok := tt.Find("[x] Done")
		if !ok {
			t.Fatal("missing task checkbox")
		}
		sample := func() color.RGBA { return color.RGBA(tt.Image().RGBAAt(int(rect.X)+3, int(rect.Y)+3)) }
		want := color.RGBA(reference.Image().RGBAAt(3, 3))
		if sample() != want {
			t.Fatalf("task checkbox is dimmed: got %v want %v", sample(), want)
		}
		tt.ClickAt(rect.X+rect.W/2, rect.Y+rect.H/2)
		if tt.Focused("[x] Done") {
			t.Fatal("pointer focused a pass-through checkbox")
		}
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeySpace)
		tt.Key(0, ui.KeyEnter)
		if sample() != want {
			t.Fatal("task checkbox appearance changed after input")
		}
	}
}

func TestMarkdownCopyCodeOnHover(t *testing.T) {
	for _, source := range []string{"```go\nhello 世界\n```", "    hello 世界\n"} {
		tt := ui.NewTester(func(c *ui.Context) { Markdown(c, source) }, 400, 200)
		tt.Move(390, 190)
		if _, ok := tt.Find("Copy code"); ok {
			t.Fatal("copy button shown without hover")
		}
		code, ok := tt.Find("hello 世界")
		if !ok {
			t.Fatal(tt.Texts())
		}
		tt.Move(code.X+2, code.Y+code.H/2)
		if _, ok := tt.Find("Copy code"); !ok {
			t.Fatal("missing hover copy button")
		}
		if err := tt.Click("Copy code"); err != nil {
			t.Fatal(err)
		}
		if tt.Clipboard() != "hello 世界\n" {
			t.Fatalf("code copy: %q", tt.Clipboard())
		}
		tt.Move(390, 190)
		if _, ok := tt.Find("Copy code"); ok {
			t.Fatal("copy button remained after leaving block")
		}
	}
}

func TestMarkdownPreparedHighlighting(t *testing.T) {
	d := ParseMarkdown("```Go extra=info\nconst n = 42\n```\n\n```unknown\nraw\n```\n\n    plain\n")
	highlighted := 0
	for _, prepared := range d.code {
		if prepared.text.Highlighted() {
			highlighted++
		}
	}
	if len(d.code) != 3 || highlighted != 1 {
		t.Fatalf("blocks=%d highlighted=%d", len(d.code), highlighted)
	}
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		d.Render(c, MarkdownOptions{CodeBlocks: map[string]func(*ui.Context, MarkdownCodeBlock){
			"Go": func(c *ui.Context, block MarkdownCodeBlock) {
				calls++
				if block.Info != "Go extra=info" || block.Code != "const n = 42\n" || block.Index != 0 {
					t.Fatalf("block: %+v", block)
				}
				components.Text(c, "replacement")
			},
		}})
	}, 400, 300)
	if calls == 0 || !tt.HasText("replacement") || tt.HasText("const n = 42") {
		t.Fatal(tt.Texts())
	}
}
