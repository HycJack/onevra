package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// A small ICNS containing one PNG-encoded 16px representation; no Apple tools
// are needed to test copying and container validation.
func testICNS(t *testing.T) []byte {
	t.Helper()
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, image.NewNRGBA(image.Rect(0, 0, 16, 16))); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 16+pngData.Len())
	copy(data, "icns")
	binary.BigEndian.PutUint32(data[4:8], uint32(len(data)))
	copy(data[8:12], "icp4")
	binary.BigEndian.PutUint32(data[12:16], uint32(8+pngData.Len()))
	copy(data[16:], pngData.Bytes())
	return data
}

func TestICNS(t *testing.T) {
	root := t.TempDir()
	resources := filepath.Join(root, "Resources")
	if err := os.Mkdir(resources, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(root, "icon.icns")
	valid := testICNS(t)
	put(t, src, string(valid))
	c := &config{root: root}
	c.MacOS.Icon = "icon.icns"
	info, err := compileIcon(context.Background(), c, root, resources)
	if err != nil || info["CFBundleIconFile"] != "AppIcon.icns" {
		t.Fatalf("compile ICNS: %v, %v", info, err)
	}
	got, err := os.ReadFile(filepath.Join(resources, "AppIcon.icns"))
	if err != nil || !bytes.Equal(got, valid) {
		t.Fatal("ICNS was not copied unchanged")
	}

	badElement := bytes.Clone(valid)
	binary.BigEndian.PutUint32(badElement[12:16], 0)
	oversizeElement := bytes.Clone(valid)
	binary.BigEndian.PutUint32(oversizeElement[12:16], uint32(len(valid)))
	for _, bad := range [][]byte{[]byte("not an ICNS"), valid[:len(valid)-1], badElement, oversizeElement} {
		put(t, src, string(bad))
		if err := copyICNS(src, filepath.Join(root, "out.icns")); err == nil {
			t.Fatal("accepted malformed ICNS")
		}
	}
	c.MacOS.Icon = "icon.png"
	if _, err := compileIcon(context.Background(), c, root, resources); err == nil {
		t.Fatal("macOS accepted PNG")
	}
}
