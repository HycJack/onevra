import {
  ComponentBrowser,
  componentsMeta,
} from "~/components/component-browser"

import type { Route } from "./+types/zh-components"

export const meta: Route.MetaFunction = ({ params }) =>
  componentsMeta("zh", params.slug)

export default function ZhComponents() {
  return <ComponentBrowser locale="zh" />
}
