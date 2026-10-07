# Development

Use Go 1.27.1 or later. Run Go commands from `source/`; dependencies are pinned in [go.mod](../source/go.mod). No sibling project or workspace file is required.

Native windows require a desktop session. Linux also needs GTK 3, Pango and OpenGL system libraries. See the [MyGo native guides](https://github.com/egoist/mygo/tree/v0.2.12/docs/ui) for platform setup.

## App bundles

The local [packer](../source/tools/packer/main.go) builds desktop apps for selected platforms and runs the host version. Run it from `source/`:

```sh
go run ./tools/packer run
go run ./tools/packer build
```

Packer accepts only the `run` and `build` commands, with no options or application-argument forwarding. It reads `packer.json` from the working directory; all packaging settings belong there.

The top-level `platforms` array in [packer.json](../source/packer.json) selects release targets using Go's `GOOS/GOARCH` spelling: `darwin`, `windows` or `linux`, each with `arm64` or `amd64`. The gallery configuration lists all six targets; edit this list to choose a subset, for example `"platforms": ["darwin/arm64", "windows/amd64"]`. If omitted or empty, Packer builds the host platform. Unsupported targets fail before building, and duplicate targets are built once.

`run` always uses the host and ignores the release target list. It compiles once with ordinary Go optimizations and caching, attaches terminal I/O and waits for exit. On macOS it retains `dist/ONEO Go Gallery Dev.app` with a separate `.dev` bundle identifier; Windows and Linux retain a `dist/ONEO Go Gallery Dev/` directory. The next run replaces the development output. There is no file watcher or hot reload.

`build` uses `-trimpath` and linker `-s -w` flags to omit symbol tables and debug information. Windows releases also use the GUI subsystem. Packer does not inject framework-specific variables or change application runtime modes. Outputs go into `dist/releases/<version>/`, with no platform subdirectories or retained release application directories:

| Target | Artifacts |
| --- | --- |
| macOS | `oneo-go-gallery-darwin-<arch>-<version>.dmg` and `.zip`, containing `ONEO Go Gallery.app`; the DMG opens in a large-icon drag-to-install layout with the app on the left and an Applications shortcut on the right |
| Windows | `oneo-go-gallery-windows-<arch>-<version>.zip`, containing the app directory and `.exe` |
| Linux | `oneo-go-gallery-linux-<arch>-<version>.tar.gz`, containing the app directory and executable |

`<version>` comes from `version` in `packer.json`. The version determines the release directory and archive suffix, not the contained application or development output. For example, version `0.2026.1` writes archives to `dist/releases/0.2026.1/`. Different versions coexist in separate directories; rebuilding the same target and version replaces its archives.

Each target finishes compilation and packaging in staging before publication, with backups for rollback. Targets build sequentially: a later target's failure does not remove an earlier target's successful archives. Development outputs and other targets remain untouched.

[packer.json](../source/packer.json) owns application packaging configuration; Packer does not read `mygo.json` or import MyGo Go packages. macOS application metadata goes into Info.plist. Both Go compilation and dependency discovery use the selected target's `GOOS`/`GOARCH` with cgo disabled. For native dependencies, Packer uses `go list` to locate the selected application's imported packages and reads their `mygo-plugin.json` manifests as data. Versions, download URLs and SHA-256 hashes remain owned by the pinned upstream dependency, not duplicated in application configuration. Missing target builds, invalid manifests and conflicting library definitions fail the build. Unknown configuration fields also fail; the supported shape is defined in the [configuration type](../source/tools/packer/main.go).

Cross-compilation requires a pure-Go-compatible application. Windows and Linux archives can be built on any supported host. macOS packaging requires a macOS host for Apple tools, even when only changing architectures. Release DMGs also require a logged-in desktop session and permission for the invoking terminal or IDE to control Finder (System Settings → Privacy & Security → Automation); Finder briefly opens the staging image to save its layout before compression. Without this permission, packaging fails rather than publishing an unstyled image. Requests for macOS on another host fail before any target builds. The scope stays small: no universal binaries, Windows installers, Linux desktop registration or distribution packages, resource mappings, frontend pipeline, Developer ID signing, notarization or publishing. macOS bundles are ad-hoc signed; Windows executables are not Authenticode-signed. The configured version is a local placeholder, not a library release. VS Code's [tasks](../.vscode/tasks.json) provide `run` and a default `build` task, both using Packer from `source/`.

Icon sources live beside `packer.json`, outside `resources/`. Packer discovers `icon.icon` or `icon.icns` for macOS, `icon.ico` for Windows, and `icon.png` or `icon.svg` for Linux, in that order of preference. Platform-specific `macos.icon`, `windows.icon` and `linux.icon` values can override discovery; `macos.minimumSystemVersion` also belongs in the platform section. Windows ICOs are embedded in the staged EXE, preserving existing resources and supplying a DPI-aware manifest if the application has none. Linux PNG/SVG icons are included beside the executable; the archive does not install a desktop entry. If no icon is configured or discovered for a target, it builds without a custom icon.

macOS accepts only a prebuilt `.icns` file or an Icon Composer `.icon` directory, not PNG. ICNS files are checked for valid container framing and copied unchanged as the bundle icon. `.icon` requires a full Xcode installation with Icon Composer support selected by `xcode-select`; Apple `actool` compiles layered assets and supplies the icon metadata merged into Info.plist without flattening to PNG. The resulting icon files and asset catalogs go into `Contents/Resources`, separately from application payload resources. Select a non-default icon through its platform's `icon` property in `packer.json`.

Paths in the configuration are relative to its directory. The gallery uses the layered [Icon Composer project](../source/icon.icon) on macOS, a multi-size [ICO](../source/icon.ico) on Windows and a 1024 × 1024 [PNG](../source/icon.png) on Linux. The ICO and PNG are static exports of the layered icon’s default appearance; refresh both when changing the design. The contents of an optional sibling `resources/` directory preserve their relative paths: under `Contents/Resources` on macOS, beside the executable on Windows and Linux. Resource symlinks, special files and collisions with generated files are rejected. Native libraries go into `Contents/Frameworks` on macOS and are signed before the enclosing app; DLLs and shared objects live beside the executable on Windows and Linux. Only macOS packaging invokes Apple's signing and archive tools.

On macOS, the gallery's Terminal page requires macOS 13 or later. Packer verifies manifest-provided downloads and keeps an independent `oneo-packer/natives` user cache. The first build may need network access; packaged apps do not download missing libraries. MyGo's current terminal loader does not search `Contents/Frameworks`, so Gallery startup selects its bundled library through `MYGO_GHOSTTY_VT`, preserving an explicit user override. This application-side adaptation is not injected by Packer. Direct `go run` and ordinary test binaries retain the plugin's normal lookup and may download the terminal library on first use. A missing bundled library is shown as an error on the page rather than silently using a developer's cache.

For an unpackaged run on any supported desktop, use `go run ./examples/gallery -theme dark`.

## Development checks

```sh
gofmt -l .
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build -o ../dist/oneo-go-gallery ./examples/gallery
(cd .. && git diff --check)
```

The formatting check should produce no output. Use a `.exe` output suffix on Windows. `GOWORK=off` checks that the module works independently of external workspaces.

## Native review

From `source/`, in a macOS desktop session:

```sh
ONEO_GO_NATIVE_SMOKE=1 ONEO_GO_SCREENSHOTS="$PWD/../dist/preview" \
  go test ./examples/gallery/internal/gallery -run '^TestNative(Gallery|MyGo|Controls|Terminal)$' -count=1 -v
```

The tests open windows, capture appearance variants (including every oneo-go control and MyGo example in light, dark and narrow layouts), exercise terminal input and cleanup, check that replacement pages preserve window bounds, and close the windows. Omit `ONEO_GO_SCREENSHOTS` to use temporary output.

For packaging changes, run `ONEO_PACKER_SMOKE=1 GOWORK=off go test ./tools/packer -run '^Test(Native|Cross)Packaging$' -count=1 -v`. Cross tests inspect Windows PE headers, embedded ICO/manifest resources, Linux ELF headers, executable permissions and target-specific native payloads for both architectures; they do not execute foreign binaries. The native packaging test uses the real Apple tools on macOS and skips elsewhere. This checks prebuilt ICNS packaging and Icon Composer compilation, payload copying, bundle metadata, signing, versioned release archive outputs, DMG integrity and the retained development bundle using an independent Go fixture. To verify the terminal library from a built bundle, set `ONEO_GALLERY_BUNDLE` to an absolute path to a development app or an extracted release app and run `go test ./examples/gallery -run '^TestBundledTerminalLoad$' -count=1 -v`. This loads a relocated copy from Frameworks rather than the cache. Also run the packaged gallery in a desktop session; archive checks do not establish Finder or Dock appearance.

For UI changes, manually review both themes, narrow layouts, text scaling, window controls, keyboard navigation and accessibility. Verify that sidebar dragging respects configured width limits and double-clicking restores the default width without resizing the window. Inspect background blending over a contrasting backdrop. Captures do not replace interaction review; Windows and Linux need separate desktop checks. See the [gallery verification approach](architecture/gallery.md#verification-approach).

## Updating Lucide icons

Ordinary builds use checked-in generated Go assets and do not need Node or a Lucide checkout. To update, obtain a clean checkout of [Lucide](https://github.com/lucide-icons/lucide) and explicitly select its full commit hash. From `source/`:

```sh
go run ./tools/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH
# Repeat without writing to check reproducibility:
go run ./tools/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH -check
```

The generator validates every regular SVG with MyGo and records the commit and input digest in [upstream.json](../source/icons/lucide/upstream.json). It copies the full upstream license and reports accessor additions and removals. Review the generated diff for renames or removals before accepting an update: upstream changes can affect this library's public API. The test-only collection exercises every accessor without introducing a runtime registry.

After updating the accessors, run `go generate ./examples/gallery/internal/gallery` to refresh the gallery-only catalog from the checked-in assets. Run the development checks and review the gallery's icon grid and icon samples in both appearances, including small sizes and complex strokes. Parsing tests alone do not prove visual fidelity. When changing the cache or generated representation, compare minimal named-icon and full-collection builds to check unused-resource elimination.

## Publication

The library is currently local and unpublished. Consumers may use a local Go workspace containing their application and this repository's `source/` directory; do not commit machine-specific paths.

Before remote publication, resolve the module-path or repository-layout choice described in the [architecture overview](../ARCHITECTURE.md#repository-landmarks).

For an explicitly authorized release, review the shipping scope in [CHANGELOG.md](../CHANGELOG.md), move only shipping user-facing entries into the confirmed version, and retain `Unreleased` for pending changes. There is no release automation. Versioning, tagging, publishing and pushing require separate authorization.
