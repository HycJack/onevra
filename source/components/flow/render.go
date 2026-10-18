package flow

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/gesture"
)

const (
	defaultNodeGap   = 40
	defaultRankGap   = 72
	edgeGap          = 24
	fallbackWidth    = 200
	fallbackHeight   = 72
	planeRoom        = 1 << 14
	mutedNodeOpacity = .45
	mutedEdgeOpacity = .35
	minZoom          = .1
	maxZoom          = 2.5
	zoomStep         = 1.25
	fitPadding       = 32
	cameraResponse   = 90 * time.Millisecond
	maxFrameStep     = 48 * time.Millisecond
	positionEpsilon  = .05
	zoomEpsilon      = .0005
	clickTolerance   = 4
	revealMargin     = 16
	defaultEdgeWidth = 1.25
	arrowLength      = 7
	arrowHalfWidth   = 4
	dashLength       = 5
	dashGap          = 4
	ringGap          = 3
	ringWidth        = 2
	labelSize        = 11
	labelPadX        = 6
	labelPadY        = 2
)

type camera struct{ x, y, zoom float64 }

type view struct {
	initialized bool
	nodes       []renderNode
	edges       []renderEdge
	index       map[string]int
	direction   Direction
	gaps        spacing
	layoutKey   uint64
	result      layoutResult

	reduceMotion bool
	selected     string
	selectFn     func(string)
	hovered      string
	focused      string

	ready, userMoved bool
	// camera is where the view is going and shown where it is this frame.
	camera, shown camera
	width, height float64
	lastFrame     time.Time

	panning, moved bool
	pressClicks    int
	startX, startY float64
	lastX, lastY   float64
	navigation     gesture.Navigation
}

func (v *view) init() {
	if v.initialized {
		return
	}
	v.initialized = true
	v.camera, v.shown = camera{zoom: 1}, camera{zoom: 1}
}

// resize keeps a diagram the user has not moved fitted to the view, and
// otherwise keeps the same point at the view's center.
func (v *view) resize(width, height float64) {
	if width <= 0 || height <= 0 || (width == v.width && height == v.height) {
		return
	}
	if v.width == 0 || !v.ready {
		v.width, v.height = width, height
		return
	}
	dx, dy := (width-v.width)/2, (height-v.height)/2
	v.width, v.height = width, height
	if !v.userMoved {
		v.fit(true)
		return
	}
	v.camera.x, v.camera.y = v.camera.x+dx, v.camera.y+dy
	v.shown.x, v.shown.y = v.shown.x+dx, v.shown.y+dy
}

// advanceCamera eases the shown camera toward its target and reports
// whether it is still moving.
func (v *view) advanceCamera(now time.Time) bool {
	elapsed := 16 * time.Millisecond
	if !v.lastFrame.IsZero() {
		elapsed = min(max(0, now.Sub(v.lastFrame)), maxFrameStep)
	}
	v.lastFrame = now
	dx, dy, dz := v.camera.x-v.shown.x, v.camera.y-v.shown.y, v.camera.zoom-v.shown.zoom
	if v.reduceMotion || (math.Abs(dx) < positionEpsilon && math.Abs(dy) < positionEpsilon && math.Abs(dz) < zoomEpsilon) {
		v.shown = v.camera
		return false
	}
	amount := 1 - math.Exp(-float64(elapsed)/float64(cameraResponse))
	v.shown = camera{v.shown.x + dx*amount, v.shown.y + dy*amount, v.shown.zoom + dz*amount}
	return true
}

func (v *view) fit(immediate bool) {
	if len(v.nodes) == 0 || v.width == 0 {
		return
	}
	w, h := max(1, v.result.width), max(1, v.result.height)
	zoom := clamp(min(max(1, v.width-2*fitPadding)/w, max(1, v.height-2*fitPadding)/h), minZoom, 1)
	v.camera = camera{v.width/2 - w/2*zoom, v.height/2 - h/2*zoom, zoom}
	v.userMoved = false
	if immediate || v.reduceMotion {
		v.shown = v.camera
	}
}

// zoomAt scales the view around a point of the view, immediately or toward
// the camera's target.
func (v *view) zoomAt(x, y, factor float64, animated bool) {
	from := v.shown
	if animated {
		from = v.camera
	}
	zoom := clamp(from.zoom*factor, minZoom, maxZoom)
	gx, gy := (x-from.x)/from.zoom, (y-from.y)/from.zoom
	v.camera = camera{x - gx*zoom, y - gy*zoom, zoom}
	v.userMoved = true
	if !animated || v.reduceMotion {
		v.shown = v.camera
	}
}

func (v *view) reveal(id string) {
	i, ok := v.index[id]
	if !ok || i >= len(v.result.x) || v.width == 0 {
		return
	}
	w, h := v.nodes[i].size()
	cx, cy := v.result.x[i]+w/2, v.result.y[i]+h/2
	v.camera.x, v.camera.y = v.width/2-cx*v.camera.zoom, v.height/2-cy*v.camera.zoom
	v.userMoved = true
	if v.reduceMotion {
		v.shown = v.camera
	}
}

// revealIfHidden reveals a node that is not wholly in view.
func (v *view) revealIfHidden(id string) {
	i, ok := v.index[id]
	if !ok || i >= len(v.result.x) {
		return
	}
	w, h := v.nodes[i].size()
	z := v.camera.zoom
	x, y := v.result.x[i]*z+v.camera.x, v.result.y[i]*z+v.camera.y
	if x < revealMargin || y < revealMargin || x+w*z > v.width-revealMargin || y+h*z > v.height-revealMargin {
		v.reveal(id)
	}
}

// hit returns the node shown under a point of the view, or -1.
func (v *view) hit(x, y float64) int {
	z := v.shown.zoom
	if z <= 0 {
		return -1
	}
	gx, gy := (x-v.shown.x)/z, (y-v.shown.y)/z
	for i := len(v.nodes) - 1; i >= 0; i-- {
		n := &v.nodes[i]
		w, h := n.size()
		if gx >= n.x && gx <= n.x+w && gy >= n.y && gy <= n.y+h {
			return i
		}
	}
	return -1
}

func (v *view) cursor() ui.Cursor {
	if v.panning {
		return ui.CursorGrabbing
	}
	return ui.CursorGrab
}

// input pans when empty space is dragged, leaving presses on nodes to
// them. Clicking empty space clears the selection and double-clicking it
// zooms.
func (v *view) input(ev ui.InputEvent) bool {
	x, y := float64(ev.X), float64(ev.Y)
	switch ev.Kind {
	case ui.InputPointerDown:
		// Button 2 is the middle button.
		if !v.ready || (ev.Button != 0 && ev.Button != 2) || v.hit(x, y) >= 0 {
			return false
		}
		v.panning, v.moved, v.pressClicks = true, false, ev.Clicks
		v.startX, v.startY, v.lastX, v.lastY = x, y, x, y
		v.camera = v.shown
		return true
	case ui.InputPointerMove:
		if !v.panning {
			return false
		}
		if math.Hypot(x-v.startX, y-v.startY) > clickTolerance {
			v.moved = true
		}
		v.camera.x += x - v.lastX
		v.camera.y += y - v.lastY
		v.shown, v.lastX, v.lastY = v.camera, x, y
		v.userMoved = true
		return true
	case ui.InputPointerUp:
		if !v.panning {
			return false
		}
		v.panning = false
		if v.moved || ev.Button != 0 {
			return true
		}
		if v.pressClicks >= 2 {
			factor := 2.0
			if ev.Mods&ui.Shift != 0 {
				factor = .5
			}
			v.zoomAt(x, y, factor, true)
		} else if v.selectFn != nil && v.selected != "" {
			v.selectFn("")
		}
		return true
	}
	return false
}

// navigate pans and zooms: the wheel and pinches zoom around the pointer,
// and two fingers on a touchpad pan.
func (v *view) navigate(e gesture.Event) {
	if !v.ready {
		return
	}
	switch e.Kind {
	case gesture.Pan:
		v.camera = camera{v.shown.x + float64(e.DX), v.shown.y + float64(e.DY), v.shown.zoom}
		v.shown = v.camera
		v.userMoved = true
	case gesture.Zoom:
		v.zoomAt(float64(e.X), float64(e.Y), e.Scale, false)
	}
}

type colors struct {
	line, emphasis, ring          ui.Color
	label, labelFace, labelBorder ui.Color
}

func (v *view) palette(t *ui.Theme) colors {
	return colors{
		line: t.TextMuted, emphasis: t.Accent, ring: t.Accent,
		label: t.TextMuted, labelFace: t.Background, labelBorder: t.Border,
	}
}

// emphasized reports whether an edge touches the hovered or selected node.
func (v *view) emphasized(e renderEdge) bool {
	for _, id := range [...]string{v.hovered, v.selected} {
		if id == "" {
			continue
		}
		if i, ok := v.index[id]; ok && (e.source == i || e.target == i) {
			return true
		}
	}
	return false
}

// shows reports whether both of an edge's ends show.
func (v *view) shows(e renderEdge) bool {
	return v.nodes[e.source].measured && v.nodes[e.target].measured
}

// route returns an edge's points where its ends show this frame: while
// nodes move to a new layout, each point follows the end it is nearer
// along the route.
func (v *view) route(i int, ox, oy float64) []point {
	if i >= len(v.result.routes) {
		return nil
	}
	laid := v.result.routes[i]
	if len(laid) < 2 {
		return nil
	}
	e := v.edges[i]
	shift := func(n int) point {
		return point{v.nodes[n].x - v.result.x[n], v.nodes[n].y - v.result.y[n]}
	}
	s, t := shift(e.source), shift(e.target)
	total := 0.0
	for k := 1; k < len(laid); k++ {
		total += math.Hypot(laid[k].x-laid[k-1].x, laid[k].y-laid[k-1].y)
	}
	out := make([]point, len(laid))
	along := 0.0
	for k, p := range laid {
		if k > 0 {
			along += math.Hypot(p.x-laid[k-1].x, p.y-laid[k-1].y)
		}
		w := 0.0
		if total > 0 {
			w = along / total
		}
		out[k] = point{p.x + s.x*(1-w) + t.x*w + ox, p.y + s.y*(1-w) + t.y*w + oy}
	}
	return out
}

func (v *view) paintEdges(p *ui.Painter, r ui.Rect, c colors) {
	ox, oy := float64(r.X), float64(r.Y)
	for pass := range 2 {
		for i, e := range v.edges {
			if v.emphasized(e) != (pass == 1) || !v.shows(e) {
				continue
			}
			points := v.route(i, ox, oy)
			if points == nil {
				continue
			}
			color, width := e.color, e.width
			if width <= 0 {
				width = defaultEdgeWidth
			}
			if pass == 1 {
				width += .5
				if color == (ui.Color{}) {
					color = c.emphasis
				}
			}
			if color == (ui.Color{}) {
				color = c.line
			}
			if v.nodes[e.source].muted || v.nodes[e.target].muted {
				color = color.Alpha(mutedEdgeOpacity)
			}
			drawRoute(p, points, width, color, e.dashed)
		}
	}
	for i, e := range v.edges {
		if e.label == "" || !v.shows(e) {
			continue
		}
		points := v.route(i, ox, oy)
		if points == nil {
			continue
		}
		alpha := float32(1)
		if v.nodes[e.source].muted || v.nodes[e.target].muted {
			alpha = mutedEdgeOpacity
		}
		x, y := labelAnchor(points)
		span := ui.Span{Text: e.label, Size: labelSize, Weight: 500, Color: c.label.Alpha(alpha)}
		w, h := p.MeasureText(0, span)
		box := ui.Rect{X: float32(x) - w/2 - labelPadX, Y: float32(y) - h/2 - labelPadY, W: w + 2*labelPadX, H: h + 2*labelPadY}
		p.Fill(box, c.labelFace.Alpha(alpha), box.H/2)
		p.Stroke(box, c.labelBorder.Alpha(alpha), box.H/2, 1)
		p.RichText(box.X+labelPadX, box.Y+labelPadY, 0, span)
	}
}

// paintSelection rings the selected node where it shows.
func (v *view) paintSelection(p *ui.Painter, r ui.Rect, c colors) {
	i, ok := v.index[v.selected]
	if !ok || !v.nodes[i].measured {
		return
	}
	n := &v.nodes[i]
	w, h := n.size()
	box := ui.Rect{X: r.X + float32(n.x) - ringGap, Y: r.Y + float32(n.y) - ringGap, W: float32(w) + 2*ringGap, H: float32(h) + 2*ringGap}
	p.Stroke(box, c.ring, CardRadius+ringGap, ringWidth)
}

// labelAnchor is the middle of a route's longest segment.
func labelAnchor(points []point) (float64, float64) {
	best, length := 1, -1.0
	for k := 1; k < len(points); k++ {
		if l := math.Hypot(points[k].x-points[k-1].x, points[k].y-points[k-1].y); l > length {
			best, length = k, l
		}
	}
	a, b := points[best-1], points[best]
	return (a.x + b.x) / 2, (a.y + b.y) / 2
}

// drawRoute draws a route as axis-aligned fills, which the GPU draws
// directly, ending in an arrowhead. Segments a transition bends off the
// axes draw as paths.
func drawRoute(p *ui.Painter, points []point, width float32, c ui.Color, dashed bool) {
	n := len(points)
	end, before := points[n-1], points[n-2]
	dx, dy := end.x-before.x, end.y-before.y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}
	ux, uy := dx/length, dy/length
	head := math.Min(arrowLength, length)
	line := append([]point(nil), points...)
	line[n-1] = point{end.x - ux*head*.8, end.y - uy*head*.8}
	phase := 0.0
	for k := 1; k < n; k++ {
		a, b := line[k-1], line[k]
		if dashed {
			phase = dashes(p, a, b, phase, width, c)
		} else {
			segment(p, a, b, width, c)
		}
	}
	var arrow ui.Path
	arrow.MoveTo(float32(end.x), float32(end.y)).
		LineTo(float32(end.x-ux*head-uy*arrowHalfWidth), float32(end.y-uy*head+ux*arrowHalfWidth)).
		LineTo(float32(end.x-ux*head+uy*arrowHalfWidth), float32(end.y-uy*head-ux*arrowHalfWidth)).
		Close()
	p.FillPath(&arrow, c)
}

// segment draws a straight line with square ends, so that bends join.
func segment(p *ui.Painter, a, b point, width float32, c ui.Color) {
	half := float64(width) / 2
	switch {
	case math.Abs(a.y-b.y) < .01:
		p.Fill(ui.Rect{X: float32(min(a.x, b.x) - half), Y: float32(a.y - half), W: float32(math.Abs(b.x-a.x) + 2*half), H: width}, c, 0)
	case math.Abs(a.x-b.x) < .01:
		p.Fill(ui.Rect{X: float32(a.x - half), Y: float32(min(a.y, b.y) - half), W: width, H: float32(math.Abs(b.y-a.y) + 2*half)}, c, 0)
	default:
		var line ui.Path
		line.MoveTo(float32(a.x), float32(a.y)).LineTo(float32(b.x), float32(b.y))
		p.StrokePath(&line, width, c)
	}
}

// dashes draws a dashed segment continuing the pattern at phase, and
// returns the phase at its end.
func dashes(p *ui.Painter, a, b point, phase float64, width float32, c ui.Color) float64 {
	length := math.Hypot(b.x-a.x, b.y-a.y)
	if length == 0 {
		return phase
	}
	ux, uy := (b.x-a.x)/length, (b.y-a.y)/length
	period := float64(dashLength + dashGap)
	at := 0.0
	for at < length {
		offset := math.Mod(phase+at, period)
		if offset < dashLength {
			run := math.Min(dashLength-offset, length-at)
			from := point{a.x + ux*at, a.y + uy*at}
			to := point{a.x + ux*(at+run), a.y + uy*(at+run)}
			segment(p, from, to, width, c)
			at += run
		} else {
			at += period - offset
		}
	}
	return math.Mod(phase+length, period)
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(v, hi)) }
