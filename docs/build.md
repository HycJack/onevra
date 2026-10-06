# Running and building

Use Go 1.27.1 or later. Run Go commands from `source/` (`cd source` from the repository root); dependencies are pinned in `source/go.mod` and `source/go.sum`. No sibling project, workspace file, JavaScript toolchain or cgo is required.

```sh
go tool mygo dev
go run ./examples/gallery -theme dark
go build -o ../dist/oneo-go-gallery ./examples/gallery
```

Use `../dist/oneo-go-gallery.exe` on Windows. MyGo uses Metal on macOS, Direct3D 11 on Windows and OpenGL on Linux. Native windows need a desktop session; Linux additionally needs GTK 3, Pango and OpenGL system libraries. Refer to the pinned [MyGo native guides](https://github.com/egoist/mygo/tree/v0.2.12/docs/ui).

## App bundles

`source/mygo.json` selects `./examples/gallery` as the main package, leaving the module-root package a library. The Go tool directive pins the CLI to the same MyGo dependency as the components; no global CLI installation is needed. VS Code's `dev` and `build` tasks run these commands from `${workspaceFolder}/source`.

```sh
go tool mygo dev
go tool mygo build
# Optional: both Mac architectures in one application.
go tool mygo build -platform darwin/universal
```

On macOS, development runs `source/.mygo/dev/darwin-<arch>/oneo-go Gallery Dev.app` and rebuilds/relaunches it when inputs change. Quit the app or interrupt the command to stop. The development bundle uses the `.dev` identifier suffix. `go run` remains available for direct executable runs with `-theme` flags; it does not create a bundle.

Build defaults to the host platform. On macOS it writes `dist/darwin-<arch>/oneo-go Gallery.app` and a disk image. Local bundles are ad-hoc signed, not Developer ID signed or notarized for distribution. The configured `0.0.0` is a local gallery bundle placeholder, not a library release. MyGo picks up the app icon from `source/resources/icon.png`; set `icon` in the configuration to use another PNG. `source/.mygo/` and the repository-root `dist/` are ignored build outputs. No upload or publication is performed.

## Development checks

```sh
gofmt -w .
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build -o ../dist/oneo-go-gallery ./examples/gallery
(cd .. && git diff --check)
```

A clean formatting check is `gofmt -l .` with no output. GOWORK=off verifies independent module operation even if a contributor uses an external workspace.

Documentation changes require size-budget, Subjects-coverage, relative-link, heading-anchor and changelog checks, plus semantic review against implementation. The overview budget is 16 KiB; subjects are 24 KiB each.

## Native review

```sh
ONEO_GO_NATIVE_SMOKE=1 ONEO_GO_SCREENSHOTS="$PWD/../dist/preview" \
  go test ./internal/gallery -run '^TestNativeGallery$' -count=1 -v
```

This opens a real window, captures light/dark/narrow states and closes it. Default output uses a temporary test directory. macOS native rendering has been exercised; cross-compilation does not establish Windows/Linux acceptance. See [gallery verification](architecture/gallery.md#tests-and-captures) for manual checks.

## Updating Lucide icons

Ordinary builds use checked-in generated Go assets and do not need Node or a Lucide checkout. To update, obtain a clean checkout of [Lucide](https://github.com/lucide-icons/lucide) and explicitly select its full commit hash. From `source/`:

```sh
go run ./internal/cmd/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH
# Repeat without writing to check reproducibility:
go run ./internal/cmd/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH -check
```

The generator validates every regular SVG with MyGo and records the commit and input digest in [upstream.json](../source/icons/lucide/upstream.json). It copies the full upstream license and reports accessor additions and removals. Review the generated diff for renames or removals before accepting an update: upstream changes can affect this library's public API. The test-only collection exercises every accessor without introducing a runtime registry.

Run the development checks and review the gallery's icon samples in both appearances, including small sizes and complex strokes. Parsing tests alone do not prove visual fidelity. When changing the cache or generated representation, compare minimal named-icon and full-collection builds to check unused-resource elimination.

## Consumer integration

The module path is `github.com/oneo-me/oneo-go`, package name `oneogo`. No remote or released version is configured by this scaffold. The module currently lives in a repository subdirectory while retaining its original import path for local use. Remote publication requires a separate module-path or repository-layout decision; this layout is not yet a downloadable module.

Before publication, consumers may opt into their own local go.work referencing their module and this repository's `source/` directory. Do not commit machine-specific paths or require a particular sibling directory in this repository. Locus is not currently a consumer.

## Release notes

[CHANGELOG.md](../CHANGELOG.md) owns API and gallery changes for library users. There is no release automation. Before an explicitly authorized release, confirm the shipping scope, move completed Unreleased entries into the confirmed version, and retain Unreleased with any deferred items. Do not invent version numbers or publish as part of routine development.
