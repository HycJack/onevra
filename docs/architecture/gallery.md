# Gallery design

The gallery is a standalone example and a native review surface for the public components and icons. It must remain runnable without any consumer application.

## Separation from the library

The [entry point](../../source/examples/gallery/main.go) owns application startup, appearance selection and shutdown. The [internal gallery](../../source/examples/gallery/internal/gallery/gallery.go) owns demonstration content and in-memory state. Neither is part of the public component API.

Examples should demonstrate ordinary consumer integration, not introduce gallery-only behavior into components. The gallery does not read or write user content or application settings. Its grouped sidebar navigates between oneo-go demonstrations, a Lucide icon library with icon-only outline buttons that copy their names and MyGo's built-in styled widgets and layout primitives. The oneo-go controls group exercises the public `ui` package on separate interactive pages, with state isolated from the upstream examples and between variants. MyGo entries are grouped by type and open separate control-only pages, without source listings or tutorial prose. State and usage variants remain interactive where enabled, with independent in-memory state per page and variant. The catalog follows the pinned upstream UI guide; unstyled Base APIs remain building blocks rather than a duplicate widget catalog. The exhaustive catalog belongs only to the gallery and references public named accessors; ordinary library consumers retain no runtime registry. MyGo's virtualized grid builds visible icons on demand. The Shell demonstration also shows sample icons at multiple sizes and inherited text colors.

The Terminal plugin page starts only on explicit request and uses an in-memory echo connection, never a system shell. MyGo owns emulation and input; the gallery owns asynchronous native-library loading and closes the connection when stopped or when the window closes. Ordinary tests do not load or download the terminal's native library.

Packaging and native-library path adaptation belong to the gallery application, not the library. Startup selects the bundled Terminal library before the page can load it, without forcing eager native-library initialization. Build commands and distribution limitations are maintained in the [development guide](../build.md#app-bundles).

## Verification approach

Ordinary tests do not open windows. The opt-in native smoke test runs the actual MyGo event loop with temporary user data and captures light, dark, narrow and replacement-page views. Native bounds assertions check that replacement pages preserve window geometry. Headless tests verify sidebar width limits and double-click reset; the native pointer path requires manual review. See [native review](../build.md#native-review) for the command and manual checks.

Content captures establish rendering but omit OS-composited background blur. Review blending in a live window over a contrasting backdrop; a content capture alone cannot prove it works. Keyboard interaction, accessibility and platform-specific window behavior also require native review rather than screenshot assertions.
