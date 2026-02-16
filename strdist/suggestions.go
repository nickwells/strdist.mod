package strdist

import (
	"sort"

	"github.com/nickwells/english.mod/english"
)

// SuggestionString returns a string suggesting the supplied values or the
// empty string if there are no values.
func SuggestionString(vals []string) string {
	if len(vals) == 0 {
		return ""
	}

	sort.Strings(vals)

	return ", did you mean " + english.JoinQuoted(vals, ", ", " or ") + "?"
}

// SuggestedVals returns a slice of suggested alternative values for
// the given value
func SuggestedVals(val string, alts []string) []string {
	const alternativeCount = 3

	finder := DefaultFinders[CaseBlindAlgoNameCosine]

	return finder.FindNStrLike(alternativeCount, val, alts...)
}
