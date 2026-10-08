package gallery

import (
	"sort"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

const componentStatusPage galleryPage = 2

type componentStatusCard struct {
	Name     string
	Replaces []string
}

// Replacement variants share a card under their owning Onevra component.
func componentStatusCards(entries []mygoReviewEntry, status mygoStatus) []componentStatusCard {
	grouped := map[string][]string{}
	for _, entry := range entries {
		if entry.Status != status {
			continue
		}
		name := entry.Name
		if status == mygoReplaced {
			owner, _, _ := strings.Cut(entry.Replacement, "/")
			name = strings.TrimPrefix(strings.TrimSpace(owner), "components.")
			grouped[name] = append(grouped[name], entry.Name)
		} else {
			grouped[name] = nil
		}
	}
	cards := make([]componentStatusCard, 0, len(grouped))
	for name, replaces := range grouped {
		sort.Strings(replaces)
		cards = append(cards, componentStatusCard{Name: name, Replaces: replaces})
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Name < cards[j].Name })
	return cards
}

// ComponentStatus is a read-only name grid backed by the startup review decisions.
func (g *Gallery) ComponentStatus(c *ui.Context) {
	entries := mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled)
	t := c.Theme()
	components.ScrollView(c, components.ScrollViewOptions{}).Key("component-status").Fill().Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			for _, status := range []mygoStatus{mygoExcluded, mygoReplaced, mygoUpstream, mygoUnreviewed, mygoPrimitive} {
				cards := componentStatusCards(entries, status)
				if len(cards) == 0 {
					continue
				}
				ui.Column(c).Key(string(status)).FillWidth().MinWidth(0).Gap(t.Space(3)).Children(func() {
					components.Text(c, string(status)).Bold()
					// Fixed-width cells wrap against their actual parent width, including
					// sidebar resizing, rather than assuming the context is the content area.
					ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(3)).Children(func() {
						height := t.Rem(3)
						if status == mygoReplaced {
							height = t.Rem(6)
						}
						for _, card := range cards {
							ui.Column(c).Key(card.Name).Width(t.Rem(12)).MaxWidthPercent(100).MinWidth(0).MinHeight(height).Padding(t.Space(3)).Gap(t.Space(1)).Radius(t.Radius).Background(t.Surface).Children(func() {
								components.Text(c, card.Name).Bold().FillWidth().MinWidth(0)
								if len(card.Replaces) > 0 {
									components.Text(c, strings.Join(card.Replaces, ", ")).FillWidth().MinWidth(0).TextColor(t.TextMuted)
								}
							})
						}
					})
				})
			}
		})
	})
}
