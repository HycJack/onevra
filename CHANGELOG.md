# Changelog

User-facing updates for Onevra library consumers.

## [Unreleased]

- Add `Leading` and `Trailing` content slots to command items, for images or other custom content beside the title.
- Require MyGo 0.4.5, which paints text selections above text and span backgrounds in text boxes, code and other selectable text.

## [0.2026.6] - 2026-10-12

- Breaking: import command, code, markdown and chat from `components/<name>` instead of `plugins/<name>`.
- Add themed bar, line, area, stacked bar and donut charts with smooth updates, interactive legends, keyboard selection and custom empty/detail content.
- Require MyGo 0.4.3 and keep table resizing, range sliders, color picking, split limits and tab indicators aligned under transforms.
- Add typed macOS Info.plist extensions through Packer configuration, including usage descriptions in development and release bundles before signing.
- Fix the floating conversation's scroll-to-latest button to follow the theme background and text colors in dark mode.

## [0.2026.5] - 2026-10-11

- Require MyGo 0.4.2 and let spoke-style spinners honor custom text colors while retaining their muted default.
- Add ordinary and virtualized tables with alternating row stripes, row selection, column sorting, resizing and reordering. Applications own row data, sort order and persistable column layout; virtualized tables support horizontal scrolling.
- Let Packer include native libraries from target-specific `runtimes/` directories with macOS signing, and retain the uncompressed host application beside release archives by default.
- Let applications place compact controls immediately before Send or Stop through the composer's trailing slot.
- Let conversation headers overlay messages with customizable top spacing, and keep scrollbars clear of the floating composer as its height changes.
- Keep the composer's Send button available during input-method composition, while preserving Enter for candidate selection; a nonempty draft now replaces Stop with Send during a running request.
- Customize the leading content of sidebar items and settings categories with slots, including images, monograms or status marks; the image and SVG convenience APIs remain available.

## [0.2026.4] - 2026-10-11

- Let sidebar shells start hidden with `ShellOptions.Hidden`, so applications can prepare their window before showing or focusing it.
- Let settings rows compose their description from inline text and links, such as a count that opens a popover, with `SettingsItem.DescriptionContent`; search still matches `Description`.

## [0.2026.3] - 2026-10-11

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
