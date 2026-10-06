# Changelog

User-facing updates for oneo-go library consumers. No version has been published yet.

## [Unreleased]

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Create a native Shell with system window controls, full-height sidebars with Lucide toggle controls: an interior sidebar revealed beneath the main content and an optional inspector revealed as the window expands to the right. Sidebar content keeps its width during these animations, with overlay separators that leave no gap when collapsed. Resize the left sidebar by dragging its separator or using arrow keys when focused. Animations respect Reduced Motion. On macOS, the sidebar exposes the window's blurred background up to the top edge; content and inspector regions remain opaque.
- Replace the entire three-column page with application content such as settings, then restore the previous sidebar state.
- Try Shell in the standalone gallery, including light/dark appearances, sidebar controls and a replacement settings page.
