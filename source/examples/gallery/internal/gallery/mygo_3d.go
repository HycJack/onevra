package gallery

import (
	"image"
	"math"

	"github.com/egoist/mygo/plugins/game"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// The rasterizer draws only the objects: the buffer is transparent between
// them, so the elements float over whatever the application paints behind
// the view. sceneNear is the near plane, sceneFOV the vertical field of view
// in degrees, and sceneMaxPixels caps the CPU cost on large or high-density
// windows.
const (
	sceneNear      = 0.05
	sceneFOV       = 52
	sceneMaxPixels = 1_600_000
)

// vec3 is a point or direction in world space, right-handed with Y up.
type vec3 struct{ x, y, z float32 }

func v3(x, y, z float32) vec3 { return vec3{x, y, z} }

func (a vec3) add(b vec3) vec3    { return v3(a.x+b.x, a.y+b.y, a.z+b.z) }
func (a vec3) sub(b vec3) vec3    { return v3(a.x-b.x, a.y-b.y, a.z-b.z) }
func (a vec3) mul(s float32) vec3 { return v3(a.x*s, a.y*s, a.z*s) }
func (a vec3) dot(b vec3) float32 { return a.x*b.x + a.y*b.y + a.z*b.z }
func (a vec3) cross(b vec3) vec3 {
	return v3(a.y*b.z-a.z*b.y, a.z*b.x-a.x*b.z, a.x*b.y-a.y*b.x)
}

func (a vec3) unit() vec3 {
	l := float32(math.Sqrt(float64(a.dot(a))))
	if l == 0 {
		return a
	}
	return a.mul(1 / l)
}

// mat3 is a rotation, kept as its three columns for transforming points.
type mat3 struct{ c0, c1, c2 vec3 }

// euler builds the rotation Rx*Ry*Rz from the angles in radians.
func euler(a vec3) mat3 {
	sx, cx := sincos(a.x)
	sy, cy := sincos(a.y)
	sz, cz := sincos(a.z)
	return mat3{
		c0: v3(cy*cz, cx*sz+sx*sy*cz, sx*sz-cx*sy*cz),
		c1: v3(-cy*sz, cx*cz-sx*sy*sz, sx*cz+cx*sy*sz),
		c2: v3(sy, -sx*cy, cx*cy),
	}
}

func (m mat3) rotate(v vec3) vec3 {
	return m.c0.mul(v.x).add(m.c1.mul(v.y)).add(m.c2.mul(v.z))
}

// sceneMesh is a triangle mesh with a shared normal per vertex. Winding is
// made consistent with the normals by orient, so back-face culling can trust
// the geometric normal.
type sceneMesh struct {
	pos    []vec3
	nrm    []vec3
	idx    []int
	albedo [3]float32
	spec   float32
	shine  int
}

func (m *sceneMesh) vert(p, n vec3) int {
	m.pos = append(m.pos, p)
	m.nrm = append(m.nrm, n)
	return len(m.pos) - 1
}

func (m *sceneMesh) tri(a, b, c int) { m.idx = append(m.idx, a, b, c) }

func (m *sceneMesh) quad(a, b, c, d int) {
	m.tri(a, b, c)
	m.tri(a, c, d)
}

// orient flips every triangle when the first one faces away from its vertex
// normals, so generators may be written in either handedness.
func (m *sceneMesh) orient() {
	if len(m.idx) < 3 {
		return
	}
	a, b, c := m.idx[0], m.idx[1], m.idx[2]
	g := m.pos[b].sub(m.pos[a]).cross(m.pos[c].sub(m.pos[a]))
	n := m.nrm[a].add(m.nrm[b]).add(m.nrm[c])
	if g.dot(n) >= 0 {
		return
	}
	for i := 0; i+2 < len(m.idx); i += 3 {
		m.idx[i+1], m.idx[i+2] = m.idx[i+2], m.idx[i+1]
	}
}

func cubeMesh(albedo [3]float32, spec float32, shine int) *sceneMesh {
	m := &sceneMesh{albedo: albedo, spec: spec, shine: shine}
	add := func(n, a, b, c, d vec3) {
		m.quad(m.vert(a, n), m.vert(b, n), m.vert(c, n), m.vert(d, n))
	}
	const h = 0.5
	add(v3(0, 0, 1), v3(-h, -h, h), v3(h, -h, h), v3(h, h, h), v3(-h, h, h))
	add(v3(0, 0, -1), v3(h, -h, -h), v3(-h, -h, -h), v3(-h, h, -h), v3(h, h, -h))
	add(v3(1, 0, 0), v3(h, -h, h), v3(h, -h, -h), v3(h, h, -h), v3(h, h, h))
	add(v3(-1, 0, 0), v3(-h, -h, -h), v3(-h, -h, h), v3(-h, h, h), v3(-h, h, -h))
	add(v3(0, 1, 0), v3(-h, h, h), v3(h, h, h), v3(h, h, -h), v3(-h, h, -h))
	add(v3(0, -1, 0), v3(-h, -h, -h), v3(h, -h, -h), v3(h, -h, h), v3(-h, -h, h))
	m.orient()
	return m
}

// sphereMesh builds a unit UV sphere.
func sphereMesh(stacks, slices int, albedo [3]float32, spec float32, shine int) *sceneMesh {
	m := &sceneMesh{albedo: albedo, spec: spec, shine: shine}
	for i := 0; i <= stacks; i++ {
		theta := math.Pi * float64(i) / float64(stacks)
		st, ct := math.Sincos(theta)
		for j := 0; j <= slices; j++ {
			phi := 2 * math.Pi * float64(j) / float64(slices)
			sp, cp := math.Sincos(phi)
			n := v3(float32(st*cp), float32(ct), float32(st*sp))
			m.vert(n, n)
		}
	}
	row := slices + 1
	for i := 0; i < stacks; i++ {
		for j := 0; j < slices; j++ {
			a := i*row + j
			m.quad(a, a+1, a+row+1, a+row)
		}
	}
	m.orient()
	return m
}

// torusMesh builds a torus around the Y axis.
func torusMesh(segments, sides int, major, minor float32, albedo [3]float32, spec float32, shine int) *sceneMesh {
	m := &sceneMesh{albedo: albedo, spec: spec, shine: shine}
	for i := 0; i <= segments; i++ {
		u := 2 * math.Pi * float64(i) / float64(segments)
		su, cu := math.Sincos(u)
		for j := 0; j <= sides; j++ {
			v := 2 * math.Pi * float64(j) / float64(sides)
			sv, cv := math.Sincos(v)
			r := float64(major) + float64(minor)*cv
			p := v3(float32(r*cu), float32(float64(minor)*sv), float32(r*su))
			n := v3(float32(cv*cu), float32(sv), float32(cv*su))
			m.vert(p, n)
		}
	}
	row := sides + 1
	for i := 0; i < segments; i++ {
		for j := 0; j < sides; j++ {
			a := i*row + j
			m.quad(a, a+row, a+row+1, a+1)
		}
	}
	m.orient()
	return m
}

// sceneObject is a mesh placed in the world, spinning about its own axes.
type sceneObject struct {
	mesh  *sceneMesh
	pos   vec3
	angle vec3
	spin  vec3
	// bob lifts the object by bob*sin(t*2+phase) along Y.
	bob   float32
	phase float32
}

// sceneLight is a directional light.
type sceneLight struct {
	dir   vec3
	color [3]float32
}

var sceneLights = []sceneLight{
	{dir: v3(0.55, 0.78, 0.32).unit(), color: [3]float32{1.0, 0.95, 0.86}},
	{dir: v3(-0.62, 0.28, -0.55).unit(), color: [3]float32{0.24, 0.30, 0.44}},
}

// scene3D is a self-contained software 3D renderer over MyGo's bitmap and
// game loop: it transforms triangles, clips them at the near plane,
// rasterizes them into a transparent image.RGBA with a depth buffer, and
// shades every covered pixel. The bitmap is uploaded once a frame and drawn
// as any other image.
type scene3D struct {
	cube, sphere, torus *sceneMesh
	objs                []sceneObject

	w, h   int
	buf    *image.RGBA
	depth  []float32
	bitmap *ui.Bitmap

	yaw, pitch, dist float32
	target           vec3
	dragging         bool
	lastX, lastY     float32
	t                float32
}

func (d *scene3D) init() {
	if d.bitmap != nil {
		return
	}
	d.cube = cubeMesh([3]float32{0.95, 0.52, 0.24}, 0.7, 48)
	d.sphere = sphereMesh(28, 18, [3]float32{0.22, 0.72, 0.82}, 1.0, 64)
	d.torus = torusMesh(44, 18, 0.82, 0.30, [3]float32{0.62, 0.40, 0.92}, 0.8, 40)
	d.objs = []sceneObject{
		{mesh: d.cube, pos: v3(-1.9, -0.1, 0.4), angle: v3(0.15, 0.5, 0), spin: v3(0.35, 0.5, 0.2)},
		{mesh: d.sphere, pos: v3(1.7, -0.05, -0.2), spin: v3(0.2, 0.4, 0), bob: 0.18, phase: 0.7},
		{mesh: d.torus, pos: v3(0.05, 0.5, -1.5), angle: v3(0.9, 0, 0), spin: v3(0.9, 0, 0.5), bob: 0.22, phase: 2.1},
	}
	d.ensure(480, 300)
	d.bitmap = ui.NewBitmap(d.buf)
	d.yaw, d.pitch, d.dist, d.target = 0.62, 0.30, 6.9, v3(0, 0, -0.5)
}

// ensure makes the buffer and depth map at least w×h, reusing their storage
// when a resize keeps the same or a smaller size.
func (d *scene3D) ensure(w, h int) {
	if d.buf != nil && d.w == w && d.h == h {
		return
	}
	if d.buf == nil || cap(d.buf.Pix) < w*h*4 {
		d.buf = image.NewRGBA(image.Rect(0, 0, w, h))
	} else {
		d.buf.Pix = d.buf.Pix[:w*h*4]
		d.buf.Stride = w * 4
		d.buf.Rect = image.Rect(0, 0, w, h)
	}
	if cap(d.depth) < w*h {
		d.depth = make([]float32, w*h)
	} else {
		d.depth = d.depth[:w*h]
	}
	d.w, d.h = w, h
}

// size returns the render resolution for r: its device pixels, capped so the
// CPU renderer stays real-time. The aspect is kept.
func (d *scene3D) size(p *ui.Painter, r ui.Rect) (int, int) {
	s := p.Scale()
	w := max(1, int(math.Ceil(float64(r.W*s))))
	h := max(1, int(math.Ceil(float64(r.H*s))))
	if w*h > sceneMaxPixels {
		f := math.Sqrt(float64(sceneMaxPixels) / float64(w*h))
		w, h = max(1, int(float64(w)*f)), max(1, int(float64(h)*f))
	}
	return w, h
}

// Update orbits the camera with the pointer and keys and spins the objects.
func (d *scene3D) Update(in *game.Input) {
	const step = 1.0 / 60.0
	d.t += step
	rate := float32(0.035)
	// Keys mirror the drag: right spins the objects right, down shows their
	// top, as a trackball turns.
	if in.Down(ui.KeyLeft) || in.Down(ui.KeyA) {
		d.yaw += rate
	}
	if in.Down(ui.KeyRight) || in.Down(ui.KeyD) {
		d.yaw -= rate
	}
	if in.Down(ui.KeyUp) || in.Down(ui.KeyW) {
		d.pitch = clampf(d.pitch-rate, 0.06, 1.35)
	}
	if in.Down(ui.KeyDown) || in.Down(ui.KeyS) {
		d.pitch = clampf(d.pitch+rate, 0.06, 1.35)
	}
	x, y, _ := in.Pointer()
	if in.PointerPressed(0) {
		d.dragging, d.lastX, d.lastY = true, x, y
	}
	if in.PointerReleased(0) {
		d.dragging = false
	}
	if d.dragging && in.PointerDown(0) {
		// Dragging right turns the objects right and dragging down shows
		// their top, as the object follows the pointer.
		d.yaw -= (x - d.lastX) * 0.012
		d.pitch = clampf(d.pitch+(y-d.lastY)*0.012, 0.06, 1.35)
		d.lastX, d.lastY = x, y
	} else {
		d.yaw -= 0.0035 // a slow idle orbit, so the scene reads as three-dimensional
	}
	// Scrolling up zooms in, as a document zooms under the wheel.
	if _, dy := in.Wheel(); dy != 0 {
		d.dist = clampf(d.dist+dy*0.25, 4, 13)
	}
	for i := range d.objs {
		o := &d.objs[i]
		o.angle = o.angle.add(o.spin.mul(step))
	}
}

// Draw renders at the element's device resolution and shows the bitmap 1:1,
// so the elements are as sharp as the display.
func (d *scene3D) Draw(p *ui.Painter, r ui.Rect, alpha float64) {
	if w, h := d.size(p, r); w != d.w || h != d.h {
		d.ensure(w, h)
	}
	d.render()
	d.bitmap.Update(d.buf)
	p.Image(d.bitmap, r, ui.FillBox)
}

// cvert is a vertex in camera space along with the world data the shader
// interpolates: position and normal.
type cvert struct {
	c  vec3 // x right, y up, z forward from the eye
	wp vec3
	wn vec3
}

func (d *scene3D) render() {
	clear(d.buf.Pix) // transparent between the objects
	for i := range d.depth {
		d.depth[i] = 3.4e38
	}

	// Camera basis.
	sp, cp := sincos(d.pitch)
	sy, cy := sincos(d.yaw)
	eye := d.target.add(v3(cp*sy, sp, cp*cy).mul(d.dist))
	forward := v3(-cp*sy, -sp, -cp*cy)
	right := forward.cross(v3(0, 1, 0)).unit()
	up := right.cross(forward).unit()
	cx, cyy := float32(d.w)/2, float32(d.h)/2
	k := cyy / float32(math.Tan(sceneFOV*math.Pi/360))

	for oi := range d.objs {
		o := &d.objs[oi]
		rot := euler(o.angle)
		base := o.pos
		if o.bob != 0 {
			base.y += o.bob * float32(math.Sin(float64(d.t*2+o.phase)))
		}
		for t := 0; t+2 < len(o.mesh.idx); t += 3 {
			a := d.cameraVert(o, rot, base, o.mesh.idx[t], eye, right, up, forward)
			b := d.cameraVert(o, rot, base, o.mesh.idx[t+1], eye, right, up, forward)
			c := d.cameraVert(o, rot, base, o.mesh.idx[t+2], eye, right, up, forward)
			d.triangle(o.mesh, a, b, c, eye, cx, cyy, k)
		}
	}
}

func (d *scene3D) cameraVert(o *sceneObject, rot mat3, base vec3, i int, eye, right, up, forward vec3) cvert {
	wp := rot.rotate(o.mesh.pos[i]).add(base)
	rel := wp.sub(eye)
	return cvert{
		c:  v3(rel.dot(right), rel.dot(up), rel.dot(forward)),
		wp: wp,
		wn: rot.rotate(o.mesh.nrm[i]),
	}
}

// triangle culls back faces, clips at the near plane and rasterizes the result.
func (d *scene3D) triangle(m *sceneMesh, a, b, c cvert, eye vec3, cx, cy, k float32) {
	// The camera basis (right, up, forward) is left-handed, so the
	// coordinate cross product is the negative of the true normal: a face
	// is front-facing when it points the same way as the vector to the eye.
	g := b.c.sub(a.c).cross(c.c.sub(a.c))
	if g.dot(a.c) <= 0 {
		return
	}
	var in [3]cvert
	in[0], in[1], in[2] = a, b, c
	var out [4]cvert
	n := 0
	for i := 0; i < 3; i++ {
		j := (i + 1) % 3
		vi, vj := in[i], in[j]
		ii, ij := vi.c.z >= sceneNear, vj.c.z >= sceneNear
		if ii {
			out[n] = vi
			n++
		}
		if ii != ij {
			t := (sceneNear - vi.c.z) / (vj.c.z - vi.c.z)
			out[n] = cvert{
				c:  vi.c.add(vj.c.sub(vi.c).mul(t)),
				wp: vi.wp.add(vj.wp.sub(vi.wp).mul(t)),
				wn: vi.wn.add(vj.wn.sub(vi.wn).mul(t)),
			}
			n++
		}
	}
	for i := 1; i+1 < n; i++ {
		d.raster(m, out[0], out[i], out[i+1], eye, cx, cy, k)
	}
}

// raster fills the triangle's pixels, testing each against the depth buffer.
func (d *scene3D) raster(m *sceneMesh, a, b, c cvert, eye vec3, cx, cy, k float32) {
	ax, ay := cx+a.c.x*k/a.c.z, cy-a.c.y*k/a.c.z
	bx, by := cx+b.c.x*k/b.c.z, cy-b.c.y*k/b.c.z
	dx, dy := cx+c.c.x*k/c.c.z, cy-c.c.y*k/c.c.z
	area := (bx-ax)*(dy-ay) - (dx-ax)*(by-ay)
	if area == 0 {
		return
	}
	minX := int(math.Floor(float64(max(0, min(ax, min(bx, dx))))))
	maxX := int(math.Floor(float64(min(float32(d.w-1), max(ax, max(bx, dx))))))
	minY := int(math.Floor(float64(max(0, min(ay, min(by, dy))))))
	maxY := int(math.Floor(float64(min(float32(d.h-1), max(ay, max(by, dy))))))
	if minX > maxX || minY > maxY {
		return
	}
	invArea := 1 / area
	for py := minY; py <= maxY; py++ {
		fy := float32(py) + 0.5
		row := d.depth[py*d.w : py*d.w+d.w]
		for px := minX; px <= maxX; px++ {
			fx := float32(px) + 0.5
			wa := ((bx-fx)*(dy-fy) - (dx-fx)*(by-fy)) * invArea
			wb := ((dx-fx)*(ay-fy) - (ax-fx)*(dy-fy)) * invArea
			wc := 1 - wa - wb
			if wa < 0 || wb < 0 || wc < 0 {
				continue
			}
			pa, pb, pc := wa/a.c.z, wb/b.c.z, wc/c.c.z
			denom := pa + pb + pc
			depth := 1 / denom
			if depth >= row[px] {
				continue
			}
			row[px] = depth
			wp := a.wp.mul(pa).add(b.wp.mul(pb)).add(c.wp.mul(pc)).mul(1 / denom)
			wn := a.wn.mul(pa).add(b.wn.mul(pb)).add(c.wn.mul(pc))
			d.shadePixel(py*d.w+px, m, wp, wn, eye)
		}
	}
}

func (d *scene3D) shadePixel(idx int, m *sceneMesh, wp, wn, eye vec3) {
	n := wn.unit()
	v := eye.sub(wp).unit()
	out := [3]float32{0.16, 0.18, 0.22}
	for i := range out {
		out[i] *= m.albedo[i]
	}
	for _, l := range sceneLights {
		ndl := n.dot(l.dir)
		if ndl <= 0 {
			continue
		}
		spec := float32(0)
		if m.spec > 0 {
			h := l.dir.add(v).unit()
			if s := n.dot(h); s > 0 {
				spec = powf(s, m.shine) * m.spec
			}
		}
		for i := range out {
			out[i] += m.albedo[i]*ndl*l.color[i] + spec*l.color[i]
		}
	}
	pix := d.buf.Pix[idx*4:]
	for i := range out {
		pix[i] = byte(clampf(out[i], 0, 1)*255 + 0.5)
	}
	pix[3] = 255
}

// powf raises x to a small non-negative integer power by squaring, cheaper
// than math.Pow in the per-pixel specular term.
func powf(x float32, n int) float32 {
	result := float32(1)
	for n > 0 {
		if n&1 == 1 {
			result *= x
		}
		x *= x
		n >>= 1
	}
	return result
}

func clampf(v, lo, hi float32) float32 {
	return min(max(v, lo), hi)
}

// sincos returns the sine and cosine of x, in float32.
func sincos(x float32) (float32, float32) {
	s, c := math.Sincos(float64(x))
	return float32(s), float32(c)
}

// example3D is the gallery page: a retained software renderer run by the MyGo
// game loop, driven by fixed steps and the plugin's input. It paints no
// background of its own, so the elements sit on the page behind them.
func example3D(c *ui.Context, s *demoState) {
	t := c.Theme()
	d := &s.scene3d
	d.init()
	components.Text(c, tr(c, "gallery.scene3d_hint")).TextColor(t.TextMuted).Wrap()
	game.View(c, d).Height(360).FillWidth().Label(tr(c, "gallery.scene3d_view"))
}
