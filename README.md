# oneo-go

**English** | [简体中文](README.zh-CN.md)

Experimental, product-neutral Go components built on [MyGo](https://github.com/egoist/mygo). Currently provides a native Shell component, Lucide icons for MyGo and a standalone gallery.

The library is not yet published as a downloadable module, and its API is not stable.

## Try the gallery

Requires Go 1.27.1 or later and a native desktop session. Linux also requires GTK 3, Pango and OpenGL.

From the repository root:

```sh
cd source
go tool mygo dev
```

The sidebar includes Shell, the Lucide icon library and MyGo components. MyGo components are grouped by type; each page displays interactive component states and variations without code or tutorial text.

On macOS this runs a live-reloading development app. To preview a specific appearance, use `go run ./examples/gallery -theme light` or `-theme dark`.

## Basic usage

Until publication, use a local Go workspace that includes your application and this repository's `source/` directory. Import `github.com/oneo-me/oneo-go/components`:

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    "github.com/oneo-me/oneo-go/components"
)

func main() {
    mygo.App.WhenReady(func() {
        shell := components.NewShell(components.ShellOptions{
            Title: "Hello",
            Left: func(c *ui.Context) {
                ui.Text(c, "Navigation")
            },
            Content: func(c *ui.Context) {
                ui.Text(c, "Hello, Shell")
            },
        })
        shell.Window().OnClosed(func() { mygo.App.Quit() })
    })
    if err := mygo.App.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Create a Shell when the app is ready; build its columns with MyGo widgets. Supply `Right` to enable the optional inspector. Use `shell.SetPage(view)` to replace the entire three-column page (for example, with settings), and `shell.SetPage(nil)` to restore it. Your application owns business state and lifecycle handling. See the [runnable gallery](source/examples/gallery/main.go) for another example.

## Lucide icons

Import `github.com/oneo-me/oneo-go/icons/lucide` and use named accessors directly with MyGo:

```go
ui.Icon(c, lucide.Search()).Size(20, 20)
ui.SidebarItem(c, "settings", lucide.Settings(), "Settings")
```

Icons are parsed on first use and reused. MyGo controls sizing, text-color inheritance and accessibility; no component wrapper is required. The package includes the pinned upstream's regular icons, not lab icons or aliases. See the [package contract](source/icons/lucide/doc.go) and [update workflow](docs/build.md#updating-lucide-icons).

## Contributing

Read the [architecture overview](ARCHITECTURE.md) before changing code. Run the basic checks from `source/`:

```sh
GOWORK=off go test ./...
GOWORK=off go vet ./...
```

See [development](docs/build.md) for packaging and native review, and [release notes](CHANGELOG.md) for user-facing updates.

## License

Project code: [MIT](LICENSE). Lucide assets retain their [upstream ISC and Feather MIT notices](source/icons/lucide/LICENSE).
