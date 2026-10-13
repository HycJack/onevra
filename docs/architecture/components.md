# Component design

Each shell style lives in its own package under `shells/`, without a shared facade or style selector. The public [sidebar shell package](../../source/shells/sidebar/doc.go) owns a replaceable sidebar-and-content container; the [plain shell package](../../source/shells/plain/shell.go) owns a single-content window. Both own native presentation, not application lifetime.

The [components package](../../source/components/doc.go) provides common controls independently of Shell. Sidebar Shell also provides icon buttons for titlebar and toolbar actions. Both packages compose MyGo directly.

Shell reuses the common Icon component; components remain independent of Shell.

Common controls expose at most one visual variant axis, such as button style. Each control has one theme-scaled standard size rather than a separate size enumeration. Interaction and value states are not visual variants.

## Composition and ownership

Shell creates and configures the MyGo window because transparency, title-bar style and traffic-light placement must be set at creation time. It exposes MyGo's original handle for application events and lifecycle operations, rather than introducing a second event or shutdown system. Both shell styles can start hidden so applications can prepare content and lifecycle handlers before choosing whether to show or focus the window. Applications should not replace the window's content behind Shell.

Applications own business state, localized text, actions and page selection. The reusable [Sidebar](../../source/components/sidebar.go) composes grouped navigation independently of Shell; it supplies item appearance, while applications handle activation and scrolling and Shell owns the window's sidebar geometry. Content customization uses product-neutral slots rather than domain-specific options: navigation can accept application-built leading content while retaining its own selection and activation.

Content builders and Shell state changes run synchronously on the UI thread; builders must not launch repeated background work during rendering. MyGo's value-based Element handles belong to the current build pass, not application state. Use a bound MyGo Handle for identity or geometry queries between builds.

Shell owns only presentation state, does not persist it, and uses MyGo's rendering, layout and input systems. An application-chosen `StateKey` passes through to MyGo, which owns capturing, restoring and storing window geometry; sidebar width, expansion and pages remain application state.

Components read the frame theme and style their own elements; they must not call `Context.SetTheme`, which also resets the window root's appearance. Only application-level theme installation owns that operation. Collection controls compose their own rows and component-level selection/navigation on MyGo's input and focus APIs.

Virtualized collections reuse non-selecting MyGo list and grid containers for scrolling and layout, rather than copying its layout engine or overriding its styled row theme. Their application-owned state is defined in the [collection API](../../source/components/list.go). Tables share MyGo's public column, sort and layout types and its sideways row-width scrolling; Onevra owns the header's sorting, moving and resizing, while applications sort rows and may persist the layout.

Double-click column fitting and pinned section headers are not provided. Public role, selection, level and active-descendant metadata describe the rows, but this composition does not claim full parity with MyGo's private collection accessibility actions; native assistive-technology review is still required.

For basic usage, see the [README](../../README.md#getting-started); public contracts and defaults belong to the [API source](../../source/shells/sidebar/shell.go).

## Complex component composition

Complex components each live in an independent package under `source/components/<name>/`, with implementation, tests and package documentation kept together. These subpackages may compose common controls; the root `components` package must not depend on them. MyGo plugins continue to be used directly.

### Charts

The [chart component](../../source/components/chart/chart.go) composes categorical bars, lines, areas, stacked bars and donuts from MyGo's scales, paths, marks and selection. Onevra supplies muted theme colors, light grids, gradients, adaptive category labels, legends and selection details. The views use ordinary element layout and own no windows or background work.

Applications retain data and optional interaction state, provide stable series or slice identities, and supply localized labels and value formatting. Cartesian legends toggle series; donut legends select slices. Keyboard and pointer selection expose formatted values to assistive technology. Empty and detail slots let applications supply content without adding business-specific variants. See the [runnable usage example](../../source/components/chart/example_test.go) and API comments for data and state contracts.

Numeric updates transition together through MyGo's animation channel, including Cartesian domains and donut boundaries. Retargeting continues from the current picture; Reduced Motion skips the transition. Data-shape changes appear immediately, and selection details report the latest application values.

These views target bounded categorical datasets. Preparation and path construction remain proportional to the supplied data; applications own aggregation or sampling. Continuous/time axes, zooming and specialized charts can compose upstream primitives. The component reuses MyGo's Cartesian selection, with a coordinate adaptation for transformed plots in 0.4.3; radial selection composes MyGo's shortcuts and slice hit testing.

### Commands

The [command component](../../source/components/command/command.go) provides an embeddable palette and a Dialog composition. Applications retain query and active-command state, supply stable command identities and localized metadata, register window shortcuts, and handle activation. Hosts can show results immediately or start with a compact search field that expands after a committed query and collapses when cleared. Local metadata filtering is optional; asynchronous search, stale-result rejection, command execution and history remain application-owned. Loading and errors replace actionable results.

Search keeps keyboard focus while arrows move the active result through Onevra's non-selecting List; MyGo owns virtual layout and text editing. The dialog lets the modal scope establish initial focus so the overlay retains its original opener for restoration. Composition-aware dismissal prevents Escape ending an input-method composition from also closing the dialog. Public accessibility metadata exposes the active command, but native assistive-technology review remains necessary. See the [compilable usage example](../../source/components/command/example_test.go).

On macOS, palettes can optionally select an already enabled keyboard input source through the existing MyGo responder’s [AppKit input context](https://developer.apple.com/documentation/appkit/nstextinputcontext/selectedkeyboardinputsource). This small native adapter fills an upstream API gap without intercepting events or implementing text input. Selection waits for search focus and committed composition, runs once per opening, and normally restores the preceding source on dismissal if the active context and selected source still match.

Hosts may retain the selected source instead. Unavailable sources leave search usable and expose an error through retained state. Embedded hosts release the session when removing the palette; all hosts release it before window destruction.

See the [input-source lifecycle API](../../source/components/command/input_source.go).

### Read-only code

The [code component](../../source/components/code/code.go) tokenizes explicitly named languages with Chroma and renders retained spans through Onevra RichText and MyGo's text engine. Applications retain immutable documents outside frame builders; parsing can run in application-owned background work. Theme changes select prepared palettes without re-tokenization. Adjacent tokens sharing an appearance are merged, and a code block remains one text node rather than a UI element per token. Unknown languages, failed lexers and inputs exceeding highlighting limits retain their complete plain text.

This path targets code snippets and Markdown blocks. It preserves MyGo selection across the document, so it does not virtualize individual lines. Highlighting limits bound lexer input and styling, not the layout cost of arbitrarily large files or a hard regex execution deadline. Large-file viewers need a separate viewport strategy and an explicit whole-document selection contract.

### Markdown

[Markdown component](../../source/components/markdown/doc.go) parses CommonMark and GFM syntax into Onevra text and MyGo layout primitives. MyGo owns a single selection across the rendered document, including default code blocks. Parsing prepares fenced-code highlighting through the code component; indented code remains plain text.

Hosts retain parsed documents between source changes and provide ordinary scrolling so offscreen paragraphs remain selectable. Raw HTML stays literal. Standalone image paragraphs display application-supplied decoded bitmaps, shrinking to fit the document width; missing images and inline images retain alternative text.

Image loading and caching remain application-owned, and rendering performs no file or network I/O.

Applications may replace fenced code blocks with synchronous view builders selected by language, while retaining state and asynchronous work outside rendering. Custom text participates in the surrounding selection unless its builder chooses a separate selection scope or excludes it. Optional code headers expose language labels and keyboard-reachable copying without including the header in document selection.

Applications also own link activation; the component never opens destinations automatically. See the [gallery example](../../source/examples/gallery/internal/gallery/control_markdown.go) for an interactive embedded view.

### Conversations

The [chat component](../../source/components/chat/chat.go) composes conversation scrolling, turns, activity groups and a multiline composer whose proportions follow Waku's transcript. Applications own ordered steps with stable identities, stream updates, localized labels and summaries, retained Markdown documents and request lifecycles. A turn shows the prompt, then intermediate prose and activity groups, then an explicit final response; a completed tool never implicitly completes its turn. Builders run synchronously and perform no network or file work.

An activity group opens while it is the live tail of a running turn and closes once it settles, unless the user chose otherwise. Each activity is a compact row with an optional bounded detail viewport, prose or monospace, that keeps its own disclosure and scroll state. Successful turns with a final response fold their steps behind a centered summary row without building hidden content; running turns end with a pulsing working row.

Manual choices survive updates and outer folding; focus inside the steps or reading above the live edge defers automatic folding. Once folded, scrolling back does not reopen a turn. Footers and response actions remain outside the fold.

MyGo owns keyboard disclosure semantics, input methods and scrolling. Enter sends and Shift+Enter inserts a newline: MyGo's multiline editor consumes Enter before shortcuts, so the composer takes plain Enter through the editor's input handler and leaves active compositions alone. While running, an empty draft shows Stop and a nonempty draft shows only Send, leaving the application to decide what a new message does to the run.

Composition does not disable the Send button; Enter still belongs to the input method. The card accepts dropped file paths or host-selected native data formats without preparing them. A paste callback lets hosts attach clipboard content while preserving normal text editing when declined.

Decoded image preview tiles have optional removal actions; slots hold attachments, accessory controls, non-shrinking trailing controls beside Send or Stop, and a footer row. Hosts can disable sending while preparing attachments, without blocking editing. Applications accept the draft and perform cancellation themselves; model choice, permission policy and artifact data are not embedded in the component.

The conversation follows the bottom until the user scrolls away, with a jump control to resume. An optional floating layout lets the transcript scroll behind a glass composer while reserving clearance for the latest turn and fading to an opaque footer for readability. The composer applies MyGo's glass material only to its card, preserving the host's theme and root transparency.

Ordinary scrolling preserves document selection, but this version does not virtualize history or preserve a row anchor when history is prepended or arbitrary content above the viewport changes height. Applications bound transcript retention and tool previews, merge incoming updates before UI publication and cache parsing outside builders. Folded turns cost nothing per hidden step; expanded content and retained data still consume proportional work and memory.

The [gallery example](../../source/examples/gallery/internal/gallery/control_chat.go) demonstrates a static conversation without contacting a model.

## Localization

Applications own language choice and translation catalogs. `components.WithLocalization` supplies component-owned labels and date presentation within a synchronous, window-local build scope; components retain English defaults outside it. Explicit control labels still belong to the caller. Nested scopes restore their enclosing values and do not introduce global language state or change element identity. Calendar presentation can change without changing stored dates, time zones or selection semantics; numeric editors retain their canonical numeric syntax.

Shell accepts updated control labels without recreating its window. Applications must keep keys and selection values independent of translated text and request a new frame after changing language. This scope does not translate MyGo's internal widgets, OS surfaces, native menus or user-provided content. It provides no automatic RTL layout mirroring.

## Page container and layout

Plain Shell hides the system title bar and automatically places its content below a draggable strip with native macOS/Linux controls or Onevra caption buttons on Windows. Its opaque Onevra background covers the whole window, including the strip. Applications do not add another header inset. Live appearance and caption-label updates preserve the window and its content. Dialogs inherit the strip's drag behavior and clearance. The starter applications use this shell for main, settings and introduction windows; their transient tray panels remain application-owned frameless MyGo windows.

The sidebar page container reserves space for system window controls and overlays a transparent draggable header with a Lucide sidebar control. The toggle keeps additional clearance from macOS traffic lights. Its default page contains a collapsible left sidebar and a flexible main region.

Both regions and their application-built content extend to the top of the window behind the header. Applications own the top inset: use Shell's frame-aware header height to keep interactive content clear while painting backgrounds through the titlebar area. The header retains window dragging and control interaction above the body.

Applications can add trailing sidebar-header actions beside the toggle, clipped to the remaining sidebar safe area and omitted while collapsed or animating. Applications can supply the leading main-content titlebar area for a title, breadcrumbs or actions. It aligns with the main surface, clears the sidebar toggle when collapsed, and clips within the remaining native-control safe area.

The optional left bottom toolbar remains beneath the sidebar body and disappears with the sidebar or on replacement pages. It composes MyGo's toolbar, retaining its focus navigation and overflow behavior. Applications supply actions and labels; the reusable [icon button](../../source/shells/sidebar/icon_button.go) uses MyGo button semantics with transparent resting faces and translucent text-color interaction feedback.

The left sidebar keeps its content at the full sidebar width while its visible slice animates, so the opaque main region appears to cover or uncover it without resizing the window. A thin separator overlays the main region's left edge, takes no layout space and disappears when the sidebar is fully closed. Its overlaid drag handle spans the full window height, taking precedence over window dragging at the separator, and supports keyboard arrow adjustment; double-clicking restores the application-configured default width.

Applications can bound resizing with minimum and maximum widths. Shell retains the chosen width across toggles and page changes and removes the handle while collapsed or transitioning. The main region takes the remaining space rather than reserving a fixed share of the window.

Narrow layouts only constrain the sidebar to keep its resize handle inside the viewport, without discarding the chosen width; widening the window restores it.

A replacement page occupies the full window behind the shared floating header, with the same application-owned inset and hides the default page's titlebar content and sidebar controls. Returning restores the previous sidebar state without changing window geometry. Shell provides this switching mechanism, not a settings implementation or business router.

MyGo's animation values drive the sidebar reveal, without a separate timer or animation loop. Reversing a toggle mid-animation continues from the current visible width; Reduced Motion skips intermediate geometry.

## Settings composition

The public [Settings component](../../source/components/settings.go) renders application-defined categories, groups and setting rows independently of Shell or Dialog. Applications provide stable identities, localized metadata and synchronous control builders, and retain both presentation state and setting values.

Settings owns navigation, searchable presentation and row layout, including optional group header actions, which are not searched, and descriptions that applications may compose from inline text and links while search still matches their plain text; it does not load configuration, persist changes or prescribe when values take effect. Hosts choose the narrow presentation from the space allocated to the surface, so embedded settings do not incorrectly use the full window width.

Search spans all categories using metadata rather than building hidden controls. Results retain category and group context and allow editing in place; clearing search restores the selected category. Hosts can reuse the same metadata search without rendering the settings surface. Dialog and standalone-window lifecycles remain application-owned. See the [compilable examples](../../source/components/settings_example_test.go) for composition.

## Dialog composition

Dialog composes MyGo overlay, modal focus scope and transition primitives rather than a separate modal manager. Applications own open state and build layers in back-to-front order; nested dialogs retain the underlying page. Each layer controls backdrop and Escape dismissal independently.

AlertDialog composes the same surface and Onevra buttons, emitting a frame-local result instead of blocking the UI thread. Applications own localized labels, retain results when needed and provide scrolling for long dialog pages. MyGo owns focus trapping, restoration and reduced-motion behavior.

Shell supplies a window-local, synchronous build scope so every Onevra dialog, including nested AlertDialogs and dialogs on replacement pages, inherits its native window-drag strip. Standalone hosts can provide the same scope without depending on Shell; individual dialogs can override or disable dragging. The scope restores its previous value on return and does not introduce a modal manager.

Each modal layer reserves clearance and builds its own strip. This keeps underlying header actions blocked rather than forwarding pointer events through the modal; the strip does not dismiss the backdrop. Dragging uses MyGo's focus behavior, which may reset the focused input within the modal.

## Toast composition

The [toast component](../../source/components/toast.go) styles MyGo's toast queue through its unstyled viewport rather than keeping a separate notification manager. Applications add, update and close toasts with MyGo's context methods and own their text, actions and asynchronous work; MyGo owns timing, stacking, pausing, Escape dismissal and announcements.

Both shells build the viewport after their content, inside a localization scope that carries the Shell's close-button label, because application localization scopes end with their builders. Shell also owns placement, adding its header height to top placements so toasts clear the window controls. A view that builds its own viewport earlier in the frame keeps it. Status types select an icon; `Data` may carry a leading content slot instead of type-specific variants.

## Native appearance

MyGo and the operating system provide native window interaction semantics. On Windows, Shell uses MyGo's frameless window and composes rounded toolbar-style minimize, maximize/restore and close buttons in the shared header, including replacement pages. Labels are application-localizable; close requests retain MyGo's cancellable application lifecycle.

Maximize/restore follows the actual window state. Shell applies Onevra's built-in light and dark themes (see the [theme package](../../source/theme/doc.go)) through MyGo's own theme mechanism, so every MyGo widget keeps working; there is no separate theme manager or custom window chrome. `ShellOptions.Appearance` selects `theme.Auto` (the default, follows the system appearance), `theme.Light` or `theme.Dark`. A pinned theme ignores later appearance changes, while window chrome still follows the operating system.

These UI buttons do not implement Windows' native maximize-hover Snap Layouts flyout.

Sidebar Shell requests a sidebar material on macOS and Mica on Windows. MyGo supports Windows system materials on Windows 11 22H2 and later; it owns compositor setup and native opaque fallback. Shell uses `Context.Vibrancy()` each frame, rather than the operating system or requested material, to decide whether to expose the backdrop.

When effective, the root and sidebar preserve transparency without additional tint in either theme. Otherwise, including Linux or unavailable composition, Shell paints an opaque theme background and contrasting sidebar surface through the header and bottom toolbar. Main content and replacement pages always remain opaque.

Applications can change or disable the material through the exposed MyGo window's `SetVibrancy`; MyGo rebuilds the view and Shell reapplies the same policy. Windows material rendering still requires target-platform native review.

Application builders can obscure the material by painting their own backgrounds. Preserve transparency through the window and sidebar when changing styling; content captures alone do not include the OS-composited blur.
