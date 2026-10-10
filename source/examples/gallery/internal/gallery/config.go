package gallery

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"go.yaml.in/yaml/v3"
)

// Config holds application preferences, separate from transient window state.
type Config struct {
	Theme            mygo.ThemeSource `yaml:"theme"`
	Language         string           `yaml:"language,omitempty"`
	DisablePixelGrid bool             `yaml:"disablePixelGrid,omitempty"`
	// Zero values preserve automatic ABC selection and initially visible results.
	DisableCommandABC     bool `yaml:"disableCommandABC,omitempty"`
	CommandExpandOnSearch bool `yaml:"commandExpandOnSearch,omitempty"`
}

func validTheme(source mygo.ThemeSource) bool {
	return source == mygo.ThemeSystem || source == mygo.ThemeLight || source == mygo.ThemeDark
}

func LoadConfig(path string) (Config, error) {
	fallback := Config{Theme: mygo.ThemeSystem}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return fallback, nil
	}
	if err != nil {
		return fallback, err
	}
	config := fallback
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fallback, err
	}
	if !validTheme(config.Theme) {
		return fallback, fmt.Errorf("unknown theme %q", config.Theme)
	}
	if !validLanguage(config.Language) {
		return Config{Theme: config.Theme}, fmt.Errorf("unknown language %q", config.Language)
	}
	return config, nil
}

func SaveConfig(path string, config Config) error {
	if !validTheme(config.Theme) {
		return fmt.Errorf("unknown theme %q", config.Theme)
	}
	if !validLanguage(config.Language) {
		return fmt.Errorf("unknown language %q", config.Language)
	}
	return saveYAML(path, config)
}

// Configure loads preferences before window creation. A nonempty override
// applies only to this launch; explicit menu choices are saved immediately.
func (g *Gallery) Configure(path string, override mygo.ThemeSource) {
	config, err := LoadConfig(path)
	if err != nil {
		log.Printf("Load gallery config: %v", err)
	}
	g.config, g.configPath = config, path
	g.systemLanguage = mygo.App.Locale()
	g.applyLanguage()
	source := config.Theme
	if validTheme(override) {
		source = override
	}
	mygo.Theme.SetSource(source)
}

func (g *Gallery) chooseTheme(source mygo.ThemeSource) {
	if !validTheme(source) {
		return
	}
	mygo.Theme.SetSource(source)
	g.config.Theme = source
	if g.configPath != "" {
		if err := SaveConfig(g.configPath, g.config); err != nil {
			log.Printf("Save gallery config: %v", err)
		}
	}
}

func themeIcon() *ui.SVG {
	switch mygo.Theme.Source() {
	case mygo.ThemeLight:
		return lucide.Sun()
	case mygo.ThemeDark:
		return lucide.Moon()
	default:
		return lucide.Monitor()
	}
}

func (g *Gallery) themeMenu(m *ui.Menu) {
	for _, option := range []struct {
		label  string
		source mygo.ThemeSource
	}{
		{"System", mygo.ThemeSystem},
		{"Light", mygo.ThemeLight},
		{"Dark", mygo.ThemeDark},
	} {
		if m.Item(g.text(option.label)).Checked(mygo.Theme.Source() == option.source).Chosen() {
			g.chooseTheme(option.source)
		}
	}
}

func (g *Gallery) savePreferences() {
	if g.configPath != "" {
		if err := SaveConfig(g.configPath, g.config); err != nil {
			log.Printf("Save gallery config: %v", err)
		}
	}
}
