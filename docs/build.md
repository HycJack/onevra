# Development

Use Go 1.27.1 or later. Run Go commands from `source/`; dependencies are pinned in [go.mod](../source/go.mod). No sibling project or workspace file is required.

Native windows require a desktop session. Linux also needs GTK 3, Pango and OpenGL system libraries. See the [MyGo native guides](https://github.com/egoist/mygo/tree/v0.2.12/docs/ui) for platform setup.

## App bundles

The local [packer](../source/tools/packer/main.go) builds desktop apps for selected platforms and runs the host version. Run it from `source/`:

```sh
go run ./tools/packer run
go run ./tools/packer build
```

Packer accepts only the `run` and `build` commands, with no options or application-argument forwarding. It reads `packer.yaml` from the working directory; all packaging settings belong there.

Packer does not run code generators. The gallery uses its checked-in [MyGo component review inventory](architecture/gallery.md#component-review).

The top-level `platforms` array in [packer.yaml](../source/packer.yaml) selects release targets using `GOOS/GOARCH[/FORMAT]`: `darwin`, `windows` or `linux`, each with `arm64` or `amd64`. The gallery configuration selects macOS arm64, Windows amd64 (x64), and Linux amd64 and arm64. Edit this list to change release targets. If omitted or empty, Packer builds the host platform. An optional format selects `dmg` or `zip` on macOS, `zip` on Windows, or `tar.gz` on Linux. For example, separate `darwin/arm64/dmg` and `darwin/arm64/zip` entries produce both archives from one build. Entries for the same OS and architecture merge their formats; duplicates are emitted once. Without a format, macOS selects DMG and ZIP, Windows ZIP, and Linux tar.gz. Unsupported targets or format combinations fail before building.

`run` always uses the host and ignores the release target list. It compiles once with ordinary Go optimizations and caching, attaches terminal I/O and waits for exit. On macOS it retains `dist/Onevra Gallery Dev.app` with a separate `.dev` bundle identifier; Windows and Linux retain a `dist/Onevra Gallery Dev/` directory. The next run replaces the development output. There is no file watcher or hot reload.

`build` uses `-trimpath` and linker `-s -w` flags to omit symbol tables and debug information. Both Windows development and release executables use the GUI subsystem so launching them does not open a console window. Packer does not inject framework-specific variables or change application runtime modes. Outputs go into `dist/releases/<version>/`, with no platform subdirectories or retained release application directories:

| Target | Artifacts |
| --- | --- |
| macOS | `onevra-gallery-darwin-arm64-<version>.dmg`, containing `Onevra Gallery.app` and an Applications shortcut |
| Windows | `onevra-gallery-windows-amd64-<version>.zip`, containing the app directory and `.exe` |
| Linux | `onevra-gallery-linux-<arch>-<version>.tar.gz`, containing the app directory and executable |

`<version>` comes from `version` in `packer.yaml`. The version determines the release directory and archive suffix, not the contained application or development output. For example, version `0.2026.1` writes archives to `dist/releases/0.2026.1/`. Different versions coexist in separate directories; rebuilding the same target and version replaces its archives.

After all release targets build successfully, Packer opens the `releases/` directory under the configured output path in the host's file manager. Set top-level `openReleases: false` in `packer.yaml` to disable this behavior, for example in CI or a headless session. Failure to open the directory produces a warning without failing the completed build; `run` does not open it.

Each target finishes compilation and packaging in staging before publication, with backups for rollback. Targets build sequentially: a later target's failure does not remove an earlier target's successful archives. Development outputs and other targets remain untouched.

[packer.yaml](../source/packer.yaml) owns application packaging configuration; Packer does not read `mygo.json` or import MyGo Go packages. macOS application metadata goes into Info.plist. Both Go compilation and dependency discovery use the selected target's `GOOS`/`GOARCH` with cgo disabled. For native dependencies, Packer uses `go list` to locate the selected application's imported packages and reads their `mygo-plugin.json` manifests as data. Versions, download URLs and SHA-256 hashes remain owned by the pinned upstream dependency, not duplicated in application configuration. Missing target builds, invalid manifests and conflicting library definitions fail the build. Unknown configuration fields also fail; the supported shape is defined in the [configuration type](../source/tools/packer/main.go).

Cross-compilation requires a pure-Go-compatible application. Windows and Linux archives can be built on any supported host. macOS packaging requires a macOS host for Apple tools, even when only changing architectures. The gallery selects `darwin/arm64/dmg` to produce only a DMG. Selecting `darwin/arm64/zip` instead skips disk-image creation and Finder automation, but still requires macOS for icon compilation and signing. When enabled, release DMGs open in a drag-to-install layout with an Applications shortcut and require a logged-in desktop session and permission for the invoking terminal or IDE to control Finder (System Settings → Privacy & Security → Automation); Finder briefly opens the staging image to save its layout before compression. Without this permission, packaging fails rather than publishing an unstyled image. Requests for macOS on another host fail before any target builds. The scope stays small: no universal binaries, Windows installers, Linux desktop registration or distribution packages, resource mappings, frontend pipeline, Developer ID signing, notarization or publishing. macOS bundles are ad-hoc signed; Windows executables are not Authenticode-signed. The configured version is a local placeholder, not a library release. VS Code's [tasks](../.vscode/tasks.json) provide `run` and a default `build` task, both using Packer from `source/`.

Icon sources live at the root of `resources/` beside `packer.yaml`. The default icon source names (including the entire `icon.icon` directory) and explicitly configured icon paths are excluded from payload copying on every platform; only the target's compiled or installed icon is shipped. Nested payload images such as `resources/images/avatar.jpeg` are unaffected. Packer discovers `icon.icon` or `icon.icns` for macOS, `icon.ico` for Windows, and `icon.png` or `icon.svg` for Linux, in that order of preference. Platform-specific `macos.icon`, `windows.icon` and `linux.icon` values can override discovery; `macos.minimumSystemVersion` also belongs in the platform section. Windows ICOs are embedded in the staged EXE, preserving existing resources and supplying a DPI-aware manifest if the application has none. Linux PNG/SVG icons are included beside the executable; the archive does not install a desktop entry. If no icon is configured or discovered for a target, it builds without a custom icon.

macOS accepts only a prebuilt `.icns` file or an Icon Composer `.icon` directory, not PNG. ICNS files are checked for valid container framing and copied unchanged as the bundle icon. `.icon` requires a full Xcode installation with Icon Composer support selected by `xcode-select`; Apple `actool` compiles layered assets and supplies the icon metadata merged into Info.plist without flattening to PNG. The resulting icon files and asset catalogs go into `Contents/Resources`, separately from application payload resources. Select a non-default icon through its platform's `icon` property in `packer.yaml`.

Paths in the configuration are relative to its directory. The gallery uses the layered [Icon Composer project](../source/resources/icon.icon) on macOS, a multi-size [ICO](../source/resources/icon.ico) on Windows and a 1024 × 1024 [PNG](../source/resources/icon.png) on Linux. The ICO and PNG are static exports of the layered icon’s default appearance; refresh both when changing the design. The gallery’s Image example loads the author’s avatar from `resources/images/avatar.jpeg` through MyGo’s resource path; direct `go run` uses the module’s resource directory. The contents of a sibling `resources/` directory preserve their relative paths: under `Contents/Resources` on macOS, beside the executable on Windows and Linux. Resource symlinks, special files and collisions with generated files are rejected. Native libraries go into `Contents/Frameworks` on macOS and are signed before the enclosing app; DLLs and shared objects live beside the executable on Windows and Linux. Only macOS packaging invokes Apple's signing and archive tools.

On macOS, the gallery's Terminal page requires macOS 13 or later. Packer verifies manifest-provided downloads and keeps an independent `onevra-packer/natives` user cache. The first build may need network access; packaged apps do not download missing libraries. MyGo's current terminal loader does not search `Contents/Frameworks`, so Gallery startup selects its bundled library through `MYGO_GHOSTTY_VT`, preserving an explicit user override. This application-side adaptation is not injected by Packer. Direct `go run` and ordinary test binaries retain the plugin's normal lookup and may download the terminal library on first use. A missing bundled library is shown as an error on the page rather than silently using a developer's cache.

For an unpackaged run on any supported desktop, use `go run ./examples/gallery -theme dark`.

## Starter applications

Three independent applications provide small starting points alongside Gallery. From `source/`, run one of:

```sh
GOWORK=off go run ./examples/desktop
GOWORK=off go run ./examples/menubar
GOWORK=off go run ./examples/hybrid
```

| Example | Startup and close behavior | Settings host |
| --- | --- | --- |
| [Desktop](../source/examples/desktop/main.go) | Shows the main window. On macOS, closing it keeps the app in the Dock by default; Settings can change this to quit. Other platforms quit on close. | Dialog inside the main window |
| [Menu bar](../source/examples/menubar/main.go) | Shows a first-launch introduction, then runs from its tray icon without a main window. The quick-note panel is optional. | Independent window |
| [Hybrid](../source/examples/hybrid/main.go) | Shows the main window and tray icon. Optional silent startup leaves the window uncreated until needed. Closing can retain the tray, retain both tray and macOS Dock, or quit. | Dialog inside the main window |

The main, settings and introduction windows use `shells/plain`, with native window controls over an integrated draggable header instead of a default title bar. Tray panels remain frameless.

Opening an already running application activates its existing instance. Desktop and hybrid restore and focus the main window even after silent startup or minimization. Menu bar shows its introduction until acknowledged, then its panel or settings. Explicit Quit always exits, regardless of window-close preferences. The examples use MyGo's single-instance lock and native activation events; they do not implement their own IPC or native runtime.

On macOS and Windows, left-click toggles the quick-note panel and right-click opens settings and quit actions; hybrid also offers Open main window. Disabling the panel makes left-click open settings, an application-owned action that can be replaced when adapting the template. Escape, close and loss of focus hide the panel. Linux uses AppIndicator menu actions instead of direct clicks and requires the desktop support described in [Tray example](#tray-example). If icon creation fails, an accessible window displays the error; closing it exits rather than leaving an inaccessible background process.

English, Simplified Chinese and Follow system apply to application content, settings, window titles and application-owned menus. System language is resolved at startup, with unsupported languages falling back to English. Preferences and the menu bar introduction acknowledgement are saved in `config.yaml` under `os.UserConfigDir()/onevra-examples/<example>/`; each example has its own directory. Writes replace the file atomically. Read failures use defaults and save failures leave changes active for the current session, with feedback in the interface. Note text is intentionally session-only. OS-provided menu items keep their system language.

Each example has a local `packer.yaml`. From `source/`, build Packer once and invoke it beside each configuration:

```sh
GOWORK=off go build -o ../dist/packer ./tools/packer
(cd examples/desktop && ../../../dist/packer build)
(cd examples/menubar && ../../../dist/packer build)
(cd examples/hybrid && ../../../dist/packer build)
```

Replace `build` with `run` for a retained development bundle. Configurations select macOS arm64 ZIPs and write to separate `dist/<example>/` directories; edit their target lists for other platforms. They need no Gallery resources or custom app icon. See [App bundles](#app-bundles) for packaging requirements.

To start a project, copy one complete example directory, including `locales/` and `packer.yaml`, into the application's Go module. Change the application name, package identifier, preference directory and output path; provide module dependencies using this repository's [module guidance](../ARCHITECTURE.md#repository-landmarks). The examples have no imports from Gallery or sibling examples. Translation IDs and stored option values remain stable across language changes; extend both catalogs together.

For native lifecycle checks and settings captures, run from `source/` in a desktop session:

```sh
ONEVRA_NATIVE_SMOKE=1 ONEVRA_SCREENSHOTS="$PWD/../dist/template-preview" \
  GOWORK=off go test -p 1 ./examples/desktop ./examples/menubar ./examples/hybrid \
  -run '^TestNative' -count=1 -timeout 45s
```

The tests exercise startup, close-to-hide, restoration, settings hosting, language updates, panel-free and silent startup, and quit cleanup. In a packaged app, also review actual Dock and tray clicks, opening the application again, right-click dismissal followed by left-click, keyboard shortcuts, focus, accessibility, both languages and themes, and panel placement across displays. Programmatic callbacks and content captures do not establish those native interactions. Windows and Linux require their own desktop review.

## Tray example

Run Gallery normally and select **Desktop / Tray**, then **Enable tray**. On macOS and Windows, left-click toggles the notes window and right-click opens a native menu with Open notes, Show Gallery and Quit actions. Escape, window close and loss of focus hide the panel. Disable tray removes the icon and panel; closing Gallery also releases both. Notes remain in memory until Gallery exits.

Linux needs `libayatana-appindicator3` and a desktop supporting AppIndicator. AppIndicator provides menu actions, not direct click events, so use **Open notes**. The panel centers when the platform cannot report tray bounds. The tray icon is generated in memory and needs no additional packaged resource.

From `source/`, run `ONEVRA_NATIVE_SMOKE=1 GOWORK=off go test ./examples/gallery/internal/gallery -run '^TestNativeTray$' -count=1 -timeout 30s` in a desktop session for the lifecycle smoke check. It checks enable/disable, show/hide, close-to-hide and cleanup, not pointer interaction or appearance. Manually verify repeated left-click toggling, right-click menu dismissal followed by another left-click, menu actions, Escape, click-away dismissal, retained text, keyboard focus, both system themes, display edges and multiple displays. Windows and Linux require separate native review.

## Dock menu example

On macOS, select **Desktop / Dock menu** in Gallery and click **Enable Dock menu**. Right-click the application's Dock icon to find Show Gallery, Open Welcome and Open Settings above the standard macOS items. Disable Dock menu removes the demo shortcuts; closing Gallery also restores the previous menu. Windows and Linux show a platform explanation rather than an enable action.

From `source/`, run `ONEVRA_NATIVE_SMOKE=1 GOWORK=off go test ./examples/gallery/internal/gallery -run '^TestNativeDockMenu$' -count=1 -timeout 30s` in a macOS desktop session. This checks menu installation, callback behavior and cleanup, not Dock pointer interaction. Manually verify the actual Dock menu, all three shortcuts with Gallery visible, hidden or minimized, the Settings dialog, and removal of the custom items without changing standard macOS actions. A packaged Gallery run is recommended to review the Dock icon's application identity.

## Notification example

Select **Desktop / System notifications** and click **Send notification**. On macOS, use `go run ./tools/packer run` so Gallery has an application bundle, and allow notifications when prompted. An unpackaged executable reports that notifications are unavailable. Linux requires a desktop notification service; Windows uses notification-area balloons.

Review delivery with sound enabled and disabled, denied permission, repeated sends, withdrawal, and clicking a notification while Gallery is hidden or minimized or showing another page. Clicking restores the notification page; closing Gallery withdraws the demo notification. System notification settings and Do Not Disturb can suppress presentation even after successful submission. Check both interface languages and themes, narrow layouts and keyboard navigation. Windows and Linux need their own native review.

## Development checks

```sh
gofmt -l .
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build -o ../dist/onevra-gallery ./examples/gallery
(cd .. && git diff --check)
```

The formatting check should produce no output. Use a `.exe` output suffix on Windows. `GOWORK=off` checks that the module works independently of external workspaces.

## Native review

From `source/`, in a macOS desktop session:

```sh
ONEVRA_NATIVE_SMOKE=1 ONEVRA_SCREENSHOTS="$PWD/../dist/preview" \
  go test ./examples/gallery/internal/gallery -run '^TestNative(Gallery|MyGo|Controls|Terminal)$' -count=1 -v
```

The tests open windows, capture appearance variants (including every Onevra control and MyGo example in light, dark and narrow layouts), exercise terminal input and cleanup, check that the settings dialog preserves window bounds, and close the windows. Omit `ONEVRA_SCREENSHOTS` to use temporary output.

For command-palette captures, use the same environment with `-run '^TestNativeCommand(Expansion|InputSource)?$'`. The macOS input-source test requires ABC and another enabled input source, temporarily switches between them, and restores the original source; it checks automatic selection, restoration, manual overrides and retaining ABC on close. In Gallery, review Command+E (Ctrl+E on other platforms), typing and Chinese input-method composition, arrow navigation, single-click and Enter activation, Escape dismissal, restored focus, and opening Settings from a command. Check both interface languages and themes, narrow windows and screen-reader announcements; captures alone do not verify these interactions.

For packaging changes, run `ONEVRA_PACKER_SMOKE=1 GOWORK=off go test ./tools/packer -run '^Test(Native|Cross)Packaging$' -count=1 -v`. Cross tests inspect Windows PE headers, embedded ICO/manifest resources, Linux ELF headers, executable permissions and target-specific native payloads for both architectures; they do not execute foreign binaries. The native packaging test uses the real Apple tools on macOS and skips elsewhere. This checks prebuilt ICNS packaging and Icon Composer compilation, payload copying, bundle metadata, signing, versioned release archive outputs, DMG integrity and the retained development bundle using an independent Go fixture. To verify the terminal library from a built bundle, set `ONEVRA_GALLERY_BUNDLE` to an absolute path to a development app or an extracted release app and run `go test ./examples/gallery -run '^TestBundledTerminalLoad$' -count=1 -v`. This loads a relocated copy from Frameworks rather than the cache. Also run the packaged gallery in a desktop session; archive checks do not establish Finder or Dock appearance.

For UI changes, manually review both themes, narrow layouts, text scaling, window controls, keyboard navigation and accessibility. Verify that sidebar dragging respects configured width limits and double-clicking restores the default width without resizing the window. Inspect background blending over a contrasting backdrop. Captures do not replace interaction review; Windows and Linux need separate desktop checks. See the [gallery verification approach](architecture/gallery.md#verification-approach).

## Updating Lucide icons

Ordinary builds use checked-in generated Go assets and do not need Node or a Lucide checkout. To update, obtain a clean checkout of [Lucide](https://github.com/lucide-icons/lucide) and explicitly select its full commit hash. From `source/`:

```sh
go run ./tools/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH
# Repeat without writing to check reproducibility:
go run ./tools/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH -check
```

The generator validates every regular SVG with MyGo and records the commit and input digest in [upstream.json](../source/icons/lucide/upstream.json). It copies the full upstream license and reports accessor additions and removals. Review the generated diff for renames or removals before accepting an update: upstream changes can affect this library's public API. The test-only collection exercises every accessor without introducing a runtime registry.

The same invocation generates the gallery-only catalog directly from the upstream icon list, in filename order; `-check` verifies it alongside the library assets. The catalog is checked in and is not regenerated during ordinary gallery runs or packaging. Run the development checks and review the gallery's icon grid and icon samples in both appearances, including small sizes and complex strokes. Parsing tests alone do not prove visual fidelity. When changing the cache or generated representation, compare minimal named-icon and full-collection builds to check unused-resource elimination.

## Publication

The library is currently local and unpublished. Consumers may use a local Go workspace containing their application and this repository's `source/` directory; do not commit machine-specific paths.

Before remote publication, resolve the module-path or repository-layout choice described in the [architecture overview](../ARCHITECTURE.md#repository-landmarks).

For an explicitly authorized release, review the shipping scope in [CHANGELOG.md](../CHANGELOG.md), move only shipping user-facing entries into the confirmed version, and retain `Unreleased` for pending changes. There is no release automation. Versioning, tagging, publishing and pushing require separate authorization.
