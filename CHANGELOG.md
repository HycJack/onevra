# Changelog

User-facing updates for Onevra library consumers. No version has been published yet.

## [Unreleased]

- Support MyGo 0.3 value-based Element handles for common controls, selection popups and Shell icon buttons; persistent identity uses MyGo Handle bindings.

- Add themed controls for text composition and editing with string or indexed-buffer storage, numeric input with prefix scrubbing and optional step controls, selection with date and time popups, RGB/HSB/HSL color editing with native format menus and copying with customizable popup swatches, ratings, status with ring and spoke loading indicators, navigation, virtualized list and grid collections, and layout, usable alongside MyGo widgets independently of Shell.
- Add theme-tinted glass surfaces for text Tooltips, anchored Popovers and selection popups, plus composable layered glass Dialogs and result-bearing AlertDialogs with optional cancellation; dialogs inherit window dragging from Shell or a standalone host scope, with per-dialog overrides.
- Add a native sidebar-style application Shell through `shells/sidebar` with a collapsible, resizable sidebar, customizable titlebar and toolbar actions, and replaceable page content; expose system materials on macOS and supported Windows systems with opaque theme backgrounds when unavailable.
- Add selectable Markdown documents with code copying, application-supplied images and custom code-block views.
- Add reusable settings surfaces with application-defined controls, grouped categories, narrow layouts and cross-category search, embeddable in pages or dialogs.
- Add built-in Onevra light and dark themes with warm neutrals and a shared terracotta accent that support system appearance or an application-selected mode and work with MyGo widgets.
- Add Lucide assets through `icons/lucide` and a composable Icon component accepting named icon accessors or parsed SVGs.
