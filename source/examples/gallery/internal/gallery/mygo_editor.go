package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/oneo-me/onevra/source/components"
)

// Editable fixture content stays in English, as other sample documents do.
const editorSample = `package main

import "fmt"

// greet returns a greeting for name.
func greet(name string) string {
	return "Hello, " + name
}

func main() {
	for _, name := range []string{"MyGo", "Onevra"} {
		fmt.Println(greet(name))
	}
}
`

const editorRevision = `package main

import (
	"fmt"
	"strings"
)

// greet returns a greeting for each name.
func greet(names ...string) string {
	return "Hello, " + strings.Join(names, " and ")
}

func main() {
	fmt.Println(greet("MyGo", "Onevra"))
}
`

func editorBuffer(text string) *editor.Buffer {
	b := editor.NewBuffer(text)
	b.SetLanguage(editor.LanguageFor("main.go"))
	return b
}

// exampleCodeEditor shows edits against the original text in the gutter, or a
// highlighted read-only snippet as tall as its lines.
// Buffers live in page memory and are never written to disk.
func exampleCodeEditor(c *ui.Context, s *demoState) {
	if s.variant == "Read only" {
		if s.code == nil {
			s.code = &editor.State{Buffer: editorBuffer(editorSample), ReadOnly: true, FitContent: true,
				LineNumbers: editor.LineNumbersOff, NoFolding: true, NoActiveLine: true, NoIndentGuides: true}
		}
		s.code.Label = tr(c, "codeEditorLabel")
		editor.View(c, s.code).FillWidth().Radius(8).Clip()
		return
	}
	if s.code == nil {
		s.code = &editor.State{Buffer: editorBuffer(editorSample), DiffBase: editorBuffer(editorSample), WrapGuides: []int{80}}
	}
	s.code.Label = tr(c, "codeEditorLabel")
	editor.View(c, s.code).Height(320).FillWidth().Radius(8).Clip()
	pos := s.code.CursorPos()
	components.Text(c, msg(c).CodeEditorPosition(pos.Line+1, pos.Col+1)).TextColor(c.Theme().TextMuted)
}

func exampleDiffViewer(c *ui.Context, s *demoState) {
	if s.diff == nil {
		s.diff = &editor.DiffState{Old: editorBuffer(editorSample), New: editorBuffer(editorRevision), Editable: true, HunkActions: true}
		if s.variant == "Unified" {
			s.diff.Layout = editor.DiffUnified
		}
	}
	editor.DiffView(c, s.diff).Height(360).FillWidth().Radius(8).Clip()
	added, deleted := s.diff.Stats()
	components.Text(c, msg(c).DiffViewerStats(added, deleted)).TextColor(c.Theme().TextMuted)
}
