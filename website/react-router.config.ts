import type { Config } from "@react-router/dev/config"

import manifest from "./app/lib/preview-manifest.json"

// GitHub Pages serves the site from https://<owner>.github.io/<repo>/.
// Override BASE_PATH to deploy under a custom domain or path.
const basePath = process.env.BASE_PATH ?? "/onevra/"

export default {
  basename: basePath,
  ssr: false,
  // Every component page ships as its own static HTML file, so a shared link
  // opens directly with its own title instead of through the 404 fallback.
  async prerender({ getStaticPaths }) {
    const pages = manifest.flatMap(({ slug }) => [
      `/components/${slug}`,
      `/zh/components/${slug}`,
    ])
    return [
      ...new Set([
        ...getStaticPaths(),
        "/components",
        "/zh/components",
        ...pages,
      ]),
    ]
  },
} satisfies Config
