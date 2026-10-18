import { type RouteConfig, index, route } from "@react-router/dev/routes"

export default [
  index("routes/home.tsx"),
  route("components", "routes/components.tsx"),
  route("zh", "routes/zh.tsx"),
  route("zh/components", "routes/zh-components.tsx"),
] satisfies RouteConfig
