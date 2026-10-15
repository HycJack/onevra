package gallery

import (
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// defaultShortcut avoids Alt+Space and Alt+Shift+Space, which launchers take:
// macOS gives a shortcut to every app that registers it.
const defaultShortcut = "CmdOrCtrl+Alt+Shift+G"

// registerShortcut moves the demonstration's global shortcut to accelerator,
// keeping the previous one when the system refuses the new one.
func (d *systemDemo) registerShortcut(accelerator string) {
	accelerator = strings.TrimSpace(accelerator)
	if !d.native || accelerator == "" {
		return
	}
	previous := d.shortcut
	d.unregisterShortcut()
	if err := mygo.GlobalShortcut.Register(accelerator, d.shortcutPressed); err != nil {
		d.shortcutErr = err.Error()
		if previous != "" && mygo.GlobalShortcut.Register(previous, d.shortcutPressed) == nil {
			d.shortcut = previous
		}
		return
	}
	d.shortcut, d.shortcutErr = accelerator, ""
}

func (d *systemDemo) unregisterShortcut() {
	if d.shortcut != "" {
		mygo.GlobalShortcut.Unregister(d.shortcut)
		d.shortcut = ""
	}
}

// shortcutPressed runs on the main thread, also while another app is focused.
func (d *systemDemo) shortcutPressed() {
	d.update(func() {
		d.shortcutPresses++
		d.shortcutLast = time.Now()
		if w := d.window(); d.shortcutFront && w != nil {
			if w.IsMinimized() {
				w.Restore()
			}
			w.Show()
			w.Focus()
		}
	})
}

func exampleGlobalShortcuts(c *ui.Context, s *demoState) {
	d, m, t := &s.system, msg(c), c.Theme()
	if d.shortcutDraft == "" && d.shortcut == "" {
		d.shortcutDraft = defaultShortcut
	}
	components.TextBox(c, &d.shortcutDraft, components.TextBoxOptions{}).Label(m.ShortcutAccelerator()).FillWidth()
	ui.Row(c).Gap(8).Children(func() {
		draft := strings.TrimSpace(d.shortcutDraft)
		if components.Button(c, m.ShortcutRegister(), components.ButtonPrimary).
			Disabled(draft == "" || draft == d.shortcut).Clicked() {
			d.registerShortcut(draft)
		}
		if components.Button(c, m.ShortcutUnregister(), components.ButtonOutline).Disabled(d.shortcut == "").Clicked() {
			d.unregisterShortcut()
		}
	})
	components.Checkbox(c, &d.shortcutFront, m.ShortcutBringToFront())
	systemRow(c, m.ShortcutRegistered(), d.shortcut)
	systemRow(c, m.ShortcutPresses(), strconv.Itoa(d.shortcutPresses))
	last := ""
	if !d.shortcutLast.IsZero() {
		last = i18n.FormatDate(c.Locale(), d.shortcutLast, i18n.DateOptions{Time: i18n.MediumStyle})
	}
	systemRow(c, m.ShortcutLastPressed(), last)
	if d.shortcutErr != "" {
		components.Text(c, m.SystemError(d.shortcutErr)).Wrap().TextColor(t.Danger)
	}
	components.Text(c, m.ShortcutHint()).Wrap().TextColor(t.TextMuted)
}
