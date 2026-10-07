// Package ui provides product-neutral UI compositions and common controls built
// on MyGo. NewShell creates a native window with a replaceable page container
// and a default sidebar-and-content layout. Callers own application lifetime,
// business state, localized text and actions.
//
// Common controls use MyGo's unstyled bases, which own pointer, keyboard, focus
// and accessibility semantics. Colors and sizes come from the frame theme (see
// the theme package); interaction states are read from the base every frame,
// and Draw and DrawOver provide oneo's appearance. Disabled controls use the
// style's opacity treatment. Shell installs the theme before content builds.
//
// The package covers everyday controls, not MyGo's full catalog. For anything
// else, compose MyGo's styled widgets or its bases directly. When importing both
// UI packages, callers can alias this package as oneoui and retain ui for MyGo.
package ui
