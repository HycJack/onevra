// Package gallery holds the standalone showcase, not public component API.
package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/oneo-go"
	"github.com/oneo-me/oneo-go/icons/lucide"
)

// Gallery owns only demonstration state.
type Gallery struct {
	Disabled bool
	Created  bool
}

// iconSamples uses only named accessors, not a full-library runtime registry.
func iconSamples(c *ui.Context) {
	ui.Text(c, "Lucide icons").Bold()
	ui.Text(c, "Original strokes at 16, 24 and 32 DIP; colors inherit from MyGo.").Wrap()
	for _, sample := range []struct {
		name string
		get  func() *ui.SVG
	}{
		{"Search", lucide.Search},
		{"Settings", lucide.Settings},
		{"Cloud rain", lucide.CloudRain},
		{"Fingerprint", lucide.FingerprintPattern},
		{"Accessibility", lucide.Accessibility},
	} {
		ui.Column(c).Gap(4).Children(func() {
			ui.Text(c, sample.name)
			ui.Row(c).Gap(12).Children(func() {
				for _, size := range []float32{16, 24, 32} {
					ui.Icon(c, sample.get()).Size(size, size)
				}
			})
		})
	}
	ui.Row(c).Gap(12).TextColor(c.Theme().TextMuted).Children(func() {
		ui.Icon(c, lucide.Sun())
		ui.Icon(c, lucide.Moon())
		ui.Text(c, "Inherited muted color")
	})
}

// View builds the same showcase for native windows and headless checks.
func (g *Gallery) View(c *ui.Context) {
	t := c.Theme()
	ui.Scroll(c).Fill().Background(t.Background).Children(func() {
		ui.Column(c).FillWidth().MaxWidth(t.Rem(52)).Margin(0, ui.Auto).Padding(t.Space(6)).Gap(t.Space(6)).Children(func() {
			ui.Text(c, "oneo-go").FontSize(t.Rem(2)).Bold()
			ui.Text(c, "Native compositions. Your state, your application.").TextColor(t.TextMuted)
			ui.Divider(c)
			oneogo.Section(c, "Empty state", "A starting point with a caller-owned action.", func() {
				ui.Checkbox(c, &g.Disabled, "Disable example action")
				oneogo.EmptyState(c, oneogo.EmptyStateOptions{
					Title:       "Nothing here yet",
					Description: "Create your first item to get started. 内容由使用方提供，组件不绑定业务。",
					Actions: func() {
						if ui.PrimaryButton(c, "Create item").Disabled(g.Disabled).Clicked() {
							g.Created = true
						}
					},
				})
				if g.Created {
					oneogo.Notice(c, oneogo.NoticeOptions{Kind: oneogo.NoticeSuccess, Title: "Item created", Description: "This is demo state only; no files were written.", Actions: func() {
						if ui.Button(c, "Reset example").Clicked() {
							g.Created = false
						}
					}})
				}
			})
			ui.Divider(c)
			oneogo.Section(c, "Notices", "Persistent feedback without taking over the page.", func() {
				oneogo.Notice(c, oneogo.NoticeOptions{Title: "Information", Description: "Theme, focus and text scale come from MyGo."})
				oneogo.Notice(c, oneogo.NoticeOptions{Kind: oneogo.NoticeWarning, Title: "Needs attention", Description: "Long descriptions wrap naturally as the window narrows. 调整窗口宽度可检查中英文混排。"})
				oneogo.Notice(c, oneogo.NoticeOptions{Kind: oneogo.NoticeError, Title: "Could not complete the action", Description: "Applications supply recovery actions and manage their own asynchronous work."})
			})
			ui.Divider(c)
			oneogo.Section(c, "Composition", "Section accepts ordinary MyGo widgets; there is no second widget system.", func() {
				ui.Text(c, "Buttons, inputs, focus and keyboard behavior remain upstream.")
			})
			iconSamples(c)
		})
	})
}
