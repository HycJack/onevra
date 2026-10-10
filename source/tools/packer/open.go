package main

import (
	"context"
	"fmt"
	"path/filepath"
)

func openDirectory(ctx context.Context, path, hostOS string) error {
	// Absolute paths cannot be interpreted as flags by the desktop opener.
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	name, err := directoryOpener(hostOS)
	if err != nil {
		return err
	}
	_, err = command(ctx, name, path)
	return err
}

func directoryOpener(hostOS string) (string, error) {
	switch hostOS {
	case "darwin":
		return "open", nil
	case "windows":
		return "explorer.exe", nil
	case "linux":
		return "xdg-open", nil
	default:
		return "", fmt.Errorf("unsupported desktop host %q", hostOS)
	}
}
