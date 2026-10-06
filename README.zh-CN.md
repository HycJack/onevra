# oneo-go

[English](README.md) | **简体中文**

基于 MyGo v0.2.12 的实验性通用组合组件。包名为 `oneogo`，module 路径为 `github.com/oneo-me/oneo-go`。使用独立仓库开发，尚未发布为可下载的 module，不承诺公共 API 稳定性。

## 运行展示程序

需要 Go 1.27.1 或以上版本，从仓库根目录进入源码模块：

```sh
cd source
go tool mygo dev
go run ./examples/gallery -theme dark
go run ./examples/gallery -theme light
```

在 macOS 上，`go tool mygo dev` 以开发版 `.app` 运行并支持自动重载。`go tool mygo build` 在仓库根目录的 `dist/` 下生成本地 `.app` 和 `.dmg`。上面的 `go run` 命令用于直接运行可执行文件以检查主题，不生成应用包。参见[打包说明](docs/build.md#app-bundles)。

需要原生桌面会话；Linux 还需要 GTK 3、Pango 和 OpenGL。展示程序不写入应用设置或内容文件。

## 组件

- `Section(c, title, description, content)`：标题与可选内容区。
- `EmptyState(c, EmptyStateOptions{...})`：空状态说明与可选操作。
- `Notice(c, NoticeOptions{...})`：持续显示的信息、警告、错误或成功提示。

组件返回 `*ui.Element`，支持原生布局修饰和键值。回调在构建视图时同步执行。使用方持有状态、本地化文案和操作。重复组件应使用不同的键值；内部内容沿用 MyGo 的位置标识机制。

```go
oneogo.EmptyState(c, oneogo.EmptyStateOptions{
    Title: "No items",
    Description: "Create an item to get started.",
    Actions: func() {
        if ui.Button(c, "Create").Clicked() {
            // Update caller-owned state.
        }
    },
})
```

导入 `github.com/egoist/mygo/ui`，并将 `github.com/oneo-me/oneo-go` 别名设为 `oneogo`。所有组件使用当前 MyGo 主题，通过 MyGo 配置，不另建主题管理器。Notice 是带标签的分组，不是自动播报的实时警报；应提供明确的状态标题，避免仅通过颜色传达含义。

## Lucide 图标

导入 `github.com/oneo-me/oneo-go/icons/lucide`，直接将具名访问函数用于 MyGo：

```go
ui.Icon(c, lucide.Search()).Size(20, 20)
ui.SidebarItem(c, "settings", lucide.Settings(), "Settings")
```

图标首次使用时解析，之后复用。尺寸、文本颜色继承和无障碍语义由 MyGo 管理，无需额外的组件封装。图标包包含固定上游版本的正式图标，不包含 lab 图标或别名。参见[包契约](source/icons/lucide/doc.go)和[更新流程](docs/build.md#updating-lucide-icons)。

## 检查

在 `source/` 目录运行：

```sh
GOWORK=off go test ./...
GOWORK=off go test -race ./...
GOWORK=off go vet ./...
GOWORK=off go build -o ../dist/oneo-go-gallery ./examples/gallery
ONEO_GO_NATIVE_SMOKE=1 go test ./internal/gallery -run '^TestNativeGallery$' -count=1 -v
```

原生冒烟检查需显式启用，捕获浅色、深色和窄窗口后关闭窗口。设置 `ONEO_GO_SCREENSHOTS` 可保留 PNG。无头测试覆盖组合、操作、禁用控件、主题、文字缩放与窄布局。截图证明渲染可用，不代表完整键盘、屏幕阅读器或其他平台验收。

参见[架构](ARCHITECTURE.md)、[开发说明](docs/build.md)和[发布说明](CHANGELOG.md)。项目代码采用 [MIT](LICENSE) 协议。Lucide 资源保留[上游 ISC 和 Feather MIT 声明](source/icons/lucide/LICENSE)。
