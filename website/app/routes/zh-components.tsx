import { ComponentBrowser } from "~/components/component-browser"
import { componentsPage } from "~/lib/content"

import type { Route } from "./+types/zh-components"

export const meta: Route.MetaFunction = () => [
  { title: componentsPage.zh.meta.title },
  { name: "description", content: componentsPage.zh.meta.description },
  { property: "og:title", content: componentsPage.zh.meta.title },
  { property: "og:description", content: componentsPage.zh.meta.description },
  { property: "og:type", content: "website" },
  { name: "twitter:card", content: "summary_large_image" },
]

export default function ZhComponents() {
  return <ComponentBrowser locale="zh" />
}
