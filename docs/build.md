# Development

Use Go 1.27.1 or later. Run Go commands from `source/`; dependencies and the MyGo CLI are pinned in [go.mod](../source/go.mod). No sibling project or workspace file is required.

Native windows require a desktop session. Linux also needs GTK 3, Pango and OpenGL system libraries. See the [MyGo native guides](https://github.com/egoist/mygo/tree/v0.2.12/docs/ui) for platform setup.

## Dependency preparation

After checkout, and after changes to the pinned MyGo patch, run:

```sh
go run ./internal/cmd/prepare
```

This bootstrap command needs only the standard library and repository patch files. It downloads pinned MyGo, checks patch context exactly, and prepares the ignored `source/.mygo/mygo/` dependency used by the relative replacement in `go.mod`. It does not modify the module cache or vendor upstream source. An unchanged patch is reused; failed patch application leaves the previous dependency intact. Run preparation before, not concurrently with, builds.

Gallery, tests and MyGo packaging then use ordinary commands with no special build tag or alternative example. For local consumers, include this prepared module in their Go workspace, or select the same patched MyGo dependency in their main module.

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
  go test ./internal/gallery -run '^TestNative(Gallery|ShellResize|MyGo|Terminal)$' -count=1 -v
```

The tests open windows, capture appearance variants (including every MyGo example in light, dark and narrow layouts), exercise terminal input and cleanup, check that replacement pages preserve window bounds, and close the windows. On macOS ARM64, the sidebar resize test sends native pointer events to verify configured width limits and double-click reset without resizing the window. Omit `ONEO_GO_SCREENSHOTS` to use temporary output.

For UI changes, manually review both themes, narrow layouts, text scaling, window controls, keyboard navigation and accessibility. Inspect background blending over a contrasting backdrop. Captures do not replace interaction review; Windows and Linux need separate desktop checks. See the [gallery verification approach](architecture/gallery.md#verification-approach).

## Hybrid experiment

Open the existing **Experiments / Native + WebView** page with `go tool mygo dev` or `go run ./examples/gallery` after [dependency preparation](#dependency-preparation). This page uses the new GPU-composited WebView directly; there is no separate hybrid example or legacy build mode.

The bridge has macOS and Windows amd64/arm64 implementations. macOS ARM64 has native click, basic typing, overlay blocking, modal focus recovery, clipping, state/lifecycle and appearance checks. Windows and macOS amd64 require target-desktop review; compilation is not runtime acceptance. Windows needs WebView2 Runtime. [Remaining limitations](architecture/components.md#embedded-webview) apply to both bridges.

Run native checks sequentially on macOS ARM64:

```sh
ONEO_GO_NATIVE_SMOKE=1 go test ./components -run '^TestNative(WebView|HybridComposition)$' -count=1 -v
ONEO_GO_NATIVE_SMOKE=1 go test ./internal/gallery -run '^TestNativeHybrid$' -count=1 -v
go test github.com/egoist/mygo/ui github.com/egoist/mygo/internal/gpu/metal github.com/egoist/mygo/internal/raster
```

Optionally set `ONEO_GO_SCREENSHOTS="$PWD/../dist/gallery-webview"` on the gallery test. OS-composited captures use macOS `screencapture` and require Screen Recording permission. MyGo's `CapturePage` alone omits the browser.

For Windows native review, in PowerShell from `source/`:

```powershell
$env:ONEO_GO_NATIVE_SMOKE = '1'
go test ./components -run '^TestNativeHybridComposition$' -count=1 -v
Remove-Item Env:ONEO_GO_NATIVE_SMOKE
go test github.com/egoist/mygo/internal/gpu/d3d11 -run '^TestNativeCompositionGPU$' -count=1 -v
```

The Windows native test moves the pointer and sends real clicks; run it in an otherwise idle desktop session. After preparation on the host, cross-compilation uses ordinary Go commands:

```sh
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ../dist/oneo-go-gallery-windows-amd64.exe ./examples/gallery
```

Review sidebar resizing, both inputs, scrolling, appearance, Go-to-JavaScript updates, hide/show, navigation and Settings. The two alerts compare a live browser underneath the modal with explicit hiding. The live-browser alert retains its eight-second escape hatch. Check IME, pointer capture during drags, keyboard navigation, accessibility, DPI changes and device-loss recovery on each target; screenshots and headless tests do not establish these behaviors.

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
