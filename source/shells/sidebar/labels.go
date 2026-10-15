package sidebar

// Labels supplies application-localized names for Shell's own controls.
// Empty fields take Onevra's catalog in the window's language. It does not affect native OS menus.
type Labels struct {
	LeftToggle, LeftResize             string
	Minimize, Maximize, Restore, Close string
	CloseToast                         string
}

// SetLabels changes control names without recreating the window or its content.
// Call on the UI thread. Sidebar geometry, page, and focus identities are retained.
func (s *Shell) SetLabels(labels Labels) {
	s.options.LeftToggleLabel = labels.LeftToggle
	s.options.LeftResizeLabel = labels.LeftResize
	s.options.MinimizeLabel = labels.Minimize
	s.options.MaximizeLabel = labels.Maximize
	s.options.RestoreLabel = labels.Restore
	s.options.CloseLabel = labels.Close
	s.options.CloseToastLabel = labels.CloseToast
	s.refresh()
}
