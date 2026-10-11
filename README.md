# Onevra

**English** | [简体中文](README.zh-CN.md)

Onevra is a Go component library for native desktop applications, built on [MyGo](https://github.com/egoist/mygo). It provides themed controls, window shells, Lucide icons and complex components. The name combines `One` from **ONEO** with the coined ending `vra`.

![Onevra Gallery main interface with navigation and interactive component examples](docs/images/gallery.png)

## Additional components

Building on MyGo's native UI capabilities, Onevra adds composed components for common application scenarios:

- **Conversation interfaces**: bring message turns, tool activity, folded summaries and attachment input together into a complete chat interface.
- **Markdown rendering**: provide selectable document views with code highlighting, images and custom code-block content.
- **Code highlighting**: provide read-only code views with light and dark themes, usable independently or within Markdown and conversations.

- **Lucide icons**: provide SVG icons for MyGo through named accessors, loaded on demand with configurable sizing and inherited text color.

## Examples

Follow [MyGo’s official requirements](https://mygo.egoist.dev/docs) for the development environment and system dependencies.

From the repository root, start the interactive Gallery:

```sh
cd source
go run ./examples/gallery
```

Append `-theme light` or `-theme dark` to preview an appearance. Three smaller applications are also available; run these commands from `source/`:

| Example | Purpose | Command |
| --- | --- | --- |
| Desktop | A main window with settings | `go run ./examples/desktop` |
| Menu bar | A tray application with a quick-note panel | `go run ./examples/menubar` |
| Hybrid | A main window and tray, with optional background startup | `go run ./examples/hybrid` |

Each template includes localization, preferences and its own packaging configuration. See [starter applications](docs/build.md#starter-applications) for copying and adapting one. Windows and Linux still require their own [native review](docs/build.md#native-review).

## Getting started

Install **0.2026.6** in your application module with:

```sh
go get github.com/oneo-me/onevra/source@v0.2026.6
```

Save this minimal application as `main.go` and run `go run .` from your application module:

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    "github.com/oneo-me/onevra/source/components"
    "github.com/oneo-me/onevra/source/shells/plain"
)

func main() {
    mygo.App.WhenReady(func() {
        shell := plain.NewShell(plain.ShellOptions{
            Title: "Hello Onevra",
            Content: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(24).Children(func() {
                    components.Text(c, "Hello, Onevra!")
                })
            },
        })
        shell.Window().OnClosed(func() { mygo.App.Quit() })
    })
    if err := mygo.App.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Create windows after MyGo is ready. Build content with Onevra controls and MyGo layout primitives; your application owns state, persistence and asynchronous work.

## Packer

Onevra provides a standalone [Packer](source/tools/packer/main.go) as an alternative to MyGo's official packaging commands. It reads `packer.yaml` and builds native application bundles with icons, resources and plugin libraries. Your application does not need to use Onevra components to use Packer.

To run or package Gallery, execute these commands from `source/`:

```sh
go run ./tools/packer run
go run ./tools/packer build
```

`run` builds and launches the host application; `build` creates archives for the configured targets. To use Packer in another project, run `go install ./tools/packer` from `source/`, add Go's binary installation directory to `PATH`, then run `packer run` or `packer build` beside that project's `packer.yaml`.

Supported archives are macOS DMG/ZIP, Windows ZIP and Linux tar.gz. macOS packaging requires a Mac; signing for distribution and notarization remain separate steps. See [packaging configuration and requirements](docs/build.md#app-bundles).

## Shells

Shells create native windows and apply Onevra's theme while leaving application lifecycle and business state to the caller.

| Package | Layout |
| --- | --- |
| [`shells/plain`](source/shells/plain/shell.go) | A single content area with automatic clearance for window controls; used in the example above |
| [`shells/sidebar`](source/shells/sidebar/shell.go) | A collapsible, resizable sidebar with a content area, header actions and replaceable pages |

Import either package under `github.com/oneo-me/onevra/source/`. Both support system, light and dark appearance through the [theme package](source/theme/doc.go). For Sidebar, use `sidebar.HeaderHeight(c)` when content needs to clear the floating header; see [layout guidance](docs/architecture/components.md#page-container-and-layout).

## Icons

The [Lucide package](source/icons/lucide/doc.go) provides named SVG accessors. Import `github.com/oneo-me/onevra/source/icons/lucide` alongside `components`, then use them in a content builder:

```go
components.Icon(c, lucide.Search).Size(20, 20)
```

Icons are parsed on first use and reused. They inherit text color, and the Icon component also accepts a parsed MyGo SVG.

## Components

Import `github.com/oneo-me/onevra/source/components` for themed inputs, selection, navigation, collections and overlays. Inside a content builder, use retained application state such as `agreed` and `wifi`:

```go
if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
    // Save application state.
}
components.Checkbox(c, &agreed, "I agree")
components.Switch(c, &wifi, "Wi-Fi")
```

Components work independently of Shell when the host applies Onevra's theme. See the [component contract](source/components/doc.go) and [theme usage](source/theme/doc.go). Complex components each have a separate package under `components/`: `chart`, `command`, `code`, `markdown` and `chat`. See [complex component composition](docs/architecture/components.md#complex-component-composition).

## Contributing

Read [AGENTS.md](AGENTS.md) and the [architecture overview](ARCHITECTURE.md) before making changes. From `source/`, run:

```sh
gofmt -l .
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build ./...
```

Formatting should produce no output. See [development checks](docs/build.md#development-checks) and [native review](docs/build.md#native-review) for verification. Keep both READMEs equivalent and record consumer-visible changes in [CHANGELOG.md](CHANGELOG.md). Versioning and release procedures are in [Publication](docs/build.md#publication).

## License

Project code: [MIT](LICENSE). Components adapted from MyGo retain the [upstream MIT notice](source/components/LICENSE.mygo). Lucide assets retain their [upstream ISC and Feather MIT notices](source/icons/lucide/LICENSE).
