package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// controlStepper walks a checkout flow: index is the current step, and on
// fails its payment, which holds the flow there until it passes.
func controlStepper(c *ui.Context, s *controlState, variant string) {
	steps := []components.StepperStep{
		{Label: tr(c, "account"), Description: tr(c, "stepperAccountHint")},
		{Label: tr(c, "stepperShipping"), Description: tr(c, "stepperShippingHint")},
		{Label: tr(c, "stepperPayment"), Description: tr(c, "stepperPaymentHint")},
		{Label: tr(c, "stepperReview"), Description: tr(c, "stepperReviewHint")},
	}
	if s.on {
		steps[2].Error = tr(c, "stepperCardDeclined")
	}
	options := components.StepperOptions{Selectable: true}
	switch variant {
	case "Vertical":
		options.Orientation = components.Vertical
	case "Compact":
		options.Compact = true
	}
	stepper := components.Stepper(c, steps, s.index, options)
	if variant == "Vertical" {
		stepper.Width(260)
	}
	if i, ok := stepper.Chosen(); ok {
		s.index = i
	}
	ui.Row(c).Gap(8).Wrap().AlignItems(ui.Center).Children(func() {
		if components.Button(c, tr(c, "stepperBack"), components.ButtonSecondary).Disabled(s.index == 0).Clicked() {
			s.index--
		}
		if components.Button(c, tr(c, "stepperNext"), components.ButtonPrimary).Disabled(s.index == len(steps) || s.on && s.index >= 2).Clicked() {
			s.index++
		}
		if components.Button(c, tr(c, "timelineRestart"), components.ButtonGhost).Disabled(s.index == 0).Clicked() {
			s.index = 0
		}
		components.Switch(c, &s.on, tr(c, "stepperPaymentFailed"))
	})
}
