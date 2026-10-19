# Changelog

User-facing updates for Onevra library consumers.

## [Unreleased]

- Add `components.HoneycombItem.Menu`.
- Add `components.IconButtonCompact`.
- Add `components.SettingsCategory.Placeholder`.
- Add `components.SettingsOptions.SectionActions`.
- Add `components.SidebarSectionWithActions`.
- Add `components.TabsWith`.
- Add `components.TabsWithLeading`.
- Localize component and Shell labels in Traditional Chinese, Japanese and Korean.

## [0.2026.12] - 2026-10-18

- Add `components.Avatar` and `components.AvatarGroup`.
- Add `components.BadgeOptions` with tones and a leading icon.
- Add `components.CardStack`.
- Add `components.Carousel`.
- Add `chart.Bubble`.
- Add `chart.SplitBar`.
- Add `components/flow`.
- Add `components.CollapsibleOptions.Leading`.
- Add `components.HoldToConfirm`.
- Add `components.FloatingButtonGroup`.
- Add `components.HoneycombItem.Overlay` and `components.HoneycombOptions.OnPress`.
- Add `components.ImageCompare`.
- Add `components.SidebarItemWithSlots` and `components.SidebarItemWithContent`.
- Add `components.Skeleton` and `components.SkeletonContent`.
- Add `components.SplitOptions.Default`.
- Add `components.Stepper`.
- Add `components.TagInput`.
- Add `shaders.Ripple`.
- Add `components.TreeItemWithLeading`.
- Add `treemap.View`.
- Let `chat.Shimmer` shrink and end in an ellipsis.
- Animate `components.Badge` text and icon changes.
- Stack `components.Toasts` behind the newest toast and fan them out on hover or focus.
- Fix development-build panics when a closing `components.Dialog` contains `components.Checkbox`, `components.Radio`, `components.Accordion`, `components.Slider` or color controls.

## [0.2026.11] - 2026-10-16

- Breaking: remove the hybrid application template.
- Breaking: remove `components.Localization` and `components.WithLocalization`.
- Breaking: remove Shell label overrides `sidebar.Labels`, `plain.Labels` and their `SetLabels` and label options.
- Breaking: remove `Header`, `HeaderOverlay`, `HeaderTopInset` and `TopInset` from `chat.FloatingConversationOptions`.
- Breaking: rename `command.CommandPalette` to `command.CommandView`.
- Add `command.Options.Filter` and `command.State.Focus`.
- Add `HideSearch`, `HideSidebar` and `Close` to `components.SettingsOptions`.
- Add `Leading`, `Stacked` and `Content` to `components.SettingsItem`.
- Add `components.SplitPercent`.
- Add `components.TabsHighlight`.
- Add `shaders.SignalField`.
- Add `components.ListState.KeepOffset`.
- Add `chat.Chip`, `chat.SwitchChip`, `chat.Picker`, `chat.Attachment`, `chat.ContextMeter`, `chat.Quota`, `chat.QuotaSummary` and `chat.Shimmer`.
- Require MyGo 0.5.5.
- Change the default `sidebar.ShellOptions` window height and minimum height to 740.
- Return command results to the top when they change.
- Animate `components.Dialog` and `components.AlertDialog` as they open and close.
- Animate `components.Honeycomb` items as they are added or removed.
- Fix `components.Honeycomb` centering on first display and after its view resizes.

## [0.2026.10] - 2026-10-15

- Breaking: require MyGo 0.5.1.
- Breaking: replace `components/markdown` and `components/code` with MyGo's `ui/markdown` and `ui/editor`.
- Breaking: add a `TabsVariant` parameter to `components.Tabs`.
- Add `components.IconButton`.
- Add `components.ButtonToolbar`.
- Add `components.Honeycomb`.
- Add `HeaderLeading` and `HeaderActions` to `plain.ShellOptions`.
- Add the Lawn defense and Launcher examples.
- Add the official website.
- Show the window title in the plain Shell's header.
- Localize component and Shell labels in English and Simplified Chinese.
- Follow regional date formats in calendars.

## [0.2026.9] - 2026-10-14

- Breaking: replace Shell appearance options with `mygo.Theme.SetSource`.
- Breaking: replace Packer with MyGo's CLI.
- Add `components/graph`.
- Add `components/gesture`.
- Follow the desktop's text scale and contrast settings.

## [0.2026.8] - 2026-10-13

- Breaking: require MyGo 0.4.8.
- Add toasts.
- Add 13 Lucide icons.

## [0.2026.7] - 2026-10-12

- Breaking: require MyGo 0.4.6.
- Add `Leading` and `Trailing` slots to command items.
- Add `TopInset` to floating conversations.
- Add `StateKey` to remember window bounds.

## [0.2026.6] - 2026-10-12

- Breaking: require MyGo 0.4.3.
- Breaking: move command, code, markdown and chat from `plugins/` to `components/`.
- Add charts.
- Add macOS Info.plist entries to Packer.
- Fix the scroll-to-latest button in dark mode.

## [0.2026.5] - 2026-10-11

- Breaking: require MyGo 0.4.2.
- Add tables.
- Add a trailing slot to the composer.
- Add leading slots to sidebar items and settings categories.
- Bundle native libraries with Packer.
- Let conversation headers overlay messages.
- Keep Send available during input-method composition.

## [0.2026.4] - 2026-10-11

- Add `ShellOptions.Hidden`.
- Add `SettingsItem.DescriptionContent`.

## [0.2026.3] - 2026-10-11

- Breaking: require MyGo 0.4.1.
- Fix the composer's Send button not drawing in some themes.
- Fix a composer crash on paste.

## [0.2026.2] - 2026-10-11

- Breaking: require MyGo 0.4.0.
- Add an action slot to settings group headers.

## [0.2026.1] - 2026-10-10

Initial experimental release.

- Add themed controls, popovers, dialogs and settings.
- Add sidebar and plain window Shells.
- Add command palettes, Markdown and syntax highlighting.
- Add conversation views and a message composer.
- Add Lucide icons.
- Add desktop, menu bar and hybrid application templates.
