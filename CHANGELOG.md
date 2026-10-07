# Changelog

User-facing updates for oneo-go library consumers. No version has been published yet.

## [Unreleased]

- Embed experimental WebView components in native MyGo windows on macOS ARM64, including multiple independent browser slots, MyGo Page navigation/JavaScript APIs, document-preserving hide/show and automatic cleanup with the host. Try them in Gallery's Native + WebView page. System views remain above MyGo-rendered overlays; scroll clipping, arbitrary transforms and unified composition are not supported.

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Create a native Shell with system window controls, full-height sidebars with Lucide toggle controls: an interior sidebar revealed beneath the main content and an optional inspector revealed as the window expands to the right. Sidebar content keeps its width during these animations; the left and main layout remain stationary during unconstrained inspector expansion, including in narrow windows. Overlay separators leave no gap when collapsed. Resize the left sidebar by dragging its separator or using arrow keys when focused. Animations respect Reduced Motion. On macOS, the left sidebar exposes the window's blurred background up to the top edge, the inspector is lightly translucent, and main content remains opaque. The left toggle keeps clearance from traffic lights, and the inspector toggle stays at the main content's right edge.
- Add icon-only toolbar actions with accessible labels, tooltips and translucent text-color hover/press feedback. Place application actions in Shell's fixed left bottom toolbar, and supply a title, breadcrumbs or actions in the leading main-content titlebar area.
- Replace the entire three-column page with application content such as settings, then restore the previous sidebar state.
- Try Shell in the standalone gallery, including light/dark appearances, sidebar controls and a replacement settings page.
