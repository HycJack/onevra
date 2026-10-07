//go:build !darwin && (!windows || (!amd64 && !arm64))

package components

import "github.com/egoist/mygo"

func attachWebView(_ *WebView, _ mygo.NativeComposition) (attachment, error) {
	return nil, ErrWebViewUnsupported
}
