package gallery

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/egoist/mygo"
	oneoshell "github.com/oneo-me/oneo-go/shell"
)

// State contains only presentation state, never interactive demo values.
type State struct {
	Window       *WindowState `json:"window,omitempty"`
	SidebarWidth float32      `json:"sidebar_width,omitempty"`
	SelectedItem string       `json:"selected_item,omitempty"`
}

type WindowState struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
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
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

// SaveState replaces the file only after a complete write in the same directory.
func SaveState(path string, state State) error {
	return saveJSON(path, state)
}

// saveJSON atomically replaces an application-owned JSON file.
func saveJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
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
	if _, err = f.Write(append(data, '\n')); err != nil {
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
func (g *Gallery) NewPersistent(path string) *oneoshell.Shell {
	state, err := LoadState(path)
	if err != nil {
		log.Printf("Load gallery state: %v", err)
	}
	g.page = pageFromKey(state.SelectedItem)
	s := g.New()
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
