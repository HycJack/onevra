package gallery

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestNotificationsHeadless(t *testing.T) {
	for _, language := range []string{"en", "zh-Hans"} {
		g := &Gallery{page: notificationsPage, language: newTranslator(language, "en")}
		tt := ui.NewTester(g.View, 360, 600)
		if !tt.HasText(g.text("gallery.system_notifications")) {
			t.Fatal("missing localized notification page", tt.Texts())
		}
		_ = tt.Click(g.text("gallery.notification_send"))
		g.sendNotification()
		if g.notification.current != nil || g.notification.pending || g.notification.status != "" {
			t.Fatal("headless rendering and actions must not send native notifications")
		}
		g.Close()
		g.Close()
	}
	if pageFromKey(pageKey(notificationsPage)) != notificationsPage {
		t.Fatal("notification route does not round trip")
	}
}

// Capture the actual desktop view without posting unsolicited OS notifications.
// Delivery and permission prompts need a packaged app and manual review.
func TestNativeNotifications(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	g := &Gallery{page: notificationsPage}
	mygo.RunOnMain(func() { g.newWindow(nativeBackground()); g.enterGallery() })
	win := g.Shell.Window()
	defer win.Destroy()
	for _, tc := range []struct {
		name, language string
		theme          mygo.ThemeSource
		width          int
	}{
		{"notifications-light", "en", mygo.ThemeLight, 820},
		{"notifications-dark", "zh-Hans", mygo.ThemeDark, 820},
		{"notifications-narrow", "zh-Hans", mygo.ThemeLight, 360},
	} {
		mygo.RunOnMain(func() {
			g.chooseLanguage(tc.language)
			mygo.Theme.SetSource(tc.theme)
			win.SetMinimumSize(360, 480)
			win.SetSize(tc.width, 600)
			win.Invalidate()
		})
		time.Sleep(400 * time.Millisecond)
		// Capture the native content surface.
		data, err := win.CapturePage()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, tc.name+".png"), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
