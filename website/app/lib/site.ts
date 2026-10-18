export const site = {
  name: "Onevra",
  github: "https://github.com/oneo-me/onevra",
  issues: "https://github.com/oneo-me/onevra/issues",
  pkg: "https://pkg.go.dev/github.com/oneo-me/onevra/source",
  architecture: "https://github.com/oneo-me/onevra/blob/main/ARCHITECTURE.md",
  examples:
    "https://github.com/oneo-me/onevra/tree/main/source/examples/gallery",
  gallery:
    "https://github.com/oneo-me/onevra/blob/main/docs/architecture/gallery.md",
  changelog: "https://github.com/oneo-me/onevra/blob/main/CHANGELOG.md",
  license: "https://github.com/oneo-me/onevra/blob/main/LICENSE",
} as const

export type Locale = "en" | "zh"

/** Derives the locale from a router pathname with the basename stripped. */
export function localeFromPathname(pathname: string): Locale {
  const path = pathname.replace(/\/+$/, "")
  return path === "/zh" || path.startsWith("/zh/") ? "zh" : "en"
}
