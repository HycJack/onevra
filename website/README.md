# Onevra website

The Onevra marketing site: a statically prerendered React Router app that is
deployed to GitHub Pages at <https://oneo-me.github.io/onevra/>. Its
[components page](app/routes/components.tsx) runs the Gallery's demonstrations
in the browser as WebAssembly.

- **Framework:** React Router in framework mode with `ssr: false` and full
  prerendering, so every route ships as static HTML.
- **Styling:** Tailwind CSS v4 with the `nova` shadcn/ui preset on Base UI.
- **Content:** English and Simplified Chinese, defined in
  [`app/lib/content.ts`](app/lib/content.ts).

## Requirements

- [Bun](https://bun.sh) 1.4 or newer.
- [Go](https://go.dev) 1.27.1 or newer, which compiles the component preview to
  WebAssembly.

## Commands

```sh
bun install           # install dependencies
bun run dev           # start the dev server at http://localhost:5173/onevra/
bun run typecheck     # generate route types and typecheck
bun run build:pages   # build the static artifact into build/pages
bun run preview       # serve the built artifact the way GitHub Pages does
bun run preview:manifest # regenerate the component catalog from Go
bun run format        # format with Prettier
```

The site is served under the `/onevra/` base path, so the dev server uses
<http://localhost:5173/onevra/> rather than the root.

> `dev` runs the React Router CLI through Bun with the `development`
> condition. Invoking `react-router dev` directly makes its self-restart loop
> under Bun, so keep the explicit invocation.

## Deployment

The [Website workflow](../.github/workflows/website.yml) builds the site and
publishes it with GitHub Pages whenever `website/**` or `source/**` changes on
`main`. It also checks that the committed component catalog matches the Go
source.

Enable it once under **Settings → Pages → Build and deployment → Source →
GitHub Actions**. GitHub Pages is only available for private repositories on
paid plans; the site keeps working unchanged once the repository is public.

## Base path

`BASE_PATH` overrides the deploy path in both `vite.config.ts` and
`react-router.config.ts`. It defaults to `/onevra/`, which matches the
`oneo-me/onevra` repository. Set it to `/` when serving from a custom domain.

React Router writes prerendered pages under `build/client/onevra/`, while Vite
keeps assets at the root of `build/client`. `scripts/prepare-pages.mjs` merges
both into `build/pages` and adds a `404.html` fallback, which is the directory
uploaded to Pages.

## Component previews

`/components` and `/zh/components` list the Gallery's component
demonstrations and frame a live preview. The preview is the Gallery's own Go app
for browsers,
[`source/examples/gallery/web`](../source/examples/gallery/web/main.go), built by
[`scripts/preview-wasm.ts`](scripts/preview-wasm.ts).

- The Vite plugin compiles the app to WebAssembly during development and in the
  Pages build, and emits `preview/index.html`, Go's `wasm_exec.js` and a
  content-hashed `.wasm` file. Go must be installed; the site does not ship a
  checked-in build.
- The Go app selects a page from its command-line arguments. Changing components
  changes only the frame's fragment, so the running app redraws the next page
  instead of loading the WebAssembly again. The frame follows the site's light or
  dark appearance, and the language comes from the page route.
- [`app/lib/preview-manifest.json`](app/lib/preview-manifest.json) is generated
  from `go run ./examples/gallery/web -manifest`. Run
  `bun run preview:manifest` after adding or renaming a Gallery component page,
  and commit the result.
- The previews share the Gallery's demonstrations rather than duplicating them.
  Native window chrome, menus, trays, notifications and file dialogs report that
  a browser cannot provide them.

## Structure

| Path | Purpose |
| --- | --- |
| `app/root.tsx` | Document shell, theme bootstrap and error boundary. |
| `app/routes.ts` | Route table: `/`, `/components` and their `/zh` pages. |
| `app/components/landing.tsx` | The landing page sections. |
| `app/components/component-browser.tsx` | The components page and its preview frame. |
| `app/lib/content.ts` | English and Simplified Chinese copy and code samples. |
| `app/lib/preview.ts` | Component catalog lookup, slugs and the frame URL. |
| `app/lib/preview-manifest.json` | Generated component pages from the Go Gallery. |
| `app/lib/site.ts` | Shared links and the locale helper. |
| `app/components/ui` | shadcn/ui components. |
| `preview/index.html` | The page that loads the WebAssembly preview app. |
| `scripts/preview-wasm.ts` | Vite plugin that compiles and serves the preview. |
| `scripts/preview-manifest.mjs` | Regenerates the component catalog from Go. |
| `app/assets/icon.png` | The app icon, used as the header/footer logo and favicon. |

`app/assets/icon.png` is downscaled from `../source/assets/icon.png`.
Regenerate it rather than editing it by hand.

Add or update shadcn/ui components with the project's Bun runner:

```sh
bun x --bun shadcn@latest add <component>
```
