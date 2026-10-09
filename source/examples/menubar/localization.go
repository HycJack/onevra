package main

import (
	"embed"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.yaml.in/yaml/v3"
	"golang.org/x/text/language"
)

//go:embed locales/*.yaml
var localeFiles embed.FS

var messages = func() *i18n.Bundle {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	for _, code := range []string{"en", "zh-Hans"} {
		data, err := localeFiles.ReadFile("locales/" + code + ".yaml")
		if err != nil {
			panic(err)
		}
		if _, err := b.ParseMessageFileBytes(data, code+".yaml"); err != nil {
			panic(err)
		}
	}
	return b
}()

func languageCode(choice, system string) string {
	if choice == "en" || choice == "zh-Hans" {
		return choice
	}
	tag, err := language.Parse(system)
	if err == nil {
		base, _ := tag.Base()
		script, _ := tag.Script()
		if base.String() == "zh" && script.String() == "Hans" {
			return "zh-Hans"
		}
	}
	return "en"
}

func (a *application) text(id string) string {
	localizer := i18n.NewLocalizer(messages, languageCode(a.config.Language, a.locale), "en")
	value, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil && value == "" {
		return id
	}
	return value
}
