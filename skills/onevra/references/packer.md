# Project layout and Packer

Packer is Onevra's standalone desktop runner and packager. Applications need not use Onevra components to use it. It replaces MyGo's packaging commands: do not add `mygo.json` or MyGo packaging configuration.

## Packer requirements

Packer reads `packer.yaml` from its working directory and accepts only `run` or `build`, with no flags or forwarded application arguments. Unknown fields fail. Paths are relative to the configuration's directory.

- `packer.yaml` sits at the Go module root, beside `go.mod`.
- `main` selects the main package (default `.`). Compilation and dependency discovery use the target `GOOS`/`GOARCH` with cgo disabled, so the application must be pure-Go-compatible for cross-builds.
- `resources/` beside `packer.yaml` is the payload, copied with relative paths preserved into `Contents/Resources` on macOS or beside the executable elsewhere. Symlinks, special files and collisions with generated files are rejected.
- Icons live at the root of `resources/` and are discovered by name: `icon.icon` (Icon Composer, needs full Xcode) or `icon.icns` for macOS, `icon.ico` for Windows, `icon.png` or `icon.svg` for Linux. They are not copied as payload. Override with `macos.icon`, `windows.icon` or `linux.icon`.
- Plugin native libraries come from imported packages' `mygo-plugin.json` manifests; do not list them in configuration.
- Application-owned compiled libraries go in `runtimes/<GOOS>/<GOARCH>/` beside `packer.yaml`, flat files only, with no YAML declaration. Packer copies the selected target's files beside the executable or into `Contents/Frameworks` and signs them on macOS. Build them before packaging; Packer does not compile them, rewrite loader paths or collect transitive dependencies. Keep them out of `resources/`.

## Layout

Follow Onevra's repository shape: documentation at the repository root, and everything Go, including packaging configuration, in `source/`.

```text
<app>/
├── .gitignore            # /dist/
├── AGENTS.md
├── ARCHITECTURE.md
├── README.md
├── docs/
└── source/               # Go module root; run Go and Packer here
    ├── go.mod            # module <repository path>/source
    ├── packer.yaml
    ├── main.go           # or a main package selected by `main`
    ├── locales/          # embedded translation catalogs
    ├── runtimes/         # optional compiled native libraries
    │   └── darwin/arm64/libexample.dylib
    ├── resources/
    │   ├── icon.icon/    # macOS
    │   ├── icon.ico      # Windows
    │   ├── icon.png      # Linux
    │   └── images/…      # runtime payload, loaded through MyGo's resource path
    └── internal/…        # application packages
```

When the module is published, its path includes `/source` and version tags use the `source/` prefix, such as `source/v0.1.0`.

## Configuration

```yaml
name: My App                      # application and bundle name; not a path
identifier: com.example.myapp
version: "0.1.0"
main: .
out: ../dist                      # repository-root dist/, ignored by Git
platforms:                        # GOOS/GOARCH[/FORMAT]; omitted means the host only
  - darwin/arm64/dmg
  - windows/amd64
  - linux/amd64
macos:
  minimumSystemVersion: "13.0"
```

Formats are `dmg` or `zip` on macOS, `zip` on Windows and `tar.gz` on Linux; without one, macOS produces both DMG and ZIP. `openReleases: false` keeps a build from opening the release folder, as CI and headless sessions need. `keepUncompressed: false` keeps only archives. The supported fields are the `config` type in `tools/packer/main.go` of the Onevra module.

Declare macOS usage descriptions and other application-owned metadata in `macos.infoPlist`; Packer writes it before signing in both modes. Generated bundle and icon fields are protected. The mapping preserves scalar and nested collection types; see the [metadata contract](https://github.com/oneo-me/onevra/blob/main/docs/build.md#additional-macos-metadata) for accepted values and conflicts. Applications still own runtime authorization; adding a usage description does not grant permission.

## Commands

Install once, then run beside `packer.yaml`:

```sh
go install github.com/oneo-me/onevra/source/tools/packer@latest
cd source
packer run      # host-only development bundle "<name> Dev", launched with terminal I/O
packer build    # release archives in <out>/releases/<version>/
```

- `run` uses a separate `.dev` identifier, replaces the previous development bundle and has no watcher. `build` strips symbols, marks production mode so `mygo.IsDev()` is false, and clears that version's release directory first.
- macOS targets require a macOS host. DMG layout needs a logged-in desktop session and Automation permission to control Finder; choose `/zip` when that is unavailable. Bundles are ad-hoc signed; Developer ID signing, notarization, installers and publishing are outside Packer.
- Keep the application's displayed version aligned with `version` in `packer.yaml`.

For full behavior, see the Onevra [development guide](https://github.com/oneo-me/onevra/blob/main/docs/build.md#app-bundles).
