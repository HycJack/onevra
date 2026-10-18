package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// carouselDemoSlides are the destinations the carousel demonstration browses,
// by title and caption keys.
var carouselDemoSlides = []struct {
	title, caption string
	icon           func() *ui.SVG
	color          ui.Color
}{
	{"carouselMountains", "carouselMountainsCaption", lucide.Mountain, ui.RGB(88, 128, 196)},
	{"carouselForest", "carouselForestCaption", lucide.Trees, ui.RGB(72, 168, 128)},
	{"carouselDesert", "carouselDesertCaption", lucide.Sun, ui.RGB(226, 150, 60)},
	{"carouselCamp", "carouselCampCaption", lucide.TentTree, ui.RGB(196, 108, 148)},
	{"carouselHarbour", "carouselHarbourCaption", lucide.GalleryHorizontal, ui.RGB(66, 160, 190)},
}

// carouselDemoInterval is how long the rotating demonstration shows a slide.
const carouselDemoInterval = 5 * time.Second

func controlCarousel(c *ui.Context, s *controlState, variant string) {
	t := c.Theme()
	options := components.CarouselOptions{
		Label: tr(c, "carouselDestinations"),
		Count: len(carouselDemoSlides),
		Slide: func(c *ui.Context, i int) {
			slide := carouselDemoSlides[i]
			ui.Column(c).FillWidth().Height(220).Radius(t.Radius*2).Padding(t.Space(5)).Gap(t.Space(1)).
				Gradient(slide.color.Alpha(.32), slide.color.Alpha(.1), 135).Border(1, slide.color.Alpha(.35)).Children(func() {
				components.Icon(c, slide.icon).Size(t.Space(8), t.Space(8)).TextColor(slide.color)
				ui.Box(c).Grow(1)
				components.Text(c, tr(c, slide.title)).FontWeight(600).FontSize(t.FontSize * 1.25)
				components.Text(c, tr(c, slide.caption)).TextColor(t.TextMuted)
			})
		},
	}
	if variant == "Rotating" {
		options.Interval = carouselDemoInterval
	}
	// The demonstration fills the page width, as a carousel needs a definite one.
	components.Carousel(c, &s.carousel, options)
	components.Text(c, msg(c).SelectedS(tr(c, carouselDemoSlides[s.carousel.Index].title))).TextColor(t.TextMuted)
}
