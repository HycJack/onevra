import { ExternalLink } from "lucide-react"
import { Link, useSearchParams } from "react-router"

import { SiteFooter } from "~/components/site-footer"
import { SiteHeader } from "~/components/site-header"
import { Badge } from "~/components/ui/badge"
import { Button } from "~/components/ui/button"
import { componentsPage, content } from "~/lib/content"
import {
  defaultPreviewSlug,
  findPreview,
  previewEntries,
  previewFrameUrl,
  previewGroups,
  previewName,
  type PreviewEntry,
} from "~/lib/preview"
import { site, type Locale } from "~/lib/site"
import { cn } from "~/lib/utils"

export function ComponentBrowser({ locale }: { locale: Locale }) {
  const landing = content[locale]
  const page = componentsPage[locale]
  const groups = previewGroups(locale)
  const [params, setParams] = useSearchParams()
  const requested = params.get("c")
  const selected =
    findPreview(requested) ??
    findPreview(defaultPreviewSlug) ??
    previewEntries[0]
  const title = previewName(selected, locale)
  const frameSrc = previewFrameUrl(locale, selected.slug)

  function choose(slug: string) {
    const next = new URLSearchParams(params)
    next.set("c", slug)
    setParams(next, { preventScrollReset: true })
  }

  return (
    <div className="flex min-h-svh flex-col">
      <SiteHeader
        locale={locale}
        content={landing}
        localeHref={locale === "en" ? "/zh/components" : "/components"}
      />

      <main className="flex-1">
        <div className="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6 sm:py-14">
          <header className="flex max-w-3xl flex-col gap-4">
            <Badge variant="secondary" className="w-fit">
              {page.badge}
            </Badge>
            <h1 className="font-heading text-4xl font-semibold tracking-tight text-balance sm:text-5xl">
              {page.title}
            </h1>
            <p className="text-lg text-pretty text-muted-foreground">
              {page.description}
            </p>
            <div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-muted-foreground">
              <span>
                {previewEntries.length} {locale === "zh" ? "个页面" : "pages"}
              </span>
              <span aria-hidden>·</span>
              <span>{page.hint}</span>
            </div>
          </header>

          <section className="mt-10 grid gap-8 lg:grid-cols-[15rem_minmax(0,1fr)]">
            <aside className="lg:sticky lg:top-20 lg:max-h-[calc(100svh-6rem)] lg:self-start lg:overflow-y-auto lg:pr-2">
              <label className="flex flex-col gap-2 lg:hidden">
                <span className="text-xs font-medium text-muted-foreground">
                  {page.choose}
                </span>
                <select
                  value={selected.slug}
                  onChange={(event) => choose(event.target.value)}
                  className="h-9 w-full rounded-lg border border-border bg-background px-2.5 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
                >
                  {groups.map((group) => (
                    <optgroup key={group.key} label={group.label}>
                      {group.entries.map((entry) => (
                        <option key={entry.slug} value={entry.slug}>
                          {previewName(entry, locale)}
                        </option>
                      ))}
                    </optgroup>
                  ))}
                </select>
              </label>

              <nav
                aria-label={page.nav}
                className="hidden flex-col gap-5 lg:flex"
              >
                {groups.map((group) => (
                  <div key={group.key} className="flex flex-col gap-1">
                    <h2 className="px-2.5 font-heading text-xs font-medium tracking-wide text-muted-foreground uppercase">
                      {group.label}
                    </h2>
                    <ul className="flex flex-col gap-0.5">
                      {group.entries.map((entry) => (
                        <li key={entry.slug}>
                          <SidebarLink
                            entry={entry}
                            locale={locale}
                            selected={entry.slug === selected.slug}
                          />
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </nav>
            </aside>

            <div className="flex min-w-0 flex-col gap-3">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="flex flex-col gap-2">
                  <h2 className="font-heading text-xl font-semibold">
                    {title}
                  </h2>
                  <div className="flex flex-wrap gap-1.5">
                    {selected.components.map((name) => (
                      <code
                        key={name}
                        className="rounded-md border bg-muted/40 px-2 py-0.5 font-mono text-xs text-muted-foreground"
                      >
                        {name}
                      </code>
                    ))}
                  </div>
                </div>
                <div className="flex flex-wrap items-center gap-1">
                  <Button
                    variant="outline"
                    size="sm"
                    render={
                      <a href={frameSrc} target="_blank" rel="noreferrer" />
                    }
                  >
                    {page.open}
                    <ExternalLink data-icon="inline-end" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    render={
                      <a href={site.pkg} target="_blank" rel="noreferrer" />
                    }
                  >
                    {page.docs}
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    render={
                      <a
                        href={site.examples}
                        target="_blank"
                        rel="noreferrer"
                      />
                    }
                  >
                    {page.source}
                  </Button>
                </div>
              </div>

              <div className="overflow-hidden rounded-xl border bg-muted/20">
                <div className="h-[min(76svh,860px)] min-h-[420px] bg-background">
                  <iframe
                    src={frameSrc}
                    title={`${page.frameTitle}: ${title}`}
                    className="h-full w-full border-0"
                    allow="clipboard-write"
                  />
                </div>
              </div>
              <p className="text-xs text-muted-foreground">{page.note}</p>
            </div>
          </section>
        </div>
      </main>

      <SiteFooter content={landing} />
    </div>
  )
}

function SidebarLink({
  entry,
  locale,
  selected,
}: {
  entry: PreviewEntry
  locale: Locale
  selected: boolean
}) {
  return (
    <Link
      to={`?c=${entry.slug}`}
      preventScrollReset
      aria-current={selected ? "page" : undefined}
      className={cn(
        "block rounded-md px-2.5 py-1.5 text-sm transition-colors",
        selected
          ? "bg-primary/10 font-medium text-primary"
          : "text-muted-foreground hover:bg-muted hover:text-foreground"
      )}
    >
      {previewName(entry, locale)}
    </Link>
  )
}
