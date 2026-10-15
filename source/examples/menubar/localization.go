package main

import (
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/oneo-me/onevra/source/examples/menubar/locales"
)

// language returns the catalog language the app shows: the saved choice, or
// the best match for the user's preferred language.
func (a *application) language() string {
	if a.config.Language != "system" {
		return a.config.Language
	}
	if tag := i18n.Match([]string{a.locale}, locales.Catalog.Languages()); tag != "" {
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

// text says the message of a key chosen at run time, such as a setting's
// value or an error.
func (a *application) text(key string) string {
	return locales.Catalog.Format(a.language(), key, nil)
}

// optionKey names the message of a saved setting value, as "zh-Hans".
func optionKey(value string) string {
	return strings.ReplaceAll(value, "-", "")
}
