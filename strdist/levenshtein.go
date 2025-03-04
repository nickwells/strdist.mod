package strdist

import (
	"unicode/utf8"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// LevenshteinAlgo encapsulates the details needed to provide the Levenshtein
// distance.
type LevenshteinAlgo struct{}

// Name returns the algorithm name
func (LevenshteinAlgo) Name() string {
	return AlgoNameLevenshtein
}

// Desc returns a string describing the algorithm configuration
func (LevenshteinAlgo) Desc() string {
	return ""
}

// Dist for a LevenshteinAlgo will calculate the Levenshtein distance between
// the two strings
func (LevenshteinAlgo) Dist(s1, s2 string) float64 {
	return float64(LevenshteinDistance(s1, s2))
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
