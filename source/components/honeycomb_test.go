package components

import (
	"fmt"
	"image"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func honeycombItems(n int) []HoneycombItem {
	bitmap := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	items := make([]HoneycombItem, n)
	for i := range items {
		id := fmt.Sprintf("p%02d", i)
		items[i] = HoneycombItem{ID: id, Label: id, Image: bitmap}
	}
	return items
}

func honeycombTestMetrics() honeycombMetrics {
	return honeycombMetricsFor(theme.Palette(false))
}

func honeycombTester(items []HoneycombItem, selected *string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(false))
		// The column lays out like an item pane, the honeycomb growing to fill it.
		ui.Column(c).Fill().Label("Honeycomb").AlignItems(ui.Center).Children(func() {
			if id := Honeycomb(c, HoneycombOptions{Items: items, Selected: *selected, SelectedLabel: "Item"}); id != "" {
				*selected = id
			}
		})
	}, 600, 300)
}

func settleHoneycomb(tt *ui.Tester) {
	time.Sleep(honeycombDuration + 50*time.Millisecond)
	tt.Frame()
}

func middle(r ui.Rect) (float32, float32) { return r.X + r.W/2, r.Y + r.H/2 }

func TestHoneycombLaysOutRowsAroundTheSelection(t *testing.T) {
	items := honeycombItems(12)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	tt.Frame()
	m := honeycombTestMetrics()
	first, ok := tt.Find("Item")
	if !ok || first.W != m.large {
		t.Fatalf("selected item not large: %+v", first)
	}
	cx, cy := middle(first)
	right, _ := tt.Find("p01")
	if x, y := middle(right); math.Abs(float64(x-cx-m.magnified(m.spacing))) > 1 || y != cy || right.W >= first.W {
		t.Fatalf("neighbor not a smaller item beside the selected one: %+v %+v", right, first)
	}
	// Items fill rows above and below, none smaller than the smallest item.
	above, below := false, false
	for _, item := range items[1:] {
		r, _ := tt.Find(item.ID)
		if r.W < m.tiny-.01 {
			t.Fatalf("%s smaller than the smallest item: %+v", item.ID, r)
		}
		_, y := middle(r)
		above = above || y < cy-m.spacing/2
		below = below || y > cy+m.spacing/2
	}
	if !above || !below {
		t.Fatalf("items not in rows around the selection (above %v, below %v)", above, below)
	}
	// Clear of the edge, the further from the center, the smaller, and gradually
	// so; nearer the edge, items shrink further. An item on the ball is an
	// ellipse, so its larger side is the size it would have facing the viewer.
	view, _ := tt.Find("Honeycomb")
	sizeOf := func(r ui.Rect) float32 { return max(r.W, r.H) }
	clear := func(r ui.Rect) bool {
		x, y := middle(r)
		return r.X >= view.X && r.Y >= view.Y && r.X+r.W <= view.X+view.W && r.Y+r.H <= view.Y+view.H &&
			x-view.X > m.edge && view.X+view.W-x > m.edge && y-view.Y > m.edge && view.Y+view.H-y > m.edge
	}
	edged := false
	for _, a := range items[1:] {
		ra, _ := tt.Find(a.ID)
		if !clear(ra) {
			edged = edged || sizeOf(ra) < m.small*.8
			continue
		}
		for _, b := range items[1:] {
			rb, _ := tt.Find(b.ID)
			ax, ay := middle(ra)
			bx, by := middle(rb)
			if math.Hypot(float64(ax-cx), float64(ay-cy)) < math.Hypot(float64(bx-cx), float64(by-cy))-1 && sizeOf(ra) < sizeOf(rb) {
				t.Fatalf("%s nearer the center but smaller than %s: %+v %+v", a.ID, b.ID, ra, rb)
			}
		}
	}
	if r, _ := tt.Find("p07"); r.W <= m.tiny || r.W >= right.W {
		t.Fatalf("second ring not between the first and the far size: %+v %+v", r, right)
	}
	if !edged {
		t.Fatal("no item shrunk near the edge")
	}
}

func TestHoneycombCentersTheClickedItem(t *testing.T) {
	items := honeycombItems(9)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	m := honeycombTestMetrics()
	tt.Frame()
	center, _ := tt.Find("Item")
	cx, cy := middle(center)
	target, _ := tt.Find("p04")
	last := time.Now()
	if err := tt.Click("p04"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if selected != "p04" {
		t.Fatalf("clicked item not chosen: %q", selected)
	}
	// The honeycomb pans over, growing the chosen item on the way.
	moving := false
	// The tester runs frames back to back while the pan goes on, so a test
	// starved of time, as when the race detector runs every package at once, may
	// finish the pan before it looks.
	gap := time.Since(last)
	last = time.Now()
	for deadline := time.Now().Add(2 * honeycombDuration); time.Now().Before(deadline); {
		if r, ok := tt.Find("Item"); ok {
			x, _ := middle(r)
			tx, _ := middle(target)
			moving = moving || (x != tx && x != cx)
		}
		time.Sleep(10 * time.Millisecond)
		tt.Frame()
		gap, last = max(gap, time.Since(last)), time.Now()
	}
	if !moving && gap > honeycombDuration/2 {
		t.Logf("frames %v apart, too far apart to see the pan", gap)
		moving = true
	}
	r, _ := tt.Find("Item")
	if x, y := middle(r); !moving || r.W != m.large || math.Abs(float64(x-cx)) > 1 || math.Abs(float64(y-cy)) > 1 {
		t.Fatalf("clicked item not centered smoothly (moving %v): %+v, center %v,%v", moving, r, cx, cy)
	}
}

func TestHoneycombSettlesOnTheItemNearestTheMiddle(t *testing.T) {
	items := honeycombItems(19)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	m := honeycombTestMetrics()
	// Without motion, the ring does not keep the harness repainting, so the
	// scroll below reaches its settle check on real time as in a window.
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	center, _ := tt.Find("Item")
	cx, cy := middle(center)
	start, _ := tt.Find("p01")
	// Dragging moves the magnified middle with the pointer, growing the item it
	// brings there, and chooses nothing yet.
	x, y := middle(start)
	reach := m.spacing * m.large / m.small
	tt.Press(x, y)
	tt.Frame()
	for i := 1; i <= 6; i++ {
		tt.Move(x-float32(i)*reach/6, y)
		tt.Frame()
	}
	if r, _ := tt.Find("p01"); func() bool { mx, my := middle(r); return math.Abs(float64(mx-cx)) > .5 || my != cy }() || r.W != m.large {
		t.Fatalf("items not following the pointer: %+v, center %v,%v", r, cx, cy)
	}
	if r, _ := tt.Find("Item"); r.W >= m.large-10 {
		t.Fatalf("item out of the middle still large: %+v", r)
	}
	time.Sleep(100 * time.Millisecond)
	tt.Release(x-reach, y)
	tt.Frame()
	if selected != "p01" {
		t.Fatalf("drag did not choose the item in the middle: %q", selected)
	}
	// Scrolling pans too, and once it stops settles on the item nearest the
	// middle, here the one below and right.
	scrolledAt := time.Now()
	tt.Scroll(x, y, 0, 40*m.large/m.small)
	tt.Frame()
	// The settle waits honeycombSettle after the scroll is over. A frame slow
	// enough to outlast that wait has rightly settled already, so only check the
	// wait when it passed.
	if time.Since(scrolledAt) < honeycombSettle && selected != "p01" {
		t.Fatalf("scrolling chose an item before it stopped: %q", selected)
	}
	time.Sleep(honeycombSettle + 20*time.Millisecond)
	tt.Frame()
	if selected != "p02" {
		t.Fatalf("scrolling did not choose the item nearest the middle: %q", selected)
	}
	settleHoneycomb(tt)
	if r, _ := tt.Find("Item"); func() bool { mx, my := middle(r); return math.Abs(float64(mx-cx)) > 1 || math.Abs(float64(my-cy)) > 1 }() || r.W != m.large {
		t.Fatalf("chosen item not centered: %+v, center %v,%v", r, cx, cy)
	}
}

func TestHoneycombSettlesInside(t *testing.T) {
	items := honeycombItems(7)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	m := honeycombTestMetrics()
	tt.Frame()
	center, _ := tt.Find("Item")
	cx, cy := middle(center)
	// Past the edge of the honeycomb, items follow the pointer less.
	tt.Press(cx, cy)
	tt.Frame()
	for i := 1; i <= 10; i++ {
		tt.Move(cx+float32(i)*40, cy)
		tt.Frame()
	}
	r, _ := tt.Find("Item")
	if x, _ := middle(r); x-cx >= 400-1 || x-cx < m.spacing {
		t.Fatalf("drag past the edge not resisted: %+v", r)
	}
	tt.Release(cx+400, cy)
	tt.Frame()
	settleHoneycomb(tt)
	// The leftmost cell is the furthest the pan goes.
	cells := honeycombCells(len(items), m.spacing)
	lo, _ := honeycombBounds(cells)
	i := slices.IndexFunc(items, func(item HoneycombItem) bool { return item.ID == selected })
	r, _ = tt.Find("Item")
	if x, _ := middle(r); cells[i].x != lo.x || math.Abs(float64(x-cx)) > 1 {
		t.Fatalf("honeycomb not settled on its leftmost item: %q at %+v, %+v", selected, cells[i], r)
	}
}

func TestHoneycombCellsFillAHoneycomb(t *testing.T) {
	m := honeycombTestMetrics()
	cells := honeycombCells(40, m.spacing)
	if cells[0] != (honeycombPoint{}) {
		t.Fatalf("first cell not in the middle: %+v", cells[0])
	}
	for i, a := range cells {
		for _, b := range cells[i+1:] {
			if d := math.Hypot(float64(a.x-b.x), float64(a.y-b.y)); d < float64(m.spacing)-.01 {
				t.Fatalf("cells overlap: %+v %+v", a, b)
			}
		}
	}
	// Round, filling each ring before the next.
	far := 0.0
	for _, c := range cells {
		far = max(far, math.Hypot(float64(c.x), float64(c.y)))
	}
	for _, c := range honeycombCells(60, m.spacing)[40:] {
		if d := math.Hypot(float64(c.x), float64(c.y)); d < far-.01 {
			t.Fatalf("cell %+v nearer than the farthest of the first, %v", c, far)
		}
	}
	if !slices.Equal(cells[:7], honeycombCells(7, m.spacing)) {
		t.Fatal("cells change with the number of items")
	}
}

func TestHoneycombTiltsItemsOnABall(t *testing.T) {
	items := honeycombItems(19)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	m := honeycombTestMetrics()
	settleHoneycomb(tt)
	center, ok := tt.Find("Item")
	if !ok {
		t.Fatal("item not found")
	}
	cx, cy := middle(center)
	// The item in the middle faces the viewer, so it stays round and full size,
	// unlike the tilted ones around it.
	if center.W != center.H || center.W != m.large {
		t.Fatalf("the middle item is not round and full size: %+v", center)
	}
	// An item beside the middle turns around the vertical axis: it stays as
	// tall, but is narrower than it is tall.
	side, _ := tt.Find("p01")
	if sx, _ := middle(side); sx <= cx || side.W >= side.H {
		t.Fatalf("the item to the right did not turn away: %+v", side)
	}
	// An item directly above turns around the horizontal axis: it stays as wide,
	// but is shorter than it is wide.
	above := ui.Rect{}
	for _, item := range items[1:] {
		r, _ := tt.Find(item.ID)
		if x, y := middle(r); math.Abs(float64(x-cx)) <= 1 && y < cy-1 {
			above = r
			break
		}
	}
	if above.W == 0 {
		t.Fatal("no item directly above the middle")
	}
	if above.H >= above.W {
		t.Fatalf("the item above did not turn away: %+v", above)
	}
}

func TestHoneycombBuildsContent(t *testing.T) {
	var sizes []float32
	content := func(c *ui.Context, size float32) {
		sizes = append(sizes, size)
		ui.Box(c).Size(size, size)
	}
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(false))
		sizes = sizes[:0]
		Honeycomb(c, HoneycombOptions{Items: []HoneycombItem{{ID: "a", Label: "A", Content: content},
			{ID: "b", Label: "B", Content: content}}, Selected: "a", SelectedLabel: "Item"})
	}, 600, 300)
	m := honeycombTestMetrics()
	tt.Frame()
	// Content fills the cell's circle itself.
	if len(sizes) != 2 || math.Abs(float64(sizes[0]-m.large)) > .01 || sizes[1] >= sizes[0] {
		t.Fatalf("content not built at each item's size: %v", sizes)
	}
	if err := tt.Click("B"); err != nil {
		t.Fatal(err)
	}
}

func TestHoneycombCentersItsContent(t *testing.T) {
	// A mark in the middle of the content finds whether the item is drawn
	// centered in its round cell.
	content := func(c *ui.Context, size float32) {
		ui.Box(c).Fill().Center().Children(func() {
			ui.Box(c).Size(size/2, size/2).Background(ui.RGB(255, 0, 0))
		})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(false))
		Honeycomb(c, HoneycombOptions{Items: []HoneycombItem{{ID: "a", Label: "A", Content: content},
			{ID: "b", Label: "B", Content: content}}, Selected: "a", SelectedLabel: "Item"})
	}, 600, 300)
	tt.Frame()
	itemLogo, ok := tt.Find("Item")
	if !ok {
		t.Fatal("item not found")
	}
	img := tt.Image()
	minX, minY, maxX, maxY, found := 0, 0, 0, 0, false
	for y := int(itemLogo.Y); y < int(itemLogo.Y+itemLogo.H); y++ {
		for x := int(itemLogo.X); x < int(itemLogo.X+itemLogo.W); x++ {
			if p := img.RGBAAt(x, y); p.R > 200 && p.G < 60 && p.B < 60 {
				if !found {
					minX, minY, maxX, maxY, found = x, y, x, y, true
				} else {
					minX, minY = min(minX, x), min(minY, y)
					maxX, maxY = max(maxX, x), max(maxY, y)
				}
			}
		}
	}
	if !found {
		t.Fatal("the content's mark was not drawn")
	}
	mx, my := float32(minX+maxX+1)/2, float32(minY+maxY+1)/2
	cx, cy := middle(itemLogo)
	if math.Abs(float64(mx-cx)) > 1 || math.Abs(float64(my-cy)) > 1 {
		t.Fatalf("content off center: mark %v,%v, item %v,%v", mx, my, cx, cy)
	}
}

func TestHoneycombClicksOnlyInsideTheCircle(t *testing.T) {
	items := honeycombItems(7)
	selected := "p00"
	tt := honeycombTester(items, &selected)
	tt.Frame()
	r, _ := tt.Find("p01")
	tt.ClickAt(r.X+2, r.Y+2)
	tt.Frame()
	if selected != "p00" {
		t.Fatalf("click outside the round item chose %q", selected)
	}
	tt.ClickAt(middle(r))
	tt.Frame()
	if selected != "p01" {
		t.Fatalf("click inside the round item chose %q", selected)
	}
}

func TestHoneycombKeyboardChoosesItems(t *testing.T) {
	for _, dark := range []bool{false, true} {
		items := honeycombItems(5)
		selected := "p00"
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(dark))
			ui.Column(c).Fill().Children(func() {
				if id := Honeycomb(c, HoneycombOptions{Items: items, Selected: selected, SelectedLabel: "Item"}); id != "" {
					selected = id
				}
			})
		}, 600, 300)
		tt.Frame()
		// Tab reaches each item in order; Enter chooses the focused one.
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyTab)
		tt.Frame()
		tt.Key(0, ui.KeyEnter)
		tt.Frame()
		if selected != "p01" {
			t.Fatalf("dark=%v: keyboard chose %q, want p01", dark, selected)
		}
	}
}

func TestHoneycombRingsInTheItemsColor(t *testing.T) {
	ring := ui.RGB(0, 200, 0)
	rings := func(items []HoneycombItem) int {
		tt := ui.NewTester(func(c *ui.Context) {
			c.SetTheme(theme.Palette(false))
			Honeycomb(c, HoneycombOptions{Items: items, Selected: "a", SelectedLabel: "Item"})
		}, 600, 300)
		settleHoneycomb(tt)
		itemLogo, ok := tt.Find("Item")
		if !ok {
			t.Fatal("item not found")
		}
		img, n := tt.Image(), 0
		for y := int(itemLogo.Y); y < int(itemLogo.Y+itemLogo.H); y++ {
			for x := int(itemLogo.X); x < int(itemLogo.X+itemLogo.W); x++ {
				// The ring is translucent, so look for its tint rather than its
				// solid color.
				if p := img.RGBAAt(x, y); p.G > 40 && int(p.G) > int(p.R)+25 && int(p.G) > int(p.B)+25 {
					n++
				}
			}
		}
		return n
	}
	items := honeycombItems(2)
	items[0].ID, items[0].Label = "a", "A"
	if n := rings(items); n != 0 {
		t.Fatalf("the accent ring drew %d pixels of the item's color", n)
	}
	items[0].Ring = ring
	if n := rings(items); n < 100 {
		t.Fatalf("the ring drew only %d pixels of the item's color", n)
	}
}

func TestHoneycombScalesWithThemeSpacing(t *testing.T) {
	measure := func(spacing float32) ui.Rect {
		tt := ui.NewTester(func(c *ui.Context) {
			th := theme.Palette(false)
			th.Spacing = spacing
			c.SetTheme(th)
			ui.Column(c).Fill().Children(func() {
				Honeycomb(c, HoneycombOptions{Items: honeycombItems(1), Selected: "p00", SelectedLabel: "Item"})
			})
		}, 600, 300)
		tt.Frame()
		r, ok := tt.Find("Item")
		if !ok {
			t.Fatal("item not found")
		}
		return r
	}
	base, wide := measure(4), measure(8)
	if wide.W != 2*base.W || wide.H != 2*base.H {
		t.Fatalf("items did not scale with theme spacing: %+v %+v", base, wide)
	}
}

func TestHoneycombWithoutSelectedItemDrawsNothing(t *testing.T) {
	drawn := false
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(false))
		Honeycomb(c, HoneycombOptions{Items: honeycombItems(3), Selected: "missing", SelectedLabel: "Item"})
		drawn = true
	}, 600, 300)
	tt.Frame()
	if !drawn {
		t.Fatal("frame not built")
	}
	if _, ok := tt.Find("Item"); ok {
		t.Fatal("an unknown selection still drew items")
	}
}
