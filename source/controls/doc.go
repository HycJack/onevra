// Package controls provides oneo's common controls, built on MyGo's
// unstyled bases rather than on MyGo's styled widgets.
//
// Each control is a base (which owns the pointer, the keyboard, the focus
// and the accessibility semantics) plus oneo's look: colors always come
// from the frame theme (see the theme package), sizes step with the theme,
// interaction states are read from the base every frame, and Draw and
// DrawOver do the custom painting. Disabled controls use the style's
// opacity treatment. Shell installs the theme before content builds;
// controls must not invent their own colors.
//
// The package stays product-neutral and minimal: it covers the everyday
// controls, not MyGo's full catalog. For anything else, compose MyGo's
// styled widgets or its bases directly.
package controls
