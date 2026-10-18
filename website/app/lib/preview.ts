import manifest from "./preview-manifest.json"
import type { Locale } from "./site"

export interface PreviewEntry {
  slug: string
  key: string
  name: string
  nameZh: string
  group: string
  groupZh: string
  variants?: string[]
  components: string[]
}

export interface PreviewGroup {
  key: string
  label: string
  entries: PreviewEntry[]
}

/** Every page the Gallery can render for the website. */
export const previewEntries: PreviewEntry[] = manifest

/** The page shown when the URL names none, or an unknown one. */
export const defaultPreviewSlug = "button"

/**
 * Pages by their own slug and by every component they demonstrate, so a
 * link to /components/calendar opens the date and time pickers page.
 */
const previewAliases = new Map<string, PreviewEntry>()
for (const entry of previewEntries) {
  previewAliases.set(entry.slug, entry)
  for (const component of entry.components) {
    previewAliases.set(previewAlias(component), entry)
  }
}

export function findPreview(slug: string | null): PreviewEntry | undefined {
  return slug ? previewAliases.get(slug.toLowerCase()) : undefined
}

/** The page demonstrating a component, by its Go name such as "TextBox". */
export function findComponentPreview(name: string): PreviewEntry | undefined {
  return previewAliases.get(previewAlias(name))
}

function previewAlias(value: string): string {
  return value
    .replace(/([a-z0-9])([A-Z])/g, "$1-$2")
    .replace(/([A-Z]+)([A-Z][a-z])/g, "$1-$2")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
}

/** Groups pages in their Gallery order, labels in the requested language. */
export function previewGroups(locale: Locale): PreviewGroup[] {
  const groups = new Map<string, PreviewGroup>()
  for (const entry of previewEntries) {
    let group = groups.get(entry.group)
    if (!group) {
      group = {
        key: entry.group,
        label: locale === "zh" ? entry.groupZh : entry.group,
        entries: [],
      }
      groups.set(entry.group, group)
    }
    group.entries.push(entry)
  }
  return [...groups.values()]
}

export function previewName(entry: PreviewEntry, locale: Locale): string {
  return locale === "zh" ? entry.nameZh : entry.name
}

/** The site path of a component page, or of the page list without a slug. */
export function componentsHref(locale: Locale, slug?: string): string {
  const base = locale === "zh" ? "/zh/components" : "/components"
  return slug ? `${base}/${slug}` : base
}

/**
 * The frame's URL. Only its fragment changes between component pages, so the
 * running WebAssembly app redraws instead of loading again.
 */
export function previewFrameUrl(locale: Locale, slug: string): string {
  const language = locale === "zh" ? "zh-Hans" : "en"
  return `${import.meta.env.BASE_URL}preview/index.html?lang=${language}#${encodeURIComponent(slug)}`
}
