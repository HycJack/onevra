package gallery

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestCommandShortcutAndNavigation(t *testing.T) {
	g := &Gallery{}
	tt := ui.NewTester(g.View, 820, 600)
	mods, _ := commandShortcut()
	tt.Key(mods, ui.KeyF)
	if !g.commands.open || !tt.Focused("Command palette") {
		t.Fatal("opening shortcut did not focus command search")
	}
	tt.Type("TextBox")
	tt.Key(0, ui.KeyEnter)
	if g.commands.open || pageKey(g.page) != "controls/TextBox" {
		t.Fatalf("command navigation: %s", pageKey(g.page))
	}
	tt.Key(mods, ui.KeyF)
	if g.commands.palette.Query != "" {
		t.Fatal("gallery did not reset search on reopening")
	}
	tt.Type("settings")
	tt.Key(0, ui.KeyEnter)
	if !g.settings.open || g.commands.open {
		t.Fatal("settings command did not hand over dialog")
	}
	tt.Key(mods, ui.KeyF)
	if g.commands.open {
		t.Fatal("shortcut escaped settings modal")
	}
}

func TestCommandLanguageAndSidebarShortcut(t *testing.T) {
	g := &Gallery{page: welcomePage}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() {
			ui.Column(c).Width(220).FillHeight().Children(func() { g.Left(c) })
			ui.Column(c).Fill().Children(func() { g.View(c) })
		})
	}, 820, 600)
	if err := tt.Click("Search components"); err != nil {
		t.Fatal(err)
	}
	if !g.commands.open || !tt.Focused("Command palette") {
		t.Fatal("sidebar trigger did not focus command search")
	}
	tt.Key(0, ui.KeyEscape)
	if !tt.Focused("Search components") {
		t.Fatal("closing palette did not restore trigger focus")
	}
	mods, _ := commandShortcut()
	tt.Key(mods, ui.KeyF)
	if !g.commands.open {
		t.Fatal("shortcut did not reopen from sidebar trigger")
	}
	g.chooseLanguage("zh-Hans")
	tt.Frame()
	if !tt.Focused("命令面板") || !tt.HasText("设置") {
		t.Fatal("language switch lost command focus or labels")
	}
	tt.Type("settings")
	tt.Key(0, ui.KeyEnter)
	if !g.settings.open {
		t.Fatal("English alias failed in Chinese UI")
	}
}
