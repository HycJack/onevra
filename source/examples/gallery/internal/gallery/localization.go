package gallery

import (
	_ "embed"
	"encoding/json"
	"log"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery/locales"
)

// Message keys stay stable when English copy changes. Static UI uses keys;
// English source lookup is only for existing demonstration metadata and
// variant names. Neither form is used as a control key or persisted value.
//
//go:embed locales/en.json
var englishCatalog []byte

// sourceKeys maps English source text to its message key.
var sourceKeys = func() map[string]string {
	var catalog map[string]any
	if err := json.Unmarshal(englishCatalog, &catalog); err != nil {
		panic(err)
	}
	keys := map[string]string{}
	for key, message := range catalog {
		text, ok := message.(string)
		// Contextual messages are addressed only by key.
		if !ok || key[0] == '@' || key == "clearGlassStyle" {
			continue
		}
		keys[text] = key
	}
	return keys
}()

// A translator says the gallery's messages in the language Settings chose.
type translator struct{ code string }

func validLanguage(code string) bool { return code == "" || code == "en" || code == "zh-Hans" }

// newTranslator resolves a choice, "" following the user's language, against
// the languages of the catalog.
func newTranslator(choice, system string) *translator {
	code := choice
	if !validLanguage(code) || code == "" {
		code = i18n.Match([]string{system}, locales.Catalog.Languages())
		if code == "" {
			code = locales.Catalog.Source()
		}
	}
	return &translator{code: code}
}
func (t *translator) messages() locales.Messages {
	if t == nil {
		return locales.In("")
	}
	return locales.In(t.code)
}

// text says the message of a key, or of English source text, without arguments.
func (t *translator) text(key string) string {
	if sourceKey, ok := sourceKeys[key]; ok {
		key = sourceKey
	}
	if _, _, ok := locales.Catalog.Message("", key); !ok {
		return key
	}
	code := ""
	if t != nil {
		code = t.code
	}
	return locales.Catalog.Format(code, key, nil)
}
func (g *Gallery) text(source string) string  { return g.language.text(source) }
func (g *Gallery) messages() locales.Messages { return g.language.messages() }

type languageScopeKey struct{}

func languageSlot(c *ui.Context) **translator {
	return ui.Local(c.Root(), languageScopeKey{}, func() *translator { return nil })
}

// msg returns the typed messages in the language of the enclosing gallery scope.
func msg(c *ui.Context) locales.Messages { return (*languageSlot(c)).messages() }

func tr(c *ui.Context, source string) string { return (*languageSlot(c)).text(source) }

func (g *Gallery) withLanguage(c *ui.Context, build func()) {
	slot := languageSlot(c)
	previous := *slot
	*slot = g.language
	defer func() { *slot = previous }()
	build()
}
func (g *Gallery) applyLanguage() {
	g.language = newTranslator(g.config.Language, g.systemLanguage)
	// MyGo's and Onevra's own strings follow the same choice.
	mygo.App.SetLocale(g.config.Language)
	g.updateLanguageWindows()
}
func (g *Gallery) chooseLanguage(choice string) {
	if !validLanguage(choice) {
		return
	}
	g.config.Language = choice
	g.applyLanguage()
	if g.configPath != "" {
		if err := SaveConfig(g.configPath, g.config); err != nil {
			log.Printf("Save gallery config: %v", err)
		}
	}
}
func (g *Gallery) updateLanguageWindows() {
	if g.Shell != nil {
		if win := g.Shell.Window(); win != nil {
			win.SetTitle(g.text("onevraGallery"))
			win.Invalidate()
		}
	}
	g.tray.localize = g.withLanguage
	g.tray.translate = g.text
	if g.tray.tray != nil {
		g.tray.tray.SetToolTip(g.text("trayTooltip"))
	}
	if g.tray.menu != nil {
		g.tray.menu.ItemByID("notes").SetLabel(g.text("openNotes"))
		g.tray.menu.ItemByID("gallery").SetLabel(g.text("showGallery"))
		g.tray.menu.ItemByID("quit").SetLabel(g.text("quit"))
		// AppIndicator retains a native snapshot rather than reading the menu on click.
		if runtime.GOOS == "linux" && g.tray.tray != nil {
			g.tray.tray.SetMenu(g.tray.menu)
		}
	}
	if g.appMenu != nil {
		g.appMenu.ItemByID("settings").SetLabel(g.text("settingsMenu"))
	}
	if g.dock.menu != nil {
		g.dock.menu.ItemByID("show-gallery").SetLabel(g.text("showGallery"))
		g.dock.menu.ItemByID("welcome").SetLabel(g.text("openWelcome"))
		g.dock.menu.ItemByID("settings").SetLabel(g.text("openSettings"))
		// Refresh the native snapshot only while Gallery still owns the Dock menu.
		if mygo.App.Dock.Menu() == g.dock.menu {
			mygo.App.Dock.SetMenu(g.dock.menu)
		}
	}
	if g.tray.win != nil {
		g.tray.win.SetTitle(g.text("galleryTrayNotes"))
		g.tray.win.Invalidate()
	}
}
