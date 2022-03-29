package strdist

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// DfltScaledLevFinder is a Finder with some default values suitable
// for a Scaled Levenshtein algorithm already set.
var DfltScaledLevFinder *Finder

// CaseBlindScaledLevFinder is a Finder with some default values
// suitable for a Scaled Levenshtein algorithm already set. CaseMod is set to
// ForceToLower.
var CaseBlindScaledLevFinder *Finder

func init() {
	var err error
	DfltScaledLevFinder, err =
		NewScaledLevFinder(DfltMinStrLen, DfltScaledLevThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default ScaledLevFinder: " + err.Error())
	}
	CaseBlindScaledLevFinder, err =
		NewScaledLevFinder(DfltMinStrLen, DfltScaledLevThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind ScaledLevFinder: " +
			err.Error())
	}
}

// DfltScaledLevThreshold is a default value for deciding whether a distance
// between two strings is sufficiently small for them to be considered
// similar
const DfltScaledLevThreshold = 0.33

// ScaledLevAlgo encapsulates the details needed to provide the ScaledLev
// distance.
type ScaledLevAlgo struct {
	s string
}

// NewScaledLevFinder returns a new Finder having a ScaledLev algo and an
// error which will be non-nil if the parameters are invalid - see NewFinder
// for details.
func NewScaledLevFinder(minStrLen int, threshold float64, cm CaseMod) (
	*Finder, error,
) {
	return NewFinder(minStrLen, threshold, cm,
		&ScaledLevAlgo{})
}

// Prep for a ScaledLevAlgo will pre-calculate the lower-case equivalent for
// the target string if the caseMod is set to ForceToLower
func (a *ScaledLevAlgo) Prep(s string, cm CaseMod) {
	if cm == ForceToLower {
		a.s = strings.ToLower(s)
		return
	}
	a.s = s
}

// Dist for a ScaledLevAlgo will calculate the ScaledLev distance between
// the two strings
func (a *ScaledLevAlgo) Dist(_, s string, cm CaseMod) float64 {
	if cm == ForceToLower {
		return ScaledLevDistance(a.s, strings.ToLower(s))
	}

	return ScaledLevDistance(a.s, s)
}

// ScaledLevDistance calculates the Scaled Levenshtein distance between
// strings a and b. This is the Levenshtein distance divided by the max of
// the lengths of the two strings. Two zero-length strings are taken as
// identical (with a zero distance between them)
func ScaledLevDistance(a, b string) float64 {
	aLen := utf8.RuneCountInString(a)
	bLen := utf8.RuneCountInString(b)

	if aLen == 0 && bLen == 0 {
		return 0.0
	}

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

	return float64(d[aLen][bLen]) / math.Max(float64(aLen), float64(bLen))
}
