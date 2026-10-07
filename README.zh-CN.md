# oneo-go

[English](README.md) | **简体中文**

基于 [MyGo](https://github.com/egoist/mygo) 的实验性通用 Go 组件库。目前提供原生 Shell 组件、适用于 MyGo 的 Lucide 图标和独立展示程序。

尚未发布为可下载的 module，API 暂不稳定。

`main` 使用未经定制补丁修改的上游 MyGo。WebView 合成实验因依赖定制版 MyGo，保留在 `feature/webview` 分支，不属于主线组件。

## 体验展示程序

需要 Go 1.27.1 或以上版本和原生桌面会话；Linux 还需要 GTK 3、Pango 和 OpenGL。

从仓库根目录运行：

```sh
cd source
go tool mygo dev
```

侧边栏提供 Shell、Lucide 图标库、oneo-go 控件和 MyGo 组件入口。**oneo-go controls** 分组为每个公共控件提供可交互的状态示例，包括全部按钮样式和尺寸。MyGo 组件按类型分组，各页面仅展示可交互的组件状态和用法变体，不包含代码或教程说明。Plugins / Terminal 页面提供 ANSI 输出和输入回显，不启动系统 Shell；详见[运行依赖](docs/build.md#app-bundles)。

在 macOS 上，这会启动支持自动重载的开发版应用。要预览指定外观，可运行 `go run ./examples/gallery -theme light` 或 `-theme dark`。

## 基本使用

发布前，请使用包含你的应用和本仓库 `source/` 目录的本地 Go workspace。导入 `github.com/oneo-me/oneo-go/ui`：

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    oneoui "github.com/oneo-me/oneo-go/ui"
)

func main() {
    mygo.App.WhenReady(func() {
        shell := oneoui.NewShell(oneoui.ShellOptions{
            Title: "Hello",
            Left: func(c *ui.Context) {
                ui.Text(c, "Navigation")
            },
            Content: func(c *ui.Context) {
                ui.Text(c, "Hello, Shell")
            },
        })
        shell.Window().OnClosed(func() { mygo.App.Quit() })
    })
    if err := mygo.App.Run(); err != nil {
        log.Fatal(err)
    }
}
```

在应用就绪后创建 Shell，使用 MyGo 控件构建各栏内容。使用 `shell.SetPage(view)` 替换整个侧边栏与主内容页面（例如显示设置），使用 `shell.SetPage(nil)` 恢复。由应用管理业务状态和生命周期。另可参考[可运行的展示程序](source/examples/gallery/main.go)。

## 外观

Shell 使用 oneo-go 内置的浅色与深色主题，通过 MyGo 自身的主题机制应用，因此所有 MyGo 控件都能正常工作。`ShellOptions.Appearance` 可选择 `theme.Auto`（默认，跟随系统外观）、`theme.Light` 或 `theme.Dark`：

```go
import "github.com/oneo-me/oneo-go/theme"

shell := oneoui.NewShell(oneoui.ShellOptions{
    Title:      "Hello",
    Appearance: theme.Dark,
    // ...
})
```

色板遵循 [oneo-me/desgin-styles](https://github.com/oneo-me/desgin-styles) 的 Light and Shadow（light-shadow）设计风格：暖白色浅色模式与深海军蓝深色模式。参见 [theme 包](source/theme/doc.go)与[外观设计](docs/architecture/components.md#native-appearance)。

## 控件

导入 `github.com/oneo-me/oneo-go/ui` 使用基于 MyGo 无样式基座的 oneo 风格常用控件：按钮、复选框、开关、单选按钮、文本框与文本域、下拉选择、滑块和选项卡。沿用上面的 `oneoui` 导入别名，与 MyGo 的 `ui` 区分。颜色取自当前主题：

```go
if oneoui.Button(c, "Save", oneoui.ButtonPrimary, oneoui.ButtonMedium).Clicked() {
    // ...
}

oneoui.Checkbox(c, &agreed, "I agree")
oneoui.Switch(c, &wifi, "Wi-Fi")
```

键盘、指针与无障碍行为由 MyGo 基座提供。参见[包契约](source/ui/doc.go)。

## Lucide 图标

导入 `github.com/oneo-me/oneo-go/icons/lucide`，直接将具名访问函数用于 MyGo：

```go
ui.Icon(c, lucide.Search()).Size(20, 20)
ui.SidebarItem(c, "settings", lucide.Settings(), "Settings")
```

图标首次使用时解析，之后复用。尺寸、文本颜色继承和无障碍语义由 MyGo 管理，无需额外的组件封装。图标包包含固定上游版本的正式图标，不包含 lab 图标或别名。参见[包契约](source/icons/lucide/doc.go)和[更新流程](docs/build.md#updating-lucide-icons)。

## 参与开发

修改代码前请阅读[架构概览](ARCHITECTURE.md)。在 `source/` 目录运行基本检查：

```sh
GOWORK=off go test ./...
GOWORK=off go vet ./...
```

打包与原生界面检查见[开发说明](docs/build.md)，面向用户的更新见[发布说明](CHANGELOG.md)。

## 协议

项目代码采用 [MIT](LICENSE) 协议。Lucide 资源保留[上游 ISC 和 Feather MIT 声明](source/icons/lucide/LICENSE)。
