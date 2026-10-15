package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Preferences belong to this application, not to Onevra or its window shell.
type config struct {
	Language    string `json:"language"`
	Appearance  string `json:"appearance"`
	CloseAction string `json:"closeAction"`
	Silent      bool   `json:"silent"`
	Panel       bool   `json:"panel"`
	Welcomed    bool   `json:"welcomed"`
}

func defaultConfig() config {
	action := "dock"
	if appKind == "hybrid" {
		action = "tray"
	}
	return config{Language: "system", Appearance: "system", CloseAction: action, Panel: true}
}

func (c *config) normalize() {
	if c.Language != "en" && c.Language != "zh-Hans" {
		c.Language = "system"
	}
	if c.Appearance != "light" && c.Appearance != "dark" {
		c.Appearance = "system"
	}
	valid := c.CloseAction == "quit" || c.CloseAction == "dock" && appKind == "desktop" || (c.CloseAction == "tray" || c.CloseAction == "both") && appKind == "hybrid"
	if !valid {
		c.CloseAction = defaultConfig().CloseAction
	}
}

func loadConfig(path string) (config, error) {
	c := defaultConfig()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return defaultConfig(), err
	}
	c.normalize()
	return c, nil
}

// Replace the file atomically so an interrupted save leaves the old preferences.
func saveConfig(path string, c config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".preferences-*")
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

// Closing a window and explicitly quitting are different actions. Other desktop
// platforms cannot use a macOS Dock entry as a way back into a hidden app.
func keepRunning(c config, platform string, trayAvailable bool) bool {
	if appKind == "menubar" {
		return trayAvailable
	}
	if appKind == "desktop" {
		return platform == "darwin" && c.CloseAction == "dock"
	}
	return c.CloseAction != "quit" && trayAvailable
}
