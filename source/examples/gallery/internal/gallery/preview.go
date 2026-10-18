package gallery

import (
	"image"
	"image/color"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

// PreviewEntry describes one component page the website can show. The
// standalone Gallery and the website preview share the same demonstrations.
type PreviewEntry struct {
	Slug string `json:"slug"`
	// Key is the stable Gallery page key, such as controls/Button.
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	NameZh     string   `json:"nameZh"`
	Group      string   `json:"group"`
	GroupZh    string   `json:"groupZh"`
	Variants   []string `json:"variants,omitempty"`
	Components []string `json:"components"`

	page  galleryPage
	build func(*Preview, *ui.Context)
}

// Preview shows one component page without the Gallery shell. The website
// builds it for browsers, where the page's own site provides navigation.
// It creates no window and no native resources.
type Preview struct {
	g     Gallery
	page  galleryPage
	title string
	build func(*Preview, *ui.Context)
	// These open the Gallery's own surfaces on the first frame only, so the
	// background can reopen them after they close.
	settingsOpened, commandOpened bool
}

// PreviewEntries lists the component pages in their Gallery order.
func PreviewEntries() []PreviewEntry {
	en := newTranslator("en", "en")
	zh := newTranslator("zh-Hans", "zh-Hans")
	entries := make([]PreviewEntry, 0, len(controlCatalog))
	for i, example := range controlCatalog {
		components := pageComponents(example)
		if len(components) == 0 {
			components = []string{example.name}
		}
		entries = append(entries, PreviewEntry{
			Slug:       previewSlug(example.name),
			Key:        "controls/" + example.name,
			Name:       en.text(example.name),
			NameZh:     zh.text(example.name),
			Group:      en.text(example.group),
			GroupZh:    zh.text(example.group),
			Variants:   example.variants,
			Components: components,
			page:       controlPageStart() + galleryPage(i),
		})
	}
	// These public surfaces live in the Gallery's own chrome rather than a
	// control page, so the website shows their existing demonstrations directly.
	entries = append(entries,
		PreviewEntry{
			Slug: "icon", Key: "application/icon",
			Name: "Icon", NameZh: "Icon",
			Group: "Actions & navigation", GroupZh: zh.text("actionsNavigation"),
			Components: []string{"Icon"},
			build:      previewIcon,
		},
		PreviewEntry{
			Slug: "pixel-grid", Key: "application/pixel-grid",
			Name: "PixelGrid", NameZh: "PixelGrid",
			Group: "Layout & containers", GroupZh: zh.text("layoutContainers"),
			Components: []string{"PixelGrid"},
			build:      previewPixelGrid,
		},
		PreviewEntry{
			Slug: "settings", Key: "application/settings",
			Name: "Settings", NameZh: zh.text("settings"),
			Group: "Feedback & overlays", GroupZh: zh.text("feedbackOverlays"),
			Components: []string{"Settings", "SearchSettings"},
			build:      previewSettings,
		},
		PreviewEntry{
			Slug: "command-palette", Key: "application/command",
			Name: "Command palette", NameZh: zh.text("commandPalette"),
			Group: "Feedback & overlays", GroupZh: zh.text("feedbackOverlays"),
			Components: []string{"CommandView", "CommandDialog"},
			build:      previewCommand,
		},
	)
	return entries
}

// NewPreview prepares a page for the given language. An unknown slug shows
// the first page.
func NewPreview(slug, language string) *Preview {
	p := &Preview{}
	if language != "zh-Hans" {
		language = "en"
	}
	p.g.config.Language = language
	p.g.language = newTranslator(language, "")
	previewImage = previewAvatar
	if !p.Show(slug) {
		p.selectEntry(PreviewEntries()[0])
	}
	return p
}

// Show selects a page by slug or by the name of a component it demonstrates,
// and reports whether it exists.
func (p *Preview) Show(slug string) bool {
	slug = strings.ToLower(slug)
	for _, entry := range PreviewEntries() {
		if entry.Slug == slug {
			p.selectEntry(entry)
			return true
		}
		for _, component := range entry.Components {
			if previewSlug(component) == slug {
				p.selectEntry(entry)
				return true
			}
		}
	}
	return false
}

func (p *Preview) selectEntry(entry PreviewEntry) {
	p.page = entry.page
	p.build = entry.build
	p.title = p.g.text(entry.Name)
}

// Title is the selected page's name in the preview's language.
func (p *Preview) Title() string { return p.title }

// View builds the selected page. The website's frame supplies the theme and
// the space, so the preview only draws the demonstration itself.
func (p *Preview) View(c *ui.Context) {
	theme.Apply(c)
	ui.Box(c).Fill().MinWidth(0).MinHeight(0).Background(c.Theme().Background).Children(func() {
		if p.build != nil {
			p.build(p, c)
			return
		}
		if example, ok := controlExampleFor(p.page); ok {
			p.g.Controls(c, example)
			return
		}
		ui.Column(c).Fill().Center().Padding(c.Theme().Space(6)).Children(func() {
			ui.Text(c, "This component preview is unavailable.")
		})
	})
	// The Gallery's Shell hosts toasts; the website's frame has no Shell.
	components.Toasts(c, p.g.toasts)
}

// previewIcon reuses the Gallery's icon samples to show sizes and inherited
// colors without a full-library registry.
func previewIcon(_ *Preview, c *ui.Context) {
	t := c.Theme()
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
		ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(4)).Children(func() {
			iconSamples(c)
		})
	})
}

// previewPixelGrid shows the themed grid surface behind sample content.
func previewPixelGrid(_ *Preview, c *ui.Context) {
	t := c.Theme()
	components.PixelGrid(c, components.PixelGridOptions{Spacing: 22, MajorEvery: 5}).Fill().Children(func() {
		ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
			components.Text(c, "PixelGrid").FontSize(t.Rem(1.8)).Bold()
			components.Text(c, tr(c, "layoutContainers")).TextColor(t.TextMuted)
		})
	})
}

// previewSettings opens the Gallery's Settings dialog over a page that can
// reopen it after closing.
func previewSettings(p *Preview, c *ui.Context) {
	if !p.settingsOpened {
		p.g.settings.open = true
		p.settingsOpened = true
	}
	t := c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
		components.Text(c, tr(c, "settings")).FontSize(t.Rem(1.4)).Bold()
		if components.Button(c, tr(c, "openSettings"), components.ButtonPrimary).Clicked() {
			p.g.settings.open = true
		}
	})
	p.g.Settings(c)
}

// previewCommand opens the Gallery's command palette over a page that can
// reopen it after closing.
func previewCommand(p *Preview, c *ui.Context) {
	if !p.commandOpened {
		p.g.openCommands()
		p.commandOpened = true
	}
	t := c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
		components.Text(c, tr(c, "commandPalette")).FontSize(t.Rem(1.4)).Bold()
		if components.Button(c, tr(c, "open"), components.ButtonPrimary).Clicked() {
			p.g.openCommands()
		}
	})
	p.g.Commands(c)
	// The palette's Settings command opens this surface too.
	p.g.Settings(c)
}

// previewImage replaces the resource-backed Image demonstration in the
// browser, where a page has no bundle resource directory.
var previewImage func() (*ui.Bitmap, error)

// previewAvatar paints a small warm image so every Image fit stays visible
// without file or network access.
func previewAvatar() (*ui.Bitmap, error) {
	const width, height = 360, 240
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			t := float64(x+y) / float64(width+height)
			r := uint8(240 - 55*t)
			g := uint8(163 - 45*t)
			b := uint8(126 - 38*t)
			dx, dy := float64(x-width/2), float64(y-height/2)
			if dx*dx+dy*dy < 64*64 {
				r, g, b = 252, 241, 230
			}
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return ui.NewBitmap(img), nil
}

// previewSlug turns a page or component name into a URL fragment, as Button
// is button and IconButton is icon-button.
func previewSlug(name string) string {
	var b strings.Builder
	runes := []rune(name)
	dash := false
	for i, r := range runes {
		upper := r >= 'A' && r <= 'Z'
		if upper && i > 0 {
			previous := runes[i-1]
			nextLower := i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z'
			if previous >= 'a' && previous <= 'z' || previous >= '0' && previous <= '9' ||
				previous >= 'A' && previous <= 'Z' && nextLower {
				if !dash {
					b.WriteByte('-')
					dash = true
				}
			}
		}
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		case upper:
			b.WriteRune(r + ('a' - 'A'))
			dash = false
		case !dash:
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
