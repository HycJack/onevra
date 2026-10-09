package gallery

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/shell"
	"go.yaml.in/yaml/v3"
)

// State contains only presentation state, never interactive demo values.
type State struct {
	Window           *WindowState `yaml:"window,omitempty"`
	SidebarWidth     float32      `yaml:"sidebar_width,omitempty"`
	SidebarCollapsed bool         `yaml:"sidebar_collapsed"`
	SelectedItem     string       `yaml:"selected_item,omitempty"`
}

type WindowState struct {
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

func LoadState(path string) (State, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	var state State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

// SaveState replaces the file only after a complete write in the same directory.
func SaveState(path string, state State) error {
	return saveYAML(path, state)
}

// saveYAML atomically replaces an application-owned YAML file.
func saveYAML(path string, value any) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".gallery-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func pageKey(page galleryPage) string {
	if page == trayPage {
		return "desktop/tray"
	}
	if page == welcomePage {
		return "welcome"
	}
	if page == componentStatusPage {
		return "component-status"
	}
	if page == shellPage {
		return "shell"
	}
	if e, ok := controlExampleFor(page); ok {
		return "controls/" + e.name
	}
	if e, ok := mygoExampleFor(page); ok {
		return "mygo/" + e.name
	}
	return "lucide"
}

func pageFromKey(key string) galleryPage {
	if key == "desktop/tray" {
		return trayPage
	}
	// Restore the grouped page for selections saved before the merge.
	switch key {
	case "controls/TextField", "controls/TextArea", "controls/TextBox (multiline)":
		key = "controls/TextBox"
	case "controls/Checkbox", "controls/Switch", "controls/Radio":
		key = "controls/Selection"
	}
	if key == "welcome" {
		return welcomePage
	}
	if key == "component-status" {
		return componentStatusPage
	}
	if key == "shell" {
		return shellPage
	}
	for i := range mygoCatalog {
		page := mygoPageStart + galleryPage(i)
		if pageKey(page) == key {
			return page
		}
	}
	for i := range controlCatalog {
		page := controlPageStart() + galleryPage(i)
		if pageKey(page) == key {
			return page
		}
	}
	return lucidePage
}

// NewPersistent is the application-owned startup path. New remains disk-free for
// isolated demonstrations and tests. Normal close and application quit save state.
func (g *Gallery) NewPersistent(path string) *shell.Shell {
	state, err := LoadState(path)
	if err != nil {
		log.Printf("Load gallery state: %v", err)
	}
	g.page = pageFromKey(state.SelectedItem)
	s := g.New()
	s.SetLeftExpanded(!state.SidebarCollapsed)
	w := s.Window()
	if state.SidebarWidth > 0 {
		s.SetLeftWidth(state.SidebarWidth)
	}
	if b := state.Window; b != nil && b.Width > 0 && b.Height > 0 {
		bounds := mygo.Rectangle{X: b.X, Y: b.Y, Width: max(820, b.Width), Height: max(600, b.Height)}
		area := mygo.Screen.DisplayMatching(bounds).WorkArea
		if area.Width > 0 && area.Height > 0 {
			bounds.Width = max(820, min(bounds.Width, area.Width))
			bounds.Height = max(600, min(bounds.Height, area.Height))
			bounds.X = max(area.X, min(bounds.X, area.X+area.Width-bounds.Width))
			bounds.Y = max(area.Y, min(bounds.Y, area.Y+area.Height-bounds.Height))
		}
		w.SetBounds(bounds)
	}
	capture := func() {
		if w.IsMinimized() || w.IsMaximized() || w.IsFullScreen() {
			return
		}
		b := w.Bounds()
		if b.Width > 0 && b.Height > 0 {
			state.Window = &WindowState{b.X, b.Y, b.Width, b.Height}
		}
	}
	capture()
	var timer *time.Timer
	closed := false
	changed := func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(300*time.Millisecond, func() {
			mygo.RunOnMain(func() {
				if !closed {
					capture()
				}
			})
		})
	}
	w.OnMove(changed)
	w.OnResize(changed)
	w.OnClose(func(*mygo.CloseEvent) {
		capture()
		state.SidebarWidth = s.LeftWidth()
		state.SidebarCollapsed = !s.LeftExpanded()
		state.SelectedItem = pageKey(g.page)
		if err := SaveState(path, state); err != nil {
			log.Printf("Save gallery state: %v", err)
		}
	})
	w.OnClosed(func() {
		closed = true
		if timer != nil {
			timer.Stop()
		}
	})
	return s
}
