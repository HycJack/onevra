# Component design

The public [components package](../../source/components/doc.go) provides native UI compositions, not an application framework. Shell owns window presentation and a replaceable page container.

## Composition and ownership

Shell creates and configures the MyGo window because transparency, title-bar style and traffic-light placement must be set at creation time. It exposes MyGo's original handle for application events and lifecycle operations, rather than introducing a second event or shutdown system. Applications should not replace the window's content behind Shell.

Applications own business state, localized text, actions and page selection. Content builders and Shell state changes run synchronously on the UI thread; builders must not launch repeated background work during rendering. Shell owns only presentation state, does not persist it, and uses MyGo's rendering, layout and input systems.

For basic usage, see the [README](../../README.md#basic-usage); public contracts and defaults belong to the [API source](../../source/components/shell.go).

## Page container and layout

The page container reserves space for system window controls and overlays a transparent draggable header with Lucide sidebar controls. The left toggle keeps additional clearance from macOS traffic lights; the inspector toggle sits at the main content's right edge, rather than following the inspector's outer edge. In tight windows, native-control clearance takes priority. Its default page contains a collapsible left sidebar, a flexible main region and an optional right inspector. All three regions extend to the top of the window behind the header; their content stays below it. Without inspector content, there is no inspector toggle. Applications can supply the leading main-content titlebar area for a title, breadcrumbs or actions. It aligns with the main surface, clears the leading sidebar toggle when collapsed, and clips within the space left by the controls.

The optional left bottom toolbar remains beneath the sidebar body and disappears with the sidebar or on replacement pages. It composes MyGo's toolbar, retaining its focus navigation and overflow behavior. Applications supply actions and labels; the reusable [icon button](../../source/components/toolbar.go) uses MyGo button semantics with transparent resting faces and translucent text-color interaction feedback.

The left sidebar keeps its content at the full sidebar width while its visible slice animates, so the opaque main region appears to cover or uncover it without resizing the window. The right inspector likewise reveals full-width content through the growing window edge, preserving the window's origin and height and, when space allows, the existing main region. Neither sidebar reflows its content to the animated slice width. Thin separators overlay the main region's left edge and the inspector's left edge; they take no layout space and disappear when the corresponding sidebar is fully closed. The left separator has an overlaid drag handle, below the window-control region, with keyboard arrow adjustment. Shell retains the chosen width across toggles and page changes, constrains it to available room, and removes the handle while collapsed or transitioning. Shell records the actual width added, so contraction does not subtract space the OS refused to grant. User resizing while the inspector is visible remains reflected in the base window size. Narrow layouts constrain sidebars to leave room for the main region. The left width limit uses the viewport without the inspector's native extension, so opening the inspector does not widen a constrained left sidebar.

A replacement page occupies the entire body beneath the shared window header and hides the default page's titlebar content and both sidebar controls. Entering it retracts the inspector's extension without discarding sidebar expansion state; returning restores the default layout. Shell provides this switching mechanism, not a settings implementation or business router.

MyGo's animation values drive sidebar reveals, without a separate timer or animation loop. Native window changes are coalesced and dispatched to the UI thread outside the current build/layout/paint transaction. Calling native resize within a frame would leave that frame painting an old viewport into the resized window and prevent a synchronous native redraw. Layout therefore follows the applied resize step, not the animation's next requested width. The target is published before resizing, and synchronous native redraws measure the actual granted extension, keeping the left and main surfaces stationary when the window can grow freely.

State changes request a new target; geometry settles across subsequent frames. Reversing a toggle or switching pages mid-animation continues from the current width. Reduced Motion skips intermediate geometry; native changes still use the deferred UI-thread dispatch. MyGo and the OS own presentation and compositing; frame geometry checks do not establish compositor behavior on every display or platform.

Maximized and full-screen windows cannot grow outward. The inspector may display inside their available space; width adjustments are deferred until a subsequent normal-window frame. Screen bounds and native resizing constraints can also prevent exact width preservation.

## Native appearance

MyGo and the operating system provide window controls and interaction semantics. Shell uses the active theme rather than creating a separate theme manager or custom window chrome.

On macOS, a transparent native window and transparent UI root expose a blurred system material. The left sidebar paints no background. The inspector uses a lightly translucent theme background; main content and replacement pages remain opaque, extending behind the unpainted header. Other platforms retain an opaque theme background rather than claiming unverified blur support.

Application builders can obscure the material by painting their own backgrounds. Preserve transparency through the window and sidebars when changing styling; content captures alone do not include the OS-composited blur.
