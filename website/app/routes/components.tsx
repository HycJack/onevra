import {
  ComponentBrowser,
  componentsMeta,
} from "~/components/component-browser"

import type { Route } from "./+types/components"

export const meta: Route.MetaFunction = ({ params }) =>
  componentsMeta("en", params.slug)

export default function Components() {
  return <ComponentBrowser locale="en" />
}
