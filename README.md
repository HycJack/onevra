# oneo-go

**English** | [简体中文](README.zh-CN.md)

Experimental, product-neutral Go components built on [MyGo](https://github.com/egoist/mygo). Currently provides a native Shell component, Lucide icons for MyGo and a standalone gallery.

The library is not yet published as a downloadable module, and its API is not stable.

`main` uses upstream MyGo without custom patches. The WebView composition experiment is preserved on `feature/webview` because it requires a customized MyGo dependency; it is not part of the mainline components.

## Try the gallery

Requires Go 1.27.1 or later and a native desktop session. Linux also requires GTK 3, Pango and OpenGL.

From the repository root:

```sh
cd source
go run ./examples/gallery
```

The sidebar includes Shell, the Lucide icon library, oneo-go controls and MyGo components. The **oneo-go controls** group demonstrates each public control with interactive state variants, including all button styles. MyGo components are grouped by type; each page displays interactive component states and variations without code or tutorial text. The Plugins / Terminal page offers ANSI output and input echo without launching a shell; see [runtime requirements](docs/build.md#app-bundles).

To preview a specific appearance, append `-theme light` or `-theme dark`. `go run ./tools/packer run` compiles and starts a packaged development app without hot reload; `go run ./tools/packer build` produces release archives for the targets selected in `source/packer.json`. See [packaging and icon requirements](docs/build.md#app-bundles).

## Basic usage

Until publication, use a local Go workspace that includes your application and this repository's `source/` directory. Import `github.com/oneo-me/oneo-go/shell`:

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    oneoshell "github.com/oneo-me/oneo-go/shell"
)

func main() {
    mygo.App.WhenReady(func() {
        shell := oneoshell.NewShell(oneoshell.ShellOptions{
            Title: "Hello",
            Left: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(oneoshell.HeaderHeight(c), 0, 0, 0).Children(func() {
                    ui.Text(c, "Navigation")
                })
            },
            Content: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(oneoshell.HeaderHeight(c), 0, 0, 0).Children(func() {
                    ui.Text(c, "Hello, Shell")
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

Create a Shell when the app is ready; build its columns with MyGo widgets. Use `shell.SetPage(view)` to replace the entire sidebar-and-content page (for example, with settings), and `shell.SetPage(nil)` to restore it. Your application owns business state and lifecycle handling. See the [runnable gallery](source/examples/gallery/main.go) for another example.

Content extends behind Shell's floating header. Use `oneoshell.HeaderHeight(c)` for a top inset when needed; paint the padded container's background to cover the titlebar area without placing controls under it. This also applies to replacement pages.

## Appearance

Shell uses oneo-go's built-in light and dark themes, applied through MyGo's own theme mechanism so every MyGo widget keeps working. `ShellOptions.Appearance` selects `theme.Auto` (the default, follows the system appearance), `theme.Light` or `theme.Dark`:

```go
import "github.com/oneo-me/oneo-go/theme"

shell := oneoshell.NewShell(oneoshell.ShellOptions{
    Title:      "Hello",
    Appearance: theme.Dark,
    // ...
})
```

The palettes follow the Light and Shadow (light-shadow) design style from [oneo-me/desgin-styles](https://github.com/oneo-me/desgin-styles): a warm white light mode and a deep navy dark mode. See the [theme package](source/theme/doc.go) and [appearance design](docs/architecture/components.md#native-appearance).

## Controls

Import `github.com/oneo-me/oneo-go/ui` for oneo-styled common controls built on MyGo's unstyled bases: buttons, check boxes, switches, radio buttons, text fields and areas, selects, sliders and tabs. Use the `oneoui` import alias to distinguish it from MyGo's `ui`. Colors come from the active theme:

```go
if oneoui.Button(c, "Save", oneoui.ButtonPrimary).Clicked() {
    // ...
}

oneoui.Checkbox(c, &agreed, "I agree")
oneoui.Switch(c, &wifi, "Wi-Fi")
```

Keyboard, pointer and accessibility behavior come from MyGo's bases. See the [package contract](source/ui/doc.go).

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
