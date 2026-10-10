# Onevra

[English](README.md) | **简体中文**

基于 [MyGo](https://github.com/egoist/mygo) 的实验性通用 Go 组件库。提供主题控件、原生窗口 Shell、命令面板、Markdown 与会话视图、Lucide 图标及独立展示程序。

首个计划版本为 **0.2026.1**，尚未发布为可下载的 module，API 暂不稳定。版本规则见[发布说明](docs/build.md#publication)。

`main` 使用未经定制补丁修改的上游 MyGo。WebView 合成实验因依赖定制版 MyGo，保留在 `experiment/webview` 分支，不属于主线组件。

## 名字由来

**Onevra** 是本项目的自造词名称：保留 **Oneo** 开头的 `One`，再组合 `vra`，形成简短且有辨识度的读音与拼写。词尾是音形上的选择，不是缩写，也没有赋予它特定的词根含义。

命名时优先考虑与 Oneo 的联系和较低的重名概率，而不是采用常见英文单词；同时避免将组件库的身份绑定到 Go、MyGo 或某一种视觉风格。初步 GitHub 仓库名称检索未发现 `onevra` 的匹配结果，但这不代表全网无人使用，也不等于商标可用性确认。

## 体验展示程序

![Onevra Gallery 主界面，包含侧边导航和可交互组件示例](docs/images/gallery.png)

需要 Go 1.27.1 或以上版本和原生桌面会话；Linux 还需要 GTK 3、Pango 和 OpenGL。

从仓库根目录运行：

```sh
cd source
go run ./examples/gallery
```

展示程序包含原生 Shell、Lucide 图标库、Onevra 控件与插件，以及 MyGo 示例。组件页面按用途分组，提供可交互的状态示例；名称末尾的圆点标记 MyGo 内置组件。Plugins / Terminal 页面在首次进入时启动默认登录 Shell；详见[运行依赖](docs/build.md#app-bundles)。

要预览指定外观，可追加 `-theme light` 或 `-theme dark`。`go run ./tools/packer run` 编译并启动打包后的开发版应用，不进行热重载；`go run ./tools/packer build` 根据 `source/packer.yaml` 中选择的目标平台生成发布归档。详见[打包与图标要求](docs/build.md#app-bundles)。

在 **Desktop / Tray** 页面启用托盘图标后，macOS 和 Windows 上可左键切换便笺小窗口、右键打开原生菜单。Linux 改为通过菜单项打开窗口，需要 AppIndicator。参见[托盘示例依赖与检查](docs/build.md#tray-example)。macOS 上还可在 **Desktop / Dock menu** 中为 Dock 图标的右键菜单添加 Gallery 快捷操作，参见[Dock 菜单检查](docs/build.md#dock-menu-example)。

在 macOS 上按 Command+E，其他平台按 Ctrl+E，即可搜索 Gallery 页面和操作。应用可以复用[命令面板与弹窗](source/plugins/command/example_test.go)，提供自己的命令和执行逻辑。

macOS 提供原生检查流程；Windows 和 Linux 尚需各自的桌面验证。编译和无窗口测试通过不代表已完成键盘、输入法或无障碍验证，详见[原生检查](docs/build.md#native-review)。

## 基本使用

模块路径为 `github.com/oneo-me/onevra/source`。**以下安装命令需要等首个版本发布后才能使用**：

```sh
go get github.com/oneo-me/onevra/source@v0.2026.1
```

发布前，请使用包含你的应用和本仓库 `source/` 目录的本地 Go workspace；工作区配置不要提交到仓库。导入 `github.com/oneo-me/onevra/source/shells/sidebar`：

```go
package main

import (
    "log"

    "github.com/egoist/mygo"
    "github.com/egoist/mygo/ui"
    "github.com/oneo-me/onevra/source/components"
    "github.com/oneo-me/onevra/source/shells/sidebar"
)

func main() {
    mygo.App.WhenReady(func() {
        appShell := sidebar.NewShell(sidebar.ShellOptions{
            Title: "Hello",
            Left: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() {
                    components.Text(c, "Navigation")
                })
            },
            Content: func(c *ui.Context) {
                ui.Column(c).Fill().Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() {
                    components.Text(c, "Hello, Shell")
                })
            },
        })
        appShell.Window().OnClosed(func() { mygo.App.Quit() })
    })
    if err := mygo.App.Run(); err != nil {
        log.Fatal(err)
    }
}
```

在应用就绪后创建 Shell，使用 Onevra 控件和 MyGo 布局原语构建各栏内容。使用 `appShell.SetPage(view)` 替换整个侧边栏与主内容页面（例如显示设置），使用 `appShell.SetPage(nil)` 恢复。由应用管理业务状态和生命周期。另可参考[可运行的展示程序](source/examples/gallery/main.go)。

内容默认延伸到 Shell 的悬浮标题栏下方。需要避让时，使用 `sidebar.HeaderHeight(c)` 设置顶部内边距；在这个容器上绘制背景，即可覆盖标题栏区域而不遮挡内容控件。替换页也遵循这一规则。

## 外观

Shell 使用 Onevra 内置的浅色与深色主题，通过 MyGo 自身的主题机制应用，因此所有 MyGo 控件都能正常工作。`ShellOptions.Appearance` 可选择 `theme.Auto`（默认，跟随系统外观）、`theme.Light` 或 `theme.Dark`：

```go
import "github.com/oneo-me/onevra/source/theme"

appShell := sidebar.NewShell(sidebar.ShellOptions{
    Title:      "Hello",
    Appearance: theme.Dark,
    // ...
})
```

Onevra 设计风格以暖白色浅色界面与暖炭灰深色界面搭配同一色系的陶土橙强调色。参见 [theme 包](source/theme/doc.go)与[外观设计](docs/architecture/components.md#native-appearance)。

## 控件

导入 `github.com/oneo-me/onevra/source/components` 使用基于 MyGo 无样式基座的 Onevra 风格常用控件：按钮、复选框、开关、单选按钮、文本框与文本域、下拉选择、滑块和选项卡。包名为 `components`，可以与 MyGo 的 `ui` 同时导入，无需别名。颜色取自当前主题：

```go
if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
    // ...
}

components.Checkbox(c, &agreed, "I agree")
components.Switch(c, &wifi, "Wi-Fi")
```

键盘、指针与无障碍行为由 MyGo 基座提供。参见[包契约](source/components/doc.go)。

## Lucide 图标

导入 `github.com/oneo-me/onevra/source/components` 和 `github.com/oneo-me/onevra/source/icons/lucide`，将具名访问函数传给 Onevra 的 Icon 组件：

```go
components.Icon(c, lucide.Search).Size(20, 20)
components.Icon(c, lucide.Settings).Size(20, 20)
```

图标首次使用时解析，之后复用。Icon 也接受其他图标包提供的已解析 `*ui.SVG` 资源。尺寸、文本颜色继承和无障碍语义由 MyGo 管理。图标包包含固定上游版本的正式图标，不包含 lab 图标或别名。参见[包契约](source/icons/lucide/doc.go)和[更新流程](docs/build.md#updating-lucide-icons)。

## 参与开发

需要小型应用模板时，可参考支持多语言的[普通、菜单栏和混合应用示例](docs/build.md#starter-applications)，每个示例都有自己的启动行为、偏好设置与打包配置。

修改代码前请阅读[架构概览](ARCHITECTURE.md)。在 `source/` 目录运行基本检查：

```sh
GOWORK=off go test ./...
GOWORK=off go vet ./...
```

打包与原生界面检查见[开发说明](docs/build.md)，面向用户的更新见[发布说明](CHANGELOG.md)。

## 协议

项目代码采用 [MIT](LICENSE) 协议。改编自 MyGo 的控件保留[上游 MIT 声明](source/components/LICENSE.mygo)。Lucide 资源保留[上游 ISC 和 Feather MIT 声明](source/icons/lucide/LICENSE)。
