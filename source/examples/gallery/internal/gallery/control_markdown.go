package gallery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/plugins/markdown"
)

// Cache source parsing independently of frames, including resource errors.
var galleryMarkdown = sync.OnceValues(func() (*markdown.MarkdownDocument, error) {
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		return nil, err
	}
	return loadGalleryMarkdown(dir)
})

func loadGalleryMarkdown(dir string) (*markdown.MarkdownDocument, error) {
	source, err := os.ReadFile(filepath.Join(dir, "markdown.md"))
	if err != nil {
		return nil, err
	}
	return markdown.ParseMarkdown(string(source)), nil
}

func controlMarkdown(c *ui.Context, s *controlState, _ string) {
	document, err := galleryMarkdown()
	if err != nil {
		components.Text(c, "Unable to load markdown.md").TextColor(c.Theme().TextMuted)
		return
	}
	bitmap, _ := galleryImage() // Missing images use Markdown's alt-text fallback.
	document.Render(c, markdown.MarkdownOptions{
		Images: map[string]*ui.Bitmap{"images/avatar.jpeg": bitmap},
		OnLink: func(string) { s.index++ },
		CodeBlocks: map[string]func(*ui.Context, markdown.MarkdownCodeBlock){
			"counter": func(c *ui.Context, block markdown.MarkdownCodeBlock) {
				ui.Column(c).Gap(c.Theme().Space(2)).Children(func() {
					components.Text(c, strings.TrimSuffix(block.Code, "\n")).Bold()
					if components.Button(c, "Increment", components.ButtonSecondary).Clicked() {
						s.index++
					}
					components.Text(c, fmt.Sprintf("Count: %d", s.index))
				})
			},
		},
	}).FillWidth().MaxWidth(640)
}
