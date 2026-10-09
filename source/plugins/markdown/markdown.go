package markdown

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
	"github.com/oneo-me/onevra/icons/lucide"
	"github.com/oneo-me/onevra/plugins/code"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// MarkdownCodeBlock describes a fenced block. Code preserves its trailing
// newline; Info is the full decoded info string, Language its first word.
// Index is the zero-based fenced-block position in this document, including
// unhandled blocks. Hosts may use it to identify application-owned view state.
type MarkdownCodeBlock struct {
	Language, Info, Code string
	Index                int
}

// MarkdownOptions configures application-owned interactions and block views.
type MarkdownOptions struct {
	// CopyCodeLabel names the code-block copy button and its tooltip.
	// Empty uses "Copy code". Copy preserves the source's trailing newline.
	CopyCodeLabel string
	// CodeBlocks maps exact, case-sensitive language names to synchronous UI-thread
	// builders. A nil or missing builder uses the default code view. Builders run
	// every frame, must not retain Elements or start asynchronous work, and own
	// their state. Their ordinary text joins the document selection; use
	// Unselectable or a nested Selectable container to override that behavior.
	CodeBlocks map[string]func(*ui.Context, MarkdownCodeBlock)
	// Images supplies decoded bitmaps keyed by the decoded Markdown destination.
	// A paragraph containing only an image displays it at its intrinsic size,
	// shrinking proportionally to the available width. Its alt text supplies the
	// accessible label. Missing/nil images and images mixed into inline content
	// display their alt text instead. Applications own loading and caching;
	// rendering never reads files or fetches URLs.
	Images map[string]*ui.Bitmap
	// OnLink handles link activation, including relative and fragment targets.
	// Nil leaves links as styled text. No URL is opened automatically.
	OnLink func(destination string)
}

// MarkdownDocument is a parsed, immutable CommonMark document with GFM tables,
// strikethrough, autolinks and read-only task checkboxes. Retain it across frames and replace
// it with ParseMarkdown when source changes. Its zero value is an empty document.
type MarkdownDocument struct {
	source []byte
	root   ast.Node
	code   map[ast.Node]markdownCode
}

// ParseMarkdown parses source without creating UI resources or loading URLs.
// Fenced code is highlighted once using its declared language. Retain the result
// across frames; parsing can run in application-owned background work.
// Raw HTML is displayed literally; images use supplied bitmaps or alt text.
func ParseMarkdown(source string) *MarkdownDocument {
	data := []byte(source)
	d := &MarkdownDocument{source: data, root: goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser().Parse(text.NewReader(data)), code: make(map[ast.Node]markdownCode)}
	index := 0
	_ = ast.Walk(d.root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		b := MarkdownCodeBlock{}
		switch n := n.(type) {
		case *ast.FencedCodeBlock:
			b.Index = index
			index++
			if n.Info != nil {
				b.Info = markdownDecode(n.Info.Text(data))
				if fields := strings.Fields(b.Info); len(fields) > 0 {
					b.Language = fields[0]
				}
			}
		case *ast.CodeBlock:
		default:
			return ast.WalkContinue, nil
		}
		b.Code = string(n.Lines().Value(data))
		d.code[n] = markdownCode{block: b, text: code.Parse(strings.TrimSuffix(b.Code, "\n"), b.Language)}
		return ast.WalkContinue, nil
	})
	return d
}

// Markdown parses and renders source. For repeatedly rendered or large documents,
// retain ParseMarkdown's result and call its Render method instead.
func Markdown(c *ui.Context, source string, options ...MarkdownOptions) ui.Element {
	return ParseMarkdown(source).Render(c, options...)
}

// Render builds a selectable document using Onevra text and MyGo layout/input
// primitives. Selection spans paragraphs and default code blocks. Copy produces
// rendered text, not Markdown source. The host owns scrolling; use a regular,
// non-virtualized ScrollView to retain whole-document selection. At most one
// options value is accepted. The returned container supports normal styling.
func (d *MarkdownDocument) Render(c *ui.Context, options ...MarkdownOptions) ui.Element {
	if len(options) > 1 {
		panic("markdown.Markdown: at most one options value")
	}
	var opts MarkdownOptions
	if len(options) == 1 {
		opts = options[0]
	}
	r := markdownRenderer{c: c, document: d, options: opts}
	return ui.Column(c).MinWidth(0).Gap(c.Theme().Space(3)).LineHeight(1.6).Selectable().Children(func() {
		if d != nil && d.root != nil {
			r.blocks(d.root)
		}
	})
}

type markdownCode struct {
	block MarkdownCodeBlock
	text  *code.Document
}

type markdownRenderer struct {
	c        *ui.Context
	document *MarkdownDocument
	options  MarkdownOptions
}

func (r *markdownRenderer) blocks(parent ast.Node) {
	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		r.block(n)
	}
}

func (r *markdownRenderer) block(n ast.Node) {
	c, t := r.c, r.c.Theme()
	switch n := n.(type) {
	case *ast.Paragraph, *ast.TextBlock:
		// Reference definitions can leave an empty TextBlock in the AST.
		// They have no visible content and must not reserve paragraph spacing.
		if n.FirstChild() == nil {
			return
		}
		if image, ok := n.FirstChild().(*ast.Image); ok && image.NextSibling() == nil && r.image(image) {
			return
		}
		r.paragraph(n)
	case *ast.Heading:
		r.paragraph(n).FontSize(t.Rem([]float32{1.5, 1.25, 1.125, 1, 1, 1}[n.Level-1])).
			FontWeight(600).LineHeight(1.4).Margin(t.Space([]float32{0, 2, 1, 0, 0, 0}[n.Level-1]), 0, 0, 0)
	case *ast.FencedCodeBlock:
		prepared := r.document.code[n]
		b := prepared.block
		ui.Column(c).MinWidth(0).Key(b.Index).Children(func() {
			if build := r.options.CodeBlocks[b.Language]; build != nil {
				build(c, b)
			} else {
				r.code(prepared)
			}
		})
	case *ast.CodeBlock:
		r.code(r.document.code[n])
	case *ast.Blockquote:
		ui.Column(c).MinWidth(0).Padding(0, 0, 0, t.Space(4)).Gap(t.Space(1)).BorderColor(t.Border).BorderWidth(0, 0, 0, t.Space(.75)).Children(func() { r.blocks(n) })
	case *ast.List:
		ui.Column(c).MinWidth(0).Gap(t.Space(1)).Children(func() {
			index := n.Start
			for item := n.FirstChild(); item != nil; item = item.NextSibling() {
				marker := "•"
				if n.IsOrdered() {
					marker = fmt.Sprintf("%d%c", index, n.Marker)
					index++
				}
				ui.Row(c).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
					if task := markdownTask(item); task != nil {
						checked := task.IsChecked
						// Keep the source immutable and the paragraph selectable;
						// the checkbox ignores pointer and keyboard activation.
						ui.Box(c).PassThrough().Unselectable().Shrink(0).
							Margin((t.Rem(1)*1.6-t.Space(4))/2, 0, 0, 0).Children(func() {
							components.Checkbox(c, &checked, "").PassThrough().
								HandleInput(func(ev ui.InputEvent) bool {
									return ev.Kind == ui.InputKeyDown && (ev.Key == ui.KeySpace || ev.Key == ui.KeyEnter)
								}).Label(markdownDecode(item.FirstChild().Text(r.document.source)))
						})
					} else {
						components.Text(c, marker).MinWidth(t.Rem(.75)).TextColor(t.TextMuted).Unselectable()
					}
					ui.Column(c).MinWidth(0).Grow(1).Gap(t.Space(1)).Children(func() { r.blocks(item) })
				})
			}
		})
	case *ast.ThematicBreak:
		ui.Box(c).Height(1).MarginY(t.Space(2)).Background(t.Border).Unselectable()
	case *ast.HTMLBlock:
		value := string(n.Lines().Value(r.document.source))
		if n.HasClosure() {
			value += string(n.ClosureLine.Value(r.document.source))
		}
		components.Text(c, strings.TrimSuffix(value, "\n"))
	case *east.Table:
		ui.Column(c).MinWidth(0).FontSize(t.Rem(.875)).Border(1, t.Border).Radius(t.Rem(.5)).Clip().Children(func() { r.blocks(n) })
	case *east.TableHeader, *east.TableRow:
		row := ui.Row(c).MinWidth(0).AlignItems(ui.Stretch)
		if _, header := n.(*east.TableHeader); header {
			row.FontWeight(600).Background(t.Text.Alpha(.06))
		}
		row.Children(func() {
			for cell := n.FirstChild(); cell != nil; cell = cell.NextSibling() {
				e := ui.Box(c).Width(0).Grow(1).MinWidth(0).Padding(t.Space(2), t.Space(3))
				var bottom, left float32
				if n.NextSibling() != nil {
					bottom = 1
				}
				if cell.PreviousSibling() != nil {
					left = 1
				}
				e.BorderColor(t.Border).BorderWidth(0, 0, bottom, left)
				e.Children(func() {
					p := r.paragraph(cell)
					if cell, ok := cell.(*east.TableCell); ok {
						switch cell.Alignment {
						case east.AlignCenter:
							p.TextAlign(ui.Center)
						case east.AlignRight:
							p.TextAlign(ui.End)
						}
					}
				})
			}
		})
	default:
		r.blocks(n)
	}
}

func markdownTask(item ast.Node) *east.TaskCheckBox {
	if block := item.FirstChild(); block != nil {
		task, _ := block.FirstChild().(*east.TaskCheckBox)
		return task
	}
	return nil
}

func (r *markdownRenderer) code(prepared markdownCode) {
	c, t := r.c, r.c.Theme()
	box := ui.Box(c).MinWidth(0).Padding(t.Space(3), t.Space(10), t.Space(3), t.Space(3.5)).
		Radius(t.Rem(.75)).Border(1, t.Border).Background(t.Text.Alpha(.06))
	hovered := box.Hovered()
	box.Children(func() {
		prepared.text.Render(c).FontSize(t.Rem(.85)).LineHeight(1.6)
		if hovered {
			label := r.options.CopyCodeLabel
			if label == "" {
				label = "Copy code"
			}
			button := components.Button(c, "", components.ButtonGhost).
				Absolute().Top(t.Space(1.5)).Right(t.Space(1.5)).
				Size(t.Space(6), t.Space(6)).Padding(0).Label(label).Children(func() {
				components.Icon(c, lucide.Copy).Size(t.Rem(1), t.Rem(1))
			})
			components.Tooltip(c, button, label)
			if button.Clicked() {
				c.WriteClipboard(prepared.block.Code)
			}
		}
	})
}

func (r *markdownRenderer) image(n *ast.Image) bool {
	bitmap := r.options.Images[markdownDecode(n.Destination)]
	if bitmap == nil {
		return false
	}
	w, h := bitmap.Size()
	if w <= 0 || h <= 0 {
		return false
	}
	components.Image(r.c, bitmap).Width(float32(w)).MaxWidthPercent(100).MinWidth(0).
		AspectRatio(float32(w)/float32(h)).Radius(r.c.Theme().Rem(.75)).Border(1, r.c.Theme().Border).Label(markdownDecode(n.Text(r.document.source)))
	return true
}

func (r *markdownRenderer) paragraph(n ast.Node) ui.Element {
	return components.RichText(r.c).MinWidth(0).Children(func() { r.inlines(n) })
}

func markdownDecode(value []byte) string {
	return string(util.ResolveEntityNames(util.ResolveNumericReferences(util.UnescapePunctuations(value))))
}

func (r *markdownRenderer) inlines(parent ast.Node) {
	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		r.inline(n)
	}
}

func (r *markdownRenderer) link(n ast.Node, label, destination string) {
	if r.options.OnLink == nil {
		components.Text(r.c, label).TextColor(r.c.Theme().Accent).FontWeight(500).Underline().DecorationColor(r.c.Theme().Border).Children(func() { r.inlines(n) })
		return
	}
	e := components.Button(r.c, label, components.ButtonLink).FontWeight(500).Underline().DecorationColor(r.c.Theme().Border).Children(func() { r.inlines(n) })
	if e.Hovered() {
		e.DecorationColor(r.c.Theme().Accent)
	}
	if e.Clicked() {
		r.options.OnLink(destination)
	}
}

func (r *markdownRenderer) inline(n ast.Node) {
	c, source := r.c, r.document.source
	switch n := n.(type) {
	case *ast.Text:
		value := string(n.Segment.Value(source))
		if !n.IsRaw() {
			value = markdownDecode([]byte(value))
		}
		if n.HardLineBreak() {
			value += "\n"
		} else if n.SoftLineBreak() {
			value += " "
		}
		components.Text(c, value)
	case *ast.String:
		value := string(n.Value)
		if !n.IsRaw() && !n.IsCode() {
			value = markdownDecode(n.Value)
		}
		components.Text(c, value)
	case *ast.Emphasis:
		e := components.RichText(c).Children(func() { r.inlines(n) })
		if n.Level == 2 {
			e.Bold()
		} else {
			e.Italic()
		}
	case *east.Strikethrough:
		components.RichText(c).Strikethrough().Children(func() { r.inlines(n) })
	case *ast.CodeSpan:
		var value strings.Builder
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			value.Write(child.(*ast.Text).Segment.Value(source))
		}
		components.Text(c, strings.ReplaceAll(value.String(), "\n", " ")).Font("monospace").FontSize(c.Theme().Rem(.85)).Background(c.Theme().Text.Alpha(.06))
	case *ast.Link:
		r.link(n, "", markdownDecode(n.Destination))
	case *ast.AutoLink:
		r.link(n, string(n.Label(source)), string(n.URL(source)))
	case *ast.Image:
		r.inlines(n)
	case *ast.RawHTML:
		components.Text(c, string(n.Segments.Value(source)))
	case *east.TaskCheckBox:
		// The containing list row renders the Onevra checkbox outside rich text.
	default:
		r.inlines(n)
	}
}
