package code

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestLanguagesAndExactSource(t *testing.T) {
	samples := map[string]string{
		"go":     "package main\n// 世界\nfunc main() { println(42) }\n",
		"JS":     "const greeting = `Hello ${name}`; // comment\n",
		"ts":     "interface User { name: string }; const n: number = 42;",
		"tsx":    "const App = () => <div title=\"hello\">World</div>;",
		"py":     "def greet(name):\n\treturn f'Hello {name}'\n",
		"rust":   "fn main() { let x = 42; println!(\"{}\", x); }",
		"c":      "#include <stdio.h>\nint main() { return 0; }",
		"cpp":    "#include <string>\nstd::string name = \"hello\";",
		"csharp": "public class A { string Name = \"Hello\"; }",
		"java":   "class A { public static void main(String[] args) {} }",
		"swift":  "let name: String = \"Hello\" // comment",
		"kotlin": "fun main() { val name = \"Hello\" }",
		"sh":     "#!/bin/sh\necho \"hello $USER\"\n",
		"json":   "{\"message\": \"你好\", \"count\": 42, \"ok\": true}",
		"yaml":   "message: hello\ncount: 42\n",
		"html":   "<div class=\"greeting\">Hello</div>",
		"css":    ".greeting { color: red; margin: 2px; }",
		"sql":    "SELECT name FROM users WHERE id = 42;",
	}
	for lang, source := range samples {
		t.Run(lang, func(t *testing.T) {
			for _, source := range []string{source, strings.ReplaceAll(source, "\n", "\r\n"), source + "\n\n"} {
				d := Parse(source, lang)
				if !d.Highlighted() {
					t.Fatalf("no highlighting: %q", source)
				}
				var joined strings.Builder
				colored := false
				for _, span := range d.light {
					joined.WriteString(span.Text)
					colored = colored || span.Color != (ui.Color{})
				}
				if joined.String() != source || d.Text() != source || !colored {
					t.Fatalf("lost source or colors: %q", joined.String())
				}
			}
		})
	}
}

func TestFallbackAndLimits(t *testing.T) {
	for _, tc := range []struct{ source, language string }{
		{"", "go"}, {"unchanged\t\r\n\n", ""}, {"unknown", "no-such-language"},
		{"const x = 1", "text"}, {"\xff", "go"},
		{strings.Repeat("x", MaxHighlightBytes+1), "go"},
		{strings.Repeat("x=1;", MaxHighlightRuns), "go"},
	} {
		d := Parse(tc.source, tc.language)
		if d.Highlighted() || d.Text() != tc.source {
			t.Fatalf("fallback failed for %q (%d bytes)", tc.language, len(tc.source))
		}
	}
	var d *Document
	if d.Text() != "" || d.Highlighted() {
		t.Fatal("nil document")
	}
}

func TestSelectionAndThemeChange(t *testing.T) {
	source := "// 世界\nconst n = 42\nprintln(\"hello\")"
	d := Parse(source, "go")
	mode := theme.Light
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, mode)
		d.Render(c).Selectable()
	}, 200, 250)
	light := bytes.Clone(tt.Image().Pix)
	rect, ok := tt.Find(source)
	if !ok {
		t.Fatal(tt.Texts())
	}
	tt.ClickAt(rect.X+2, rect.Y+rect.H/2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != source {
		t.Fatalf("copy: %q", tt.Clipboard())
	}
	mode = theme.Dark
	tt.Key(0, ui.KeyEscape)
	if bytes.Equal(light, tt.Image().Pix) {
		t.Fatal("theme did not change")
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(ui.Cmd, ui.KeyC)
	if tt.Clipboard() != source {
		t.Fatalf("dark copy: %q", tt.Clipboard())
	}
}

func TestConcurrentParse(t *testing.T) {
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 10 {
				if !Parse("const x = 42", "js").Highlighted() {
					t.Error("parse failed")
				}
			}
		})
	}
	wg.Wait()
}

func BenchmarkParse(b *testing.B) {
	for _, lines := range []int{20, 200, 500} {
		source := strings.Repeat("const answer = 42 // hello\n", lines)
		b.Run(fmt.Sprint(lines), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(source)))
			if !Parse(source, "go").Highlighted() {
				b.Fatal("benchmark exceeded highlight limits")
			}
			for b.Loop() {
				Parse(source, "go")
			}
		})
	}
}

func BenchmarkRetainedRender(b *testing.B) {
	d := Parse(strings.Repeat("const answer = 42 // hello\n", 200), "go")
	tt := ui.NewTester(func(c *ui.Context) { d.Render(c) }, 600, 400)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tt.Frame()
	}
}

func ExampleDocument_Render() {
	// Retain this result in application state, outside the frame builder.
	document := Parse("fmt.Println(42)", "go")
	_ = func(c *ui.Context) {
		document.Render(c).Selectable().FontSize(13)
	}
}
