# Changelog

User-facing updates for oneo-go library consumers. No version has been published yet.

## [Unreleased]

- Embed experimental live WebViews in native MyGo windows, retaining MyGo Page APIs, document-preserving hide/show and host-owned cleanup. GPU composition supports native overlays, modal input blocking, rectangular ancestor clipping and automatic hiding when a slot is omitted; these behaviors have macOS ARM64 native checks. Gallery's existing Native + WebView page uses this implementation directly; see [build requirements and platform limitations](docs/build.md#hybrid-experiment).

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Create a native Shell with system window controls and a full-height left sidebar with a Lucide toggle. The sidebar reveals beneath the main content without resizing the window and keeps its content width during animations. An overlay separator leaves no gap when collapsed; drag it or use arrow keys when focused to resize the sidebar within configurable minimum and maximum widths. Set a default width and restore it by double-clicking the separator. Animations respect Reduced Motion. On macOS, the sidebar exposes the window's blurred background up to the top edge, while main content remains opaque. The toggle keeps clearance from traffic lights.
- Add icon-only toolbar actions with accessible labels, tooltips and translucent text-color hover/press feedback. Place application actions in Shell's fixed left bottom toolbar, and supply a title, breadcrumbs or actions in the leading main-content titlebar area.
- Replace the entire sidebar-and-content page with application content such as settings, then restore the previous sidebar state.
- Try Shell in the standalone gallery, including light/dark appearances, sidebar controls and a replacement settings page.
