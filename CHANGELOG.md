# Changelog

User-facing updates for Onevra library consumers.

## [Unreleased]

- Require MyGo 0.4.1, including corrected Linux outer-window and content bounds.
- Paint the composer's Send button in themes that leave `Inverse` unset, including Onevra's palettes; it previously reserved space beside Stop without being drawn.
- Stop the composer from crashing on paste when no `OnPaste` hook is supplied.

## [0.2026.2] - 2026-10-11

- Let settings groups place application actions, such as a refresh button, at the end of their header row.
- Require MyGo 0.4.0. Popover elements report Escape or outside-press dismissal through `Dismissed`, and Select popups now keep their wheel scroll position while revealing options reached with the keyboard.

## [0.2026.1] - 2026-10-10

Initial experimental release.

- Add themed controls for editing, selection, navigation, virtualized collections and layout, with application-supplied localization.
- Add composable popovers, dialogs and settings surfaces with application-owned state and actions.
- Add sidebar and single-content native window shells with light, dark and system appearance modes.
- Add searchable command palettes with keyboard navigation and optional macOS input-source selection.
- Add read-only syntax highlighting and selectable Markdown with application-supplied images and custom code-block views.
- Add composable conversation views with activity groups, response folding and a message composer; applications own networking and asynchronous work.
- Add Lucide icons usable independently or through the common Icon component.
- Add standalone multilingual desktop, menu bar and hybrid application templates with configurable lifecycle behavior and persistent preferences.
