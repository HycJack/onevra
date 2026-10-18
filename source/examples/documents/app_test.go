package main

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/source/examples/documents/locales"
)

func TestStatus(t *testing.T) {
	for _, lang := range []string{"en", "zh-Hans", "zh-Hant", "ja", "ko"} {
		m := locales.In(lang)
		for state, want := range map[mygo.DocumentState]string{
			{}:                             m.NotSaved(),
			{Path: "a.md", ReadOnly: true}: m.ReadOnly(),
			{Path: "a.md", Edited: true}:   m.Edited(),
			{Path: "a.md"}:                 m.Saved(),
		} {
			if got := status(m, state); got != want {
				t.Errorf("%s %+v: %q, want %q", lang, state, got, want)
			}
		}
	}
	if locales.In("zh-Hans").Words(3) != "3 个词" || locales.In("en").Words(1) != "1 word" {
		t.Fatal("word counts")
	}
}
