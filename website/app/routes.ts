import { type RouteConfig, index, route } from "@react-router/dev/routes"

export default [
  index("routes/home.tsx"),
  route("components/:slug?", "routes/components.tsx"),
  route("zh", "routes/zh.tsx"),
  route("zh/components/:slug?", "routes/zh-components.tsx"),
] satisfies RouteConfig
