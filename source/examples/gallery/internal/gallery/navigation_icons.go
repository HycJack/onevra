package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// Keys are stable example names, not translated navigation labels. Keep SVGs lazy.
var navigationIcons = map[string]func() *ui.SVG{
	"Shader backgrounds": lucide.Orbit,
	"Data charts":        lucide.ChartColumn,
	"Relationship graph": lucide.Waypoints,

	"Transforms":          lucide.Rotate3d,
	"Motion":              lucide.Activity,
	"Date & time pickers": lucide.CalendarClock,
	"Calendar":            lucide.CalendarDays,
	"Color pickers":       lucide.SwatchBook,
	"ColorPicker":         lucide.Palette,
	"Tooltip":             lucide.MessageCircle,
	"Popover":             lucide.PanelTopDashed,
	"Dialogs":             lucide.AppWindow,
	"AlertDialog":         lucide.TriangleAlert,
	"Badge":               lucide.Badge,
	"Chat":                lucide.MessagesSquare,
	"Markdown":            lucide.FileText,
	"RichText":            lucide.Bold,
	"Breadcrumbs":         lucide.ChevronsRight,
	"Rating":              lucide.Star,
	"Text":                lucide.Type,
	"TextBox":             lucide.TextCursorInput,
	"Button":              lucide.MousePointerClick,
	"Selection":           lucide.SquareCheck,
	"Select":              lucide.ListFilter,
	"Menu":                lucide.Menu,
	"NumberInput":         lucide.Hash,
	"Slider":              lucide.SlidersHorizontal,
	"Progress":            lucide.RectangleEllipsis,
	"Spinner":             lucide.LoaderCircle,
	"Tabs":                lucide.PanelsTopLeft,
	"Accordion":           lucide.ListChevronsDownUp,
	"Collapsible":         lucide.ChevronsUpDown,
	"Disclosure":          lucide.ChevronsUpDown,
	"Image":               lucide.Image,
	"List":                lucide.List,
	"GridView":            lucide.LayoutGrid,
	"Tree":                lucide.ListTree,
	"Sidebar":             lucide.PanelLeft,
	"Honeycomb":           lucide.Hexagon,
	"Split":               lucide.Columns2,
	"Split panes":         lucide.Columns2,
	"ScrollView":          lucide.ArrowDownUp,
	"CheckboxGroup":       lucide.ListChecks,
	"RadioGroup":          lucide.CircleDot,
	"StepSlider":          lucide.SlidersHorizontal,
	"RangeSlider":         lucide.SlidersHorizontal,
	"Stepper":             lucide.ListOrdered,
	"TextInputBuffer":     lucide.TextCursorInput,
	"TextAreaBuffer":      lucide.NotepadText,
	"TokenField":          lucide.Tags,
	"EditableText":        lucide.PencilLine,
	"DateInput":           lucide.CalendarClock,
	"TimeInput":           lucide.Clock,
	"ColorWell":           lucide.Pipette,
	"Table":               lucide.Table2,
	"Outline":             lucide.ListTree,
	"OutlineTable":        lucide.TableProperties,
	"Scroll":              lucide.ArrowDownUp,
	"ScrollHorizontal":    lucide.ArrowLeftRight,
	"ScrollBoth":          lucide.Move,
	"Grid":                lucide.Grid2x2,
	"SplitVertical":       lucide.Rows2,
	"Modal":               lucide.AppWindow,
	"Toast":               lucide.MessageSquare,
	"Box":                 lucide.Square,
	"Row":                 lucide.Columns3,
	"Column":              lucide.Rows3,
	"Spacer":              lucide.AlignHorizontalSpaceAround,
	"Divider":             lucide.SeparatorHorizontal,
	"Terminal":            lucide.Terminal,
	"Flyout":              lucide.PanelsTopLeft,
	"WebView":             lucide.Globe,
	"Glass":               lucide.GlassWater,
	"Audio":               lucide.AudioLines,
	"System information":  lucide.Info,
	"Recent documents":    lucide.FileClock,
	"File watching":       lucide.FolderSync,
	"Global shortcuts":    lucide.Keyboard,
	"Canvas":              lucide.Brush,
	"Code editor":         lucide.SquareCode,
	"Rich text editor":    lucide.NotebookPen,
	"Markdown stream":     lucide.MessageSquareText,
	"Diff viewer":         lucide.Diff,
	"Game loop":           lucide.Gamepad2,
	"3D renderer":         lucide.Cuboid,

	"Gradients & blur": lucide.Blend,
	"Shaders":          lucide.Sparkles,
	"Live chart":       lucide.ChartSpline,
}

func (entry navigationEntry) icon() *ui.SVG {
	if entry.mygo {
		return lucide.Boxes()
	}
	switch entry.page {
	case welcomePage:
		return lucide.House()
	case lucidePage:
		return lucide.Shapes()
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
