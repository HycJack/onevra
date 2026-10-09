package gallery

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/plugins/chat"
	"github.com/oneo-me/onevra/plugins/markdown"
	"github.com/oneo-me/onevra/shells/sidebar"
	"github.com/oneo-me/onevra/theme"
)

func TestNativeChat(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	dir := os.Getenv("ONEVRA_SCREENSHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []theme.Mode{theme.Light, theme.Dark} {
		for _, width := range []int{760, 360} {
			for _, status := range []chat.Status{chat.Running, chat.Completed, chat.Failed} {
				var shell *sidebar.Shell
				mygo.RunOnMain(func() {
					g := &Gallery{page: pageFromKey("controls/Chat")}
					if width == 360 {
						g.language = newTranslator("zh-Hans", "en")
					}
					s := &chatDemoState{model: "Demo model", effort: "Medium"}
					s.start("Review this component and explain the results.", time.Now().Add(time.Hour))
					s.turns[0].step, s.turns[0].status = 4, status
					s.turns[0].started = time.Now().Add(-561 * time.Second)
					if status != chat.Running {
						s.turns[0].ended = time.Now()
					}
					s.turns[0].answer = markdown.ParseMarkdown(chatDemoAnswer)
					g.controlState("Chat", "Default").chat = s
					height := 820
					if width == 760 {
						height = 1000
					}
					shell = sidebar.NewShell(sidebar.ShellOptions{
						Title: "Chat review", Width: width, Height: height, MinWidth: 280, MinHeight: 400,
						LeftCollapsed: true, Appearance: mode,
						Content: func(c *ui.Context) {
							ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() { g.body(c) })
						},
					})
				})
				win := shell.Window()
				time.Sleep(250 * time.Millisecond)
				data, err := win.CapturePage()
				win.Destroy()
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(dir, fmt.Sprintf("chat-%v-%d-%d.png", mode, width, status))
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	t.Logf("Chat captures: %s", dir)
}
