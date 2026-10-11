package gallery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/markdown"
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
		components.Text(c, tr(c, "gallery.markdown_load_error")).TextColor(c.Theme().TextMuted)
		return
	}
	bitmap, _ := galleryImage() // Missing images use Markdown's alt-text fallback.
	document.Render(c, markdown.MarkdownOptions{
		CopyCodeLabel: tr(c, "gallery.copy_code"),
		Images:        map[string]*ui.Bitmap{"images/avatar.jpeg": bitmap},
		OnLink:        func(string) { s.index++ },
		CodeBlocks: map[string]func(*ui.Context, markdown.MarkdownCodeBlock){
			"counter": func(c *ui.Context, block markdown.MarkdownCodeBlock) {
				ui.Column(c).Gap(c.Theme().Space(2)).Children(func() {
					components.Text(c, strings.TrimSuffix(block.Code, "\n")).Bold()
					if components.Button(c, tr(c, "gallery.increment"), components.ButtonSecondary).Clicked() {
						s.index++
					}
					components.Text(c, fmt.Sprintf(tr(c, "gallery.count_d"), s.index))
				})
			},
		},
	}).FillWidth().MaxWidth(640)
}
