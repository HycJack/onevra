package components

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// Localization supplies application-owned text for component internals. Empty
// fields keep English defaults. It does not select languages or load catalogs.
// FormatDate receives a Go layout identifying the requested presentation:
// "January 2006", "January 2, 2006", "Jan", or "2006-01-02".
// NumberAction formats a complete accessible name for "adjust", "increase",
// or "decrease", including the application-supplied field label.
// Weekdays are Sunday-first; FirstWeekday is optional and defaults to Monday.
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
	Weekdays                  [7]string
	FirstWeekday              *time.Weekday
	FormatDate                func(time.Time, string) string
	NumberAction              func(action, label string) string
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
	if words.Minimum == "" {
		words.Minimum = "minimum"
	}
	if words.Maximum == "" {
		words.Maximum = "maximum"
	}
	if words.Today == "" {
		words.Today = "Today"
	}
	if words.ChooseMonth == "" {
		words.ChooseMonth = "Choose month"
	}
	if words.ChooseYear == "" {
		words.ChooseYear = "Choose year"
	}
	if words.PreviousMonth == "" {
		words.PreviousMonth = "Previous month"
	}
	if words.NextMonth == "" {
		words.NextMonth = "Next month"
	}
	if words.PreviousYear == "" {
		words.PreviousYear = "Previous year"
	}
	if words.NextYear == "" {
		words.NextYear = "Next year"
	}
	if words.PreviousYears == "" {
		words.PreviousYears = "Previous years"
	}
	if words.NextYears == "" {
		words.NextYears = "Next years"
	}
	if words.MonthSelection == "" {
		words.MonthSelection = "Month selection"
	}
	if words.YearSelection == "" {
		words.YearSelection = "Year selection"
	}
	if words.ChooseTime == "" {
		words.ChooseTime = "Choose time"
	}
	if words.Hours == "" {
		words.Hours = "Hours"
	}
	if words.Minutes == "" {
		words.Minutes = "Minutes"
	}
	if words.Divider == "" {
		words.Divider = "Divider"
	}
	if words.SaturationBrightness == "" {
		words.SaturationBrightness = "Saturation and brightness"
	}
	if words.Hue == "" {
		words.Hue = "Hue"
	}
	if words.Opacity == "" {
		words.Opacity = "Opacity"
	}
	if words.ColorFormat == "" {
		words.ColorFormat = "Color format"
	}
	if words.CopyColor == "" {
		words.CopyColor = "Copy color"
	}
	if words.Red == "" {
		words.Red = "Red"
	}
	if words.Green == "" {
		words.Green = "Green"
	}
	if words.Blue == "" {
		words.Blue = "Blue"
	}
	if words.HexColor == "" {
		words.HexColor = "Hex color"
	}
	if words.OpacityPercent == "" {
		words.OpacityPercent = "Opacity percent"
	}
	if words.HueDegrees == "" {
		words.HueDegrees = "Hue degrees"
	}
	if words.SaturationPercent == "" {
		words.SaturationPercent = "Saturation percent"
	}
	if words.BrightnessPercent == "" {
		words.BrightnessPercent = "Brightness percent"
	}
	if words.LightnessPercent == "" {
		words.LightnessPercent = "Lightness percent"
	}
	if words.Settings == "" {
		words.Settings = "Settings"
	}
	if words.SearchSettings == "" {
		words.SearchSettings = "Search settings"
	}
	if words.SearchSettingsPlaceholder == "" {
		words.SearchSettingsPlaceholder = "Search settings…"
	}
	if words.SettingsCategory == "" {
		words.SettingsCategory = "Settings category"
	}
	if words.SearchResults == "" {
		words.SearchResults = "Search results"
	}
	if words.NoSettings == "" {
		words.NoSettings = "No matching settings. Try another search."
	}
	for i, day := range [7]string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"} {
		if words.Weekdays[i] == "" {
			words.Weekdays[i] = day
		}
	}
	return words
}
func (l Localization) date(value time.Time, layout string) string {
	if l.FormatDate != nil {
		return l.FormatDate(value, layout)
	}
	return value.Format(layout)
}
func (l Localization) numberAction(action, label string) string {
	if l.NumberAction != nil {
		return l.NumberAction(action, label)
	}
	switch action {
	case "adjust":
		return strings.TrimSpace("Adjust " + label)
	case "increase":
		return strings.TrimSpace(label + " Increase")
	default:
		return strings.TrimSpace(label + " Decrease")
	}
}
func (l Localization) firstWeekday() time.Weekday {
	if l.FirstWeekday != nil && *l.FirstWeekday >= time.Sunday && *l.FirstWeekday <= time.Saturday {
		return *l.FirstWeekday
	}
	return time.Monday
}
