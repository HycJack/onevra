# Components

The module-root `oneogo` package in `source/` contains product-neutral compositions of MyGo widgets. It depends only on MyGo, with no consumer application imports. Native controls remain usable directly alongside these components.

## Public API

- `Section(c, title, description, content)` builds a heading, optional description and synchronous content callback.
- `EmptyState(c, EmptyStateOptions)` builds a leading-aligned explanation and optional caller-supplied actions.
- `Notice(c, NoticeOptions)` builds persistent feedback with information, warning, error or success heading colors and optional actions.

Each returns `*ui.Element` for ordinary MyGo layout modifiers and keys. Optional descriptions and nil callbacks are omitted. Titles should be meaningful caller-localized text. Zero-valued NoticeKind means information; unknown kinds use the default text color.

Callers build actions with native controls and handle their events. Components do not launch work, dismiss themselves or own business state. Repeated components need distinct caller-provided keys; internal child identity follows MyGo's positional mechanism.

## Styling and semantics

Components derive dimensions from the active theme's Space and Rem helpers and colors from its semantic roles. They neither mutate nor replace the theme. Desktop scaling and built-in control preferences follow MyGo; no parallel theme or global density state exists.

Components use labeled groups. Notice is not an automatically announced live alert; callers must provide explicit status text instead of communicating severity only by color. Native controls supply their keyboard, focus and accessibility behavior.

Components impose no application window, article width, file format, language or persistence policy. Root elements can be further styled by consumers.

## Verification

Module-root headless tests verify composition, optional content, caller-owned button activation, disabled controls, theme identity, text scaling and narrow bounds across both themes. The [gallery](gallery.md) supplies native rendering coverage; full assistive-technology and keyboard review remains manual. [README usage](../../README.md#components) shows the public call pattern.
