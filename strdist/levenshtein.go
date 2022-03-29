package strdist

import (
	"strings"
	"unicode/utf8"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// DfltLevenshteinFinder is a Finder with some suitable default values
// suitable for a Levenshtein algorithm already set.
var DfltLevenshteinFinder *Finder

// CaseBlindLevenshteinFinder is a Finder with some suitable default
// values suitable for a Levenshtein algorithm already set. CaseMod is set to
// ForceToLower.
var CaseBlindLevenshteinFinder *Finder

func init() {
	var err error
	DfltLevenshteinFinder, err =
		NewLevenshteinFinder(
			DfltMinStrLen, DfltLevenshteinThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default LevenshteinFinder: " + err.Error())
	}
	CaseBlindLevenshteinFinder, err =
		NewLevenshteinFinder(
			DfltMinStrLen, DfltLevenshteinThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind LevenshteinFinder: " +
			err.Error())
	}
}

// DfltLevenshteinThreshold is a default value for deciding whether a distance
// between two strings is sufficiently small for them to be considered
// similar
const DfltLevenshteinThreshold = 5.0

// LevenshteinAlgo encapsulates the details needed to provide the Levenshtein
// distance.
type LevenshteinAlgo struct {
	s string
}

// NewLevenshteinFinder returns a new Finder having a Levenshtein algo
// and an error which will be non-nil if the parameters are invalid - see
// NewFinder for details.
func NewLevenshteinFinder(minStrLen int, threshold float64, cm CaseMod) (
	*Finder, error,
) {
	return NewFinder(minStrLen, threshold, cm,
		&LevenshteinAlgo{})
}

// Prep for a LevenshteinAlgo will pre-calculate the lower-case equivalent for
// the target string if the caseMod is set to ForceToLower
func (a *LevenshteinAlgo) Prep(s string, cm CaseMod) {
	if cm == ForceToLower {
		a.s = strings.ToLower(s)
		return
	}
	a.s = s
}

// Dist for a LevenshteinAlgo will calculate the Levenshtein distance between
// the two strings
func (a *LevenshteinAlgo) Dist(_, s string, cm CaseMod) float64 {
	if cm == ForceToLower {
		return float64(LevenshteinDistance(a.s, strings.ToLower(s)))
	}

	return float64(LevenshteinDistance(a.s, s))
}

// LevenshteinDistance calculates the Levenshtein distance between strings a
// and b
func LevenshteinDistance(a, b string) int {
	aLen := utf8.RuneCountInString(a)
	bLen := utf8.RuneCountInString(b)
	d := make([][]int, aLen+1)
	for i := range d {
		d[i] = make([]int, bLen+1)
		d[i][0] = i
	}

	for i := 1; i <= bLen; i++ {
		d[0][i] = i
	}

	for j, bRune := range b {
		for i, aRune := range a {
			var subsCost int
			if aRune != bRune {
				subsCost = 1
			}

			del := d[i][j+1] + 1
			ins := d[i+1][j] + 1
			sub := d[i][j] + subsCost

			d[i+1][j+1] = mathutil.MinOfInt(del, ins, sub)
		}
	}

	return d[aLen][bLen]
}
