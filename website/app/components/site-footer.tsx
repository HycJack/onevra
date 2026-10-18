import { Link, useLocation } from "react-router"

import { Logo } from "~/components/logo"
import type { LandingContent } from "~/lib/content"
import { localeFromPathname, site } from "~/lib/site"

export function SiteFooter({ content }: { content: LandingContent }) {
  const locale = localeFromPathname(useLocation().pathname)
  const homeHref = locale === "zh" ? "/zh" : "/"
  const linkHref = (href: string) => {
    if (href.startsWith("#")) return `${homeHref}${href}`
    if (href.startsWith("/")) return locale === "zh" ? `/zh${href}` : href
    return href
  }

  return (
    <footer className="border-t border-border/60 bg-muted/20">
      <div className="mx-auto grid w-full max-w-6xl gap-10 px-4 py-14 sm:px-6 md:grid-cols-[2fr_repeat(3,1fr)]">
        <div className="flex flex-col gap-3">
          <Link
            to={homeHref}
            className="flex items-center gap-2 font-heading text-base font-semibold"
          >
            <Logo />
            {site.name}
          </Link>
          <p className="max-w-xs text-sm text-muted-foreground">
            {content.footer.description}
          </p>
        </div>

        {content.footer.groups.map((group) => (
          <div key={group.title} className="flex flex-col gap-3">
            <h2 className="font-heading text-sm font-medium">{group.title}</h2>
            <ul className="flex flex-col gap-2">
              {group.links.map((link) => (
                <li key={link.label}>
                  <a
                    href={linkHref(link.href)}
                    target={link.href.startsWith("http") ? "_blank" : undefined}
                    rel={
                      link.href.startsWith("http") ? "noreferrer" : undefined
                    }
                    className="text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      <div className="border-t border-border/60">
        <div className="mx-auto w-full max-w-6xl px-4 py-6 text-xs text-muted-foreground sm:px-6">
          {content.footer.builtWith}
        </div>
      </div>
    </footer>
  )
}
