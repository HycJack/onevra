# Changelog

User-facing updates for Onevra library consumers.

## [Unreleased]

- Add a Launcher starter application under `source/examples/launcher`: a Spotlight-style panel that searches the installed apps with `mygo.Shell.InstalledApps`, draws their icons and opens the chosen one, summoned by a configurable global shortcut and living in the tray without a Dock icon. Its Onevra content composes `components.TextBox`, `components.List`, `components.Image` and `components.Settings`, and its messages are a MyGo catalog in English and Simplified Chinese.
- Breaking: `components.Tabs` takes a `TabsVariant` before its labels. Pass `components.TabsUnderline` for the previous accent-underline look, or `components.TabsSegmented` for a macOS-style segmented control: a raised, shadowed thumb over a capsule track with a heavier chosen label, replacing `ui.Segmented`.
- Add `components.ButtonToolbar`, a raised, shadowed control that matches `TabsSegmented`'s chosen segment, for actions beside a segmented tab list.
- Add `components.IconButton`, an icon-only button with an accessible label and tooltip, in a plain or raised toolbar style; `sidebar.IconButton` now shares it. The starter applications use raised icon actions in their toolbars and settings header and a plain `×` to close Settings.
- Add `plain.ShellOptions.HeaderLeading` and `HeaderActions`, leading and trailing header controls. The desktop and hybrid examples put their segmented Note and About tabs in `HeaderLeading` and their Settings action in `HeaderActions`.
- Show the window title in `shells/plain`'s header, so a window without the default title bar still displays it. The header follows `Window().SetTitle`, including later language changes.
- Add the official website at <https://oneo-me.github.io/onevra/>, covering the component set, theming, installation and common questions in English and Simplified Chinese.
- Add `components.Honeycomb`, a picker that arranges round images or generated content in a magnified honeycomb on the surface of a ball: drag or scroll to pan, or Tab through the items, and the item nearest the middle is chosen and ringed in its own color or the theme's accent. Applications own the items, decoded images, labels, selected ID and the action a choice triggers.
- Breaking: require MyGo 0.5.1. Bound controls now write the user's edit through their pointer before the constructor returns, so code after a control reads the new value in the same frame, and an edit applies when the app makes the control read-only or disabled in that frame. MyGo also adds a rich text editor (`ui/richtext`), streaming Markdown (`ui/markdown`), timelines, document-based apps, launchers and app localization for applications to use directly.
- Components and Shells now name their own controls in the window's language, from a MyGo catalog shipped in English and Simplified Chinese: set the language with `mygo.App.SetLocale`, and change or add translations with `onevra.<key>` messages in the application's MyGo catalogs. `components.WithLocalization` and Shell labels still override them.
- Calendars take their month and day names and first weekday from the user's regional formats, and `CalendarPicker` shows the regional numeric date instead of ISO `2006-01-02`; `Localization.FormatDate` and `FirstWeekday` still override them.

## [0.2026.9] - 2026-10-14

- Breaking: Shells no longer choose the theme. Replace `ShellOptions.Appearance` and `plain.Shell.SetAppearance` with `mygo.Theme.SetSource`, which also switches native window chrome, and pass `Theme: theme.Onevra` for Onevra's palettes or a function returning your own `*ui.Theme`; without it, Shells keep MyGo's theme. `theme.Mode` is removed and `theme.Apply(c)` takes no mode.
- Onevra's palettes now follow the desktop's text scale and higher-contrast settings; `theme.Palette(dark)` returns either palette directly.
- Add `components/graph`, an interactive force-directed relationship graph: drag or two-finger touchpad scroll to pan, drag nodes to move them, use the wheel or a pinch to zoom around the pointer, hover to highlight a node's neighbourhood and double-click to open nodes. Applications supply node groups, sizes and colors, edge styles, palettes, group colors and forces, with empty and overlay content slots. Layout runs across processors for graphs of thousands of nodes.
- Add `components/gesture` for canvases that pan and zoom: `Classify` tells mouse wheels from touchpads and maps scrolling to pan or zoom, and `Navigation.Handle` also delivers touchpad pinches, which MyGo does not report, on macOS.
- Remove Packer and `packer.yaml`; package applications with MyGo's CLI (`go tool mygo dev` and `go tool mygo build`) and `mygo.json` instead. MyGo bundles plugin native libraries among the app's resources, and application-owned libraries go in platform directories of `resources/` rather than `runtimes/`.

## [0.2026.8] - 2026-10-13

- Add Onevra-styled toasts on MyGo's toast queue: Shell shows `c.Toast`, `ToastAction` and `AddToast` toasts with status icons for `components.ToastInfo`, `ToastSuccess`, `ToastWarning` and `ToastError`, an optional action, a close button and a `components.ToastContent` leading slot. Standalone hosts call `components.Toasts`; localize the close button with `Labels.CloseToast` or `Localization.CloseToast`.
- Place toasts at the top or bottom edge, centered or at either side, with `components.ToastOptions`, through `ShellOptions.Toasts` and `SetToasts` or the standalone `Toasts` call. Top placements stack downward, newest first, just below the Shell header; `Inset` adds further edge clearance.
- Require MyGo 0.4.8, which adds a code editor and diff viewer (`ui/editor`), offscreen canvases, painter transforms, pixelated bitmaps and the game plugin for applications to use directly.
- Update the Lucide icons with 13 additions, including `lucide.MessageCircleText`, `ServerNas`, `ShieldHouse`, `UserRoundStar` and `WifiLock`, and refreshed artwork for `PenTool`, `Salad`, `Soup` and `ToolCase`.

## [0.2026.7] - 2026-10-12

- Add `Leading` and `Trailing` content slots to command items, for images or other custom content beside the title.
- Add `TopInset` to floating conversations, so transcripts scroll beneath a floating window header while their scrollbar, Header and first turn start below it.
- Add `StateKey` to sidebar and plain Shell options, so MyGo remembers and restores a window's bounds and maximized or full-screen state, in its own file or a store set with `mygo.App.SetWindowStateStore`.
- Require MyGo 0.4.6, which paints text selections above text and span backgrounds in text boxes, code and other selectable text.

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
