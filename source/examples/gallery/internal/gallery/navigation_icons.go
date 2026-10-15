package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// Keys are stable example names, not translated navigation labels. Keep SVGs lazy.
var navigationIcons = map[string]func() *ui.SVG{
	"Bar chart":          lucide.ChartColumn,
	"Line chart":         lucide.ChartLine,
	"Area chart":         lucide.ChartArea,
	"Stacked bar chart":  lucide.ChartColumnStacked,
	"Donut chart":        lucide.ChartPie,
	"Relationship graph": lucide.Waypoints,

	"Transforms":         lucide.Rotate3d,
	"Motion":             lucide.Activity,
	"TimePicker":         lucide.Clock,
	"Calendar":           lucide.CalendarDays,
	"CalendarPicker":     lucide.CalendarCheck,
	"ColorView":          lucide.SwatchBook,
	"ColorPicker":        lucide.Palette,
	"Tooltip":            lucide.MessageCircle,
	"Popover":            lucide.PanelTopDashed,
	"Dialog":             lucide.AppWindow,
	"AlertDialog":        lucide.TriangleAlert,
	"Badge":              lucide.Badge,
	"Chat":               lucide.MessagesSquare,
	"Markdown":           lucide.FileText,
	"RichText":           lucide.Bold,
	"Breadcrumbs":        lucide.ChevronsRight,
	"Rating":             lucide.Star,
	"Text":               lucide.Type,
	"TextBox":            lucide.TextCursorInput,
	"TextBoxBuffer":      lucide.TextCursor,
	"Button":             lucide.MousePointerClick,
	"Selection":          lucide.SquareCheck,
	"Select":             lucide.ListFilter,
	"Menu":               lucide.Menu,
	"NumberInput":        lucide.Hash,
	"Slider":             lucide.SlidersHorizontal,
	"Progress":           lucide.RectangleEllipsis,
	"Spinner":            lucide.LoaderCircle,
	"Tabs":               lucide.PanelsTopLeft,
	"Accordion":          lucide.ListChevronsDownUp,
	"Collapsible":        lucide.ChevronsUpDown,
	"Icon":               lucide.Shapes,
	"Image":              lucide.Image,
	"List":               lucide.List,
	"GridView":           lucide.LayoutGrid,
	"Tree":               lucide.ListTree,
	"Sidebar":            lucide.PanelLeft,
	"PixelGrid":          lucide.Grid3x3,
	"Honeycomb":          lucide.Hexagon,
	"Split":              lucide.Columns2,
	"ScrollView":         lucide.ArrowDownUp,
	"CheckboxGroup":      lucide.ListChecks,
	"RadioGroup":         lucide.CircleDot,
	"StepSlider":         lucide.SlidersHorizontal,
	"RangeSlider":        lucide.SlidersHorizontal,
	"Stepper":            lucide.ListOrdered,
	"TextInputBuffer":    lucide.TextCursorInput,
	"TextAreaBuffer":     lucide.NotepadText,
	"TokenField":         lucide.Tags,
	"EditableText":       lucide.PencilLine,
	"DateInput":          lucide.CalendarClock,
	"TimeInput":          lucide.Clock,
	"ColorWell":          lucide.Pipette,
	"Table":              lucide.Table2,
	"Outline":            lucide.ListTree,
	"OutlineTable":       lucide.TableProperties,
	"Scroll":             lucide.ArrowDownUp,
	"ScrollHorizontal":   lucide.ArrowLeftRight,
	"ScrollBoth":         lucide.Move,
	"Grid":               lucide.Grid2x2,
	"SplitVertical":      lucide.Rows2,
	"Modal":              lucide.AppWindow,
	"Toast":              lucide.MessageSquare,
	"Box":                lucide.Square,
	"Row":                lucide.Columns3,
	"Column":             lucide.Rows3,
	"Spacer":             lucide.AlignHorizontalSpaceAround,
	"Divider":            lucide.SeparatorHorizontal,
	"Terminal":           lucide.Terminal,
	"Flyout":             lucide.PanelsTopLeft,
	"WebView":            lucide.Globe,
	"Glass":              lucide.GlassWater,
	"Audio":              lucide.AudioLines,
	"System information": lucide.Info,
	"Recent documents":   lucide.FileClock,
	"File watching":      lucide.FolderSync,
	"Canvas":             lucide.Brush,
	"Code editor":        lucide.SquareCode,
	"Rich text editor":   lucide.NotebookPen,
	"Markdown stream":    lucide.MessageSquareText,
	"Diff viewer":        lucide.Diff,
	"Game loop":          lucide.Gamepad2,
	"3D renderer":        lucide.Cuboid,
}

func (entry navigationEntry) icon() *ui.SVG {
	switch entry.page {
	case welcomePage:
		return lucide.House()
	case lucidePage:
		return lucide.Shapes()
	case trayPage:
		return lucide.Inbox()
	case dockPage:
		return lucide.PanelBottom()
	case notificationsPage:
		return lucide.Bell()
	}
	name := ""
	if example, ok := controlExampleFor(entry.page); ok {
		name = example.name
	} else if example, ok := mygoExampleFor(entry.page); ok {
		name = example.name
	}
	if icon := navigationIcons[name]; icon != nil {
		return icon()
	}
	return lucide.Component()
}
