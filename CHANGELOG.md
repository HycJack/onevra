# Changelog

User-facing updates for oneo-go library consumers. No version has been published yet.

## [Unreleased]

- Use oneo-styled common controls from `ui`: buttons (primary, default, outline, ghost and link variants in three sizes), check boxes, switches, radio buttons, text fields and areas, drop-down selects, sliders and tab rows. They are built on MyGo's unstyled bases, so pointer, keyboard and accessibility behavior come from MyGo, while colors follow oneo-go's built-in themes. Preview every control in the gallery's **oneo-go controls** group, including button styles and sizes, editable and read-only inputs, and selection states.

- Use 1,866 Lucide icons directly with MyGo through `icons/lucide`, with named accessors, first-use parsing and inherited text colors. Preview representative icons at multiple sizes in the gallery.

- Create a native Shell through `shell.NewShell` with system window controls on macOS/Linux and rounded, localizable toolbar-style minimize, maximize/restore and close buttons on Windows and a full-height left sidebar with a Lucide toggle. The sidebar reveals beneath the main content without resizing the window and keeps its content width during animations. An overlay separator leaves no gap when collapsed; drag it or use arrow keys when focused to resize the sidebar within configurable minimum and maximum widths. Set a default width and restore it by double-clicking the separator. Animations respect Reduced Motion. On macOS, the sidebar exposes the window's blurred background up to the top edge, while main content remains opaque. The toggle keeps clearance from traffic lights. On Windows, a theme-aware gray sidebar background distinguishes navigation from the main content.
- Add icon-only toolbar actions through `shell.IconButton` with accessible labels, tooltips and translucent text-color hover/press feedback. Place application actions in Shell's fixed left bottom toolbar, and supply a title, breadcrumbs or actions in the leading main-content titlebar area.
- Replace the entire sidebar-and-content page with application content such as settings, then restore the previous sidebar state.
- Try Shell in the standalone gallery, including light/dark appearances, sidebar controls and a replacement settings page.
- Use oneo-go's built-in light and dark themes, applied through MyGo's theme mechanism so all MyGo widgets keep working. `ShellOptions.Appearance` selects `theme.Auto` (the default, follows the system appearance), `theme.Light` or `theme.Dark`. The palettes follow the Light and Shadow (light-shadow) design style from oneo-me/desgin-styles: a warm white light mode and a deep navy dark mode; preview both in the gallery with `-theme light` or `-theme dark`.
