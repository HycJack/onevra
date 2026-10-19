package gallery

import "math/rand/v2"

// bubbleLabels name the sample's categories, of which the first count show.
var bubbleLabels = []string{
	"chartDesign", "chartDevelopment", "chartResearch", "chartPlanning", "chartReview",
	"chartSupport", "chartTesting", "chartOperations", "chartDocumentation", "chartMarketing",
}

const bubbleInitialCount = 7

// bubbleValues returns the sample values of an update round, the same for
// every category added later, so updates and additions are independent.
func bubbleValues(round int) []float64 {
	values := []float64{38, 82, 46, 20, 27, 14, 61, 9, 33, 52}
	if round == 0 {
		return values
	}
	r := rand.New(rand.NewPCG(uint64(round), 1))
	for i := range values {
		values[i] = float64(5 + r.IntN(95))
	}
	return values
}
