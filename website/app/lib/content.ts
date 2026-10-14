import { site, type Locale } from "~/lib/site"

export interface NavItem {
  label: string
  href: string
}

export interface CodeSample {
  name: string
  code: string
}

export interface FaqItem {
  question: string
  answer: string
}

export interface ComponentGroup {
  title: string
  items: string[]
}

export interface ViewItem {
  name: string
  exports: string[]
  description: string
}

export interface FooterGroup {
  title: string
  links: NavItem[]
}

export interface LandingContent {
  meta: { title: string; description: string }
  nav: NavItem[]
  header: { language: string; theme: string; github: string }
  hero: {
    badge: string
    title: string
    accent: string
    highlight: string
    description: string
    primary: string
    secondary: string
    hint: string
  }
  catalog: {
    title: string
    description: string
    groups: ComponentGroup[]
    note: string
  }
  views: { title: string; description: string; items: ViewItem[] }
  code: {
    title: string
    description: string
    copy: string
    copied: string
    samples: CodeSample[]
  }
  theming: {
    title: string
    description: string
    bullets: { title: string; description: string }[]
  }
  install: {
    title: string
    description: string
    steps: { label: string; code: string }[]
    requirements: string
  }
  faq: { title: string; description: string; items: FaqItem[] }
  cta: {
    title: string
    description: string
    primary: string
    secondary: string
  }
  footer: { description: string; groups: FooterGroup[]; builtWith: string }
}

const samples: CodeSample[] = [
  {
    name: "main.go",
    code: `package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

func main() {
	mygo.App.WhenReady(func() {
		shell := plain.NewShell(plain.ShellOptions{
			Title: "Hello Onevra",
			Theme: theme.Onevra,
			Content: func(c *ui.Context) {
				ui.Column(c).Fill().Padding(24).Children(func() {
					components.Text(c, "Hello, Onevra!")
				})
			},
		})
		shell.Window().OnClosed(func() { mygo.App.Quit() })
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}`,
  },
  {
    name: "window.go",
    code: `package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/sidebar"
	"github.com/oneo-me/onevra/source/theme"
)

func newWindow() *sidebar.Shell {
	return sidebar.NewShell(sidebar.ShellOptions{
		Title:    "Onevra",
		StateKey: "main",
		Theme:    theme.Onevra,
		Left:     navigation,
		Content:  libraryPage,
	})
}

func libraryPage(c *ui.Context) {
	// Build application content with Onevra controls.
}`,
  },
  {
    name: "controls.go",
    code: `package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

type settings struct {
	agreed bool
	wifi   bool
}

func settingsPage(c *ui.Context, s *settings) {
	ui.Column(c).Fill().Padding(24).Gap(12).Children(func() {
		if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
			// Save application state.
		}
		components.Checkbox(c, &s.agreed, "I agree")
		components.Switch(c, &s.wifi, "Wi-Fi")
		components.Icon(c, lucide.Search).Size(20, 20)
	})
}`,
  },
  {
    name: "markdown.go",
    code: `package main

import (
	"github.com/egoist/mygo/ui"
	markdown "github.com/oneo-me/onevra/source/components/markdown"
)

func documentPage(c *ui.Context) {
	source := "# Onevra\\n\\nA component library for native desktop apps.\\n"
	markdown.Markdown(c, source)
}`,
  },
]

export const content: Record<Locale, LandingContent> = {
  en: {
    meta: {
      title: "Onevra — the native desktop component library",
      description:
        "Onevra is the MyGo component library: buttons, inputs, selection, navigation, collections, overlays, charts, graphs, a command palette, code and Markdown rendering, conversations and window shells.",
    },
    nav: [
      { label: "Components", href: "#components" },
      { label: "Views", href: "#views" },
      { label: "Code", href: "#code" },
      { label: "Theming", href: "#theming" },
      { label: "Install", href: "#install" },
      { label: "FAQ", href: "#faq" },
    ],
    header: {
      language: "Switch language",
      theme: "Toggle theme",
      github: "View on GitHub",
    },
    hero: {
      badge: "MyGo · Native desktop · Light and dark",
      title: "The native desktop component library",
      accent: "Onevra",
      highlight: "one focused toolkit",
      description:
        "Onevra is the MyGo component library: themed controls, window shells, Lucide icons and complex views for chat, Markdown, code, charts and graphs — one focused toolkit you can adopt piece by piece.",
      primary: "Browse the components",
      secondary: "View on GitHub",
      hint: "go get github.com/oneo-me/onevra/source",
    },
    catalog: {
      title: "Every component already in the library",
      description:
        "Onevra ships the controls below today, each with package documentation and tests. Import them from github.com/oneo-me/onevra/source/components.",
      groups: [
        {
          title: "Actions",
          items: ["Button", "ToggleButton", "Badge", "Icon"],
        },
        {
          title: "Text and input",
          items: [
            "Text",
            "RichText",
            "TextBox",
            "TextBoxBuffer",
            "NumberInput",
          ],
        },
        {
          title: "Selection",
          items: [
            "Checkbox",
            "Switch",
            "Radio",
            "Rating",
            "Select",
            "NativeSelect",
          ],
        },
        {
          title: "Values and pickers",
          items: [
            "Slider",
            "Calendar",
            "CalendarPicker",
            "TimePicker",
            "ColorPicker",
            "ColorView",
          ],
        },
        {
          title: "Navigation and layout",
          items: [
            "Tabs",
            "Accordion",
            "Collapsible",
            "Breadcrumbs",
            "Sidebar",
            "Split",
            "ScrollView",
          ],
        },
        {
          title: "Collections",
          items: [
            "List",
            "GridView",
            "Tree",
            "TreeVirtualized",
            "Table",
            "TableVirtualized",
          ],
        },
        {
          title: "Overlays and feedback",
          items: [
            "Dialog",
            "AlertDialog",
            "Popover",
            "Tooltip",
            "Toasts",
            "Progress",
            "ProgressRing",
            "Spinner",
          ],
        },
        {
          title: "Surfaces and media",
          items: ["Settings", "SearchSettings", "Image", "PixelGrid"],
        },
        {
          title: "Shells and icons",
          items: ["plain", "sidebar", "Lucide icons"],
        },
      ],
      note: "Complex views live in their own packages, listed in the next section.",
    },
    views: {
      title: "Complex views, in their own packages",
      description:
        "Beyond single controls, Onevra implements the complete views an application needs. Each one is a separate package with its own documentation and tests.",
      items: [
        {
          name: "components/chart",
          exports: ["Bar", "Line", "Area", "StackedBar", "Donut"],
          description:
            "Categorical charts with smooth updates, interactive legends, keyboard selection and custom empty and detail content.",
        },
        {
          name: "components/graph",
          exports: ["View"],
          description:
            "An interactive force-directed relationship graph with pan, zoom, hover and draggable nodes.",
        },
        {
          name: "components/command",
          exports: ["CommandPalette", "CommandDialog"],
          description:
            "A searchable, keyboard-driven command palette, embeddable or hosted in a dialog.",
        },
        {
          name: "components/code",
          exports: ["Parse", "Document"],
          description:
            "Read-only syntax highlighting with light and dark palettes, usable alone or inside Markdown.",
        },
        {
          name: "components/markdown",
          exports: ["Markdown", "ParseMarkdown"],
          description:
            "Selectable CommonMark and GFM documents with application-supplied images and custom code-block views.",
        },
        {
          name: "components/chat",
          exports: ["Conversation", "Composer", "Turn", "ActivityRow"],
          description:
            "Conversations with turns, activity groups, folded summaries, attachments and a composer.",
        },
        {
          name: "components/gesture",
          exports: ["Classify", "Navigation"],
          description:
            "Wheel, pinch and touchpad pan-and-zoom handling for custom canvases.",
        },
      ],
    },
    code: {
      title: "Less ceremony, more product",
      description:
        "Create a window, build content with Onevra controls, and keep state in your own values.",
      copy: "Copy",
      copied: "Copied",
      samples,
    },
    theming: {
      title: "A theme that follows the desktop",
      description:
        "Onevra's palettes cover the controls, the shells and the window chrome. Switch appearance once and the whole application follows.",
      bullets: [
        {
          title: "Warm neutrals, terracotta accent",
          description:
            "Light and dark palettes share one accent family, with dark mode lifting the accent for legibility.",
        },
        {
          title: "System-aware",
          description:
            "The preset follows the desktop's light or dark appearance, text scale and higher-contrast settings.",
        },
        {
          title: "One switch, two moods",
          description:
            "Choose light, dark or system and every component and window edge changes with it.",
        },
      ],
    },
    install: {
      title: "Up and running in a minute",
      description:
        "Onevra requires Go 1.27.1 or later. Add the module, then run the gallery.",
      steps: [
        {
          label: "Add the library",
          code: "go get github.com/oneo-me/onevra/source@latest",
        },
        {
          label: "Run the gallery",
          code: "cd source && go run ./examples/gallery",
        },
        {
          label: "Try a starter application",
          code: "cd source && go run ./examples/desktop",
        },
      ],
      requirements:
        "Each starter application includes localization, preferences and packaging configuration.",
    },
    faq: {
      title: "Frequently asked questions",
      description: "Short answers before you open the docs.",
      items: [
        {
          question: "What does Onevra include?",
          answer:
            "The controls listed above, two window shells, Lucide icons, an optional light and dark palette, and seven complex-view packages for charts, graphs, commands, code, Markdown, chat and gestures.",
        },
        {
          question: "Which platforms are supported?",
          answer:
            "macOS, Windows and Linux. macOS has the most native review; Windows and Linux still need their own desktop review before you rely on them.",
        },
        {
          question: "Is the API stable?",
          answer:
            "Not yet. Onevra is experimental and its public API can change. Releases use the 0.YEAR.SEQUENCE calendar convention.",
        },
        {
          question: "Can I adopt only part of the library?",
          answer:
            "Yes. The theme, window shells, common controls and icons are independent, so take only what the application needs.",
        },
        {
          question: "Do the components ship their own text?",
          answer:
            "Components read application-supplied labels and can supply their own interface text within a localization scope, so your keys, stored values and translations stay yours.",
        },
        {
          question: "Is Onevra free?",
          answer:
            "Yes. Onevra is released under the MIT license, and assets adapted from upstream projects keep their notices.",
        },
      ],
    },
    cta: {
      title: "Ready to build with Onevra?",
      description:
        "Explore the components, run the gallery, and shape the library with us.",
      primary: "Read the architecture",
      secondary: "Browse the source",
    },
    footer: {
      description:
        "Onevra is a component library for native desktop applications: controls, window shells, icons and complex views.",
      groups: [
        {
          title: "Product",
          links: [
            { label: "Components", href: "#components" },
            { label: "Views", href: "#views" },
            { label: "Install", href: "#install" },
            { label: "Gallery", href: site.gallery },
          ],
        },
        {
          title: "Resources",
          links: [
            { label: "Architecture", href: site.architecture },
            { label: "Examples", href: site.examples },
            { label: "Changelog", href: site.changelog },
            { label: "License", href: site.license },
          ],
        },
        {
          title: "Community",
          links: [
            { label: "GitHub", href: site.github },
            { label: "Issues", href: site.issues },
          ],
        },
      ],
      builtWith:
        "Built with React Router and shadcn/ui. Released under the MIT license.",
    },
  },
  zh: {
    meta: {
      title: "Onevra — 原生桌面组件库",
      description:
        "Onevra 是 MyGo 组件库：按钮、输入、选择、导航、集合、浮层，以及图表、关系图、命令面板、代码与 Markdown 渲染、对话和窗口外壳。",
    },
    nav: [
      { label: "组件", href: "#components" },
      { label: "复杂视图", href: "#views" },
      { label: "代码", href: "#code" },
      { label: "主题", href: "#theming" },
      { label: "安装", href: "#install" },
      { label: "常见问题", href: "#faq" },
    ],
    header: {
      language: "切换语言",
      theme: "切换主题",
      github: "在 GitHub 上查看",
    },
    hero: {
      badge: "MyGo · 原生桌面 · 明暗主题",
      title: "原生桌面组件库",
      accent: "Onevra",
      highlight: "一套专注的工具包",
      description:
        "Onevra 是 MyGo 的组件库：主题化控件、窗口外壳、Lucide 图标，以及聊天、Markdown、代码、图表与关系图等复杂视图——一套专注的工具包，可按需取用。",
      primary: "浏览组件",
      secondary: "在 GitHub 上查看",
      hint: "go get github.com/oneo-me/onevra/source",
    },
    catalog: {
      title: "组件库中已经实现的组件",
      description:
        "以下控件现在都已实现，每个都附带包文档与测试，统一从 github.com/oneo-me/onevra/source/components 引入。",
      groups: [
        {
          title: "操作",
          items: ["Button", "ToggleButton", "Badge", "Icon"],
        },
        {
          title: "文本与输入",
          items: [
            "Text",
            "RichText",
            "TextBox",
            "TextBoxBuffer",
            "NumberInput",
          ],
        },
        {
          title: "选择",
          items: [
            "Checkbox",
            "Switch",
            "Radio",
            "Rating",
            "Select",
            "NativeSelect",
          ],
        },
        {
          title: "数值与选择器",
          items: [
            "Slider",
            "Calendar",
            "CalendarPicker",
            "TimePicker",
            "ColorPicker",
            "ColorView",
          ],
        },
        {
          title: "导航与布局",
          items: [
            "Tabs",
            "Accordion",
            "Collapsible",
            "Breadcrumbs",
            "Sidebar",
            "Split",
            "ScrollView",
          ],
        },
        {
          title: "集合",
          items: [
            "List",
            "GridView",
            "Tree",
            "TreeVirtualized",
            "Table",
            "TableVirtualized",
          ],
        },
        {
          title: "浮层与反馈",
          items: [
            "Dialog",
            "AlertDialog",
            "Popover",
            "Tooltip",
            "Toasts",
            "Progress",
            "ProgressRing",
            "Spinner",
          ],
        },
        {
          title: "界面与媒体",
          items: ["Settings", "SearchSettings", "Image", "PixelGrid"],
        },
        {
          title: "窗口外壳与图标",
          items: ["plain", "sidebar", "Lucide icons"],
        },
      ],
      note: "复杂视图位于各自独立的包中，见下一节。",
    },
    views: {
      title: "复杂视图，各自独立成包",
      description:
        "除了单个控件，Onevra 还实现了应用真正需要的完整视图。每个视图都是独立包，自带文档与测试。",
      items: [
        {
          name: "components/chart",
          exports: ["Bar", "Line", "Area", "StackedBar", "Donut"],
          description:
            "分类图表支持平滑更新、交互图例、键盘选中，以及自定义空状态与详情内容。",
        },
        {
          name: "components/graph",
          exports: ["View"],
          description: "可交互的力导向关系图，支持平移、缩放、悬停与拖动节点。",
        },
        {
          name: "components/command",
          exports: ["CommandPalette", "CommandDialog"],
          description: "可搜索、键盘驱动的命令面板，可嵌入使用或置于对话框中。",
        },
        {
          name: "components/code",
          exports: ["Parse", "Document"],
          description:
            "只读语法高亮，提供明暗两套配色，可单独使用或嵌入 Markdown。",
        },
        {
          name: "components/markdown",
          exports: ["Markdown", "ParseMarkdown"],
          description:
            "可选择的 CommonMark 与 GFM 文档，支持应用提供的图片与自定义代码块视图。",
        },
        {
          name: "components/chat",
          exports: ["Conversation", "Composer", "Turn", "ActivityRow"],
          description:
            "对话包含消息轮次、工具活动分组、折叠摘要、附件与输入区。",
        },
        {
          name: "components/gesture",
          exports: ["Classify", "Navigation"],
          description: "面向自定义画布的滚轮、捏合与触控板平移缩放。",
        },
      ],
    },
    code: {
      title: "更少样板，更多产品",
      description:
        "创建窗口，用 Onevra 控件构建内容，状态始终由你自己的值持有。",
      copy: "复制",
      copied: "已复制",
      samples,
    },
    theming: {
      title: "跟随桌面的主题",
      description:
        "Onevra 的调色板覆盖控件、外壳与窗口边框。切换一次外观，整个应用随之改变。",
      bullets: [
        {
          title: "暖色中性色，陶土强调色",
          description:
            "明暗两套调色板共享同一强调色系，深色模式会提亮强调色以保证可读性。",
        },
        {
          title: "感知系统设置",
          description: "预设跟随桌面的明暗外观、文字缩放与高对比度设置。",
        },
        {
          title: "一次切换，两种氛围",
          description:
            "选择浅色、深色或跟随系统，所有组件与原生窗口边框同步切换。",
        },
      ],
    },
    install: {
      title: "一分钟即可上手",
      description:
        "Onevra 需要 Go 1.27.1 或更高版本。添加模块后即可运行 Gallery。",
      steps: [
        {
          label: "添加依赖",
          code: "go get github.com/oneo-me/onevra/source@latest",
        },
        {
          label: "运行 Gallery",
          code: "cd source && go run ./examples/gallery",
        },
        {
          label: "体验起始应用",
          code: "cd source && go run ./examples/desktop",
        },
      ],
      requirements: "每个起始应用都自带本地化、偏好设置与打包配置。",
    },
    faq: {
      title: "常见问题",
      description: "在打开文档之前，先看几个简短的回答。",
      items: [
        {
          question: "Onevra 包含哪些内容？",
          answer:
            "上文列出的全部控件、两种窗口外壳、Lucide 图标、可选的明暗调色板，以及图表、关系图、命令、代码、Markdown、聊天与手势七个复杂视图包。",
        },
        {
          question: "支持哪些平台？",
          answer:
            "macOS、Windows 与 Linux。macOS 的原生验证最充分；Windows 与 Linux 仍需各自的桌面复核。",
        },
        {
          question: "API 稳定吗？",
          answer:
            "尚未稳定。Onevra 仍是实验性项目，公开 API 可能变化；版本号采用 0.年份.序号 的日历约定。",
        },
        {
          question: "可以只使用其中一部分吗？",
          answer:
            "可以。主题、窗口外壳、通用控件与图标彼此独立，按应用需要取用即可。",
        },
        {
          question: "组件会自带界面文案吗？",
          answer:
            "组件读取应用提供的标签，并可在本地化作用域内提供自身界面文案，键值、存储值与翻译都由应用掌握。",
        },
        {
          question: "Onevra 是免费的吗？",
          answer:
            "是的。Onevra 基于 MIT 许可证发布；改编自上游项目的资源保留其声明。",
        },
      ],
    },
    cta: {
      title: "准备好用 Onevra 构建了吗？",
      description: "浏览组件，运行 Gallery，与我们一起打磨这个组件库。",
      primary: "阅读架构文档",
      secondary: "浏览源代码",
    },
    footer: {
      description:
        "Onevra 是面向原生桌面应用的组件库：控件、窗口外壳、图标与复杂视图。",
      groups: [
        {
          title: "产品",
          links: [
            { label: "组件", href: "#components" },
            { label: "复杂视图", href: "#views" },
            { label: "安装", href: "#install" },
            { label: "Gallery", href: site.gallery },
          ],
        },
        {
          title: "资源",
          links: [
            { label: "架构", href: site.architecture },
            { label: "示例", href: site.examples },
            { label: "更新日志", href: site.changelog },
            { label: "许可证", href: site.license },
          ],
        },
        {
          title: "社区",
          links: [
            { label: "GitHub", href: site.github },
            { label: "问题反馈", href: site.issues },
          ],
        },
      ],
      builtWith: "由 React Router 与 shadcn/ui 构建，基于 MIT 许可证发布。",
    },
  },
}
