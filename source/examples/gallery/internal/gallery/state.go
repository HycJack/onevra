package gallery

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/source/shells/sidebar"
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
	X          int  `yaml:"x"`
	Y          int  `yaml:"y"`
	Width      int  `yaml:"width"`
	Height     int  `yaml:"height"`
	Maximized  bool `yaml:"maximized,omitempty"`
	FullScreen bool `yaml:"full_screen,omitempty"`
}

const mainWindowStateKey = "main"

// windowStateStore keeps MyGo's state of the main window in state.yaml.
// MyGo calls it on the main thread, as are the gallery's own saves.
type windowStateStore struct {
	path  string
	state *State
}

func (s windowStateStore) LoadWindowStates() (map[string]mygo.WindowState, error) {
	w := s.state.Window
	if w == nil {
		return nil, nil
	}
	return map[string]mygo.WindowState{mainWindowStateKey: mygo.WindowState(*w)}, nil
}

func (s windowStateStore) SaveWindowStates(states map[string]mygo.WindowState) error {
	w, ok := states[mainWindowStateKey]
	if !ok {
		return nil
	}
	state := WindowState(w)
	s.state.Window = &state
	return SaveState(s.path, *s.state)
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
	if page == notificationsPage {
		return "desktop/notifications"
	}
	if page == dockPage {
		return "desktop/dock-menu"
	}
	if page == trayPage {
		return "desktop/tray"
	}
	if page == welcomePage {
		return "welcome"
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
	if key == "desktop/notifications" {
		return notificationsPage
	}
	if key == "desktop/dock-menu" {
		return dockPage
	}
	if key == "desktop/tray" {
		return trayPage
	}
	// Restore the grouped page for selections saved before the merge.
	switch key {
	case "mygo/Charts":
		key = "controls/Bar chart"
	case "mygo/Web data":
		key = "mygo/WebView"
	case "controls/TextField", "controls/TextArea", "controls/TextBox (multiline)":
		key = "controls/TextBox"
	case "controls/Checkbox", "controls/Switch", "controls/Radio":
		key = "controls/Selection"
	}
	if key == "welcome" {
		return welcomePage
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
func (g *Gallery) NewPersistent(path string) *sidebar.Shell {
	state, err := LoadState(path)
	if err != nil {
		log.Printf("Load gallery state: %v", err)
	}
	g.page = pageFromKey(state.SelectedItem)
	// MyGo captures, restores and saves the window's geometry through the
	// store, which shares this state so either save writes the whole file.
	mygo.App.SetWindowStateStore(windowStateStore{path, &state})
	g.stateKey = mainWindowStateKey
	s := g.New()
	s.SetLeftExpanded(!state.SidebarCollapsed)
	w := s.Window()
	if state.SidebarWidth > 0 {
		s.SetLeftWidth(state.SidebarWidth)
	}
	w.OnClose(func(*mygo.CloseEvent) {
		state.SidebarWidth = s.LeftWidth()
		state.SidebarCollapsed = !s.LeftExpanded()
		state.SelectedItem = pageKey(g.page)
		if err := SaveState(path, state); err != nil {
			log.Printf("Save gallery state: %v", err)
		}
	})
	return s
}
