//go:build !darwin

package gesture

func nativeScrollTouch() (touch, ok bool)   { return false, false }
func currentPointer() (nativePointer, bool) { return nativePointer{}, false }
func installPinch()                         {}
