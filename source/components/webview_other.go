//go:build !darwin || !arm64

package components

import "github.com/egoist/mygo"

func supported() error {
	return ErrWebViewUnsupported
}

func attach(_, _ *mygo.Window) (attachment, error) { return nil, supported() }
