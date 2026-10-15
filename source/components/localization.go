package components

import (
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/internal/l10n"
)

// Localization overrides the text of component internals. Empty fields take
// Onevra's catalog in the language the window shows (ui.Context.Locale),
// which applications also change with "onevra.<key>" keys in their MyGo
// catalogs. Dates and the first weekday default to the user's regional
// formats, as MyGo's date controls follow them.
// FormatDate receives a Go layout identifying the requested presentation:
// "January 2006", "January 2, 2006", "Jan", or "2006-01-02".
// NumberAction formats a complete accessible name for "adjust", "increase",
// or "decrease", including the application-supplied field label.
// Weekdays are Sunday-first.
type Localization struct {
	Minimum, Maximum          string
	Today                     string
	ChooseMonth               string
	ChooseYear                string
	PreviousMonth             string
	NextMonth                 string
	PreviousYear              string
	NextYear                  string
	PreviousYears             string
	NextYears                 string
	MonthSelection            string
	YearSelection             string
	ChooseTime                string
	Hours                     string
	Minutes                   string
	Divider                   string
	SaturationBrightness      string
	Hue                       string
	Opacity                   string
	ColorFormat               string
	CopyColor                 string
	Red                       string
	Green                     string
	Blue                      string
	HexColor                  string
	OpacityPercent            string
	HueDegrees                string
	SaturationPercent         string
	BrightnessPercent         string
	LightnessPercent          string
	Settings                  string
	SearchSettings            string
	SearchSettingsPlaceholder string
	SettingsCategory          string
	SearchResults             string
	NoSettings                string
	CloseToast                string
	Weekdays                  [7]string
	FirstWeekday              *time.Weekday
	FormatDate                func(time.Time, string) string
	NumberAction              func(action, label string) string

	messages     l10n.Messages
	region       i18n.Region
	formatLocale string
}

type localizationKey struct{}

func localizationSlot(c *ui.Context) *Localization {
	return ui.Local(c.Root(), localizationKey{}, func() Localization { return Localization{} })
}

// WithLocalization supplies text to components built synchronously by build.
// The scope is window-local, adds no elements, and restores the enclosing scope
// even after a panic. Call during a UI build, on the UI thread; do not retain c.
func WithLocalization(c *ui.Context, words Localization, build func()) {
	current := localizationSlot(c)
	previous := *current
	*current = words
	defer func() { *current = previous }()
	if build != nil {
		build()
	}
}

func localized(c *ui.Context) Localization {
	words := *localizationSlot(c)
	m := l10n.In(c.Locale())
	for _, field := range []struct {
		value *string
		key   string
	}{
		{&words.Minimum, "minimum"}, {&words.Maximum, "maximum"}, {&words.Today, "today"},
		{&words.ChooseMonth, "chooseMonth"}, {&words.ChooseYear, "chooseYear"},
		{&words.PreviousMonth, "previousMonth"}, {&words.NextMonth, "nextMonth"},
		{&words.PreviousYear, "previousYear"}, {&words.NextYear, "nextYear"},
		{&words.PreviousYears, "previousYears"}, {&words.NextYears, "nextYears"},
		{&words.MonthSelection, "monthSelection"}, {&words.YearSelection, "yearSelection"},
		{&words.ChooseTime, "chooseTime"}, {&words.Hours, "hours"}, {&words.Minutes, "minutes"},
		{&words.Divider, "divider"}, {&words.SaturationBrightness, "saturationBrightness"},
		{&words.Hue, "hue"}, {&words.Opacity, "opacity"}, {&words.ColorFormat, "colorFormat"},
		{&words.CopyColor, "copyColor"}, {&words.Red, "red"}, {&words.Green, "green"}, {&words.Blue, "blue"},
		{&words.HexColor, "hexColor"}, {&words.OpacityPercent, "opacityPercent"}, {&words.HueDegrees, "hueDegrees"},
		{&words.SaturationPercent, "saturationPercent"}, {&words.BrightnessPercent, "brightnessPercent"},
		{&words.LightnessPercent, "lightnessPercent"}, {&words.Settings, "settings"},
		{&words.SearchSettings, "searchSettings"}, {&words.SearchSettingsPlaceholder, "searchSettingsPlaceholder"},
		{&words.SettingsCategory, "settingsCategory"}, {&words.SearchResults, "searchResults"},
		{&words.NoSettings, "noSettings"}, {&words.CloseToast, "closeToast"},
	} {
		if *field.value == "" {
			*field.value = m.Text(field.key)
		}
	}
	words.messages = m
	words.region, words.formatLocale = formats(c)
	names := i18n.Weekdays(words.formatLocale, i18n.Short, true)
	for i := range words.Weekdays {
		if words.Weekdays[i] == "" {
			words.Weekdays[i] = names[i]
		}
	}
	return words
}

// formats returns the user's regional formats and the locale that writes
// dates in them and the window's language, as MyGo's date controls do. Where
// the system reports none, as in a headless Tester, CLDR's for the language.
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

// Date presentations, named by the Go layouts FormatDate receives.
var dateSkeletons = map[string]string{"January 2006": "yMMMM", "January 2, 2006": "yMMMMd", "2006-01-02": "yMd"}

func (l Localization) date(value time.Time, layout string) string {
	if l.FormatDate != nil {
		return l.FormatDate(value, layout)
	}
	if layout == "Jan" {
		return i18n.Months(l.formatLocale, i18n.Abbreviated, true)[value.Month()-1]
	}
	if skeleton, ok := dateSkeletons[layout]; ok {
		return i18n.FormatDate(l.formatLocale, value, i18n.DateOptions{Skeleton: skeleton})
	}
	return value.Format(layout)
}
func (l Localization) numberAction(action, label string) string {
	if l.NumberAction != nil {
		return l.NumberAction(action, label)
	}
	key := map[string]string{"adjust": "numberAdjust", "increase": "numberIncrease"}[action]
	if key == "" {
		key = "numberDecrease"
	}
	return l.messages.Label(key, label)
}
func (l Localization) firstWeekday() time.Weekday {
	if l.FirstWeekday != nil && *l.FirstWeekday >= time.Sunday && *l.FirstWeekday <= time.Saturday {
		return *l.FirstWeekday
	}
	return l.region.FirstWeekday
}
