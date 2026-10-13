# Project layout and packaging

Package applications with MyGo's own CLI and `mygo.json`; Onevra has no packager of its own. MyGo's [configuration](https://mygo.egoist.dev/docs/configuration) and [distribution](https://mygo.egoist.dev/docs/distribution) guides are authoritative; this page records only the Onevra project conventions.

## Requirements

- Pin the CLI as a Go tool in the application module: `go get -tool github.com/egoist/mygo/cmd/mygo`, at the MyGo version the module already requires.
- `mygo.json` sits at the Go module root, beside `go.mod`. `main` selects the main package (default `.`). Builds disable cgo, so the application must be pure-Go-compatible for cross-builds.
- Everything in `resources/` ships as payload, including platform directories such as `resources/darwin-arm64/` for compiled libraries and helpers of one target. Keep icon sources outside it.
- `icon` is a square PNG, ideally 1024 × 1024, from which MyGo derives every platform's icon. `macos.icon` adds an Icon Composer `.icon` (compiled with Xcode 26 or later) or a prebuilt `Assets.car`.
- Plugin native libraries come from imported packages' `mygo-plugin.json` manifests; do not list them in configuration.

## Layout

Follow Onevra's repository shape: documentation at the repository root, and everything Go, including packaging configuration, in `source/`.

```text
<app>/
├── .gitignore            # /dist/
├── AGENTS.md
├── ARCHITECTURE.md
├── README.md
├── docs/
└── source/               # Go module root; run Go and MyGo here
    ├── go.mod            # module <repository path>/source; tool github.com/egoist/mygo/cmd/mygo
    ├── mygo.json
    ├── main.go           # or a main package selected by `main`
    ├── assets/
    │   ├── icon.png
    │   └── icon.icon/    # optional layered macOS icon
    ├── locales/          # embedded translation catalogs
    ├── resources/        # runtime payload, loaded through MyGo's resource path
    └── internal/…        # application packages
```

When the module is published, its path includes `/source` and version tags use the `source/` prefix, such as `source/v0.1.0`.

## Configuration

```json
{
  "name": "My App",
  "identifier": "com.example.myapp",
  "version": "0.1.0",
  "icon": "assets/icon.png",
  "out": "../dist",
  "macos": {
    "icon": "assets/icon.icon",
    "minimumSystemVersion": "13.0"
  }
}
```

Declare macOS usage descriptions and other application-owned metadata in `macos.infoPlist`. Applications still own runtime authorization; adding a usage description does not grant permission.

## Commands

```sh
cd source
go tool mygo dev                                          # "<name> Dev" with live reload
go tool mygo build -platform darwin/arm64,windows/amd64   # outputs in <out>/<os>-<arch>/
```

Keep the application's displayed version aligned with `version` in `mygo.json`.
