package main

import (
	"errors"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

// Preferences belong to this application, not to Onevra or its window shell.
type config struct {
	Language   string `yaml:"language"`
	Appearance string `yaml:"appearance"`
	Shortcut   string `yaml:"shortcut"`
}

func defaultConfig() config {
	return config{Language: "system", Appearance: "system", Shortcut: defaultShortcut}
}

func (c *config) normalize() {
	if c.Language != "en" && c.Language != "zh-Hans" {
		c.Language = "system"
	}
	if c.Appearance != "light" && c.Appearance != "dark" {
		c.Appearance = "system"
	}
	if acceleratorFor(c.Shortcut) == "" {
		c.Shortcut = defaultShortcut
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
	if err := yaml.Unmarshal(data, &c); err != nil {
		return defaultConfig(), err
	}
	c.normalize()
	return c, nil
}

// Replace the file atomically so an interrupted save leaves the old preferences.
func saveConfig(path string, c config) error {
	data, err := yaml.Marshal(c)
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
