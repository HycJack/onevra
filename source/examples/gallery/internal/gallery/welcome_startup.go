package gallery

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// Keep the application version aligned with packer.yaml; a test checks both.
const galleryVersion = "0.2026.7"

// welcomeCounts summarizes the catalog: component pages and their sections.
func welcomeCounts() (pages, sections int) {
	for _, group := range navigationCatalog() {
		if slices.Contains(navigationGroups, group.label) {
			pages += len(group.entries)
			sections++
		}
	}
	return pages, sections
}

// WelcomePage is a standalone startup entrance, not a component demonstration.
func (g *Gallery) WelcomePage(c *ui.Context) {
	g.withLanguage(c, func() { g.welcomePageLocalized(c) })
}

func (g *Gallery) welcomePageLocalized(c *ui.Context) {
	t := c.Theme()
	width, height := c.Size()
	header := sidebar.HeaderHeight(c)
	root := ui.Column(c).Fill().Padding(header, 0, 0, 0)
	welcomeBackdrop(c, root,
		welcomeGlow{x: .05, y: .1, radius: .3, color: t.Accent, alpha: .22},
		welcomeGlow{x: .78, y: .45, radius: .3, color: brandSky, alpha: .32},
		welcomeGlow{x: .45, y: 1, radius: .28, color: brandMint, alpha: .28},
	).Children(func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
			ui.Column(c).FillWidth().MinHeight(height-header).
				Padding(t.Space(4), t.Space(8), t.Space(6), t.Space(8)).Gap(t.Space(6)).Children(func() {
				ui.Row(c).FillWidth().Grow(1).MaxWidth(t.Space(230)).AlignSelf(ui.Center).
					AlignItems(ui.Center).Justify(ui.Center).Gap(t.Space(12)).Children(func() {
					g.welcomeIntro(c, width < 760)
					if width >= 760 {
						welcomeShowcase(c)
					}
				})
				welcomeEnter(c, ui.Row(c).FillWidth().Wrap().Justify(ui.Center).GapX(t.Space(3)).GapY(t.Space(1)), 7).Children(func() {
					components.Text(c, fmt.Sprintf(tr(c, "gallery.version_s"), galleryVersion)).FontSize(t.Rem(.75)).TextColor(t.TextMuted)
					components.Text(c, tr(c, "gallery.built_with_onevra_powered_by_mygo_written_in_go")).
						FontSize(t.Rem(.75)).TextColor(t.TextMuted).TextAlign(ui.Center).Wrap()
				})
			})
		})
	})
	g.Commands(c)
}

func (g *Gallery) welcomeIntro(c *ui.Context, centered bool) {
	t := c.Theme()
	align, textAlign := ui.Start, ui.Start
	if centered {
		align, textAlign = ui.Center, ui.Center
	}
	pages, _ := welcomeCounts()
	ui.Column(c).Basis(t.Space(90)).Grow(1).MinWidth(0).MaxWidth(t.Space(112)).
		AlignItems(align).Gap(t.Space(5)).Children(func() {
		welcomeEnter(c, ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3)), 0).Children(func() {
			welcomeBrandMark(c, t.Space(12))
			ui.Column(c).Gap(t.Space(.5)).Children(func() {
				components.Text(c, "Onevra Gallery").Bold().SingleLine()
				components.Text(c, tr(c, "gallery.welcome_tagline")).FontSize(t.Rem(.8)).TextColor(t.TextMuted).SingleLine()
			})
		})
		welcomeEnter(c, ui.Column(c).FillWidth().AlignItems(align).Gap(t.Space(3)), 1).Children(func() {
			components.Text(c, tr(c, "gallery.welcome_to_onevra")).FontSize(t.Rem(2.6)).Bold().
				LineHeight(1.1).LetterSpacing(-.5).TextAlign(textAlign).Wrap()
			components.Text(c, tr(c, "gallery.welcome_entrance_summary")).FontSize(t.Rem(1.05)).
				TextColor(t.TextMuted).LineHeight(1.45).TextAlign(textAlign).Wrap()
		})
		welcomeEnter(c, ui.Column(c).AlignItems(ui.Start).Gap(t.Space(2.5)), 2).Children(func() {
			for _, feature := range []struct {
				icon  func() *ui.SVG
				color ui.Color
				text  string
			}{
				{lucide.LayoutDashboard, t.Accent, fmt.Sprintf(tr(c, "gallery.welcome_feature_pages_d"), pages)},
				{lucide.Shapes, brandSky, fmt.Sprintf(tr(c, "gallery.welcome_feature_icons_d"), len(iconCatalog))},
				{lucide.SunMoon, t.Success, tr(c, "gallery.welcome_feature_themes")},
			} {
				ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
					welcomeIconTile(c, feature.icon, feature.color, t.Space(7))
					components.Text(c, feature.text).FontSize(t.Rem(.9)).Wrap()
				})
			}
		})
		welcomeEnter(c, ui.Row(c).Wrap().AlignItems(ui.Center).Justify(align).GapX(t.Space(4)).GapY(t.Space(3)), 3).Children(func() {
			enter := welcomeAction(c, "gallery.enter_gallery", components.ButtonPrimary, lucide.ArrowRight).
				Padding(t.Space(2.5), t.Space(5)).Radius(t.Space(3))
			lift := ui.Pose{}
			if enter.Hovered() {
				lift = ui.Pose{Y: -1, Scale: 1.02}
			}
			enter.Shadow(0, t.Space(1.5), t.Space(4), -t.Space(1), t.Accent.Alpha(.35)).AnimatePose(lift, ui.Snappy)
			if enter.Clicked() {
				g.enterGallery()
			}
			_, shortcut := commandShortcut()
			ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
				welcomeKeycap(c, shortcut)
				components.Text(c, tr(c, "gallery.welcome_search_hint")).FontSize(t.Rem(.8)).TextColor(t.TextMuted)
			})
		})
	})
}

// welcomeShowcase is a decorative collage of interface fragments: a miniature
// Shell window with floating cards, drawn from boxes rather than live controls.
func welcomeShowcase(c *ui.Context) {
	t := c.Theme()
	card := func(e ui.Element) ui.Element {
		return e.Background(t.Surface).Border(1, t.Border).
			Shadow(0, t.Space(3), t.Space(8), -t.Space(2), t.Text.Alpha(.16))
	}
	line := func(width, alpha float32) {
		ui.Box(c).Width(width).Height(t.Space(1.5)).Radius(t.Space(1)).Background(t.Text.Alpha(alpha))
	}
	ui.Box(c).Size(t.Space(92), t.Space(84)).Shrink(0).Unselectable().
		Draw(func(p *ui.Painter, r ui.Rect) {
			welcomeOrbits(p, r, r.W*.55, r.H*.48, t.Space(16), 4, t.Accent)
		}).Children(func() {
		// The window.
		welcomeEnter(c, card(ui.Column(c).Absolute().Left(t.Space(4)).Top(t.Space(10)).
			Size(t.Space(76), t.Space(56)).Radius(t.Space(4)).Clip()), 2).Children(func() {
			ui.Row(c).FillWidth().Height(t.Space(7)).Padding(0, t.Space(3)).AlignItems(ui.Center).Gap(t.Space(1.5)).
				Background(t.Text.Alpha(.03)).Children(func() {
				for _, dot := range []ui.Color{ui.RGB(0xff, 0x5f, 0x57), ui.RGB(0xfe, 0xbc, 0x2e), ui.RGB(0x28, 0xc8, 0x40)} {
					ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(t.Space(1)).Background(dot.Alpha(.85))
				}
			})
			ui.Row(c).FillWidth().Grow(1).Children(func() {
				ui.Column(c).Width(t.Space(20)).FillHeight().Padding(t.Space(2.5)).Gap(t.Space(2.5)).
					Background(t.Text.Alpha(.025)).Children(func() {
					ui.Box(c).FillWidth().Height(t.Space(4)).Radius(t.Space(1.5)).Background(t.Accent.Alpha(.16))
					for _, w := range []float32{11, 9, 12, 8, 10} {
						line(t.Space(w), .1)
					}
				})
				ui.Column(c).Grow(1).Padding(t.Space(3.5)).Gap(t.Space(3)).Children(func() {
					line(t.Space(22), .2)
					line(t.Space(34), .08)
					ui.Row(c).FillWidth().Height(t.Space(24)).AlignItems(ui.End).Gap(t.Space(1.5)).Padding(t.Space(2), 0, 0, 0).Children(func() {
						for i, h := range []float32{.45, .7, .55, .9, .62, .8, .4} {
							color := brandSky
							if i == 3 {
								color = t.Accent
							}
							ui.Box(c).Grow(1).HeightPercent(h * 100).Radius(t.Space(1)).Background(color.Alpha(.75))
						}
					})
				})
			})
		})
		// A floating switch.
		welcomeEnter(c, card(ui.Row(c).Absolute().Right(0).Top(t.Space(3)).
			Padding(t.Space(2.5), t.Space(3)).Gap(t.Space(2.5)).AlignItems(ui.Center).Radius(t.Space(3))).Rotate(4), 4).Children(func() {
			welcomeIconTile(c, lucide.Moon, brandSky, t.Space(6))
			ui.Box(c).Size(t.Space(9), t.Space(5)).Padding(t.Space(.5)).Radius(t.Space(2.5)).
				Background(t.Accent).Justify(ui.End).Row().Children(func() {
				ui.Box(c).Size(t.Space(4), t.Space(4)).Radius(t.Space(2)).Background(ui.RGB(0xff, 0xff, 0xff))
			})
		})
		// A floating progress card.
		welcomeEnter(c, card(ui.Row(c).Absolute().Left(0).Bottom(t.Space(2)).
			Padding(t.Space(3)).Gap(t.Space(3)).AlignItems(ui.Center).Radius(t.Space(3.5))).Rotate(-3), 5).Children(func() {
			components.ProgressRing(c, .72).Size(t.Space(9), t.Space(9)).Label(tr(c, "gallery.welcome_showcase_progress"))
			ui.Column(c).Gap(t.Space(1.5)).Children(func() {
				line(t.Space(14), .2)
				line(t.Space(10), .08)
			})
		})
		// A floating palette.
		welcomeEnter(c, card(ui.Row(c).Absolute().Right(t.Space(2)).Bottom(0).
			Padding(t.Space(2.5)).Gap(t.Space(1.5)).Radius(t.Space(3))).Rotate(2), 6).Children(func() {
			for _, color := range []ui.Color{t.Accent, brandDeep, brandSky, brandMint, t.Warning} {
				ui.Box(c).Size(t.Space(5), t.Space(5)).Radius(t.Space(2.5)).Background(color).Border(2, t.Surface)
			}
		})
	})
}

func (g *Gallery) enterGallery() {
	if g.Shell != nil {
		g.Shell.SetPage(nil)
	}
}

// welcomeAction is a button with a trailing icon; label is a catalog ID.
func welcomeAction(c *ui.Context, label string, variant components.ButtonVariant, icon func() *ui.SVG) ui.Element {
	b := components.Button(c, "", variant).Label(tr(c, label)).Gap(c.Theme().Space(2))
	b.Children(func() {
		components.Text(c, tr(c, label)).Wrap()
		components.Icon(c, icon).Size(c.Theme().Space(4), c.Theme().Space(4)).Shrink(0)
	})
	return b
}
