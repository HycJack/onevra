# Changelog

User-facing updates for oneo-go library consumers. No version has been published yet.

## [Unreleased]

- Embed experimental WebView components in native MyGo windows on macOS ARM64, including multiple independent browser slots, MyGo Page navigation/JavaScript APIs, document-preserving hide/show and automatic cleanup with the host. Try them in Gallery's Native + WebView page. System views remain above MyGo-rendered overlays; scroll clipping, arbitrary transforms and unified composition are not supported.

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Create a native Shell with system window controls and a full-height left sidebar with a Lucide toggle. The sidebar reveals beneath the main content without resizing the window and keeps its content width during animations. An overlay separator leaves no gap when collapsed; drag it or use arrow keys when focused to resize the sidebar. Animations respect Reduced Motion. On macOS, the sidebar exposes the window's blurred background up to the top edge, while main content remains opaque. The toggle keeps clearance from traffic lights.
- Add icon-only toolbar actions with accessible labels, tooltips and translucent text-color hover/press feedback. Place application actions in Shell's fixed left bottom toolbar, and supply a title, breadcrumbs or actions in the leading main-content titlebar area.
- Replace the entire sidebar-and-content page with application content such as settings, then restore the previous sidebar state.
- Try Shell in the standalone gallery, including light/dark appearances, sidebar controls and a replacement settings page.
