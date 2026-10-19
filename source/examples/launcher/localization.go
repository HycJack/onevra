package main

import (
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/oneo-me/onevra/source/examples/launcher/locales"
)

// language returns the catalog language the app shows: the saved choice, or
// the best match for the user's preferred languages.
func (a *application) language() string {
	if a.config.Language != "system" {
		return a.config.Language
	}
	if tag := i18n.Match(a.preferred, locales.Catalog.Languages()); tag != "" {
		return tag
	}
	return locales.Catalog.Source()
}

// applyLanguage makes MyGo's and Onevra's own strings follow the saved choice;
// "system" follows the user's languages again.
func (a *application) applyLanguage() {
	tag := a.config.Language
	if tag == "system" {
		tag = ""
	}
	mygo.App.SetLocale(tag)
}

// followSystemLanguages reads the user's languages and keeps reading them
// when they change while the app runs, so that its own messages follow
// "system" as MyGo's and Onevra's strings do.
func (a *application) followSystemLanguages() {
	a.preferred = mygo.App.PreferredLanguages()
	mygo.App.OnLocaleChanged(func(string) {
		before := a.language()
		a.preferred = mygo.App.PreferredLanguages()
		if a.ready && a.language() != before {
			a.refresh()
		}
	})
}

// text says the message of a key chosen at run time, such as a setting's
// value or an error.
func (a *application) text(key string) string {
	return a.format(key, nil)
}

// format says a message with arguments, such as the shortcut that opens the
// launcher.
func (a *application) format(key string, args map[string]any) string {
	return locales.Catalog.Format(a.language(), key, args)
}

// optionKey names the message of a saved setting value, as "zh-Hans".
func optionKey(value string) string {
	return strings.ReplaceAll(value, "-", "")
}
