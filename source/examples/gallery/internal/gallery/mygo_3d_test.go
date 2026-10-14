package gallery

import (
	"bytes"
	"testing"
)

func newTestScene() *scene3D {
	d := &scene3D{}
	d.init()
	return d
}

func TestScene3DRenders(t *testing.T) {
	d := newTestScene()
	d.render()

	opaque, transparent := 0, 0
	colors := map[uint32]bool{}
	for i := range d.depth {
		pix := d.buf.Pix[i*4:]
		colors[uint32(pix[0])<<16|uint32(pix[1])<<8|uint32(pix[2])] = true
		switch pix[3] {
		case 0:
			transparent++
		case 255:
			opaque++
		}
	}
	// The background is transparent, so the elements composite over the page.
	if transparent == 0 || opaque == 0 {
		t.Fatalf("not composited transparently: %d opaque, %d transparent", opaque, transparent)
	}
	if opaque < len(d.depth)/20 || opaque > len(d.depth)*3/5 {
		t.Fatalf("scene covers %d of %d pixels", opaque, len(d.depth))
	}
	if len(colors) < 300 {
		t.Fatalf("scene has %d colors, want shading", len(colors))
	}

	// The cube and sphere sit in front of the torus and occlude part of it.
	far := newTestScene()
	far.objs = far.objs[2:]
	far.render()
	nearer := 0
	for i := range d.depth {
		if d.depth[i] < far.depth[i]-1 {
			nearer++
		}
	}
	if nearer < 200 {
		t.Fatalf("near objects occlude only %d pixels", nearer)
	}

	// A frame is a pure function of the retained state.
	again := newTestScene()
	again.render()
	if !bytes.Equal(d.buf.Pix, again.buf.Pix) {
		t.Fatal("render is not deterministic")
	}
}
