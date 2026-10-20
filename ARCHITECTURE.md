# Onevra architecture

Onevra is an experimental, product-neutral Go component library built on MyGo. It provides reusable native UI compositions and Lucide SVG assets without introducing an application framework. Its public API is not stable.

## System boundaries

- MyGo owns rendering, layout, input, native controls and desktop integration. The command component has a narrow AppKit adapter for optional input-source selection, the gesture package one for scroll devices and pinches, and the switch paints the Glass plugin's Clear effect with a thinner bezel than `glass.Paint` gives a knob-sized lens, which upstream does not expose; see [complex component composition](docs/architecture/components.md#complex-component-composition).
- Onevra composes MyGo primitives, rather than wrapping every upstream API or maintaining a parallel UI system.
- Onevra restyles, merges and simplifies MyGo's built-in components and adds components MyGo lacks. When MyGo officially provides an equivalent of an Onevra feature or component, Onevra may remove its own version in favor of MyGo's.
- Applications own business state, localized text, persistence, background work and application lifetime.
- Applications own the theme: `mygo.Theme.SetSource` selects light or dark for the whole app, including native chrome, and the application supplies the frame's `ui.Theme`. Onevra's light and dark palettes are an optional preset in the `theme` package. Components read the frame theme. MyGo and the operating system own native window behavior; Shell composes caption buttons on Windows.

## Components and flow

```text
Application / standalone gallery
  ├─ Onevra components → MyGo → native desktop
  └─ MyGo directly
```

Shell styles live in independent subpackages under `shells/`; `sidebar` provides the sidebar-and-content layout, while `plain` provides a single content region below a draggable window-control area. Both hide the default window title bar. An application creates its Shell after MyGo is ready and supplies synchronous UI-thread content builders.

Shell configures and creates a MyGo window, installs its content container; Sidebar Shell also supports an application-supplied replacement page. It exposes the native window handle for application-owned events and lifecycle handling. Layout and appearance ownership are described in [Components](docs/architecture/components.md).

The `icons/lucide` package supplies lazily parsed, shared MyGo SVGs independently of UI components. The `components.Icon` component accepts these named accessors or parsed SVGs while MyGo owns rendering, sizing, color and accessibility. Generated named accessors have no full-library runtime registry or eager parsing, so ordinary consumers need not retain every icon. Upstream resources and licenses are pinned and checked in; a Lucide checkout is only needed for [maintainer updates](docs/build.md#updating-lucide-icons), not consumer builds.

The gallery is an independent consumer of the public API, not a dependency of the library. It keeps interactive demonstration values in memory, owns its presentation-state persistence and provides a place to review native behavior. See [Gallery](docs/architecture/gallery.md) for storage and verification. Its `web` entry compiles the same control demonstrations for the website using MyGo's browser backend; the website supplies navigation and the frame, while native window chrome and desktop integrations are unavailable. See [browser preview](docs/architecture/gallery.md#browser-preview).

The desktop, menu bar, launcher, documents and lawn defense starter applications are independent consumers under `source/examples/`. Each keeps its own lifecycle, preferences where it has them, MyGo message catalog and packaging configuration so it can be copied without importing another example. Their small application helpers are intentionally local, not a shared application framework.

They use MyGo directly for single-instance activation, Dock visibility, tray integration, panels and global shortcuts, and compose Onevra controls for their content. Settings is an application action: desktop hosts it in a Dialog, while the menu bar and launcher examples need no main window and use a settings window. See [starter applications](docs/build.md#starter-applications) for running, copying and native review.

## Repository landmarks

There is one Go module under `source/`, with module path `github.com/oneo-me/onevra/source`:

- `source/shells/`: native window presentation styles, `sidebar` and `plain`.
- `source/components/`: common controls and independent subpackages for complex components, usable without Shell; see [complex component composition](docs/architecture/components.md#complex-component-composition).
- `source/icons/lucide/`: public icon assets and upstream provenance.
- `source/tools/`: library maintenance commands such as the Lucide importer; these share the Go module but are not runtime library dependencies. Applications are packaged with MyGo's CLI, pinned as a Go tool.
- `source/examples/gallery/`: standalone example, with presentation, native review tests and MyGo inventory tooling under its own `internal/` directory; its `web` entry serves the website's browser preview.
- [Development guide](docs/build.md): build, verification and release procedures.
- `website/`: the official website, built with React Router and shadcn/ui and deployed to GitHub Pages; its `/components` page runs the Gallery's demonstrations as WebAssembly; see [website/README.md](website/README.md).

The repository-root [`skills/onevra/`](skills/onevra/SKILL.md) is an Agent skill for consumer applications, covering usage, component choice, hidden tests and project layout. It is not part of the Go module.

No sibling application or shared workspace is required. The module path includes `source` so the repository layout also works for remote consumers. Version tags include the same directory prefix; see [publication](docs/build.md#publication).

## Subjects

| Document | When to read |
| --- | --- |
| [Components](docs/architecture/components.md) | Changing composition boundaries, lifecycle ownership or platform styling |
| [Gallery](docs/architecture/gallery.md) | Changing the standalone example's role or native verification approach |

## Platform confidence

Headless tests and opt-in native captures serve different purposes: the former check component behavior, while the latter exercise desktop startup and rendering. Neither alone establishes complete keyboard or accessibility support. Windows and Linux require their own native review; do not infer support from compilation alone.
