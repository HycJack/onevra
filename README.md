# oneo-go

**English** | [简体中文](README.zh-CN.md)

Experimental, product-neutral compositions for MyGo v0.2.12. The package name is `oneogo`; the module path is `github.com/oneo-me/oneo-go`. It is an independent repository and is not yet published as a downloadable module. Public API stability is not promised.

## Run the gallery

With Go 1.27.1 or later, enter the source module from the repository root:

```sh
cd source
go tool mygo dev
go run ./examples/gallery -theme dark
go run ./examples/gallery -theme light
```

On macOS, `go tool mygo dev` runs a development `.app` with live reload. `go tool mygo build` produces a local `.app` and `.dmg` under the repository-root `dist/`. The `go run` commands above are direct executable runs for theme checks, not app bundles. See [build details](docs/build.md#app-bundles).

A native desktop session is required. Linux additionally requires GTK 3, Pango and OpenGL. The gallery writes no application settings or content files.

## Components

- `Section(c, title, description, content)`: heading and optional content.
- `EmptyState(c, EmptyStateOptions{...})`: explanation and optional actions.
- `Notice(c, NoticeOptions{...})`: persistent information, warning, error or success feedback.

Components return `*ui.Element` for native layout modifiers and keys. Callbacks build content synchronously during the view. Callers own state, localized text and actions. Use distinct keys for repeated components; internal content uses MyGo's normal positional identity.

```go
oneogo.EmptyState(c, oneogo.EmptyStateOptions{
    Title: "No items",
    Description: "Create an item to get started.",
    Actions: func() {
        if ui.Button(c, "Create").Clicked() {
            // Update caller-owned state.
        }
    },
})
```

Import `github.com/egoist/mygo/ui` and alias `github.com/oneo-me/oneo-go` as `oneogo`. All components use the current MyGo theme; configure it through MyGo, not a parallel theme manager. Notice is a labeled group, not an automatic live alert. Supply a meaningful status title rather than relying on color alone.

## Checks

Run from `source/`:

```sh
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build -o ../dist/oneo-go-gallery ./examples/gallery
ONEO_GO_NATIVE_SMOKE=1 go test ./internal/gallery -run '^TestNativeGallery$' -count=1 -v
```

Native smoke is opt-in and captures light, dark and narrow windows, then closes them. Set `ONEO_GO_SCREENSHOTS` to retain PNGs. Headless tests cover composition, actions, disabled controls, themes, text scaling and narrow layouts. Captures establish rendering, not complete keyboard, screen-reader or other-platform acceptance.

See [architecture](ARCHITECTURE.md), [development](docs/build.md) and [release notes](CHANGELOG.md). Licensed under [MIT](LICENSE).
