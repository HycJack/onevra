package gallery

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// cardStackDestination is one card of the gallery's shortlist. Its text is
// kept as locale keys, so switching language updates the cards in place.
type cardStackDestination struct {
	id, key string
	icon    func() *ui.SVG
	tint    ui.Color
}

var cardStackDestinations = []cardStackDestination{
	{"lisbon", "destinationLisbon", lucide.TramFront, ui.RGB(226, 132, 72)},
	{"kyoto", "destinationKyoto", lucide.Landmark, ui.RGB(206, 86, 110)},
	{"reykjavik", "destinationReykjavik", lucide.Snowflake, ui.RGB(72, 150, 214)},
	{"cape-town", "destinationCapeTown", lucide.Mountain, ui.RGB(72, 168, 128)},
	{"marrakesh", "destinationMarrakesh", lucide.Sunset, ui.RGB(214, 150, 56)},
	{"queenstown", "destinationQueenstown", lucide.MountainSnow, ui.RGB(126, 112, 214)},
}

func controlCardStack(c *ui.Context, s *controlState, _ string) {
	t := c.Theme()
	items := make([]components.CardStackItem, len(cardStackDestinations))
	names := map[string]string{}
	for i, d := range cardStackDestinations {
		name, place, note, icon, tint := tr(c, d.key), tr(c, d.key+"Place"), tr(c, d.key+"Note"), d.icon, d.tint
		names[d.id] = name
		items[i] = components.CardStackItem{ID: d.id, Label: name, Content: func(c *ui.Context, radius float32) {
			ui.Column(c).Fill().Children(func() {
				// A tinted band carries the destination's mark above its text,
				// rounded to the card's top corners.
				ui.Box(c).FillWidth().Height(t.Space(36)).Center().Radius(radius, radius, 0, 0).
					LinearGradient(ui.LinearGradient{From: tint.Alpha(.32), To: tint.Alpha(.08), Angle: 135}).
					Children(func() {
						components.Icon(c, icon).Size(t.Space(14), t.Space(14)).TextColor(tint)
					})
				ui.Column(c).FillWidth().Padding(t.Space(5)).Gap(t.Space(1.5)).Children(func() {
					components.Text(c, name).FontSize(t.Rem(1.25)).FontWeight(600)
					components.Text(c, place).TextColor(t.TextMuted)
					components.Text(c, note).Wrap()
				})
			})
		}}
	}
	ui.Column(c).FillWidth().MinWidth(0).Gap(t.Space(5)).Children(func() {
		components.CardStack(c, &s.cards, components.CardStackOptions{
			Items: items, Label: tr(c, "cardStackDestinations"),
			LeftLabel: tr(c, "cardStackSkip"), RightLabel: tr(c, "cardStackShortlist"),
			LeftOutcome: tr(c, "cardStackSkipped"), RightOutcome: tr(c, "cardStackShortlisted"),
		})
		var kept []string
		for _, d := range s.cards.Decisions {
			if d.Decision == components.CardRight {
				kept = append(kept, names[d.ID])
			}
		}
		summary := tr(c, "cardStackShortlistEmpty")
		if len(kept) > 0 {
			sep := ", "
			if strings.HasPrefix(c.Locale(), "zh") {
				sep = "、"
			}
			summary = msg(c).CardStackShortlistS(strings.Join(kept, sep))
		}
		components.Text(c, summary).TextColor(t.TextMuted).TextAlign(ui.Center).FillWidth()
	})
}
