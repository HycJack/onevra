# Changelog

User-facing updates for Onevra library consumers.

## [Unreleased]

## [0.2026.10] - 2026-10-15

- Breaking: require MyGo 0.5.1, whose bound controls write an edit through their pointer before the constructor returns.
- Breaking: remove `components/markdown` and `components/code`; use MyGo's `ui/markdown` and a read-only `ui/editor` instead.
- Breaking: `components.Tabs` takes a `TabsVariant` first; pass `TabsUnderline` for the previous look or `TabsSegmented` for a segmented control.
- Add `components.IconButton`, an icon-only button with an accessible label and tooltip, which `sidebar.IconButton` now uses.
- Add `components.ButtonToolbar` for actions beside segmented tabs.
- Add `components.Honeycomb`, a picker that arranges round items in a magnified honeycomb.
- Add `plain.ShellOptions.HeaderLeading` and `HeaderActions` for controls in the plain Shell's header, which now shows the window title.
- Label component and Shell controls in the app's language, with English and Simplified Chinese built in and `onevra.<key>` catalog messages as overrides.
- Take calendar names, the first weekday and `CalendarPicker` dates from the user's regional formats.
- Add the Lawn defense game and Launcher examples under `source/examples/`.
- Add the official website at <https://oneo-me.github.io/onevra/>.

## [0.2026.9] - 2026-10-14

- Breaking: Shells no longer set the theme; replace `ShellOptions.Appearance`, `SetAppearance` and `theme.Mode` with `mygo.Theme.SetSource` and `Theme: theme.Onevra`.
- Breaking: remove Packer; package applications with `go tool mygo build` and `mygo.json`.
- Follow the desktop's text scale and higher-contrast settings in Onevra's palettes.
- Add `components/graph`, an interactive force-directed relationship graph.
- Add `components/gesture` to tell mouse wheels from touchpads and handle pan, zoom and macOS pinch gestures.

## [0.2026.8] - 2026-10-13

- Add Onevra-styled toasts with status icons, an optional action and configurable placement, shown by Shells or `components.Toasts`.
- Require MyGo 0.4.8.
- Add 13 Lucide icons and refresh four others.

## [0.2026.7] - 2026-10-12

- Add `Leading` and `Trailing` content slots to command items.
- Add `TopInset` to floating conversations so transcripts can scroll beneath a floating header.
- Add `StateKey` to sidebar and plain Shells to restore a window's bounds and maximized or full-screen state.
- Require MyGo 0.4.6, which draws text selections above text backgrounds.

## [0.2026.6] - 2026-10-12

- Breaking: import command, code, markdown and chat from `components/<name>` instead of `plugins/<name>`.
- Add bar, line, area, stacked bar and donut charts.
- Require MyGo 0.4.3, which keeps interactive controls aligned under transforms.
- Add typed macOS Info.plist entries to Packer configuration.
- Fix the floating conversation's scroll-to-latest button colors in dark mode.

## [0.2026.5] - 2026-10-11

- Add tables, ordinary and virtualized, with row selection, sorting and resizable, reorderable columns.
- Add a composer trailing slot for controls before Send or Stop.
- Add leading content slots to sidebar items and settings categories.
- Let conversation headers overlay messages, and keep scrollbars clear of the floating composer.
- Keep Send available during input-method composition, and show it instead of Stop when a running request has a draft.
- Let Packer bundle native libraries from `runtimes/` and keep the uncompressed app beside release archives.
- Require MyGo 0.4.2, which lets spoke spinners use custom text colors.

## [0.2026.4] - 2026-10-11

- Add `ShellOptions.Hidden` to start sidebar Shells hidden.
- Add `SettingsItem.DescriptionContent` for settings descriptions with inline links.

## [0.2026.3] - 2026-10-11

- Require MyGo 0.4.1, which fixes Linux window bounds.
- Fix the composer's Send button not drawing in themes without `Inverse`, including Onevra's.
- Fix a composer crash on paste without an `OnPaste` hook.

## [0.2026.2] - 2026-10-11

- Add an action slot to settings group headers.
- Require MyGo 0.4.0, which reports popover dismissal through `Dismissed`.

## [0.2026.1] - 2026-10-10

Initial experimental release.

- Add themed controls for editing, selection, navigation, collections and layout.
- Add popovers, dialogs and settings surfaces.
- Add sidebar and single-content window Shells with light, dark and system appearance.
- Add searchable command palettes.
- Add syntax highlighting and selectable Markdown.
- Add conversation views with activity groups, response folding and a message composer.
- Add Lucide icons.
- Add multilingual desktop, menu bar and hybrid application templates.
