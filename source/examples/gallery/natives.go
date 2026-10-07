package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// MyGo's terminal does not search Contents/Frameworks. Point it at the bundled
// library before the gallery can start a terminal, preserving explicit overrides.
func configureTerminalLibrary() error {
	if runtime.GOOS != "darwin" || os.Getenv("MYGO_GHOSTTY_VT") != "" {
		return nil
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}
	return useBundledTerminalLibrary(exe)
}

func useBundledTerminalLibrary(exe string) error {
	if os.Getenv("MYGO_GHOSTTY_VT") != "" {
		return nil
	}
	macos := filepath.Dir(exe)
	contents := filepath.Dir(macos)
	if filepath.Base(macos) != "MacOS" || filepath.Base(contents) != "Contents" ||
		!strings.HasSuffix(filepath.Dir(contents), ".app") {
		return nil // Direct go run and tests retain the plugin's normal lookup.
	}
	// Set even if missing, so a broken bundle reports the missing bundled file
	// instead of accidentally working only on a developer's cached library.
	return os.Setenv("MYGO_GHOSTTY_VT", filepath.Join(contents, "Frameworks", "libghostty-vt.dylib"))
}
