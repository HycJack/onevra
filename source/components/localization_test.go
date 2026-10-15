package components

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestLocalizationScopes(t *testing.T) {
	a, b := 0.0, 0.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			WithLocalization(c, Localization{Settings: "设置", NumberAction: func(action, label string) string { return action + ":" + label }}, func() {
				NumberInput(c, &a, 0, 10, 1).Label("数量")
				WithLocalization(c, Localization{Settings: "Inner"}, func() {
					if localized(c).Settings != "Inner" {
						t.Fatal("nested scope")
					}
				})
				if localized(c).Settings != "设置" {
					t.Fatal("outer scope not restored")
				}
			})
			NumberInput(c, &b, 0, 10, 1).Label("Count")
			if localized(c).Settings != "Settings" {
				t.Fatal("default not restored")
			}
		})
	}, 500, 400)
	if err := tt.Click("increase:数量"); err != nil {
		t.Fatal(err)
	}
	if a != 1 || b != 0 {
		t.Fatal(a, b)
	}
	if err := tt.Click("Count Increase"); err != nil {
		t.Fatal(err)
	}
	if a != 1 || b != 1 {
		t.Fatal(a, b)
	}
}

func TestLocalizedCalendarNavigation(t *testing.T) {
	day := time.Date(2026, 10, 9, 12, 30, 0, 0, time.UTC)
	sunday := time.Sunday
	words := Localization{Today: "今天", NextMonth: "下个月", Weekdays: [7]string{"日", "一", "二", "三", "四", "五", "六"}, FirstWeekday: &sunday,
		FormatDate: func(d time.Time, layout string) string {
			if layout == "January 2006" {
				return fmt.Sprintf("%d年%d月", d.Year(), d.Month())
			}
			if layout == "January 2, 2006" {
				return fmt.Sprintf("%d年%d月%d日", d.Year(), d.Month(), d.Day())
			}
			return d.Format(layout)
		}}
	tt := ui.NewTester(func(c *ui.Context) { WithLocalization(c, words, func() { Calendar(c, &day) }) }, 500, 400)
	if !tt.HasText("2026年10月") || !tt.HasText("日") {
		t.Fatal(tt.Texts())
	}
	if err := tt.Click("下个月"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("2026年11月") || day.Month() != time.October {
		t.Fatal("navigation changed value")
	}
	if err := tt.Click("2026年11月1日"); err != nil {
		t.Fatal(err)
	}
	if day.Month() != time.November || day.Day() != 1 || day.Hour() != 12 {
		t.Fatal(day)
	}
}

func TestCatalogFollowsWindowLanguage(t *testing.T) {
	day := time.Date(2026, 10, 9, 12, 30, 0, 0, time.UTC)
	count := 0.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			Calendar(c, &day)
			NumberInput(c, &count, 0, 10, 1).Label("数量")
		})
	}, 500, 500)
	tt.SetLocale("zh-Hans")
	for _, want := range []string{"2026年10月", "今天", "一", "下个月", "增加数量"} {
		if _, ok := tt.Find(want); !ok && !tt.HasText(want) {
			t.Errorf("zh-Hans: missing %q in %q", want, tt.Texts())
		}
	}
	tt.SetLocale("en")
	if _, ok := tt.Find("Next month"); !ok || !tt.HasText("October 2026") {
		t.Fatalf("en: %q", tt.Texts())
	}
}
