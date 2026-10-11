---
name: onevra
description: Build, test and package native Go desktop applications with Onevra components on MyGo. Use when writing or reviewing UI that imports github.com/oneo-me/onevra/source, choosing between Onevra and MyGo controls, writing headless or hidden-window UI tests, or setting up packer.yaml and the source/ project layout. Excludes changes to MyGo itself.
---

# Onevra

Onevra (`github.com/oneo-me/onevra/source`) is a product-neutral component library built on MyGo. MyGo owns rendering, layout, input, windows and plugins; Onevra supplies themed controls, window shells, Lucide icons, complex plugins and the Packer packaging tool. The application owns business state, localized text, persistence and asynchronous work.

Before using an API, read its public signature in the version-matched module source rather than guessing:

```sh
go list -m -f '{{.Dir}}' github.com/oneo-me/onevra/source   # Onevra source and examples
go list -m -f '{{.Dir}}' github.com/egoist/mygo              # MyGo docs/ and ui/
```

Run these inside the application module (`source/`); add `GOWORK=off` there when checking that the module builds independently. The Onevra module directory also contains runnable starter applications under `examples/desktop`, `examples/menubar` and `examples/hybrid`; copy their structure rather than inventing a lifecycle.

## Basic usage

```sh
go get github.com/oneo-me/onevra/source@latest
```

```go
func main() {
	mygo.App.WhenReady(func() {
		shell := plain.NewShell(plain.ShellOptions{
			Title: "Hello",
			Content: func(c *ui.Context) {
				ui.Column(c).Fill().Padding(24).Gap(12).Children(func() {
					components.Text(c, "Hello, Onevra!")
					if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
						// Act on application state.
					}
					components.Icon(c, lucide.Search).Size(20, 20)
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

- Create windows only after `mygo.App.WhenReady`. Choose `shells/plain` for a single content area or `shells/sidebar` for sidebar navigation; use `shell.Window()` for native events instead of adding a second lifecycle. Pass `Hidden: true` when lifecycle handlers or content must be prepared before showing.
- Builders run synchronously on the UI thread every frame. Keep state in retained application values passed by pointer (`components.Checkbox(c, &agreed, "I agree")`); never start background work, I/O or parsing inside a builder. Do background work elsewhere and apply results on the main thread, then request a frame.
- Element handles are valid only for the current build pass. Bind a `ui.Handle` for identity or geometry across frames.
- Shells apply Onevra's theme. Without a Shell, call `theme.Apply(c, mode)` at the top of the root builder. Components never call `Context.SetTheme`.
- Pass translated labels explicitly and wrap content in `components.WithLocalization` for component-owned labels. Keep keys and stored values independent of translated text.
- Complex views are in `plugins/chart`, `plugins/command`, `plugins/code`, `plugins/markdown` and `plugins/chat`. Keep their documents and results in retained state, prepared outside builders.

## Choosing controls

Use Onevra's control whenever one exists; it replaces MyGo's styled widget of the same purpose. Check `components/`, `shells/` and `plugins/` in the Onevra module before choosing a control.

| Need | Use |
| --- | --- |
| Button, link, primary button, menu button | `components.Button` with a variant; `.Menu(...)` for menus |
| Text, editable text, rich inline text | `components.Text`, `components.RichText` |
| Text input, text area, search field | `components.TextBox`, `components.TextBoxBuffer` |
| Checkbox, switch, radio, toggle | `components.Checkbox`, `Switch`, `Radio`, `ToggleButton` |
| Select, combobox, autocomplete | `components.Select` (searchable), `components.NativeSelect` |
| Number, slider, date, time, color | `NumberInput`, `Slider`, `Calendar`/`CalendarPicker`, `TimePicker`, `ColorPicker`/`ColorView` |
| Scroll, list, grid, tree, outline, table | `ScrollView`, `List`, `GridView`, `Tree`, `TreeVirtualized`, `Table`, `TableVirtualized` |
| Tabs, accordion, collapsible, split, breadcrumbs, sidebar | The Onevra component of the same name |
| Modal, alert, popover, tooltip | `Dialog`, `AlertDialog`, `Popover`, `Tooltip` |
| Progress, spinner, badge, rating, icon, image | The Onevra component of the same name; icons from `icons/lucide` |
| Categorical charts and donuts | `plugins/chart`: `Bar`, `Line`, `Area`, `StackedBar`, `Donut`; see the [usage example](../../source/plugins/chart/example_test.go) |
| Settings surface | `components.Settings` |

Use MyGo directly for:

- Layout primitives: `ui.Box`, `Column`, `Row`, `Spacer`, `Grid`, `Overlay`, and `ui.Textf` for formatted display text.
- Specialized charts and continuous/time axes from `ui/chart`, and element transforms and motion methods. Use `LayoutBounds` for layout calculations and local pointer coordinates; `Element.Bounds` includes visual transforms.
- Native menus and context menus, `Flyout`, `WebView`, and every MyGo plugin, including Glass. Do not write Onevra-style replacements.
- MyGo's unstyled bases when building a new reusable control.
- A MyGo widget with no Onevra equivalent. Onevra intentionally omits Toast, Divider, Avatar, Meter, Form, Field, Toolbar, OutlineTable, Segmented and TokenField: prefer contextual feedback over toasts and compose dividers from Box or borders before reaching for the upstream widget.

Charts accept application-owned data and formatting, with generic empty and detail slots. Keep series/slice IDs stable and state separate per chart; `chart.State{Index: -1}` starts unselected. Prepare large datasets outside builders; these views do not downsample or virtualize data.

The authoritative replacement and exclusion map, with reasons, is `examples/gallery/internal/gallery/mygo_review.go` in the Onevra module.

## Tests

Write tests that never show a window by default. Read [references/testing.md](references/testing.md) before adding UI or native tests.

## Project layout and packaging

Package applications with Onevra Packer and its `packer.yaml`, not MyGo's packaging commands or `mygo.json`. Read [references/packer.md](references/packer.md) before creating a project, changing its layout or packaging configuration.

## Checks

From the Go module directory:

```sh
gofmt -l .
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build ./...
```

`gofmt -l .` must print nothing. Headless tests and successful cross-builds do not prove native behavior: review UI changes in a desktop session on each claimed platform, covering both themes, narrow windows, keyboard navigation, input methods and accessibility, and report what remains unreviewed.
