// Documents is a notes editor of Markdown and text files on mygo.NewDocuments,
// built from Onevra components: MyGo owns the File menu, a window per note,
// autosave, crash recovery and files changed by other apps.
package main

import (
	"io"
	"log"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/examples/documents/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

// note is the window of a document and the text its editor shows.
type note struct {
	doc  *mygo.Document[string]
	text string
}

func (n *note) view(c *ui.Context) {
	t, m := c.Theme(), locales.In(c.Locale())
	state := n.doc.State()
	ui.Column(c).Fill().Gap(t.Space(2)).Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() {
		if state.Conflict {
			ui.Row(c).Padding(t.Space(2.5)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Warning.Alpha(.15)).AlignItems(ui.Center).Children(func() {
				components.Text(c, m.Conflict()).Wrap().Grow(1)
				if components.Button(c, m.Revert(), components.ButtonOutline).Clicked() {
					go n.doc.Revert()
				}
				if components.Button(c, m.Save(), components.ButtonPrimary).Clicked() {
					go n.doc.SaveAnyway()
				}
			})
		}
		editor := components.TextBox(c, &n.text, components.TextBoxOptions{Multiline: true}).
			Label(m.Note()).Placeholder(m.Placeholder()).Grow(1).FillWidth().AutoFocus()
		if editor.Changed() {
			n.doc.Set(n.text)
		}
		ui.Row(c).Gap(t.Space(3)).Children(func() {
			components.Text(c, m.Words(len(strings.Fields(n.text)))).FontSize(12).TextColor(t.TextMuted).Grow(1)
			components.Text(c, status(m, state)).FontSize(12).TextColor(t.TextMuted)
		})
	})
}

func status(m locales.Messages, state mygo.DocumentState) string {
	switch {
	case state.Path == "":
		return m.NotSaved()
	case state.ReadOnly:
		return m.ReadOnly()
	case state.Edited:
		return m.Edited()
	}
	return m.Saved()
}

func main() {
	mygo.NewDocuments(mygo.DocumentsOptions[string]{
		Types: []mygo.DocumentType{
			{Name: "Markdown Document", Extensions: []string{"md", "markdown"}},
			{Name: "Plain Text", Extensions: []string{"txt"}},
		},
		Read: func(r io.Reader) (string, error) {
			b, err := io.ReadAll(r)
			return string(b), err
		},
		Write: func(w io.Writer, text string) error {
			_, err := io.WriteString(w, text)
			return err
		},
		// Undoing back to the file's text leaves the note unedited.
		Equal:  func(a, b string) bool { return a == b },
		Window: newNoteWindow,
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func newNoteWindow(doc *mygo.Document[string]) *mygo.Window {
	n := &note{doc: doc, text: doc.Content()}
	win := plain.NewShell(plain.ShellOptions{Width: 640, Height: 480, StateKey: "note", Theme: theme.Onevra, Content: n.view}).Window()
	// Reverted, or changed by another app: show the file's text.
	doc.OnReload(func() { win.Update(func() { n.text = doc.Content() }) })
	doc.OnStateChange(func(mygo.DocumentState) { win.Update(func() {}) })
	return win
}
