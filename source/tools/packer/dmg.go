package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// createDMG lets Finder author its own layout metadata on a writable image,
// then detaches it before compression so the final image retains the layout.
func createDMG(ctx context.Context, stage, root, name, appName, output string) (err error) {
	background := filepath.Join(root, ".background", "install.png")
	if err := writeDMGBackground(background); err != nil {
		return err
	}
	writable := filepath.Join(stage, "layout.dmg")
	if _, err := command(ctx, "hdiutil", "create", "-quiet", "-fs", "HFS+", "-format", "UDRW", "-volname", name, "-srcfolder", root, writable); err != nil {
		return err
	}
	// Finder identifies mounted disks by their mountpoint name; make it unique
	// across targets and concurrent builds to avoid stale window references.
	mount, err := os.MkdirTemp(stage, "mount-")
	if err != nil {
		return err
	}
	if _, err := command(ctx, "hdiutil", "attach", "-quiet", "-nobrowse", "-noautoopen", "-mountpoint", mount, writable); err != nil {
		return err
	}
	mounted := true
	defer func() {
		if mounted {
			// Cancellation must not leave a mounted staging image behind.
			cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if _, detachErr := command(cleanup, "hdiutil", "detach", "-quiet", "-force", mount); detachErr != nil {
				err = fmt.Errorf("%w; cleanup: %v", err, detachErr)
			}
		}
	}()
	// Arguments are passed separately, never interpolated into AppleScript.
	layoutCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	if _, err := command(layoutCtx, "osascript", "-e", dmgLayoutScript, mount, appName); err != nil {
		return fmt.Errorf("save DMG layout (requires a desktop session and permission to control Finder): %w", err)
	}
	info, err := os.Stat(filepath.Join(mount, ".DS_Store"))
	if err != nil {
		return fmt.Errorf("Finder did not save the DMG layout: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("Finder saved an empty DMG layout")
	}
	if _, err := command(ctx, "hdiutil", "detach", "-quiet", mount); err != nil {
		return err
	}
	mounted = false
	_, err = command(ctx, "hdiutil", "convert", "-quiet", writable, "-format", "UDZO", "-o", output)
	return err
}

const dmgLayoutScript = `on run argv
	set imageFolder to POSIX file (item 1 of argv) as alias
	set appName to item 2 of argv
	tell application "Finder"
		set imageFolder to folder imageFolder
		open imageFolder
		delay 1
		tell imageFolder
			set current view of container window to icon view
			set toolbar visible of container window to false
			set statusbar visible of container window to false
			set bounds of container window to {160, 120, 880, 580}
			set options to icon view options of container window
			set arrangement of options to not arranged
			set icon size of options to 160
			set text size of options to 14
			set label position of options to bottom
			set shows item info of options to false
			set shows icon preview of options to false
			set background picture of options to file ".background:install.png" of imageFolder
			set position of item appName to {180, 185}
			set position of item "Applications" to {540, 185}
		end tell
		update imageFolder without registering applications
		delay 2
		close container window of imageFolder
	end tell
	-- Finder writes .DS_Store asynchronously after the window closes.
	delay 2
end run`

// The background contains only a direction cue; both installation targets are
// real Finder items, not pictures baked into the background.
func writeDMGBackground(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Finder displays PNG backgrounds at their pixel dimensions, so use the
	// window's logical dimensions rather than relying on PNG DPI metadata.
	img := image.NewRGBA(image.Rect(0, 0, 720, 440))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	ink := color.RGBA{R: 145, G: 151, B: 163, A: 255}
	for y := 165; y <= 205; y++ {
		for x := 338; x <= 382; x++ {
			d := y - 185
			if d < 0 {
				d = -d
			}
			if (x <= 377 && d <= 3) || (x >= 360 && x <= 382 && d <= 382-x && d >= 375-x) {
				img.SetRGBA(x, y, ink)
			}
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
