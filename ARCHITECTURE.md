# oneo-go architecture

oneo-go is an experimental, product-neutral Go component library based on MyGo native UI. The Go module in `source/` is `github.com/oneo-me/oneo-go`; its Go package is `oneogo`. This repository owns component implementation, documentation, examples and tests. No consumer application or sibling repository is required.

Read this overview before changes, then choose the relevant Subjects. [README.md](README.md) introduces usage and [development](docs/build.md) owns verification commands.

## Boundaries

1. MyGo owns rendering, layout, input, focus, native controls and desktop integration.
2. oneo-go composes those primitives; it does not wrap every MyGo API or impose an application framework.
3. Callers own business state, localized text, actions, persistence and background work. Components build synchronously on the UI thread.
4. Components use the current MyGo theme without mutating it. Product-specific page widths and window chrome remain outside the library.
5. Public API is experimental. The module is local and unpublished; its path does not imply a downloadable remote repository or released version.

## Modules and flow

```text
Consumer / standalone gallery
  ├─ oneogo components → MyGo ui
  └─ MyGo directly
```

The module-root package in `source/` exports EmptyState, Section and Notice. Each returns a native `*ui.Element`, enabling layout modifiers and keys. Caller callbacks build optional content and native controls; their events update caller state.

The gallery executable uses `source/internal/gallery` for its view and demonstration state. It does not write content or application settings. Its native smoke test isolates user data in a temporary directory. There is no application dependency, filesystem service, persistence layer or network client in public components.

Dependencies and the Go toolchain requirement are pinned in `source/go.mod` and `source/go.sum`. There is one Go module and no required `go.work`, sibling-module replacement or nested repository.

## Repository layout

```text
AGENTS.md                 Contributor constraints
ARCHITECTURE.md           Overview and subject index
README.md / README.zh-CN.md Usage and quick start
CHANGELOG.md              Pending API/gallery changes
LICENSE                   MIT license
source/
  go.mod / go.sum         Module and pinned dependencies
  *.go                    Public components and headless tests
  mygo.json               Gallery app identity and MyGo build configuration
  examples/gallery/       Standalone native executable
  internal/gallery/       Demonstration view and its tests
  .mygo/                  Ignored development app bundles
docs/architecture/        Owning subject documents
docs/build.md             Checks, platforms and release boundaries
dist/                     Ignored local binaries and captures
```

## Subjects

| Document | Read before changing |
| --- | --- |
| [Components](docs/architecture/components.md) | Public component contracts, ownership, theme usage and composition |
| [Gallery](docs/architecture/gallery.md) | Example state, executable, native capture and acceptance boundaries |

## Verification and releases

Headless tests exercise component composition, caller-owned actions, disabled controls, scaling, narrow layout and theme preservation. Opt-in macOS native captures exercise light, dark and narrow gallery windows. These do not establish complete keyboard, accessibility, IME or Windows/Linux acceptance.

[CHANGELOG.md](CHANGELOG.md) owns consumer-visible changes. There is no release automation or published version. Publishing a remote repository, assigning versions, creating tags and pushing require separate authorization.
