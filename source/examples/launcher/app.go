package main

import (
	"cmp"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

const (
	width, height = 680, 420
	iconSize      = 28 // DIPs
	iconLimit     = 24 // apps whose icons are drawn ahead of the first frame
)

// application holds what the panel shows. It is read and changed on the main
// thread, by the view and the callbacks of the app, the tray and the shortcut;
// the installed apps and their icons come from goroutines through panel.Update.
type application struct {
	config      config
	path        string
	locale      string
	native      bool
	ready       bool
	quitting    bool
	pendingShow bool

	errorKey, errorDetail string

	panel, settings *mygo.Window
	tray            *mygo.Tray
	trayMenu        *mygo.Menu
	settingsState   components.SettingsState

	shortcut, accelerator string
	login                 bool

	apps        []mygo.InstalledApp
	results     []mygo.InstalledApp
	query       string
	chosen      int
	opening     bool
	loading     bool
	list        components.ListState
	input       ui.Handle
	icons       map[string]*ui.Bitmap
	iconPending map[string]bool
}

func (a *application) report(key string, err error) {
	a.errorKey, a.errorDetail = key, err.Error()
	log.Printf("%s: %v", key, err)
}

// themeSource selects light or dark for the whole app, including native
// window chrome; config values match MyGo's theme sources.
func (a *application) themeSource() mygo.ThemeSource { return mygo.ThemeSource(a.config.Appearance) }

func (a *application) start() {
	if preferred := mygo.App.PreferredLanguages(); len(preferred) > 0 {
		a.locale = preferred[0]
	}
	a.icons, a.iconPending = map[string]*ui.Bitmap{}, map[string]bool{}
	a.login = mygo.App.OpenAtLogin()
	a.createPanel()
	if err := a.createTray(); err != nil {
		a.report("trayError", err)
	}
	a.startShortcut()
	a.native, a.ready = true, true
	a.refresh()
	switch {
	case a.accelerator == "" && a.tray == nil:
		// Without a shortcut or a tray icon, Settings is the way in.
		a.openSettings()
	case !mygo.App.WasOpenedAtLogin():
		a.show() // started by the user: show what it does
	}
	if a.pendingShow {
		a.pendingShow = false
		a.show()
	}
}

// createPanel makes the launcher panel hidden and once: it floats above other
// apps and is hidden, not closed, so it shows again at once, as it was.
func (a *application) createPanel() {
	a.panel = mygo.NewWindow(mygo.WindowOptions{
		Title:         a.text("title"),
		Panel:         true,
		Hidden:        true,
		Frameless:     true,
		Transparent:   true,
		Vibrancy:      mygo.VibrancyPopover, // Acrylic on Windows 11
		Width:         width,
		Height:        height,
		DisableResize: true,
		Content:       ui.View(a.panelView),
	})
	// An Escape that the view leaves, Cmd-W or Alt-F4 hides the panel; quitting
	// destroys it.
	a.panel.OnClose(func(e *mygo.CloseEvent) {
		if a.quitting {
			return
		}
		e.PreventDefault()
		a.panel.Hide()
	})
	a.panel.OnBlur(a.hide)
}

// show shows the panel on the display with the pointer, a third of the way
// down, with an empty query, and gives its search field the keyboard.
func (a *application) show() {
	if a.panel == nil {
		return
	}
	d := mygo.Screen.DisplayNearestPoint(mygo.Screen.CursorScreenPoint())
	w, h := a.panel.Size()
	a.panel.SetPosition(d.WorkArea.X+(d.WorkArea.Width-w)/2, d.WorkArea.Y+(d.WorkArea.Height-h)/3)
	a.query, a.opening = "", false
	a.filter()
	a.load()
	a.input.Focus()
	a.panel.Show()
}

func (a *application) hide() {
	if a.panel != nil {
		a.panel.Hide()
	}
}

func (a *application) toggle() {
	if a.panel != nil && a.panel.IsVisible() {
		a.hide()
	} else {
		a.show()
	}
}

// load lists the installed apps on a goroutine, as at each show, where apps
// may have come or gone since.
func (a *application) load() {
	a.loading = true
	go func() {
		apps, err := mygo.Shell.InstalledApps()
		if err != nil {
			log.Print(err)
		}
		a.panel.Update(func() {
			a.apps, a.loading = apps, false
			a.filter()
		})
	}()
}

// filter ranks the apps matching the query: names starting with it first, then
// those with a word starting with it, those containing it, and those with its
// letters in order, as "vsc" finds Visual Studio Code.
func (a *application) filter() {
	q := strings.ToLower(strings.TrimSpace(a.query))
	type ranked struct {
		app  mygo.InstalledApp
		rank int
	}
	var found []ranked
	for _, app := range a.apps {
		if r := rank(strings.ToLower(app.Name), q); r >= 0 {
			found = append(found, ranked{app, r})
		}
	}
	slices.SortStableFunc(found, func(x, y ranked) int { return cmp.Compare(x.rank, y.rank) })
	a.results = a.results[:0]
	for _, f := range found {
		a.results = append(a.results, f.app)
	}
	a.choose(0)
}

func rank(name, q string) int {
	switch {
	case q == "" || strings.HasPrefix(name, q):
		return 0
	case strings.Contains(name, " "+q) || strings.Contains(name, "-"+q):
		return 1
	case strings.Contains(name, q):
		return 2
	}
	rest := name
	for _, r := range q {
		i := strings.IndexRune(rest, r)
		if i < 0 {
			return -1
		}
		rest = rest[i+len(string(r)):]
	}
	return 3
}

func (a *application) choose(i int) {
	a.chosen = max(0, min(i, len(a.results)-1))
	if len(a.results) > 0 {
		a.list.ScrollIntoView(a.chosen)
	}
}

// open launches the chosen app and hides the panel, which gives the keyboard
// back to the app the user was in, before the launched app takes it.
func (a *application) open() {
	if a.opening || a.chosen < 0 || a.chosen >= len(a.results) {
		return
	}
	a.opening = true
	app := a.results[a.chosen]
	a.hide()
	if err := mygo.Shell.LaunchApp(app); err != nil {
		log.Printf("launching %s: %v", app.Name, err)
	}
}

// prefetchIcons draws the icons of the rows that show ahead of their first
// frame, on a goroutine, and keeps them by path; a row without one shows its
// initial.
func (a *application) prefetchIcons() {
	if !a.native || len(a.results) == 0 {
		return
	}
	var todo []string
	for i, app := range a.results {
		if i >= iconLimit {
			break
		}
		if _, ok := a.icons[app.Path]; ok || a.iconPending[app.Path] {
			continue
		}
		a.iconPending[app.Path] = true
		todo = append(todo, app.Path)
	}
	if len(todo) == 0 {
		return
	}
	scale := 1.0
	for _, d := range mygo.Screen.Displays() {
		scale = max(scale, d.ScaleFactor)
	}
	size := int(iconSize * scale)
	go func() {
		for _, path := range todo {
			var bitmap *ui.Bitmap
			if data, err := mygo.Shell.FileIcon(path, size); err == nil {
				bitmap, _ = ui.DecodeBitmap(data)
			}
			a.panel.Update(func() { a.icons[path] = bitmap })
		}
	}()
}

func (a *application) icon(app mygo.InstalledApp) *ui.Bitmap {
	return a.icons[app.Path]
}

func (a *application) openSettings() {
	a.hide()
	if a.settings == nil {
		a.settings = a.newWindow(plain.ShellOptions{
			Title: a.text("settingsTitle"), Width: 620, Height: 520, MinWidth: 380, MinHeight: 420,
			Hidden: true, Content: func(c *ui.Context) {
				a.scope(c, func() { w, _ := c.Size(); a.settingsView(c, w) })
			},
		})
		a.settings.OnClose(func(e *mygo.CloseEvent) {
			if a.quitting {
				return
			}
			e.PreventDefault()
			a.settings.Hide()
		})
	}
	show(a.settings)
}

func (a *application) changed() {
	if a.path != "" {
		if err := saveConfig(a.path, a.config); err != nil {
			a.report("saveError", err)
		} else if a.errorKey == "saveError" || a.errorKey == "loadError" {
			a.errorKey, a.errorDetail = "", ""
		}
	}
	if a.ready {
		a.refresh()
	}
}

func (a *application) refresh() {
	mygo.Theme.SetSource(a.themeSource())
	a.applyLanguage()
	a.installMenus()
	for _, item := range []struct {
		win   *mygo.Window
		title string
	}{
		{a.panel, "title"}, {a.settings, "settingsTitle"},
	} {
		if item.win != nil {
			item.win.SetTitle(a.text(item.title))
			item.win.Invalidate()
		}
	}
	if a.tray != nil {
		a.tray.SetToolTip(a.text("title"))
	}
}

func (a *application) newWindow(options plain.ShellOptions) *mygo.Window {
	options.Theme = theme.Onevra
	shell := plain.NewShell(options)
	return shell.Window()
}

func show(win *mygo.Window) {
	if win.IsMinimized() {
		win.Restore()
	}
	win.Show()
	win.Focus()
}

// where says where an app is, shortly: its folder or its ID.
func where(app mygo.InstalledApp) string {
	switch runtime.GOOS {
	case "darwin":
		dir := filepath.Dir(app.Path)
		if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(dir, home) {
			return "~" + dir[len(home):]
		}
		for _, system := range []string{"/System/Applications", "/System/Library/CoreServices/Applications", "/System/Library/CoreServices"} {
			if rest, ok := strings.CutPrefix(dir, system); ok {
				return "System" + rest
			}
		}
		return strings.TrimPrefix(dir, "/")
	case "windows":
		// A desktop app's program, or a packaged app's ID.
		if i := strings.LastIndex(app.ID, `\`); i >= 0 {
			return app.ID[i+1:]
		}
		return "App"
	}
	return strings.TrimSuffix(app.ID, ".desktop")
}

// initial is an app's first letter, for a row whose icon is not drawn yet.
func initial(name string) string {
	for _, r := range strings.TrimSpace(name) {
		return strings.ToUpper(string(r))
	}
	return "?"
}

// divider draws a hairline between regions without reaching for a widget
// Onevra intentionally omits.
func divider(c *ui.Context) {
	ui.Box(c).FillWidth().Height(1).Background(c.Theme().Border)
}

// keys spells a shortcut as the platform does: ⌥⇧Space on macOS.
func keys(shortcut string) string {
	if runtime.GOOS != "darwin" {
		return strings.NewReplacer("CmdOrCtrl+", "Ctrl+", "Cmd+", "Ctrl+").Replace(shortcut)
	}
	return strings.NewReplacer("CmdOrCtrl+", "⌘", "Cmd+", "⌘", "Ctrl+", "⌃", "Alt+", "⌥", "Shift+", "⇧").Replace(shortcut)
}

// appIcon draws an app's icon when it is ready, and its initial otherwise.
func (a *application) appIcon(c *ui.Context, app mygo.InstalledApp) {
	t := c.Theme()
	if bitmap := a.icon(app); bitmap != nil {
		components.Image(c, bitmap).Size(iconSize, iconSize).Shrink(0)
		return
	}
	ui.Box(c).Size(iconSize, iconSize).Radius(t.Radius).Background(t.Surface).Center().Shrink(0).Children(func() {
		components.Text(c, initial(app.Name)).FontSize(iconSize * 0.45).TextColor(t.TextMuted)
	})
}
