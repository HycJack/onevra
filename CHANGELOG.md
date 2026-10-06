# Changelog

Pending API and gallery changes for oneo-go consumers. This module is not yet published;

## [Unreleased]

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Run the gallery as a live-reloading macOS development app from `source/` with `go tool mygo dev`, or build local app bundles with `go tool mygo build`; matching VS Code tasks are included. Go source, tests, module files and MyGo configuration live in `source/`; build output remains in the repository-root `dist/`.
- Compose empty states, sections and persistent notices using native MyGo widgets, caller-owned state and the current theme.
- Preview components in a standalone gallery with light/dark appearances and interactive create/reset and disabled-action examples.
