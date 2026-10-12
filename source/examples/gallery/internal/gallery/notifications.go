package gallery

import (
	"errors"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

const notificationsPage galleryPage = dockPage + 1

type notificationDemo struct {
	current *mygo.Notification
	pending bool
	silent  bool
	status  string
	err     error
}

func (d *notificationDemo) close() {
	if d.current != nil {
		d.current.Close()
	}
	d.current = nil
	d.pending = false
}

func (g *Gallery) sendNotification() {
	d := &g.notification
	if g.Shell == nil || d.pending {
		return
	}
	d.err = nil
	if !mygo.NotificationsSupported() {
		d.status = "gallery.notification_unsupported"
		return
	}
	d.close()
	n := mygo.NewNotification(mygo.NotificationOptions{
		Title:  g.text("gallery.notification_title"),
		Body:   g.text("gallery.notification_body"),
		Silent: d.silent,
	})
	d.current, d.pending = n, true
	d.status = "gallery.notification_sending"
	n.OnClick(func() {
		if d.current != n {
			return
		}
		d.close()
		d.status = "gallery.notification_clicked"
		g.page = notificationsPage
		g.settings.open = false
		g.showFromDock()
		g.Shell.Window().Invalidate()
	})
	// Show may wait for the macOS permission prompt. Keep the UI responsive;
	// all gallery state and completion handling remain on the UI thread.
	go func() {
		err := n.Show()
		mygo.RunOnMain(func() {
			if d.current != n {
				// Withdrawal or Gallery shutdown may precede delivery.
				n.Close()
				return
			}
			d.pending = false
			switch {
			case errors.Is(err, mygo.ErrNotificationsDenied):
				d.close()
				d.status = "gallery.notification_denied"
			case err != nil:
				d.close()
				d.status, d.err = "gallery.notification_failed", err
			default:
				d.status = "gallery.notification_sent"
			}
			g.Shell.Window().Invalidate()
		})
	}()
}

func (g *Gallery) Notifications(c *ui.Context) {
	g.withLanguage(c, func() {
		d := &g.notification
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
			ui.Column(c).FillWidth().Padding(c.Theme().Space(6)).Gap(16).Children(func() {
				components.Text(c, tr(c, "gallery.system_notifications")).FontSize(22).Bold()
				components.Text(c, tr(c, "gallery.notification_description")).Wrap()
				components.Text(c, tr(c, "gallery.notification_requirements")).Wrap().TextColor(c.Theme().TextMuted)
				components.Checkbox(c, &d.silent, tr(c, "gallery.notification_silent"))
				components.Button(c, tr(c, "gallery.notification_send"), components.ButtonPrimary).
					Disabled(g.Shell == nil || d.pending).OnClick(g.sendNotification)
				components.Button(c, tr(c, "gallery.notification_withdraw"), components.ButtonOutline).
					Disabled(d.current == nil).OnClick(func() {
					d.close()
					d.err = nil
					d.status = "gallery.notification_withdrawn"
				})
				if d.status != "" {
					components.Text(c, tr(c, d.status)).Wrap()
				}
				if d.err != nil {
					components.Text(c, d.err.Error()).Wrap().TextColor(c.Theme().Danger)
				}
			})
		})
	})
}
