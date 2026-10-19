package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// badgeStatuses are the steps the live Badge example cycles through.
var badgeStatuses = []struct {
	key  string
	tone components.BadgeTone
	icon func() *ui.SVG
}{
	{"badgeQueued", components.BadgeNeutral, lucide.Clock},
	{"badgeSyncing", components.BadgeAccent, lucide.RefreshCw},
	{"badgeSynced", components.BadgeSuccess, lucide.Check},
	{"badgeSyncFailed", components.BadgeDanger, lucide.CircleX},
}

func controlBadge(c *ui.Context, s *controlState, variant string) {
	switch variant {
	case "Tones":
		ui.Row(c).Gap(8).Wrap().Children(func() {
			components.Badge(c, tr(c, "badgeNeutral"))
			components.Badge(c, tr(c, "accent"), components.BadgeOptions{Tone: components.BadgeAccent})
			components.Badge(c, tr(c, "badgeSuccess"), components.BadgeOptions{Tone: components.BadgeSuccess})
			components.Badge(c, tr(c, "badgeWarning"), components.BadgeOptions{Tone: components.BadgeWarning})
			components.Badge(c, tr(c, "badgeDanger"), components.BadgeOptions{Tone: components.BadgeDanger})
		})
	case "Icons":
		ui.Row(c).Gap(8).Wrap().Children(func() {
			components.Badge(c, tr(c, "new"), components.BadgeOptions{Icon: lucide.Sparkles()})
			components.Badge(c, tr(c, "badgeVerified"), components.BadgeOptions{Tone: components.BadgeSuccess, Icon: lucide.BadgeCheck()})
			components.Badge(c, tr(c, "badgeBeta"), components.BadgeOptions{Tone: components.BadgeAccent, Icon: lucide.FlaskConical()})
			components.Badge(c, tr(c, "badgeOffline"), components.BadgeOptions{Tone: components.BadgeWarning, Icon: lucide.WifiOff()})
		})
	case "Live updates":
		// Keeping each badge built in place lets a change morph it.
		status := badgeStatuses[s.badge.status%len(badgeStatuses)]
		ui.Row(c).Gap(8).AlignItems(ui.Center).Wrap().Children(func() {
			components.Badge(c, tr(c, status.key), components.BadgeOptions{Tone: status.tone, Icon: status.icon()})
			if components.Button(c, tr(c, "badgeNextStatus"), components.ButtonSecondary).Clicked() {
				s.badge.status++
			}
		})
		ui.Row(c).Gap(8).AlignItems(ui.Center).Wrap().Children(func() {
			tone := components.BadgeNeutral
			if s.badge.unread > 0 {
				tone = components.BadgeAccent
			}
			components.Badge(c, msg(c).BadgeUnread(s.badge.unread), components.BadgeOptions{Tone: tone, Icon: lucide.Mail()})
			if components.Button(c, tr(c, "badgeAddUnread"), components.ButtonSecondary).Clicked() {
				s.badge.unread = s.badge.unread*3 + 1
			}
		})
	}
}

func controlRichText(c *ui.Context, s *controlState, _ string) {
	components.RichText(c, ui.Span{Text: tr(c, "styled"), Weight: 600}).Children(func() {
		components.Text(c, tr(c, "textAndAn"))
		if components.Button(c, tr(c, "inlineAction"), components.ButtonLink).Clicked() {
			s.index++
		}
		ui.Text(c, msg(c).WrapAsOneParagraphClicksD(s.index))
	}).MaxWidth(320)
}

func controlBreadcrumbs(c *ui.Context, s *controlState, variant string) {
	if variant == "Long path" {
		// Twenty levels collapse their middle into an ellipsis menu.
		path := []string{tr(c, "home")}
		for i := 1; i < 20; i++ {
			path = append(path, msg(c).LongPathFolder(i))
		}
		b := components.Breadcrumbs(c, path, &s.index)
		b.Label(tr(c, "path"))
		if b.Changed() {
			s.text = path[s.index]
		}
		if s.text != "" {
			components.Text(c, s.text).TextColor(c.Theme().TextMuted)
		}
		return
	}
	path := []string{tr(c, "home"), tr(c, "components"), "Breadcrumbs"}
	b := components.Breadcrumbs(c, path, &s.index)
	b.Label(tr(c, "path"))
	if b.Changed() {
		s.text = []string{"Home", "Components", "Breadcrumbs"}[s.index]
	}
	if s.text != "" {
		components.Text(c, tr(c, s.text)).TextColor(c.Theme().TextMuted)
	}
}

func controlRating(c *ui.Context, s *controlState, _ string) {
	components.Rating(c, &s.index, 5).Label(tr(c, "rating"))
	ui.Text(c, msg(c).RatingOutOfFive(s.index)).TextColor(c.Theme().TextMuted)
}
