# Gallery design

The gallery is a standalone example and a native review surface for the public components and icons. It must remain runnable without any consumer application.

## Separation from the library

The [entry point](../../source/examples/gallery/main.go) owns application startup, appearance selection and shutdown. The [internal gallery](../../source/internal/gallery/gallery.go) owns demonstration content and in-memory state. Neither is part of the public component API.

Examples should demonstrate ordinary consumer integration, not introduce gallery-only behavior into components. The gallery does not read or write user content or application settings. Its grouped sidebar navigates between oneo-go demonstrations, a searchable Lucide icon library and MyGo's built-in styled widgets and layout primitives. MyGo entries are grouped by type and open separate control-only pages, without source listings or tutorial prose. State and usage variants remain interactive where enabled, with independent in-memory state per page and variant. The catalog follows the pinned upstream UI guide; unstyled Base APIs remain building blocks rather than a duplicate widget catalog. The exhaustive catalog belongs only to the gallery and references public named accessors; ordinary library consumers retain no runtime registry. MyGo's virtualized grid builds visible icons on demand. The Shell demonstration also shows sample icons at multiple sizes and inherited text colors.

The Terminal plugin page starts only on explicit request and uses an in-memory echo connection, never a system shell. MyGo owns emulation and input; the gallery owns asynchronous native-library loading and closes the connection when stopped or when the window closes. Ordinary tests do not load or download the terminal's native library.

Packaging belongs to the gallery application, not the library. Build commands and distribution limitations are maintained in the [development guide](../build.md#app-bundles).

## Hybrid experiment

The **Experiments / Native + WebView** page consumes the public [WebView component](components.md#embedded-webview) in Gallery's main Shell window. It creates the browser lazily outside view construction, loads only local demonstration HTML, hides it when another page or Settings is shown, and restores the same document on return. Manual hide state is preserved across navigation. Unsupported platforms show an explanation rather than opening another window.

Two MyGo AlertDialog actions compare a live browser beneath the dialog with explicit hiding until dismissal. The live-browser alert automatically dismisses after eight seconds. It covers and blocks the live browser using the component's GPU composition path. The existing page displays asynchronous attachment errors; there is no alternate Gallery or legacy composition mode. Leaving the page dismisses the alert. Gallery owns this demonstration policy, not the WebView component.

The macOS opt-in component tests cover multiple views, geometry, AppKit hit testing, web mouse and basic key input, focus recovery, state retention and cleanup. Hybrid checks add overlay click blocking, modal focus takeover while the page keeps running, automatic unmounting and ancestor clipping. The Windows hybrid component and D3D11 composition tests are available for target-desktop review, not a substitute for it. The gallery test covers navigation, Settings, Go-to-JavaScript updates and both alert modes. Optional screenshots capture the OS-composited window, not MyGo's Surface-only content capture. See [running and reviewing the experiment](../build.md#hybrid-experiment).

## Verification approach

Ordinary tests do not open windows. The opt-in native smoke test runs the actual MyGo event loop with temporary user data and captures light, dark, narrow and replacement-page views. Native bounds assertions check that replacement pages preserve window geometry. A macOS ARM64 pointer test verifies sidebar width limits and double-click reset through the native event path. See [native review](../build.md#native-review) for the command and manual checks.

Content captures establish rendering but omit OS-composited background blur. Review blending in a live window over a contrasting backdrop; a content capture alone cannot prove it works. Keyboard interaction, accessibility and platform-specific window behavior also require native review rather than screenshot assertions.
