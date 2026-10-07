// Package ui provides product-neutral common controls built on MyGo.
// Callers own business state, localized text and actions. Native window and
// page composition belong to the separate shell package.
//
// Common controls use MyGo's unstyled bases, which own pointer, keyboard, focus
// and accessibility semantics. Colors and sizes come from the frame theme (see
// the theme package); interaction states are read from the base every frame,
// and Draw and DrawOver provide oneo's appearance. Disabled controls use the
// style's opacity treatment. The shell package installs the theme before content
// builds; without Shell, callers can apply it through the theme package.
//
// The package covers everyday controls, not MyGo's full catalog. For anything
// else, compose MyGo's styled widgets or its bases directly. When importing both
// UI packages, callers can alias this package as oneoui and retain ui for MyGo.
package ui
