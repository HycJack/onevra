//go:build !js

package main

import (
	"github.com/egoist/mygo"
	"github.com/oneo-me/onevra/source/examples/gallery/internal/gallery"
)

// installBridge is only for browsers: a desktop preview selects its page at
// startup. The parameters keep both builds' call sites the same.
func installBridge(*gallery.Preview, *mygo.Window) {}
