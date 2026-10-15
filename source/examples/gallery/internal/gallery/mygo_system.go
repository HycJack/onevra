package gallery

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery/locales"
)

// systemDemo holds the system-service examples. Native is false in headless
// views, which show the controls without calling the running app; results of
// background work return to the UI thread through the gallery window.
type systemDemo struct {
	native bool
	window func() *mygo.Window

	network    mygo.NetworkStatus
	networkOff func()

	secretKey, secretValue, secretStatus string

	recentStatus string

	watchPath   string
	watchEvents []mygo.FileEvent
	watchErr    string
	watchStop   context.CancelFunc

	apps        []mygo.InstalledApp
	appsLoading bool
	appsErr     string
	appQuery    string
	icons       map[string]*ui.Bitmap
	panel       *mygo.Window

	fontQuery string
	monospace bool

	cookies      []mygo.Cookie
	cookieStatus string
}

// update applies fn on the UI thread and redraws the gallery.
func (d *systemDemo) update(fn func()) {
	if w := d.window(); w != nil {
		w.Update(fn)
	}
}

// background runs work off the UI thread, as dialogs and system stores block.
func (d *systemDemo) background(work func() func()) {
	if !d.native {
		return
	}
	go func() {
		apply := work()
		d.update(apply)
	}()
}

func (d *systemDemo) close() {
	if d.networkOff != nil {
		d.networkOff()
		d.networkOff = nil
	}
	d.stopWatching()
	if d.panel != nil {
		d.panel.Close()
		d.panel = nil
	}
}

func systemRow(c *ui.Context, label, value string) {
	if value == "" {
		value = "—"
	}
	ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
		components.Text(c, label).TextColor(c.Theme().TextMuted).Width(160)
		components.Text(c, value).Grow(1).Wrap()
	})
}

func yesNo(m locales.Messages, v bool) string {
	if v {
		return m.SystemYes()
	}
	return m.SystemNo()
}

func exampleNetwork(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	if d.native && d.networkOff == nil {
		d.network = mygo.Network.Status()
		d.networkOff = mygo.Network.OnChanged(func(status mygo.NetworkStatus) {
			d.update(func() { d.network = status })
		})
	}
	status := m.NetworkOffline()
	if d.network.Online {
		status = m.NetworkOnline()
	}
	components.Badge(c, status)
	systemRow(c, m.NetworkMetered(), yesNo(m, d.network.Metered))
	components.Text(c, m.NetworkHint()).Wrap().TextColor(c.Theme().TextMuted)
}

func exampleSecrets(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	if d.secretKey == "" {
		d.secretKey = "gallery-demo"
	}
	components.TextBox(c, &d.secretKey, components.TextBoxOptions{}).Label(m.SecretsKey()).FillWidth()
	components.TextBox(c, &d.secretValue, components.TextBoxOptions{}).Label(m.SecretsValue()).FillWidth()
	key, value := d.secretKey, d.secretValue
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, m.SecretsSave(), components.ButtonPrimary).Disabled(key == "").Clicked() {
			d.background(func() func() {
				err := mygo.Secrets.Set(key, []byte(value))
				return func() { d.secretStatus = secretResult(m, err, m.SecretsSaved()) }
			})
		}
		if components.Button(c, m.SecretsRead(), components.ButtonOutline).Disabled(key == "").Clicked() {
			d.background(func() func() {
				got, err := mygo.Secrets.Get(key)
				return func() {
					if err == nil {
						d.secretValue = string(got)
					}
					d.secretStatus = secretResult(m, err, m.SecretsLoaded())
				}
			})
		}
		if components.Button(c, m.Delete(), components.ButtonOutline).Disabled(key == "").Clicked() {
			d.background(func() func() {
				err := mygo.Secrets.Delete(key)
				return func() { d.secretStatus = secretResult(m, err, m.SecretsDeleted()) }
			})
		}
	})
	if d.secretStatus != "" {
		components.Text(c, d.secretStatus).Wrap().TextColor(c.Theme().TextMuted)
	}
	components.Text(c, m.SecretsHint()).Wrap().TextColor(c.Theme().TextMuted)
}

func secretResult(m locales.Messages, err error, ok string) string {
	switch {
	case err == nil:
		return ok
	case errors.Is(err, mygo.ErrSecretNotFound):
		return m.SecretsNotFound()
	default:
		return m.SystemError(err.Error())
	}
}

func exampleRecentDocuments(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, m.RecentChoose(), components.ButtonPrimary).Clicked() {
			d.background(func() func() {
				paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{Title: m.RecentChoose()})
				if err == nil && len(paths) > 0 {
					err = mygo.App.AddRecentDocument(paths[0])
				}
				return func() {
					if err != nil {
						d.recentStatus = m.SystemError(err.Error())
					}
				}
			})
		}
		if components.Button(c, m.Clear(), components.ButtonOutline).Clicked() && d.native {
			if err := mygo.App.ClearRecentDocuments(); err != nil {
				d.recentStatus = m.SystemError(err.Error())
			}
		}
		if runtime.GOOS == "windows" && components.Button(c, m.RecentJumpList(), components.ButtonOutline).Clicked() && d.native {
			err := mygo.App.JumpList.SetTasks([]mygo.JumpListTask{{Title: m.RecentJumpListTask(), Description: m.RecentJumpListHint()}})
			if err != nil {
				d.recentStatus = m.SystemError(err.Error())
			}
		}
	})
	var recent []string
	if d.native {
		recent = mygo.App.RecentDocuments()
	}
	if len(recent) == 0 {
		components.Text(c, m.RecentEmpty()).TextColor(c.Theme().TextMuted)
	}
	for i, path := range recent {
		ui.Column(c).Key(i).FillWidth().Children(func() {
			components.Text(c, filepath.Base(path)).Bold()
			components.Text(c, path).TextColor(c.Theme().TextMuted).FontSize(12).Wrap()
		})
	}
	if d.recentStatus != "" {
		components.Text(c, d.recentStatus).Wrap().TextColor(c.Theme().Danger)
	}
	components.Text(c, m.RecentHint()).Wrap().TextColor(c.Theme().TextMuted)
}

func (d *systemDemo) stopWatching() {
	if d.watchStop != nil {
		d.watchStop()
		d.watchStop = nil
	}
}

// watch follows a directory until stopped, keeping the latest events.
func (d *systemDemo) watch(path string) {
	d.stopWatching()
	ctx, cancel := context.WithCancel(context.Background())
	d.watchStop, d.watchPath, d.watchEvents, d.watchErr = cancel, path, nil, ""
	go func() {
		w, err := mygo.WatchFiles(ctx, path, mygo.FileWatchOptions{Recursive: true,
			Ignore: func(p string, dir bool) bool { return dir && (p == ".git" || p == "node_modules") }})
		if err != nil {
			d.update(func() { d.watchErr = err.Error() })
			return
		}
		defer w.Close()
		for {
			events, err := w.Next(ctx)
			if err != nil {
				if ctx.Err() == nil {
					d.update(func() { d.watchErr = err.Error() })
				}
				return
			}
			d.update(func() {
				d.watchEvents = append(d.watchEvents, events...)
				if extra := len(d.watchEvents) - 50; extra > 0 {
					d.watchEvents = d.watchEvents[extra:]
				}
			})
		}
	}()
}

func exampleFileWatching(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, m.WatchChoose(), components.ButtonPrimary).Clicked() {
			d.background(func() func() {
				paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{Title: m.WatchChoose(), Directory: true})
				return func() {
					if err != nil {
						d.watchErr = err.Error()
					} else if len(paths) > 0 {
						d.watch(paths[0])
					}
				}
			})
		}
		if components.Button(c, m.WatchStop(), components.ButtonOutline).Disabled(d.watchStop == nil).Clicked() {
			d.stopWatching()
		}
	})
	if d.watchPath != "" {
		systemRow(c, m.WatchFolder(), d.watchPath)
	}
	if d.watchErr != "" {
		components.Text(c, m.SystemError(d.watchErr)).Wrap().TextColor(c.Theme().Danger)
	}
	if len(d.watchEvents) == 0 {
		components.Text(c, m.WatchEmpty()).TextColor(c.Theme().TextMuted)
	}
	for i := len(d.watchEvents) - 1; i >= 0; i-- {
		e := d.watchEvents[i]
		ui.Row(c).Key(i).Gap(8).Children(func() {
			components.Badge(c, string(e.Op))
			components.Text(c, e.Path).Grow(1).Wrap()
		})
	}
}

func exampleInstalledApps(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, m.AppsLoad(), components.ButtonPrimary).Disabled(d.appsLoading).Clicked() && d.native {
			d.appsLoading = true
			d.background(func() func() {
				apps, err := mygo.Shell.InstalledApps()
				return func() {
					d.apps, d.appsLoading = apps, false
					d.appsErr = ""
					if err != nil {
						d.appsErr = err.Error()
					}
				}
			})
		}
		label := m.AppsShowPanel()
		if d.panel != nil {
			label = m.AppsHidePanel()
		}
		if components.Button(c, label, components.ButtonOutline).Clicked() && d.native {
			d.togglePanel(m)
		}
	})
	if d.appsErr != "" {
		components.Text(c, m.SystemError(d.appsErr)).Wrap().TextColor(c.Theme().Danger)
	}
	if len(d.apps) == 0 {
		components.Text(c, m.AppsHint()).Wrap().TextColor(c.Theme().TextMuted)
		return
	}
	components.TextBox(c, &d.appQuery, components.TextBoxOptions{}).Label(m.AppsSearch()).Placeholder(m.AppsSearch()).FillWidth()
	shown := 0
	for _, app := range d.apps {
		if shown == 30 {
			break
		}
		if d.appQuery != "" && !strings.Contains(strings.ToLower(app.Name), strings.ToLower(d.appQuery)) {
			continue
		}
		shown++
		ui.Row(c).Key(app.ID).FillWidth().Gap(10).AlignItems(ui.Center).Children(func() {
			if icon := d.icon(app.Path); icon != nil {
				components.Image(c, icon).Size(24, 24)
			} else {
				ui.Box(c).Size(24, 24)
			}
			components.Text(c, app.Name).Grow(1).SingleLine()
			if components.Button(c, m.Open(), components.ButtonGhost).Label(m.AppsLaunchNamed(app.Name)).Clicked() {
				go mygo.Shell.LaunchApp(app)
			}
		})
	}
}

// icon returns an app's cached icon, loading it in the background once.
func (d *systemDemo) icon(path string) *ui.Bitmap {
	if d.icons == nil {
		d.icons = map[string]*ui.Bitmap{}
	}
	if icon, ok := d.icons[path]; ok {
		return icon
	}
	d.icons[path] = nil
	d.background(func() func() {
		data, err := mygo.Shell.FileIcon(path, 48)
		var icon *ui.Bitmap
		if err == nil {
			icon, _ = ui.DecodeBitmap(data)
		}
		return func() { d.icons[path] = icon }
	})
	return nil
}

// togglePanel shows a launcher-style panel, which floats above other apps.
func (d *systemDemo) togglePanel(m locales.Messages) {
	if d.panel != nil {
		d.panel.Close()
		d.panel = nil
		return
	}
	var panel *mygo.Window
	panel = mygo.NewWindow(mygo.WindowOptions{Title: m.AppsPanelTitle(), Width: 420, Height: 160, Panel: true,
		Content: ui.View(func(c *ui.Context) {
			c.OnShortcut(0, ui.KeyEscape, func() { panel.Close() })
			ui.Column(c).Fill().Padding(20).Gap(10).Children(func() {
				components.Text(c, m.AppsPanelTitle()).Bold().FontSize(18)
				components.Text(c, m.AppsPanelHint()).Wrap().TextColor(c.Theme().TextMuted)
			})
		})})
	panel.OnClosed(func() {
		if d.panel == panel {
			d.panel = nil
			d.update(func() {})
		}
	})
	d.panel = panel
	panel.Show()
}

func exampleFonts(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	components.TextBox(c, &d.fontQuery, components.TextBoxOptions{}).Label(m.FontsSearch()).Placeholder(m.FontsSearch()).FillWidth()
	components.Checkbox(c, &d.monospace, m.FontsMonospaced())
	families := ui.FontFamilies()
	shown := 0
	for _, family := range families {
		if d.monospace && !family.Monospaced || d.fontQuery != "" && !strings.Contains(strings.ToLower(family.Name), strings.ToLower(d.fontQuery)) {
			continue
		}
		if shown++; shown > 40 {
			break
		}
		ui.Column(c).Key(family.Name).FillWidth().Gap(2).Children(func() {
			components.Text(c, family.Name).TextColor(c.Theme().TextMuted).FontSize(12)
			components.Text(c, m.FontsSample()).Font(family.Name).FontSize(18).SingleLine()
		})
	}
	components.Text(c, m.FontsCount(len(families))).TextColor(c.Theme().TextMuted)
}

func exampleSystemAppearance(c *ui.Context, s *demoState) {
	m, p, t := msg(c), c.Preferences(), c.Theme()
	ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
		components.Text(c, m.AppearanceAccent()).TextColor(t.TextMuted).Width(160)
		if p.Accent.A == 0 {
			components.Text(c, m.AppearanceNoAccent())
		} else {
			ui.Box(c).Size(20, 20).Radius(10).Background(p.Accent).Label(m.AppearanceAccent())
		}
	})
	systemRow(c, m.AppearanceReduceMotion(), yesNo(m, p.ReduceMotion))
	systemRow(c, m.AppearanceHighContrast(), yesNo(m, p.HighContrast))
	scale := p.TextScale
	if scale == 0 {
		scale = 1
	}
	systemRow(c, m.AppearanceTextScale(), i18n.FormatNumber(c.Locale(), scale, i18n.NumberOptions{Percent: true}))
	components.Text(c, m.AppearanceHint()).Wrap().TextColor(t.TextMuted)
}

func exampleLanguages(c *ui.Context, s *demoState) {
	m, r := msg(c), c.Region()
	_, format := formats(c)
	locale := c.Locale()
	if locale == "" {
		locale = "en"
	}
	systemRow(c, m.LanguagesShown(), i18n.Autonym(locale))
	var preferred []string
	if s.system.native {
		preferred = mygo.App.PreferredLanguages()
	}
	names := make([]string, 0, len(preferred))
	for _, tag := range preferred {
		names = append(names, i18n.Autonym(tag))
	}
	systemRow(c, m.LanguagesPreferred(), strings.Join(names, ", "))
	systemRow(c, m.LanguagesRegion(), r.Locale)
	systemRow(c, m.LanguagesFirstWeekday(), i18n.Weekdays(format, i18n.Wide, true)[r.FirstWeekday])
	systemRow(c, m.LanguagesNumber(), i18n.FormatNumber(format, 1234567.891, i18n.NumberOptions{}))
	systemRow(c, m.LanguagesDate(), i18n.FormatDate(format, time.Date(2026, 10, 15, 13, 7, 0, 0, time.Local), i18n.DateOptions{Date: i18n.FullStyle, Time: i18n.ShortStyle}))
	components.Text(c, m.LanguagesHint()).Wrap().TextColor(c.Theme().TextMuted)
}

// formats mirrors the regional formats MyGo's date controls use, with CLDR's
// defaults for the language where the system reports none.
func formats(c *ui.Context) (i18n.Region, string) {
	lang := c.Locale()
	if lang == "" {
		lang = "en"
	}
	r := c.Region()
	if r.HourCycle != 12 && r.HourCycle != 24 {
		r = i18n.RegionFor(lang)
		r.Locale = ""
	}
	return r, r.FormatLocale(lang)
}

func exampleWebData(c *ui.Context, s *demoState) {
	d, m := &s.system, msg(c)
	available := false
	if d.native {
		available = mygo.WebViewAvailable()
	}
	systemRow(c, m.WebDataAvailable(), yesNo(m, available))
	ui.Row(c).Gap(8).Children(func() {
		if components.Button(c, m.WebDataList(), components.ButtonPrimary).Clicked() {
			d.background(func() func() {
				cookies, err := mygo.Cookies.List(mygo.CookieFilter{})
				return func() { d.cookies, d.cookieStatus = cookies, cookieResult(m, err) }
			})
		}
		if components.Button(c, m.WebDataSet(), components.ButtonOutline).Clicked() {
			d.background(func() func() {
				err := mygo.Cookies.Set(mygo.Cookie{Name: "onevra_gallery", Value: "demo", Domain: "example.com", Expires: time.Now().Add(time.Hour)})
				cookies, listErr := mygo.Cookies.List(mygo.CookieFilter{})
				return func() { d.cookies, d.cookieStatus = cookies, cookieResult(m, errors.Join(err, listErr)) }
			})
		}
		if components.Button(c, m.WebDataClear(), components.ButtonOutline).Clicked() {
			d.background(func() func() {
				err := mygo.Cookies.Clear()
				return func() { d.cookies, d.cookieStatus = nil, cookieResult(m, err) }
			})
		}
	})
	if d.cookieStatus != "" {
		components.Text(c, d.cookieStatus).Wrap().TextColor(c.Theme().TextMuted)
	}
	cookies := slices.Clone(d.cookies)
	slices.SortFunc(cookies, func(a, b mygo.Cookie) int { return strings.Compare(a.Domain+a.Name, b.Domain+b.Name) })
	for i, cookie := range cookies {
		if i == 30 {
			break
		}
		systemRow(c, cookie.Domain, cookie.Name)
	}
}

func cookieResult(m locales.Messages, err error) string {
	if err != nil {
		return m.SystemError(err.Error())
	}
	return ""
}
