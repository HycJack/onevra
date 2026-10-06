# Gallery

`source/examples/gallery/` is the standalone native entry point. `source/internal/gallery/` contains its view and demonstration state, not public component API.

The executable accepts `-theme system`, `-theme light` or `-theme dark` and rejects other values. It opens an 820 × 900 window with a 360 × 480 minimum. Closing the window quits the app. The view uses a scrolling, bounded column with headings and dividers; it demonstrates empty state, notices and section composition.

Create item and Reset example change in-memory demonstration state only. A checkbox disables the create action through MyGo. Mixed English/Chinese descriptions exercise wrapping. The gallery does not read or write user content or application settings. Its icon samples use named public accessors at multiple sizes and inherited text colors; they are not an exhaustive runtime catalog.

## Application packaging

`source/mygo.json` points MyGo's pinned Go CLI tool at `examples/gallery` without adding another module or changing the library package. Run Go commands from `source/`. `go tool mygo dev` runs a live-reloading development app bundle on macOS; `go tool mygo build` creates host-platform artifacts under the repository-root `dist/`. VS Code tasks use the same entry points. Application identity, the local-only bundle version and output paths belong to this configuration, not public components. See [app bundles](../build.md#app-bundles) for paths and signing limitations.

## Tests and captures

Headless tests exercise create, reset and disabled-action behavior. The opt-in native test runs MyGo's actual event loop with isolated temporary user data and captures light, dark and 360-pixel-wide windows, then destroys the window. Enable it with `ONEO_GO_NATIVE_SMOKE=1`; `ONEO_GO_SCREENSHOTS` selects a directory for retained PNGs. Ordinary tests never open windows.

See [development commands](../build.md#native-review). Native PNG decoding and bounds checks establish startup and rendering, not complete interaction acceptance. Visually inspect both themes, long text, narrow widths and enlarged text; manually review keyboard traversal, disabled controls, screen-reader behavior and scrolling. Windows and Linux need their own desktop sessions.
