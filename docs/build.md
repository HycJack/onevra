# Development

Use Go 1.27.1 or later. Run Go commands from `source/`; dependencies and the MyGo CLI are pinned in [go.mod](../source/go.mod). No sibling project or workspace file is required.

Native windows require a desktop session. Linux also needs GTK 3, Pango and OpenGL system libraries. See the [MyGo native guides](https://github.com/egoist/mygo/tree/v0.2.12/docs/ui) for platform setup.

## App bundles

```sh
go tool mygo dev
go tool mygo build
# Optional universal macOS application:
go tool mygo build -platform darwin/universal
```

On macOS, development runs a live-reloading app; builds produce an app bundle and disk image under the repository-root `dist/`. For a direct run without packaging, use `go run ./examples/gallery -theme dark`.

[mygo.json](../source/mygo.json) owns gallery packaging configuration. Local macOS bundles are ad-hoc signed, not Developer ID signed or notarized for distribution. The configured version is a local placeholder, not a library release. These commands do not publish anything.

The gallery's Terminal page uses MyGo's terminal plugin and requires macOS 13 or later. The MyGo CLI bundles its pinned `libghostty-vt` library. With direct `go run` or opt-in native tests, the plugin may download and verify that library into MyGo's user cache on first use; `MYGO_GHOSTTY_VT` can select an existing compatible library. A missing library is shown as an error on the page, not replaced with a simulated terminal.

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
  go test ./internal/gallery -run '^TestNative(Gallery|ShellResize|MyGo|Controls|Terminal)$' -count=1 -v
```

The tests open windows, capture appearance variants (including every oneo-go control and MyGo example in light, dark and narrow layouts), exercise terminal input and cleanup, check that replacement pages preserve window bounds, and close the windows. On macOS ARM64, the sidebar resize test sends native pointer events to verify configured width limits and double-click reset without resizing the window. Omit `ONEO_GO_SCREENSHOTS` to use temporary output.

For UI changes, manually review both themes, narrow layouts, text scaling, window controls, keyboard navigation and accessibility. Inspect background blending over a contrasting backdrop. Captures do not replace interaction review; Windows and Linux need separate desktop checks. See the [gallery verification approach](architecture/gallery.md#verification-approach).

## Updating Lucide icons

Ordinary builds use checked-in generated Go assets and do not need Node or a Lucide checkout. To update, obtain a clean checkout of [Lucide](https://github.com/lucide-icons/lucide) and explicitly select its full commit hash. From `source/`:

```sh
go run ./internal/cmd/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH
# Repeat without writing to check reproducibility:
go run ./internal/cmd/lucidegen -upstream /path/to/lucide -commit FULL_COMMIT_HASH -check
```

The generator validates every regular SVG with MyGo and records the commit and input digest in [upstream.json](../source/icons/lucide/upstream.json). It copies the full upstream license and reports accessor additions and removals. Review the generated diff for renames or removals before accepting an update: upstream changes can affect this library's public API. The test-only collection exercises every accessor without introducing a runtime registry.

After updating the accessors, run `go generate ./internal/gallery` to refresh the gallery-only catalog from the checked-in assets. Run the development checks and review the gallery's searchable icon grid and icon samples in both appearances, including small sizes and complex strokes. Parsing tests alone do not prove visual fidelity. When changing the cache or generated representation, compare minimal named-icon and full-collection builds to check unused-resource elimination.

## Publication

The library is currently local and unpublished. Consumers may use a local Go workspace containing their application and this repository's `source/` directory; do not commit machine-specific paths.

Before remote publication, resolve the module-path or repository-layout choice described in the [architecture overview](../ARCHITECTURE.md#repository-landmarks).

For an explicitly authorized release, review the shipping scope in [CHANGELOG.md](../CHANGELOG.md), move only shipping user-facing entries into the confirmed version, and retain `Unreleased` for pending changes. There is no release automation. Versioning, tagging, publishing and pushing require separate authorization.
