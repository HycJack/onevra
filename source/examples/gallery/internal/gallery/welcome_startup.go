package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// Keep the application version aligned with packer.yaml; a test checks both.
const galleryVersion = "0.2026.1"

// WelcomePage is a standalone startup entrance, not a component demonstration.
func (g *Gallery) WelcomePage(c *ui.Context) {
	g.withLanguage(c, func() { g.welcomePageLocalized(c) })
}

func (g *Gallery) welcomePageLocalized(c *ui.Context) {
	t := c.Theme()
	_, height := c.Size()
	ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).
		Gradient(t.Background.Mix(t.Accent, .06), t.Background, 140).Children(func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
			ui.Column(c).FillWidth().MinHeight(height - sidebar.HeaderHeight(c)).
				Padding(t.Space(8)).Center().Children(func() {
				ui.Column(c).FillWidth().MaxWidth(t.Space(100)).AlignItems(ui.Center).Gap(t.Space(5)).Children(func() {
					ui.Box(c).Size(t.Space(16), t.Space(16)).Radius(t.Space(5)).Center().
						Background(t.Accent.Alpha(.1)).TextColor(t.Accent).Children(func() {
						components.Icon(c, lucide.Sparkles).Size(t.Space(8), t.Space(8))
					})
					components.Text(c, tr(c, "gallery.welcome_to_onevra")).FontSize(t.Rem(2)).Bold().TextAlign(ui.Center).Wrap()
					components.Text(c, tr(c, "gallery.a_collection_of_native_components_ready_to_explore")).TextColor(t.TextMuted).TextAlign(ui.Center).Wrap()
					if components.Button(c, tr(c, "gallery.enter_gallery"), components.ButtonPrimary).Clicked() {
						g.enterGallery()
					}
					components.Text(c, fmt.Sprintf(tr(c, "gallery.version_s"), galleryVersion)).FontSize(t.Rem(.8)).TextColor(t.TextMuted)
				})
			})
		})
	})
	g.Commands(c)
}

func (g *Gallery) enterGallery() {
	if g.Shell != nil {
		g.Shell.SetPage(nil)
	}
}
