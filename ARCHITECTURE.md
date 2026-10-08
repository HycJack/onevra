# Onevra architecture

Onevra is an experimental, product-neutral Go component library built on MyGo. It provides reusable native UI compositions and Lucide SVG assets without introducing an application framework. The library is unpublished and its public API is not stable.

## System boundaries

- MyGo owns rendering, layout, input, native controls and desktop integration.
- Onevra composes MyGo primitives, rather than wrapping every upstream API or maintaining a parallel UI system.
- Applications own business state, localized text, persistence, background work and application lifetime.
- Components use Onevra's built-in light and dark themes, applied per frame through MyGo's theme mechanism without mutating global appearance. MyGo and the operating system own native window behavior; Shell composes caption buttons on Windows.

## Components and flow

```text
Application / standalone gallery
  ├─ Onevra components → MyGo → native desktop
  └─ MyGo directly
```

An application creates a Shell after MyGo is ready and supplies synchronous UI-thread content builders. Shell configures and creates a MyGo window, preprocesses its content through a shared page container, and presents either the default sidebar-and-content page or an application-supplied replacement. It exposes the native window handle for application-owned events and lifecycle handling. Layout and appearance ownership are described in [Components](docs/architecture/components.md).

The `icons/lucide` package supplies lazily parsed, shared MyGo SVGs independently of UI components. The `components.Icon` component accepts these named accessors or parsed SVGs while MyGo owns rendering, sizing, color and accessibility. Generated named accessors have no full-library runtime registry or eager parsing, so ordinary consumers need not retain every icon. Upstream resources and licenses are pinned and checked in; a Lucide checkout is only needed for [maintainer updates](docs/build.md#updating-lucide-icons), not consumer builds.

The gallery is an independent consumer of the public API, not a dependency of the library. It keeps interactive demonstration values in memory, owns its presentation-state persistence and provides a place to review native behavior. See [Gallery](docs/architecture/gallery.md) for storage and verification.

## Repository landmarks

There is one Go module under `source/`, with module path `github.com/oneo-me/onevra`:

- `source/shell/`: native window presentation and the replaceable sidebar-and-content page container.
- `source/components/`: common controls built on MyGo's unstyled bases, usable independently of Shell.
- `source/icons/lucide/`: public icon assets and upstream provenance.
- `source/tools/`: Packer, the minimal desktop app runner/packager, and library maintenance commands such as the Lucide importer; these share the Go module but are not runtime library dependencies.
- `source/examples/gallery/`: standalone example, with presentation, native review tests and MyGo inventory tooling under its own `internal/` directory.
- [Development guide](docs/build.md): build, verification and release procedures.

No sibling application or shared workspace is required. The source subdirectory and current module path support local development; remote publication still requires a repository-layout or module-path decision.

## Subjects

| Document | When to read |
| --- | --- |
| [Components](docs/architecture/components.md) | Changing composition boundaries, lifecycle ownership or platform styling |
| [Gallery](docs/architecture/gallery.md) | Changing the standalone example's role or native verification approach |

## Platform confidence

Headless tests and opt-in native captures serve different purposes: the former check component behavior, while the latter exercise desktop startup and rendering. Neither alone establishes complete keyboard or accessibility support. Windows and Linux require their own native review; do not infer support from compilation alone.
