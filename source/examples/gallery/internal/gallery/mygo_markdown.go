package gallery

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/oneo-me/onevra/source/components"
)

// Read the resource once, independently of frames, including its error.
var galleryMarkdown = sync.OnceValues(func() (string, error) {
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		return "", err
	}
	return loadGalleryMarkdown(dir)
})

func loadGalleryMarkdown(dir string) (string, error) {
	source, err := os.ReadFile(filepath.Join(dir, "markdown.md"))
	return string(source), err
}

// exampleMarkdown shows a finished document with resource images, gallery-owned
// link activation and an application-built code block.
func exampleMarkdown(c *ui.Context, s *demoState) {
	source, err := galleryMarkdown()
	if err != nil {
		components.Text(c, tr(c, "markdownLoadError")).TextColor(c.Theme().TextMuted)
		return
	}
	if s.markdown.LoadImage == nil {
		// Missing images keep their alt text.
		if dir, err := mygo.App.Path(mygo.PathResources); err == nil {
			s.markdown.LoadImage = markdown.FileImages(dir)
		}
		s.markdown.OnLink = func(string) { s.index++ }
	}
	s.markdown.Code = map[string]func(*ui.Context, *markdown.CodeBlockView){
		"counter": func(c *ui.Context, block *markdown.CodeBlockView) {
			ui.Column(c).Gap(c.Theme().Space(2)).Children(func() {
				components.Text(c, strings.TrimSuffix(block.Code(), "\n")).Bold()
				if components.Button(c, tr(c, "increment"), components.ButtonSecondary).Clicked() {
					s.index++
				}
				components.Text(c, msg(c).CountD(s.index))
			})
		},
	}
	markdown.View(c, source, &s.markdown).FillWidth().MaxWidth(640)
}
