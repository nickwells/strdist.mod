package strdist

import (
	"strings"
	"unicode/utf8"
)

// DfltHammingFinder is a HammingFinder with some suitable default values
// already set.
var DfltHammingFinder *Finder

// CaseBlindHammingFinder is a HammingFinder with some suitable default
// values already set.
var CaseBlindHammingFinder *Finder

func init() {
	var err error
	DfltHammingFinder, err =
		NewHammingFinder(DfltMinStrLen, DfltHammingThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default HammingFinder: " + err.Error())
	}
	CaseBlindHammingFinder, err =
		NewHammingFinder(DfltMinStrLen, DfltHammingThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind HammingFinder: " +
			err.Error())
	}
}

// DfltHammingThreshold is a default value for deciding whether a distance
// between two strings is sufficiently small for them to be considered
// similar
const DfltHammingThreshold = 5.0

// HammingAlgo encapsulates the details needed to provide the Hamming distance.
type HammingAlgo struct {
	s string
}

// NewHammingFinder returns a new Finder having a Hamming algo and an
// error which will be non-nil if the parameters are invalid - see
// NewFinder for details.
func NewHammingFinder(minStrLen int, threshold float64, cm CaseMod) (*Finder, error) {
	return NewFinder(minStrLen, threshold, cm,
		&HammingAlgo{})
}

// Prep for a HammingAlgo will pre-calculate the lower-case equivalent for
// the target string if the caseMod is set to ForceToLower
func (a *HammingAlgo) Prep(s string, cm CaseMod) {
	if cm == ForceToLower {
		a.s = strings.ToLower(s)
		return
	}
	a.s = s
}

// Dist for a HammingAlgo will calculate the Hamming distance between the two
// strings
func (a *HammingAlgo) Dist(_, s string, cm CaseMod) float64 {
	if cm == ForceToLower {
		return HammingDistance(a.s, strings.ToLower(s))
	}

	return HammingDistance(a.s, s)
}

// HammingDistance returns the Hamming distance of the two strings. if the
// two strings are of different length then the Hamming distance is increased
// by the difference in lengths. Note that it compares runes rather than
// characters or chars
func HammingDistance(a, b string) float64 {
	var d = utf8.RuneCountInString(b) - utf8.RuneCountInString(a)
	if d < 0 {
		d *= -1
		a, b = b, a // a is longer than b so swap
	}

	var offset int
	for _, aRune := range a {
		bRune, width := utf8.DecodeRuneInString(b[offset:])
		offset += width
		if bRune != aRune {
			d++
		}
	}

	return float64(d)
}
