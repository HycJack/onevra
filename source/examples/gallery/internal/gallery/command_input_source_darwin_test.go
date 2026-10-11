package gallery

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ebitengine/purego/objc"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

func TestNativeCommandInputSource(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_FOREGROUND") != "1" {
		t.Skip("requires ONEVRA_NATIVE_FOREGROUND=1")
	}
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	var shell *sidebar.Shell
	var state command.State
	open, text := false, ""
	keep := false
	source := command.InputSourceABC
	mygo.RunOnMain(func() {
		shell = sidebar.NewShell(sidebar.ShellOptions{Title: "Input source review", Content: func(c *ui.Context) {
			components.TextBox(c, &text, components.TextBoxOptions{}).Label("Original editor").AutoFocus()
			command.CommandDialog(c, &open, &state, nil, command.Options{InputSource: source, KeepInputSource: keep}, components.DialogOptions{})
		}})
		shell.Window().Focus()
	})
	win := shell.Window()
	defer win.Destroy()
	time.Sleep(200 * time.Millisecond)
	context := func() objc.ID {
		return objc.ID(objc.GetClass("NSApplication")).Send(objc.RegisterName("sharedApplication")).Send(objc.RegisterName("keyWindow")).Send(objc.RegisterName("firstResponder")).Send(objc.RegisterName("inputContext"))
	}
	current := func() string {
		s := context().Send(objc.RegisterName("selectedKeyboardInputSource"))
		if s == 0 {
			return ""
		}
		return objc.Send[string](s, objc.RegisterName("UTF8String"))
	}
	set := func(source string) {
		name := objc.ID(objc.GetClass("NSString")).Send(objc.RegisterName("stringWithUTF8String:"), source)
		context().Send(objc.RegisterName("setSelectedKeyboardInputSource:"), name)
	}
	var original, alternate string
	var hasABC bool
	mygo.RunOnMain(func() {
		original = current()
		sources := context().Send(objc.RegisterName("keyboardInputSources"))
		for i, count := uintptr(0), uintptr(sources.Send(objc.RegisterName("count"))); i < count; i++ {
			id := sources.Send(objc.RegisterName("objectAtIndex:"), i)
			name := objc.Send[string](id, objc.RegisterName("UTF8String"))
			if name == command.InputSourceABC {
				hasABC = true
			} else if alternate == "" {
				alternate = name
			}
		}
	})
	if !hasABC || original == "" || alternate == "" {
		t.Skipf("requires ABC and another enabled source; current=%q ABC=%v alternate=%q", original, hasABC, alternate)
	}
	defer mygo.RunOnMain(func() { state.ReleaseInputSource(); set(original) })
	waitFor := func(want string) {
		t.Helper()
		deadline := time.Now().Add(2 * time.Second)
		var selected string
		var selectionErr error
		for time.Now().Before(deadline) {
			mygo.RunOnMain(func() { selected, selectionErr = current(), state.InputSourceError })
			if selectionErr != nil {
				t.Fatal(selectionErr)
			}
			if selected == want {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		t.Fatalf("selected=%q want %q", selected, want)
	}
	mygo.RunOnMain(func() { set(alternate); open = true; win.Invalidate() })
	waitFor(command.InputSourceABC)
	mygo.RunOnMain(func() { open = false; win.Invalidate() })
	waitFor(alternate)
	mygo.RunOnMain(func() { open = true; win.Invalidate() })
	waitFor(command.InputSourceABC)
	mygo.RunOnMain(func() { set(alternate); win.Invalidate() })
	time.Sleep(200 * time.Millisecond)
	waitFor(alternate)
	mygo.RunOnMain(func() { open = false; win.Invalidate() })
	time.Sleep(200 * time.Millisecond)
	waitFor(alternate)
	mygo.RunOnMain(func() { keep = true; open = true; win.Invalidate() })
	waitFor(command.InputSourceABC)
	mygo.RunOnMain(func() { open = false; win.Invalidate() })
	time.Sleep(200 * time.Millisecond)
	waitFor(command.InputSourceABC)
	mygo.RunOnMain(func() { source = "onevra.invalid-source"; open = true; win.Invalidate() })
	deadline := time.Now().Add(2 * time.Second)
	var selectionErr error
	var selected string
	for time.Now().Before(deadline) {
		mygo.RunOnMain(func() { selectionErr, selected = state.InputSourceError, current() })
		if selectionErr != nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !errors.Is(selectionErr, command.ErrInputSourceUnavailable) || selected != command.InputSourceABC {
		t.Fatalf("unavailable source: selected=%q error=%v", selected, selectionErr)
	}
}
