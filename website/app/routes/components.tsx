import { ComponentBrowser } from "~/components/component-browser"
import { componentsPage } from "~/lib/content"

import type { Route } from "./+types/components"

export const meta: Route.MetaFunction = () => [
  { title: componentsPage.en.meta.title },
  { name: "description", content: componentsPage.en.meta.description },
  { property: "og:title", content: componentsPage.en.meta.title },
  { property: "og:description", content: componentsPage.en.meta.description },
  { property: "og:type", content: "website" },
  { name: "twitter:card", content: "summary_large_image" },
]

export default function Components() {
  return <ComponentBrowser locale="en" />
}
