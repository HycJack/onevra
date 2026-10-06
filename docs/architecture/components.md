# Component design

The public [components package](../../source/components/doc.go) provides native UI compositions, not an application framework. Shell owns window presentation and a replaceable page container.

## Composition and ownership

Shell creates and configures the MyGo window because transparency, title-bar style and traffic-light placement must be set at creation time. It exposes MyGo's original handle for application events and lifecycle operations, rather than introducing a second event or shutdown system. Applications should not replace the window's content behind Shell.

Applications own business state, localized text, actions and page selection. Content builders and Shell state changes run synchronously on the UI thread; builders must not launch repeated background work during rendering. Shell owns only presentation state, does not persist it, and uses MyGo's rendering, layout and input systems.

For basic usage, see the [README](../../README.md#basic-usage); public contracts and defaults belong to the [API source](../../source/components/shell.go).

## Page container and layout

The page container reserves space for system window controls and overlays a transparent draggable header with Lucide sidebar controls. Its default page contains a collapsible left sidebar, a flexible main region and an optional right inspector. All three regions extend to the top of the window behind the header; their content stays below it. Without inspector content, there is no inspector toggle.

The left sidebar keeps its content at the full sidebar width while its visible slice animates, so the opaque main region appears to cover or uncover it without resizing the window. The right inspector likewise reveals full-width content through the growing window edge, preserving the window's origin and height and, when space allows, the existing main region. Neither sidebar reflows its content to the animated slice width. Thin separators overlay the main region's left edge and the inspector's left edge; they take no layout space and disappear when the corresponding sidebar is fully closed. The left separator has an overlaid drag handle, below the window-control region, with keyboard arrow adjustment. Shell retains the chosen width across toggles and page changes, constrains it to available room, and removes the handle while collapsed or transitioning. Shell records the actual width added, so contraction does not subtract space the OS refused to grant. User resizing while the inspector is visible remains reflected in the base window size. Narrow layouts constrain sidebars to leave room for the main region.

A replacement page occupies the entire body beneath the shared window header and hides both sidebar controls. Entering it retracts the inspector's extension without discarding sidebar expansion state; returning restores the default layout. Shell provides this switching mechanism, not a settings implementation or business router.

MyGo's animation values drive sidebar reveals and native window width adjustments on the UI thread, without a separate timer or animation loop. Native resizing and UI viewport updates are not frame-synchronized, so right-inspector transitions can visibly stretch or reflow the main content; the current animation does not eliminate that limitation. State changes request a new target; geometry settles across subsequent frames. Reversing a toggle or switching pages mid-animation continues from the current width. Reduced Motion applies target geometry immediately, including when enabled during a transition.

Maximized and full-screen windows cannot grow outward. The inspector may display inside their available space; width adjustments are deferred until a subsequent normal-window frame. Screen bounds and native resizing constraints can also prevent exact width preservation.

## Native appearance

MyGo and the operating system provide window controls and interaction semantics. Shell uses the active theme rather than creating a separate theme manager or custom window chrome.

On macOS, a transparent native window and transparent UI root expose a blurred system material. The left sidebar paints no background. The main content, inspector and replacement pages use opaque theme backgrounds extending behind the unpainted header. Other platforms retain an opaque theme background rather than claiming unverified blur support.

Application builders can obscure the material by painting their own backgrounds. Preserve transparency through both the window and left sidebar when changing styling; content captures alone do not include the OS-composited blur.
