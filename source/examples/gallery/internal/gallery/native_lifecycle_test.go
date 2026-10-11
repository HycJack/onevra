package gallery

import (
	"os"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// Closing from a content builder must end the frame before presentation, as
// happens when Shell's Windows caption close button activates during a build.
func TestNativeCloseDuringFrame(t *testing.T) {
	if os.Getenv("ONEVRA_NATIVE_SMOKE") != "1" {
		t.Skip("requires a desktop session")
	}
	closed := make(chan struct{}, 1)
	returned := make(chan struct{}, 1)
	var win *mygo.Window
	mygo.RunOnMain(func() {
		s := sidebar.NewShell(sidebar.ShellOptions{Hidden: nativeBackground()})
		win = s.Window()
		win.OnClosed(func() { closed <- struct{}{} })
		s.SetPage(func(c *ui.Context) {
			win.Close()
		})
	})
	defer win.Destroy()
	select {
	case <-closed:
	case <-time.After(10 * time.Second):
		t.Fatal("window did not close during its content build")
	}
	// A subsequent main-thread task must run after the closing frame unwinds.
	go func() {
		mygo.RunOnMain(func() { returned <- struct{}{} })
	}()
	select {
	case <-returned:
	case <-time.After(10 * time.Second):
		t.Fatal("closing frame did not return to the event loop")
	}
}
