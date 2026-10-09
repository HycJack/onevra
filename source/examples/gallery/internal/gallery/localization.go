package gallery

import (
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/shells/sidebar"
	"go.yaml.in/yaml/v3"
	"golang.org/x/text/language"
)

// Catalog IDs stay stable when English copy changes. Static UI uses IDs;
// source lookup is only for existing demonstration metadata and variant names.
// Neither form is used as a control key or persisted option value.
//
//go:embed locales/*.yaml
var localeFiles embed.FS

var englishMessages, sourceMessageIDs = loadEnglishMessages()
var messages = loadMessages()

func loadEnglishMessages() (map[string]string, map[string]string) {
	data, err := localeFiles.ReadFile("locales/en.yaml")
	if err != nil {
		panic(err)
	}
	var catalog map[string]struct {
		Other string `yaml:"other"`
	}
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		panic(err)
	}
	english, sources := map[string]string{}, map[string]string{}
	for id, message := range catalog {
		english[id] = message.Other
		// Contextual messages are addressed only by ID.
		if id != "gallery.clear_glass_style" {
			sources[message.Other] = id
		}
	}
	return english, sources
}

func loadMessages() *i18n.Bundle {
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
}

type translator struct {
	code      string
	localizer *i18n.Localizer
}

func validLanguage(code string) bool { return code == "" || code == "en" || code == "zh-Hans" }
func newTranslator(choice, system string) *translator {
	code := choice
	if !validLanguage(code) || code == "" {
		code = "en"
		tag, err := language.Parse(system)
		if err == nil {
			base, _ := tag.Base()
			script, _ := tag.Script()
			// Do not silently equate Traditional Chinese with Simplified Chinese.
			if base.String() == "zh" && script.String() == "Hans" {
				code = "zh-Hans"
			}
		}
	}
	return &translator{code: code, localizer: i18n.NewLocalizer(messages, code, "en")}
}
func (t *translator) text(key string) string {
	id := key
	if sourceID, ok := sourceMessageIDs[key]; ok {
		id = sourceID
	}
	fallback, ok := englishMessages[id]
	if !ok {
		return key
	}
	if t == nil || t.code == "en" {
		return fallback
	}
	text, err := t.localizer.Localize(&i18n.LocalizeConfig{MessageID: id, DefaultMessage: &i18n.Message{ID: id, Other: fallback}})
	if err != nil && text == "" {
		return fallback
	}
	return text
}
func (g *Gallery) text(source string) string { return g.language.text(source) }

type languageScopeKey struct{}

func languageSlot(c *ui.Context) **translator {
	return ui.Local(c.Root(), languageScopeKey{}, func() *translator { return nil })
}
func tr(c *ui.Context, source string) string { return (*languageSlot(c)).text(source) }

func (g *Gallery) withLanguage(c *ui.Context, build func()) {
	slot := languageSlot(c)
	previous := *slot
	*slot = g.language
	defer func() { *slot = previous }()
	components.WithLocalization(c, g.componentWords(), build)
}
func (g *Gallery) applyLanguage() {
	g.language = newTranslator(g.config.Language, g.systemLanguage)
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
		g.Shell.SetLabels(sidebar.Labels{
			LeftToggle: g.text("gallery.toggle_sidebar"), LeftResize: g.text("gallery.resize_sidebar"),
			Minimize: g.text("gallery.minimize_window"), Maximize: g.text("gallery.maximize_window"),
			Restore: g.text("gallery.restore_window"), Close: g.text("gallery.close_window"),
		})
		if win := g.Shell.Window(); win != nil {
			win.SetTitle(g.text("gallery.onevra_gallery"))
			win.Invalidate()
		}
	}
	g.tray.localize = g.withLanguage
	if g.tray.win != nil {
		g.tray.win.SetTitle(g.text("gallery.gallery_tray_notes"))
		g.tray.win.Invalidate()
	}
}

func (g *Gallery) componentWords() components.Localization {
	words := components.Localization{
		Minimum: g.text("gallery.minimum"), Maximum: g.text("gallery.maximum"),
		Today:                     g.text("gallery.today"),
		ChooseMonth:               g.text("gallery.choose_month"),
		ChooseYear:                g.text("gallery.choose_year"),
		PreviousMonth:             g.text("gallery.previous_month"),
		NextMonth:                 g.text("gallery.next_month"),
		PreviousYear:              g.text("gallery.previous_year"),
		NextYear:                  g.text("gallery.next_year"),
		PreviousYears:             g.text("gallery.previous_years"),
		NextYears:                 g.text("gallery.next_years"),
		MonthSelection:            g.text("gallery.month_selection"),
		YearSelection:             g.text("gallery.year_selection"),
		ChooseTime:                g.text("gallery.choose_time"),
		Hours:                     g.text("gallery.hours"),
		Minutes:                   g.text("gallery.minutes"),
		Divider:                   g.text("gallery.divider"),
		SaturationBrightness:      g.text("gallery.saturation_and_brightness"),
		Hue:                       g.text("gallery.hue"),
		Opacity:                   g.text("gallery.opacity"),
		ColorFormat:               g.text("gallery.color_format"),
		CopyColor:                 g.text("gallery.copy_color"),
		Red:                       g.text("gallery.red"),
		Green:                     g.text("gallery.green"),
		Blue:                      g.text("gallery.blue"),
		HexColor:                  g.text("gallery.hex_color"),
		OpacityPercent:            g.text("gallery.opacity_percent"),
		HueDegrees:                g.text("gallery.hue_degrees"),
		SaturationPercent:         g.text("gallery.saturation_percent"),
		BrightnessPercent:         g.text("gallery.brightness_percent"),
		LightnessPercent:          g.text("gallery.lightness_percent"),
		Settings:                  g.text("gallery.settings"),
		SearchSettings:            g.text("gallery.search_settings"),
		SearchSettingsPlaceholder: g.text("gallery.search_settings_2"),
		SettingsCategory:          g.text("gallery.settings_category"),
		SearchResults:             g.text("gallery.search_results"),
		NoSettings:                g.text("gallery.no_matching_settings_try_another_search"),
	}
	words.NumberAction = func(action, label string) string {
		source := map[string]string{"adjust": "gallery.adjust_s", "increase": "gallery.s_increase", "decrease": "gallery.s_decrease"}[action]
		return strings.TrimSpace(fmt.Sprintf(g.text(source), label))
	}
	if g.language != nil && g.language.code == "zh-Hans" {
		words.Weekdays = [7]string{"日", "一", "二", "三", "四", "五", "六"}
		words.FormatDate = func(value time.Time, layout string) string {
			switch layout {
			case "January 2006":
				return fmt.Sprintf("%d年%d月", value.Year(), value.Month())
			case "January 2, 2006":
				return fmt.Sprintf("%d年%d月%d日", value.Year(), value.Month(), value.Day())
			case "Jan":
				return fmt.Sprintf("%d月", value.Month())
			default:
				return value.Format(layout)
			}
		}
	}
	return words
}
