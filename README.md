# Onevra

**English** | [简体中文](README.zh-CN.md)

Experimental, product-neutral Go components built on [MyGo](https://github.com/egoist/mygo). Provides themed controls, native window shells, command palettes, Markdown and conversation views, Lucide icons and a standalone gallery.

The first planned release is **0.2026.1**. The library is not yet published as a downloadable module, and its API is not stable. See the [versioning rules](docs/build.md#publication).

`main` uses upstream MyGo without custom patches. The WebView composition experiment is preserved on `experiment/webview` because it requires a customized MyGo dependency; it is not part of the mainline components.

## Name

**Onevra** is a coined name for this project. It retains the opening `One` from **Oneo** and combines it with `vra` for a compact, distinctive sound and spelling. The ending is a phonetic choice, not an acronym or a word root with an assigned meaning.

The naming approach favors a recognizable connection to Oneo and fewer naming collisions over a common English word. It deliberately avoids tying the library's identity to Go, MyGo or a particular visual style. An initial GitHub repository-name search returned no matches for `onevra`; this is not a guarantee of worldwide availability or trademark clearance.

## Try the gallery

![Onevra Gallery main interface with navigation and interactive component examples](docs/images/gallery.png)

Requires Go 1.27.1 or later and a native desktop session. Linux also requires GTK 3, Pango and OpenGL.

From the repository root:

```sh
cd source
go run ./examples/gallery
```

The gallery includes native shells, the Lucide icon library, Onevra controls and plugins, and MyGo examples. Component pages are grouped by purpose and offer interactive state variants; a trailing dot marks built-in MyGo entries. The Plugins / Terminal page starts the default login shell on first entry; see [runtime requirements](docs/build.md#app-bundles).

To preview a specific appearance, append `-theme light` or `-theme dark`. `go run ./tools/packer run` compiles and starts a packaged development app without hot reload; `go run ./tools/packer build` produces release archives for the targets selected in `source/packer.yaml`. See [packaging and icon requirements](docs/build.md#app-bundles).

The **Desktop / Tray** page enables a tray icon: left-click toggles a small notes window and right-click opens a native menu on macOS and Windows. Linux uses a menu action instead and requires AppIndicator. See [tray review and requirements](docs/build.md#tray-example). On macOS, **Desktop / Dock menu** adds Gallery shortcuts to the Dock icon's right-click menu; see [Dock menu review](docs/build.md#dock-menu-example).

Press Command+E on macOS or Ctrl+E on other platforms to search Gallery pages and actions. Applications can reuse the [command palette and dialog](source/plugins/command/example_test.go), supplying their own commands and execution handlers.

Native review procedures are available for macOS; Windows and Linux still need their own desktop verification. Compilation and headless tests do not establish keyboard, input-method or accessibility acceptance; see [native review](docs/build.md#native-review).

## Basic usage

The module path is `github.com/oneo-me/onevra/source`. **This installation command becomes available only after the first release**:

```sh
go get github.com/oneo-me/onevra/source@v0.2026.1
```

Until publication, use a local Go workspace that includes your application and this repository's `source/` directory; do not commit the workspace configuration. Import `github.com/oneo-me/onevra/source/shells/sidebar`:

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    "github.com/oneo-me/onevra/source/components"
    "github.com/oneo-me/onevra/source/shells/sidebar"
)

func main() {
    mygo.App.WhenReady(func() {
        appShell := sidebar.NewShell(sidebar.ShellOptions{
            Title: "Hello",
            Left: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() {
                    components.Text(c, "Navigation")
                })
            },
            Content: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() {
                    components.Text(c, "Hello, Shell")
                })
            },
        })
        appShell.Window().OnClosed(func() { mygo.App.Quit() })
    })
    if err := mygo.App.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Create a Shell when the app is ready; build its columns with Onevra controls and MyGo layout primitives. Use `appShell.SetPage(view)` to replace the entire sidebar-and-content page (for example, with settings), and `appShell.SetPage(nil)` to restore it. Your application owns business state and lifecycle handling. See the [runnable gallery](source/examples/gallery/main.go) for another example.

Content extends behind Shell's floating header. Use `sidebar.HeaderHeight(c)` for a top inset when needed; paint the padded container's background to cover the titlebar area without placing controls under it. This also applies to replacement pages.

## Appearance

Shell uses Onevra's built-in light and dark themes, applied through MyGo's own theme mechanism so every MyGo widget keeps working. `ShellOptions.Appearance` selects `theme.Auto` (the default, follows the system appearance), `theme.Light` or `theme.Dark`:

```go
import "github.com/oneo-me/onevra/source/theme"

appShell := sidebar.NewShell(sidebar.ShellOptions{
    Title:      "Hello",
    Appearance: theme.Dark,
    // ...
})
```

The Onevra design style pairs warm white light surfaces and warm charcoal dark surfaces with a shared terracotta accent. See the [theme package](source/theme/doc.go) and [appearance design](docs/architecture/components.md#native-appearance).

## Controls

Import `github.com/oneo-me/onevra/source/components` for Onevra-styled common controls built on MyGo's unstyled bases: buttons, check boxes, switches, radio buttons, text fields and areas, selects, sliders and tabs. The package is named `components`, so it can be imported alongside MyGo's `ui` without aliases. Colors come from the active theme:

```go
if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
    // ...
}

components.Checkbox(c, &agreed, "I agree")
components.Switch(c, &wifi, "Wi-Fi")
```

Keyboard, pointer and accessibility behavior come from MyGo's bases. See the [package contract](source/components/doc.go).

## Lucide icons

Import `github.com/oneo-me/onevra/source/components` and `github.com/oneo-me/onevra/source/icons/lucide` to use named accessors with Onevra's Icon component:

```go
components.Icon(c, lucide.Search).Size(20, 20)
components.Icon(c, lucide.Settings).Size(20, 20)
```

Icons are parsed on first use and reused. Icon also accepts parsed `*ui.SVG` assets from other packs. MyGo controls sizing, text-color inheritance and accessibility. The package includes the pinned upstream's regular icons, not lab icons or aliases. See the [package contract](source/icons/lucide/doc.go) and [update workflow](docs/build.md#updating-lucide-icons).

## Contributing

For small application templates, see the multilingual [desktop, menu bar and hybrid examples](docs/build.md#starter-applications), each with its own startup behavior, preferences and packaging configuration.

Read the [architecture overview](ARCHITECTURE.md) before changing code. Run the basic checks from `source/`:

```sh
GOWORK=off go test ./...
GOWORK=off go vet ./...
```

See [development](docs/build.md) for packaging and native review, and [release notes](CHANGELOG.md) for user-facing updates.

## License

Project code: [MIT](LICENSE). Components adapted from MyGo retain the [upstream MIT notice](source/components/LICENSE.mygo). Lucide assets retain their [upstream ISC and Feather MIT notices](source/icons/lucide/LICENSE).
