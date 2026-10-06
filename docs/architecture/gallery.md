# Gallery design

The gallery is a standalone example and a native review surface for the public components and icons. It must remain runnable without any consumer application.

## Separation from the library

The [entry point](../../source/examples/gallery/main.go) owns application startup, appearance selection and shutdown. The [internal gallery](../../source/internal/gallery/gallery.go) owns demonstration content and in-memory state. Neither is part of the public component API.

Examples should demonstrate ordinary consumer integration, not introduce gallery-only behavior into components. The gallery does not read or write user content or application settings. Its grouped sidebar navigates between component demonstrations and a searchable Lucide icon library. The exhaustive catalog belongs only to the gallery and references public named accessors; ordinary library consumers retain no runtime registry. MyGo's virtualized grid builds visible icons on demand. The Shell demonstration also shows sample icons at multiple sizes and inherited text colors.

Packaging belongs to the gallery application, not the library. Build commands and distribution limitations are maintained in the [development guide](../build.md#app-bundles).

## Verification approach

Ordinary tests do not open windows. The opt-in native smoke test runs the actual MyGo event loop with temporary user data and captures light, dark, narrow, inspector and replacement-page views. Native bounds assertions check that the inspector extends the window to the right and that replacement pages retract that extension. A separate native frame test samples every painted layout during expansion and collapse in normal and narrow windows, checking viewport/native-size agreement and stable left/main bounds; it does not sample the OS compositor. See [native review](../build.md#native-review) for the command and manual checks.

Content captures establish rendering but omit OS-composited background blur. Review blending in a live window over a contrasting backdrop; a content capture alone cannot prove it works. Keyboard interaction, accessibility and platform-specific window behavior also require native review rather than screenshot assertions.
