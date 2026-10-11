//go:build !darwin

package gallery

// Other platforms retain visible native review until hidden frame delivery has
// been verified against their window systems.
func nativeBackground() bool    { return false }
func startNativeFrames() func() { return func() {} }
