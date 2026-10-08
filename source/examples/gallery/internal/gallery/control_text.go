package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlText(c *ui.Context, _ *controlState, _ string) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(c.Theme().Space(4)).Children(func() {
		components.Text(c, "A heading").FontSize(c.Theme().Rem(1.5)).Bold()
		components.Text(c, "Hello, world. 你好，世界。こんにちは 🌍")
		components.Text(c, "Secondary text").TextColor(c.Theme().TextMuted)
		components.Text(c, "Long text wraps within the available width, including in narrow windows.").FillWidth()
		components.Text(c, "A single-line label").SingleLine()
		components.Textf(c, "%s: %d", "Count", 42)
		ui.RichText(c).Children(func() {
			components.Text(c, "Inline ")
			components.Text(c, "bold text").Bold()
		})
	})
}
