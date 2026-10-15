package gallery

import (
	"sync"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/code"
)

type codeSample struct {
	language, source string
	document         *code.Document
}

// Retain parsed samples across navigation and frames, independently of UI state.
var galleryCode = sync.OnceValue(func() []codeSample {
	samples := []codeSample{
		{language: "Go", source: "// Unicode: 世界\nfunc greet(name string) string {\n    return \"Hello, \" + name\n}\n"},
		{language: "TypeScript", source: "interface User { name: string; active: boolean }\nconst welcome = (user: User) => `Hello, ${user.name}`;\n"},
		{language: "Python", source: "def greet(name):\n    return f\"Hello, {name}\"\n"},
		{language: "JSON", source: "{\n    \"name\": \"Onevra\",\n    \"readonly\": true,\n    \"count\": 42\n}\n"},
		{language: "SQL", source: "SELECT name, created_at\nFROM users\nWHERE active = true\nORDER BY created_at DESC;\n"},
		{language: "Bash", source: "#!/bin/bash\nfor name in Onevra MyGo; do\n    echo \"Hello, $name\"\ndone\n"},
		{language: "Plaintext", source: "Unknown languages and plain text remain selectable.\n    Indentation and source text are preserved.\n"},
	}
	for i := range samples {
		samples[i].document = code.Parse(samples[i].source, samples[i].language)
	}
	return samples
})

func controlCode(c *ui.Context, _ *controlState, _ string) {
	t := c.Theme()
	ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(4)).Children(func() {
		for _, sample := range galleryCode() {
			ui.Column(c).Key(sample.language).FillWidth().MinWidth(0).Gap(t.Space(2)).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
					components.Text(c, sample.language).Grow(1).TextColor(t.TextMuted)
					if components.Button(c, tr(c, "copyCode"), components.ButtonGhost).Clicked() {
						c.WriteClipboard(sample.document.Text())
					}
				})
				sample.document.Render(c).Selectable().FillWidth().FontSize(t.Rem(.85)).LineHeight(1.6).
					Padding(t.Space(3)).Radius(t.Rem(.75)).Border(1, t.Border).Background(t.Background)
			})
		}
	})
}
