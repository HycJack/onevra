# Component design

The public [components package](../../source/components/doc.go) provides native UI compositions and experimental system-view hosting, not an application framework. Shell owns window presentation and a replaceable page container.

## Composition and ownership

Shell creates and configures the MyGo window because transparency, title-bar style and traffic-light placement must be set at creation time. It exposes MyGo's original handle for application events and lifecycle operations, rather than introducing a second event or shutdown system. Applications should not replace the window's content behind Shell.

Applications own business state, localized text, actions and page selection. Content builders and Shell state changes run synchronously on the UI thread; builders must not launch repeated background work during rendering. Shell owns only presentation state, does not persist it, and uses MyGo's rendering, layout and input systems.

For basic usage, see the [README](../../README.md#basic-usage); public contracts and defaults belong to the [API source](../../source/components/shell.go).

## Page container and layout

The page container reserves space for system window controls and overlays a transparent draggable header with a Lucide sidebar control. The toggle keeps additional clearance from macOS traffic lights. Its default page contains a collapsible left sidebar and a flexible main region. Both regions extend to the top of the window behind the header; their content stays below it. Applications can supply the leading main-content titlebar area for a title, breadcrumbs or actions. It aligns with the main surface, clears the sidebar toggle when collapsed, and clips within the remaining native-control safe area.

The optional left bottom toolbar remains beneath the sidebar body and disappears with the sidebar or on replacement pages. It composes MyGo's toolbar, retaining its focus navigation and overflow behavior. Applications supply actions and labels; the reusable [icon button](../../source/components/toolbar.go) uses MyGo button semantics with transparent resting faces and translucent text-color interaction feedback.

The left sidebar keeps its content at the full sidebar width while its visible slice animates, so the opaque main region appears to cover or uncover it without resizing the window. A thin separator overlays the main region's left edge, takes no layout space and disappears when the sidebar is fully closed. Its overlaid drag handle sits below the window-control region and supports keyboard arrow adjustment; double-clicking restores the application-configured default width. Applications can bound resizing with minimum and maximum widths. Shell retains the chosen width across toggles and page changes and removes the handle while collapsed or transitioning. Narrow layouts prioritize room for the main region over the configured minimum without discarding the chosen width; widening the window restores it.

A replacement page occupies the entire body beneath the shared window header and hides the default page's titlebar content and sidebar controls. Returning restores the previous sidebar state without changing window geometry. Shell provides this switching mechanism, not a settings implementation or business router.

MyGo's animation values drive the sidebar reveal, without a separate timer or animation loop. Reversing a toggle mid-animation continues from the current visible width; Reduced Motion skips intermediate geometry.

## Embedded WebView

The experimental [WebView component](../../source/components/webview.go) places a MyGo-owned system browser in a native window's rectangular layout slot. Applications create and retain instances outside view construction, supply content through MyGo's original Page API, and explicitly hide a view before omitting its slot on navigation. Hiding preserves the document; closing releases the browser without closing the host. Host shutdown closes attached instances automatically. Multiple instances can occupy separate slots.

On macOS ARM64, the component reparents the content view of a hidden MyGo web window into the visible host. The hidden window remains the browser's lifecycle owner and must not be manipulated by callers. MyGo still owns page loading, rendering and input. Geometry changes are coalesced through MyGo's UI queue outside paint; cleanup restores the donor's view before destroying it. This bridge relies on the pinned backend's AppKit hierarchy and needs review on upstream upgrades. Other platforms return an unsupported error without creating a browser.

This is OS view composition, not a shared GPU scene. The browser sits above MyGo's native Surface, so MyGo overlays cannot cover it and their modal input barriers do not block it. Applications must explicitly hide affected WebViews for native in-surface dialogs. Ancestor scrolling/clipping of the browser slot, arbitrary transforms and unified Tab navigation are not supported. Window-bound browser dialogs, menus, popups and fullscreen still target the hidden owner and are outside the component's scope. Input methods, edit shortcuts, accessibility and live-resize synchronization need further native review. MyGo's content capture omits these system views; review an OS-composited capture instead.

See the [usage example](../../source/components/webview_example_test.go) and [gallery experiment](gallery.md#hybrid-experiment). The latter deliberately demonstrates both overlay occlusion and explicit hiding; neither is presented as unified composition.

## Native appearance

MyGo and the operating system provide window controls and interaction semantics. Shell uses the active theme rather than creating a separate theme manager or custom window chrome.

On macOS, a transparent native window and transparent UI root expose a blurred system material. The left sidebar paints no background. Main content and replacement pages remain opaque, extending behind the unpainted header. Other platforms retain an opaque theme background rather than claiming unverified blur support.

Application builders can obscure the material by painting their own backgrounds. Preserve transparency through the window and sidebar when changing styling; content captures alone do not include the OS-composited blur.
